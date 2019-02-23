//
// Date: 2019-02-22
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package backtesting

import (
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/brokers/types"
	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/optionscafe/options-cafe-cli/helpers"
)

//
// CloseMultiLegCredit - Close positions
//
func (t *Base) CloseMultiLegCredit(today time.Time, underlyingLast float64, backtest *models.Backtest, chains map[time.Time]types.OptionsChain) {
	// Expire positions
	t.expirePositions(today, backtest)

	// Close if we touch the short leg
	t.closeOnShortTouch(today, underlyingLast, backtest, chains)
}

//
// closeOnShortTouch - If our trade touches the short leg we close
//
func (t *Base) closeOnShortTouch(today time.Time, underlyingLast float64, backtest *models.Backtest, chains map[time.Time]types.OptionsChain) {

	// TODO(spicer): make this work from configs
	lots := 10

	// Loop for expired postions
	for key, row := range backtest.Positions {

		if row.Status == "Closed" {
			continue
		}

		// TODO(Spicer): Currently this only works for PCS. We assume the second leg is the short leg.
		if underlyingLast <= row.Legs[1].OptionStrike {
			ed := helpers.ParseDateNoError(row.Legs[0].OptionExpire.Format("2006-01-02"))

			var leg1Chain types.OptionsChainItem
			var leg2Chain types.OptionsChainItem

			// Loop through until we find first symbol
			for _, row2 := range chains[ed].Puts {
				if row2.Symbol != row.Legs[0].ShortName {
					continue
				}

				// We found it
				leg1Chain = row2
				break
			}

			// Loop through until we find second symbol
			for _, row2 := range chains[ed].Puts {
				if row2.Symbol != row.Legs[1].ShortName {
					continue
				}

				// We found it
				leg2Chain = row2
				break
			}

			// Set closing Closing Price
			closingPrice := ((leg2Chain.Ask * 100.00 * float64(lots)) - (leg1Chain.Bid * 100.00 * float64(lots))) * -1

			// Close trade
			backtest.Positions[key].Status = "Closed"
			backtest.Positions[key].ClosePrice = closingPrice
			backtest.Positions[key].CloseDate = models.Date{today}
			backtest.Positions[key].Note = "Trade touched the short leg."
			backtest.EndingBalance += closingPrice
		}
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
