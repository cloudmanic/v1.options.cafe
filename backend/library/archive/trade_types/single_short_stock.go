//
// Date: 9/20/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package trade_types

import (
	"app.options.cafe/library/helpers"
	"app.options.cafe/models"
)

//
// Calculate max loss of the trade. Return -1 for unlimited risk
//
func GetSingleShortStockRiskProfile(positions *[]models.Position) (float64, float64) {
	// Unlimited risk & gain
	return -1.00, 0.00
}

//
// Detect if this trade is single stock trade
//
func IsSingleShortStock(positions *[]models.Position) bool {

	// Most only be 1 leg
	if len(*positions) != 1 {
		return false
	}

	// Parse the option symbol
	for _, row := range *positions {

		// Make sure there are no long positions
		if row.OrgQty > 0 {
			return false
		}

		// Check to see if this is an option
		_, err := helpers.OptionParse(row.Symbol.ShortName)

		// If we can't parse the option we assume it is a stock
		if err != nil {
			return true
		}
	}

	// If we made it here we know it is not a stock
	return false
}

/* End File */
