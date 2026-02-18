package yahoo

import (
	"net/http"
	"strings"
	"testing"
)

func TestGetMarketSummary_Success(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		if !strings.Contains(req.URL.String(), "/v6/finance/quote/marketSummary") {
			t.Errorf("unexpected URL: %s", req.URL.String())
		}
		if m := req.URL.Query().Get("market"); m != "US" {
			t.Errorf("market = %q, want %q", m, "US")
		}
		return jsonResponse(200, `{
			"marketSummaryResponse": {
				"result": [
					{
						"symbol": "^GSPC",
						"shortName": "S&P 500",
						"exchange": "SNP",
						"marketState": "REGULAR",
						"regularMarketPrice": 5021.84,
						"regularMarketChange": 25.60,
						"regularMarketChangePercent": 0.51
					},
					{
						"symbol": "^DJI",
						"shortName": "Dow Jones Industrial Average",
						"exchange": "DJI",
						"marketState": "REGULAR",
						"regularMarketPrice": 38996.39,
						"regularMarketChange": 134.21,
						"regularMarketChangePercent": 0.35
					}
				],
				"error": null
			}
		}`), nil
	})

	results, err := client.GetMarketSummary("US")
	if err != nil {
		t.Fatalf("GetMarketSummary() error: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("got %d results, want 2", len(results))
	}

	sp := results[0]
	if sp.Symbol != "^GSPC" {
		t.Errorf("results[0].Symbol = %q, want %q", sp.Symbol, "^GSPC")
	}
	if sp.ShortName != "S&P 500" {
		t.Errorf("results[0].ShortName = %q, want %q", sp.ShortName, "S&P 500")
	}
	if sp.RegularMarketPrice != 5021.84 {
		t.Errorf("results[0].RegularMarketPrice = %v, want 5021.84", sp.RegularMarketPrice)
	}
	if sp.RegularMarketChange != 25.60 {
		t.Errorf("results[0].RegularMarketChange = %v, want 25.60", sp.RegularMarketChange)
	}
	if sp.MarketState != "REGULAR" {
		t.Errorf("results[0].MarketState = %q, want %q", sp.MarketState, "REGULAR")
	}

	dji := results[1]
	if dji.Symbol != "^DJI" {
		t.Errorf("results[1].Symbol = %q, want %q", dji.Symbol, "^DJI")
	}
}

func TestGetMarketSummary_YahooError(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		return jsonResponse(200, `{
			"marketSummaryResponse": {
				"result": null,
				"error": {"code": "Bad Request", "description": "Invalid market key"}
			}
		}`), nil
	})

	_, err := client.GetMarketSummary("INVALID")
	if err == nil {
		t.Fatal("expected error for yahoo error response")
	}
	if !strings.Contains(err.Error(), "Invalid market key") {
		t.Errorf("error should contain yahoo error description, got: %v", err)
	}
}

func TestGetMarketSummary_HTTPError(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		return jsonResponse(500, `Internal Server Error`), nil
	})

	_, err := client.GetMarketSummary("US")
	if err == nil {
		t.Fatal("expected error for HTTP 500")
	}
}

func TestGetMarketStatus_Success(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		if !strings.Contains(req.URL.String(), "/v6/finance/markettime") {
			t.Errorf("unexpected URL: %s", req.URL.String())
		}
		if m := req.URL.Query().Get("market"); m != "US" {
			t.Errorf("market = %q, want %q", m, "US")
		}
		return jsonResponse(200, `{
			"finance": {
				"marketTimes": [
					{
						"marketTime": [
							{
								"id": "us_market",
								"name": "US Markets",
								"status": "open",
								"message": "Market is open",
								"open": "09:30",
								"close": "16:00",
								"time": "2024-01-15T12:00:00Z",
								"timezone": [
									{"gmtoffset": "-05:00", "short": "EST"}
								]
							}
						]
					}
				],
				"error": null
			}
		}`), nil
	})

	groups, err := client.GetMarketStatus("US")
	if err != nil {
		t.Fatalf("GetMarketStatus() error: %v", err)
	}

	if len(groups) != 1 {
		t.Fatalf("got %d groups, want 1", len(groups))
	}

	if len(groups[0].MarketTime) != 1 {
		t.Fatalf("MarketTime length = %d, want 1", len(groups[0].MarketTime))
	}

	mt := groups[0].MarketTime[0]
	if mt.ID != "us_market" {
		t.Errorf("ID = %q, want %q", mt.ID, "us_market")
	}
	if mt.Name != "US Markets" {
		t.Errorf("Name = %q, want %q", mt.Name, "US Markets")
	}
	if mt.Status != "open" {
		t.Errorf("Status = %q, want %q", mt.Status, "open")
	}
	if mt.Open != "09:30" {
		t.Errorf("Open = %q, want %q", mt.Open, "09:30")
	}
	if mt.Close != "16:00" {
		t.Errorf("Close = %q, want %q", mt.Close, "16:00")
	}
	if len(mt.Timezone) != 1 {
		t.Fatalf("Timezone length = %d, want 1", len(mt.Timezone))
	}
	if mt.Timezone[0].Short != "EST" {
		t.Errorf("Timezone[0].Short = %q, want %q", mt.Timezone[0].Short, "EST")
	}
}

func TestGetMarketStatus_YahooError(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		return jsonResponse(200, `{
			"finance": {
				"marketTimes": null,
				"error": {"code": "Bad Request", "description": "Invalid market"}
			}
		}`), nil
	})

	_, err := client.GetMarketStatus("INVALID")
	if err == nil {
		t.Fatal("expected error for yahoo error response")
	}
	if !strings.Contains(err.Error(), "Invalid market") {
		t.Errorf("error should contain yahoo error description, got: %v", err)
	}
}

func TestGetMarketStatus_HTTPError(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		return jsonResponse(500, `Internal Server Error`), nil
	})

	_, err := client.GetMarketStatus("US")
	if err == nil {
		t.Fatal("expected error for HTTP 500")
	}
}

func TestValidMarketKeys(t *testing.T) {
	expected := []string{
		"US", "GB", "ASIA", "EUROPE",
		"RATES", "COMMODITIES", "CURRENCIES", "CRYPTOCURRENCIES",
	}
	if len(ValidMarketKeys) != len(expected) {
		t.Fatalf("ValidMarketKeys length = %d, want %d", len(ValidMarketKeys), len(expected))
	}
	for i, key := range expected {
		if ValidMarketKeys[i] != key {
			t.Errorf("ValidMarketKeys[%d] = %q, want %q", i, ValidMarketKeys[i], key)
		}
	}
}
