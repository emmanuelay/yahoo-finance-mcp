package yahoo

import (
	"fmt"
	"net/url"
)

// GetRecommendations fetches analyst recommendation trends for a symbol.
func (c *Client) GetRecommendations(symbol string) (*RecommendationTrendData, error) {
	params := url.Values{
		"modules": {"recommendationTrend"},
	}

	var resp QuoteSummaryResponse
	path := fmt.Sprintf("/v10/finance/quoteSummary/%s", url.PathEscape(symbol))
	if err := c.GetJSON(path, params, true, &resp); err != nil {
		return nil, fmt.Errorf("get recommendations: %w", err)
	}

	if resp.QuoteSummary.Error != nil {
		return nil, fmt.Errorf("yahoo error: %s", resp.QuoteSummary.Error.Description)
	}

	if len(resp.QuoteSummary.Result) == 0 {
		return nil, fmt.Errorf("no data found for symbol %q", symbol)
	}

	return resp.QuoteSummary.Result[0].RecommendationTrend, nil
}
