//
// Date: 2022-05-15
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Copyright: 2022 Cloudmanic Labs, LLC. All rights reserved.
//

package backtesting

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"app.options.cafe/brokers/eod"
	"app.options.cafe/brokers/types"
	"app.options.cafe/library/helpers"
	"app.options.cafe/models"
	"app.options.cafe/screener"

	screenerCache "app.options.cafe/screener/cache"
)

//
// LongCallButterflySpreadPlaceTrades managed trades. Call this after all possible trades are found.
//
func (t *Base) LongCallButterflySpreadPlaceTrades(today time.Time, backtest *models.Backtest, results []screener.Result, options []types.OptionsChainItem, underlyingLast float64, benchmarkQuotes []types.HistoryQuote) {
	// See if we have any positions to close
	t.CloseLongCallButterflySpread(today, underlyingLast, backtest, options)

	// Make sure we have at least one result
	if len(results) <= 0 {
		return
	}

	// Figure which result to open
	result, err := t.LongCallButterflySpreadSelectTrade(today, backtest, results)

	// Open trade
	if err == nil {
		t.LongCallButterflySpreadOpenTrade(today, "long-call-butterfly-spread", backtest, result)
	}

	return
}

//
// CloseLongCallButterflySpread - Close positions
//
func (t *Base) CloseLongCallButterflySpread(today time.Time, underlyingLast float64, backtest *models.Backtest, options []types.OptionsChainItem) {

	// // Close if we touch the short leg
	// t.closeOnShortTouch(today, underlyingLast, backtest, options)

	// Close if we hit a particular debit
	t.CloseOnDebit(today, backtest, options)

	// Expire positions
	t.LongCallButterflySpreadExpirePositions(today, underlyingLast, backtest, options)
}

//
// LongCallButterflySpreadExpirePositions - Expire trades
//
func (t *Base) LongCallButterflySpreadExpirePositions(today time.Time, underlyingLast float64, backtest *models.Backtest, options []types.OptionsChainItem) {
	// Get the benchmark last
	benchmarkLast := t.getBenchmarkByDate(today)

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

			// Benchmark stuff
			investedBenchmark := math.Floor(backtest.StartingBalance / backtest.BenchmarkStart)
			investedBenchmarkLeftOver := backtest.StartingBalance - (investedBenchmark * backtest.BenchmarkStart)

			// Shared for all strats
			backtest.TradeGroups[key].Status = "Closed"
			backtest.TradeGroups[key].CloseDate = models.Date{today}
			backtest.TradeGroups[key].BenchmarkLast = benchmarkLast
			backtest.TradeGroups[key].ReturnFromStart = helpers.Round((((backtest.EndingBalance - backtest.StartingBalance) / backtest.StartingBalance) * 100), 2)
			backtest.TradeGroups[key].BenchmarkBalance = helpers.Round((investedBenchmark*backtest.TradeGroups[key].BenchmarkLast)+investedBenchmarkLeftOver, 2)
			backtest.TradeGroups[key].BenchmarkReturn = helpers.Round((((backtest.TradeGroups[key].BenchmarkLast - backtest.BenchmarkStart) / backtest.BenchmarkStart) * 100), 2)

			// Get closing price
			closePrice := (t.getClosedPrice(row, options) * 100)

			// See if we have a profit.
			profit := closePrice - row.OpenPrice

			if profit > 0 {
				backtest.EndingBalance = backtest.EndingBalance + closePrice
				backtest.TradeGroups[key].ClosePrice = closePrice
				backtest.TradeGroups[key].Note = "Took profit on expire day."
				backtest.TradeGroups[key].Balance = backtest.EndingBalance
				backtest.TradeGroups[key].BenchmarkLast = benchmarkLast
				backtest.TradeGroups[key].ReturnFromStart = helpers.Round((((backtest.EndingBalance - backtest.StartingBalance) / backtest.StartingBalance) * 100), 2)
				backtest.TradeGroups[key].ReturnPercent = helpers.Round(((closePrice-row.OpenPrice)/row.OpenPrice)*100, 2)
			} else {
				backtest.TradeGroups[key].ClosePrice = 0.00
				backtest.TradeGroups[key].Note = "Expired worthless."
				backtest.TradeGroups[key].ReturnPercent = -100.00
			}
		}
	}
}

//
// closeOnDebit - Close a trade if it hits our debit trigger
//
func (t *Base) CloseOnDebit(today time.Time, backtest *models.Backtest, options []types.OptionsChainItem) {
	// Take profit amount.
	takeProfitLessThanPercent := 0.00
	takeProfitGreaterThanPercent := 0.00

	// Set up a new screener so we can use it's Functions. Figure out when to take profit.
	screenObj := screener.NewScreen(t.DB, &eod.Api{})
	items := screenObj.FindFilterItemsByKey("take-profit-percent", backtest.Screen)

	if len(items) < 0 {
		return
	}

	for _, row := range items {
		if row.Operator == ">" {
			takeProfitGreaterThanPercent = (row.ValueNumber / 100)
		} else if row.Operator == "<" {
			takeProfitLessThanPercent = (row.ValueNumber / 100)
		}
	}

	// Get the benchmark last
	benchmarkLast := t.getBenchmarkByDate(today)

	// Loop for expired postions
	for key, row := range backtest.TradeGroups {
		// Already closed?
		if row.Status == "Closed" {
			continue
		}

		// Get closing price
		closePrice := (t.getClosedPrice(row, options) * 100)

		// See if we have a profit.
		profit := closePrice - row.OpenPrice

		// Trigger profit
		gtePrice := 0.00

		if takeProfitGreaterThanPercent > 0 {
			gtePrice = row.OpenPrice * takeProfitGreaterThanPercent
		}

		ltePrice := row.OpenPrice * 1000000 // random big number
		if takeProfitLessThanPercent > 0 {
			ltePrice = row.OpenPrice * takeProfitLessThanPercent
		}

		// Benchmark stuff
		investedBenchmark := math.Floor(backtest.StartingBalance / backtest.BenchmarkStart)
		investedBenchmarkLeftOver := backtest.StartingBalance - (investedBenchmark * backtest.BenchmarkStart)

		// If we have a good profit close the trade.
		if (profit > gtePrice) && (profit < ltePrice) {
			backtest.EndingBalance = backtest.EndingBalance + closePrice
			backtest.TradeGroups[key].Status = "Closed"
			backtest.TradeGroups[key].ClosePrice = closePrice
			backtest.TradeGroups[key].CloseDate = models.Date{today}
			backtest.TradeGroups[key].Note = "Triggered at profit amount."
			backtest.TradeGroups[key].Balance = backtest.EndingBalance
			backtest.TradeGroups[key].BenchmarkLast = benchmarkLast
			backtest.TradeGroups[key].ReturnFromStart = helpers.Round((((backtest.EndingBalance - backtest.StartingBalance) / backtest.StartingBalance) * 100), 2)
			backtest.TradeGroups[key].ReturnPercent = helpers.Round(((closePrice-row.OpenPrice)/row.OpenPrice)*100, 2)
			backtest.TradeGroups[key].BenchmarkBalance = helpers.Round((investedBenchmark*backtest.TradeGroups[key].BenchmarkLast)+investedBenchmarkLeftOver, 2)
			backtest.TradeGroups[key].BenchmarkReturn = helpers.Round((((backtest.TradeGroups[key].BenchmarkLast - backtest.BenchmarkStart) / backtest.BenchmarkStart) * 100), 2)
		}
	}

}

//
// LongCallButterflySpreadSelectTrade will figure out which trade we are placing today.
//
func (t *Base) LongCallButterflySpreadSelectTrade(today time.Time, backtest *models.Backtest, results []screener.Result) (screener.Result, error) {
	// Result we return.
	winner := screener.Result{}

	// Temp holding after filtering out.
	tempResults := []screener.Result{}

	// Loop through and filter out results we do not need.
	for _, row := range results {
		// First see if we already have this position
		if !t.checkForCurrentPosition(backtest, row) {
			continue
		}

		tempResults = append(tempResults, row)
	}

	// Make sure we have some results.
	if len(tempResults) == 0 {
		return winner, errors.New("no results found")
	}

	// Loop through temp list to find the result we want based on TradeSelect
	switch backtest.TradeSelect {

	// // Least days to expire.
	// case "least-days-to-expire":
	// 	daysTracker := 10000000 // Random big number

	// 	for _, row := range tempResults {
	// 		expire, _ := time.Parse("2006-01-02", row.Legs[0].OptionExpire.Format("2006-01-02"))
	// 		daysToExpire := int(today.Sub(expire).Hours() / 24 * -1)

	// 		if daysTracker > daysToExpire {
	// 			daysTracker = daysToExpire
	// 			winner = row
	// 		}
	// 	}
	// 	break

	// // Find the highest midpoint
	// case "highest-midpoint":
	// 	for _, row := range tempResults {
	// 		// first lap Midpoint will equal zero
	// 		if row.MidPoint > winner.MidPoint {
	// 			winner = row
	// 		}
	// 	}
	// 	break

	// Find the lowest ask
	case "lowest-ask":
		for _, row := range tempResults {
			// first lap Ask will equal zero
			if row.Ask < winner.Ask {
				winner = row
			}
		}
		break

		// // Find the highest percent away
		// case "highest-percent-away":
		// 	for _, row := range tempResults {
		// 		if row.PutPrecentAway > winner.PutPrecentAway {
		// 			winner = row
		// 		}
		// 	}
		// 	break

		// // Find the shortest percent away
		// case "shortest-percent-away":
		// 	for key, row := range tempResults {
		// 		// winner is unknown on first lap
		// 		if key == 0 {
		// 			winner = row
		// 		}

		// 		if row.PutPrecentAway < winner.PutPrecentAway {
		// 			winner = row
		// 		}
		// 	}
		// 	break
	}

	return winner, nil
}

//
// LongCallButterflySpreadResults - Find possible trades for this strategy.
//
func (t *Base) LongCallButterflySpreadResults(db models.Datastore, today time.Time, backtest *models.Backtest, underlyingLast float64, options []types.OptionsChainItem, cache screenerCache.Cache) ([]screener.Result, error) {
	// Results that we return.
	results := []screener.Result{}

	// Set up a new screener so we can use it's Functions
	screenObj := screener.NewScreen(t.DB, &eod.Api{})

	// Take complete list of options and return a list of expiration dates.
	expireDates := t.GetExpirationDatesFromOptions(options)

	// Loop through the expire dates
	for _, row := range expireDates {
		// Expire Date.
		expireDate, _ := time.Parse("2006-01-02", row.Format("2006-01-02"))

		// Filter for expire dates
		if !screenObj.FilterDaysToExpire(today, backtest.Screen, expireDate) {
			continue
		}

		// Get the options and just pull out the CALL options for this expire date.
		callOptions := t.GetOptionsByExpirationType(row, "Call", options)

		// Core screen logic from our screener.
		for _, row2 := range screenObj.LongCallButterflySpreadCoreScreen(today, backtest.Screen, callOptions, underlyingLast, cache) {
			results = append(results, row2)
		}
	}

	// Return happy with results.
	return results, nil
}

//
// LongCallButterflySpreadOpenTrade - Open a new spread by adding a position
//
func (t *Base) LongCallButterflySpreadOpenTrade(today time.Time, strategy string, backtest *models.Backtest, result screener.Result) {
	// Default values
	lots := 1

	// // First see if we already have this position
	// if !t.checkForCurrentPosition(backtest, result) {
	// 	return
	// }

	// Figure out position size (We just trade one lot when we do this, so set the account balance to like 200% of one lot)
	if backtest.PositionSize == "one-at-time" {
		// Get the count of open positions
		posCount := t.openPositionsCount(backtest)

		// Only open one position at a time.
		if posCount > 0 {
			return
		}
	} else if strings.Contains(backtest.PositionSize, "percent") { // percent of trade
		totalToTrade := t.percentOfAccount(backtest, backtest.PositionSize)
		lots = int(math.Floor(totalToTrade / (result.Bid * 100.00)))
	}

	// Figure out open price.
	openPrice := result.Bid * 100 * float64(lots)

	// we can configure to use midpoint.
	if backtest.Midpoint {
		openPrice = result.MidPoint * 100 * float64(lots)
	}

	// // See if we are allowed to only have one strike
	// st, err := screenObj.FindFilterItemValue("allow-more-than-one-strike", backtest.Screen)

	// // Make sure we do not already have a trade on at this strike.
	// if (err == nil) && (st.ValueString == "no") {
	// 	for _, row := range backtest.TradeGroups {
	// 		// Only open trades
	// 		if row.Status != "Open" {
	// 			continue
	// 		}

	// 		for _, row2 := range row.Positions {
	// 			if row2.Symbol.OptionStrike == result.Legs[0].OptionStrike {
	// 				return
	// 			}

	// 			if row2.Symbol.OptionStrike == result.Legs[1].OptionStrike {
	// 				return
	// 			}
	// 		}
	// 	}
	// }

	// // See if we are allowed to only have one expire
	// st2, err := screenObj.FindFilterItemValue("allow-more-than-one-expire", backtest.Screen)

	// // Make sure we do not already have a trade on at this expire.
	// if (err == nil) && (st2.ValueString == "no") {
	// 	for _, row := range backtest.TradeGroups {
	// 		// Only open trades
	// 		if row.Status != "Open" {
	// 			continue
	// 		}

	// 		for _, row2 := range row.Positions {
	// 			if row2.Symbol.OptionExpire == result.Legs[0].OptionExpire {
	// 				return
	// 			}

	// 			if row2.Symbol.OptionExpire == result.Legs[1].OptionExpire {
	// 				return
	// 			}
	// 		}
	// 	}
	// }

	// Build legs
	legs := []models.BacktestPosition{}

	for key, row := range result.Legs {
		qty := lots

		// Short leg
		if key == 1 {
			qty = qty * 2 * -1
		}

		legs = append(legs, models.BacktestPosition{
			UserId:   backtest.UserId,
			Status:   "Open",
			SymbolId: row.Id,
			Symbol:   row,
			OpenDate: today,
			Qty:      qty,
			OrgQty:   qty,
		})
	}

	// Spread text
	spreadText := fmt.Sprintf("%s %s $%.2f / $%.2f / $%.2f", legs[0].Symbol.OptionUnderlying, legs[0].Symbol.OptionExpire.Format("01/02/2006"), legs[0].Symbol.OptionStrike, legs[1].Symbol.OptionStrike, legs[2].Symbol.OptionStrike)

	// Add position
	backtest.TradeGroups = append(backtest.TradeGroups, models.BacktestTradeGroup{
		UserId:          backtest.UserId,
		Strategy:        strategy,
		Status:          "Open",
		SpreadText:      spreadText,
		OpenDate:        models.Date{today},
		OpenPrice:       openPrice,
		Margin:          0.00, // No margin this is a debit trade
		Positions:       legs,
		Lots:            lots,
		Credit:          (openPrice / float64(lots)) / 100,
		PutPrecentAway:  result.PutPrecentAway,
		CallPrecentAway: result.CallPrecentAway,
		Balance:         (backtest.EndingBalance - openPrice),
	})

	// Update ending balance
	backtest.EndingBalance = backtest.EndingBalance - openPrice

	//fmt.Println(today.Format("2006-01-02"), " : ", backtest.EndingBalance, " / ", spreadText, " - ", openPrice)
}

/* End File */
