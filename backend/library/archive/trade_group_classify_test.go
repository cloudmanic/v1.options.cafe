//
// Date: 2/9/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package archive

import (
	"testing"

	"app.options.cafe/models"
	"github.com/nbio/st"
)

//
// Test trade risked in trade - Single stock trade
//
func TestGetAmountRiskedInTrade01(t *testing.T) {
	// Test put credit spread
	positions := &[]models.Position{
		{
			Symbol:       models.Symbol{ShortName: "SPY"},
			OrgQty:       9,
			CostBasis:    2034.00,
			AvgOpenPrice: 226.12,
		},
	}

	// Get the classification of this trade group
	class, _ := GetAmountRiskedInTrade(positions)

	// Verify the data was return as expected
	st.Expect(t, class, 2034.00)
}

//
// Test trade risked in trade - Single option trade
//
func TestGetAmountRiskedInTrade02(t *testing.T) {
	// Test single option
	positions := &[]models.Position{
		{
			Symbol:       models.Symbol{ShortName: "VXX180223C00055000"},
			OrgQty:       2,
			CostBasis:    982.00,
			AvgOpenPrice: 4.91,
		},
	}

	// Get the classification of this trade group
	class, _ := GetAmountRiskedInTrade(positions)

	// Verify the data was return as expected
	st.Expect(t, class, 982.00)
}

//
// Test trade risked in trade - Put Credit Spread trade
//
func TestGetAmountRiskedInTrade03(t *testing.T) {
	// Test put credit spread
	positions := &[]models.Position{
		{
			Symbol:       models.Symbol{ShortName: "SPY180221P00241000"},
			OrgQty:       9,
			CostBasis:    2034.00,
			AvgOpenPrice: 2.26,
		},

		{
			Symbol:       models.Symbol{ShortName: "SPY180221P00243000"},
			OrgQty:       -9,
			CostBasis:    -2286.00,
			AvgOpenPrice: 2.54,
		},
	}

	// Get the classification of this trade group
	class, _ := GetAmountRiskedInTrade(positions)

	// Verify the data was return as expected
	st.Expect(t, class, 1548.00)
}

//
// Test trade risked in trade - Call Credit Spread trade
//
func TestGetAmountRiskedInTrade04(t *testing.T) {
	// Test call credit spread
	positions := &[]models.Position{
		{
			Symbol:       models.Symbol{ShortName: "VXX180223C00055000"},
			OrgQty:       2,
			CostBasis:    982.00,
			AvgOpenPrice: 4.91,
		},

		{
			Symbol:       models.Symbol{ShortName: "VXX180223C00050000"},
			OrgQty:       -2,
			CostBasis:    -1270.00,
			AvgOpenPrice: 6.35,
		},
	}

	// Get the classification of this trade group
	class, _ := GetAmountRiskedInTrade(positions)

	// Verify the data was return as expected
	st.Expect(t, class, 712.00)
}

//
// Test trade classification - Single stock trade
//
func TestClassifyTradeGroupSingleStock01(t *testing.T) {
	// Test put credit spread
	positions := &[]models.Position{
		{
			Symbol:       models.Symbol{ShortName: "SPY"},
			OrgQty:       9,
			CostBasis:    2034.00,
			AvgOpenPrice: 226.12,
		},
	}

	// Get the classification of this trade group
	class := ClassifyTradeGroup(positions)

	// Verify the data was return as expected
	st.Expect(t, class, "Equity")
}

//
// Test trade classification - Single short stock trade
//
func TestClassifyTradeGroupSingleShortStock01(t *testing.T) {
	// Test put credit spread
	positions := &[]models.Position{
		{
			Symbol:       models.Symbol{ShortName: "SPY"},
			OrgQty:       -9,
			CostBasis:    -2034.00,
			AvgOpenPrice: 226.12,
		},
	}

	// Get the classification of this trade group
	class := ClassifyTradeGroup(positions)

	// Verify the data was return as expected
	st.Expect(t, class, "Short Equity")
}

//
// Test trade classification - Single option trade
//
func TestClassifyTradeGroupSingleOption01(t *testing.T) {
	// Test put credit spread
	positions := &[]models.Position{
		{
			Symbol:       models.Symbol{ShortName: "SPY180221P00241000"},
			OrgQty:       9,
			CostBasis:    2034.00,
			AvgOpenPrice: 2.26,
		},
	}

	// Get the classification of this trade group
	class := ClassifyTradeGroup(positions)

	// Verify the data was return as expected
	st.Expect(t, class, "Option")
}

//
// Test trade classification - Put Credit Spread
//
func TestClassifyTradeGroupPutCreditSpread01(t *testing.T) {
	// Test put credit spread
	positions := &[]models.Position{
		{
			Symbol:       models.Symbol{ShortName: "SPY180221P00241000"},
			OrgQty:       9,
			CostBasis:    2034.00,
			AvgOpenPrice: 2.26,
		},

		{
			Symbol:       models.Symbol{ShortName: "SPY180221P00243000"},
			OrgQty:       -9,
			CostBasis:    -2286.00,
			AvgOpenPrice: 2.54,
		},
	}

	// Get the classification of this trade group
	class := ClassifyTradeGroup(positions)

	// Verify the data was return as expected
	st.Expect(t, class, "Put Credit Spread")
}

//
// Test trade classification - Put Debit Spread
//
func TestClassifyTradeGroupPutDebitSpread01(t *testing.T) {
	// Test put credit spread
	positions := &[]models.Position{
		{
			Symbol:       models.Symbol{ShortName: "SPY180221P00241000"},
			OrgQty:       -9,
			CostBasis:    -2034.00,
			AvgOpenPrice: 2.26,
		},

		{
			Symbol:       models.Symbol{ShortName: "SPY180221P00243000"},
			OrgQty:       9,
			CostBasis:    2286.00,
			AvgOpenPrice: 2.54,
		},
	}

	// Get the classification of this trade group
	class := ClassifyTradeGroup(positions)

	// Verify the data was return as expected
	st.Expect(t, class, "Put Debit Spread")
}

//
// Test trade classification - Call Credit Spread
//
func TestClassifyTradeGroupCallCreditSpread01(t *testing.T) {
	// Test put credit spread
	positions := &[]models.Position{
		{
			Symbol:       models.Symbol{ShortName: "SPY180221C00241000"},
			OrgQty:       9,
			CostBasis:    2034.00,
			AvgOpenPrice: 2.26,
		},

		{
			Symbol:       models.Symbol{ShortName: "SPY180221C00243000"},
			OrgQty:       -9,
			CostBasis:    -2286.00,
			AvgOpenPrice: 2.54,
		},
	}

	// Get the classification of this trade group
	class := ClassifyTradeGroup(positions)

	// Verify the data was return as expected
	st.Expect(t, class, "Call Credit Spread")
}

//
// Test trade classification - Call Debit Spread
//
func TestClassifyTradeGroupCallDebitSpread01(t *testing.T) {
	// Test put credit spread
	positions := &[]models.Position{
		{
			Symbol:       models.Symbol{ShortName: "SPY180221C00241000"},
			OrgQty:       -9,
			CostBasis:    -2034.00,
			AvgOpenPrice: 2.26,
		},

		{
			Symbol:       models.Symbol{ShortName: "SPY180221C00243000"},
			OrgQty:       9,
			CostBasis:    2286.00,
			AvgOpenPrice: 2.54,
		},
	}

	// Get the classification of this trade group
	class := ClassifyTradeGroup(positions)

	// Verify the data was return as expected
	st.Expect(t, class, "Call Debit Spread")
}

/* End File */
