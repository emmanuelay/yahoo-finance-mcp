package main

import (
	"fmt"
	"log"

	"github.com/emmanuelay/yahoo-finance-mcp/tools"
	"github.com/emmanuelay/yahoo-finance-mcp/yahoo"
	"github.com/mark3labs/mcp-go/server"
)

// These will be set at build time
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	client := yahoo.NewClient()
	handlers := tools.NewHandlers(client)

	s := server.NewMCPServer(
		"yahoo-finance",
		fmt.Sprintf("%s (%s) %s", version, commit, date),
		server.WithToolCapabilities(true),
	)

	s.AddTool(tools.GetQuoteTool(), handlers.HandleGetQuote)
	s.AddTool(tools.GetChartTool(), handlers.HandleGetChart)
	s.AddTool(tools.SearchTool(), handlers.HandleSearch)
	s.AddTool(tools.GetFinancialsTool(), handlers.HandleGetFinancials)
	s.AddTool(tools.GetOptionsTool(), handlers.HandleGetOptions)
	s.AddTool(tools.GetRecommendationsTool(), handlers.HandleGetRecommendations)
	s.AddTool(tools.GetNewsTool(), handlers.HandleGetNews)
	s.AddTool(tools.GetProfileTool(), handlers.HandleGetProfile)
	s.AddTool(tools.GetBulkQuotesTool(), handlers.HandleGetBulkQuotes)
	s.AddTool(tools.GetBulkSparkTool(), handlers.HandleGetBulkSpark)
	s.AddTool(tools.GetSectorTool(), handlers.HandleGetSector)
	s.AddTool(tools.GetIndustryTool(), handlers.HandleGetIndustry)
	s.AddTool(tools.GetMarketSummaryTool(), handlers.HandleGetMarketSummary)
	s.AddTool(tools.GetMarketStatusTool(), handlers.HandleGetMarketStatus)

	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
