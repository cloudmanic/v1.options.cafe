//
// Date: 2019-02-22
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package backtesting

import (
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/brokers/types"
	"github.com/cloudmanic/app.options.cafe/backend/models"
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
	for key, row := range backtest.Positions {

		if row.Status == "Closed" {
			continue
		}

		// Get closing price
		closePrice := t.getClosedPrice(row, options)

		// Close trade at the debitAmount
		if closePrice <= debitAmount {
			backtest.Positions[key].Status = "Closed"
			backtest.Positions[key].ClosePrice = closePrice * 100 * float64(row.Lots)
			backtest.Positions[key].CloseDate = models.Date{today}
			backtest.Positions[key].Note = "Triggered at debit amount."
			backtest.EndingBalance = (backtest.EndingBalance - backtest.Positions[key].ClosePrice)
			backtest.Positions[key].Balance = backtest.EndingBalance
			backtest.Positions[key].BenchmarkLast = benchmarkLast
			backtest.Positions[key].ReturnFromStart = (((backtest.EndingBalance - backtest.StartingBalance) / backtest.StartingBalance) * 100)
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
	for key, row := range backtest.Positions {

		if row.Status == "Closed" {
			continue
		}

		// TODO(Spicer): Currently this only works for PCS. We assume the second leg is the short leg.
		if underlyingLast <= row.Legs[1].OptionStrike {

			// Set closing Closing Price
			closingPrice := (t.getClosedPrice(row, options) * 100.00 * float64(row.Lots)) - 1

			// Make sure close price is not bigger than our max risk
			if closingPrice > row.Margin {
				continue
			}

			// Close trade
			backtest.Positions[key].Status = "Closed"
			backtest.Positions[key].ClosePrice = closingPrice
			backtest.Positions[key].CloseDate = models.Date{today}
			backtest.Positions[key].Note = "Trade touched the short leg."
			backtest.Positions[key].Balance -= closingPrice
			backtest.EndingBalance -= closingPrice
			backtest.Positions[key].BenchmarkLast = benchmarkLast
			backtest.Positions[key].ReturnFromStart = (((backtest.EndingBalance - backtest.StartingBalance) / backtest.StartingBalance) * 100)
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
	for key, row := range backtest.Positions {

		if row.Status == "Closed" {
			continue
		}

		expired := false

		// See if any of the legs are expired
		for _, row2 := range row.Legs {
			if today.Format("2006-01-02") == row2.OptionExpire.Format("2006-01-02") ||
				today.After(helpers.ParseDateNoError(row2.OptionExpire.Format("2006-01-02"))) {
				expired = true
				break
			}
		}

		// If expired close out trade
		if expired && row.Status == "Open" {

			// Shared for all strats
			backtest.Positions[key].Status = "Closed"
			backtest.Positions[key].CloseDate = models.Date{today}
			backtest.Positions[key].BenchmarkLast = benchmarkLast
			backtest.Positions[key].ReturnFromStart = (((backtest.EndingBalance - backtest.StartingBalance) / backtest.StartingBalance) * 100)

			// Figure out how to close based on strategy
			switch row.Strategy {

			// Put credit spread
			case "put-credit-spread":
				diff := underlyingLast - row.Legs[1].OptionStrike

				// Expired worthless or in the money
				if diff > 0 {
					backtest.Positions[key].ClosePrice = 0.00
					backtest.Positions[key].Note = "Expired worthless."
				} else {
					spread := row.Legs[1].OptionStrike - row.Legs[0].OptionStrike

					if (diff * -1) > spread {
						backtest.Positions[key].ClosePrice = (spread * float64(row.Lots) * 100)
					} else {
						backtest.Positions[key].ClosePrice = ((diff * -1) * float64(row.Lots) * 100)
					}

					backtest.Positions[key].Note = "Expired in the money."
				}

			// Unknown strategy
			default:
				backtest.Positions[key].ClosePrice = 0.00
				backtest.Positions[key].Note = "Expired unknown."

			}

		}
	}

}

// -------------- Helper Functions ---------------- //

//
// getClosedPrice - Figure out how much it would be to close this trade now
//
func (t *Base) getClosedPrice(position models.BacktestPosition, options []types.OptionsChainItem) float64 {

	// TODO(spicer): Make this work for everything. Currently just works for PCS
	var leg1Chain types.OptionsChainItem
	var leg2Chain types.OptionsChainItem

	// Loop through until we find first symbol
	for _, row := range options {
		if row.Symbol != position.Legs[0].ShortName {
			continue
		}

		// We found it
		leg1Chain = row
		break
	}

	// Loop through until we find second symbol
	for _, row := range options {
		if row.Symbol != position.Legs[1].ShortName {
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
