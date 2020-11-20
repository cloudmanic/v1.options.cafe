//
// Date: 10/27/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package trade_types

import (
	"app.options.cafe/library/helpers"
	"app.options.cafe/models"
)

//
// Calculate max loss of the trade.
//
func GetIronCondorRiskProfile(positions *[]models.Position) (float64, float64) {

	var qty int = 0
	var cost float64 = 0.00
	var maxRisk float64 = 0.00
	var callShortStrike float64 = 0
	var putShortStrike float64 = 0
	var callLongStrike float64 = 0
	var putLongStrike float64 = 0

	// Loop through the different positions and get some summary data
	for _, row := range *positions {

		// Set Cost
		cost += row.CostBasis

		// Parse the option
		option, _ := helpers.OptionParse(row.Symbol.ShortName)

		// Assign the known strikes.
		if (option.Type == "Call") && (row.OrgQty > 0) {
			callLongStrike = option.Strike
		}

		if (option.Type == "Call") && (row.OrgQty < 0) {
			callShortStrike = option.Strike
		}

		if (option.Type == "Put") && (row.OrgQty > 0) {
			putLongStrike = option.Strike
		}

		if (option.Type == "Put") && (row.OrgQty < 0) {
			putShortStrike = option.Strike
		}

		// Get the qty
		if row.OrgQty > 0 {
			qty = row.OrgQty
		}
	}

	// Get call diff
	callDiff := callLongStrike - callShortStrike

	// Get put diff
	putDiff := putShortStrike - putLongStrike

	// Get max risk before credit
	if callDiff > putDiff {
		maxRisk = float64(qty) * callDiff * 100.00
	} else {
		maxRisk = float64(qty) * putDiff * 100.00
	}

	// Return happy.
	return maxRisk + cost, (cost * -1.00)
}

//
// Detect if this trade is an Iron Condor
//
func IsIronCondor(positions *[]models.Position) bool {

	tradeCost := 0.00
	puts := 0
	calls := 0
	shortCalls := 0
	shortPuts := 0

	var callShortStrike float64 = 0
	var putShortStrike float64 = 0
	var callLongStrike float64 = 0
	var putLongStrike float64 = 0

	var firstLeg helpers.OptionParts

	// Most only be 2 legs
	if len(*positions) != 4 {
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

		// Assign the known strikes.
		if (option.Type == "Call") && (row.OrgQty > 0) {
			callLongStrike = option.Strike
		}

		if (option.Type == "Call") && (row.OrgQty < 0) {
			callShortStrike = option.Strike
		}

		if (option.Type == "Put") && (row.OrgQty > 0) {
			putLongStrike = option.Strike
		}

		if (option.Type == "Put") && (row.OrgQty < 0) {
			putShortStrike = option.Strike
		}

		// Count number of calls
		if option.Type == "Call" {
			calls++

			if row.CostBasis > 0 {
				shortCalls++
			}
		}

		// Count number of puts
		if option.Type == "Put" {
			puts++

			if row.CostBasis > 0 {
				shortPuts++
			}
		}

		// Make sure the expire date is the same
		if firstLeg.Expire != option.Expire {
			return false
		}

	}

	// Make sure this is not a reverse iron condor
	if callShortStrike > callLongStrike {
		return false
	}

	if putShortStrike < putLongStrike {
		return false
	}

	// make sure we have 2 calls
	if calls != 2 {
		return false
	}

	// make sure we have 2 puts
	if puts != 2 {
		return false
	}

	// make sure we have 1 short call
	if shortCalls != 1 {
		return false
	}

	// make sure we have 1 short put
	if shortPuts != 1 {
		return false
	}

	// Total cost must be positive (the trade was a credit)
	if tradeCost > 0 {
		return false
	}

	// If we made it here we know it is an Iron Condor
	return true
}

/* End File */
