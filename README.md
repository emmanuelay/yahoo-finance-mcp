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

## Install

### Desktop Extension (.mcpb) for Claude Desktop (recommended)

Download the **`.mcpb`** file for your platform from the [latest release](https://github.com/emmanuelay/yahoo-finance-mcp/releases/latest):

| Platform | File |
|----------|------|
| macOS (Apple Silicon) | `yahoo-finance-mcp_<version>_darwin_arm64.mcpb` |
| macOS (Intel) | `yahoo-finance-mcp_<version>_darwin_amd64.mcpb` |
| Windows (x64) | `yahoo-finance-mcp_<version>_win32_amd64.mcpb` |
| Windows (ARM64) | `yahoo-finance-mcp_<version>_win32_arm64.mcpb` |
| Linux (x64 / ARM64) | `yahoo-finance-mcp_<version>_linux_<arch>.mcpb` |

Then open Claude Desktop, go to **Settings → Extensions**, and install the file. On macOS you can also double-click the `.mcpb` file.

> [!NOTE]
> The `.zip` and `.tar.gz` files on the release page are plain binary archives for Homebrew and manual installs. They are **not** extensions. Dropping one into Claude Desktop's *Upload local plugin* dialog fails with `Invalid plugin: missing .claude-plugin/plugin.json`, because that dialog is for Claude plugins rather than MCP extensions. Use the `.mcpb` file instead.

### Homebrew (macOS / Linux)

```sh
brew tap emmanuelay/tap
brew install yahoo-finance-mcp
```

This puts `yahoo-finance-mcp` on your `PATH`, so the config examples below work as written.

### Pre-built binary (manual)

Download the `.zip` (Windows) or `.tar.gz` (macOS / Linux) archive for your platform from the [latest release](https://github.com/emmanuelay/yahoo-finance-mcp/releases/latest) and extract it. The archive contains a single binary, `yahoo-finance-mcp` (`yahoo-finance-mcp.exe` on Windows). Either put it somewhere on your `PATH`, or reference it by its full path in the config examples below.

### From source

```sh
git clone https://github.com/emmanuelay/yahoo-finance-mcp.git
cd yahoo-finance-mcp
make build
```

## Usage as a Go package

The `yahoo` package is a standalone client with no dependency on MCP. It
imports only the standard library, so you can use it in your own Go code
without pulling in the server.

```sh
go get github.com/emmanuelay/yahoo-finance-mcp/yahoo
```

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/emmanuelay/yahoo-finance-mcp/yahoo"
)

func main() {
	client := yahoo.NewClient()
	ctx := context.Background()

	price, detail, err := client.GetQuote(ctx, "AAPL")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s %.2f %s (P/E %.1f)\n",
		price.Symbol, price.RegularMarketPrice.Raw, price.Currency, detail.TrailingPE.Raw)

	chart, err := client.GetChart(ctx, "VOLV-B.ST", "5d", "1d")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(chart.Timestamps), "bars in", chart.Meta.ExchangeTimezoneName)
}
```

Every tool in the table above has a matching method on `Client`, all taking a
`context.Context` first. Authentication (Yahoo's cookie and crumb handshake)
happens lazily on the first call that needs it and is retried once on a
401 or 403. Requests have a 30 second timeout by default.

To supply your own transport, proxy or timeout, pass an `http.Client`:

```go
client := yahoo.NewClientWithHTTP(&http.Client{
	Timeout: 10 * time.Second,
	Jar:     jar, // a cookie jar is required for the crumb handshake
})
```

For endpoints the package does not wrap, `Get`, `GetJSON`, `GetAbsolute`
and `GetAbsoluteJSON` are the low-level request methods the typed methods
are built on.

The full API is on [pkg.go.dev](https://pkg.go.dev/github.com/emmanuelay/yahoo-finance-mcp/yahoo).
The module is at v0, so the API may change between releases.

## Enable in Anthropics Claude Desktop / Claude Code

### Claude Desktop

If you installed the `.mcpb` extension you are done; skip this section.

Otherwise add the server to your `claude_desktop_config.json`. If the binary is on your `PATH` (Homebrew), the bare name is enough:

```json
{
  "mcpServers": {
    "yahoo-finance": {
      "command": "yahoo-finance-mcp"
    }
  }
}
```

If you downloaded the binary manually, use its full path. On Windows, double the backslashes:

```json
{
  "mcpServers": {
    "yahoo-finance": {
      "command": "C:\\Users\\you\\Tools\\yahoo-finance-mcp\\yahoo-finance-mcp.exe"
    }
  }
}
```

```json
{
  "mcpServers": {
    "yahoo-finance": {
      "command": "/Users/you/Tools/yahoo-finance-mcp/yahoo-finance-mcp"
    }
  }
}
```

The config file is located at:

- macOS: `~/Library/Application Support/Claude/claude_desktop_config.json`
- Windows: `%APPDATA%\Claude\claude_desktop_config.json`

Restart Claude Desktop after saving the configuration.

> [!TIP]
> `MCP yahoo-finance: spawn yahoo-finance-mcp ENOENT` means Claude Desktop could not find the binary. Either it is not on your `PATH`, or the path in the config is wrong. Using the full path to the binary fixes it.

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
make tags       # List all git tags
make mcpb       # Build .mcpb bundles from goreleaser dist/
```

## License

MIT
