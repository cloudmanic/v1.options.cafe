//
// Date: 2/9/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package archive

import (
	"app.options.cafe/library/archive/trade_types"
	"app.options.cafe/models"
)

//
// Pass in some positions and return the amount risked in this trade.
//
func GetAmountRiskedInTrade(positions *[]models.Position) (float64, float64) {

	// Get the trade type
	tradeType := ClassifyTradeGroup(positions)

	// Based on the trade type we call different functions
	switch tradeType {

	case "Equity":
		return trade_types.GetSingleStockRiskProfile(positions)

	case "Option":
		return trade_types.GetSingleOptionRiskProfile(positions)

	case "Long Call Butterfly":
		return trade_types.GetLongCallButterflyRiskProfile(positions)

	case "Long Put Butterfly":
		return trade_types.GetLongPutButterflyRiskProfile(positions)

	case "Put Credit Spread":
		return trade_types.GetPutCreditSpreadRiskProfile(positions)

	case "Call Credit Spread":
		return trade_types.GetCallCreditSpreadRiskProfile(positions)

	case "Iron Condor":
		return trade_types.GetIronCondorRiskProfile(positions)

	case "Reverse Iron Condor":
		return trade_types.GetReverseIronCondorRiskProfile(positions)

	}

	// Should never make it here
	return 0.00, 0.00
}

//
// Loop through the positions and try to figure out what type of trade this is.
//
// Note: The order we check positions in matters.
//
func ClassifyTradeGroup(positions *[]models.Position) string {

	// single stock trade
	if trade_types.IsSingleStock(positions) {
		return "Equity"
	}

	// single short stock trade
	if trade_types.IsSingleShortStock(positions) {
		return "Short Equity"
	}

	// single option trade
	if trade_types.IsSingleOption(positions) {
		return "Option"
	}

	// long call butterfly trade
	if trade_types.IsLongCallButterfly(positions) {
		return "Long Call Butterfly"
	}

	// long put butterfly trade
	if trade_types.IsLongPutButterfly(positions) {
		return "Long Put Butterfly"
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

	// iron condor
	if trade_types.IsIronCondor(positions) {
		return "Iron Condor"
	}

	// reverse iron condor
	if trade_types.IsReverseIronCondor(positions) {
		return "Reverse Iron Condor"
	}

	// We could not figure out what this trade group was.
	return "Other"
}

/* End File */
