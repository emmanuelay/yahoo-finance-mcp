package yahoo

import (
	"context"
	"fmt"
	"net/url"
)

const marketBaseURL = "https://query1.finance.yahoo.com/v6/finance"

// MarketSummaryResponse from the marketSummary endpoint.
type MarketSummaryResponse struct {
	MarketSummaryResponse struct {
		Result []MarketSummaryItem `json:"result"`
		Error  *YahooError         `json:"error"`
	} `json:"marketSummaryResponse"`
}

type MarketSummaryItem struct {
	Symbol                     string  `json:"symbol"`
	ShortName                  string  `json:"shortName"`
	Exchange                   string  `json:"exchange"`
	MarketState                string  `json:"marketState"`
	RegularMarketPrice         float64 `json:"regularMarketPrice"`
	RegularMarketChange        float64 `json:"regularMarketChange"`
	RegularMarketChangePercent float64 `json:"regularMarketChangePercent"`
}

// MarketStatusResponse from the markettime endpoint.
type MarketStatusResponse struct {
	Finance struct {
		MarketTimes []MarketTimeGroup `json:"marketTimes"`
		Error       *YahooError       `json:"error"`
	} `json:"finance"`
}

type MarketTimeGroup struct {
	MarketTime []MarketTime `json:"marketTime"`
}

type MarketTime struct {
	ID       string           `json:"id"`
	Name     string           `json:"name"`
	Status   string           `json:"status"`
	Message  string           `json:"message"`
	Open     string           `json:"open"`
	Close    string           `json:"close"`
	Time     string           `json:"time"`
	Timezone []MarketTimezone `json:"timezone"`
}

type MarketTimezone struct {
	GMTOffset string `json:"gmtoffset"`
	Short     string `json:"short"`
}

// ValidMarketKeys lists all recognized market keys.
var ValidMarketKeys = []string{
	"US", "GB", "ASIA", "EUROPE",
	"RATES", "COMMODITIES", "CURRENCIES", "CRYPTOCURRENCIES",
}

// GetMarketSummary fetches market summary (indices/benchmarks) for a market.
func (c *Client) GetMarketSummary(ctx context.Context, market string) ([]MarketSummaryItem, error) {
	params := url.Values{
		"fields":    {"shortName,regularMarketPrice,regularMarketChange,regularMarketChangePercent"},
		"formatted": {"false"},
		"lang":      {"en-US"},
		"market":    {market},
	}

	var resp MarketSummaryResponse
	if err := c.GetAbsoluteJSON(ctx, marketBaseURL+"/quote/marketSummary", params, false, &resp); err != nil {
		return nil, fmt.Errorf("get market summary %q: %w", market, err)
	}

	if resp.MarketSummaryResponse.Error != nil {
		return nil, fmt.Errorf("yahoo error: %s", resp.MarketSummaryResponse.Error.Description)
	}

	return resp.MarketSummaryResponse.Result, nil
}

// GetMarketStatus fetches market open/close times and timezone for a market.
func (c *Client) GetMarketStatus(ctx context.Context, market string) ([]MarketTimeGroup, error) {
	params := url.Values{
		"formatted": {"true"},
		"key":       {"finance"},
		"lang":      {"en-US"},
		"market":    {market},
	}

	var resp MarketStatusResponse
	if err := c.GetAbsoluteJSON(ctx, marketBaseURL+"/markettime", params, false, &resp); err != nil {
		return nil, fmt.Errorf("get market status %q: %w", market, err)
	}

	if resp.Finance.Error != nil {
		return nil, fmt.Errorf("yahoo error: %s", resp.Finance.Error.Description)
	}

	return resp.Finance.MarketTimes, nil
}
