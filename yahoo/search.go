package yahoo

import (
	"fmt"
	"net/url"
	"strconv"
)

// Search finds symbols and companies matching the query.
func (c *Client) Search(query string, limit int) (*SearchResponse, error) {
	if limit <= 0 {
		limit = 10
	}

	params := url.Values{
		"q":           {query},
		"quotesCount": {strconv.Itoa(limit)},
		"newsCount":   {"0"},
		"enableFuzzyQuery": {"false"},
		"quotesQueryId": {"tss_match_phrase_query"},
	}

	var resp SearchResponse
	if err := c.GetJSON("/v1/finance/search", params, false, &resp); err != nil {
		return nil, fmt.Errorf("search: %w", err)
	}

	return &resp, nil
}
