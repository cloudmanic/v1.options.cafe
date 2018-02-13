//
// Date: 2/9/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package archive

import (
	"github.com/cloudmanic/app.options.cafe/backend/library/archive/trade_types"
	"github.com/cloudmanic/app.options.cafe/backend/models"
)

//
// Loop through the positions and try to figure out what type of trade this is.
//
// Note: The order we check positions in matters.
//
func ClassifyTradeGroup(positions *[]models.Position) string {

	// single stock trade
	if trade_types.IsSingleStock(positions) {
		return "Stock"
	}

	// single option trade
	if trade_types.IsSingleOption(positions) {
		return "Option"
	}

	// put credit spread
	if trade_types.IsPutCreditSpread(positions) {
		return "Put Credit Spread"
	}

	// put debit spread
	if trade_types.IsPutDebitSpread(positions) {
		return "Put Debit Spread"
	}

	// call credit spread
	if trade_types.IsCallCreditSpread(positions) {
		return "Call Credit Spread"
	}

	// call debit spread
	if trade_types.IsCallDebitSpread(positions) {
		return "Call Debit Spread"
	}

	// We could not figure out what this trade group was.
	return "Other"
}

/* End File */
