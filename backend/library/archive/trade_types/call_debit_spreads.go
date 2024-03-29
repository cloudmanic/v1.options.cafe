//
// Date: 2/9/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package trade_types

import (
	"app.options.cafe/library/helpers"
	"app.options.cafe/models"
)

//
// Detect if this trade is an call debit spread
//
func IsCallDebitSpread(positions *[]models.Position) bool {

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
	if tradeCost <= 0 {
		return false
	}

	// If we made it here we know it is a put credit spread
	return true
}

/* End File */
