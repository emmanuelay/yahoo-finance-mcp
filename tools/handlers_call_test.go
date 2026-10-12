package tools

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/emmanuelay/yahoo-finance-mcp/yahoo"
	"github.com/mark3labs/mcp-go/mcp"
)

// fakeYahoo routes requests to canned JSON by URL substring. Auth endpoints
// are always served so the client's cookie/crumb flow succeeds.
type fakeYahoo struct {
	routes map[string]string // URL substring -> JSON body
	status int               // status for matched routes; 0 means 200
	calls  []string
}

func (f *fakeYahoo) RoundTrip(req *http.Request) (*http.Response, error) {
	u := req.URL.String()
	f.calls = append(f.calls, u)

	if req.URL.Host == "fc.yahoo.com" {
		return resp(404, ""), nil
	}
	if strings.Contains(u, "/v1/test/getcrumb") {
		return resp(200, "fake-crumb"), nil
	}
	for sub, body := range f.routes {
		if strings.Contains(u, sub) {
			st := f.status
			if st == 0 {
				st = 200
			}
			return resp(st, body), nil
		}
	}
	return resp(404, `{"error":"no route"}`), nil
}

func resp(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     http.Header{"Content-Type": []string{"application/json"}},
	}
}

func newHandlers(f *fakeYahoo) *Handlers {
	return NewHandlers(yahoo.NewClientWithHTTP(&http.Client{Transport: f}))
}

func call(args map[string]any) mcp.CallToolRequest {
	var req mcp.CallToolRequest
	req.Params.Arguments = args
	return req
}

// text extracts the single text payload from a tool result.
func text(t *testing.T, res *mcp.CallToolResult) string {
	t.Helper()
	if len(res.Content) != 1 {
		t.Fatalf("expected 1 content item, got %d", len(res.Content))
	}
	tc, ok := res.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("content is %T, want TextContent", res.Content[0])
	}
	return tc.Text
}

func wantOK(t *testing.T, res *mcp.CallToolResult, err error, substrs ...string) string {
	t.Helper()
	if err != nil {
		t.Fatalf("handler returned Go error: %v", err)
	}
	out := text(t, res)
	if res.IsError {
		t.Fatalf("unexpected tool error: %s", out)
	}
	for _, s := range substrs {
		if !strings.Contains(out, s) {
			t.Errorf("output missing %q:\n%s", s, out)
		}
	}
	return out
}

func wantToolError(t *testing.T, res *mcp.CallToolResult, err error, substr string) {
	t.Helper()
	if err != nil {
		t.Fatalf("handler returned Go error: %v", err)
	}
	if !res.IsError {
		t.Fatalf("expected tool error, got success:\n%s", text(t, res))
	}
	if out := text(t, res); !strings.Contains(out, substr) {
		t.Errorf("error %q missing %q", out, substr)
	}
}

// --- fixtures ---

const quoteSummaryJSON = `{"quoteSummary":{"result":[{
	"price":{"symbol":"AAPL","longName":"Apple Inc.","currency":"USD","exchangeName":"NasdaqGS","marketState":"REGULAR",
		"regularMarketPrice":{"raw":178.72},"regularMarketChange":{"raw":-1.28},"regularMarketChangePercent":{"raw":-0.0071},
		"regularMarketVolume":{"raw":54000000},"marketCap":{"raw":2800000000000},
		"regularMarketOpen":{"raw":179.5},"regularMarketDayHigh":{"raw":180.1},"regularMarketDayLow":{"raw":177.9},"regularMarketPreviousClose":{"raw":180.0}},
	"summaryDetail":{"trailingPE":{"raw":29.5},"forwardPE":{"raw":27.1},"fiftyTwoWeekLow":{"raw":124.17},"fiftyTwoWeekHigh":{"raw":199.62},
		"dividendYield":{"raw":0.0055},"beta":{"raw":1.29},"fiftyDayAverage":{"raw":175.2},"twoHundredDayAverage":{"raw":170.8}},
	"assetProfile":{"sector":"Technology","industry":"Consumer Electronics","website":"https://www.apple.com","phone":"408 996 1010",
		"address1":"One Apple Park Way","city":"Cupertino","state":"CA","country":"United States","fullTimeEmployees":161000,
		"longBusinessSummary":"Apple designs things.",
		"companyOfficers":[{"name":"Tim Cook","title":"CEO"},{"name":"A","title":"t"},{"name":"B","title":"t"},{"name":"C","title":"t"},{"name":"D","title":"t"},{"name":"E","title":"t"}]},
	"quoteType":{"symbol":"AAPL","longName":"Apple Inc."},
	"recommendationTrend":{"trend":[{"period":"0m","strongBuy":10,"buy":20,"hold":8,"sell":1,"strongSell":0},{"period":"-1m","strongBuy":9,"buy":21,"hold":8,"sell":1,"strongSell":1}]}
}],"error":null}}`

const yahooErrJSON = `{"quoteSummary":{"result":null,"error":{"code":"Not Found","description":"Quote not found for ticker symbol: NOPE"}}}`

const chartJSON = `{"chart":{"result":[{
	"meta":{"currency":"USD","symbol":"AAPL","exchangeName":"NMS","dataGranularity":"1d","range":"5d",
		"exchangeTimezoneName":"America/New_York","timezone":"EST","gmtoffset":-18000},
	"timestamp":[1700058600,1700145000],
	"indicators":{"quote":[{"open":[175.0,176.5],"high":[177.0,178.0],"low":[174.0,175.5],"close":[176.5,null],"volume":[50000000,45000000]}]}
}],"error":null}}`

const searchJSON = `{"quotes":[
	{"symbol":"AAPL","shortname":"Apple Inc.","longname":"Apple Inc.","exchange":"NMS","quoteType":"EQUITY","sector":"Technology","industry":"Consumer Electronics"},
	{"symbol":"APLE","shortname":"Apple Hospitality","exchange":"NYQ","quoteType":"EQUITY"}
],"news":[{"title":"Apple hits record","publisher":"Reuters","link":"https://r.example/1","providerPublishTime":1700058600}],"count":2}`

const financialsJSON = `{"timeseries":{"result":[
	{"meta":{"symbol":["AAPL"],"type":["annualTotalRevenue"]},"annualTotalRevenue":[
		{"asOfDate":"2022-09-30","reportedValue":{"raw":394328000000,"fmt":"394.33B"},"currencyCode":"USD"},
		{"asOfDate":"2023-09-30","reportedValue":{"raw":383285000000,"fmt":"383.29B"},"currencyCode":"USD"}]},
	{"meta":{"symbol":["AAPL"],"type":["annualNetIncome"]},"annualNetIncome":[
		{"asOfDate":"2023-09-30","reportedValue":{"raw":96995000000,"fmt":"97.00B"},"currencyCode":"USD"}]}
],"error":null}}`

const optionsJSON = `{"optionChain":{"result":[{
	"underlyingSymbol":"AAPL","expirationDates":[1700006400,1700611200,1701216000,1701820800,1702425600,1703030400,1703635200],
	"quote":{"symbol":"AAPL","currency":"USD","regularMarketPrice":178.72},
	"options":[{"expirationDate":1700006400,
		"calls":[{"strike":175,"lastPrice":5.2,"bid":5.1,"ask":5.3,"volume":100,"openInterest":2000,"impliedVolatility":0.25,"inTheMoney":true}],
		"puts":[{"strike":180,"lastPrice":4.1,"bid":4.0,"ask":4.2,"volume":50,"openInterest":900,"impliedVolatility":0.3,"inTheMoney":true}]}]
}],"error":null}}`

const bulkQuoteJSON = `{"quoteResponse":{"result":[
	{"symbol":"AAPL","longName":"Apple Inc.","currency":"USD","regularMarketPrice":178.72,"regularMarketChange":-1.28,"regularMarketChangePercent":-0.71,"regularMarketVolume":54000000,"marketCap":2800000000000},
	{"symbol":"VOLV-B.ST","shortName":"Volvo, AB ser. B with a very long name indeed","currency":"SEK","regularMarketPrice":317,"regularMarketChange":1.6,"regularMarketChangePercent":0.507,"regularMarketVolume":651619,"marketCap":650000000000}
],"error":null}}`

const sparkJSON = `{"AAPL":{"symbol":"AAPL","timestamp":[1700058600,1700062200],"close":[178.5,179.1],"chartPreviousClose":176.57,"dataGranularity":3600}}`

const sectorJSON = `{"data":{"name":"Technology","symbol":"^YH311",
	"overview":{"companiesCount":800,"industriesCount":12,"marketCap":{"raw":15000000000000},"marketWeight":{"raw":0.29},"employeeCount":{"raw":9000000},"description":"Tech sector."},
	"topCompanies":[{"symbol":"AAPL","name":"Apple Inc.","marketWeight":{"raw":0.18}}],
	"topETFs":[{"symbol":"XLK","name":"Technology Select Sector SPDR"}],
	"topMutualFunds":[{"symbol":"FSPTX","name":"Fidelity Select Technology"}],
	"industries":[{"key":"consumer-electronics","name":"Consumer Electronics","symbol":"^YH31101","marketWeight":{"raw":0.2}}]}}`

const industryJSON = `{"data":{"sectorKey":"technology","sectorName":"Technology","name":"Consumer Electronics","symbol":"^YH31101",
	"overview":{"companiesCount":40,"marketCap":{"raw":3000000000000},"marketWeight":{"raw":0.2},"employeeCount":{"raw":500000},"description":"Gadgets."},
	"topCompanies":[{"symbol":"AAPL","name":"Apple Inc.","marketWeight":{"raw":0.9}}],
	"topPerformingCompanies":[{"symbol":"SONO","name":"Sonos Inc. with an unnecessarily long name","ytdReturn":{"raw":0.12},"lastPrice":{"raw":14.5},"targetPrice":{"raw":18}}],
	"topGrowthCompanies":[{"symbol":"GPRO","name":"GoPro","ytdReturn":{"raw":-0.05},"growthEstimate":{"raw":0.3}}]}}`

const marketSummaryJSON = `{"marketSummaryResponse":{"result":[
	{"symbol":"^GSPC","shortName":"S&P 500","regularMarketPrice":4500.5,"regularMarketChange":12.3,"regularMarketChangePercent":0.27},
	{"symbol":"^DJI","shortName":"Dow Jones Industrial Average Index Long Name","regularMarketPrice":35000,"regularMarketChange":-50,"regularMarketChangePercent":-0.14}
],"error":null}}`

const marketStatusJSON = `{"finance":{"marketTimes":[{"marketTime":[
	{"id":"us","name":"U.S. markets","status":"open","message":"U.S. markets close in 2 hours","open":"2023-11-15T14:30:00Z","close":"2023-11-15T21:00:00Z","timezone":[{"gmtoffset":"-18000","short":"EST"}]},
	{"id":"nasdaq","status":"closed","open":"","close":""}
]}],"error":null}}`

// --- handler tests ---

func TestHandleGetQuote(t *testing.T) {
	f := &fakeYahoo{routes: map[string]string{"/v10/finance/quoteSummary/AAPL": quoteSummaryJSON}}
	h := newHandlers(f)

	res, err := h.HandleGetQuote(context.Background(), call(map[string]any{"symbol": "aapl"}))
	out := wantOK(t, res, err, "=== AAPL (Apple Inc.)", "Price:           $178.72", "Change:          -1.28 (-0.71%)",
		"Market Cap:      $2.80T", "P/E Ratio:       29.50", "52-Week Range:   $124.17 - $199.62", "Dividend Yield:  0.55%")
	if strings.Contains(out, "Pre-Market") || strings.Contains(out, "Post-Market") {
		t.Error("REGULAR session should not print pre/post market blocks")
	}

	wantToolError(t, mustCall(h.HandleGetQuote, nil), nil, "symbol is required")

	f.routes["/v10/finance/quoteSummary/NOPE"] = yahooErrJSON
	wantToolError(t, mustCall(h.HandleGetQuote, map[string]any{"symbol": "NOPE"}), nil, "Quote not found")
}

func TestHandleGetChart(t *testing.T) {
	f := &fakeYahoo{routes: map[string]string{"/v8/finance/chart/AAPL": chartJSON}}
	h := newHandlers(f)

	res, err := h.HandleGetChart(context.Background(), call(map[string]any{"symbol": "AAPL", "range": "5d"}))
	wantOK(t, res, err, "=== AAPL Chart Data", "Timezone: America/New_York", "2023-11-15 09:30", "N/A", "Showing 2 of 2 data points")

	if !strings.Contains(strings.Join(f.calls, "\n"), "range=5d") || !strings.Contains(strings.Join(f.calls, "\n"), "interval=1d") {
		t.Errorf("expected range/interval params in request, got %v", f.calls)
	}

	wantToolError(t, mustCall(h.HandleGetChart, map[string]any{}), nil, "symbol is required")

	f.status = 500
	wantToolError(t, mustCall(h.HandleGetChart, map[string]any{"symbol": "AAPL"}), nil, "Failed to get chart for AAPL")
}

func TestHandleSearch(t *testing.T) {
	f := &fakeYahoo{routes: map[string]string{"/v1/finance/search": searchJSON}}
	h := newHandlers(f)

	res, err := h.HandleSearch(context.Background(), call(map[string]any{"query": "Apple", "limit": 2}))
	wantOK(t, res, err, "=== Search Results (2)", "1. AAPL - Apple Inc.", "Sector: Technology | Industry: Consumer Electronics", "2. APLE - Apple Hospitality")

	if !strings.Contains(f.calls[len(f.calls)-1], "quotesCount=2") {
		t.Errorf("limit not forwarded: %s", f.calls[len(f.calls)-1])
	}

	wantToolError(t, mustCall(h.HandleSearch, map[string]any{}), nil, "query is required")

	f.status = 503
	wantToolError(t, mustCall(h.HandleSearch, map[string]any{"query": "x"}), nil, "Search failed")
}

func TestHandleGetFinancials(t *testing.T) {
	f := &fakeYahoo{routes: map[string]string{"/ws/fundamentals-timeseries/v1/finance/timeseries/AAPL": financialsJSON}}
	h := newHandlers(f)

	res, err := h.HandleGetFinancials(context.Background(), call(map[string]any{"symbol": "AAPL"}))
	wantOK(t, res, err, "=== AAPL Annual Financial Data (Income)", "Total Revenue:", "2023-09-30: $383.29B", "Net Income:", "$97.00B")

	res, err = h.HandleGetFinancials(context.Background(), call(map[string]any{"symbol": "AAPL", "statement": "balance", "quarterly": true}))
	wantOK(t, res, err, "Quarterly Financial Data (Balance)")

	wantToolError(t, mustCall(h.HandleGetFinancials, map[string]any{}), nil, "symbol is required")
	wantToolError(t, mustCall(h.HandleGetFinancials, map[string]any{"symbol": "AAPL", "statement": "bogus"}), nil, "invalid statement type")
}

func TestHandleGetOptions(t *testing.T) {
	f := &fakeYahoo{routes: map[string]string{"/v7/finance/options/AAPL": optionsJSON}}
	h := newHandlers(f)

	res, err := h.HandleGetOptions(context.Background(), call(map[string]any{"symbol": "AAPL", "expiration": "1700006400"}))
	wantOK(t, res, err, "=== AAPL Options Chain", "Underlying Price: $178.72", "... and 1 more", "CALLS (Exp: 2023-11-15)", "PUTS (Exp: 2023-11-15)", "* = In the money")

	if !strings.Contains(f.calls[len(f.calls)-1], "date=1700006400") {
		t.Errorf("expiration not forwarded: %s", f.calls[len(f.calls)-1])
	}

	wantToolError(t, mustCall(h.HandleGetOptions, map[string]any{}), nil, "symbol is required")

	f.status = 500
	wantToolError(t, mustCall(h.HandleGetOptions, map[string]any{"symbol": "AAPL"}), nil, "Failed to get options")
}

func TestHandleGetRecommendations(t *testing.T) {
	f := &fakeYahoo{routes: map[string]string{"/v10/finance/quoteSummary/AAPL": quoteSummaryJSON}}
	h := newHandlers(f)

	res, err := h.HandleGetRecommendations(context.Background(), call(map[string]any{"symbol": "AAPL"}))
	wantOK(t, res, err, "=== AAPL Analyst Recommendations", "0m", "Current consensus (0m): 30 bullish, 8 neutral, 1 bearish out of 39 analysts")

	wantToolError(t, mustCall(h.HandleGetRecommendations, map[string]any{}), nil, "symbol is required")

	f.routes["/v10/finance/quoteSummary/NOPE"] = yahooErrJSON
	wantToolError(t, mustCall(h.HandleGetRecommendations, map[string]any{"symbol": "NOPE"}), nil, "Failed to get recommendations")
}

func TestHandleGetNews(t *testing.T) {
	f := &fakeYahoo{routes: map[string]string{"/v1/finance/search": searchJSON}}
	h := newHandlers(f)

	res, err := h.HandleGetNews(context.Background(), call(map[string]any{"symbol": "AAPL", "count": 1}))
	wantOK(t, res, err, "=== AAPL Recent News", "1. Apple hits record", "Publisher: Reuters | 2023-11-15 14:30 UTC", "Link: https://r.example/1")

	if !strings.Contains(f.calls[len(f.calls)-1], "newsCount=1") {
		t.Errorf("count not forwarded: %s", f.calls[len(f.calls)-1])
	}

	wantToolError(t, mustCall(h.HandleGetNews, map[string]any{}), nil, "symbol is required")

	f.status = 500
	wantToolError(t, mustCall(h.HandleGetNews, map[string]any{"symbol": "AAPL"}), nil, "Failed to get news")
}

func TestHandleGetProfile(t *testing.T) {
	f := &fakeYahoo{routes: map[string]string{"/v10/finance/quoteSummary/AAPL": quoteSummaryJSON}}
	h := newHandlers(f)

	res, err := h.HandleGetProfile(context.Background(), call(map[string]any{"symbol": "AAPL"}))
	wantOK(t, res, err, "=== Apple Inc. (AAPL) Company Profile", "Sector:     Technology", "Employees:  161,000",
		"Location:   One Apple Park Way, Cupertino, CA, United States", "Apple designs things.", "- Tim Cook: CEO", "... and 1 more")

	wantToolError(t, mustCall(h.HandleGetProfile, map[string]any{}), nil, "symbol is required")

	f.routes["/v10/finance/quoteSummary/NOPE"] = yahooErrJSON
	wantToolError(t, mustCall(h.HandleGetProfile, map[string]any{"symbol": "NOPE"}), nil, "Failed to get profile")
}

func TestHandleGetBulkQuotes(t *testing.T) {
	f := &fakeYahoo{routes: map[string]string{"/v7/finance/quote": bulkQuoteJSON}}
	h := newHandlers(f)

	res, err := h.HandleGetBulkQuotes(context.Background(), call(map[string]any{"symbols": " aapl, volv-b.st ,, "}))
	out := wantOK(t, res, err, "=== Bulk Quotes (2 symbols)", "AAPL", "$2.80T", "SEK 650.00B", "Volvo, AB ser. B with..")
	if !strings.Contains(out, "+1.60") {
		t.Errorf("expected positive change with sign, got:\n%s", out)
	}
	if !strings.Contains(f.calls[len(f.calls)-1], "symbols=AAPL%2CVOLV-B.ST") {
		t.Errorf("symbols not normalised and joined: %s", f.calls[len(f.calls)-1])
	}

	wantToolError(t, mustCall(h.HandleGetBulkQuotes, map[string]any{}), nil, "symbols is required")
	wantToolError(t, mustCall(h.HandleGetBulkQuotes, map[string]any{"symbols": " , , "}), nil, "at least one symbol is required")

	many := strings.Repeat("A,", 51)
	wantToolError(t, mustCall(h.HandleGetBulkQuotes, map[string]any{"symbols": many}), nil, "too many symbols")
}

func TestHandleGetBulkSpark(t *testing.T) {
	f := &fakeYahoo{routes: map[string]string{"/v8/finance/spark": sparkJSON}}
	h := newHandlers(f)

	res, err := h.HandleGetBulkSpark(context.Background(), call(map[string]any{"symbols": "AAPL,MSFT", "range": "1d", "interval": "1h"}))
	wantOK(t, res, err, "=== Bulk Price History (1 symbols)", "--- AAPL (interval: 1h) ---", "Previous Close: 176.57", "--- MSFT: no data ---", "(2 of 2 data points)")

	wantToolError(t, mustCall(h.HandleGetBulkSpark, map[string]any{}), nil, "symbols is required")
	wantToolError(t, mustCall(h.HandleGetBulkSpark, map[string]any{"symbols": ","}), nil, "at least one symbol is required")

	f.status = 500
	wantToolError(t, mustCall(h.HandleGetBulkSpark, map[string]any{"symbols": "AAPL"}), nil, "Failed to get bulk spark data")
}

func TestHandleGetSector(t *testing.T) {
	f := &fakeYahoo{routes: map[string]string{"/v1/finance/sectors/technology": sectorJSON}}
	h := newHandlers(f)

	res, err := h.HandleGetSector(context.Background(), call(map[string]any{"key": " Technology "}))
	wantOK(t, res, err, "=== Sector: Technology", "Tech sector.", "Companies:    800", "Market Cap:   $15.00T", "Market Weight: 29.00%",
		"Employees:    9,000,000", "AAPL     Apple Inc. (18.00%)", "key: consumer-electronics", "XLK", "FSPTX")

	wantToolError(t, mustCall(h.HandleGetSector, map[string]any{}), nil, "key is required")

	f.status = 404
	wantToolError(t, mustCall(h.HandleGetSector, map[string]any{"key": "technology"}), nil, `Failed to get sector "technology"`)
}

func TestHandleGetIndustry(t *testing.T) {
	f := &fakeYahoo{routes: map[string]string{"/v1/finance/industries/consumer-electronics": industryJSON}}
	h := newHandlers(f)

	res, err := h.HandleGetIndustry(context.Background(), call(map[string]any{"key": "consumer-electronics"}))
	wantOK(t, res, err, "=== Industry: Consumer Electronics", "Sector: Technology (technology)", "Gadgets.", "Market Cap:   $3.00T",
		"Top Performing Companies", "SONO     Sonos Inc. with an unn..", "12.0%", "Top Growth Companies", "GPRO", "-5.0%", "30.0%")

	wantToolError(t, mustCall(h.HandleGetIndustry, map[string]any{}), nil, "key is required")

	f.status = 404
	wantToolError(t, mustCall(h.HandleGetIndustry, map[string]any{"key": "consumer-electronics"}), nil, "Failed to get industry")
}

func TestHandleGetMarketSummary(t *testing.T) {
	f := &fakeYahoo{routes: map[string]string{"/v6/finance/quote/marketSummary": marketSummaryJSON}}
	h := newHandlers(f)

	res, err := h.HandleGetMarketSummary(context.Background(), call(map[string]any{"market": " us "}))
	wantOK(t, res, err, "=== US Market Summary", "^GSPC", "S&P 500", "+12.30", "+0.27%", "Dow Jones Industrial Averag..", "-50.00", "-0.14%")

	if !strings.Contains(f.calls[len(f.calls)-1], "market=US") {
		t.Errorf("market not uppercased: %s", f.calls[len(f.calls)-1])
	}

	wantToolError(t, mustCall(h.HandleGetMarketSummary, map[string]any{}), nil, "market is required")

	f.status = 500
	wantToolError(t, mustCall(h.HandleGetMarketSummary, map[string]any{"market": "US"}), nil, "Failed to get market summary")
}

func TestHandleGetMarketStatus(t *testing.T) {
	f := &fakeYahoo{routes: map[string]string{"/v6/finance/markettime": marketStatusJSON}}
	h := newHandlers(f)

	res, err := h.HandleGetMarketStatus(context.Background(), call(map[string]any{"market": "US"}))
	wantOK(t, res, err, "=== US Market Status", "--- U.S. markets ---", "Status:  open", "Info:    U.S. markets close in 2 hours",
		"Timezone: EST", "--- nasdaq ---", "Status:  closed")

	wantToolError(t, mustCall(h.HandleGetMarketStatus, map[string]any{}), nil, "market is required")

	f.status = 500
	wantToolError(t, mustCall(h.HandleGetMarketStatus, map[string]any{"market": "US"}), nil, "Failed to get market status")
}

// mustCall invokes a handler and panics on a Go-level error, which handlers
// never return by design (errors are reported as tool results).
func mustCall(h func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error), args map[string]any) *mcp.CallToolResult {
	res, err := h(context.Background(), call(args))
	if err != nil {
		panic(err)
	}
	return res
}
