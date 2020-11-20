//
// Date: 9/20/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package trade_types

import (
	"testing"

	"app.options.cafe/models"
	"github.com/nbio/st"
)

//
// Test get mass loss for Single Stock - Long
//
func TestGetSingleShortStockRiskProfile01(t *testing.T) {

	// Test one stock
	positions := &[]models.Position{
		{
			Symbol:       models.Symbol{ShortName: "IBM"},
			OrgQty:       -10,
			CostBasis:    -1495.80,
			AvgOpenPrice: 149.58,
		},
	}

	// Get max loss
	loss, _ := GetSingleShortStockRiskProfile(positions)

	// Verify the data was return as expected
	st.Expect(t, loss, -1.00)
}

//
// Test trade classification - Single stock trade
//
func TestIsSingleShortStock01(t *testing.T) {

	// Test put credit spread
	positions := &[]models.Position{
		{
			Symbol:       models.Symbol{ShortName: "IBM"},
			OrgQty:       -10,
			CostBasis:    -1495.80,
			AvgOpenPrice: 149.58,
		},
	}

	// Get the classification of this trade group
	class := IsSingleShortStock(positions)

	// Verify the data was return as expected
	st.Expect(t, class, true)

}

//
// Test trade classification - Single stock trade
//
func TestIsSingleShortStock02(t *testing.T) {

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
	class := IsSingleShortStock(positions)

	// Verify the data was return as expected
	st.Expect(t, class, false)

}

//
// Test trade classification - Single long stock trade
//
func TestIsSingleShortStock03(t *testing.T) {

	// Test put credit spread
	positions := &[]models.Position{
		{
			Symbol:       models.Symbol{ShortName: "IBM"},
			OrgQty:       10,
			CostBasis:    1495.80,
			AvgOpenPrice: 149.58,
		},
	}

	// Get the classification of this trade group
	class := IsSingleShortStock(positions)

	// Verify the data was return as expected
	st.Expect(t, class, false)

}

/* End File */
