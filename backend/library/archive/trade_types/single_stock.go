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
// Detect if this trade is single stock trade
//
func IsSingleStock(positions *[]models.Position) bool {

	// Most only be 1 leg
	if len(*positions) != 1 {
		return false
	}

	// Parse the option symbol
	for _, row := range *positions {
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
