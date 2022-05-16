//
// Date: 2022-05-15
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Copyright: 2022 Cloudmanic Labs, LLC. All rights reserved.
//

package singleoption

import (
	"time"

	"app.options.cafe/backtesting/helpers"
	"app.options.cafe/brokers/eod"
	"app.options.cafe/brokers/types"
	"app.options.cafe/models"
	"app.options.cafe/screener"
	"app.options.cafe/screener/cache"
)

//
// Results will return results per day that we then decide if we trade or not.
//
func Results(db models.Datastore, today time.Time, backtest *models.Backtest, underlyingLast float64, options []types.OptionsChainItem, cache cache.Cache) ([]screener.Result, error) {
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

		// Core screen logic from our screener.
		for _, row2 := range screenObj.SingleOptionCoreScreen(today, expireDate, backtest.Screen, options, underlyingLast, cache) {
			results = append(results, row2)
		}
	}

	// Return happy with results.
	return results, nil
}
