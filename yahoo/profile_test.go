package yahoo

import (
	"net/http"
	"strings"
	"testing"
)

func TestGetProfile_Success(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		if !strings.Contains(req.URL.Path, "/v10/finance/quoteSummary/AAPL") {
			t.Errorf("unexpected path: %s", req.URL.Path)
		}
		modules := req.URL.Query().Get("modules")
		if modules != "assetProfile,quoteType" {
			t.Errorf("modules = %q, want %q", modules, "assetProfile,quoteType")
		}
		return jsonResponse(200, `{
			"quoteSummary": {
				"result": [{
					"assetProfile": {
						"address1": "One Apple Park Way",
						"city": "Cupertino",
						"state": "CA",
						"zip": "95014",
						"country": "United States",
						"phone": "408 996 1010",
						"website": "https://www.apple.com",
						"industry": "Consumer Electronics",
						"sector": "Technology",
						"longBusinessSummary": "Apple Inc. designs, manufactures, and markets smartphones.",
						"fullTimeEmployees": 164000,
						"companyOfficers": [
							{
								"name": "Mr. Timothy D. Cook",
								"title": "CEO & Director",
								"age": 63,
								"totalPay": {"raw": 16000000, "fmt": "16M"}
							}
						]
					},
					"quoteType": {
						"symbol": "AAPL",
						"shortName": "Apple Inc.",
						"longName": "Apple Inc.",
						"quoteType": "EQUITY",
						"exchange": "NMS"
					}
				}]
			}
		}`), nil
	})

	profile, quoteType, err := client.GetProfile("AAPL")
	if err != nil {
		t.Fatalf("GetProfile() error: %v", err)
	}

	if profile.City != "Cupertino" {
		t.Errorf("profile.City = %q, want %q", profile.City, "Cupertino")
	}
	if profile.Sector != "Technology" {
		t.Errorf("profile.Sector = %q, want %q", profile.Sector, "Technology")
	}
	if profile.FullTimeEmployees != 164000 {
		t.Errorf("profile.FullTimeEmployees = %d, want 164000", profile.FullTimeEmployees)
	}
	if len(profile.CompanyOfficers) != 1 {
		t.Fatalf("officers count = %d, want 1", len(profile.CompanyOfficers))
	}
	if profile.CompanyOfficers[0].Name != "Mr. Timothy D. Cook" {
		t.Errorf("officer name = %q, want %q", profile.CompanyOfficers[0].Name, "Mr. Timothy D. Cook")
	}

	if quoteType.Symbol != "AAPL" {
		t.Errorf("quoteType.Symbol = %q, want %q", quoteType.Symbol, "AAPL")
	}
	if quoteType.QuoteType != "EQUITY" {
		t.Errorf("quoteType.QuoteType = %q, want %q", quoteType.QuoteType, "EQUITY")
	}
}

func TestGetProfile_YahooError(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		return jsonResponse(200, `{
			"quoteSummary": {
				"result": null,
				"error": {"code": "Not Found", "description": "Symbol not found"}
			}
		}`), nil
	})

	_, _, err := client.GetProfile("INVALID")
	if err == nil {
		t.Fatal("expected error for yahoo error response")
	}
	if !strings.Contains(err.Error(), "Symbol not found") {
		t.Errorf("error should contain description, got: %v", err)
	}
}

func TestGetProfile_EmptyResults(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		return jsonResponse(200, `{"quoteSummary": {"result": []}}`), nil
	})

	_, _, err := client.GetProfile("UNKNOWN")
	if err == nil {
		t.Fatal("expected error for empty results")
	}
	if !strings.Contains(err.Error(), "no data found") {
		t.Errorf("error should mention no data found, got: %v", err)
	}
}
