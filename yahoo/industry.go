package yahoo

import (
	"context"
	"fmt"
	"net/url"
)

// IndustryResponse from the industries endpoint.
type IndustryResponse struct {
	Data IndustryData `json:"data"`
}

type IndustryData struct {
	SectorKey              string              `json:"sectorKey"`
	SectorName             string              `json:"sectorName"`
	Name                   string              `json:"name"`
	Symbol                 string              `json:"symbol"`
	Overview               SectorOverview      `json:"overview"`
	TopCompanies           []TopCompany        `json:"topCompanies"`
	TopPerformingCompanies []PerformingCompany `json:"topPerformingCompanies"`
	TopGrowthCompanies     []GrowthCompany     `json:"topGrowthCompanies"`
}

type PerformingCompany struct {
	Symbol      string     `json:"symbol"`
	Name        string     `json:"name"`
	YtdReturn   YahooValue `json:"ytdReturn"`
	LastPrice   YahooValue `json:"lastPrice"`
	TargetPrice YahooValue `json:"targetPrice"`
}

type GrowthCompany struct {
	Symbol         string     `json:"symbol"`
	Name           string     `json:"name"`
	YtdReturn      YahooValue `json:"ytdReturn"`
	GrowthEstimate YahooValue `json:"growthEstimate"`
}

// GetIndustry fetches industry data by key (e.g., "consumer-electronics", "semiconductors").
func (c *Client) GetIndustry(ctx context.Context, key string) (*IndustryData, error) {
	params := url.Values{
		"formatted":   {"true"},
		"withReturns": {"true"},
		"lang":        {"en-US"},
		"region":      {"US"},
	}

	apiURL := fmt.Sprintf("%s/industries/%s", query1BaseURL, url.PathEscape(key))

	var resp IndustryResponse
	if err := c.GetAbsoluteJSON(ctx, apiURL, params, true, &resp); err != nil {
		return nil, fmt.Errorf("get industry %q: %w", key, err)
	}

	return &resp.Data, nil
}
