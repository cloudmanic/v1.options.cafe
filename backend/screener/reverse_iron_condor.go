//
// Date: 2018-10-27
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-11-06
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package screener

import (
	"os"
	"strings"
	"time"

	"app.options.cafe/library/helpers"
	"app.options.cafe/library/services"
	"app.options.cafe/models"
)

//
// Run a put credit spread screen
//
func (t *Base) RunReverseIronCondor(screen models.Screener) ([]Result, error) {

	result := []Result{}

	// Set today's date
	today := time.Now()

	// Change today's date for unit testing.
	if strings.HasSuffix(os.Args[0], ".test") {
		today = helpers.ParseDateNoError("2018-10-18").UTC()
	}

	// Make call to get current quote.
	quote, err := t.GetQuote(screen.Symbol)

	if err != nil {
		return result, err
	}

	// Get all possible expire dates.
	expires, err := t.Broker.GetOptionsExpirationsBySymbol(screen.Symbol)

	if err != nil {
		services.Info(err)
		return result, err
	}

	// Loop through the expire dates
	for _, row := range expires {

		// Expire Date.
		expireDate, _ := time.Parse("2006-01-02", row)

		// Filter for expire dates
		if !t.FilterDaysToExpire(today, screen, expireDate) {
			continue
		}

		// Get options Chain
		chain, err := t.Broker.GetOptionsChainByExpiration(screen.Symbol, row)

		if err != nil {
			continue
		}

		// Get all possible Put Legs
		putLegs := t.GetPossibleVerticalSpreads(screen, quote, chain.Puts, "put", "debit")

		if len(putLegs) <= 0 {
			continue
		}

		// Get all possible Call Legs
		callLegs := t.GetPossibleVerticalSpreads(screen, quote, chain.Calls, "call", "debit")

		if len(callLegs) <= 0 {
			continue
		}

		trades := t.FindPossibleReverseIronCondorTrades(screen, putLegs, callLegs)

		// Add trades to the results
		for _, row := range trades {

			// We only want the first 100
			if len(result) >= 100 {
				return result, nil
			}

			// Add in Symbol Object - Put Short leg
			symbPutShortLeg, err := t.DB.CreateNewSymbol(row.PutShort.Symbol, row.PutShort.Description, "Option")

			if err != nil {
				continue
			}

			// Add in Symbol Object - Put Long leg
			symbPutLongLeg, err := t.DB.CreateNewSymbol(row.PutLong.Symbol, row.PutLong.Description, "Option")

			if err != nil {
				continue
			}

			// Add in Symbol Object - Call Long leg
			symbCallLongLeg, err := t.DB.CreateNewSymbol(row.CallLong.Symbol, row.CallLong.Description, "Option")

			if err != nil {
				continue
			}

			// Add in Symbol Object - Call Short leg
			symbCallShortLeg, err := t.DB.CreateNewSymbol(row.CallShort.Symbol, row.CallShort.Description, "Option")

			if err != nil {
				continue
			}

			// Figure out the amounts.
			debitCost := (row.PutLong.Ask - row.PutShort.Bid) + (row.CallLong.Ask - row.CallShort.Bid)
			creditCost := (row.PutLong.Bid - row.PutShort.Ask) + (row.CallLong.Bid - row.CallShort.Ask)
			midPoint := (creditCost + debitCost) / 2

			// Percent away - We show the lowest percent away
			putPercentAway := ((1 - (row.PutLong.Strike / quote.Last)) * 100)
			callPercentAway := ((1 - (quote.Last / row.CallLong.Strike)) * 100)

			// We have a winner
			result = append(result, Result{
				Debit:           helpers.Round(debitCost, 2),
				MidPoint:        helpers.Round(midPoint, 2),
				PutPrecentAway:  helpers.Round(putPercentAway, 2),
				CallPrecentAway: helpers.Round(callPercentAway, 2),
				Legs:            []models.Symbol{symbPutShortLeg, symbPutLongLeg, symbCallLongLeg, symbCallShortLeg},
			})

		}

	}

	// Return happy
	return result, nil
}

// ------------------------ Helper Functions -------------------------- //

//
// Find possible Trades
//
func (t *Base) FindPossibleReverseIronCondorTrades(screen models.Screener, putLegs []Spread, callLegs []Spread) []IronCondor {

	rt := []IronCondor{}
	mapIndex := make(map[string]IronCondor)

	// Loop through the PUT legs
	for _, row := range putLegs {

		side1Cost := row.Long.Ask - row.Short.Bid

		for _, row2 := range callLegs {

			side2Cost := row2.Long.Ask - row2.Short.Bid

			// Get total cost of the trade
			totalCost := side1Cost + side2Cost

			if !t.FilterOpenDebit(screen, totalCost) {
				continue
			}

			// if we made it here store as possible trade
			indexKey := helpers.FloatToString(row.Short.Strike) + "/" + helpers.FloatToString(row.Long.Strike) + "/" + helpers.FloatToString(row2.Long.Strike) + "/" + helpers.FloatToString(row2.Short.Strike)

			mapIndex[indexKey] = IronCondor{
				CallShort: row2.Short,
				CallLong:  row2.Long,
				PutShort:  row.Short,
				PutLong:   row.Long,
			}

		}

	}

	// Loop through the CALL legs
	for _, row := range callLegs {

		side1Cost := row.Long.Ask - row.Short.Bid

		for _, row2 := range putLegs {
			side2Cost := row2.Long.Ask - row2.Short.Bid

			// Get total cost of the trade
			totalCost := side1Cost + side2Cost

			if !t.FilterOpenDebit(screen, totalCost) {
				continue
			}

			// if we made it here store as possible trade
			indexKey := helpers.FloatToString(row.Short.Strike) + "/" + helpers.FloatToString(row.Long.Strike) + "/" + helpers.FloatToString(row2.Long.Strike) + "/" + helpers.FloatToString(row2.Short.Strike)

			mapIndex[indexKey] = IronCondor{
				CallShort: row.Short,
				CallLong:  row.Long,
				PutShort:  row2.Short,
				PutLong:   row2.Long,
			}

		}

	}

	// We used a hash map to remove duplicates now normalize it.
	for _, row := range mapIndex {
		rt = append(rt, row)
	}

	// Return happy
	return rt
}

/* End File */
