package tools

import (
	"strings"
	"testing"

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
