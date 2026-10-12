package yahoo

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// GetNews fetches recent news articles for a symbol.
func (c *Client) GetNews(ctx context.Context, symbol string, count int) ([]SearchNews, error) {
	if count <= 0 {
		count = 5
	}

	params := url.Values{
		"q":           {symbol},
		"quotesCount": {"0"},
		"newsCount":   {strconv.Itoa(count)},
	}

	var resp SearchResponse
	if err := c.GetJSON(ctx, "/v1/finance/search", params, false, &resp); err != nil {
		return nil, fmt.Errorf("get news: %w", err)
	}

	return resp.News, nil
}
