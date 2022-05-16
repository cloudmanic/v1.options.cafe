//
// Date: 2022-05-15
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Copyright: 2022 Cloudmanic Labs, LLC. All rights reserved.
//

package helpers

import (
	"math"
	"time"

	"app.options.cafe/brokers/types"
	"app.options.cafe/library/helpers"
	"app.options.cafe/models"
)

//
// CloseTrade will close an open trade.
// tradeId is just the index to the trade group.
//
func CloseTrade(tradeId int, today time.Time, note string, backtest *models.Backtest, options []types.OptionsChainItem, benchmarkQuotes []types.HistoryQuote) {
	// Get benchmark last
	benchmarkLast := GetBenchmarkByDate(today, benchmarkQuotes)

	// Get closing price
	closePrice := (GetClosedPrice(backtest.TradeGroups[tradeId], options) * 100)

	// TradeGroup ending balance
	backtest.EndingBalance = backtest.EndingBalance + closePrice

	// Shared for all strats
	backtest.TradeGroups[tradeId].Note = note
	backtest.TradeGroups[tradeId].Status = "Closed"
	backtest.TradeGroups[tradeId].CloseDate = models.Date{today}
	backtest.TradeGroups[tradeId].ClosePrice = closePrice
	backtest.TradeGroups[tradeId].Balance = backtest.EndingBalance
	backtest.TradeGroups[tradeId].ReturnFromStart = helpers.Round((((backtest.EndingBalance - backtest.StartingBalance) / backtest.StartingBalance) * 100), 2)
	backtest.TradeGroups[tradeId].ReturnPercent = helpers.Round(((closePrice-backtest.TradeGroups[tradeId].OpenPrice)/backtest.TradeGroups[tradeId].OpenPrice)*100, 2)

	// Benchmark stuff
	investedBenchmark := math.Floor(backtest.StartingBalance / backtest.BenchmarkStart)
	investedBenchmarkLeftOver := backtest.StartingBalance - (investedBenchmark * backtest.BenchmarkStart)

	backtest.TradeGroups[tradeId].BenchmarkLast = benchmarkLast
	backtest.TradeGroups[tradeId].BenchmarkBalance = helpers.Round((investedBenchmark*backtest.TradeGroups[tradeId].BenchmarkLast)+investedBenchmarkLeftOver, 2)
	backtest.TradeGroups[tradeId].BenchmarkReturn = helpers.Round((((backtest.TradeGroups[tradeId].BenchmarkLast - backtest.BenchmarkStart) / backtest.BenchmarkStart) * 100), 2)
}

//
// CloseOnExpire will close a trade if it has expired. Designs so all backtests can use this.
//
func CloseOnExpire(today time.Time, underlyingLast float64, backtest *models.Backtest, options []types.OptionsChainItem, benchmarkQuotes []types.HistoryQuote) {
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
			// Get closing price
			closePrice := (GetClosedPrice(row, options) * 100)

			// See if we have a profit.
			profit := closePrice - row.OpenPrice

			// Close the trade
			if profit > 0 {
				CloseTrade(key, today, "Took profit on expire day.", backtest, options, benchmarkQuotes)
			} else {
				CloseTrade(key, today, "Expired worthless.", backtest, options, benchmarkQuotes)
				backtest.TradeGroups[key].ClosePrice = 0.00
				backtest.TradeGroups[key].ReturnPercent = -100.00
			}
		}
	}
}

//
// GetClosedPrice - Figure out how much it would be to close this trade now
//
func GetClosedPrice(tradegroup models.BacktestTradeGroup, options []types.OptionsChainItem) float64 {
	// Total to close
	closePrice := 0.00

	// Loop through the different options and come up with a close price
	for _, row := range tradegroup.Positions {
		// Loop through until we find first symbol
		for _, row2 := range options {
			// Not found.
			if row2.Symbol != row.Symbol.ShortName {
				continue
			}

			// If Short or long?
			if row.Qty > 0 {
				closePrice = closePrice + (row2.Bid * float64(row.Qty))
				break
			} else {
				closePrice = closePrice + (row2.Ask * float64(row.Qty)) // Qty is negative
				break
			}
		}
	}

	return closePrice
}
