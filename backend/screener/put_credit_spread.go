//
// Date: 2018-07-17
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package screener

import (
	"fmt"
	"os"
	"strings"
	"time"

	"app.options.cafe/brokers/types"
	"app.options.cafe/library/helpers"
	"app.options.cafe/library/services"
	"app.options.cafe/models"

	screenerCache "app.options.cafe/screener/cache"
)

//
// RunPutCreditSpread will run a put credit spread screen
//
func (t *Base) RunPutCreditSpread(screen models.Screener) ([]Result, error) {

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
		services.Info(err)
		return result, err
	}

	// Build the cache for this screen.
	cache := screenerCache.New(t.DB)

	// Loop through the expire dates
	for _, row := range expires {

		// Expire Date.
		expireDate, _ := time.Parse("2006-01-02", row)

		// Filter for expire dates
		if !t.FilterDaysToExpire(today, screen, expireDate) {
			continue
		}

		// Get options Chain
		chain, err := t.Broker.GetOptionsChainByExpiration(screen.Symbol, row)

		if err != nil {
			continue
		}

		// Core screen logic from our screener.
		for _, row2 := range t.PutCreditSpreadCoreScreen(today, screen, chain.Puts, quote.Last, cache) {
			result = append(result, row2)
		}

	}

	// Return happy
	return result, nil
}

//
// PutCreditSpreadCoreScreen will run our core screen and return results.
//
func (t *Base) PutCreditSpreadCoreScreen(today time.Time, screen models.Screener, puts []types.OptionsChainItem, underlyingLast float64, cache screenerCache.Cache) []Result {
	result := []Result{}

	// Set params
	spreadWidth := t.GetPutCreditSpreadParms(screen, underlyingLast)

	// Loop through all the puts.
	for _, row := range puts {
		// No need to pay attention to open interest of zero
		if row.OpenInterest == 0 {
			continue
		}

		// Skip strikes that are higher than our min strike. Based on percent away.
		if !t.FilterStrikeByPercentDown("short-strike-percent-away", screen, row.Strike, underlyingLast) {
			continue
		}

		// We only want the first 100
		if len(result) >= 100 {
			return result
		}

		// Find the strike that is x points away.
		buyLeg, err := t.FindByStrike(puts, (row.Strike - spreadWidth))

		if err != nil {
			continue
		}

		// See if there is enough credit
		credit := row.Bid - buyLeg.Ask

		if !t.FilterOpenCredit(screen, credit) {
			continue
		}

		// Figure out the credit spread amount.
		buyCost := row.Ask - buyLeg.Bid
		midPoint := (credit + buyCost) / 2

		// Add in Symbol Object - Buy leg
		symbBuyLeg, err := cache.GetSymbol(buyLeg.Symbol, buyLeg.Description, "Option")

		if err != nil {
			fmt.Println("HERE 1")
			continue
		}

		// Add in Symbol Object - Sell leg
		symbSellLeg, err := cache.GetSymbol(row.Symbol, row.Description, "Option")

		if err != nil {
			fmt.Println("HERE 2")
			continue
		}

		// We have a winner
		result = append(result, Result{
			Day:            types.Date{today},
			Credit:         helpers.Round(credit, 2),
			Bid:            helpers.Round(buyCost, 2),
			Ask:            helpers.Round(credit, 2),
			MidPoint:       helpers.Round(midPoint, 2),
			UnderlyingLast: underlyingLast,
			PutPrecentAway: helpers.Round(((1 - row.Strike/underlyingLast) * 100), 2),
			Legs:           []models.Symbol{symbBuyLeg, symbSellLeg},
		})
	}

	// Return happy
	return result
}

//
// GetPutCreditSpreadParms will set Parms we need.
//
func (t *Base) GetPutCreditSpreadParms(screen models.Screener, lastQuote float64) float64 {

	spreadWidth := 5.00

	// See if we have a spread width
	sw, err := t.FindFilterItemValue("spread-width", screen)

	if err == nil {
		spreadWidth = sw.ValueNumber
	}

	// Return values
	return spreadWidth
}

/* End File */
