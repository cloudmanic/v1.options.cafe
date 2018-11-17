//
// Date: 2018-11-16
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-11-17
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package analyze

import (
	"github.com/cloudmanic/app.options.cafe/backend/library/helpers"
)

const NumProfitLossByUnderlyingPricePoints = 500
const NumProfitLossByUnderlyingPriceRangePercent = .05

//
// Analyze and options trade. We pass in an array
// of symbols and return a list of underlying prices based on
// profit and loss.
//
func OptionsProfitLossByUnderlyingPrice(trade Trade) []Result {

	// Get min and max strikes
	minStrike, maxStrike := GetMinMaxStrikesTradeLegs(trade)

	// Figure out our result range
	startRange := minStrike - (minStrike * NumProfitLossByUnderlyingPriceRangePercent)
	endRange := maxStrike * (1 + NumProfitLossByUnderlyingPriceRangePercent)

	// Get a list of underlying results with prices to compare our options against
	results := GetRangeOfUnderlyingResults(startRange, endRange, trade.OpenCost, NumProfitLossByUnderlyingPricePoints)

	// Loop through the different legs and get results based on price list
	for _, row := range trade.Legs {

		for key2 := range results {

			// Get the profit / Loss for this leg
			OptionsProfitLossByUnderlyingPriceManageLeg(row, results, key2)

		}

	}

	// Return happy
	return results
}

//
// Manage Leg
//
func OptionsProfitLossByUnderlyingPriceManageLeg(leg TradeLegs, results []Result, key int) {

	// Is this a long call?
	if (leg.Symbol.OptionType == "Call") && (leg.Qty > 0) {

		// If the OptionStrike less than the UnderlyingPrice we have a profit otherwise expired worthless
		if leg.Symbol.OptionStrike < results[key].UnderlyingPrice {
			results[key].Profit = helpers.Round(results[key].Profit+(((results[key].UnderlyingPrice-leg.Symbol.OptionStrike)*100.00)*float64(leg.Qty)), 2)
		}

	}

	// Is this a short call?
	if (leg.Symbol.OptionType == "Call") && (leg.Qty < 0) {

		// If the OptionStrike is less than the UnderlyingPrice we have a loss otherwise expired worthless. (Reminder: QTY will be negative)
		if leg.Symbol.OptionStrike < results[key].UnderlyingPrice {
			results[key].Profit = helpers.Round(results[key].Profit+(((results[key].UnderlyingPrice-leg.Symbol.OptionStrike)*100.00)*float64(leg.Qty)), 2)
		}

	}

	// Is this a long put?
	if (leg.Symbol.OptionType == "Put") && (leg.Qty > 0) {

		// If the Underlying Price is less than the Option Strike we have a gain
		if leg.Symbol.OptionStrike > results[key].UnderlyingPrice {
			results[key].Profit = helpers.Round(results[key].Profit+(((leg.Symbol.OptionStrike-results[key].UnderlyingPrice)*100.00)*float64(leg.Qty)), 2)
		}

	}

	// Is this a short put?
	if (leg.Symbol.OptionType == "Put") && (leg.Qty < 0) {

		// If the Underlying Price is less than the Option Strike we have a loss (Reminder: QTY will be negative)
		if leg.Symbol.OptionStrike > results[key].UnderlyingPrice {
			results[key].Profit = helpers.Round(results[key].Profit+((leg.Symbol.OptionStrike-results[key].UnderlyingPrice)*100*float64(leg.Qty)), 2)
		}

	}

}

/* End File */
