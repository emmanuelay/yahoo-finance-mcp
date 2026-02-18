package yahoo

import (
	"fmt"
	"net/url"
)

// GetProfile fetches company profile information for a symbol.
func (c *Client) GetProfile(symbol string) (*AssetProfileData, *QuoteTypeData, error) {
	params := url.Values{
		"modules": {"assetProfile,quoteType"},
	}

	var resp QuoteSummaryResponse
	path := fmt.Sprintf("/v10/finance/quoteSummary/%s", url.PathEscape(symbol))
	if err := c.GetJSON(path, params, true, &resp); err != nil {
		return nil, nil, fmt.Errorf("get profile: %w", err)
	}

	if resp.QuoteSummary.Error != nil {
		return nil, nil, fmt.Errorf("yahoo error: %s", resp.QuoteSummary.Error.Description)
	}

	if len(resp.QuoteSummary.Result) == 0 {
		return nil, nil, fmt.Errorf("no data found for symbol %q", symbol)
	}

	result := resp.QuoteSummary.Result[0]
	return result.AssetProfile, result.QuoteType, nil
}
