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
// Test get max loss for Long Put Butterfly
//
func TestGetLongPutButterflyRiskProfile01(t *testing.T) {

	// Test put credit spread
	positions := &[]models.Position{
		{
			Symbol:       models.Symbol{ShortName: "SPY181116P00275000"}, // SPY Nov 16 2018 $275.00 Put
			OrgQty:       1,
			CostBasis:    323.00,
			AvgOpenPrice: 3.23,
		},

		{
			Symbol:       models.Symbol{ShortName: "SPY181116P00270000"}, // SPY Nov 16 2018 $270.00 Put
			OrgQty:       -2,
			CostBasis:    -154.00,
			AvgOpenPrice: 0.77,
		},

		{
			Symbol:       models.Symbol{ShortName: "SPY181116P00265000"}, // SPY Nov 16, 2018 $265.00 Put
			OrgQty:       1,
			CostBasis:    15.00,
			AvgOpenPrice: 0.15,
		},
	}

	// Get max loss
	loss, _ := GetLongPutButterflyRiskProfile(positions)

	// Verify the data was return as expected
	st.Expect(t, loss, 184.00)
}

//
// Test trade classification - Long Put Butterfly
//
func TestIsLongPutButterfly01(t *testing.T) {

	// Test put credit spread
	positions := &[]models.Position{
		{
			Symbol:       models.Symbol{ShortName: "SPY181116P00275000"}, // SPY Nov 16 2018 $275.00 Put
			OrgQty:       1,
			CostBasis:    323.00,
			AvgOpenPrice: 3.23,
		},

		{
			Symbol:       models.Symbol{ShortName: "SPY181116P00270000"}, // SPY Nov 16 2018 $270.00 Put
			OrgQty:       -2,
			CostBasis:    -154.00,
			AvgOpenPrice: 0.77,
		},

		{
			Symbol:       models.Symbol{ShortName: "SPY181116P00265000"}, // SPY Nov 16, 2018 $265.00 Put
			OrgQty:       1,
			CostBasis:    15.00,
			AvgOpenPrice: 0.15,
		},
	}

	// Get the classification of this trade group
	class := IsLongPutButterfly(positions)

	// Verify the data was return as expected
	st.Expect(t, class, true)

}

//
// Test trade classification - Long Put Butterfly
//
func TestIsLongPutButterfly02(t *testing.T) {

	// Test put credit spread
	positions := &[]models.Position{
		{
			Symbol:       models.Symbol{ShortName: "SPY181221C00250000"}, // SPY Dec 21, 2018 $250.00 Call
			OrgQty:       1,
			CostBasis:    577.00,
			AvgOpenPrice: 5.77,
		},

		{
			Symbol:       models.Symbol{ShortName: "SPY181221C00260000"}, // SPY Dec 21, 2018 $260.00 Call
			OrgQty:       -2,
			CostBasis:    -2480.00,
			AvgOpenPrice: 12.40,
		},

		{
			Symbol:       models.Symbol{ShortName: "SPY181221C00270000"}, // SPY Dec 21, 2018 $270.00 Call
			OrgQty:       1,
			CostBasis:    2060.00,
			AvgOpenPrice: 20.60,
		},
	}

	// Get the classification of this trade group
	class := IsLongPutButterfly(positions)

	// Verify the data was return as expected
	st.Expect(t, class, false)

}

/* End File */
