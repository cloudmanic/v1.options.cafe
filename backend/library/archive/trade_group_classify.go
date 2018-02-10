//
// Date: 2/9/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package archive

import (
	"github.com/cloudmanic/app.options.cafe/backend/library/helpers"
	"github.com/cloudmanic/app.options.cafe/backend/models"
)

//
// Loop through the positions and try to figure out what type of trade this is.
//
// Note: The order we check positions in matters.
//
func ClassifyTradeGroup(positions *[]models.Position) string {

	// single stock trade
	if IsSingleStock(positions) {
		return "Stock"
	}

	// single option trade
	if IsSingleOption(positions) {
		return "Option"
	}

	// put credit spread
	if IsPutCreditSpread(positions) {
		return "Put Credit Spread"
	}

	// put debit spread
	if IsPutDebitSpread(positions) {
		return "Put Debit Spread"
	}

	// call credit spread
	if IsCallCreditSpread(positions) {
		return "Call Credit Spread"
	}

	// call debit spread
	if IsCallDebitSpread(positions) {
		return "Call Debit Spread"
	}

	// We could not figure out what this trade group was.
	return "Other"
}

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
		_, err := helpers.OptionParse(row.Symbol)

		// If we can't parse the option we assume it is a stock
		if err != nil {
			return true
		}
	}

	// If we made it here we know it is not a stock
	return false
}

//
// Detect if this trade is single option trade
//
func IsSingleOption(positions *[]models.Position) bool {

	// Most only be 1 leg
	if len(*positions) != 1 {
		return false
	}

	// Parse the option symbol
	for _, row := range *positions {
		_, err := helpers.OptionParse(row.Symbol)

		// If we can't parse the option we assume it is not an option
		if err != nil {
			return false
		}
	}

	// If we made it here we know it is an option
	return true
}

//
// Detect if this trade is an put credit spread
//
func IsPutCreditSpread(positions *[]models.Position) bool {

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
		option, _ := helpers.OptionParse(row.Symbol)

		// Store the first leg
		if key == 0 {
			firstLeg = option
		}

		// If this is a call we know this is not a put credit spread
		if option.Type != "Put" {
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

//
// Detect if this trade is an put debit spread
//
func IsPutDebitSpread(positions *[]models.Position) bool {

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
		option, _ := helpers.OptionParse(row.Symbol)

		// Store the first leg
		if key == 0 {
			firstLeg = option
		}

		// If this is a call we know this is not a put credit spread
		if option.Type != "Put" {
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
		option, _ := helpers.OptionParse(row.Symbol)

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
		option, _ := helpers.OptionParse(row.Symbol)

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
