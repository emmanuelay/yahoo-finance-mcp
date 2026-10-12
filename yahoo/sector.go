package yahoo

import (
	"context"
	"fmt"
	"net/url"
)

const query1BaseURL = "https://query1.finance.yahoo.com/v1/finance"

// SectorResponse from the sectors endpoint.
type SectorResponse struct {
	Data SectorData `json:"data"`
}

type SectorData struct {
	Name           string             `json:"name"`
	Symbol         string             `json:"symbol"`
	Overview       SectorOverview     `json:"overview"`
	TopCompanies   []TopCompany       `json:"topCompanies"`
	TopETFs        []TopFund          `json:"topETFs"`
	TopMutualFunds []TopFund          `json:"topMutualFunds"`
	Industries     []IndustryListItem `json:"industries"`
}

type SectorOverview struct {
	CompaniesCount  int        `json:"companiesCount"`
	MarketCap       YahooValue `json:"marketCap"`
	Description     string     `json:"description"`
	IndustriesCount int        `json:"industriesCount"`
	MarketWeight    YahooValue `json:"marketWeight"`
	EmployeeCount   YahooValue `json:"employeeCount"`
}

type TopCompany struct {
	Symbol       string     `json:"symbol"`
	Name         string     `json:"name"`
	Rating       any        `json:"rating"`
	MarketWeight YahooValue `json:"marketWeight"`
}

type TopFund struct {
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
}

type IndustryListItem struct {
	Key          string     `json:"key"`
	Name         string     `json:"name"`
	Symbol       string     `json:"symbol"`
	MarketWeight YahooValue `json:"marketWeight"`
}

// ValidSectorKeys lists all recognized sector keys.
var ValidSectorKeys = []string{
	"basic-materials",
	"communication-services",
	"consumer-cyclical",
	"consumer-defensive",
	"energy",
	"financial-services",
	"healthcare",
	"industrials",
	"real-estate",
	"technology",
	"utilities",
}

// GetSector fetches sector data by key (e.g., "technology", "healthcare").
func (c *Client) GetSector(ctx context.Context, key string) (*SectorData, error) {
	params := url.Values{
		"formatted":   {"true"},
		"withReturns": {"true"},
		"lang":        {"en-US"},
		"region":      {"US"},
	}

	apiURL := fmt.Sprintf("%s/sectors/%s", query1BaseURL, url.PathEscape(key))

	var resp SectorResponse
	if err := c.GetAbsoluteJSON(ctx, apiURL, params, true, &resp); err != nil {
		return nil, fmt.Errorf("get sector %q: %w", key, err)
	}

	return &resp.Data, nil
}
