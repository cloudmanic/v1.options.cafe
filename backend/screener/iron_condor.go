//
// Date: 2018-10-27
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-11-07
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package screener

import (
	"flag"
	"sort"
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/library/helpers"
	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/cloudmanic/app.options.cafe/backend/models"
)

//
// Run Iron Condor screen
//
func (t *Base) RunIronCondor(screen models.Screener) ([]Result, error) {

	result := []Result{}

	// Set today's date
	today := time.Now()

	// Change today's date for unit testing.
	if flag.Lookup("test.v") != nil {
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
		services.Warning(err)
		return result, err
	}

	// Add default values
	t.IronCondorFillDefault(&screen)

	// Loop through the expire dates
	for _, row := range expires {

		// Expire Date.
		expireDate, _ := time.Parse("2006-01-02", row)

		// Filter for expire dates
		if !t.FilterDaysToExpireDaysToExpire(today, screen, expireDate) {
			continue
		}

		// Get options Chain
		chain, err := t.Broker.GetOptionsChainByExpiration(screen.Symbol, row)

		if err != nil {
			continue
		}

		// Get all possible Put Legs
		putLegs := t.GetPossibleVerticalSpreads(screen, quote, chain.Puts, "put", "credit")

		if len(putLegs) <= 0 {
			continue
		}

		// Get all possible Call Legs
		callLegs := t.GetPossibleVerticalSpreads(screen, quote, chain.Calls, "call", "credit")

		if len(callLegs) <= 0 {
			continue
		}

		trades := t.FindPossibleIronCondorTrades(screen, putLegs, callLegs)

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
			closeCost := (row.PutShort.Ask - row.PutLong.Bid) + (row.CallShort.Ask - row.CallLong.Bid)
			openCost := (row.PutShort.Bid - row.PutLong.Ask) + (row.CallShort.Bid - row.CallLong.Ask)
			midPoint := (openCost + closeCost) / 2

			// Percent away - We show the lowest percent away
			putPercentAway := ((1 - (row.PutShort.Strike / quote.Last)) * 100)
			callPercentAway := ((1 - (quote.Last / row.CallShort.Strike)) * 100)

			// We have a winner
			result = append(result, Result{
				Expired:         symbPutLongLeg.OptionExpire,
				Credit:          helpers.Round(openCost, 2),
				MidPoint:        helpers.Round(midPoint, 2),
				PutPrecentAway:  helpers.Round(putPercentAway, 2),
				CallPrecentAway: helpers.Round(callPercentAway, 2),
				Legs:            []models.Symbol{symbPutLongLeg, symbPutShortLeg, symbCallShortLeg, symbCallLongLeg},
			})

		}

	}

	// Sort the results expire in asc order.
	sort.Slice(result, func(i, j int) bool {

		// Deal with tied sorts
		if result[i].Expired.Unix() == result[j].Expired.Unix() {
			return result[i].MidPoint < result[j].MidPoint
		}

		return result[i].Expired.Unix() < result[j].Expired.Unix()
	})

	// Return happy
	return result, nil
}

// ------------------------ Helper Functions -------------------------- //

//
// Setup default values. We need to make sure we have at least these params to run a backtest.
//
func (t *Base) IronCondorFillDefault(screen *models.Screener) {

	// Map found
	found := map[string]bool{}

	// Fields that are required
	required := map[string]models.ScreenerItem{
		"open-credit":           {Key: "open-credit", Operator: ">", ValueNumber: 0.10},
		"put-leg-width":         {Key: "put-leg-width", Operator: "=", ValueNumber: 2.00},
		"call-leg-width":        {Key: "call-leg-width", Operator: "=", ValueNumber: 2.00},
		"put-leg-percent-away":  {Key: "put-leg-percent-away", Operator: ">", ValueNumber: 2.00},
		"call-leg-percent-away": {Key: "call-leg-percent-away", Operator: ">", ValueNumber: 2.00},
	}

	// Loop through and identify items we already have
	for _, row := range screen.Items {

		if _, ok := required[row.Key]; ok {
			found[row.Key] = true
		}

	}

	// Add default values
	for key, row := range required {

		if _, ok := found[key]; !ok {
			screen.Items = append(screen.Items, row)
		}

	}

}

//
// Find possible Trades
//
func (t *Base) FindPossibleIronCondorTrades(screen models.Screener, putLegs []Spread, callLegs []Spread) []IronCondor {

	rt := []IronCondor{}
	mapIndex := make(map[string]IronCondor)

	// Loop through the PUT legs
	for _, row := range putLegs {

		side1Cost := row.Short.Bid - row.Long.Ask

		for _, row2 := range callLegs {

			side2Cost := row2.Short.Bid - row2.Long.Ask

			// Get total cost of the trade
			totalCost := side1Cost + side2Cost

			if !t.FilterOpenCredit(screen, totalCost) {
				continue
			}

			// if we made it here store as possible trade
			indexKey := helpers.FloatToString(row.Long.Strike) + "/" + helpers.FloatToString(row.Short.Strike) + "/" + helpers.FloatToString(row2.Short.Strike) + "/" + helpers.FloatToString(row2.Long.Strike)

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

		side1Cost := row.Short.Bid - row.Long.Ask

		for _, row2 := range putLegs {
			side2Cost := row2.Short.Bid - row2.Long.Ask

			// Get total cost of the trade
			totalCost := side1Cost + side2Cost

			if !t.FilterOpenCredit(screen, totalCost) {
				continue
			}

			// if we made it here store as possible trade
			indexKey := helpers.FloatToString(row2.Long.Strike) + "/" + helpers.FloatToString(row2.Short.Strike) + "/" + helpers.FloatToString(row.Short.Strike) + "/" + helpers.FloatToString(row.Long.Strike)

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
