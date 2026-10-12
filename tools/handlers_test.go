package tools

import (
	"strings"
	"testing"
	"time"

	"github.com/emmanuelay/yahoo-finance-mcp/yahoo"
)

func TestCurrencySymbol(t *testing.T) {
	cases := map[string]string{
		"":    "$",
		"USD": "$",
		"EUR": "€",
		"GBP": "£",
		"GBp": "£",
		"JPY": "¥",
		"CAD": "C$",
		"AUD": "A$",
		"SEK": "SEK ",
		"XYZ": "XYZ ",
	}
	for in, want := range cases {
		if got := currencySymbol(in); got != want {
			t.Errorf("currencySymbol(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestFmtPrice(t *testing.T) {
	if got := fmtPrice(123.456, "EUR"); got != "€123.46" {
		t.Errorf("fmtPrice EUR = %q", got)
	}
	if got := fmtPrice(99, "SEK"); got != "SEK 99.00" {
		t.Errorf("fmtPrice SEK = %q", got)
	}
}

func TestFmtLargeNumber(t *testing.T) {
	cases := []struct {
		val      float64
		currency string
		want     string
	}{
		{2.5e12, "USD", "$2.50T"},
		{1.25e9, "SEK", "SEK 1.25B"},
		{3e6, "", "$3.00M"},
		{1500, "EUR", "€1.50K"},
		{42, "USD", "$42.00"},
		{-7.5e9, "USD", "-$7.50B"},
	}
	for _, c := range cases {
		if got := fmtLargeNumber(c.val, c.currency); got != c.want {
			t.Errorf("fmtLargeNumber(%v, %q) = %q, want %q", c.val, c.currency, got, c.want)
		}
	}
}

func TestFmtInt(t *testing.T) {
	cases := map[int64]string{
		0:          "0",
		999:        "999",
		1000:       "1,000",
		1234567:    "1,234,567",
		-1234567:   "-1,234,567",
		1000000000: "1,000,000,000",
	}
	for in, want := range cases {
		if got := fmtInt(in); got != want {
			t.Errorf("fmtInt(%d) = %q, want %q", in, got, want)
		}
	}
}

func TestFormatBulkSpark_NullClose(t *testing.T) {
	v1, v3 := 178.5, 178.72
	results := yahoo.SparkResponse{
		"AAPL": {
			Symbol:             "AAPL",
			Timestamps:         []int64{1700000000, 1700003600, 1700007200},
			Close:              []*float64{&v1, nil, &v3},
			ChartPreviousClose: 176.57,
			DataGranularity:    3600,
		},
	}

	out := formatBulkSpark([]string{"AAPL", "MISSING"}, results)

	if !strings.Contains(out, "N/A") {
		t.Errorf("expected N/A for null close, got:\n%s", out)
	}
	if strings.Contains(out, "-100.00%") {
		t.Errorf("null close must not render as -100%%, got:\n%s", out)
	}
	if !strings.Contains(out, "178.50") || !strings.Contains(out, "178.72") {
		t.Errorf("expected real close values in output, got:\n%s", out)
	}
	if !strings.Contains(out, "MISSING: no data") {
		t.Errorf("expected missing-symbol marker, got:\n%s", out)
	}
	if !strings.Contains(out, "interval: 1h") {
		t.Errorf("expected 1h granularity label, got:\n%s", out)
	}
}

func TestFormatBulkSpark_ShortCloseArray(t *testing.T) {
	v := 10.0
	results := yahoo.SparkResponse{
		"X": {
			Symbol:             "X",
			Timestamps:         []int64{1700000000, 1700003600},
			Close:              []*float64{&v},
			ChartPreviousClose: 9,
			DataGranularity:    3600,
		},
	}
	// Must not panic when close is shorter than timestamps.
	out := formatBulkSpark([]string{"X"}, results)
	if !strings.Contains(out, "N/A") {
		t.Errorf("expected N/A for the missing trailing close, got:\n%s", out)
	}
}

func TestChartLocation(t *testing.T) {
	ny := chartLocation(yahoo.ChartMeta{ExchangeTimezoneName: "America/New_York", Timezone: "EST", GMTOffset: -18000})
	if ny.String() != "America/New_York" {
		t.Errorf("named zone = %q, want America/New_York", ny)
	}

	fixed := chartLocation(yahoo.ChartMeta{ExchangeTimezoneName: "Not/AZone", Timezone: "CEST", GMTOffset: 7200})
	_, off := time.Unix(0, 0).In(fixed).Zone()
	if fixed.String() != "CEST" || off != 7200 {
		t.Errorf("fallback zone = %q/%d, want CEST/7200", fixed, off)
	}

	offsetOnly := chartLocation(yahoo.ChartMeta{GMTOffset: -14400})
	if offsetOnly.String() != "UTC-4" {
		t.Errorf("offset-only zone = %q, want UTC-4", offsetOnly)
	}

	if chartLocation(yahoo.ChartMeta{}) != time.UTC {
		t.Error("empty meta should fall back to UTC")
	}
}

func TestFormatChart_UsesExchangeTimezone(t *testing.T) {
	// 2023-11-15 14:30:00 UTC == 09:30 in New York (EST).
	o, h, l, c := 175.0, 177.0, 174.0, 176.5
	var v int64 = 50_000_000
	chart := &yahoo.ChartResult{
		Meta: yahoo.ChartMeta{
			Symbol: "AAPL", ExchangeName: "NMS", Currency: "USD", Range: "1d", DataGranularity: "1h",
			ExchangeTimezoneName: "America/New_York", Timezone: "EST", GMTOffset: -18000,
		},
		Timestamps: []int64{1700058600},
		Indicators: yahoo.ChartIndicators{Quote: []yahoo.ChartQuote{{
			Open: []*float64{&o}, High: []*float64{&h}, Low: []*float64{&l}, Close: []*float64{&c}, Volume: []*int64{&v},
		}}},
	}

	out := formatChart(chart)
	if !strings.Contains(out, "Timezone: America/New_York") {
		t.Errorf("expected timezone header, got:\n%s", out)
	}
	if !strings.Contains(out, "2023-11-15 09:30") {
		t.Errorf("expected row in exchange-local time, got:\n%s", out)
	}
	if !strings.Contains(out, "50,000,000") {
		t.Errorf("expected formatted volume, got:\n%s", out)
	}
	if !strings.Contains(out, "Showing 1 of 1 data points") {
		t.Errorf("expected row count footer, got:\n%s", out)
	}
}

func TestFormatChart_NoData(t *testing.T) {
	out := formatChart(&yahoo.ChartResult{Meta: yahoo.ChartMeta{Symbol: "X"}})
	if !strings.Contains(out, "No data points available") {
		t.Errorf("expected empty marker, got:\n%s", out)
	}
}

func TestFormatOptions_ExpirationDatesAreUTC(t *testing.T) {
	// Yahoo expirations are midnight UTC; 1700006400 is 2023-11-15 00:00:00 UTC.
	// In any zone west of UTC a local-time format would say 2023-11-14.
	result := &yahoo.OptionsResult{
		UnderlyingSymbol: "AAPL",
		Quote:            yahoo.OptionsQuote{Currency: "USD", RegularMarketPrice: 178.72},
		ExpirationDates:  []int64{1700006400, 1700611200},
		Options: []yahoo.OptionsChain{{
			ExpirationDate: 1700006400,
			Calls:          []yahoo.OptionContract{{Strike: 175, LastPrice: 5.2, Bid: 5.1, Ask: 5.3, Volume: 100, OpenInterest: 2000, ImpliedVolatility: 0.25, InTheMoney: true}},
			Puts:           []yahoo.OptionContract{{Strike: 180, LastPrice: 4.1, Bid: 4.0, Ask: 4.2, Volume: 50, OpenInterest: 900, ImpliedVolatility: 0.3}},
		}},
	}

	out := formatOptions(result)
	if !strings.Contains(out, "2023-11-15") || strings.Contains(out, "2023-11-14") {
		t.Errorf("expiration must render as the UTC date, got:\n%s", out)
	}
	if !strings.Contains(out, "Underlying Price: $178.72") {
		t.Errorf("expected underlying price, got:\n%s", out)
	}
	if !strings.Contains(out, "175.00*") {
		t.Errorf("expected ITM marker on call strike, got:\n%s", out)
	}
	if !strings.Contains(out, "25.0%") {
		t.Errorf("expected IV as percent, got:\n%s", out)
	}
}

func TestFormatOptions_NoChain(t *testing.T) {
	out := formatOptions(&yahoo.OptionsResult{UnderlyingSymbol: "X"})
	if !strings.Contains(out, "No options data available") {
		t.Errorf("expected empty marker, got:\n%s", out)
	}
}

func TestFormatNews_UTC(t *testing.T) {
	news := []yahoo.SearchNews{
		{Title: "First", Publisher: "Pub A", Link: "https://a.example", ProviderPublishTime: 1700058600},
		{Title: "Second", Publisher: "Pub B", ProviderPublishTime: 1700058600},
	}
	out := formatNews("AAPL", news)
	if !strings.Contains(out, "2023-11-15 14:30 UTC") {
		t.Errorf("expected UTC timestamp, got:\n%s", out)
	}
	if !strings.Contains(out, "Link: https://a.example") {
		t.Errorf("expected link line, got:\n%s", out)
	}
	if strings.Count(out, "Link:") != 1 {
		t.Errorf("link line should be omitted when empty, got:\n%s", out)
	}
	if !strings.Contains(formatNews("X", nil), "No recent news found") {
		t.Error("expected empty marker")
	}
}

func TestFormatBulkSpark_UTCHeader(t *testing.T) {
	v := 1.0
	out := formatBulkSpark([]string{"X"}, yahoo.SparkResponse{"X": {
		Symbol: "X", Timestamps: []int64{1700058600}, Close: []*float64{&v}, DataGranularity: 86400,
	}})
	if !strings.Contains(out, "Date (UTC)") || !strings.Contains(out, "2023-11-15 14:30") {
		t.Errorf("expected UTC-labelled rows, got:\n%s", out)
	}
	if !strings.Contains(out, "interval: 1d") {
		t.Errorf("expected 1d granularity label, got:\n%s", out)
	}
}
