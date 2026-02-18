package yahoo

import (
	"fmt"
	"net/url"
)

// GetChart fetches historical OHLCV chart data for a symbol.
func (c *Client) GetChart(symbol, rangeStr, interval string) (*ChartResult, error) {
	if rangeStr == "" {
		rangeStr = "1mo"
	}
	if interval == "" {
		interval = "1d"
	}

	params := url.Values{
		"range":    {rangeStr},
		"interval": {interval},
	}

	var resp ChartResponse
	path := fmt.Sprintf("/v8/finance/chart/%s", url.PathEscape(symbol))
	if err := c.GetJSON(path, params, false, &resp); err != nil {
		return nil, fmt.Errorf("get chart: %w", err)
	}

	if resp.Chart.Error != nil {
		return nil, fmt.Errorf("yahoo error: %s", resp.Chart.Error.Description)
	}

	if len(resp.Chart.Result) == 0 {
		return nil, fmt.Errorf("no chart data found for symbol %q", symbol)
	}

	return &resp.Chart.Result[0], nil
}
