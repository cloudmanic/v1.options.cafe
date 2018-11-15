//
// Date: 2/9/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package trade_types

import (
	"github.com/cloudmanic/app.options.cafe/backend/library/helpers"
	"github.com/cloudmanic/app.options.cafe/backend/models"
)

//
// Calculate max loss of the trade.
//
func GetLongCallButterflyRiskProfile(positions *[]models.Position) (float64, float64) {

	var cost float64 = 0.00

	// Loop through the different positions and get some summary data
	for _, row := range *positions {
		cost += row.CostBasis
	}

	// Return happy.
	return cost, 0
}

//
// Detect if this trade is an call credit spread
//
func IsLongCallButterfly(positions *[]models.Position) bool {

	tradeCost := 0.00
	count := 0

	var upLeg helpers.OptionParts
	var downLeg helpers.OptionParts
	var middleLeg helpers.OptionParts

	// Most only be 3 legs
	if len(*positions) != 3 {
		return false
	}

	// Loop through the different positions and get some summary data
	for _, row := range *positions {

		// Add up total trade costs
		tradeCost = tradeCost + row.CostBasis

		// Parse the option symbol
		option, _ := helpers.OptionParse(row.Symbol.ShortName)

		// Is this the middle leg
		if row.OrgQty < 0 {
			middleLeg = option
		}

		if (len(downLeg.Name) == 0) && (len(upLeg.Name) > 0) && (row.OrgQty > 0) {
			downLeg = option
		}

		if (len(upLeg.Name) == 0) && (row.OrgQty > 0) {
			upLeg = option
		}

		// If this is a call we know this is not a put credit spread
		if option.Type != "Call" {
			return false
		}

		// Keep track of the qty count
		count = count + row.OrgQty
	}

	// Reset the non-middle legs
	if downLeg.Strike > upLeg.Strike {
		down := upLeg.Strike
		up := downLeg.Strike
		downLeg.Strike = down
		upLeg.Strike = up
	}

	// Make sure all the expires are the same date
	if (downLeg.Expire != middleLeg.Expire) || (downLeg.Expire != upLeg.Expire) {
		return false
	}

	// Make sure the middle leg is 2x the other legs.
	if count != 0 {
		return false
	}

	// Total cost must be positive (the trade was a debit)
	if tradeCost <= 0 {
		return false
	}

	// If we made it here we know it is a put credit spread
	return true
}

/* End File */
