package tools

import "github.com/mark3labs/mcp-go/mcp"

// GetQuoteTool returns the MCP tool definition for get_quote.
func GetQuoteTool() mcp.Tool {
	return mcp.NewTool("get_quote",
		mcp.WithDescription("Get real-time stock quote including price, change, volume, market cap, P/E ratio, and 52-week range"),
		mcp.WithString("symbol",
			mcp.Description("Stock ticker symbol (e.g., AAPL, MSFT, GOOGL)"),
			mcp.Required(),
		),
	)
}

// GetChartTool returns the MCP tool definition for get_chart.
func GetChartTool() mcp.Tool {
	return mcp.NewTool("get_chart",
		mcp.WithDescription("Get historical OHLCV (Open, High, Low, Close, Volume) chart data for a stock"),
		mcp.WithString("symbol",
			mcp.Description("Stock ticker symbol (e.g., AAPL, MSFT, GOOGL)"),
			mcp.Required(),
		),
		mcp.WithString("range",
			mcp.Description("Time range: 1d, 5d, 1mo, 3mo, 6mo, 1y, 2y, 5y, 10y, ytd, max (default: 1mo)"),
		),
		mcp.WithString("interval",
			mcp.Description("Data interval: 1m, 2m, 5m, 15m, 30m, 60m, 90m, 1h, 1d, 5d, 1wk, 1mo, 3mo (default: 1d)"),
		),
	)
}

// SearchTool returns the MCP tool definition for search.
func SearchTool() mcp.Tool {
	return mcp.NewTool("search",
		mcp.WithDescription("Search for stock symbols and companies by name or ticker"),
		mcp.WithString("query",
			mcp.Description("Search query (company name or ticker symbol)"),
			mcp.Required(),
		),
		mcp.WithNumber("limit",
			mcp.Description("Maximum number of results to return (default: 10)"),
		),
	)
}

// GetFinancialsTool returns the MCP tool definition for get_financials.
func GetFinancialsTool() mcp.Tool {
	return mcp.NewTool("get_financials",
		mcp.WithDescription("Get financial statements: income statement, balance sheet, or cash flow"),
		mcp.WithString("symbol",
			mcp.Description("Stock ticker symbol (e.g., AAPL, MSFT, GOOGL)"),
			mcp.Required(),
		),
		mcp.WithString("statement",
			mcp.Description("Financial statement type: income, balance, or cashflow (default: income)"),
		),
		mcp.WithBoolean("quarterly",
			mcp.Description("If true, return quarterly data instead of annual (default: false)"),
		),
	)
}

// GetOptionsTool returns the MCP tool definition for get_options.
func GetOptionsTool() mcp.Tool {
	return mcp.NewTool("get_options",
		mcp.WithDescription("Get options chain (calls and puts) with strike prices, volume, open interest, and implied volatility"),
		mcp.WithString("symbol",
			mcp.Description("Stock ticker symbol (e.g., AAPL, MSFT, GOOGL)"),
			mcp.Required(),
		),
		mcp.WithString("expiration",
			mcp.Description("Expiration date as Unix timestamp (omit for nearest expiration)"),
		),
	)
}

// GetRecommendationsTool returns the MCP tool definition for get_recommendations.
func GetRecommendationsTool() mcp.Tool {
	return mcp.NewTool("get_recommendations",
		mcp.WithDescription("Get analyst recommendation trends (strong buy, buy, hold, sell, strong sell)"),
		mcp.WithString("symbol",
			mcp.Description("Stock ticker symbol (e.g., AAPL, MSFT, GOOGL)"),
			mcp.Required(),
		),
	)
}

// GetNewsTool returns the MCP tool definition for get_news.
func GetNewsTool() mcp.Tool {
	return mcp.NewTool("get_news",
		mcp.WithDescription("Get recent news articles for a stock symbol"),
		mcp.WithString("symbol",
			mcp.Description("Stock ticker symbol (e.g., AAPL, MSFT, GOOGL)"),
			mcp.Required(),
		),
		mcp.WithNumber("count",
			mcp.Description("Number of news articles to return (default: 5)"),
		),
	)
}

// GetProfileTool returns the MCP tool definition for get_profile.
func GetProfileTool() mcp.Tool {
	return mcp.NewTool("get_profile",
		mcp.WithDescription("Get company profile: sector, industry, description, website, and key executives"),
		mcp.WithString("symbol",
			mcp.Description("Stock ticker symbol (e.g., AAPL, MSFT, GOOGL)"),
			mcp.Required(),
		),
	)
}
