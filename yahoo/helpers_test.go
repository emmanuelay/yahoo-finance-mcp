package yahoo

import (
	"io"
	"net/http"
	"strings"
)

// roundTripFunc allows using a function as an http.RoundTripper.
type roundTripFunc func(req *http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// newTestClient creates a pre-authenticated Client with a mock transport.
func newTestClient(fn roundTripFunc) *Client {
	return &Client{
		httpClient: &http.Client{Transport: fn},
		crumb:      "test-crumb",
		authed:     true,
	}
}

// newUnauthClient creates an unauthenticated Client with a mock transport.
func newUnauthClient(fn roundTripFunc) *Client {
	return &Client{
		httpClient: &http.Client{Transport: fn},
	}
}

// jsonResponse builds an *http.Response with the given status and JSON body.
func jsonResponse(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     http.Header{"Content-Type": []string{"application/json"}},
	}
}

// textResponse builds an *http.Response with the given status and plain text body.
func textResponse(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     http.Header{"Content-Type": []string{"text/plain"}},
	}
}
