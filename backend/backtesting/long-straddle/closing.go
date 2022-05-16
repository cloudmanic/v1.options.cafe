//
// Date: 2022-05-15
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Copyright: 2022 Cloudmanic Labs, LLC. All rights reserved.
//

package longstraddle

import (
	"time"

	"app.options.cafe/brokers/types"
	"app.options.cafe/models"

	btHelpers "app.options.cafe/backtesting/helpers"
)

//
// Close - Looks for positions to close.
//
func Close(today time.Time, underlyingLast float64, backtest *models.Backtest, options []types.OptionsChainItem, benchmarkQuotes []types.HistoryQuote) {
	// Close on expire.
	btHelpers.CloseOnExpire(today, underlyingLast, backtest, options, benchmarkQuotes)

	// Close on profit.
	CloseOnProfit(today, underlyingLast, backtest, options, benchmarkQuotes)
}

//
// CloseOnProfit will look to see if we hit our profit target and close.
//
func CloseOnProfit(today time.Time, underlyingLast float64, backtest *models.Backtest, options []types.OptionsChainItem, benchmarkQuotes []types.HistoryQuote) {
	// Loop for expired postions
	for key, row := range backtest.TradeGroups {
		// If closed moved on.
		if row.Status == "Closed" {
			continue
		}

		// Get closing price
		closePrice := (btHelpers.GetClosedPrice(row, options) * 100)

		// See if we have a profit.
		profit := closePrice - row.OpenPrice

		// Beyond this point we assume we have a trade to close. TODO(spicer): Make this a config
		if profit < 100.00 {
			continue
		}

		// Close the trade
		btHelpers.CloseTrade(key, today, "Triggered at profit target", backtest, options, benchmarkQuotes)
	}
}
