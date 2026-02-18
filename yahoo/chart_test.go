package yahoo

import (
	"net/http"
	"strings"
	"testing"
)

func TestGetChart_Success(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		if !strings.Contains(req.URL.Path, "/v8/finance/chart/AAPL") {
			t.Errorf("unexpected path: %s", req.URL.Path)
		}
		return jsonResponse(200, `{
			"chart": {
				"result": [{
					"meta": {
						"currency": "USD",
						"symbol": "AAPL",
						"exchangeName": "NMS",
						"regularMarketPrice": 178.72,
						"previousClose": 176.57,
						"dataGranularity": "1d",
						"range": "1mo"
					},
					"timestamp": [1700000000, 1700086400, 1700172800],
					"indicators": {
						"quote": [{
							"open": [175.0, 176.5, 177.0],
							"high": [177.0, 178.0, 179.0],
							"low": [174.0, 175.5, 176.0],
							"close": [176.5, 177.0, 178.72],
							"volume": [50000000, 45000000, 54000000]
						}],
						"adjclose": [{
							"adjclose": [176.5, 177.0, 178.72]
						}]
					}
				}]
			}
		}`), nil
	})

	result, err := client.GetChart("AAPL", "1mo", "1d")
	if err != nil {
		t.Fatalf("GetChart() error: %v", err)
	}

	if result.Meta.Symbol != "AAPL" {
		t.Errorf("meta.Symbol = %q, want %q", result.Meta.Symbol, "AAPL")
	}
	if result.Meta.Currency != "USD" {
		t.Errorf("meta.Currency = %q, want %q", result.Meta.Currency, "USD")
	}
	if len(result.Timestamps) != 3 {
		t.Errorf("timestamps count = %d, want 3", len(result.Timestamps))
	}
	if len(result.Indicators.Quote) != 1 {
		t.Fatalf("quote indicators count = %d, want 1", len(result.Indicators.Quote))
	}
	if len(result.Indicators.Quote[0].Close) != 3 {
		t.Errorf("close prices count = %d, want 3", len(result.Indicators.Quote[0].Close))
	}
}

func TestGetChart_Defaults(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		rangeParam := req.URL.Query().Get("range")
		interval := req.URL.Query().Get("interval")
		if rangeParam != "1mo" {
			t.Errorf("default range = %q, want %q", rangeParam, "1mo")
		}
		if interval != "1d" {
			t.Errorf("default interval = %q, want %q", interval, "1d")
		}
		return jsonResponse(200, `{
			"chart": {
				"result": [{"meta": {"symbol": "AAPL"}, "timestamp": [], "indicators": {"quote": [{}]}}]
			}
		}`), nil
	})

	_, err := client.GetChart("AAPL", "", "")
	if err != nil {
		t.Fatalf("GetChart() error: %v", err)
	}
}

func TestGetChart_CustomParams(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		rangeParam := req.URL.Query().Get("range")
		interval := req.URL.Query().Get("interval")
		if rangeParam != "5d" {
			t.Errorf("range = %q, want %q", rangeParam, "5d")
		}
		if interval != "15m" {
			t.Errorf("interval = %q, want %q", interval, "15m")
		}
		return jsonResponse(200, `{
			"chart": {
				"result": [{"meta": {"symbol": "AAPL"}, "timestamp": [], "indicators": {"quote": [{}]}}]
			}
		}`), nil
	})

	_, err := client.GetChart("AAPL", "5d", "15m")
	if err != nil {
		t.Fatalf("GetChart() error: %v", err)
	}
}

func TestGetChart_YahooError(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		return jsonResponse(200, `{
			"chart": {
				"result": null,
				"error": {"code": "Not Found", "description": "No data found"}
			}
		}`), nil
	})

	_, err := client.GetChart("INVALID", "1mo", "1d")
	if err == nil {
		t.Fatal("expected error for yahoo error response")
	}
	if !strings.Contains(err.Error(), "No data found") {
		t.Errorf("error should contain description, got: %v", err)
	}
}

func TestGetChart_EmptyResults(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		return jsonResponse(200, `{"chart": {"result": []}}`), nil
	})

	_, err := client.GetChart("UNKNOWN", "1mo", "1d")
	if err == nil {
		t.Fatal("expected error for empty results")
	}
	if !strings.Contains(err.Error(), "no chart data found") {
		t.Errorf("error should mention no chart data, got: %v", err)
	}
}
