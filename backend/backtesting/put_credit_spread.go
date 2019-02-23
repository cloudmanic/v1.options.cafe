//
// Date: 2019-02-22
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package backtesting

import (
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/brokers/eod"
	"github.com/cloudmanic/app.options.cafe/backend/brokers/types"
	"github.com/cloudmanic/app.options.cafe/backend/library/helpers"
	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/cloudmanic/app.options.cafe/backend/screener"
)

//
// DoPutCreditSpread - Run a put credit spread backtest.
//
func (t *Base) DoPutCreditSpread(today time.Time, backtest *models.Backtest, underlyingLast float64, chains map[time.Time]types.OptionsChain) error {

	// See if we have any positions to close
	t.CloseMultiLegCredit(today, backtest)

	// Results that we return.
	results := []screener.Result{}

	// Set up a new screener so we can use it's Functions
	screenObj := screener.NewScreen(t.DB, &eod.Api{})

	// Set params
	spreadWidth := screenObj.GetPutCreditSpreadParms(backtest.Screen, underlyingLast)

	// Loop through the expire dates
	for _, row := range chains {

		for _, row2 := range row.Puts {

			// No need to pay attention to open interest of zero
			if row2.OpenInterest == 0 {
				continue
			}

			// Skip strikes that are higher than our min strike. Based on percent away.
			if !screenObj.FilterStrikeByPercentDown("short-strike-percent-away", backtest.Screen, row2.Strike, underlyingLast) {
				continue
			}

			// Find the strike that is x points away.
			buyLeg, err := screenObj.FindByStrike(row.Puts, (row2.Strike - spreadWidth))

			if err != nil {
				continue
			}

			// See if there is enough credit
			credit := row2.Bid - buyLeg.Ask

			if !screenObj.FilterOpenCredit(backtest.Screen, credit) {
				continue
			}

			// Figure out the credit spread amount.
			buyCost := row2.Ask - buyLeg.Bid
			midPoint := (credit + buyCost) / 2

			// Add in Symbol Object - Sell leg
			symbSellLeg, err := t.DB.CreateNewSymbol(row2.Symbol, row2.Description, "Option")

			if err != nil {
				return err
			}

			// Add in Symbol Object - Buy leg
			symbBuyLeg, err := t.DB.CreateNewSymbol(buyLeg.Symbol, buyLeg.Description, "Option")

			if err != nil {
				return err
			}

			// We have a winner
			results = append(results, screener.Result{
				Credit:         helpers.Round(credit, 2),
				MidPoint:       helpers.Round(midPoint, 2),
				PutPrecentAway: helpers.Round(((1 - row2.Strike/underlyingLast) * 100), 2),
				Legs:           []models.Symbol{symbBuyLeg, symbSellLeg},
			})
		}

	}

	// TODO(spicer): Figure which result to open
	if len(results) > 0 {
		t.OpenMultiLegCredit(today, backtest, results[0])
	}

	return nil
}

/* End File */
