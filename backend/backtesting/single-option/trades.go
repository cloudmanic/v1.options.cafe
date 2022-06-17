//
// Date: 2022-05-15
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Copyright: 2022 Cloudmanic Labs, LLC. All rights reserved.
//

package singleoption

import (
	"time"

	"app.options.cafe/backtesting/helpers"
	"app.options.cafe/brokers/types"
	"app.options.cafe/models"
	"app.options.cafe/screener"
)

//
// Trades will place our trades
//
func Trades(db models.Datastore, today time.Time, backtest *models.Backtest, results []screener.Result, options []types.OptionsChainItem, underlyingLast float64, benchmarkQuotes []types.HistoryQuote) {
	// First look if there are any trades we need to close.
	Close(today, underlyingLast, backtest, options, benchmarkQuotes)

	// Figure out qtys. QTY field in symbol is not stored in DB. Just used for this purpose
	results[0].Legs[0].Qty = 1

	// Open the trade - Only open one at a time. TODO(spicer): make this a config
	if helpers.GetOpenTradeCount(backtest) == 0 {
		helpers.OpenTrade(today, backtest, results[0])
	}

	return
}
