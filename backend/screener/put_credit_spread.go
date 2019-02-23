//
// Date: 2018-07-17
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package screener

import (
	"flag"
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/library/helpers"
	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/cloudmanic/app.options.cafe/backend/models"
)

//
// Run a put credit spread screen
//
func (t *Base) RunPutCreditSpread(screen models.Screener) ([]Result, error) {

	result := []Result{}

	// Set today's date
	today := time.Now()

	// Change today's date for unit testing.
	if flag.Lookup("test.v") != nil {
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

	// Set params
	spreadWidth := t.GetPutCreditSpreadParms(screen, quote.Last)

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

		for _, row2 := range chain.Puts {

			// No need to pay attention to open interest of zero
			if row2.OpenInterest == 0 {
				continue
			}

			// Skip strikes that are higher than our min strike. Based on percent away.
			if !t.FilterStrikeByPercentDown("short-strike-percent-away", screen, row2.Strike, quote.Last) {
				continue
			}

			// We only want the first 100
			if len(result) >= 100 {
				return result, nil
			}

			// Find the strike that is x points away.
			buyLeg, err := t.FindByStrike(chain.Puts, (row2.Strike - spreadWidth))

			if err != nil {
				continue
			}

			// See if there is enough credit
			credit := row2.Bid - buyLeg.Ask

			if !t.FilterOpenCredit(screen, credit) {
				continue
			}

			// Figure out the credit spread amount.
			buyCost := row2.Ask - buyLeg.Bid
			midPoint := (credit + buyCost) / 2

			// Add in Symbol Object - Buy leg
			symbBuyLeg, err := t.DB.CreateNewSymbol(buyLeg.Symbol, buyLeg.Description, "Option")

			if err != nil {
				continue
			}

			// Add in Symbol Object - Sell leg
			symbSellLeg, err := t.DB.CreateNewSymbol(row2.Symbol, row2.Description, "Option")

			if err != nil {
				continue
			}

			// We have a winner
			result = append(result, Result{
				Credit:         helpers.Round(credit, 2),
				MidPoint:       helpers.Round(midPoint, 2),
				PutPrecentAway: helpers.Round(((1 - row2.Strike/quote.Last) * 100), 2),
				Legs:           []models.Symbol{symbBuyLeg, symbSellLeg},
			})

		}

	}

	// Return happy
	return result, nil
}

//
// Set Parms we need.
//
func (t *Base) GetPutCreditSpreadParms(screen models.Screener, lastQuote float64) float64 {

	var spreadWidth float64 = 5.00

	// See if we have a spread width
	sw, err := t.FindFilterItemValue("spread-width", screen)

	if err == nil {
		spreadWidth = sw.ValueNumber
	}

	// Return values
	return spreadWidth
}

/* End File */
