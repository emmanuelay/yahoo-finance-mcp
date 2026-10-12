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

	tools.RegisterAll(s, handlers)

	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
