package yahoo

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestGetBulkQuotes_Success(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		if !strings.Contains(req.URL.Path, "/v7/finance/quote") {
			t.Errorf("unexpected path: %s", req.URL.Path)
		}
		syms := req.URL.Query().Get("symbols")
		if syms != "AAPL,MSFT" {
			t.Errorf("symbols = %q, want %q", syms, "AAPL,MSFT")
		}
		return jsonResponse(200, `{
			"quoteResponse": {
				"result": [
					{
						"symbol": "AAPL",
						"shortName": "Apple Inc.",
						"longName": "Apple Inc.",
						"currency": "USD",
						"exchange": "NMS",
						"fullExchangeName": "NasdaqGS",
						"quoteType": "EQUITY",
						"marketState": "REGULAR",
						"regularMarketPrice": 178.72,
						"regularMarketChange": 2.15,
						"regularMarketChangePercent": 1.22,
						"regularMarketVolume": 54321000,
						"regularMarketOpen": 176.50,
						"regularMarketDayHigh": 179.00,
						"regularMarketDayLow": 175.80,
						"regularMarketPreviousClose": 176.57,
						"marketCap": 2800000000000,
						"trailingPE": 29.5,
						"forwardPE": 27.1,
						"fiftyTwoWeekLow": 140.0,
						"fiftyTwoWeekHigh": 199.0,
						"fiftyDayAverage": 175.0,
						"twoHundredDayAverage": 170.0,
						"trailingAnnualDividendYield": 0.0055
					},
					{
						"symbol": "MSFT",
						"shortName": "Microsoft Corporation",
						"regularMarketPrice": 380.50,
						"regularMarketVolume": 22000000,
						"marketCap": 2830000000000
					}
				],
				"error": null
			}
		}`), nil
	})

	results, err := client.GetBulkQuotes([]string{"AAPL", "MSFT"})
	if err != nil {
		t.Fatalf("GetBulkQuotes() error: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("got %d results, want 2", len(results))
	}

	aapl := results[0]
	if aapl.Symbol != "AAPL" {
		t.Errorf("results[0].Symbol = %q, want %q", aapl.Symbol, "AAPL")
	}
	if aapl.RegularMarketPrice != 178.72 {
		t.Errorf("results[0].RegularMarketPrice = %v, want 178.72", aapl.RegularMarketPrice)
	}
	if aapl.MarketCap != 2800000000000 {
		t.Errorf("results[0].MarketCap = %v, want 2800000000000", aapl.MarketCap)
	}
	if aapl.TrailingPE != 29.5 {
		t.Errorf("results[0].TrailingPE = %v, want 29.5", aapl.TrailingPE)
	}
	if aapl.FiftyTwoWeekHigh != 199.0 {
		t.Errorf("results[0].FiftyTwoWeekHigh = %v, want 199.0", aapl.FiftyTwoWeekHigh)
	}

	msft := results[1]
	if msft.Symbol != "MSFT" {
		t.Errorf("results[1].Symbol = %q, want %q", msft.Symbol, "MSFT")
	}
	if msft.RegularMarketPrice != 380.50 {
		t.Errorf("results[1].RegularMarketPrice = %v, want 380.50", msft.RegularMarketPrice)
	}
}

func TestGetBulkQuotes_SingleSymbol(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		syms := req.URL.Query().Get("symbols")
		if syms != "TSLA" {
			t.Errorf("symbols = %q, want %q", syms, "TSLA")
		}
		return jsonResponse(200, `{
			"quoteResponse": {
				"result": [{"symbol": "TSLA", "regularMarketPrice": 250.0}],
				"error": null
			}
		}`), nil
	})

	results, err := client.GetBulkQuotes([]string{"TSLA"})
	if err != nil {
		t.Fatalf("GetBulkQuotes() error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("got %d results, want 1", len(results))
	}
	if results[0].Symbol != "TSLA" {
		t.Errorf("Symbol = %q, want %q", results[0].Symbol, "TSLA")
	}
}

func TestGetBulkQuotes_EmptySymbols(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		t.Fatal("should not make HTTP request for empty symbols")
		return nil, nil
	})

	_, err := client.GetBulkQuotes([]string{})
	if err == nil {
		t.Fatal("expected error for empty symbols")
	}
	if !strings.Contains(err.Error(), "at least one symbol") {
		t.Errorf("error should mention at least one symbol, got: %v", err)
	}
}

func TestGetBulkQuotes_YahooError(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		return jsonResponse(200, `{
			"quoteResponse": {
				"result": null,
				"error": {"code": "Not Found", "description": "No data found for symbols"}
			}
		}`), nil
	})

	_, err := client.GetBulkQuotes([]string{"INVALID1", "INVALID2"})
	if err == nil {
		t.Fatal("expected error for yahoo error response")
	}
	if !strings.Contains(err.Error(), "No data found") {
		t.Errorf("error should contain yahoo error description, got: %v", err)
	}
}

func TestGetBulkQuotes_TooManySymbols(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		t.Fatal("should not make HTTP request for too many symbols")
		return nil, nil
	})

	symbols := make([]string, maxBulkSymbols+1)
	for i := range symbols {
		symbols[i] = fmt.Sprintf("SYM%d", i)
	}

	_, err := client.GetBulkQuotes(symbols)
	if err == nil {
		t.Fatal("expected error for too many symbols")
	}
	if !strings.Contains(err.Error(), "too many symbols") {
		t.Errorf("error should mention too many symbols, got: %v", err)
	}
}

func TestGetBulkQuotes_HTTPError(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		return jsonResponse(500, `Internal Server Error`), nil
	})

	_, err := client.GetBulkQuotes([]string{"AAPL"})
	if err == nil {
		t.Fatal("expected error for HTTP 500")
	}
}
