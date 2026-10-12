package yahoo

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sync"
	"time"
)

const baseURL = "https://query2.finance.yahoo.com"

// defaultTimeout bounds every HTTP request as a safety net on top of the
// caller's context, so a stalled upstream never hangs a tool call forever.
const defaultTimeout = 30 * time.Second

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:121.0) Gecko/20100101 Firefox/121.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.2 Safari/605.1.15",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0",
}

// Client is a Yahoo Finance API client with cookie/crumb authentication.
type Client struct {
	httpClient *http.Client
	crumb      string
	mu         sync.RWMutex
	authed     bool
}

// NewClient creates a new Yahoo Finance client.
func NewClient() *Client {
	jar, _ := cookiejar.New(nil)
	return &Client{
		httpClient: &http.Client{
			Jar:     jar,
			Timeout: defaultTimeout,
		},
	}
}

// authenticate performs the cookie/crumb flow.
func (c *Client) authenticate(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Step 1: GET fc.yahoo.com to get cookies (expect 404)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://fc.yahoo.com", nil)
	if err != nil {
		return fmt.Errorf("creating cookie request: %w", err)
	}
	req.Header.Set("User-Agent", randomUA())
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("fetching cookies: %w", err)
	}
	resp.Body.Close()

	// Step 2: GET crumb using cookies
	req, err = http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/v1/test/getcrumb", nil)
	if err != nil {
		return fmt.Errorf("creating crumb request: %w", err)
	}
	req.Header.Set("User-Agent", randomUA())
	resp, err = c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("fetching crumb: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("crumb request failed with status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading crumb: %w", err)
	}

	c.crumb = string(body)
	c.authed = true
	return nil
}

// ensureAuth performs lazy authentication on first call.
func (c *Client) ensureAuth(ctx context.Context) error {
	c.mu.RLock()
	authed := c.authed
	c.mu.RUnlock()

	if !authed {
		return c.authenticate(ctx)
	}
	return nil
}

// getCrumb returns the current crumb in a thread-safe way.
func (c *Client) getCrumb() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.crumb
}

// Get performs a GET request to a path on the default Yahoo Finance host.
// If needsCrumb is true, the client authenticates first and appends the crumb.
func (c *Client) Get(ctx context.Context, path string, params url.Values, needsCrumb bool) ([]byte, error) {
	return c.GetAbsolute(ctx, baseURL+path, params, needsCrumb)
}

// GetJSON performs a GET on the default host and unmarshals the JSON response into v.
func (c *Client) GetJSON(ctx context.Context, path string, params url.Values, needsCrumb bool, v any) error {
	return c.GetAbsoluteJSON(ctx, baseURL+path, params, needsCrumb, v)
}

// GetAbsoluteJSON performs a GET on an absolute URL and unmarshals the JSON response into v.
func (c *Client) GetAbsoluteJSON(ctx context.Context, absoluteURL string, params url.Values, needsCrumb bool, v any) error {
	body, err := c.GetAbsolute(ctx, absoluteURL, params, needsCrumb)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, v); err != nil {
		return fmt.Errorf("parsing JSON response: %w", err)
	}
	return nil
}

// GetAbsolute performs a GET request to an absolute URL. It is the single
// request path shared by every endpoint: it attaches the crumb when needed,
// re-authenticates once on 401/403, and rejects any non-200 response.
func (c *Client) GetAbsolute(ctx context.Context, absoluteURL string, params url.Values, needsCrumb bool) ([]byte, error) {
	if needsCrumb {
		if err := c.ensureAuth(ctx); err != nil {
			return nil, err
		}
	}

	body, statusCode, err := c.doGet(ctx, buildURL(absoluteURL, params, c.crumbIf(needsCrumb)))
	if err != nil {
		return nil, err
	}

	// Retry on 401/403: re-authenticate once and retry with the fresh crumb.
	if statusCode == http.StatusUnauthorized || statusCode == http.StatusForbidden {
		if err := c.authenticate(ctx); err != nil {
			return nil, fmt.Errorf("re-authentication failed: %w", err)
		}
		body, statusCode, err = c.doGet(ctx, buildURL(absoluteURL, params, c.crumbIf(needsCrumb)))
		if err != nil {
			return nil, err
		}
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", statusCode, string(body))
	}

	return body, nil
}

// crumbIf returns the current crumb when needed, or "" otherwise.
func (c *Client) crumbIf(needsCrumb bool) string {
	if !needsCrumb {
		return ""
	}
	return c.getCrumb()
}

// buildURL encodes params onto absoluteURL, adding the crumb when non-empty.
// It copies params so the caller's map is never mutated.
func buildURL(absoluteURL string, params url.Values, crumb string) string {
	q := make(url.Values, len(params)+1)
	for k, v := range params {
		q[k] = v
	}
	if crumb != "" {
		q.Set("crumb", crumb)
	}
	if len(q) == 0 {
		return absoluteURL
	}
	return absoluteURL + "?" + q.Encode()
}

func (c *Client) doGet(ctx context.Context, url string) ([]byte, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("User-Agent", randomUA())
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("reading response: %w", err)
	}

	return body, resp.StatusCode, nil
}

func randomUA() string {
	return userAgents[rand.Intn(len(userAgents))]
}
