//
// Date: 2019-02-22
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package putcreditspread

import (
	"errors"
	"time"

	"app.options.cafe/brokers/types"
	"app.options.cafe/models"
	"app.options.cafe/screener"
)

//
// Trades managed trades. Call this after all possible trades are found.
//
func Trades(db models.Datastore, today time.Time, backtest *models.Backtest, results []screener.Result, options []types.OptionsChainItem, underlyingLast float64, benchmarkQuotes []types.HistoryQuote) {
	// Make sure we have at least one result
	if len(results) <= 0 {
		return
	}

	// See if we have any positions to close
	CloseTrades(today, results[0].UnderlyingLast, backtest, options, benchmarkQuotes)

	// Figure which result to open
	result, err := SelectTrade(today, backtest, results)

	// Open trade
	if err == nil {
		OpenMultiLegCredit(db, today, "put-credit-spread", backtest, result)
	}

	return
}

//
// SelectTrade will figure out which trade we are placing today.
//
func SelectTrade(today time.Time, backtest *models.Backtest, results []screener.Result) (screener.Result, error) {
	// Result we return.
	winner := screener.Result{}

	// Temp holding after filtering out.
	tempResults := []screener.Result{}

	// Loop through and filter out results we do not need.
	for _, row := range results {
		// First see if we already have this position
		if !checkForCurrentPosition(backtest, row) {
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

	// Find the highest percent away
	case "highest-percent-away":
		for _, row := range tempResults {
			if row.PutPrecentAway > winner.PutPrecentAway {
				winner = row
			}
		}
		break

	// Find the shortest percent away
	case "shortest-percent-away":
		for key, row := range tempResults {
			// winner is unknown on first lap
			if key == 0 {
				winner = row
			}

			if row.PutPrecentAway < winner.PutPrecentAway {
				winner = row
			}
		}
		break
	}

	return winner, nil
}
