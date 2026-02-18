package yahoo

import (
	"net/http"
	"strings"
	"testing"
)

func TestSearch_Success(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		if !strings.Contains(req.URL.Path, "/v1/finance/search") {
			t.Errorf("unexpected path: %s", req.URL.Path)
		}
		q := req.URL.Query().Get("q")
		if q != "Apple" {
			t.Errorf("q = %q, want %q", q, "Apple")
		}
		return jsonResponse(200, `{
			"quotes": [
				{
					"symbol": "AAPL",
					"shortname": "Apple Inc.",
					"longname": "Apple Inc.",
					"exchange": "NMS",
					"quoteType": "EQUITY",
					"industry": "Consumer Electronics",
					"sector": "Technology",
					"score": 100000
				},
				{
					"symbol": "APLE",
					"shortname": "Apple Hospitality REIT",
					"exchange": "NYQ",
					"quoteType": "EQUITY",
					"score": 20000
				}
			],
			"news": [],
			"count": 2
		}`), nil
	})

	resp, err := client.Search("Apple", 5)
	if err != nil {
		t.Fatalf("Search() error: %v", err)
	}

	if len(resp.Quotes) != 2 {
		t.Fatalf("quotes count = %d, want 2", len(resp.Quotes))
	}
	if resp.Quotes[0].Symbol != "AAPL" {
		t.Errorf("first quote symbol = %q, want %q", resp.Quotes[0].Symbol, "AAPL")
	}
	if resp.Quotes[0].Sector != "Technology" {
		t.Errorf("first quote sector = %q, want %q", resp.Quotes[0].Sector, "Technology")
	}
	if resp.Count != 2 {
		t.Errorf("count = %d, want 2", resp.Count)
	}
}

func TestSearch_DefaultLimit(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		quotesCount := req.URL.Query().Get("quotesCount")
		if quotesCount != "10" {
			t.Errorf("default quotesCount = %q, want %q", quotesCount, "10")
		}
		return jsonResponse(200, `{"quotes": [], "news": [], "count": 0}`), nil
	})

	_, err := client.Search("test", 0)
	if err != nil {
		t.Fatalf("Search() error: %v", err)
	}
}

func TestSearch_CustomLimit(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		quotesCount := req.URL.Query().Get("quotesCount")
		if quotesCount != "3" {
			t.Errorf("quotesCount = %q, want %q", quotesCount, "3")
		}
		return jsonResponse(200, `{"quotes": [], "news": [], "count": 0}`), nil
	})

	_, err := client.Search("test", 3)
	if err != nil {
		t.Fatalf("Search() error: %v", err)
	}
}

func TestSearch_NoCrumb(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		crumb := req.URL.Query().Get("crumb")
		if crumb != "" {
			t.Errorf("search should not include crumb, got %q", crumb)
		}
		return jsonResponse(200, `{"quotes": [], "news": [], "count": 0}`), nil
	})

	client.Search("test", 5)
}
