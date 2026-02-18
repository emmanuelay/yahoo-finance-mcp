package main

import (
	"log"

	"github.com/emmanuelay/yahoo-finance-mcp/tools"
	"github.com/emmanuelay/yahoo-finance-mcp/yahoo"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	client := yahoo.NewClient()
	handlers := tools.NewHandlers(client)

	s := server.NewMCPServer(
		"yahoo-finance",
		"1.0.0",
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

	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
