// Package yahoo is a client for the unofficial Yahoo Finance JSON API.
//
// It wraps the quote, chart, search, financials, options, recommendations,
// news, profile, bulk quote, spark, sector, industry and market endpoints
// behind typed methods on Client, and handles Yahoo's cookie and crumb
// authentication transparently. The package depends only on the standard
// library and can be used independently of the MCP server in this module.
//
//	client := yahoo.NewClient()
//	price, detail, err := client.GetQuote(ctx, "AAPL")
//
// Yahoo!, Y!Finance and Yahoo! Finance are registered trademarks of
// Yahoo, Inc. This package is not affiliated with, endorsed or vetted by
// Yahoo, Inc. It uses publicly available endpoints that are intended for
// personal use; refer to Yahoo's terms of service for your rights to the
// data.
package yahoo
