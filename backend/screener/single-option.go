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

//
// SingleOptionCoreScreen will run our core screen and return results.
//
func (t *Base) SingleOptionCoreScreen(today time.Time, expireDate time.Time, screen models.Screener, options []types.OptionsChainItem, underlyingLast float64, cache screenerCache.Cache) []Result {
	result := []Result{}

	// Figure out best strike
	atmStrike := math.RoundToEven(underlyingLast)

	putOption := t.GetOptionByExpirationDateAndStrike(expireDate, atmStrike, "Put", options)
	//callOption := t.GetOptionByExpirationDateAndStrike(expireDate, atmStrike, "Call", options)
	openCost := putOption.Ask
	closeCost := putOption.Bid
	midPoint := (openCost + closeCost) / 2

	// No need for these options.
	if putOption.Ask <= 0.00 {
		return result
	}

	// Add in Symbol Object - Leg 1
	leg, err := cache.GetSymbol(putOption.Symbol, putOption.Description, "Option")

	if err != nil {
		return result
	}

	// We have a winner
	result = append(result, Result{
		Day:             types.Date{today},
		Debit:           helpers.Round(openCost, 2),
		Bid:             helpers.Round(openCost, 2),
		Ask:             helpers.Round(closeCost, 2),
		MidPoint:        helpers.Round(midPoint, 2),
		UnderlyingLast:  underlyingLast,
		CallPrecentAway: 0.00,
		PutPrecentAway:  0.00,
		Legs:            []models.Symbol{leg},
	})

	// Return happy
	return result
}

/* End File */
