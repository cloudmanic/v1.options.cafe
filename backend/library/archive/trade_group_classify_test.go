//
// Date: 2/9/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package archive

import (
	"testing"

	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/nbio/st"
)

//
// Test trade classification - Single stock trade
//
func TestClassifyTradeGroupSingleStock01(t *testing.T) {

	// Test put credit spread
	positions := &[]models.Position{
		{
			Symbol:       "SPY",
			OrgQty:       9,
			CostBasis:    2034.00,
			AvgOpenPrice: 226.12,
		},
	}

	// Get the classification of this trade group
	class := ClassifyTradeGroup(positions)

	// Verify the data was return as expected
	st.Expect(t, class, "Stock")

}

//
// Test trade classification - Single option trade
//
func TestClassifyTradeGroupSingleOption01(t *testing.T) {

	// Test put credit spread
	positions := &[]models.Position{
		{
			Symbol:       "SPY180221P00241000",
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
			Symbol:       "SPY180221P00241000",
			OrgQty:       9,
			CostBasis:    2034.00,
			AvgOpenPrice: 2.26,
		},

		{
			Symbol:       "SPY180221P00243000",
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
			Symbol:       "SPY180221P00241000",
			OrgQty:       -9,
			CostBasis:    -2034.00,
			AvgOpenPrice: 2.26,
		},

		{
			Symbol:       "SPY180221P00243000",
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
			Symbol:       "SPY180221C00241000",
			OrgQty:       9,
			CostBasis:    2034.00,
			AvgOpenPrice: 2.26,
		},

		{
			Symbol:       "SPY180221C00243000",
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
			Symbol:       "SPY180221C00241000",
			OrgQty:       -9,
			CostBasis:    -2034.00,
			AvgOpenPrice: 2.26,
		},

		{
			Symbol:       "SPY180221C00243000",
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
