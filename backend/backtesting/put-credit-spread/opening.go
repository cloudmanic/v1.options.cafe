//
// Date: 2019-02-22
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package putcreditspread

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"app.options.cafe/backtesting/helpers"
	"app.options.cafe/brokers/eod"
	"app.options.cafe/models"
	"app.options.cafe/screener"
	"github.com/davecgh/go-spew/spew"
)

//
// OpenMultiLegCredit - Open a new spread by adding a position
//
func OpenMultiLegCredit(db models.Datastore, today time.Time, strategy string, backtest *models.Backtest, result screener.Result) {
	// Default values
	lots := 0

	// First see if we already have this position
	if !checkForCurrentPosition(backtest, result) {
		return
	}

	// Set up a new screener so we can use it's Functions
	screenObj := screener.NewScreen(db, &eod.Api{})

	if result.Credit == 0 {
		fmt.Println("#####")
		spew.Dump(result)
		fmt.Println("#####")
	}

	// Amount of margin left after trade is opened.
	diff := result.Legs[1].OptionStrike - result.Legs[0].OptionStrike

	// Figure out position size
	if backtest.PositionSize == "one-at-time" {
		// Get the count of open positions
		posCount := helpers.GetOpenTradeCount(backtest)

		// Only open one position at a time. TODO(spicer): make this a config.
		if posCount > 0 {
			return
		}
	} else if strings.Contains(backtest.PositionSize, "percent") { // percent of trade
		totalToTrade := percentOfAccount(backtest, backtest.PositionSize)
		lots = int(math.Floor(totalToTrade / (diff * 100.00)))
	}

	// If lots == 0 we have used all our margin
	if lots == 0 {
		return
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
	totalMarginNeeded := getTotalMarginUsed(backtest) + margin

	// Make sure we have enough margin to continue
	if totalMarginNeeded > backtest.EndingBalance {
		return
	}

	// See if we are allowed to only have one strike
	st, err := screenObj.FindFilterItemValue("allow-more-than-one-strike", backtest.Screen)

	// Make sure we do not already have a trade on at this strike.
	if (err == nil) && (st.ValueString == "no") {
		for _, row := range backtest.TradeGroups {
			// Only open trades
			if row.Status != "Open" {
				continue
			}

			for _, row2 := range row.Positions {
				if row2.Symbol.OptionStrike == result.Legs[0].OptionStrike {
					return
				}

				if row2.Symbol.OptionStrike == result.Legs[1].OptionStrike {
					return
				}
			}
		}
	}

	// See if we are allowed to only have one expire
	st2, err := screenObj.FindFilterItemValue("allow-more-than-one-expire", backtest.Screen)

	// Make sure we do not already have a trade on at this expire.
	if (err == nil) && (st2.ValueString == "no") {
		for _, row := range backtest.TradeGroups {
			// Only open trades
			if row.Status != "Open" {
				continue
			}

			for _, row2 := range row.Positions {
				if row2.Symbol.OptionExpire == result.Legs[0].OptionExpire {
					return
				}

				if row2.Symbol.OptionExpire == result.Legs[1].OptionExpire {
					return
				}
			}
		}
	}

	// Build legs
	legs := []models.BacktestPosition{}

	for key, row := range result.Legs {
		q := lots

		// Second leg is the short one
		if key == 1 {
			q = q * -1
		}

		legs = append(legs, models.BacktestPosition{
			UserId:   backtest.UserId,
			Status:   "Open",
			SymbolId: row.Id,
			Symbol:   row,
			OpenDate: today,
			Qty:      q,
		})
	}

	// Spread text TODO(spicer): Add a switch statement to customize this string per strategy.
	spreadText := fmt.Sprintf("%s %s $%.2f / $%.2f", legs[0].Symbol.OptionUnderlying, legs[0].Symbol.OptionExpire.Format("01/02/2006"), legs[0].Symbol.OptionStrike, legs[1].Symbol.OptionStrike)

	// Add position
	backtest.TradeGroups = append(backtest.TradeGroups, models.BacktestTradeGroup{
		UserId:          backtest.UserId,
		Strategy:        strategy,
		Status:          "Open",
		SpreadText:      spreadText,
		OpenDate:        models.Date{today},
		OpenPrice:       openPrice,
		Margin:          margin,
		Positions:       legs,
		Lots:            lots,
		Credit:          (openPrice / float64(lots)) / 100,
		PutPrecentAway:  result.PutPrecentAway,
		CallPrecentAway: 0,
		Balance:         (backtest.EndingBalance + openPrice),
	})

	// Update ending balance
	backtest.EndingBalance = backtest.EndingBalance + openPrice

	//fmt.Println(today.Format("2006-01-02"), " : ", backtest.EndingBalance, " / ", totalMarginNeeded, " / ", margin, " / ", backtest.TradeGroups[len(backtest.TradeGroups)-1].Credit)
}

// -------------- Private Helper Functions ------------------- //

//
// percentOfAccount - will return a percent of my total balance.
//
func percentOfAccount(backtest *models.Backtest, percentString string) float64 {
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
// checkForCurrentPosition - Check to make sure a position is not already on.
//
func checkForCurrentPosition(backtest *models.Backtest, result screener.Result) bool {

	// Loop through the legs
	for _, row := range result.Legs {
		// Loop through current positons and search for this leg.
		for _, row2 := range backtest.TradeGroups {

			// Ignored closed trades
			if row2.Status == "Closed" {
				continue
			}

			// Loop through the legs of the positions
			for _, row3 := range row2.Positions {
				if row3.Symbol.ShortName == row.ShortName {
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
func getTotalMarginUsed(backtest *models.Backtest) float64 {
	total := 0.00

	// Loop for expired postions
	for _, row := range backtest.TradeGroups {

		if row.Status == "Closed" {
			continue
		}

		total += row.Margin
	}

	// Return happy
	return total
}
