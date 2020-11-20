//
// Date: 2019-02-22
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package backtesting

import (
	"math"
	"strconv"
	"strings"
	"time"

	"app.options.cafe/brokers/eod"
	"app.options.cafe/models"
	"app.options.cafe/screener"
	"github.com/davecgh/go-spew/spew"
)

//
// OpenMultiLegCredit - Open a new spread by adding a position
//
func (t *Base) OpenMultiLegCredit(today time.Time, strategy string, backtest *models.Backtest, result screener.Result) {
	// Default values
	lots := 1

	// First see if we already have this position
	if !t.checkForCurrentPosition(backtest, result) {
		return
	}

	// Set up a new screener so we can use it's Functions
	screenObj := screener.NewScreen(t.DB, &eod.Api{})

	if result.Credit == 0 {
		spew.Dump(result)
	}

	// Amount of margin left after trade is opened.
	diff := result.Legs[1].OptionStrike - result.Legs[0].OptionStrike

	// Figure out position size
	if backtest.PositionSize == "one-at-time" {
		// Get the count of open positions
		posCount := t.openPositionsCount(backtest)

		// Only open one position at a time. TODO(spicer): make this a config.
		if posCount > 0 {
			return
		}
	} else if strings.Contains(backtest.PositionSize, "percent") { // percent of trade
		totalToTrade := t.percentOfAccount(backtest, backtest.PositionSize)
		lots = int(math.Floor(totalToTrade / (diff * 100.00)))
	}

	// Figure out open price.
	openPrice := result.Ask * 100 * float64(lots)

	// we can configure to use midpoint.
	if backtest.Midpoint {
		openPrice = result.MidPoint * 100 * float64(lots)
	}

	// Get margin used
	margin := (diff * 100 * float64(lots)) - openPrice

	// Get total margin needed
	totalMarginNeeded := t.getTotalMarginUsed(backtest) + margin

	// Make sure we have enough margin to continue
	if totalMarginNeeded > backtest.EndingBalance {
		return
	}

	// See if we are allowed to only have one strike
	st, err := screenObj.FindFilterItemValue("allow-more-than-one-strike", backtest.Screen)

	// Make sure we do not already have a trade on at this strike.
	if (err == nil) && (st.ValueString == "no") {
		for _, row := range backtest.Positions {
			// Only open trades
			if row.Status != "Open" {
				continue
			}

			for _, row2 := range row.Legs {
				if row2.OptionStrike == result.Legs[0].OptionStrike {
					return
				}

				if row2.OptionStrike == result.Legs[1].OptionStrike {
					return
				}
			}
		}
	}

	// See if we are allowed to only have one expire
	st2, err := screenObj.FindFilterItemValue("allow-more-than-one-expire", backtest.Screen)

	// Make sure we do not already have a trade on at this expire.
	if (err == nil) && (st2.ValueString == "no") {
		for _, row := range backtest.Positions {
			// Only open trades
			if row.Status != "Open" {
				continue
			}

			for _, row2 := range row.Legs {
				if row2.OptionExpire == result.Legs[0].OptionExpire {
					return
				}

				if row2.OptionExpire == result.Legs[1].OptionExpire {
					return
				}
			}
		}
	}

	// Add position
	backtest.Positions = append(backtest.Positions, models.BacktestPosition{
		UserId:          backtest.UserId,
		Strategy:        strategy,
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
// percentOfAccount - will return a percent of my total balance.
//
func (t *Base) percentOfAccount(backtest *models.Backtest, percentString string) float64 {
	// Split string
	y := strings.Split(percentString, "-")

	// If this is not percent return 0.00 (this should not happen)
	if y[1] != "percent" {
		return 0.00
	}

	// Convert from string to float
	percent, err := strconv.ParseFloat(y[0], 64)

	if err != nil {
		return 0.00
	}

	// Return percent of port we can trade.
	return backtest.EndingBalance * (percent / 100)
}

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
