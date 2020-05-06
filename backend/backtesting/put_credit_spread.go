//
// Date: 2019-02-22
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package backtesting

import (
	"errors"
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/brokers/eod"
	"github.com/cloudmanic/app.options.cafe/backend/brokers/types"
	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/cloudmanic/app.options.cafe/backend/screener"
	"github.com/optionscafe/options-cafe-cli/helpers"
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
func (t *Base) PutCreditSpreadResults(today time.Time, backtest *models.Backtest, underlyingLast float64, options []types.OptionsChainItem) ([]screener.Result, error) {
	// Results that we return.
	results := []screener.Result{}

	// Set up a new screener so we can use it's Functions
	screenObj := screener.NewScreen(t.DB, &eod.Api{})

	// Set params
	spreadWidth := screenObj.GetPutCreditSpreadParms(backtest.Screen, underlyingLast)

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

		// Loop through the options we need.
		for _, row2 := range putOptions {

			// No need to pay attention to open interest of zero
			if row2.OpenInterest == 0 {
				continue
			}

			// Skip strikes that are higher than our min strike. Based on percent away.
			if !screenObj.FilterStrikeByPercentDown("short-strike-percent-away", backtest.Screen, row2.Strike, underlyingLast) {
				continue
			}

			// Find the strike that is x points away.
			buyLeg, err := screenObj.FindByStrike(putOptions, (row2.Strike - spreadWidth))

			if err != nil {
				continue
			}

			// See if there is enough credit
			credit := row2.Bid - buyLeg.Ask

			if !screenObj.FilterOpenCredit(backtest.Screen, credit) {
				continue
			}

			// Figure out the credit spread amount.
			buyCost := row2.Ask - buyLeg.Bid
			midPoint := (credit + buyCost) / 2

			// Add in Symbol Object - Sell leg
			symbSellLeg, err := t.GetSymbol(row2.Symbol, row2.Description, "Option")

			if err != nil {
				return []screener.Result{}, err
			}

			// Add in Symbol Object - Buy leg
			symbBuyLeg, err := t.GetSymbol(buyLeg.Symbol, buyLeg.Description, "Option")

			if err != nil {
				return []screener.Result{}, err
			}

			// We have a winner
			results = append(results, screener.Result{
				Day:            types.Date{today},
				Credit:         helpers.Round(credit, 2),
				Bid:            helpers.Round(buyCost, 2),
				Ask:            helpers.Round(credit, 2),
				MidPoint:       helpers.Round(midPoint, 2),
				UnderlyingLast: underlyingLast,
				PutPrecentAway: helpers.Round(((1 - row2.Strike/underlyingLast) * 100), 2),
				Legs:           []models.Symbol{symbBuyLeg, symbSellLeg},
			})
		}

	}

	// Return happy with results.
	return results, nil
}

/* End File */
