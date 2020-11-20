//
// Date: 2018-11-16
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-11-17
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package analyze

import (
	"testing"

	"app.options.cafe/library/helpers"
	"app.options.cafe/models"
	"github.com/nbio/st"
)

//
// Test - OptionsProfitLossByUnderlyingPrice - 01 (Long Call Butterfly)
//
func TestOptionsProfitLossByUnderlyingPrice01(t *testing.T) {

	syb1 := models.Symbol{
		ShortName:        "SPY181221C00250000",
		Name:             "SPY Dec 21, 2018 $250.00 Call",
		Type:             "Option",
		OptionUnderlying: "SPY",
		OptionType:       "Call",
		OptionExpire:     models.Date{helpers.ParseDateNoError("2018-12-21")},
		OptionStrike:     250,
	}

	syb2 := models.Symbol{
		ShortName:        "SPY181221C00260000",
		Name:             "SPY Dec 21, 2018 $260.00 Call",
		Type:             "Option",
		OptionUnderlying: "SPY",
		OptionType:       "Call",
		OptionExpire:     models.Date{helpers.ParseDateNoError("2018-12-21")},
		OptionStrike:     260,
	}

	syb3 := models.Symbol{
		ShortName:        "SPY181221C00270000",
		Name:             "SPY Dec 21, 2018 $270.00 Call",
		Type:             "Option",
		OptionUnderlying: "SPY",
		OptionType:       "Call",
		OptionExpire:     models.Date{helpers.ParseDateNoError("2018-12-21")},
		OptionStrike:     270,
	}

	// Long Call Butterfly
	legs := []TradeLegs{
		{Symbol: syb1, Qty: 1},
		{Symbol: syb2, Qty: -2},
		{Symbol: syb3, Qty: 1},
	}

	// Get the Profit and Loss By Underlying Price
	results := OptionsProfitLossByUnderlyingPrice(Trade{
		OpenCost: 157.00,
		Legs:     legs,
	})

	// Test results
	st.Expect(t, len(results), 501)
	st.Expect(t, results[100].Profit, -157.00)
	st.Expect(t, helpers.Round(results[100].UnderlyingPrice, 2), 246.70)
	st.Expect(t, results[300].Profit, 333.00)
	st.Expect(t, helpers.Round(results[300].UnderlyingPrice, 2), 265.10)

}

//
// Test - OptionsProfitLossByUnderlyingPrice - 02 (Put Credit Spread)
//
func TestOptionsProfitLossByUnderlyingPrice02(t *testing.T) {

	syb1 := models.Symbol{
		ShortName:        "SPY181214P00260000",
		Name:             "SPY Dec 14 2018 $260.00 Put",
		Type:             "Option",
		OptionUnderlying: "SPY",
		OptionType:       "Put",
		OptionExpire:     models.Date{helpers.ParseDateNoError("2018-12-14")},
		OptionStrike:     260,
	}

	syb2 := models.Symbol{
		ShortName:        "SPY181214P00262000",
		Name:             "SPY Dec 14 2018 $262.00 Put",
		Type:             "Option",
		OptionUnderlying: "SPY",
		OptionType:       "Put",
		OptionExpire:     models.Date{helpers.ParseDateNoError("2018-12-14")},
		OptionStrike:     262,
	}

	// Long Call Butterfly
	legs := []TradeLegs{
		{Symbol: syb1, Qty: 11},
		{Symbol: syb2, Qty: -11},
	}

	// Get the Profit and Loss By Underlying Price
	results := OptionsProfitLossByUnderlyingPrice(Trade{
		OpenCost: -286.00,
		Legs:     legs,
	})

	// Test results
	st.Expect(t, len(results), 501)
	st.Expect(t, results[100].Profit, -1913.99)
	st.Expect(t, helpers.Round(results[100].UnderlyingPrice, 2), 252.62)
	st.Expect(t, results[300].Profit, 286.00)
	st.Expect(t, helpers.Round(results[300].UnderlyingPrice, 2), 263.86)

}

/* End File */
