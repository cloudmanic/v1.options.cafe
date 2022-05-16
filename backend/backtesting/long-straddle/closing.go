//
// Date: 2022-05-15
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Copyright: 2022 Cloudmanic Labs, LLC. All rights reserved.
//

package longstraddle

import (
	"os"
	"time"

	"app.options.cafe/brokers/types"
	"app.options.cafe/library/helpers"
	"app.options.cafe/models"
	"github.com/davecgh/go-spew/spew"
)

//
// Close - Looks for positions to close.
//
func Close(today time.Time, underlyingLast float64, backtest *models.Backtest, options []types.OptionsChainItem) {

	// Loop for expired postions
	for key, row := range backtest.TradeGroups {
		// If closed moved on.
		if row.Status == "Closed" {
			continue
		}

		expired := false

		// See if any of the legs are expired
		for _, row2 := range row.Positions {
			if today.Format("2006-01-02") == row2.Symbol.OptionExpire.Format("2006-01-02") ||
				today.After(helpers.ParseDateNoError(row2.Symbol.OptionExpire.Format("2006-01-02"))) {
				expired = true
				break
			}
		}

		// If expired close out trade
		if expired && row.Status == "Open" {

			spew.Dump(row)
			spew.Dump(today)
			spew.Dump(key)
			os.Exit(1)

			// // Benchmark stuff
			// investedBenchmark := math.Floor(backtest.StartingBalance / backtest.BenchmarkStart)
			// investedBenchmarkLeftOver := backtest.StartingBalance - (investedBenchmark * backtest.BenchmarkStart)

			// // Shared for all strats
			// backtest.TradeGroups[key].Status = "Closed"
			// backtest.TradeGroups[key].CloseDate = models.Date{today}
			// backtest.TradeGroups[key].BenchmarkLast = benchmarkLast
			// backtest.TradeGroups[key].ReturnFromStart = helpers.Round((((backtest.EndingBalance - backtest.StartingBalance) / backtest.StartingBalance) * 100), 2)
			// backtest.TradeGroups[key].BenchmarkBalance = helpers.Round((investedBenchmark*backtest.TradeGroups[key].BenchmarkLast)+investedBenchmarkLeftOver, 2)
			// backtest.TradeGroups[key].BenchmarkReturn = helpers.Round((((backtest.TradeGroups[key].BenchmarkLast - backtest.BenchmarkStart) / backtest.BenchmarkStart) * 100), 2)

			// // Get closing price
			// closePrice := (t.getClosedPrice(row, options) * 100)

			// // See if we have a profit.
			// profit := closePrice - row.OpenPrice

			// if profit > 0 {
			// 	backtest.EndingBalance = backtest.EndingBalance + closePrice
			// 	backtest.TradeGroups[key].ClosePrice = closePrice
			// 	backtest.TradeGroups[key].Note = "Took profit on expire day."
			// 	backtest.TradeGroups[key].Balance = backtest.EndingBalance
			// 	backtest.TradeGroups[key].BenchmarkLast = benchmarkLast
			// 	backtest.TradeGroups[key].ReturnFromStart = helpers.Round((((backtest.EndingBalance - backtest.StartingBalance) / backtest.StartingBalance) * 100), 2)
			// 	backtest.TradeGroups[key].ReturnPercent = helpers.Round(((closePrice-row.OpenPrice)/row.OpenPrice)*100, 2)
			// } else {
			// 	backtest.TradeGroups[key].ClosePrice = 0.00
			// 	backtest.TradeGroups[key].Note = "Expired worthless."
			// 	backtest.TradeGroups[key].ReturnPercent = -100.00
			// }
		}
	}

	// // Close if we touch the short leg
	// t.closeOnShortTouch(today, underlyingLast, backtest, options)

	// Close if we hit a particular debit
	//t.CloseOnDebit(today, backtest, options)

	// Expire positions
	//t.LongCallButterflySpreadExpirePositions(today, underlyingLast, backtest, options)
}
