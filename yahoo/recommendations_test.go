package yahoo

import (
	"net/http"
	"strings"
	"testing"
)

func TestGetRecommendations_Success(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		if !strings.Contains(req.URL.Path, "/v10/finance/quoteSummary/AAPL") {
			t.Errorf("unexpected path: %s", req.URL.Path)
		}
		modules := req.URL.Query().Get("modules")
		if modules != "recommendationTrend" {
			t.Errorf("modules = %q, want %q", modules, "recommendationTrend")
		}
		return jsonResponse(200, `{
			"quoteSummary": {
				"result": [{
					"recommendationTrend": {
						"trend": [
							{
								"period": "0m",
								"strongBuy": 12,
								"buy": 20,
								"hold": 8,
								"sell": 2,
								"strongSell": 1
							},
							{
								"period": "-1m",
								"strongBuy": 10,
								"buy": 18,
								"hold": 9,
								"sell": 3,
								"strongSell": 1
							}
						]
					}
				}]
			}
		}`), nil
	})

	result, err := client.GetRecommendations("AAPL")
	if err != nil {
		t.Fatalf("GetRecommendations() error: %v", err)
	}

	if len(result.Trend) != 2 {
		t.Fatalf("trend count = %d, want 2", len(result.Trend))
	}

	current := result.Trend[0]
	if current.Period != "0m" {
		t.Errorf("trend[0].Period = %q, want %q", current.Period, "0m")
	}
	if current.StrongBuy != 12 {
		t.Errorf("trend[0].StrongBuy = %d, want 12", current.StrongBuy)
	}
	if current.Buy != 20 {
		t.Errorf("trend[0].Buy = %d, want 20", current.Buy)
	}
	if current.Hold != 8 {
		t.Errorf("trend[0].Hold = %d, want 8", current.Hold)
	}
	if current.Sell != 2 {
		t.Errorf("trend[0].Sell = %d, want 2", current.Sell)
	}
	if current.StrongSell != 1 {
		t.Errorf("trend[0].StrongSell = %d, want 1", current.StrongSell)
	}
}

func TestGetRecommendations_YahooError(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		return jsonResponse(200, `{
			"quoteSummary": {
				"result": null,
				"error": {"code": "Not Found", "description": "Symbol not found"}
			}
		}`), nil
	})

	_, err := client.GetRecommendations("INVALID")
	if err == nil {
		t.Fatal("expected error for yahoo error response")
	}
	if !strings.Contains(err.Error(), "Symbol not found") {
		t.Errorf("error should contain description, got: %v", err)
	}
}

func TestGetRecommendations_EmptyResults(t *testing.T) {
	client := newTestClient(func(req *http.Request) (*http.Response, error) {
		return jsonResponse(200, `{"quoteSummary": {"result": []}}`), nil
	})

	_, err := client.GetRecommendations("UNKNOWN")
	if err == nil {
		t.Fatal("expected error for empty results")
	}
	if !strings.Contains(err.Error(), "no data found") {
		t.Errorf("error should mention no data found, got: %v", err)
	}
}
