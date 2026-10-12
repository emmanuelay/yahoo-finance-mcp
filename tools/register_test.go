package tools

import (
	"testing"

	"github.com/mark3labs/mcp-go/server"
)

func TestAllTools_UniqueNamesAndRequiredArgs(t *testing.T) {
	h := NewHandlers(nil)
	all := AllTools(h)

	if len(all) != 14 {
		t.Fatalf("AllTools() returned %d tools, want 14", len(all))
	}

	// Every tool must declare at least one required argument; this guards
	// against a definition accidentally dropping mcp.Required().
	seen := map[string]bool{}
	for _, st := range all {
		name := st.Tool.Name
		if name == "" {
			t.Error("tool with empty name")
		}
		if seen[name] {
			t.Errorf("duplicate tool name %q", name)
		}
		seen[name] = true

		if st.Handler == nil {
			t.Errorf("%s: nil handler", name)
		}
		if st.Tool.Description == "" {
			t.Errorf("%s: empty description", name)
		}
		if len(st.Tool.InputSchema.Required) == 0 {
			t.Errorf("%s: no required arguments", name)
		}
		for _, req := range st.Tool.InputSchema.Required {
			if _, ok := st.Tool.InputSchema.Properties[req]; !ok {
				t.Errorf("%s: required arg %q has no property definition", name, req)
			}
		}
	}

	for _, want := range []string{
		"get_quote", "get_chart", "search", "get_financials", "get_options",
		"get_recommendations", "get_news", "get_profile", "get_bulk_quotes",
		"get_bulk_spark", "get_sector", "get_industry", "get_market_summary",
		"get_market_status",
	} {
		if !seen[want] {
			t.Errorf("missing tool %q", want)
		}
	}
}

func TestRegisterAll(t *testing.T) {
	s := server.NewMCPServer("test", "0.0.0", server.WithToolCapabilities(true))
	n := RegisterAll(s, NewHandlers(nil))
	if n != 14 {
		t.Errorf("RegisterAll() = %d, want 14", n)
	}
}
