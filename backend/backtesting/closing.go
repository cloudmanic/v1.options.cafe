//
// Date: 2019-02-22
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package backtesting

import (
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/optionscafe/options-cafe-cli/helpers"
)

//
// CloseMultiLegCredit - Close positions
//
func (t *Base) CloseMultiLegCredit(today time.Time, underlyingLast float64, backtest *models.Backtest) {
	// Expire positions
	t.expirePositions(today, backtest)

	// Close if we touch the short leg
	t.closeOnShortTouch(today, underlyingLast, backtest)
}

//
// closeOnShortTouch - If our trade touches the short leg we close
//
func (t *Base) closeOnShortTouch(today time.Time, underlyingLast float64, backtest *models.Backtest) {

	// Loop for expired postions
	for key, row := range backtest.Positions {

		if row.Status == "Closed" {
			continue
		}

		// TODO(Spicer): Currently this only works for PCS. We assume the second leg is the short leg.
		if underlyingLast <= row.Legs[1].OptionStrike {
			backtest.Positions[key].Status = "Closed"
			backtest.Positions[key].ClosePrice = 0.00
			backtest.Positions[key].CloseDate = models.Date{today}
			backtest.Positions[key].Note = "Trade touched the short leg."
		}

		// // See if any of the legs are expired
		// for _, row2 := range row.Legs {
		// 	if today.Format("2006-01-02") == row2.OptionExpire.Format("2006-01-02") ||
		// 		today.After(helpers.ParseDateNoError(row2.OptionExpire.Format("2006-01-02"))) {
		// 		expired = true
		// 		break
		// 	}
		// }

		// // If expired close out trade
		// if expired && row.Status == "Open" {
		// 	backtest.Positions[key].Status = "Closed"
		// 	backtest.Positions[key].ClosePrice = 0.00
		// 	backtest.Positions[key].CloseDate = models.Date{today}
		// 	backtest.Positions[key].Note = "Expired worthless."
		// }
	}

}

//
// expirePositions - Loop through to see if we have any positions to expire.
//
func (t *Base) expirePositions(today time.Time, backtest *models.Backtest) {

	// Loop for expired postions
	for key, row := range backtest.Positions {

		if row.Status == "Closed" {
			continue
		}

		expired := false

		// See if any of the legs are expired
		for _, row2 := range row.Legs {
			if today.Format("2006-01-02") == row2.OptionExpire.Format("2006-01-02") ||
				today.After(helpers.ParseDateNoError(row2.OptionExpire.Format("2006-01-02"))) {
				expired = true
				break
			}
		}

		// If expired close out trade
		if expired && row.Status == "Open" {
			backtest.Positions[key].Status = "Closed"
			backtest.Positions[key].ClosePrice = 0.00
			backtest.Positions[key].CloseDate = models.Date{today}
			backtest.Positions[key].Note = "Expired worthless."
		}
	}

}

/* End File */
