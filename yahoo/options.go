package yahoo

import (
	"fmt"
	"net/url"
)

// GetOptions fetches the options chain for a symbol.
// If expiration is empty, returns the nearest expiration.
func (c *Client) GetOptions(symbol, expiration string) (*OptionsResult, error) {
	params := url.Values{}
	if expiration != "" {
		params.Set("date", expiration)
	}

	var resp OptionsResponse
	path := fmt.Sprintf("/v7/finance/options/%s", url.PathEscape(symbol))
	if err := c.GetJSON(path, params, true, &resp); err != nil {
		return nil, fmt.Errorf("get options: %w", err)
	}

	if resp.OptionChain.Error != nil {
		return nil, fmt.Errorf("yahoo error: %s", resp.OptionChain.Error.Description)
	}

	if len(resp.OptionChain.Result) == 0 {
		return nil, fmt.Errorf("no options data found for symbol %q", symbol)
	}

	return &resp.OptionChain.Result[0], nil
}
