//
// Date: 2019-02-22
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package putcreditspread

import (
	"time"

	"app.options.cafe/backtesting/helpers"
	"app.options.cafe/brokers/eod"
	"app.options.cafe/brokers/types"
	"app.options.cafe/models"
	"app.options.cafe/screener"

	screenerCache "app.options.cafe/screener/cache"
)

//
// Results - Find possible trades for this strategy.
//
func Results(db models.Datastore, today time.Time, backtest *models.Backtest, underlyingLast float64, options []types.OptionsChainItem, cache screenerCache.Cache) ([]screener.Result, error) {
	// Results that we return.
	results := []screener.Result{}

	// Set up a new screener so we can use it's Functions
	screenObj := screener.NewScreen(db, &eod.Api{})

	// Take complete list of options and return a list of expiration dates.
	expireDates := helpers.GetExpirationDatesFromOptions(options)

	// Loop through the expire dates
	for _, row := range expireDates {
		// Expire Date.
		expireDate, _ := time.Parse("2006-01-02", row.Format("2006-01-02"))

		// Filter for expire dates
		if !screenObj.FilterDaysToExpire(today, backtest.Screen, expireDate) {
			continue
		}

		// Get the options and just pull out the PUT options for this expire date.
		putOptions := helpers.GetOptionsByExpirationType(row, "Put", options)

		// Core screen logic from our screener.
		for _, row2 := range screenObj.PutCreditSpreadCoreScreen(today, backtest.Screen, putOptions, underlyingLast, cache) {
			results = append(results, row2)
		}
	}

	// Return happy with results.
	return results, nil
}
