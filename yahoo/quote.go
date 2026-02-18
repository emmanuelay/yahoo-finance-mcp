package yahoo

import (
	"fmt"
	"net/url"
)

// GetQuote fetches real-time price and summary details for a symbol.
func (c *Client) GetQuote(symbol string) (*PriceData, *SummaryDetailData, error) {
	params := url.Values{
		"modules": {"price,summaryDetail"},
	}

	var resp QuoteSummaryResponse
	path := fmt.Sprintf("/v10/finance/quoteSummary/%s", url.PathEscape(symbol))
	if err := c.GetJSON(path, params, true, &resp); err != nil {
		return nil, nil, fmt.Errorf("get quote: %w", err)
	}

	if resp.QuoteSummary.Error != nil {
		return nil, nil, fmt.Errorf("yahoo error: %s", resp.QuoteSummary.Error.Description)
	}

	if len(resp.QuoteSummary.Result) == 0 {
		return nil, nil, fmt.Errorf("no data found for symbol %q", symbol)
	}

	result := resp.QuoteSummary.Result[0]
	return result.Price, result.SummaryDetail, nil
}
