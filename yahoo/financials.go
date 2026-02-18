package yahoo

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"
)

// Financial statement type names for the timeseries endpoint.
var incomeStatementTypes = []string{
	"annualTotalRevenue", "annualCostOfRevenue", "annualGrossProfit",
	"annualOperatingExpense", "annualOperatingIncome",
	"annualNetIncome", "annualEbitda", "annualBasicEPS", "annualDilutedEPS",
}

var balanceSheetTypes = []string{
	"annualTotalAssets", "annualTotalLiabilitiesNetMinorityInterest",
	"annualStockholdersEquity", "annualCashAndCashEquivalents",
	"annualCurrentAssets", "annualCurrentLiabilities",
	"annualTotalDebt", "annualNetDebt",
}

var cashFlowTypes = []string{
	"annualOperatingCashFlow", "annualInvestingCashFlow",
	"annualFinancingCashFlow", "annualFreeCashFlow",
	"annualCapitalExpenditure",
}

var quarterlyIncomeStatementTypes = []string{
	"quarterlyTotalRevenue", "quarterlyCostOfRevenue", "quarterlyGrossProfit",
	"quarterlyOperatingExpense", "quarterlyOperatingIncome",
	"quarterlyNetIncome", "quarterlyEbitda", "quarterlyBasicEPS", "quarterlyDilutedEPS",
}

var quarterlyBalanceSheetTypes = []string{
	"quarterlyTotalAssets", "quarterlyTotalLiabilitiesNetMinorityInterest",
	"quarterlyStockholdersEquity", "quarterlyCashAndCashEquivalents",
	"quarterlyCurrentAssets", "quarterlyCurrentLiabilities",
	"quarterlyTotalDebt", "quarterlyNetDebt",
}

var quarterlyCashFlowTypes = []string{
	"quarterlyOperatingCashFlow", "quarterlyInvestingCashFlow",
	"quarterlyFinancingCashFlow", "quarterlyFreeCashFlow",
	"quarterlyCapitalExpenditure",
}

// FinancialItem represents a single financial data point.
type FinancialItem struct {
	Date           string  `json:"asOfDate"`
	ReportedValue  float64 `json:"-"`
	CurrencyCode   string  `json:"currencyCode"`
}

// FinancialResult holds parsed financial data for a metric.
type FinancialResult struct {
	Type  string
	Items []FinancialItem
}

// GetFinancials fetches financial statement data for a symbol.
// statement can be "income", "balance", or "cashflow".
// If quarterly is true, fetches quarterly data instead of annual.
func (c *Client) GetFinancials(symbol, statement string, quarterly bool) ([]FinancialResult, error) {
	types := getFinancialTypes(statement, quarterly)
	if len(types) == 0 {
		return nil, fmt.Errorf("invalid statement type %q (use income, balance, or cashflow)", statement)
	}

	params := url.Values{
		"symbol":  {symbol},
		"type":    {strings.Join(types, ",")},
		"period1": {"493590046"},
		"period2": {fmt.Sprintf("%d", time.Now().Unix())},
	}

	path := fmt.Sprintf("/ws/fundamentals-timeseries/v1/finance/timeseries/%s", url.PathEscape(symbol))
	body, err := c.Get(path, params, true)
	if err != nil {
		return nil, fmt.Errorf("get financials: %w", err)
	}

	return parseFinancialsResponse(body, types)
}

func getFinancialTypes(statement string, quarterly bool) []string {
	switch strings.ToLower(statement) {
	case "income":
		if quarterly {
			return quarterlyIncomeStatementTypes
		}
		return incomeStatementTypes
	case "balance":
		if quarterly {
			return quarterlyBalanceSheetTypes
		}
		return balanceSheetTypes
	case "cashflow", "cash_flow":
		if quarterly {
			return quarterlyCashFlowTypes
		}
		return cashFlowTypes
	default:
		return nil
	}
}

func parseFinancialsResponse(body []byte, types []string) ([]FinancialResult, error) {
	var raw struct {
		Timeseries struct {
			Result []json.RawMessage `json:"result"`
			Error  *YahooError       `json:"error"`
		} `json:"timeseries"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("parsing financials response: %w", err)
	}

	if raw.Timeseries.Error != nil {
		return nil, fmt.Errorf("yahoo error: %s", raw.Timeseries.Error.Description)
	}

	var results []FinancialResult

	for _, rawResult := range raw.Timeseries.Result {
		var resultMap map[string]json.RawMessage
		if err := json.Unmarshal(rawResult, &resultMap); err != nil {
			continue
		}

		// Find which type key exists in this result
		for _, typeName := range types {
			dataRaw, ok := resultMap[typeName]
			if !ok {
				continue
			}

			var items []struct {
				AsOfDate      string `json:"asOfDate"`
				ReportedValue struct {
					Raw float64 `json:"raw"`
					Fmt string  `json:"fmt"`
				} `json:"reportedValue"`
				CurrencyCode string `json:"currencyCode"`
			}
			if err := json.Unmarshal(dataRaw, &items); err != nil {
				continue
			}

			fr := FinancialResult{Type: typeName}
			for _, item := range items {
				fr.Items = append(fr.Items, FinancialItem{
					Date:          item.AsOfDate,
					ReportedValue: item.ReportedValue.Raw,
					CurrencyCode:  item.CurrencyCode,
				})
			}
			results = append(results, fr)
		}
	}

	return results, nil
}
