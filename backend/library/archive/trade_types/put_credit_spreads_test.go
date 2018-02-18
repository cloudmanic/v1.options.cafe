//
// Date: 2/9/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package trade_types

import (
	"testing"

	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/nbio/st"
)

//
// Test get mass loss for put credit spread
//
func TestGetPutCreditSpreadRiskProfile01(t *testing.T) {

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

	// Get max loss
	loss, credit := GetPutCreditSpreadRiskProfile(positions)

	// Verify the data was return as expected
	st.Expect(t, loss, 1548.00)
	st.Expect(t, credit, 252.00)
}

//
// Test trade classification - Put Credit Spread
//
func TestIsPutCreditSpread01(t *testing.T) {

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
	class := IsPutCreditSpread(positions)

	// Verify the data was return as expected
	st.Expect(t, class, true)
}

//
// Test trade classification - Put Credit Spread
//
func TestIsPutCreditSpread02(t *testing.T) {

	// Test put debit spread
	positions := &[]models.Position{
		{
			Symbol:       models.Symbol{ShortName: "SPY180221P00241000"},
			OrgQty:       9,
			CostBasis:    2034.00,
			AvgOpenPrice: 2.26,
		},

		{
			Symbol:       models.Symbol{ShortName: "SPY180221P00239000"},
			OrgQty:       -9,
			CostBasis:    -1886.00,
			AvgOpenPrice: 1.99,
		},
	}

	// Get the classification of this trade group
	class := IsPutCreditSpread(positions)

	// Verify the data was return as expected
	st.Expect(t, class, false)
}

/* End File */
