package tools

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/emmanuelay/yahoo-finance-mcp/yahoo"
	"github.com/mark3labs/mcp-go/mcp"
)

// Handlers holds the Yahoo Finance client and provides MCP tool handler functions.
type Handlers struct {
	client *yahoo.Client
}

// NewHandlers creates a new Handlers instance with the given Yahoo Finance client.
func NewHandlers(client *yahoo.Client) *Handlers {
	return &Handlers{client: client}
}

// HandleGetQuote handles the get_quote tool call.
func (h *Handlers) HandleGetQuote(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	symbol := req.GetString("symbol", "")
	if symbol == "" {
		return mcp.NewToolResultError("symbol is required"), nil
	}
	symbol = strings.ToUpper(symbol)

	price, detail, err := h.client.GetQuote(symbol)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get quote for %s: %v", symbol, err)), nil
	}

	return mcp.NewToolResultText(formatQuote(price, detail)), nil
}

// HandleGetChart handles the get_chart tool call.
func (h *Handlers) HandleGetChart(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	symbol := strings.ToUpper(req.GetString("symbol", ""))
	if symbol == "" {
		return mcp.NewToolResultError("symbol is required"), nil
	}

	rangeStr := req.GetString("range", "1mo")
	interval := req.GetString("interval", "1d")

	chart, err := h.client.GetChart(symbol, rangeStr, interval)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get chart for %s: %v", symbol, err)), nil
	}

	return mcp.NewToolResultText(formatChart(chart)), nil
}

// HandleSearch handles the search tool call.
func (h *Handlers) HandleSearch(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query := req.GetString("query", "")
	if query == "" {
		return mcp.NewToolResultError("query is required"), nil
	}

	limit := req.GetInt("limit", 10)

	results, err := h.client.Search(query, limit)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Search failed: %v", err)), nil
	}

	return mcp.NewToolResultText(formatSearch(results)), nil
}

// HandleGetFinancials handles the get_financials tool call.
func (h *Handlers) HandleGetFinancials(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	symbol := strings.ToUpper(req.GetString("symbol", ""))
	if symbol == "" {
		return mcp.NewToolResultError("symbol is required"), nil
	}

	statement := req.GetString("statement", "income")
	quarterly := req.GetBool("quarterly", false)

	results, err := h.client.GetFinancials(symbol, statement, quarterly)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get financials for %s: %v", symbol, err)), nil
	}

	return mcp.NewToolResultText(formatFinancials(symbol, statement, quarterly, results)), nil
}

// HandleGetOptions handles the get_options tool call.
func (h *Handlers) HandleGetOptions(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	symbol := strings.ToUpper(req.GetString("symbol", ""))
	if symbol == "" {
		return mcp.NewToolResultError("symbol is required"), nil
	}

	expiration := req.GetString("expiration", "")

	result, err := h.client.GetOptions(symbol, expiration)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get options for %s: %v", symbol, err)), nil
	}

	return mcp.NewToolResultText(formatOptions(result)), nil
}

// HandleGetRecommendations handles the get_recommendations tool call.
func (h *Handlers) HandleGetRecommendations(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	symbol := strings.ToUpper(req.GetString("symbol", ""))
	if symbol == "" {
		return mcp.NewToolResultError("symbol is required"), nil
	}

	trend, err := h.client.GetRecommendations(symbol)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get recommendations for %s: %v", symbol, err)), nil
	}

	return mcp.NewToolResultText(formatRecommendations(symbol, trend)), nil
}

// HandleGetNews handles the get_news tool call.
func (h *Handlers) HandleGetNews(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	symbol := strings.ToUpper(req.GetString("symbol", ""))
	if symbol == "" {
		return mcp.NewToolResultError("symbol is required"), nil
	}

	count := req.GetInt("count", 5)

	news, err := h.client.GetNews(symbol, count)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get news for %s: %v", symbol, err)), nil
	}

	return mcp.NewToolResultText(formatNews(symbol, news)), nil
}

// HandleGetProfile handles the get_profile tool call.
func (h *Handlers) HandleGetProfile(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	symbol := strings.ToUpper(req.GetString("symbol", ""))
	if symbol == "" {
		return mcp.NewToolResultError("symbol is required"), nil
	}

	profile, quoteType, err := h.client.GetProfile(symbol)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get profile for %s: %v", symbol, err)), nil
	}

	return mcp.NewToolResultText(formatProfile(symbol, profile, quoteType)), nil
}

// --- Text formatters ---

func formatQuote(price *yahoo.PriceData, detail *yahoo.SummaryDetailData) string {
	if price == nil {
		return "No price data available"
	}

	var b strings.Builder
	name := price.LongName
	if name == "" {
		name = price.ShortName
	}
	fmt.Fprintf(&b, "=== %s (%s) ===\n", price.Symbol, name)
	fmt.Fprintf(&b, "Exchange: %s | Currency: %s | Market: %s\n\n", price.ExchangeName, price.Currency, price.MarketState)

	fmt.Fprintf(&b, "Price:           %s\n", fmtPrice(price.RegularMarketPrice.Raw, price.Currency))
	changeSign := "+"
	if price.RegularMarketChange.Raw < 0 {
		changeSign = ""
	}
	fmt.Fprintf(&b, "Change:          %s%.2f (%s%.2f%%)\n", changeSign, price.RegularMarketChange.Raw, changeSign, price.RegularMarketChangePercent.Raw)
	fmt.Fprintf(&b, "Volume:          %s\n", fmtInt(price.RegularMarketVolume.Raw))

	if price.MarketCap.Raw > 0 {
		fmt.Fprintf(&b, "Market Cap:      %s\n", fmtLargeNumber(float64(price.MarketCap.Raw)))
	}

	fmt.Fprintf(&b, "Open:            %s\n", fmtPrice(price.RegularMarketOpen.Raw, price.Currency))
	fmt.Fprintf(&b, "Day High:        %s\n", fmtPrice(price.RegularMarketDayHigh.Raw, price.Currency))
	fmt.Fprintf(&b, "Day Low:         %s\n", fmtPrice(price.RegularMarketDayLow.Raw, price.Currency))
	fmt.Fprintf(&b, "Previous Close:  %s\n", fmtPrice(price.RegularMarketPreviousClose.Raw, price.Currency))

	if detail != nil {
		if detail.TrailingPE.Raw > 0 {
			fmt.Fprintf(&b, "P/E Ratio:       %.2f\n", detail.TrailingPE.Raw)
		}
		if detail.ForwardPE.Raw > 0 {
			fmt.Fprintf(&b, "Forward P/E:     %.2f\n", detail.ForwardPE.Raw)
		}
		if detail.FiftyTwoWeekLow.Raw > 0 || detail.FiftyTwoWeekHigh.Raw > 0 {
			fmt.Fprintf(&b, "52-Week Range:   %s - %s\n",
				fmtPrice(detail.FiftyTwoWeekLow.Raw, price.Currency),
				fmtPrice(detail.FiftyTwoWeekHigh.Raw, price.Currency))
		}
		if detail.DividendYield.Raw > 0 {
			fmt.Fprintf(&b, "Dividend Yield:  %.2f%%\n", detail.DividendYield.Raw*100)
		}
		if detail.Beta.Raw > 0 {
			fmt.Fprintf(&b, "Beta:            %.2f\n", detail.Beta.Raw)
		}
		if detail.FiftyDayAverage.Raw > 0 {
			fmt.Fprintf(&b, "50-Day Avg:      %s\n", fmtPrice(detail.FiftyDayAverage.Raw, price.Currency))
		}
		if detail.TwoHundredDayAverage.Raw > 0 {
			fmt.Fprintf(&b, "200-Day Avg:     %s\n", fmtPrice(detail.TwoHundredDayAverage.Raw, price.Currency))
		}
	}

	// Pre/Post market info
	if price.MarketState == "PRE" && price.PreMarketPrice.Raw > 0 {
		fmt.Fprintf(&b, "\n--- Pre-Market ---\n")
		fmt.Fprintf(&b, "Price:  %s (%.2f, %.2f%%)\n", fmtPrice(price.PreMarketPrice.Raw, price.Currency), price.PreMarketChange.Raw, price.PreMarketChangePercent.Raw)
	}
	if price.MarketState == "POST" && price.PostMarketPrice.Raw > 0 {
		fmt.Fprintf(&b, "\n--- Post-Market ---\n")
		fmt.Fprintf(&b, "Price:  %s (%.2f, %.2f%%)\n", fmtPrice(price.PostMarketPrice.Raw, price.Currency), price.PostMarketChange.Raw, price.PostMarketChangePercent.Raw)
	}

	return b.String()
}

func formatChart(chart *yahoo.ChartResult) string {
	var b strings.Builder

	fmt.Fprintf(&b, "=== %s Chart Data ===\n", chart.Meta.Symbol)
	fmt.Fprintf(&b, "Exchange: %s | Currency: %s\n", chart.Meta.ExchangeName, chart.Meta.Currency)
	fmt.Fprintf(&b, "Range: %s | Interval: %s\n\n", chart.Meta.Range, chart.Meta.DataGranularity)

	if len(chart.Timestamps) == 0 || len(chart.Indicators.Quote) == 0 {
		fmt.Fprintf(&b, "No data points available\n")
		return b.String()
	}

	q := chart.Indicators.Quote[0]

	// Cap at 100 rows to avoid flooding LLM context
	step := 1
	total := len(chart.Timestamps)
	if total > 100 {
		step = total / 100
	}

	fmt.Fprintf(&b, "%-20s %10s %10s %10s %10s %12s\n", "Date", "Open", "High", "Low", "Close", "Volume")
	fmt.Fprintf(&b, "%s\n", strings.Repeat("-", 82))

	for i := 0; i < total; i += step {
		ts := time.Unix(chart.Timestamps[i], 0)
		dateStr := ts.Format("2006-01-02 15:04")

		open := fmtOptFloat(q.Open, i)
		high := fmtOptFloat(q.High, i)
		low := fmtOptFloat(q.Low, i)
		close := fmtOptFloat(q.Close, i)
		vol := fmtOptInt(q.Volume, i)

		fmt.Fprintf(&b, "%-20s %10s %10s %10s %10s %12s\n", dateStr, open, high, low, close, vol)
	}

	fmt.Fprintf(&b, "\nShowing %d of %d data points\n", min(total, 100), total)

	return b.String()
}

func formatSearch(results *yahoo.SearchResponse) string {
	var b strings.Builder

	if len(results.Quotes) == 0 {
		return "No results found"
	}

	fmt.Fprintf(&b, "=== Search Results (%d) ===\n\n", len(results.Quotes))

	for i, q := range results.Quotes {
		name := q.LongName
		if name == "" {
			name = q.ShortName
		}
		fmt.Fprintf(&b, "%d. %s - %s\n", i+1, q.Symbol, name)
		fmt.Fprintf(&b, "   Exchange: %s | Type: %s\n", q.Exchange, q.QuoteType)
		if q.Sector != "" {
			fmt.Fprintf(&b, "   Sector: %s | Industry: %s\n", q.Sector, q.Industry)
		}
		if i < len(results.Quotes)-1 {
			fmt.Fprintln(&b)
		}
	}

	return b.String()
}

func formatFinancials(symbol, statement string, quarterly bool, results []yahoo.FinancialResult) string {
	var b strings.Builder

	period := "Annual"
	if quarterly {
		period = "Quarterly"
	}

	stmtTitle := strings.ToUpper(statement[:1]) + statement[1:]
	fmt.Fprintf(&b, "=== %s %s Financial Data (%s) ===\n\n", symbol, period, stmtTitle)

	if len(results) == 0 {
		fmt.Fprintf(&b, "No financial data available\n")
		return b.String()
	}

	for _, r := range results {
		// Clean up type name for display
		displayName := r.Type
		displayName = strings.TrimPrefix(displayName, "annual")
		displayName = strings.TrimPrefix(displayName, "quarterly")
		// Add spaces before capital letters
		displayName = addSpaces(displayName)

		fmt.Fprintf(&b, "%s:\n", displayName)
		for _, item := range r.Items {
			fmt.Fprintf(&b, "  %s: %s %s\n", item.Date, fmtLargeNumber(item.ReportedValue), item.CurrencyCode)
		}
		fmt.Fprintln(&b)
	}

	return b.String()
}

func formatOptions(result *yahoo.OptionsResult) string {
	var b strings.Builder

	fmt.Fprintf(&b, "=== %s Options Chain ===\n", result.UnderlyingSymbol)
	fmt.Fprintf(&b, "Underlying Price: $%.2f\n", result.Quote.RegularMarketPrice)

	if len(result.ExpirationDates) > 0 {
		fmt.Fprintf(&b, "Available Expirations: ")
		for i, d := range result.ExpirationDates {
			if i > 5 {
				fmt.Fprintf(&b, "... and %d more", len(result.ExpirationDates)-i)
				break
			}
			if i > 0 {
				fmt.Fprintf(&b, ", ")
			}
			fmt.Fprintf(&b, "%s", time.Unix(d, 0).Format("2006-01-02"))
		}
		fmt.Fprintln(&b)
	}

	if len(result.Options) == 0 {
		fmt.Fprintf(&b, "\nNo options data available\n")
		return b.String()
	}

	chain := result.Options[0]
	expDate := time.Unix(chain.ExpirationDate, 0).Format("2006-01-02")

	// Calls
	fmt.Fprintf(&b, "\n--- CALLS (Exp: %s) ---\n", expDate)
	fmt.Fprintf(&b, "%-10s %10s %10s %10s %10s %10s %8s\n", "Strike", "Last", "Bid", "Ask", "Volume", "OI", "IV")
	fmt.Fprintf(&b, "%s\n", strings.Repeat("-", 78))

	maxContracts := 25
	for i, c := range chain.Calls {
		if i >= maxContracts {
			fmt.Fprintf(&b, "... and %d more contracts\n", len(chain.Calls)-maxContracts)
			break
		}
		itm := ""
		if c.InTheMoney {
			itm = "*"
		}
		fmt.Fprintf(&b, "%-10s %10.2f %10.2f %10.2f %10d %10d %7.1f%%\n",
			fmt.Sprintf("%.2f%s", c.Strike, itm), c.LastPrice, c.Bid, c.Ask, c.Volume, c.OpenInterest, c.ImpliedVolatility*100)
	}

	// Puts
	fmt.Fprintf(&b, "\n--- PUTS (Exp: %s) ---\n", expDate)
	fmt.Fprintf(&b, "%-10s %10s %10s %10s %10s %10s %8s\n", "Strike", "Last", "Bid", "Ask", "Volume", "OI", "IV")
	fmt.Fprintf(&b, "%s\n", strings.Repeat("-", 78))

	for i, p := range chain.Puts {
		if i >= maxContracts {
			fmt.Fprintf(&b, "... and %d more contracts\n", len(chain.Puts)-maxContracts)
			break
		}
		itm := ""
		if p.InTheMoney {
			itm = "*"
		}
		fmt.Fprintf(&b, "%-10s %10.2f %10.2f %10.2f %10d %10d %7.1f%%\n",
			fmt.Sprintf("%.2f%s", p.Strike, itm), p.LastPrice, p.Bid, p.Ask, p.Volume, p.OpenInterest, p.ImpliedVolatility*100)
	}

	fmt.Fprintf(&b, "\n* = In the money\n")

	return b.String()
}

func formatRecommendations(symbol string, trend *yahoo.RecommendationTrendData) string {
	var b strings.Builder

	fmt.Fprintf(&b, "=== %s Analyst Recommendations ===\n\n", symbol)

	if trend == nil || len(trend.Trend) == 0 {
		fmt.Fprintf(&b, "No recommendation data available\n")
		return b.String()
	}

	fmt.Fprintf(&b, "%-12s %10s %6s %6s %6s %11s %6s\n", "Period", "Strong Buy", "Buy", "Hold", "Sell", "Strong Sell", "Total")
	fmt.Fprintf(&b, "%s\n", strings.Repeat("-", 63))

	for _, t := range trend.Trend {
		total := t.StrongBuy + t.Buy + t.Hold + t.Sell + t.StrongSell
		fmt.Fprintf(&b, "%-12s %10d %6d %6d %6d %11d %6d\n",
			t.Period, t.StrongBuy, t.Buy, t.Hold, t.Sell, t.StrongSell, total)
	}

	// Summary for current period
	if len(trend.Trend) > 0 {
		current := trend.Trend[0]
		total := current.StrongBuy + current.Buy + current.Hold + current.Sell + current.StrongSell
		if total > 0 {
			bullish := current.StrongBuy + current.Buy
			bearish := current.Sell + current.StrongSell
			fmt.Fprintf(&b, "\nCurrent consensus (%s): %d bullish, %d neutral, %d bearish out of %d analysts\n",
				current.Period, bullish, current.Hold, bearish, total)
		}
	}

	return b.String()
}

func formatNews(symbol string, news []yahoo.SearchNews) string {
	var b strings.Builder

	fmt.Fprintf(&b, "=== %s Recent News ===\n\n", symbol)

	if len(news) == 0 {
		fmt.Fprintf(&b, "No recent news found\n")
		return b.String()
	}

	for i, n := range news {
		pubTime := time.Unix(n.ProviderPublishTime, 0)
		fmt.Fprintf(&b, "%d. %s\n", i+1, n.Title)
		fmt.Fprintf(&b, "   Publisher: %s | %s\n", n.Publisher, pubTime.Format("2006-01-02 15:04 MST"))
		if n.Link != "" {
			fmt.Fprintf(&b, "   Link: %s\n", n.Link)
		}
		if i < len(news)-1 {
			fmt.Fprintln(&b)
		}
	}

	return b.String()
}

func formatProfile(symbol string, profile *yahoo.AssetProfileData, quoteType *yahoo.QuoteTypeData) string {
	var b strings.Builder

	name := symbol
	if quoteType != nil && quoteType.LongName != "" {
		name = quoteType.LongName
	}

	fmt.Fprintf(&b, "=== %s (%s) Company Profile ===\n\n", name, symbol)

	if profile == nil {
		fmt.Fprintf(&b, "No profile data available\n")
		return b.String()
	}

	if profile.Sector != "" {
		fmt.Fprintf(&b, "Sector:     %s\n", profile.Sector)
	}
	if profile.Industry != "" {
		fmt.Fprintf(&b, "Industry:   %s\n", profile.Industry)
	}
	if profile.Website != "" {
		fmt.Fprintf(&b, "Website:    %s\n", profile.Website)
	}
	if profile.Phone != "" {
		fmt.Fprintf(&b, "Phone:      %s\n", profile.Phone)
	}
	if profile.FullTimeEmployees > 0 {
		fmt.Fprintf(&b, "Employees:  %s\n", fmtInt(int64(profile.FullTimeEmployees)))
	}

	// Address
	var addrParts []string
	if profile.Address1 != "" {
		addrParts = append(addrParts, profile.Address1)
	}
	if profile.City != "" {
		addrParts = append(addrParts, profile.City)
	}
	if profile.State != "" {
		addrParts = append(addrParts, profile.State)
	}
	if profile.Country != "" {
		addrParts = append(addrParts, profile.Country)
	}
	if len(addrParts) > 0 {
		fmt.Fprintf(&b, "Location:   %s\n", strings.Join(addrParts, ", "))
	}

	// Business summary
	if profile.LongBusinessSummary != "" {
		fmt.Fprintf(&b, "\n--- Business Summary ---\n%s\n", profile.LongBusinessSummary)
	}

	// Key executives
	if len(profile.CompanyOfficers) > 0 {
		fmt.Fprintf(&b, "\n--- Key Executives ---\n")
		limit := 5
		if len(profile.CompanyOfficers) < limit {
			limit = len(profile.CompanyOfficers)
		}
		for _, officer := range profile.CompanyOfficers[:limit] {
			fmt.Fprintf(&b, "- %s: %s\n", officer.Name, officer.Title)
		}
		if len(profile.CompanyOfficers) > 5 {
			fmt.Fprintf(&b, "  ... and %d more\n", len(profile.CompanyOfficers)-5)
		}
	}

	return b.String()
}

// --- Formatting helpers ---

func fmtPrice(val float64, currency string) string {
	symbol := "$"
	switch currency {
	case "EUR":
		symbol = "\u20ac"
	case "GBP", "GBp":
		symbol = "\u00a3"
	case "JPY":
		symbol = "\u00a5"
	case "CHF":
		symbol = "CHF "
	case "SEK":
		symbol = "SEK "
	case "NOK":
		symbol = "NOK "
	case "DKK":
		symbol = "DKK "
	case "CAD":
		symbol = "C$"
	case "AUD":
		symbol = "A$"
	}
	return fmt.Sprintf("%s%.2f", symbol, val)
}

func fmtInt(val int64) string {
	if val == 0 {
		return "0"
	}
	s := fmt.Sprintf("%d", val)
	if val < 0 {
		s = s[1:] // remove negative sign, add back later
	}
	n := len(s)
	if n <= 3 {
		if val < 0 {
			return "-" + s
		}
		return s
	}

	var result strings.Builder
	remainder := n % 3
	if remainder > 0 {
		result.WriteString(s[:remainder])
		if n > remainder {
			result.WriteString(",")
		}
	}
	for i := remainder; i < n; i += 3 {
		result.WriteString(s[i : i+3])
		if i+3 < n {
			result.WriteString(",")
		}
	}
	if val < 0 {
		return "-" + result.String()
	}
	return result.String()
}

func fmtLargeNumber(val float64) string {
	abs := val
	sign := ""
	if val < 0 {
		abs = -val
		sign = "-"
	}

	switch {
	case abs >= 1e12:
		return fmt.Sprintf("%s$%.2fT", sign, abs/1e12)
	case abs >= 1e9:
		return fmt.Sprintf("%s$%.2fB", sign, abs/1e9)
	case abs >= 1e6:
		return fmt.Sprintf("%s$%.2fM", sign, abs/1e6)
	case abs >= 1e3:
		return fmt.Sprintf("%s$%.2fK", sign, abs/1e3)
	default:
		return fmt.Sprintf("%s$%.2f", sign, abs)
	}
}

func fmtOptFloat(vals []*float64, idx int) string {
	if idx >= len(vals) || vals[idx] == nil {
		return "N/A"
	}
	return fmt.Sprintf("%.2f", *vals[idx])
}

func fmtOptInt(vals []*int64, idx int) string {
	if idx >= len(vals) || vals[idx] == nil {
		return "N/A"
	}
	return fmtInt(*vals[idx])
}

func addSpaces(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune(' ')
		}
		result.WriteRune(r)
	}
	return result.String()
}
