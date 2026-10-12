package tools

import "github.com/mark3labs/mcp-go/server"

// AllTools returns every tool definition paired with its handler, in the
// order they are registered on the server.
func AllTools(h *Handlers) []server.ServerTool {
	return []server.ServerTool{
		{Tool: GetQuoteTool(), Handler: h.HandleGetQuote},
		{Tool: GetChartTool(), Handler: h.HandleGetChart},
		{Tool: SearchTool(), Handler: h.HandleSearch},
		{Tool: GetFinancialsTool(), Handler: h.HandleGetFinancials},
		{Tool: GetOptionsTool(), Handler: h.HandleGetOptions},
		{Tool: GetRecommendationsTool(), Handler: h.HandleGetRecommendations},
		{Tool: GetNewsTool(), Handler: h.HandleGetNews},
		{Tool: GetProfileTool(), Handler: h.HandleGetProfile},
		{Tool: GetBulkQuotesTool(), Handler: h.HandleGetBulkQuotes},
		{Tool: GetBulkSparkTool(), Handler: h.HandleGetBulkSpark},
		{Tool: GetSectorTool(), Handler: h.HandleGetSector},
		{Tool: GetIndustryTool(), Handler: h.HandleGetIndustry},
		{Tool: GetMarketSummaryTool(), Handler: h.HandleGetMarketSummary},
		{Tool: GetMarketStatusTool(), Handler: h.HandleGetMarketStatus},
	}
}

// RegisterAll adds every tool to s and returns how many were registered.
func RegisterAll(s *server.MCPServer, h *Handlers) int {
	tools := AllTools(h)
	s.AddTools(tools...)
	return len(tools)
}
