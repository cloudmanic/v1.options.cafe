//
// Date: 2019-02-22
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package backtesting

import (
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/cloudmanic/app.options.cafe/backend/screener"
)

//
// OpenMultiLegCredit - Open a new spread by adding a position
//
func (t *Base) OpenMultiLegCredit(today time.Time, backtest *models.Backtest, result screener.Result) {

	// TODO(spicer): make this work from configs
	lots := 1

	// TODO(spicer): figure which price to use to open
	openPrice := result.MidPoint * 100 * float64(lots)

	// First see if we already have this position
	if !t.checkForCurrentPosition(backtest, result) {
		return
	}

	// Get the count of open positions
	posCount := t.openPositionsCount(backtest)

	// Only open one position at a time. TODO(spicer): make this a config.
	if posCount > 0 {
		return
	}

	// Amount of margin left after trade is opened.
	diff := result.Legs[1].OptionStrike - result.Legs[0].OptionStrike
	var margin float64 = (diff * 100 * float64(lots)) - openPrice

	// Get total margin needed
	totalMarginNeeded := t.getTotalMarginUsed(backtest) + margin

	// Make sure we have enough margin to continue
	if totalMarginNeeded > backtest.EndingBalance {
		return
	}

	// Add position
	backtest.Positions = append(backtest.Positions, models.BacktestPosition{
		UserId:          backtest.UserId,
		Status:          "Open",
		OpenDate:        models.Date{today},
		OpenPrice:       openPrice,
		Margin:          margin,
		Legs:            result.Legs,
		Lots:            lots,
		PutPrecentAway:  result.PutPrecentAway,
		CallPrecentAway: 0,
		Balance:         (backtest.EndingBalance + openPrice),
	})

	// Update ending balance
	backtest.EndingBalance = backtest.EndingBalance + openPrice

	//fmt.Println(today.Format("2006-01-02"), " : ", backtest.EndingBalance, " / ", totalMarginNeeded, " / ", margin, " / ", backtest.Positions[len(backtest.Positions)-1].OpenPrice)
}

// -------------- Private Helper Functions ------------------- //

//
// openPositionsCount - Return a count of how many open trades we have.
//
func (t *Base) openPositionsCount(backtest *models.Backtest) int {
	var count int = 0

	// Loop through current positons and search for this leg.
	for _, row := range backtest.Positions {

		// Ignored closed trades
		if row.Status == "Closed" {
			continue
		}

		// Update count
		count++
	}

	// Return happy
	return count
}

//
// checkForCurrentPosition - Check to make sure a position is not already on.
//
func (t *Base) checkForCurrentPosition(backtest *models.Backtest, result screener.Result) bool {

	// Loop through the legs
	for _, row := range result.Legs {
		// Loop through current positons and search for this leg.
		for _, row2 := range backtest.Positions {

			// Ignored closed trades
			if row2.Status == "Closed" {
				continue
			}

			// Loop through the legs of the positions
			for _, row3 := range row2.Legs {
				if row3.ShortName == row.ShortName {
					return false
				}
			}

		}
	}

	// Position not found if we made it here
	return true
}

//
// getTotalMarginUsed - Return a value for the total margin being used right now.
//
func (t *Base) getTotalMarginUsed(backtest *models.Backtest) float64 {
	total := 0.00

	// Loop for expired postions
	for _, row := range backtest.Positions {

		if row.Status == "Closed" {
			continue
		}

		total += row.Margin
	}

	// Return happy
	return total
}

/* End File */
