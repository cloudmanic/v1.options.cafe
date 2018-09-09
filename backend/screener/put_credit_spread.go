//
// Date: 2018-07-17
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package screener

import (
	"math"
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
	today := time.Now()

	// Make call to get current quote.
	quote, err := GetQuote(screen.Symbol)

	if err != nil {
		return result, err
	}

	// Set params
	minDaysToExpire, maxDaysToExpire, minCredit, spreadWidth, minSellStrike := getPutCreditSpreadParms(screen, quote.Last)

	// Get all possible expire dates.
	expires, err := broker.GetOptionsExpirationsBySymbol(screen.Symbol)

	if err != nil {
		services.Warning(err)
		return result, err
	}

	// Loop through the expire dates
	for _, row := range expires {

		// Days to expire.
		then, _ := time.Parse("2006-01-02", row)
		daysToExpire := int(today.Sub(then).Hours()/24) * -1

		// Don't want to go too far out.
		if daysToExpire > maxDaysToExpire {
			continue
		}

		// Don't want to go too close out.
		if daysToExpire < minDaysToExpire {
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

			// Skip strikes that are higher than our min strike. If the user
			// did not set this value minSellStrike is set to Zero
			if row2.Strike > minSellStrike {
				continue
			}

			// Find the strike that is x points away.
			buyLeg, err := FindByStrike(chain.Puts, (row2.Strike - spreadWidth))

			if err != nil {
				continue
			}

			// See if there is enough credit
			credit := row2.Bid - buyLeg.Ask

			if credit < minCredit {
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
func getPutCreditSpreadParms(screen models.Screener, lastQuote float64) (int, int, float64, float64, float64) {

	var widthIncrment float64 = 0.50
	var minDaysToExpire int = 0
	var maxDaysToExpire int = 10000
	var minCredit float64 = 0.01
	var spreadWidth float64 = 5.00
	var minSellStrike float64 = 0.00

	// See if we have a min strike price to sell
	percentAway, err := FindFilterItemValue("short-strike-percent-away", screen)

	if err == nil {

		// Figure out the strike price that is the min we can sell.
		var tmp float64 = lastQuote - (lastQuote * (percentAway.ValueNumber / 100))
		fraction := tmp - math.Floor(tmp)

		if fraction >= widthIncrment {
			minSellStrike = (math.Floor(tmp) + widthIncrment)
		} else {
			minSellStrike = math.Floor(tmp)
		}

	}

	// See if we have a spread width
	sw, err := FindFilterItemValue("spread-width", screen)

	if err == nil {
		spreadWidth = sw.ValueNumber
	}

	// See if we have a min credit
	mc, err := FindFilterItemValue("min-credit", screen)

	if err == nil {
		minCredit = mc.ValueNumber
	}

	// See if we have max days to expire
	mde, err := FindFilterItemValue("max-days-to-expire", screen)

	if err == nil {
		maxDaysToExpire = int(mde.ValueNumber)
	}

	// See if we have min days to expire
	minde, err := FindFilterItemValue("min-days-to-expire", screen)

	if err == nil {
		minDaysToExpire = int(minde.ValueNumber)
	}

	// Return values
	return minDaysToExpire, maxDaysToExpire, minCredit, spreadWidth, minSellStrike
}

/* End File */
