//
// Date: 2018-07-17
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package screener

import (
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/library/helpers"
	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/cloudmanic/app.options.cafe/backend/models"
)

//
// Run a put credit spread screen
//
func RunPutCreditSpread(screen models.Screener, db models.Datastore) ([]Result, error) {

	result := []Result{}

	// Make call to get current quote.
	quote, err := GetQuote(screen.Symbol)

	if err != nil {
		return result, err
	}

	// Set params
	spreadWidth := getPutCreditSpreadParms(screen, quote.Last)

	// Get all possible expire dates.
	expires, err := broker.GetOptionsExpirationsBySymbol(screen.Symbol)

	if err != nil {
		services.Warning(err)
		return result, err
	}

	// Loop through the expire dates
	for _, row := range expires {

		// Expire Date.
		expireDate, _ := time.Parse("2006-01-02", row)

		// Filter for expire dates
		if !FilterDaysToExpireDaysToExpire(screen, expireDate) {
			continue
		}

		// Get options Chain
		chain, err := broker.GetOptionsChainByExpiration(screen.Symbol, row)

		if err != nil {
			continue
		}

		for _, row2 := range chain.Puts {

			// No need to pay attention to open interest of zero
			if row2.OpenInterest == 0 {
				continue
			}

			// Skip strikes that are higher than our min strike. Based on percent away.
			if !FilterStrikeByPercentDown("short-strike-percent-away", screen, row2.Strike, quote.Last) {
				continue
			}

			// Find the strike that is x points away.
			buyLeg, err := FindByStrike(chain.Puts, (row2.Strike - spreadWidth))

			if err != nil {
				continue
			}

			// See if there is enough credit
			credit := row2.Bid - buyLeg.Ask

			if !FilterOpenCredit(screen, credit) {
				continue
			}

			// Figure out the credit spread amount.
			buyCost := row2.Ask - buyLeg.Bid
			midPoint := (credit + buyCost) / 2

			// Add in Symbol Object - Buy leg
			symbBuyLeg, err := db.CreateNewSymbol(buyLeg.Symbol, buyLeg.Description, "Option")

			if err != nil {
				continue
			}

			// Add in Symbol Object - Sell leg
			symbSellLeg, err := db.CreateNewSymbol(row2.Symbol, row2.Description, "Option")

			if err != nil {
				continue
			}

			// We have a winner
			result = append(result, Result{
				Credit:      helpers.Round(credit, 2),
				MidPoint:    helpers.Round(midPoint, 2),
				PrecentAway: helpers.Round(((1 - row2.Strike/quote.Last) * 100), 2),
				Legs:        []models.Symbol{symbBuyLeg, symbSellLeg},
			})

		}

	}

	// Return happy
	return result, nil
}

// --------------------- Private Helper Functions ------------------------ //

//
// Set Parms we need.
//
func getPutCreditSpreadParms(screen models.Screener, lastQuote float64) float64 {

	var spreadWidth float64 = 5.00

	// See if we have a spread width
	sw, err := FindFilterItemValue("spread-width", screen)

	if err == nil {
		spreadWidth = sw.ValueNumber
	}

	// Return values
	return spreadWidth
}

/* End File */
