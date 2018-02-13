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
// Test get mass loss for Single Stock - Long
//
func TestSingleStockGetMaxRisked01(t *testing.T) {

	// Test one stock
	positions := &[]models.Position{
		{
			Symbol:       models.Symbol{ShortName: "SPY"},
			OrgQty:       9,
			CostBasis:    2034.12,
			AvgOpenPrice: 226.12,
		},
	}

	// Get max loss
	loss := SingleStockGetMaxRisked(positions)

	// Verify the data was return as expected
	st.Expect(t, loss, 2034.12)
}

//
// Test get mass loss for Single Stock - Short
//
func TestSingleStockGetMaxRisked02(t *testing.T) {

	// Test one stock
	positions := &[]models.Position{
		{
			Symbol:       models.Symbol{ShortName: "SPY"},
			OrgQty:       -9,
			CostBasis:    -2034.00,
			AvgOpenPrice: 226.12,
		},
	}

	// Get max loss
	loss := SingleStockGetMaxRisked(positions)

	// Verify the data was return as expected
	st.Expect(t, loss, -1.00)
}

//
// Test trade classification - Single stock trade
//
func TestIsSingleStock01(t *testing.T) {

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
	class := IsSingleStock(positions)

	// Verify the data was return as expected
	st.Expect(t, class, true)

}

//
// Test trade classification - Single stock trade
//
func TestIsSingleStock02(t *testing.T) {

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
	class := IsSingleStock(positions)

	// Verify the data was return as expected
	st.Expect(t, class, false)

}

/* End File */
