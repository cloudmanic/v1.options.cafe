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
// Test get mass loss for Single Option - Long
//
func TestSingleOptionGetMaxRisked01(t *testing.T) {

	// Test put credit spread
	positions := &[]models.Position{
		{
			Symbol:       models.Symbol{ShortName: "VXX180223C00055000"},
			OrgQty:       2,
			CostBasis:    982.00,
			AvgOpenPrice: 4.91,
		},
	}

	// Get max loss
	loss := SingleOptionGetMaxRisked(positions)

	// Verify the data was return as expected
	st.Expect(t, loss, 982.00)
}

//
// Test get mass loss for Single Option - Short
//
func TestSingleOptionGetMaxRisked02(t *testing.T) {

	// Test put credit spread
	positions := &[]models.Position{
		{
			Symbol:       models.Symbol{ShortName: "VXX180223C00050000"},
			OrgQty:       -2,
			CostBasis:    -1270.00,
			AvgOpenPrice: 6.35,
		},
	}

	// Get max loss
	loss := SingleOptionGetMaxRisked(positions)

	// Verify the data was return as expected
	st.Expect(t, loss, -1.00)
}

//
// Test trade classification - Single option trade
//
func TestIsSingleOption01(t *testing.T) {

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
	class := IsSingleOption(positions)

	// Verify the data was return as expected
	st.Expect(t, class, true)

}

//
// Test trade classification - Single option trade
//
func TestIsSingleOption02(t *testing.T) {

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
	class := IsSingleOption(positions)

	// Verify the data was return as expected
	st.Expect(t, class, false)

}

/* End File */
