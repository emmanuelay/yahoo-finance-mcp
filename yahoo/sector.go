package yahoo

import (
	"fmt"
	"net/url"
)

const query1BaseURL = "https://query1.finance.yahoo.com/v1/finance"

// SectorResponse from the sectors endpoint.
type SectorResponse struct {
	Data SectorData `json:"data"`
}

type SectorData struct {
	Name         string              `json:"name"`
	Symbol       string              `json:"symbol"`
	Overview     SectorOverview      `json:"overview"`
	TopCompanies []TopCompany        `json:"topCompanies"`
	TopETFs      []TopFund           `json:"topETFs"`
	TopMutualFunds []TopFund         `json:"topMutualFunds"`
	Industries   []IndustryListItem  `json:"industries"`
}

type SectorOverview struct {
	CompaniesCount  int        `json:"companiesCount"`
	MarketCap       RawFmt     `json:"marketCap"`
	Description     string     `json:"description"`
	IndustriesCount int        `json:"industriesCount"`
	MarketWeight    RawFmt     `json:"marketWeight"`
	EmployeeCount   RawFmt     `json:"employeeCount"`
}

type RawFmt struct {
	Raw float64 `json:"raw"`
	Fmt string  `json:"fmt"`
}

type TopCompany struct {
	Symbol       string `json:"symbol"`
	Name         string `json:"name"`
	Rating       any    `json:"rating"`
	MarketWeight RawFmt `json:"marketWeight"`
}

type TopFund struct {
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
}

type IndustryListItem struct {
	Key          string `json:"key"`
	Name         string `json:"name"`
	Symbol       string `json:"symbol"`
	MarketWeight RawFmt `json:"marketWeight"`
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
func (c *Client) GetSector(key string) (*SectorData, error) {
	params := url.Values{
		"formatted":   {"true"},
		"withReturns": {"true"},
		"lang":        {"en-US"},
		"region":      {"US"},
	}

	apiURL := fmt.Sprintf("%s/sectors/%s", query1BaseURL, url.PathEscape(key))

	var resp SectorResponse
	if err := c.GetAbsoluteJSON(apiURL, params, &resp); err != nil {
		return nil, fmt.Errorf("get sector %q: %w", key, err)
	}

	return &resp.Data, nil
}
