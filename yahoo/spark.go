package yahoo

import (
	"fmt"
	"net/url"
	"strings"
)

// maxBulkSymbols is the maximum number of symbols allowed per bulk API call.
const maxBulkSymbols = 50

// SparkResponse from v8 finance/spark endpoint.
type SparkResponse map[string]SparkResult

// SparkResult contains spark data for a single symbol.
type SparkResult struct {
	Symbol             string    `json:"symbol"`
	Timestamps         []int64   `json:"timestamp"`
	Close              []float64 `json:"close"`
	ChartPreviousClose float64   `json:"chartPreviousClose"`
	DataGranularity    int       `json:"dataGranularity"`
}

// GetBulkSpark fetches simplified chart data for multiple symbols in a single call.
func (c *Client) GetBulkSpark(symbols []string, rangeStr, interval string) (SparkResponse, error) {
	if len(symbols) == 0 {
		return nil, fmt.Errorf("at least one symbol is required")
	}
	if len(symbols) > maxBulkSymbols {
		return nil, fmt.Errorf("too many symbols: %d (max %d)", len(symbols), maxBulkSymbols)
	}
	if rangeStr == "" {
		rangeStr = "1mo"
	}
	if interval == "" {
		interval = "1d"
	}

	params := url.Values{
		"symbols":  {strings.Join(symbols, ",")},
		"range":    {rangeStr},
		"interval": {interval},
	}

	var resp SparkResponse
	if err := c.GetJSON("/v8/finance/spark", params, false, &resp); err != nil {
		return nil, fmt.Errorf("get bulk spark: %w", err)
	}

	return resp, nil
}
