//
// Date: 2/9/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package trade_types

import (
	"testing"

	"app.options.cafe/models"
	"github.com/nbio/st"
)

//
// Test get mass loss for call credit spread
//
func TestGetCallCreditSpreadRiskProfile01(t *testing.T) {

	// Test put credit spread
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

	// Get max loss
	loss, _ := GetCallCreditSpreadRiskProfile(positions)

	// Verify the data was return as expected
	st.Expect(t, loss, 712.00)
}

//
// Test trade classification - Call Credit Spread
//
func TestIsCallCreditSpread01(t *testing.T) {

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
	class := IsCallCreditSpread(positions)

	// Verify the data was return as expected
	st.Expect(t, class, true)

}

//
// Test trade classification - Call Credit Spread
//
func TestIsCallCreditSpread02(t *testing.T) {

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
	class := IsCallCreditSpread(positions)

	// Verify the data was return as expected
	st.Expect(t, class, false)

}

/* End File */
