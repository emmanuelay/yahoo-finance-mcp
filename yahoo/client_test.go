package yahoo

import (
	"net/http"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"
)

func TestAuthenticate_Success(t *testing.T) {
	client := newUnauthClient(func(req *http.Request) (*http.Response, error) {
		switch req.URL.Host {
		case "fc.yahoo.com":
			return textResponse(404, ""), nil
		default:
			if strings.Contains(req.URL.Path, "/v1/test/getcrumb") {
				return textResponse(200, "my-crumb-123"), nil
			}
			t.Fatalf("unexpected request: %s", req.URL.String())
			return nil, nil
		}
	})

	err := client.authenticate()
	if err != nil {
		t.Fatalf("authenticate() error: %v", err)
	}
	if client.crumb != "my-crumb-123" {
		t.Errorf("crumb = %q, want %q", client.crumb, "my-crumb-123")
	}
	if !client.authed {
		t.Error("authed should be true after authenticate()")
	}
}

func TestAuthenticate_CrumbFailure(t *testing.T) {
	client := newUnauthClient(func(req *http.Request) (*http.Response, error) {
		switch req.URL.Host {
		case "fc.yahoo.com":
			return textResponse(404, ""), nil
		default:
			return textResponse(403, "Forbidden"), nil
		}
	})

	err := client.authenticate()
	if err == nil {
		t.Fatal("expected error from authenticate()")
	}
	if !strings.Contains(err.Error(), "403") {
		t.Errorf("error should mention status 403, got: %v", err)
	}
}

func TestGet_AddsCrumb(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		crumb := req.URL.Query().Get("crumb")
		if crumb != "test-crumb" {
			t.Errorf("crumb = %q, want %q", crumb, "test-crumb")
		}
		return jsonResponse(200, `{"ok":true}`), nil
	})

	_, err := client.Get("/v10/finance/quoteSummary/AAPL", nil, true)
	if err != nil {
		t.Fatalf("Get() error: %v", err)
	}
}

func TestGet_NoCrumbWhenNotNeeded(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		crumb := req.URL.Query().Get("crumb")
		if crumb != "" {
			t.Errorf("crumb should be empty, got %q", crumb)
		}
		return jsonResponse(200, `{}`), nil
	})

	_, err := client.Get("/v8/finance/chart/AAPL", url.Values{"range": {"1d"}}, false)
	if err != nil {
		t.Fatalf("Get() error: %v", err)
	}
}

func TestGet_RetryOn401(t *testing.T) {
	var callCount atomic.Int32

	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		// Handle auth endpoints during re-authentication
		if req.URL.Host == "fc.yahoo.com" {
			return textResponse(404, ""), nil
		}
		if strings.Contains(req.URL.Path, "/v1/test/getcrumb") {
			return textResponse(200, "new-crumb"), nil
		}

		n := callCount.Add(1)
		if n == 1 {
			return jsonResponse(401, `{"error":"Unauthorized"}`), nil
		}
		return jsonResponse(200, `{"data":"ok"}`), nil
	})

	body, err := client.Get("/v10/finance/quoteSummary/AAPL", nil, true)
	if err != nil {
		t.Fatalf("Get() error: %v", err)
	}
	if !strings.Contains(string(body), "ok") {
		t.Errorf("expected retry response body, got %s", string(body))
	}
	if callCount.Load() != 2 {
		t.Errorf("expected 2 calls to API endpoint, got %d", callCount.Load())
	}
}

func TestGet_NonOKStatus(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		return jsonResponse(500, `Internal Server Error`), nil
	})

	_, err := client.Get("/v10/finance/quoteSummary/AAPL", nil, false)
	if err == nil {
		t.Fatal("expected error for 500 status")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("error should mention status 500, got: %v", err)
	}
}

func TestGetJSON_Success(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		return jsonResponse(200, `{"name":"test","value":42}`), nil
	})

	var result struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}
	err := client.GetJSON("/test", nil, false, &result)
	if err != nil {
		t.Fatalf("GetJSON() error: %v", err)
	}
	if result.Name != "test" || result.Value != 42 {
		t.Errorf("GetJSON() result = %+v, want {test 42}", result)
	}
}

func TestGetJSON_InvalidJSON(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		return jsonResponse(200, `not-json`), nil
	})

	var result map[string]interface{}
	err := client.GetJSON("/test", nil, false, &result)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
	if !strings.Contains(err.Error(), "parsing JSON") {
		t.Errorf("error should mention JSON parsing, got: %v", err)
	}
}

func TestEnsureAuth_SkipsIfAlreadyAuthed(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		t.Fatal("should not make any HTTP requests when already authed")
		return nil, nil
	})

	err := client.ensureAuth()
	if err != nil {
		t.Fatalf("ensureAuth() error: %v", err)
	}
}

func TestGetCrumb_ThreadSafe(t *testing.T) {
	client := newTestClient(nil)
	crumb := client.getCrumb()
	if crumb != "test-crumb" {
		t.Errorf("getCrumb() = %q, want %q", crumb, "test-crumb")
	}
}
