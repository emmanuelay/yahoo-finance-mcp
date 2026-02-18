# Yahoo finance MCP

An MCP (Model Context Protocol) server that provides access to Yahoo! Finance data. Built with [mcp-go](https://github.com/mark3labs/mcp-go). 

Heavily inspired by [Ran Aroussis](https://github.com/ranaroussi) python project [yfinance](https://github.com/ranaroussi/yfinance)

> [!IMPORTANT]  
> **Yahoo!, Y!Finance, and Yahoo! finance are registered trademarks of Yahoo, Inc.**
>
> yahoo-mcp-server is **not** affiliated, endorsed, or vetted by Yahoo, Inc. It's an open-source tool that uses Yahoo's publicly available APIs, and is intended for research and educational purposes.
> 
> **You should refer to Yahoo!'s terms of use** ([here](https://policies.yahoo.com/us/en/yahoo/terms/product-atos/apiforydn/index.htm), [here](https://legal.yahoo.com/us/en/yahoo/terms/otos/index.html), and [here](https://policies.yahoo.com/us/en/yahoo/terms/index.htm)) **for details on your rights to use the actual data downloaded.
>
> Remember - the Yahoo! finance API is intended for personal use only.**


## Example prompts

- "What's the current price of Apple stock?"
- "Compare the performance of AAPL, MSFT, and GOOGL over the past 6 months"
- "Show me Tesla's income statement for the last 4 quarters"
- "What are analysts recommending for NVIDIA right now?"
- "Find me the top companies in the semiconductor industry"
- "What options are available for AMZN expiring this month?"
- "Give me a summary of the US market today"
- "What's the latest news about Microsoft?"
- "Tell me about AMD - what do they do, what sector are they in?"
- "Is the European market open right now?"

## Tools

| Tool | Description |
|------|-------------|
| `get_quote` | Real-time stock quote with price, change, volume, market cap, P/E ratio, and 52-week range |
| `get_chart` | Historical OHLCV chart data with configurable range and interval |
| `get_bulk_quotes` | Real-time quotes for multiple stocks in a single request (max 50) |
| `get_bulk_spark` | Simplified price history for multiple stocks in a single request (max 50) |
| `search` | Search for stock symbols and companies by name or ticker |
| `get_financials` | Financial statements: income statement, balance sheet, or cash flow |
| `get_options` | Options chain with strike prices, volume, open interest, and implied volatility |
| `get_recommendations` | Analyst recommendation trends |
| `get_news` | Recent news articles for a stock symbol |
| `get_profile` | Company profile: sector, industry, description, website, and key executives |
| `get_sector` | Sector overview: market cap, top companies, ETFs, and industries |
| `get_industry` | Industry overview: top companies, top performers, and growth estimates |
| `get_market_summary` | Market summary with index prices and changes |
| `get_market_status` | Market open/close times and timezone information |

## Install binary

### Homebrew

```sh
brew tap emmanuelay/tap
brew install yahoo-finance-mcp
```

### From source

```sh
git clone https://github.com/emmanuelay/yahoo-finance-mcp.git
cd yahoo-finance-mcp
make build
```

## Enable in Anthropics Claude Desktop / Claude Code

### Claude Desktop

Add the server to your `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "yahoo-finance": {
      "command": "yahoo-finance-mcp"
    }
  }
}
```

The config file is located at:

- macOS: `~/Library/Application Support/Claude/claude_desktop_config.json`
- Windows: `%APPDATA%\Claude\claude_desktop_config.json`

Restart Claude Desktop after saving the configuration.

### Claude Code

```sh
claude mcp add yahoo-finance -- yahoo-finance-mcp
```


## Development

```sh
make help       # Show all available targets
make deps       # Download dependencies
make build      # Build the binary
make test       # Run tests
make install    # Install to $GOPATH/bin
```

## License

MIT
