//
// Date: 2019-02-22
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package putcreditspread

import (
	"math"
	"time"

	"app.options.cafe/brokers/types"
	"app.options.cafe/library/helpers"
	"app.options.cafe/models"

	bthelpers "app.options.cafe/backtesting/helpers"
)

//
// CloseTrades - Close positions
//
func CloseTrades(today time.Time, underlyingLast float64, backtest *models.Backtest, options []types.OptionsChainItem, benchmarkQuotes []types.HistoryQuote) {
	// Close if we touch the short leg
	closeOnShortTouch(today, underlyingLast, backtest, options, benchmarkQuotes)

	// Close if we hit a particular debit
	closeOnDebit(today, underlyingLast, backtest, options, benchmarkQuotes)

	// Expire positions
	expirePositions(today, underlyingLast, backtest, benchmarkQuotes)
}

//
// closeOnDebit - Close a trade if it hits our debit trigger
//
func closeOnDebit(today time.Time, underlyingLast float64, backtest *models.Backtest, options []types.OptionsChainItem, benchmarkQuotes []types.HistoryQuote) {
	// TODO(spicer): make this work from configs
	debitAmount := 0.03

	// Get the benchmark last
	benchmarkLast := getBenchmarkByDate(today, benchmarkQuotes)

	// Loop for expired postions
	for key, row := range backtest.TradeGroups {
		// If we are closed moved along
		if row.Status == "Closed" {
			continue
		}

		// Get closing price
		closePrice := bthelpers.GetClosedPrice(row, options) * -1

		// Benchmark stuff
		investedBenchmark := math.Floor(backtest.StartingBalance / backtest.BenchmarkStart)
		investedBenchmarkLeftOver := backtest.StartingBalance - (investedBenchmark * backtest.BenchmarkStart)

		// Close trade at the debitAmount
		if closePrice <= debitAmount {
			// We assume we closed at the debit price not the EOD price
			closePrice = debitAmount

			backtest.EndingBalance = (backtest.EndingBalance - backtest.TradeGroups[key].ClosePrice)
			backtest.TradeGroups[key].Status = "Closed"
			backtest.TradeGroups[key].ClosePrice = closePrice * 100 * float64(row.Lots)
			backtest.TradeGroups[key].CloseDate = models.Date{today}
			backtest.TradeGroups[key].Note = "Triggered at debit amount."
			backtest.TradeGroups[key].Balance = backtest.EndingBalance
			backtest.TradeGroups[key].BenchmarkLast = benchmarkLast
			backtest.TradeGroups[key].ReturnFromStart = helpers.Round((((backtest.EndingBalance - backtest.StartingBalance) / backtest.StartingBalance) * 100), 2)
			backtest.TradeGroups[key].ReturnPercent = helpers.Round((((backtest.TradeGroups[key].Margin + (backtest.TradeGroups[key].OpenPrice - backtest.TradeGroups[key].ClosePrice) - backtest.TradeGroups[key].Margin) / backtest.TradeGroups[key].Margin) * 100), 2)
			backtest.TradeGroups[key].BenchmarkBalance = helpers.Round((investedBenchmark*backtest.TradeGroups[key].BenchmarkLast)+investedBenchmarkLeftOver, 2)
			backtest.TradeGroups[key].BenchmarkReturn = helpers.Round((((backtest.TradeGroups[key].BenchmarkLast - backtest.BenchmarkStart) / backtest.BenchmarkStart) * 100), 2)
		}
	}

}

//
// closeOnShortTouch - If our trade touches the short leg we close
//
func closeOnShortTouch(today time.Time, underlyingLast float64, backtest *models.Backtest, options []types.OptionsChainItem, benchmarkQuotes []types.HistoryQuote) {
	// Get the benchmark last
	benchmarkLast := getBenchmarkByDate(today, benchmarkQuotes)

	// Loop for expired postions
	for key, row := range backtest.TradeGroups {

		if row.Status == "Closed" {
			continue
		}

		// TODO(Spicer): Currently this only works for PCS. We assume the second leg is the short leg.
		if underlyingLast <= row.Positions[1].Symbol.OptionStrike {

			// Set closing Closing Price
			closingPrice := (getClosedPrice(row, options) * 100.00 * float64(row.Lots)) * -1

			// Make sure close price is not bigger than our max risk
			if closingPrice > row.Margin {
				continue
			}

			// Benchmark stuff
			investedBenchmark := math.Floor(backtest.StartingBalance / backtest.BenchmarkStart)
			investedBenchmarkLeftOver := backtest.StartingBalance - (investedBenchmark * backtest.BenchmarkStart)

			// Close trade
			backtest.TradeGroups[key].Status = "Closed"
			backtest.TradeGroups[key].ClosePrice = closingPrice
			backtest.TradeGroups[key].CloseDate = models.Date{today}
			backtest.TradeGroups[key].Note = "Trade touched the short leg."
			backtest.TradeGroups[key].Balance -= closingPrice
			backtest.EndingBalance -= closingPrice
			backtest.TradeGroups[key].BenchmarkLast = benchmarkLast
			backtest.TradeGroups[key].ReturnFromStart = helpers.Round((((backtest.EndingBalance - backtest.StartingBalance) / backtest.StartingBalance) * 100), 2)
			backtest.TradeGroups[key].ReturnPercent = helpers.Round((((backtest.TradeGroups[key].Margin + (backtest.TradeGroups[key].OpenPrice - backtest.TradeGroups[key].ClosePrice) - backtest.TradeGroups[key].Margin) / backtest.TradeGroups[key].Margin) * 100), 2)
			backtest.TradeGroups[key].BenchmarkBalance = helpers.Round((investedBenchmark*backtest.TradeGroups[key].BenchmarkLast)+investedBenchmarkLeftOver, 2)
			backtest.TradeGroups[key].BenchmarkReturn = helpers.Round((((backtest.TradeGroups[key].BenchmarkLast - backtest.BenchmarkStart) / backtest.BenchmarkStart) * 100), 2)
		}
	}

}

//
// expirePositions - Loop through to see if we have any positions to expire.
//
func expirePositions(today time.Time, underlyingLast float64, backtest *models.Backtest, benchmarkQuotes []types.HistoryQuote) {
	// Get the benchmark last
	benchmarkLast := getBenchmarkByDate(today, benchmarkQuotes)

	// Loop for expired postions
	for key, row := range backtest.TradeGroups {

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

			// Benchmark stuff
			investedBenchmark := math.Floor(backtest.StartingBalance / backtest.BenchmarkStart)
			investedBenchmarkLeftOver := backtest.StartingBalance - (investedBenchmark * backtest.BenchmarkStart)

			// Cleaner to make vars
			returnPercent := helpers.Round((((backtest.TradeGroups[key].Margin + (backtest.TradeGroups[key].OpenPrice - backtest.TradeGroups[key].ClosePrice) - backtest.TradeGroups[key].Margin) / backtest.TradeGroups[key].Margin) * 100), 2)

			// Shared for all strats
			backtest.TradeGroups[key].Status = "Closed"
			backtest.TradeGroups[key].CloseDate = models.Date{today}
			backtest.TradeGroups[key].BenchmarkLast = benchmarkLast
			backtest.TradeGroups[key].ReturnFromStart = helpers.Round((((backtest.EndingBalance - backtest.StartingBalance) / backtest.StartingBalance) * 100), 2)
			backtest.TradeGroups[key].ReturnPercent = returnPercent
			backtest.TradeGroups[key].BenchmarkBalance = helpers.Round((investedBenchmark*backtest.TradeGroups[key].BenchmarkLast)+investedBenchmarkLeftOver, 2)
			backtest.TradeGroups[key].BenchmarkReturn = helpers.Round((((backtest.TradeGroups[key].BenchmarkLast - backtest.BenchmarkStart) / backtest.BenchmarkStart) * 100), 2)

			// Figure out how to close based on strategy
			switch row.Strategy {

			// Put credit spread
			case "put-credit-spread":
				diff := underlyingLast - row.Positions[1].Symbol.OptionStrike

				// Expired worthless or in the money
				if diff > 0 {
					backtest.TradeGroups[key].ClosePrice = 0.00
					backtest.TradeGroups[key].Note = "Expired worthless."
				} else {
					spread := row.Positions[1].Symbol.OptionStrike - row.Positions[0].Symbol.OptionStrike

					if (diff * -1) > spread {
						backtest.TradeGroups[key].ClosePrice = (spread * float64(row.Lots) * 100)
					} else {
						backtest.TradeGroups[key].ClosePrice = ((diff * -1) * float64(row.Lots) * 100)
					}

					backtest.TradeGroups[key].Note = "Expired in the money."
				}

			// Unknown strategy
			default:
				backtest.TradeGroups[key].ClosePrice = 0.00
				backtest.TradeGroups[key].Note = "Expired unknown."

			}

		}
	}

}

// -------------- Helper Functions ---------------- //

//
// getClosedPrice - Figure out how much it would be to close this trade now
//
func getClosedPrice(tradegroup models.BacktestTradeGroup, options []types.OptionsChainItem) float64 {
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

//
// getBenchmarkByDate will return the close price of the benchmark
//
func getBenchmarkByDate(date time.Time, benchmarkQuotes []types.HistoryQuote) float64 {
	for _, row := range benchmarkQuotes {
		if row.Date.Format("2006-01-02") == date.Format("2006-01-02") {
			return row.Close
		}
	}
	return 0.00
}
