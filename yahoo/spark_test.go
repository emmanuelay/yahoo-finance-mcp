package yahoo

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestGetBulkSpark_Success(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		if !strings.Contains(req.URL.Path, "/v8/finance/spark") {
			t.Errorf("unexpected path: %s", req.URL.Path)
		}
		q := req.URL.Query()
		if syms := q.Get("symbols"); syms != "AAPL,MSFT" {
			t.Errorf("symbols = %q, want %q", syms, "AAPL,MSFT")
		}
		if r := q.Get("range"); r != "5d" {
			t.Errorf("range = %q, want %q", r, "5d")
		}
		if iv := q.Get("interval"); iv != "1h" {
			t.Errorf("interval = %q, want %q", iv, "1h")
		}
		return jsonResponse(200, `{
			"AAPL": {
				"symbol": "AAPL",
				"timestamp": [1700000000, 1700003600, 1700007200],
				"close": [178.50, 179.10, 178.72],
				"chartPreviousClose": 176.57,
				"dataGranularity": 3600
			},
			"MSFT": {
				"symbol": "MSFT",
				"timestamp": [1700000000, 1700003600],
				"close": [379.00, 380.50],
				"chartPreviousClose": 378.20,
				"dataGranularity": 3600
			}
		}`), nil
	})

	resp, err := client.GetBulkSpark([]string{"AAPL", "MSFT"}, "5d", "1h")
	if err != nil {
		t.Fatalf("GetBulkSpark() error: %v", err)
	}

	if len(resp) != 2 {
		t.Fatalf("got %d results, want 2", len(resp))
	}

	aapl := resp["AAPL"]
	if aapl.Symbol != "AAPL" {
		t.Errorf("AAPL.Symbol = %q, want %q", aapl.Symbol, "AAPL")
	}
	if len(aapl.Timestamps) != 3 {
		t.Errorf("AAPL.Timestamps length = %d, want 3", len(aapl.Timestamps))
	}
	if len(aapl.Close) != 3 {
		t.Errorf("AAPL.Close length = %d, want 3", len(aapl.Close))
	}
	if aapl.ChartPreviousClose != 176.57 {
		t.Errorf("AAPL.ChartPreviousClose = %v, want 176.57", aapl.ChartPreviousClose)
	}
	if aapl.DataGranularity != 3600 {
		t.Errorf("AAPL.DataGranularity = %v, want 3600", aapl.DataGranularity)
	}

	msft := resp["MSFT"]
	if msft.Symbol != "MSFT" {
		t.Errorf("MSFT.Symbol = %q, want %q", msft.Symbol, "MSFT")
	}
	if len(msft.Close) != 2 {
		t.Errorf("MSFT.Close length = %d, want 2", len(msft.Close))
	}
}

func TestGetBulkSpark_Defaults(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		q := req.URL.Query()
		if r := q.Get("range"); r != "1mo" {
			t.Errorf("default range = %q, want %q", r, "1mo")
		}
		if iv := q.Get("interval"); iv != "1d" {
			t.Errorf("default interval = %q, want %q", iv, "1d")
		}
		return jsonResponse(200, `{
			"TSLA": {"symbol": "TSLA", "timestamp": [1700000000], "close": [250.0], "chartPreviousClose": 248.0, "dataGranularity": 86400}
		}`), nil
	})

	resp, err := client.GetBulkSpark([]string{"TSLA"}, "", "")
	if err != nil {
		t.Fatalf("GetBulkSpark() error: %v", err)
	}
	if resp["TSLA"].Symbol != "TSLA" {
		t.Errorf("Symbol = %q, want %q", resp["TSLA"].Symbol, "TSLA")
	}
}

func TestGetBulkSpark_EmptySymbols(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		t.Fatal("should not make HTTP request for empty symbols")
		return nil, nil
	})

	_, err := client.GetBulkSpark([]string{}, "1mo", "1d")
	if err == nil {
		t.Fatal("expected error for empty symbols")
	}
	if !strings.Contains(err.Error(), "at least one symbol") {
		t.Errorf("error should mention at least one symbol, got: %v", err)
	}
}

func TestGetBulkSpark_TooManySymbols(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		t.Fatal("should not make HTTP request for too many symbols")
		return nil, nil
	})

	symbols := make([]string, maxBulkSymbols+1)
	for i := range symbols {
		symbols[i] = fmt.Sprintf("SYM%d", i)
	}

	_, err := client.GetBulkSpark(symbols, "1mo", "1d")
	if err == nil {
		t.Fatal("expected error for too many symbols")
	}
	if !strings.Contains(err.Error(), "too many symbols") {
		t.Errorf("error should mention too many symbols, got: %v", err)
	}
}

func TestGetBulkSpark_HTTPError(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		return jsonResponse(500, `Internal Server Error`), nil
	})

	_, err := client.GetBulkSpark([]string{"AAPL"}, "1mo", "1d")
	if err == nil {
		t.Fatal("expected error for HTTP 500")
	}
}
