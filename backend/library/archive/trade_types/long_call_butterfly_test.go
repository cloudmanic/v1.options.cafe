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
// Test get max loss for Long Call Butterfly
//
func TestGetLongCallButterflyRiskProfile01(t *testing.T) {

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

	// Get max loss
	loss, _ := GetLongCallButterflyRiskProfile(positions)

	// Verify the data was return as expected
	st.Expect(t, loss, 157.00)
}

//
// Test trade classification - Long Call Butterfly
//
func TestIsLongCallButterfly01(t *testing.T) {

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
	class := IsLongCallButterfly(positions)

	// Verify the data was return as expected
	st.Expect(t, class, true)

}

//
// Test trade classification - Long Call Butterfly
//
func TestIsLongCallButterfly02(t *testing.T) {

	// Test put credit spread
	positions := &[]models.Position{

		{
			Symbol:       models.Symbol{ShortName: "SPY181221C00260000"}, // SPY Dec 21, 2018 $260.00 Call
			OrgQty:       -2,
			CostBasis:    -2480.00,
			AvgOpenPrice: 12.40,
		},

		{
			Symbol:       models.Symbol{ShortName: "SPY181221C00250000"}, // SPY Dec 21, 2018 $250.00 Call
			OrgQty:       1,
			CostBasis:    577.00,
			AvgOpenPrice: 5.77,
		},

		{
			Symbol:       models.Symbol{ShortName: "SPY181221C00270000"}, // SPY Dec 21, 2018 $270.00 Call
			OrgQty:       1,
			CostBasis:    2060.00,
			AvgOpenPrice: 20.60,
		},
	}

	// Get the classification of this trade group
	class := IsLongCallButterfly(positions)

	// Verify the data was return as expected
	st.Expect(t, class, true)

}

//
// Test trade classification - Long Call Butterfly
//
func TestIsLongCallButterfly03(t *testing.T) {

	// Test put credit spread
	positions := &[]models.Position{

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

		{
			Symbol:       models.Symbol{ShortName: "SPY181221C00250000"}, // SPY Dec 21, 2018 $250.00 Call
			OrgQty:       1,
			CostBasis:    577.00,
			AvgOpenPrice: 5.77,
		},
	}

	// Get the classification of this trade group
	class := IsLongCallButterfly(positions)

	// Verify the data was return as expected
	st.Expect(t, class, true)

}

//
// Test trade classification - Long Call Butterfly
//
func TestIsLongCallButterfly04(t *testing.T) {

	// Test put credit spread
	positions := &[]models.Position{

		{
			Symbol:       models.Symbol{ShortName: "SPY181221C00260000"}, // SPY Dec 21, 2018 $260.00 Call
			OrgQty:       2,
			CostBasis:    2480.00,
			AvgOpenPrice: 12.40,
		},

		{
			Symbol:       models.Symbol{ShortName: "SPY181221C00270000"}, // SPY Dec 21, 2018 $270.00 Call
			OrgQty:       1,
			CostBasis:    2060.00,
			AvgOpenPrice: 20.60,
		},

		{
			Symbol:       models.Symbol{ShortName: "SPY181221C00250000"}, // SPY Dec 21, 2018 $250.00 Call
			OrgQty:       1,
			CostBasis:    577.00,
			AvgOpenPrice: 5.77,
		},
	}

	// Get the classification of this trade group
	class := IsCallCreditSpread(positions)

	// Verify the data was return as expected
	st.Expect(t, class, false)

}

//
// Test trade classification - Long Call Butterfly
//
func TestIsLongCallButterfly05(t *testing.T) {

	// Test put credit spread
	positions := &[]models.Position{
		{
			Symbol:       models.Symbol{ShortName: "SPY181221C00250000"}, // SPY Dec 21, 2018 $250.00 Call
			OrgQty:       2,
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
	class := IsLongCallButterfly(positions)

	// Verify the data was return as expected
	st.Expect(t, class, false)

}

/* End File */
