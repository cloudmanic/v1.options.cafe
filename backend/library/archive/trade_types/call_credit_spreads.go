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
func CallCreditSpreadGetMaxRisked(positions *[]models.Position) (float64, float64) {

	var qty int = 0
	var cost float64 = 0.00
	var strikes []float64

	// Loop through the different positions and get some summary data
	for _, row := range *positions {

		// Parse the option
		option, _ := helpers.OptionParse(row.Symbol.ShortName)
		cost += row.CostBasis
		strikes = append(strikes, option.Strike)

		// Get the qty
		if row.OrgQty > 0 {
			qty = row.OrgQty
		}

	}

	// Get max, min. and diff
	max := FindMaxStike(strikes)
	min := FindMinStike(strikes)
	dif := max - min

	// Get max risk before credit
	maxRisk := float64(qty) * dif * 100.00

	// Return happy.
	return maxRisk + cost, 0.00
}

//
// Detect if this trade is an call credit spread
//
func IsCallCreditSpread(positions *[]models.Position) bool {

	tradeCost := 0.00

	var firstLeg helpers.OptionParts

	// Most only be 2 legs
	if len(*positions) != 2 {
		return false
	}

	// Loop through the different positions and get some summary data
	for key, row := range *positions {

		// Add up total trade costs
		tradeCost = tradeCost + row.CostBasis

		// Parse the option symbol
		option, _ := helpers.OptionParse(row.Symbol.ShortName)

		// Store the first leg
		if key == 0 {
			firstLeg = option
		}

		// If this is a call we know this is not a put credit spread
		if option.Type != "Call" {
			return false
		}

		// Make sure the expire date is the same
		if firstLeg.Expire != option.Expire {
			return false
		}
	}

	// Total cost must be negative (the trade was a credit)
	if tradeCost >= 0 {
		return false
	}

	// If we made it here we know it is a put credit spread
	return true
}

/* End File */
