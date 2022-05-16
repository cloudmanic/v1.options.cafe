//
// Date: 2022-05-15
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Copyright: 2022 Cloudmanic Labs, LLC. All rights reserved.
//

package helpers

import (
	"fmt"
	"time"

	"app.options.cafe/models"
	"app.options.cafe/screener"
)

//
// OpenTrade will open a trade
// qtyMap maps each leg with the qty we need to purchase for that leg. Negative value is a short
//
func OpenTrade(today time.Time, backtest *models.Backtest, result screener.Result) {
	// Lots. We just assume it is the lowest qty of the legs passed in.
	lots := 0

	// Build legs
	legs := []models.BacktestPosition{}

	for _, row := range result.Legs {
		legs = append(legs, models.BacktestPosition{
			UserId:   backtest.UserId,
			Status:   "Open",
			SymbolId: row.Id,
			Symbol:   row,
			OpenDate: today,
			Qty:      row.Qty,
			OrgQty:   row.Qty,
		})

		// Figure out lots of this trade
		if lots == 0 {
			lots = row.Qty
		} else if lots > row.Qty {
			lots = row.Qty
		}
	}

	// Figure out open price.
	openPrice := result.Bid * 100 * float64(lots)

	// we can configure to use midpoint.
	if backtest.Midpoint {
		openPrice = result.MidPoint * 100 * float64(lots)
	}

	// Spread text
	spreadText := ""

	switch len(result.Legs) {
	case 1:
		spreadText = fmt.Sprintf("%s %s $%.2f", legs[0].Symbol.OptionUnderlying, legs[0].Symbol.OptionExpire.Format("01/02/2006"), legs[0].Symbol.OptionStrike)
		break

	case 2:
		spreadText = fmt.Sprintf("%s %s $%.2f / $%.2f", legs[0].Symbol.OptionUnderlying, legs[0].Symbol.OptionExpire.Format("01/02/2006"), legs[0].Symbol.OptionStrike, legs[1].Symbol.OptionStrike)
		break

	case 3:
		spreadText = fmt.Sprintf("%s %s $%.2f / $%.2f / $%.2f", legs[0].Symbol.OptionUnderlying, legs[0].Symbol.OptionExpire.Format("01/02/2006"), legs[0].Symbol.OptionStrike, legs[1].Symbol.OptionStrike, legs[2].Symbol.OptionStrike)
		break

	case 4:
		spreadText = fmt.Sprintf("%s %s $%.2f / $%.2f / $%.2f / $%.2f", legs[0].Symbol.OptionUnderlying, legs[0].Symbol.OptionExpire.Format("01/02/2006"), legs[0].Symbol.OptionStrike, legs[1].Symbol.OptionStrike, legs[2].Symbol.OptionStrike, legs[3].Symbol.OptionStrike)
		break
	}

	// Figure out margin
	margin := 0.00
	if openPrice < 0 {
		margin = 0.00 // TODO(spicer): Do the math and make this work.
	}

	// Add position
	backtest.TradeGroups = append(backtest.TradeGroups, models.BacktestTradeGroup{
		UserId:          backtest.UserId,
		Strategy:        backtest.Screen.Strategy,
		Status:          "Open",
		SpreadText:      spreadText,
		OpenDate:        models.Date{today},
		OpenPrice:       openPrice,
		Margin:          margin,
		Positions:       legs,
		Lots:            lots,
		Credit:          (openPrice / float64(lots)) / 100,
		PutPrecentAway:  result.PutPrecentAway,
		CallPrecentAway: result.CallPrecentAway,
		Balance:         (backtest.EndingBalance - openPrice),
	})

	// Update ending balance
	backtest.EndingBalance = backtest.EndingBalance - openPrice
}
