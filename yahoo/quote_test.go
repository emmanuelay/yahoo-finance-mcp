package yahoo

import (
	"net/http"
	"strings"
	"testing"
)

func TestGetQuote_Success(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		if !strings.Contains(req.URL.Path, "/v10/finance/quoteSummary/AAPL") {
			t.Errorf("unexpected path: %s", req.URL.Path)
		}
		return jsonResponse(200, `{
			"quoteSummary": {
				"result": [{
					"price": {
						"symbol": "AAPL",
						"shortName": "Apple Inc.",
						"currency": "USD",
						"regularMarketPrice": {"raw": 178.72, "fmt": "178.72"},
						"regularMarketChange": {"raw": 2.15, "fmt": "2.15"},
						"regularMarketChangePercent": {"raw": 0.0122, "fmt": "1.22%"},
						"regularMarketVolume": {"raw": 54321000, "fmt": "54.32M"},
						"marketCap": {"raw": 2800000000000, "fmt": "2.8T"}
					},
					"summaryDetail": {
						"trailingPE": {"raw": 29.5, "fmt": "29.50"},
						"forwardPE": {"raw": 27.1, "fmt": "27.10"},
						"dividendYield": {"raw": 0.0055, "fmt": "0.55%"},
						"fiftyTwoWeekLow": {"raw": 140.0, "fmt": "140.00"},
						"fiftyTwoWeekHigh": {"raw": 199.0, "fmt": "199.00"}
					}
				}]
			}
		}`), nil
	})

	price, detail, err := client.GetQuote("AAPL")
	if err != nil {
		t.Fatalf("GetQuote() error: %v", err)
	}

	if price.Symbol != "AAPL" {
		t.Errorf("price.Symbol = %q, want %q", price.Symbol, "AAPL")
	}
	if price.RegularMarketPrice.Raw != 178.72 {
		t.Errorf("price.RegularMarketPrice.Raw = %v, want 178.72", price.RegularMarketPrice.Raw)
	}
	if price.MarketCap.Raw != 2800000000000 {
		t.Errorf("price.MarketCap.Raw = %v, want 2800000000000", price.MarketCap.Raw)
	}
	if detail.TrailingPE.Raw != 29.5 {
		t.Errorf("detail.TrailingPE.Raw = %v, want 29.5", detail.TrailingPE.Raw)
	}
	if detail.FiftyTwoWeekHigh.Raw != 199.0 {
		t.Errorf("detail.FiftyTwoWeekHigh.Raw = %v, want 199.0", detail.FiftyTwoWeekHigh.Raw)
	}
}

func TestGetQuote_YahooError(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		return jsonResponse(200, `{
			"quoteSummary": {
				"result": null,
				"error": {"code": "Not Found", "description": "No data found for symbol XYZ123"}
			}
		}`), nil
	})

	_, _, err := client.GetQuote("XYZ123")
	if err == nil {
		t.Fatal("expected error for yahoo error response")
	}
	if !strings.Contains(err.Error(), "No data found") {
		t.Errorf("error should contain yahoo error description, got: %v", err)
	}
}

func TestGetQuote_EmptyResults(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		return jsonResponse(200, `{
			"quoteSummary": {
				"result": []
			}
		}`), nil
	})

	_, _, err := client.GetQuote("UNKNOWN")
	if err == nil {
		t.Fatal("expected error for empty results")
	}
	if !strings.Contains(err.Error(), "no data found") {
		t.Errorf("error should mention no data found, got: %v", err)
	}
}

func TestGetQuote_RequestModules(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		modules := req.URL.Query().Get("modules")
		if modules != "price,summaryDetail" {
			t.Errorf("modules = %q, want %q", modules, "price,summaryDetail")
		}
		return jsonResponse(200, `{
			"quoteSummary": {
				"result": [{"price": {"symbol": "AAPL"}, "summaryDetail": {}}]
			}
		}`), nil
	})

	client.GetQuote("AAPL")
}
