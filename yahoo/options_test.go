package yahoo

import (
	"net/http"
	"strings"
	"testing"
)

func TestGetOptions_Success(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		if !strings.Contains(req.URL.Path, "/v7/finance/options/AAPL") {
			t.Errorf("unexpected path: %s", req.URL.Path)
		}
		return jsonResponse(200, `{
			"optionChain": {
				"result": [{
					"underlyingSymbol": "AAPL",
					"expirationDates": [1700000000, 1700604800],
					"strikes": [170.0, 175.0, 180.0],
					"quote": {
						"symbol": "AAPL",
						"regularMarketPrice": 178.72
					},
					"options": [{
						"expirationDate": 1700000000,
						"calls": [{
							"contractSymbol": "AAPL231117C00170000",
							"strike": 170.0,
							"lastPrice": 9.50,
							"volume": 1500,
							"openInterest": 5000,
							"impliedVolatility": 0.35,
							"inTheMoney": true
						}],
						"puts": [{
							"contractSymbol": "AAPL231117P00180000",
							"strike": 180.0,
							"lastPrice": 2.10,
							"volume": 800,
							"openInterest": 3000,
							"impliedVolatility": 0.30,
							"inTheMoney": true
						}]
					}]
				}]
			}
		}`), nil
	})

	result, err := client.GetOptions("AAPL", "")
	if err != nil {
		t.Fatalf("GetOptions() error: %v", err)
	}

	if result.UnderlyingSymbol != "AAPL" {
		t.Errorf("UnderlyingSymbol = %q, want %q", result.UnderlyingSymbol, "AAPL")
	}
	if len(result.ExpirationDates) != 2 {
		t.Errorf("expirationDates count = %d, want 2", len(result.ExpirationDates))
	}
	if len(result.Strikes) != 3 {
		t.Errorf("strikes count = %d, want 3", len(result.Strikes))
	}
	if len(result.Options) != 1 {
		t.Fatalf("options chains count = %d, want 1", len(result.Options))
	}
	if len(result.Options[0].Calls) != 1 {
		t.Errorf("calls count = %d, want 1", len(result.Options[0].Calls))
	}
	if result.Options[0].Calls[0].Strike != 170.0 {
		t.Errorf("call strike = %v, want 170.0", result.Options[0].Calls[0].Strike)
	}
	if !result.Options[0].Calls[0].InTheMoney {
		t.Error("call should be in the money")
	}
	if len(result.Options[0].Puts) != 1 {
		t.Errorf("puts count = %d, want 1", len(result.Options[0].Puts))
	}
}

func TestGetOptions_WithExpiration(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		date := req.URL.Query().Get("date")
		if date != "1700000000" {
			t.Errorf("date = %q, want %q", date, "1700000000")
		}
		return jsonResponse(200, `{
			"optionChain": {
				"result": [{"underlyingSymbol": "AAPL", "options": []}]
			}
		}`), nil
	})

	_, err := client.GetOptions("AAPL", "1700000000")
	if err != nil {
		t.Fatalf("GetOptions() error: %v", err)
	}
}

func TestGetOptions_NoExpiration(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		date := req.URL.Query().Get("date")
		if date != "" {
			t.Errorf("date should be empty when no expiration, got %q", date)
		}
		return jsonResponse(200, `{
			"optionChain": {
				"result": [{"underlyingSymbol": "AAPL", "options": []}]
			}
		}`), nil
	})

	_, err := client.GetOptions("AAPL", "")
	if err != nil {
		t.Fatalf("GetOptions() error: %v", err)
	}
}

func TestGetOptions_YahooError(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		return jsonResponse(200, `{
			"optionChain": {
				"result": null,
				"error": {"code": "Not Found", "description": "No options data"}
			}
		}`), nil
	})

	_, err := client.GetOptions("INVALID", "")
	if err == nil {
		t.Fatal("expected error for yahoo error response")
	}
	if !strings.Contains(err.Error(), "No options data") {
		t.Errorf("error should contain description, got: %v", err)
	}
}

func TestGetOptions_EmptyResults(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		return jsonResponse(200, `{"optionChain": {"result": []}}`), nil
	})

	_, err := client.GetOptions("UNKNOWN", "")
	if err == nil {
		t.Fatal("expected error for empty results")
	}
	if !strings.Contains(err.Error(), "no options data found") {
		t.Errorf("error should mention no options data, got: %v", err)
	}
}
