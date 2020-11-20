//
// Date: 2019-02-22
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package backtesting

import (
	"errors"
	"time"

	"app.options.cafe/brokers/eod"
	"app.options.cafe/brokers/types"
	"app.options.cafe/models"
	"app.options.cafe/screener"

	screenerCache "app.options.cafe/screener/cache"
)

//
// PutCreditSpreadPlaceTrades managed trades. Call this after all possible trades are found.
//
func (t *Base) PutCreditSpreadPlaceTrades(today time.Time, backtest *models.Backtest, results []screener.Result, options []types.OptionsChainItem) {
	// Make sure we have at least one result
	if len(results) <= 0 {
		return
	}

	// See if we have any positions to close
	t.CloseMultiLegCredit(today, results[0].UnderlyingLast, backtest, options)

	// Figure which result to open
	result, err := t.PutCreditSpreadSelectTrade(today, backtest, results)

	// Open trade
	if err == nil {
		t.OpenMultiLegCredit(today, "put-credit-spread", backtest, result)
	}

	return
}

//
// PutCreditSpreadSelectTrade will figure out which trade we are placing today.
//
func (t *Base) PutCreditSpreadSelectTrade(today time.Time, backtest *models.Backtest, results []screener.Result) (screener.Result, error) {
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

	// Least days to expire.
	case "least-days-to-expire":
		daysTracker := 10000000 // Random big number

		for _, row := range tempResults {
			expire, _ := time.Parse("2006-01-02", row.Legs[0].OptionExpire.Format("2006-01-02"))
			daysToExpire := int(today.Sub(expire).Hours() / 24 * -1)

			if daysTracker > daysToExpire {
				daysTracker = daysToExpire
				winner = row
			}
		}
		break

	// Find the highest midpoint
	case "highest-midpoint":
		for _, row := range tempResults {
			// first lap Midpoint will equal zero
			if row.MidPoint > winner.MidPoint {
				winner = row
			}
		}
		break

	// Find the highest ask
	case "highest-ask":
		for _, row := range tempResults {
			// first lap Midpoint will equal zero
			if row.Ask > winner.Ask {
				winner = row
			}
		}
		break

	}

	return winner, nil
}

//
// PutCreditSpreadResults - Find possible trades for this strategy.
//
func (t *Base) PutCreditSpreadResults(today time.Time, backtest *models.Backtest, underlyingLast float64, options []types.OptionsChainItem, cache screenerCache.Cache) ([]screener.Result, error) {
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

		// Get the options and just pull out the PUT options for this expire date.
		putOptions := t.GetOptionsByExpirationType(row, "Put", options)

		// Core screen logic from our screener.
		for _, row2 := range screenObj.PutCreditSpreadCoreScreen(today, backtest.Screen, putOptions, underlyingLast, cache) {
			results = append(results, row2)
		}
	}

	// Return happy with results.
	return results, nil
}

/* End File */
