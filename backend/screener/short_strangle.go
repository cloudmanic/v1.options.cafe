//
// Date: 2018-11-08
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-11-08
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package screener

import (
	"os"
	"strings"
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/library/helpers"
	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/cloudmanic/app.options.cafe/backend/models"
)

//
// Run ShortS Strangle screen
//
func (t *Base) RunShortStrangle(screen models.Screener) ([]Result, error) {

	result := []Result{}

	// Set today's date
	today := time.Now()

	// Change today's date for unit testing.
	if strings.HasSuffix(os.Args[0], ".test") {
		today = helpers.ParseDateNoError("2018-10-18").UTC()
	}

	// Make call to get current quote.
	quote, err := t.GetQuote(screen.Symbol)

	if err != nil {
		return result, err
	}

	// Get all possible expire dates.
	expires, err := t.Broker.GetOptionsExpirationsBySymbol(screen.Symbol)

	if err != nil {
		services.Warning(err)
		return result, err
	}

	// Add default values
	t.ShortStrangleFillDefault(&screen)

	// Loop through the expire dates
	for _, row := range expires {

		// Expire Date.
		expireDate, _ := time.Parse("2006-01-02", row)

		// Filter for expire dates
		if !t.FilterDaysToExpireDaysToExpire(today, screen, expireDate) {
			continue
		}

		// Get options Chain
		chain, err := t.Broker.GetOptionsChainByExpiration(screen.Symbol, row)

		if err != nil {
			continue
		}

		// Loop through the Put chain
		for _, row2 := range chain.Puts {

			// See if the strike is too
			if !t.FilterStrikeByPercentDown("put-leg-percent-away", screen, row2.Strike, quote.Last) {
				continue
			}

			// Loop through the call side to find possible trades.
			for _, row3 := range chain.Calls {

				// See if the strike is too
				if !t.FilterStrikeByPercentUp("call-leg-percent-away", screen, row3.Strike, quote.Last) {
					continue
				}

				// Get the credit for this trade
				credit := row2.Bid + row3.Bid

				if !t.FilterOpenCredit(screen, credit) {
					continue
				}

				// Figure out the amounts.
				closeCost := (row2.Ask + row3.Ask)
				midPoint := (credit + closeCost) / 2

				// Percent away - We show the lowest percent away
				putPercentAway := ((1 - (row2.Strike / quote.Last)) * 100)
				callPercentAway := ((row3.Strike - quote.Last) / quote.Last) * 100

				// Because of the rounding we do to find the closest .5 strike price some spreads will slip in we do one last filter.
				if !t.FilterPercentAwayResults("put-leg-percent-away", screen, putPercentAway) {
					continue
				}

				if !t.FilterPercentAwayResults("call-leg-percent-away", screen, callPercentAway) {
					continue
				}

				// Add in Symbol Object - Put Short leg
				symbPutShortLeg, err := t.DB.CreateNewSymbol(row2.Symbol, row2.Description, "Option")

				if err != nil {
					continue
				}

				// Add in Symbol Object - Put Short leg
				symbCallShortLeg, err := t.DB.CreateNewSymbol(row3.Symbol, row3.Description, "Option")

				if err != nil {
					continue
				}

				// We have a winner
				result = append(result, Result{
					Expired:         models.Date{expireDate},
					Credit:          helpers.Round(credit, 2),
					MidPoint:        helpers.Round(midPoint, 2),
					PutPrecentAway:  helpers.Round(putPercentAway, 2),
					CallPrecentAway: helpers.Round(callPercentAway, 2),
					Legs:            []models.Symbol{symbPutShortLeg, symbCallShortLeg},
				})

				// We only want the first 100
				if len(result) >= 100 {

					// Sort the results
					t.SortResults(result)

					return result, nil
				}

			}

		}

	}

	// Sort the results
	t.SortResults(result)

	// Return happy
	return result, nil

}

// ------------------------ Helper Functions -------------------------- //

//
// Setup default values. We need to make sure we have at least these params to run a screen.
//
func (t *Base) ShortStrangleFillDefault(screen *models.Screener) {

	// Map found
	found := map[string]bool{}

	// Fields that are required
	required := map[string]models.ScreenerItem{
		"open-credit":           {Key: "open-credit", Operator: ">", ValueNumber: 1.50},
		"put-leg-percent-away":  {Key: "put-leg-percent-away", Operator: ">", ValueNumber: 4.00},
		"call-leg-percent-away": {Key: "call-leg-percent-away", Operator: ">", ValueNumber: 4.00},
		"days-to-expire":        {Key: "days-to-expire", Operator: "<", ValueNumber: 30},
	}

	// Loop through and identify items we already have
	for _, row := range screen.Items {

		if _, ok := required[row.Key]; ok {
			found[row.Key] = true
		}

	}

	// Add default values
	for key, row := range required {

		if _, ok := found[key]; !ok {
			screen.Items = append(screen.Items, row)
		}

	}

}

/* End File */
