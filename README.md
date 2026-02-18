# Yahoo finance MCP

An MCP (Model Context Protocol) server that provides access to Yahoo Finance data. Built with [mcp-go](https://github.com/mark3labs/mcp-go).

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

## Install

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
