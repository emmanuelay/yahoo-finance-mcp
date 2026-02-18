package yahoo

import (
	"net/http"
	"strings"
	"testing"
)

func TestGetFinancialTypes_Income(t *testing.T) {
	types := getFinancialTypes("income", false)
	if len(types) == 0 {
		t.Fatal("expected income statement types")
	}
	if types[0] != "annualTotalRevenue" {
		t.Errorf("first income type = %q, want %q", types[0], "annualTotalRevenue")
	}
}

func TestGetFinancialTypes_IncomeQuarterly(t *testing.T) {
	types := getFinancialTypes("income", true)
	if len(types) == 0 {
		t.Fatal("expected quarterly income types")
	}
	if types[0] != "quarterlyTotalRevenue" {
		t.Errorf("first quarterly income type = %q, want %q", types[0], "quarterlyTotalRevenue")
	}
}

func TestGetFinancialTypes_Balance(t *testing.T) {
	types := getFinancialTypes("balance", false)
	if len(types) == 0 {
		t.Fatal("expected balance sheet types")
	}
	if types[0] != "annualTotalAssets" {
		t.Errorf("first balance type = %q, want %q", types[0], "annualTotalAssets")
	}
}

func TestGetFinancialTypes_BalanceQuarterly(t *testing.T) {
	types := getFinancialTypes("balance", true)
	if len(types) == 0 {
		t.Fatal("expected quarterly balance sheet types")
	}
	if types[0] != "quarterlyTotalAssets" {
		t.Errorf("first quarterly balance type = %q, want %q", types[0], "quarterlyTotalAssets")
	}
}

func TestGetFinancialTypes_CashFlow(t *testing.T) {
	types := getFinancialTypes("cashflow", false)
	if len(types) == 0 {
		t.Fatal("expected cash flow types")
	}
	if types[0] != "annualOperatingCashFlow" {
		t.Errorf("first cashflow type = %q, want %q", types[0], "annualOperatingCashFlow")
	}
}

func TestGetFinancialTypes_CashFlowAlias(t *testing.T) {
	types := getFinancialTypes("cash_flow", false)
	if len(types) == 0 {
		t.Fatal("expected cash_flow alias to return types")
	}
	if types[0] != "annualOperatingCashFlow" {
		t.Errorf("cash_flow alias first type = %q, want %q", types[0], "annualOperatingCashFlow")
	}
}

func TestGetFinancialTypes_CashFlowQuarterly(t *testing.T) {
	types := getFinancialTypes("cashflow", true)
	if len(types) == 0 {
		t.Fatal("expected quarterly cash flow types")
	}
	if types[0] != "quarterlyOperatingCashFlow" {
		t.Errorf("first quarterly cashflow type = %q, want %q", types[0], "quarterlyOperatingCashFlow")
	}
}

func TestGetFinancialTypes_CaseInsensitive(t *testing.T) {
	types := getFinancialTypes("INCOME", false)
	if len(types) == 0 {
		t.Fatal("expected case-insensitive match for INCOME")
	}
}

func TestGetFinancialTypes_Invalid(t *testing.T) {
	types := getFinancialTypes("invalid", false)
	if types != nil {
		t.Errorf("expected nil for invalid statement, got %v", types)
	}
}

func TestParseFinancialsResponse_Success(t *testing.T) {
	body := []byte(`{
		"timeseries": {
			"result": [
				{
					"meta": {"symbol": ["AAPL"], "type": ["annualTotalRevenue"]},
					"timestamp": [1609459200, 1640995200],
					"annualTotalRevenue": [
						{"asOfDate": "2021-09-25", "reportedValue": {"raw": 365817000000, "fmt": "365.82B"}, "currencyCode": "USD"},
						{"asOfDate": "2022-09-24", "reportedValue": {"raw": 394328000000, "fmt": "394.33B"}, "currencyCode": "USD"}
					]
				},
				{
					"meta": {"symbol": ["AAPL"], "type": ["annualNetIncome"]},
					"timestamp": [1609459200],
					"annualNetIncome": [
						{"asOfDate": "2021-09-25", "reportedValue": {"raw": 94680000000, "fmt": "94.68B"}, "currencyCode": "USD"}
					]
				}
			]
		}
	}`)

	types := []string{"annualTotalRevenue", "annualNetIncome"}
	results, err := parseFinancialsResponse(body, types)
	if err != nil {
		t.Fatalf("parseFinancialsResponse() error: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("results count = %d, want 2", len(results))
	}

	if results[0].Type != "annualTotalRevenue" {
		t.Errorf("results[0].Type = %q, want %q", results[0].Type, "annualTotalRevenue")
	}
	if len(results[0].Items) != 2 {
		t.Fatalf("results[0] items count = %d, want 2", len(results[0].Items))
	}
	if results[0].Items[0].Date != "2021-09-25" {
		t.Errorf("results[0].Items[0].Date = %q, want %q", results[0].Items[0].Date, "2021-09-25")
	}
	if results[0].Items[0].ReportedValue != 365817000000 {
		t.Errorf("results[0].Items[0].ReportedValue = %v, want 365817000000", results[0].Items[0].ReportedValue)
	}
	if results[0].Items[0].CurrencyCode != "USD" {
		t.Errorf("results[0].Items[0].CurrencyCode = %q, want %q", results[0].Items[0].CurrencyCode, "USD")
	}

	if results[1].Type != "annualNetIncome" {
		t.Errorf("results[1].Type = %q, want %q", results[1].Type, "annualNetIncome")
	}
	if len(results[1].Items) != 1 {
		t.Errorf("results[1] items count = %d, want 1", len(results[1].Items))
	}
}

func TestParseFinancialsResponse_YahooError(t *testing.T) {
	body := []byte(`{
		"timeseries": {
			"result": null,
			"error": {"code": "Not Found", "description": "Data not available"}
		}
	}`)

	_, err := parseFinancialsResponse(body, []string{"annualTotalRevenue"})
	if err == nil {
		t.Fatal("expected error for yahoo error response")
	}
	if !strings.Contains(err.Error(), "Data not available") {
		t.Errorf("error should contain description, got: %v", err)
	}
}

func TestParseFinancialsResponse_InvalidJSON(t *testing.T) {
	body := []byte(`not-json`)
	_, err := parseFinancialsResponse(body, []string{"annualTotalRevenue"})
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestParseFinancialsResponse_EmptyResults(t *testing.T) {
	body := []byte(`{"timeseries": {"result": []}}`)
	results, err := parseFinancialsResponse(body, []string{"annualTotalRevenue"})
	if err != nil {
		t.Fatalf("parseFinancialsResponse() error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected empty results, got %d", len(results))
	}
}

func TestGetFinancials_InvalidStatement(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		t.Fatal("should not make HTTP request for invalid statement type")
		return nil, nil
	})

	_, err := client.GetFinancials("AAPL", "invalid", false)
	if err == nil {
		t.Fatal("expected error for invalid statement type")
	}
	if !strings.Contains(err.Error(), "invalid statement type") {
		t.Errorf("error should mention invalid statement type, got: %v", err)
	}
}

func TestGetFinancials_Success(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		if !strings.Contains(req.URL.Path, "/ws/fundamentals-timeseries/v1/finance/timeseries/AAPL") {
			t.Errorf("unexpected path: %s", req.URL.Path)
		}
		typeParam := req.URL.Query().Get("type")
		if !strings.Contains(typeParam, "annualTotalRevenue") {
			t.Errorf("type param should contain annualTotalRevenue, got %q", typeParam)
		}
		return jsonResponse(200, `{
			"timeseries": {
				"result": [
					{
						"meta": {"symbol": ["AAPL"], "type": ["annualTotalRevenue"]},
						"annualTotalRevenue": [
							{"asOfDate": "2023-09-30", "reportedValue": {"raw": 383285000000}, "currencyCode": "USD"}
						]
					}
				]
			}
		}`), nil
	})

	results, err := client.GetFinancials("AAPL", "income", false)
	if err != nil {
		t.Fatalf("GetFinancials() error: %v", err)
	}

	if len(results) == 0 {
		t.Fatal("expected at least one result")
	}
	if results[0].Type != "annualTotalRevenue" {
		t.Errorf("results[0].Type = %q, want %q", results[0].Type, "annualTotalRevenue")
	}
}
