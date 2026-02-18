//go:build integration

package yahoo

import (
	"os"
	"testing"
)

const testSymbol = "AAPL"

// Shared client so authentication happens once across all integration tests.
var integrationClient *Client

func TestMain(m *testing.M) {
	integrationClient = NewClient()
	os.Exit(m.Run())
}

func TestIntegration_GetQuote(t *testing.T) {
	price, detail, err := integrationClient.GetQuote(testSymbol)
	if err != nil {
		t.Fatalf("GetQuote(%q) error: %v", testSymbol, err)
	}

	if price.Symbol != testSymbol {
		t.Errorf("price.Symbol = %q, want %q", price.Symbol, testSymbol)
	}
	if price.RegularMarketPrice.Raw <= 0 {
		t.Errorf("price.RegularMarketPrice.Raw = %v, want > 0", price.RegularMarketPrice.Raw)
	}
	if price.Currency == "" {
		t.Error("price.Currency is empty")
	}
	if detail.FiftyTwoWeekHigh.Raw <= 0 {
		t.Errorf("detail.FiftyTwoWeekHigh.Raw = %v, want > 0", detail.FiftyTwoWeekHigh.Raw)
	}
}

func TestIntegration_GetChart(t *testing.T) {
	result, err := integrationClient.GetChart(testSymbol, "5d", "1d")
	if err != nil {
		t.Fatalf("GetChart(%q) error: %v", testSymbol, err)
	}

	if result.Meta.Symbol != testSymbol {
		t.Errorf("meta.Symbol = %q, want %q", result.Meta.Symbol, testSymbol)
	}
	if result.Meta.Currency == "" {
		t.Error("meta.Currency is empty")
	}
	if len(result.Timestamps) == 0 {
		t.Error("expected at least one timestamp")
	}
	if len(result.Indicators.Quote) == 0 {
		t.Error("expected at least one quote indicator")
	}
}

func TestIntegration_Search(t *testing.T) {
	resp, err := integrationClient.Search("Apple", 5)
	if err != nil {
		t.Fatalf("Search(%q) error: %v", "Apple", err)
	}

	if len(resp.Quotes) == 0 {
		t.Fatal("expected at least one search result")
	}

	found := false
	for _, q := range resp.Quotes {
		if q.Symbol == testSymbol {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected %q in search results", testSymbol)
	}
}

func TestIntegration_GetNews(t *testing.T) {
	news, err := integrationClient.GetNews(testSymbol, 3)
	if err != nil {
		t.Fatalf("GetNews(%q) error: %v", testSymbol, err)
	}

	if len(news) == 0 {
		t.Skip("no news articles returned (may vary by time of day)")
	}
	if news[0].Title == "" {
		t.Error("first news article has empty title")
	}
	if news[0].Link == "" {
		t.Error("first news article has empty link")
	}
}

func TestIntegration_GetOptions(t *testing.T) {
	result, err := integrationClient.GetOptions(testSymbol, "")
	if err != nil {
		t.Fatalf("GetOptions(%q) error: %v", testSymbol, err)
	}

	if result.UnderlyingSymbol != testSymbol {
		t.Errorf("UnderlyingSymbol = %q, want %q", result.UnderlyingSymbol, testSymbol)
	}
	if len(result.ExpirationDates) == 0 {
		t.Error("expected at least one expiration date")
	}
	if len(result.Options) == 0 {
		t.Fatal("expected at least one options chain")
	}
	if len(result.Options[0].Calls) == 0 {
		t.Error("expected at least one call option")
	}
	if len(result.Options[0].Puts) == 0 {
		t.Error("expected at least one put option")
	}
}

func TestIntegration_GetProfile(t *testing.T) {
	profile, quoteType, err := integrationClient.GetProfile(testSymbol)
	if err != nil {
		t.Fatalf("GetProfile(%q) error: %v", testSymbol, err)
	}

	if profile.Sector == "" {
		t.Error("profile.Sector is empty")
	}
	if profile.Industry == "" {
		t.Error("profile.Industry is empty")
	}
	if profile.Country == "" {
		t.Error("profile.Country is empty")
	}
	if profile.FullTimeEmployees <= 0 {
		t.Errorf("profile.FullTimeEmployees = %d, want > 0", profile.FullTimeEmployees)
	}
	if quoteType.Symbol != testSymbol {
		t.Errorf("quoteType.Symbol = %q, want %q", quoteType.Symbol, testSymbol)
	}
}

func TestIntegration_GetFinancials(t *testing.T) {
	for _, tc := range []struct {
		statement string
		quarterly bool
	}{
		{"income", false},
		{"income", true},
		{"balance", false},
		{"cashflow", false},
	} {
		label := tc.statement
		if tc.quarterly {
			label = "quarterly " + label
		}
		t.Run(label, func(t *testing.T) {
			results, err := integrationClient.GetFinancials(testSymbol, tc.statement, tc.quarterly)
			if err != nil {
				t.Fatalf("GetFinancials(%q, %q, %v) error: %v", testSymbol, tc.statement, tc.quarterly, err)
			}
			if len(results) == 0 {
				t.Fatal("expected at least one financial result")
			}
			if results[0].Type == "" {
				t.Error("first result has empty Type")
			}
			if len(results[0].Items) == 0 {
				t.Error("first result has no data items")
			}
		})
	}
}

func TestIntegration_GetRecommendations(t *testing.T) {
	result, err := integrationClient.GetRecommendations(testSymbol)
	if err != nil {
		t.Fatalf("GetRecommendations(%q) error: %v", testSymbol, err)
	}

	if len(result.Trend) == 0 {
		t.Fatal("expected at least one recommendation trend")
	}
	total := result.Trend[0].StrongBuy + result.Trend[0].Buy + result.Trend[0].Hold + result.Trend[0].Sell + result.Trend[0].StrongSell
	if total == 0 {
		t.Error("expected at least one analyst recommendation in current period")
	}
}
