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
// Test trade classification - Put Debit Spread
//
func TestIsPutDebitSpread01(t *testing.T) {

	// Test put debit spread
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

	// Get the classification of this trade
	class := IsPutDebitSpread(positions)

	// Verify the data was return as expected
	st.Expect(t, class, true)
}

//
// Test trade classification - Put Debit Spread
//
func TestIsPutDebitSpread02(t *testing.T) {

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

	// Get the classification of this trade
	class := IsPutDebitSpread(positions)

	// Verify the data was return as expected
	st.Expect(t, class, false)
}

/* End File */
