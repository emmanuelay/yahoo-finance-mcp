package yahoo

import (
	"net/http"
	"strings"
	"testing"
)

func TestGetSector_Success(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		if !strings.Contains(req.URL.Path, "/v1/finance/sectors/technology") {
			t.Errorf("unexpected path: %s", req.URL.Path)
		}
		q := req.URL.Query()
		if q.Get("formatted") != "true" {
			t.Errorf("formatted = %q, want %q", q.Get("formatted"), "true")
		}
		if q.Get("withReturns") != "true" {
			t.Errorf("withReturns = %q, want %q", q.Get("withReturns"), "true")
		}
		return jsonResponse(200, `{
			"data": {
				"name": "Technology",
				"symbol": "technology",
				"overview": {
					"companiesCount": 950,
					"marketCap": {"raw": 15000000000000, "fmt": "15T"},
					"description": "Technology sector overview",
					"industriesCount": 12,
					"marketWeight": {"raw": 0.28, "fmt": "28%"},
					"employeeCount": {"raw": 5000000, "fmt": "5M"}
				},
				"topCompanies": [
					{"symbol": "AAPL", "name": "Apple Inc.", "rating": "Buy", "marketWeight": {"raw": 0.15, "fmt": "15%"}},
					{"symbol": "MSFT", "name": "Microsoft Corporation", "rating": "Buy", "marketWeight": {"raw": 0.12, "fmt": "12%"}}
				],
				"topETFs": [
					{"symbol": "XLK", "name": "Technology Select Sector SPDR Fund"}
				],
				"topMutualFunds": [
					{"symbol": "FTEC", "name": "Fidelity MSCI Information Technology Index ETF"}
				],
				"industries": [
					{"key": "consumer-electronics", "name": "Consumer Electronics", "symbol": "consumer-electronics", "marketWeight": {"raw": 0.20, "fmt": "20%"}},
					{"key": "semiconductors", "name": "Semiconductors", "symbol": "semiconductors", "marketWeight": {"raw": 0.18, "fmt": "18%"}}
				]
			}
		}`), nil
	})

	data, err := client.GetSector("technology")
	if err != nil {
		t.Fatalf("GetSector() error: %v", err)
	}

	if data.Name != "Technology" {
		t.Errorf("Name = %q, want %q", data.Name, "Technology")
	}
	if data.Symbol != "technology" {
		t.Errorf("Symbol = %q, want %q", data.Symbol, "technology")
	}
	if data.Overview.CompaniesCount != 950 {
		t.Errorf("Overview.CompaniesCount = %d, want 950", data.Overview.CompaniesCount)
	}
	if data.Overview.MarketCap.Raw != 15000000000000 {
		t.Errorf("Overview.MarketCap.Raw = %v, want 15000000000000", data.Overview.MarketCap.Raw)
	}
	if data.Overview.IndustriesCount != 12 {
		t.Errorf("Overview.IndustriesCount = %d, want 12", data.Overview.IndustriesCount)
	}
	if len(data.TopCompanies) != 2 {
		t.Fatalf("TopCompanies length = %d, want 2", len(data.TopCompanies))
	}
	if data.TopCompanies[0].Symbol != "AAPL" {
		t.Errorf("TopCompanies[0].Symbol = %q, want %q", data.TopCompanies[0].Symbol, "AAPL")
	}
	if len(data.TopETFs) != 1 {
		t.Fatalf("TopETFs length = %d, want 1", len(data.TopETFs))
	}
	if data.TopETFs[0].Symbol != "XLK" {
		t.Errorf("TopETFs[0].Symbol = %q, want %q", data.TopETFs[0].Symbol, "XLK")
	}
	if len(data.TopMutualFunds) != 1 {
		t.Fatalf("TopMutualFunds length = %d, want 1", len(data.TopMutualFunds))
	}
	if len(data.Industries) != 2 {
		t.Fatalf("Industries length = %d, want 2", len(data.Industries))
	}
	if data.Industries[0].Key != "consumer-electronics" {
		t.Errorf("Industries[0].Key = %q, want %q", data.Industries[0].Key, "consumer-electronics")
	}
}

func TestGetSector_HTTPError(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		return jsonResponse(500, `Internal Server Error`), nil
	})

	_, err := client.GetSector("technology")
	if err == nil {
		t.Fatal("expected error for HTTP 500")
	}
}

func TestValidSectorKeys(t *testing.T) {
	expected := []string{
		"basic-materials", "communication-services", "consumer-cyclical",
		"consumer-defensive", "energy", "financial-services", "healthcare",
		"industrials", "real-estate", "technology", "utilities",
	}
	if len(ValidSectorKeys) != len(expected) {
		t.Fatalf("ValidSectorKeys length = %d, want %d", len(ValidSectorKeys), len(expected))
	}
	for i, key := range expected {
		if ValidSectorKeys[i] != key {
			t.Errorf("ValidSectorKeys[%d] = %q, want %q", i, ValidSectorKeys[i], key)
		}
	}
}
