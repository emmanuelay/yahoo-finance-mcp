package yahoo

// QuoteSummaryResponse is the top-level response for quoteSummary endpoints.
type QuoteSummaryResponse struct {
	QuoteSummary struct {
		Result []QuoteSummaryResult `json:"result"`
		Error  *YahooError          `json:"error"`
	} `json:"quoteSummary"`
}

type QuoteSummaryResult struct {
	Price               *PriceData               `json:"price"`
	SummaryDetail       *SummaryDetailData       `json:"summaryDetail"`
	AssetProfile        *AssetProfileData        `json:"assetProfile"`
	QuoteType           *QuoteTypeData           `json:"quoteType"`
	RecommendationTrend *RecommendationTrendData `json:"recommendationTrend"`
}

type YahooError struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

// YahooValue represents a Yahoo Finance value wrapper (raw + fmt).
type YahooValue struct {
	Raw float64 `json:"raw"`
	Fmt string  `json:"fmt"`
}

type YahooLongValue struct {
	Raw int64  `json:"raw"`
	Fmt string `json:"fmt"`
}

// PriceData from quoteSummary price module.
type PriceData struct {
	Symbol                     string     `json:"symbol"`
	ShortName                  string     `json:"shortName"`
	LongName                   string     `json:"longName"`
	Currency                   string     `json:"currency"`
	Exchange                   string     `json:"exchange"`
	ExchangeName               string     `json:"exchangeName"`
	QuoteType                  string     `json:"quoteType"`
	MarketState                string     `json:"marketState"`
	RegularMarketPrice         YahooValue `json:"regularMarketPrice"`
	RegularMarketChange        YahooValue `json:"regularMarketChange"`
	RegularMarketChangePercent YahooValue `json:"regularMarketChangePercent"`
	RegularMarketVolume        YahooLongValue `json:"regularMarketVolume"`
	RegularMarketOpen          YahooValue `json:"regularMarketOpen"`
	RegularMarketDayHigh       YahooValue `json:"regularMarketDayHigh"`
	RegularMarketDayLow        YahooValue `json:"regularMarketDayLow"`
	RegularMarketPreviousClose YahooValue `json:"regularMarketPreviousClose"`
	MarketCap                  YahooLongValue `json:"marketCap"`
	PreMarketPrice             YahooValue `json:"preMarketPrice"`
	PreMarketChange            YahooValue `json:"preMarketChange"`
	PreMarketChangePercent     YahooValue `json:"preMarketChangePercent"`
	PostMarketPrice            YahooValue `json:"postMarketPrice"`
	PostMarketChange           YahooValue `json:"postMarketChange"`
	PostMarketChangePercent    YahooValue `json:"postMarketChangePercent"`
}

// SummaryDetailData from quoteSummary summaryDetail module.
type SummaryDetailData struct {
	TrailingPE       YahooValue `json:"trailingPE"`
	ForwardPE        YahooValue `json:"forwardPE"`
	DividendYield    YahooValue `json:"dividendYield"`
	DividendRate     YahooValue `json:"dividendRate"`
	ExDividendDate   YahooValue `json:"exDividendDate"`
	FiftyTwoWeekLow  YahooValue `json:"fiftyTwoWeekLow"`
	FiftyTwoWeekHigh YahooValue `json:"fiftyTwoWeekHigh"`
	FiftyDayAverage  YahooValue `json:"fiftyDayAverage"`
	TwoHundredDayAverage YahooValue `json:"twoHundredDayAverage"`
	Beta             YahooValue `json:"beta"`
	TrailingAnnualDividendYield YahooValue `json:"trailingAnnualDividendYield"`
	PayoutRatio      YahooValue `json:"payoutRatio"`
}

// ChartResponse from v8 chart endpoint.
type ChartResponse struct {
	Chart struct {
		Result []ChartResult `json:"result"`
		Error  *YahooError   `json:"error"`
	} `json:"chart"`
}

type ChartResult struct {
	Meta       ChartMeta       `json:"meta"`
	Timestamps []int64         `json:"timestamp"`
	Indicators ChartIndicators `json:"indicators"`
}

type ChartMeta struct {
	Currency             string  `json:"currency"`
	Symbol               string  `json:"symbol"`
	ExchangeName         string  `json:"exchangeName"`
	InstrumentType       string  `json:"instrumentType"`
	RegularMarketPrice   float64 `json:"regularMarketPrice"`
	PreviousClose        float64 `json:"previousClose"`
	ChartPreviousClose   float64 `json:"chartPreviousClose"`
	DataGranularity      string  `json:"dataGranularity"`
	Range                string  `json:"range"`
	ValidRanges          []string `json:"validRanges"`
}

type ChartIndicators struct {
	Quote    []ChartQuote    `json:"quote"`
	AdjClose []ChartAdjClose `json:"adjclose"`
}

type ChartQuote struct {
	Open   []*float64 `json:"open"`
	High   []*float64 `json:"high"`
	Low    []*float64 `json:"low"`
	Close  []*float64 `json:"close"`
	Volume []*int64   `json:"volume"`
}

type ChartAdjClose struct {
	AdjClose []*float64 `json:"adjclose"`
}

// SearchResponse from v1 finance/search.
type SearchResponse struct {
	Quotes []SearchQuote `json:"quotes"`
	News   []SearchNews  `json:"news"`
	Count  int           `json:"count"`
}

type SearchQuote struct {
	Symbol    string `json:"symbol"`
	ShortName string `json:"shortname"`
	LongName  string `json:"longname"`
	Exchange  string `json:"exchange"`
	QuoteType string `json:"quoteType"`
	Industry  string `json:"industry"`
	Sector    string `json:"sector"`
	Score     float64 `json:"score"`
}

type SearchNews struct {
	UUID          string `json:"uuid"`
	Title         string `json:"title"`
	Publisher     string `json:"publisher"`
	Link          string `json:"link"`
	ProviderPublishTime int64  `json:"providerPublishTime"`
}

// AssetProfileData from quoteSummary assetProfile module.
type AssetProfileData struct {
	Address1            string              `json:"address1"`
	Address2            string              `json:"address2"`
	City                string              `json:"city"`
	State               string              `json:"state"`
	Zip                 string              `json:"zip"`
	Country             string              `json:"country"`
	Phone               string              `json:"phone"`
	Website             string              `json:"website"`
	Industry            string              `json:"industry"`
	IndustryKey         string              `json:"industryKey"`
	Sector              string              `json:"sector"`
	SectorKey           string              `json:"sectorKey"`
	LongBusinessSummary string              `json:"longBusinessSummary"`
	FullTimeEmployees   int                 `json:"fullTimeEmployees"`
	CompanyOfficers     []CompanyOfficer    `json:"companyOfficers"`
}

type CompanyOfficer struct {
	Name         string     `json:"name"`
	Title        string     `json:"title"`
	Age          int        `json:"age"`
	TotalPay     YahooLongValue `json:"totalPay"`
	YearBorn     int        `json:"yearBorn"`
}

type QuoteTypeData struct {
	Symbol    string `json:"symbol"`
	ShortName string `json:"shortName"`
	LongName  string `json:"longName"`
	QuoteType string `json:"quoteType"`
	Exchange  string `json:"exchange"`
}

// RecommendationTrendData from quoteSummary recommendationTrend module.
type RecommendationTrendData struct {
	Trend []RecommendationTrend `json:"trend"`
}

type RecommendationTrend struct {
	Period     string `json:"period"`
	StrongBuy  int    `json:"strongBuy"`
	Buy        int    `json:"buy"`
	Hold       int    `json:"hold"`
	Sell       int    `json:"sell"`
	StrongSell int    `json:"strongSell"`
}

// OptionsResponse from v7 finance/options.
type OptionsResponse struct {
	OptionChain struct {
		Result []OptionsResult `json:"result"`
		Error  *YahooError     `json:"error"`
	} `json:"optionChain"`
}

type OptionsResult struct {
	UnderlyingSymbol string          `json:"underlyingSymbol"`
	ExpirationDates  []int64         `json:"expirationDates"`
	Strikes          []float64       `json:"strikes"`
	Quote            OptionsQuote    `json:"quote"`
	Options          []OptionsChain  `json:"options"`
}

type OptionsQuote struct {
	Symbol             string  `json:"symbol"`
	RegularMarketPrice float64 `json:"regularMarketPrice"`
}

type OptionsChain struct {
	ExpirationDate int64          `json:"expirationDate"`
	Calls          []OptionContract `json:"calls"`
	Puts           []OptionContract `json:"puts"`
}

type OptionContract struct {
	ContractSymbol    string  `json:"contractSymbol"`
	Strike            float64 `json:"strike"`
	Currency          string  `json:"currency"`
	LastPrice         float64 `json:"lastPrice"`
	Change            float64 `json:"change"`
	PercentChange     float64 `json:"percentChange"`
	Volume            int     `json:"volume"`
	OpenInterest      int     `json:"openInterest"`
	Bid               float64 `json:"bid"`
	Ask               float64 `json:"ask"`
	ImpliedVolatility float64 `json:"impliedVolatility"`
	InTheMoney        bool    `json:"inTheMoney"`
	Expiration        int64   `json:"expiration"`
}

// TimeseriesResponse from fundamentals-timeseries endpoint.
type TimeseriesResponse struct {
	Timeseries struct {
		Result []TimeseriesResult `json:"result"`
		Error  *YahooError        `json:"error"`
	} `json:"timeseries"`
}

type TimeseriesResult struct {
	Meta       TimeseriesMeta          `json:"meta"`
	Timestamp  []int64                 `json:"timestamp"`
	Type       string                  // populated from the key name
	DataPoints []map[string]interface{} // populated from dynamic keys
}

type TimeseriesMeta struct {
	Symbol []string `json:"symbol"`
	Type   []string `json:"type"`
}
