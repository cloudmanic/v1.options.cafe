//
// Date: 2018-11-16
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-11-16
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package analyze

import (
	"github.com/cloudmanic/app.options.cafe/backend/library/helpers"
	"github.com/cloudmanic/app.options.cafe/backend/models"
)

type Trade struct {
	OpenCost               float64
	CurrentUnderlyingPrice float64
	Legs                   []TradeLegs
}

type TradeLegs struct {
	Symbol models.Symbol
	Qty    int
}

type Result struct {
	UnderlyingPrice float64
	Profit          float64
}

//
// Return the min and max values of an array of strikes.
//
func GetMinMaxStrikesTradeLegs(trade Trade) (float64, float64) {

	strikes := []float64{}

	// Get an array of Strikes
	for _, row := range trade.Legs {
		strikes = append(strikes, row.Symbol.OptionStrike)
	}

	minStrike := helpers.MinFloat64Slice(strikes)
	maxStrike := helpers.MaxFloat64Slice(strikes)

	// Return happy
	return minStrike, maxStrike
}

//
// Return an array of ending underlying prices for us to compare our P&L against
//
func GetRangeOfUnderlyingResults(min float64, max float64, openCost float64, pointCount float64) []Result {

	list := []Result{}

	// Figure out the range step
	rangeStep := (max - min) / NumProfitLossByUnderlyingPricePoints

	// Loop through our range.
	for i := min; i <= max; i = i + rangeStep {

		list = append(list, Result{UnderlyingPrice: i, Profit: (openCost * -1)})

	}

	// Return happy
	return list
}

/* End File */
