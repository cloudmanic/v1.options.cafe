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
	"app.options.cafe/models"
	"app.options.cafe/screener"

	screenerCache "app.options.cafe/screener/cache"
)

//
// LongCallButterflySpreadPlaceTrades managed trades. Call this after all possible trades are found.
//
func (t *Base) LongCallButterflySpreadPlaceTrades(today time.Time, backtest *models.Backtest, results []screener.Result, options []types.OptionsChainItem) {
	// Make sure we have at least one result
	if len(results) <= 0 {
		return
	}

	// // See if we have any positions to close
	// t.CloseMultiLegCredit(today, results[0].UnderlyingLast, backtest, options)

	// Figure which result to open
	result, err := t.LongCallButterflySpreadSelectTrade(today, backtest, results)

	// Open trade
	if err == nil {
		t.LongCallButterflySpreadOpenTrade(today, "long-call-butterfly-spread", backtest, result)
	}

	return
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
func (t *Base) LongCallButterflySpreadResults(today time.Time, backtest *models.Backtest, underlyingLast float64, options []types.OptionsChainItem, cache screenerCache.Cache) ([]screener.Result, error) {
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

	// // Set up a new screener so we can use it's Functions
	// screenObj := screener.NewScreen(t.DB, &eod.Api{})

	// if result.Credit == 0 {
	// 	fmt.Println("#####")
	// 	spew.Dump(result)
	// 	fmt.Println("#####")
	// }

	// // Amount of margin left after trade is opened.
	// diff := result.Legs[1].OptionStrike - result.Legs[0].OptionStrike

	// Figure out position size
	if backtest.PositionSize == "one-at-time" {
		// Get the count of open positions
		posCount := t.openPositionsCount(backtest)

		// Only open one position at a time. TODO(spicer): make this a config.
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

	fmt.Println(today.Format("2006-01-02"), " : ", backtest.EndingBalance, " / ", spreadText, " - ", openPrice)
}

/* End File */
