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
	OpenCost               float64     `json:"open_cost"`
	CurrentUnderlyingPrice float64     `json:"current_underlying_price"`
	Legs                   []TradeLegs `json:"legs"`
}

type TradeLegs struct {
	SymbolStr string        `json:"symbol_str"`
	Symbol    models.Symbol `json:"symbol"`
	Qty       int           `json:"qty"`
}

type Result struct {
	UnderlyingPrice float64 `json:"underlying_price"`
	Profit          float64 `json:"profit"`
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

		list = append(list, Result{UnderlyingPrice: helpers.Round(i, 2), Profit: (openCost * -1)})

	}

	// Return happy
	return list
}

/* End File */
