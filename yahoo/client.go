package yahoo

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sync"
)

const baseURL = "https://query2.finance.yahoo.com"

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
			Jar: jar,
		},
	}
}

// authenticate performs the cookie/crumb flow.
func (c *Client) authenticate() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Step 1: GET fc.yahoo.com to get cookies (expect 404)
	req, err := http.NewRequest("GET", "https://fc.yahoo.com", nil)
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
	req, err = http.NewRequest("GET", baseURL+"/v1/test/getcrumb", nil)
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
func (c *Client) ensureAuth() error {
	c.mu.RLock()
	authed := c.authed
	c.mu.RUnlock()

	if !authed {
		return c.authenticate()
	}
	return nil
}

// getCrumb returns the current crumb in a thread-safe way.
func (c *Client) getCrumb() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.crumb
}

// Get performs an authenticated GET request to a Yahoo Finance API endpoint.
// If needsCrumb is true, the crumb parameter is appended.
func (c *Client) Get(path string, params url.Values, needsCrumb bool) ([]byte, error) {
	if needsCrumb {
		if err := c.ensureAuth(); err != nil {
			return nil, err
		}
	}

	fullURL := baseURL + path
	if needsCrumb {
		if params == nil {
			params = url.Values{}
		}
		params.Set("crumb", c.getCrumb())
	}
	if len(params) > 0 {
		fullURL += "?" + params.Encode()
	}

	body, statusCode, err := c.doGet(fullURL)
	if err != nil {
		return nil, err
	}

	// Retry on 401/403: re-authenticate once and retry
	if statusCode == http.StatusUnauthorized || statusCode == http.StatusForbidden {
		if err := c.authenticate(); err != nil {
			return nil, fmt.Errorf("re-authentication failed: %w", err)
		}
		if needsCrumb {
			params.Set("crumb", c.getCrumb())
			fullURL = baseURL + path + "?" + params.Encode()
		}
		body, statusCode, err = c.doGet(fullURL)
		if err != nil {
			return nil, err
		}
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", statusCode, string(body))
	}

	return body, nil
}

func (c *Client) doGet(url string) ([]byte, int, error) {
	req, err := http.NewRequest("GET", url, nil)
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

// GetJSON performs a GET and unmarshals the JSON response into v.
func (c *Client) GetJSON(path string, params url.Values, needsCrumb bool, v interface{}) error {
	body, err := c.Get(path, params, needsCrumb)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, v); err != nil {
		return fmt.Errorf("parsing JSON response: %w", err)
	}
	return nil
}

// GetAbsoluteJSON fetches an absolute URL with crumb auth and unmarshals the JSON response into v.
func (c *Client) GetAbsoluteJSON(absoluteURL string, params url.Values, v any) error {
	if err := c.ensureAuth(); err != nil {
		return err
	}

	if params == nil {
		params = url.Values{}
	}
	params.Set("crumb", c.getCrumb())
	fullURL := absoluteURL + "?" + params.Encode()

	body, statusCode, err := c.doGet(fullURL)
	if err != nil {
		return err
	}

	// Retry on 401/403
	if statusCode == http.StatusUnauthorized || statusCode == http.StatusForbidden {
		if err := c.authenticate(); err != nil {
			return fmt.Errorf("re-authentication failed: %w", err)
		}
		params.Set("crumb", c.getCrumb())
		fullURL = absoluteURL + "?" + params.Encode()
		body, statusCode, err = c.doGet(fullURL)
		if err != nil {
			return err
		}
	}

	if statusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status %d: %s", statusCode, string(body))
	}
	if err := json.Unmarshal(body, v); err != nil {
		return fmt.Errorf("parsing JSON response: %w", err)
	}
	return nil
}

func randomUA() string {
	return userAgents[rand.Intn(len(userAgents))]
}
