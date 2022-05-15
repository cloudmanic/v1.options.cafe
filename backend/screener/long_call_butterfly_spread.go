//
// Date: 2022-05-15
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Copyright: 2022 Cloudmanic Labs, LLC. All rights reserved.
//

package screener

import (
	"math"
	"time"

	"app.options.cafe/brokers/types"
	"app.options.cafe/library/helpers"
	"app.options.cafe/models"
	screenerCache "app.options.cafe/screener/cache"
)

// //
// // RunPutCreditSpread will run a put credit spread screen
// //
// func (t *Base) RunPutCreditSpread(screen models.Screener) ([]Result, error) {

// 	result := []Result{}

// 	// // Set today's date
// 	// today := time.Now()

// 	// // Change today's date for unit testing.
// 	// if strings.HasSuffix(os.Args[0], ".test") {
// 	// 	today = helpers.ParseDateNoError("2018-10-18").UTC()
// 	// }

// 	// // Make call to get current quote.
// 	// quote, err := t.GetQuote(screen.Symbol)

// 	// if err != nil {
// 	// 	return result, err
// 	// }

// 	// // Get all possible expire dates.
// 	// expires, err := t.Broker.GetOptionsExpirationsBySymbol(screen.Symbol)

// 	// if err != nil {
// 	// 	services.Info(err)
// 	// 	return result, err
// 	// }

// 	// // Build the cache for this screen.
// 	// cache := screenerCache.New(t.DB)

// 	// // Loop through the expire dates
// 	// for _, row := range expires {

// 	// 	// Expire Date.
// 	// 	expireDate, _ := time.Parse("2006-01-02", row)

// 	// 	// Filter for expire dates
// 	// 	if !t.FilterDaysToExpire(today, screen, expireDate) {
// 	// 		continue
// 	// 	}

// 	// 	// Get options Chain
// 	// 	chain, err := t.Broker.GetOptionsChainByExpiration(screen.Symbol, row)

// 	// 	if err != nil {
// 	// 		continue
// 	// 	}

// 	// 	// Core screen logic from our screener.
// 	// 	for _, row2 := range t.PutCreditSpreadCoreScreen(today, screen, chain.Puts, quote.Last, cache) {
// 	// 		result = append(result, row2)
// 	// 	}

// 	// }

// 	// Return happy
// 	return result, nil
// }

//
// LongCallButterflySpreadCoreScreen will run our core screen and return results.
//
func (t *Base) LongCallButterflySpreadCoreScreen(today time.Time, screen models.Screener, calls []types.OptionsChainItem, underlyingLast float64, cache screenerCache.Cache) []Result {
	result := []Result{}

	// Loop through all the puts.
	for _, row := range calls {
		// No need to pay attention to open interest of zero
		if row.OpenInterest == 0 {
			continue
		}

		// Skip strikes that are higher than our min strike. Based on percent away.
		if !t.FilterStrikeByPercentDown("left-strike-percent-away", screen, row.Strike, underlyingLast) {
			continue
		}

		// We only want the first 100
		if len(result) >= 100 {
			return result
		}

		// Find diff between this strike and the current stock price.
		roundedLastPrice := math.RoundToEven(underlyingLast)
		diff := roundedLastPrice - row.Strike
		legCStrike := roundedLastPrice + diff

		// Find the strike at current stock price.
		legB, err := t.FindByStrike(calls, roundedLastPrice)

		if err != nil {
			continue
		}

		// Find the strike that is x points away.
		legC, err := t.FindByStrike(calls, legCStrike)

		if err != nil {
			continue
		}

		// Figure out open debit
		openCost := row.Ask - (legB.Bid * 2) + legC.Ask
		closeCost := (legB.Ask * 2) - row.Bid - legC.Bid
		midPoint := (openCost + closeCost) / 2

		// See if this is too expensive
		if !t.FilterOpenDebit(screen, openCost) {
			continue
		}

		//fmt.Println(math.RoundToEven(underlyingLast), " ", row.Strike, " ", diff, " ", legCStrike, " : ", openCost)

		// Add in Symbol Object - Leg A
		legASym, err := cache.GetSymbol(row.Symbol, row.Description, "Option")

		if err != nil {
			continue
		}

		// Add in Symbol Object - Leg B
		legBSym, err := cache.GetSymbol(legB.Symbol, legB.Description, "Option")

		if err != nil {
			continue
		}

		// Add in Symbol Object - Leg C
		legCSym, err := cache.GetSymbol(legC.Symbol, legC.Description, "Option")

		if err != nil {
			continue
		}

		// We have a winner
		result = append(result, Result{
			Day:            types.Date{today},
			Debit:          helpers.Round(openCost, 2),
			Bid:            helpers.Round(openCost, 2),
			Ask:            helpers.Round(closeCost, 2),
			MidPoint:       helpers.Round(midPoint, 2),
			UnderlyingLast: underlyingLast,
			PutPrecentAway: helpers.Round(((1 - row.Strike/underlyingLast) * 100), 2),
			Legs:           []models.Symbol{legASym, legBSym, legCSym},
		})
	}

	// Return happy
	return result
}

/* End File */
