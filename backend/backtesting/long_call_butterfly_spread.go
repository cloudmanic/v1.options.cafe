//
// Date: 2022-05-15
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Copyright: 2022 Cloudmanic Labs, LLC. All rights reserved.
//

package backtesting

import (
	"fmt"
	"time"

	"app.options.cafe/brokers/eod"
	"app.options.cafe/brokers/types"
	"app.options.cafe/models"
	"app.options.cafe/screener"
	"github.com/davecgh/go-spew/spew"

	screenerCache "app.options.cafe/screener/cache"
)

//
// LongCallButterflySpreadPlaceTrades managed trades. Call this after all possible trades are found.
//
func (t *Base) LongCallButterflySpreadPlaceTrades(today time.Time, backtest *models.Backtest, results []screener.Result, options []types.OptionsChainItem) {

	spew.Dump(today)

	// // Make sure we have at least one result
	// if len(results) <= 0 {
	// 	return
	// }

	// // See if we have any positions to close
	// t.CloseMultiLegCredit(today, results[0].UnderlyingLast, backtest, options)

	// // Figure which result to open
	// result, err := t.PutCreditSpreadSelectTrade(today, backtest, results)

	// // Open trade
	// if err == nil {
	// 	t.OpenMultiLegCredit(today, "put-credit-spread", backtest, result)
	// }

	return
}

//
// LongCallButterflySpreadSelectTrade will figure out which trade we are placing today.
//
func (t *Base) LongCallButterflySpreadSelectTrade(today time.Time, backtest *models.Backtest, results []screener.Result) (screener.Result, error) {
	// Result we return.
	winner := screener.Result{}

	fmt.Println("LongCallButterflySpreadSelectTrade")

	// // Temp holding after filtering out.
	// tempResults := []screener.Result{}

	// // Loop through and filter out results we do not need.
	// for _, row := range results {
	// 	// First see if we already have this position
	// 	if !t.checkForCurrentPosition(backtest, row) {
	// 		continue
	// 	}

	// 	tempResults = append(tempResults, row)
	// }

	// // Make sure we have some results.
	// if len(tempResults) == 0 {
	// 	return winner, errors.New("no results found")
	// }

	// // Loop through temp list to find the result we want based on TradeSelect
	// switch backtest.TradeSelect {

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

	// // Find the highest ask
	// case "highest-ask":
	// 	for _, row := range tempResults {
	// 		// first lap Midpoint will equal zero
	// 		if row.Ask > winner.Ask {
	// 			winner = row
	// 		}
	// 	}
	// 	break

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
	// }

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

/* End File */
