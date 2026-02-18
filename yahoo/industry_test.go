package yahoo

import (
	"net/http"
	"strings"
	"testing"
)

func TestGetIndustry_Success(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		if !strings.Contains(req.URL.Path, "/v1/finance/industries/consumer-electronics") {
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
				"sectorKey": "technology",
				"sectorName": "Technology",
				"name": "Consumer Electronics",
				"symbol": "consumer-electronics",
				"overview": {
					"companiesCount": 45,
					"marketCap": {"raw": 3500000000000, "fmt": "3.5T"},
					"description": "Consumer electronics industry",
					"industriesCount": 0,
					"marketWeight": {"raw": 0.08, "fmt": "8%"},
					"employeeCount": {"raw": 1200000, "fmt": "1.2M"}
				},
				"topCompanies": [
					{"symbol": "AAPL", "name": "Apple Inc.", "rating": "Buy", "marketWeight": {"raw": 0.65, "fmt": "65%"}}
				],
				"topPerformingCompanies": [
					{
						"symbol": "AAPL",
						"name": "Apple Inc.",
						"ytdReturn": {"raw": 0.32, "fmt": "32%"},
						"lastPrice": {"raw": 178.72, "fmt": "178.72"},
						"targetPrice": {"raw": 200.0, "fmt": "200.00"}
					}
				],
				"topGrowthCompanies": [
					{
						"symbol": "SONY",
						"name": "Sony Group Corporation",
						"ytdReturn": {"raw": 0.15, "fmt": "15%"},
						"growthEstimate": {"raw": 0.22, "fmt": "22%"}
					}
				]
			}
		}`), nil
	})

	data, err := client.GetIndustry("consumer-electronics")
	if err != nil {
		t.Fatalf("GetIndustry() error: %v", err)
	}

	if data.SectorKey != "technology" {
		t.Errorf("SectorKey = %q, want %q", data.SectorKey, "technology")
	}
	if data.SectorName != "Technology" {
		t.Errorf("SectorName = %q, want %q", data.SectorName, "Technology")
	}
	if data.Name != "Consumer Electronics" {
		t.Errorf("Name = %q, want %q", data.Name, "Consumer Electronics")
	}
	if data.Symbol != "consumer-electronics" {
		t.Errorf("Symbol = %q, want %q", data.Symbol, "consumer-electronics")
	}
	if data.Overview.CompaniesCount != 45 {
		t.Errorf("Overview.CompaniesCount = %d, want 45", data.Overview.CompaniesCount)
	}
	if data.Overview.MarketCap.Raw != 3500000000000 {
		t.Errorf("Overview.MarketCap.Raw = %v, want 3500000000000", data.Overview.MarketCap.Raw)
	}

	if len(data.TopCompanies) != 1 {
		t.Fatalf("TopCompanies length = %d, want 1", len(data.TopCompanies))
	}
	if data.TopCompanies[0].Symbol != "AAPL" {
		t.Errorf("TopCompanies[0].Symbol = %q, want %q", data.TopCompanies[0].Symbol, "AAPL")
	}

	if len(data.TopPerformingCompanies) != 1 {
		t.Fatalf("TopPerformingCompanies length = %d, want 1", len(data.TopPerformingCompanies))
	}
	perf := data.TopPerformingCompanies[0]
	if perf.Symbol != "AAPL" {
		t.Errorf("TopPerformingCompanies[0].Symbol = %q, want %q", perf.Symbol, "AAPL")
	}
	if perf.YtdReturn.Raw != 0.32 {
		t.Errorf("TopPerformingCompanies[0].YtdReturn.Raw = %v, want 0.32", perf.YtdReturn.Raw)
	}
	if perf.TargetPrice.Raw != 200.0 {
		t.Errorf("TopPerformingCompanies[0].TargetPrice.Raw = %v, want 200.0", perf.TargetPrice.Raw)
	}

	if len(data.TopGrowthCompanies) != 1 {
		t.Fatalf("TopGrowthCompanies length = %d, want 1", len(data.TopGrowthCompanies))
	}
	growth := data.TopGrowthCompanies[0]
	if growth.Symbol != "SONY" {
		t.Errorf("TopGrowthCompanies[0].Symbol = %q, want %q", growth.Symbol, "SONY")
	}
	if growth.GrowthEstimate.Raw != 0.22 {
		t.Errorf("TopGrowthCompanies[0].GrowthEstimate.Raw = %v, want 0.22", growth.GrowthEstimate.Raw)
	}
}

func TestGetIndustry_HTTPError(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		return jsonResponse(500, `Internal Server Error`), nil
	})

	_, err := client.GetIndustry("consumer-electronics")
	if err == nil {
		t.Fatal("expected error for HTTP 500")
	}
}

func TestGetIndustry_PathEscaping(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		if !strings.Contains(req.URL.Path, "/industries/software-application") {
			t.Errorf("unexpected path: %s", req.URL.Path)
		}
		return jsonResponse(200, `{"data": {}}`), nil
	})

	client.GetIndustry("software-application")
}
