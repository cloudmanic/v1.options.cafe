//
// Date: 2019-02-22
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package backtesting

import (
	"math"
	"time"

	"app.options.cafe/brokers/types"
	"app.options.cafe/models"
	"github.com/optionscafe/options-cafe-cli/helpers"
)

//
// CloseMultiLegCredit - Close positions
//
func (t *Base) CloseMultiLegCredit(today time.Time, underlyingLast float64, backtest *models.Backtest, options []types.OptionsChainItem) {
	// Close if we touch the short leg
	t.closeOnShortTouch(today, underlyingLast, backtest, options)

	// Close if we hit a particular debit
	t.closeOnDebit(today, underlyingLast, backtest, options)

	// Expire positions
	t.expirePositions(today, underlyingLast, backtest)
}

//
// closeOnDebit - Close a trade if it hits our debit trigger
//
func (t *Base) closeOnDebit(today time.Time, underlyingLast float64, backtest *models.Backtest, options []types.OptionsChainItem) {
	// TODO(spicer): make this work from configs
	debitAmount := 0.03

	// Get the benchmark last
	benchmarkLast := t.getBenchmarkByDate(today)

	// Loop for expired postions
	for key, row := range backtest.TradeGroups {

		if row.Status == "Closed" {
			continue
		}

		// Get closing price
		closePrice := t.getClosedPrice(row, options)

		// Benchmark stuff
		investedBenchmark := math.Floor(backtest.StartingBalance / backtest.BenchmarkStart)
		investedBenchmarkLeftOver := backtest.StartingBalance - (investedBenchmark * backtest.BenchmarkStart)

		// Close trade at the debitAmount
		if closePrice <= debitAmount {
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
func (t *Base) closeOnShortTouch(today time.Time, underlyingLast float64, backtest *models.Backtest, options []types.OptionsChainItem) {
	// Get the benchmark last
	benchmarkLast := t.getBenchmarkByDate(today)

	// Loop for expired postions
	for key, row := range backtest.TradeGroups {

		if row.Status == "Closed" {
			continue
		}

		// TODO(Spicer): Currently this only works for PCS. We assume the second leg is the short leg.
		if underlyingLast <= row.Positions[1].Symbol.OptionStrike {

			// Set closing Closing Price
			closingPrice := (t.getClosedPrice(row, options) * 100.00 * float64(row.Lots)) - 1

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
func (t *Base) expirePositions(today time.Time, underlyingLast float64, backtest *models.Backtest) {
	// Get the benchmark last
	benchmarkLast := t.getBenchmarkByDate(today)

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

			// Shared for all strats
			backtest.TradeGroups[key].Status = "Closed"
			backtest.TradeGroups[key].CloseDate = models.Date{today}
			backtest.TradeGroups[key].BenchmarkLast = benchmarkLast
			backtest.TradeGroups[key].ReturnFromStart = helpers.Round((((backtest.EndingBalance - backtest.StartingBalance) / backtest.StartingBalance) * 100), 2)
			backtest.TradeGroups[key].ReturnPercent = helpers.Round((((backtest.TradeGroups[key].Margin + (backtest.TradeGroups[key].OpenPrice - backtest.TradeGroups[key].ClosePrice) - backtest.TradeGroups[key].Margin) / backtest.TradeGroups[key].Margin) * 100), 2)
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
func (t *Base) getClosedPrice(position models.BacktestTradeGroup, options []types.OptionsChainItem) float64 {

	// TODO(spicer): Make this work for everything. Currently just works for PCS
	var leg1Chain types.OptionsChainItem
	var leg2Chain types.OptionsChainItem

	// Loop through until we find first symbol
	for _, row := range options {
		if row.Symbol != position.Positions[0].Symbol.ShortName {
			continue
		}

		// We found it
		leg1Chain = row
		break
	}

	// Loop through until we find second symbol
	for _, row := range options {
		if row.Symbol != position.Positions[1].Symbol.ShortName {
			continue
		}

		// We found it
		leg2Chain = row
		break
	}

	// Get price to close.
	return leg2Chain.Ask - leg1Chain.Bid
}

/* End File */
