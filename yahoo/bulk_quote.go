package yahoo

import (
	"fmt"
	"net/url"
	"strings"
)

// BulkQuoteResponse from v7 finance/quote endpoint.
type BulkQuoteResponse struct {
	QuoteResponse struct {
		Result []BulkQuoteResult `json:"result"`
		Error  *YahooError       `json:"error"`
	} `json:"quoteResponse"`
}

// BulkQuoteResult contains quote data for a single symbol from the bulk endpoint.
type BulkQuoteResult struct {
	Symbol                     string  `json:"symbol"`
	ShortName                  string  `json:"shortName"`
	LongName                   string  `json:"longName"`
	Currency                   string  `json:"currency"`
	Exchange                   string  `json:"exchange"`
	FullExchangeName           string  `json:"fullExchangeName"`
	QuoteType                  string  `json:"quoteType"`
	MarketState                string  `json:"marketState"`
	RegularMarketPrice         float64 `json:"regularMarketPrice"`
	RegularMarketChange        float64 `json:"regularMarketChange"`
	RegularMarketChangePercent float64 `json:"regularMarketChangePercent"`
	RegularMarketVolume        int64   `json:"regularMarketVolume"`
	RegularMarketOpen          float64 `json:"regularMarketOpen"`
	RegularMarketDayHigh       float64 `json:"regularMarketDayHigh"`
	RegularMarketDayLow        float64 `json:"regularMarketDayLow"`
	RegularMarketPreviousClose float64 `json:"regularMarketPreviousClose"`
	MarketCap                  int64   `json:"marketCap"`
	TrailingPE                 float64 `json:"trailingPE"`
	ForwardPE                  float64 `json:"forwardPE"`
	FiftyTwoWeekLow            float64 `json:"fiftyTwoWeekLow"`
	FiftyTwoWeekHigh           float64 `json:"fiftyTwoWeekHigh"`
	FiftyDayAverage            float64 `json:"fiftyDayAverage"`
	TwoHundredDayAverage       float64 `json:"twoHundredDayAverage"`
	TrailingAnnualDividendYield float64 `json:"trailingAnnualDividendYield"`
}

// GetBulkQuotes fetches quotes for multiple symbols in a single API call.
func (c *Client) GetBulkQuotes(symbols []string) ([]BulkQuoteResult, error) {
	if len(symbols) == 0 {
		return nil, fmt.Errorf("at least one symbol is required")
	}
	if len(symbols) > maxBulkSymbols {
		return nil, fmt.Errorf("too many symbols: %d (max %d)", len(symbols), maxBulkSymbols)
	}

	params := url.Values{
		"symbols": {strings.Join(symbols, ",")},
	}

	var resp BulkQuoteResponse
	if err := c.GetJSON("/v7/finance/quote", params, true, &resp); err != nil {
		return nil, fmt.Errorf("get bulk quotes: %w", err)
	}

	if resp.QuoteResponse.Error != nil {
		return nil, fmt.Errorf("yahoo error: %s", resp.QuoteResponse.Error.Description)
	}

	return resp.QuoteResponse.Result, nil
}
