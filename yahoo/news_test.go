package yahoo

import (
	"net/http"
	"strings"
	"testing"
)

func TestGetNews_Success(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		if !strings.Contains(req.URL.Path, "/v1/finance/search") {
			t.Errorf("unexpected path: %s", req.URL.Path)
		}
		q := req.URL.Query().Get("q")
		if q != "AAPL" {
			t.Errorf("q = %q, want %q", q, "AAPL")
		}
		quotesCount := req.URL.Query().Get("quotesCount")
		if quotesCount != "0" {
			t.Errorf("quotesCount should be 0 for news, got %q", quotesCount)
		}
		return jsonResponse(200, `{
			"quotes": [],
			"news": [
				{
					"uuid": "abc-123",
					"title": "Apple Reports Record Quarter",
					"publisher": "Reuters",
					"link": "https://example.com/article1",
					"providerPublishTime": 1700000000
				},
				{
					"uuid": "def-456",
					"title": "Apple Launches New Product",
					"publisher": "Bloomberg",
					"link": "https://example.com/article2",
					"providerPublishTime": 1700086400
				}
			],
			"count": 2
		}`), nil
	})

	news, err := client.GetNews("AAPL", 5)
	if err != nil {
		t.Fatalf("GetNews() error: %v", err)
	}

	if len(news) != 2 {
		t.Fatalf("news count = %d, want 2", len(news))
	}
	if news[0].Title != "Apple Reports Record Quarter" {
		t.Errorf("news[0].Title = %q, want %q", news[0].Title, "Apple Reports Record Quarter")
	}
	if news[0].Publisher != "Reuters" {
		t.Errorf("news[0].Publisher = %q, want %q", news[0].Publisher, "Reuters")
	}
	if news[1].UUID != "def-456" {
		t.Errorf("news[1].UUID = %q, want %q", news[1].UUID, "def-456")
	}
}

func TestGetNews_DefaultCount(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		newsCount := req.URL.Query().Get("newsCount")
		if newsCount != "5" {
			t.Errorf("default newsCount = %q, want %q", newsCount, "5")
		}
		return jsonResponse(200, `{"quotes": [], "news": [], "count": 0}`), nil
	})

	_, err := client.GetNews("AAPL", 0)
	if err != nil {
		t.Fatalf("GetNews() error: %v", err)
	}
}

func TestGetNews_CustomCount(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		newsCount := req.URL.Query().Get("newsCount")
		if newsCount != "10" {
			t.Errorf("newsCount = %q, want %q", newsCount, "10")
		}
		return jsonResponse(200, `{"quotes": [], "news": [], "count": 0}`), nil
	})

	_, err := client.GetNews("TSLA", 10)
	if err != nil {
		t.Fatalf("GetNews() error: %v", err)
	}
}

func TestGetNews_EmptyResults(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		return jsonResponse(200, `{"quotes": [], "news": [], "count": 0}`), nil
	})

	news, err := client.GetNews("UNKNOWN", 5)
	if err != nil {
		t.Fatalf("GetNews() error: %v", err)
	}
	if len(news) != 0 {
		t.Errorf("expected empty news, got %d items", len(news))
	}
}
