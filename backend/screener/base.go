//
// Date: 2018-07-17
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package screener

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/brokers/tradier"
	"github.com/cloudmanic/app.options.cafe/backend/brokers/types"
	"github.com/cloudmanic/app.options.cafe/backend/library/cache"
	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/cloudmanic/app.options.cafe/backend/screener/put_credit_spread"
)

var (
	broker tradier.Api
)

type Result struct {
	Credit      float64         `json:"credit"`
	MidPoint    float64         `json:"midpoint"`
	PrecentAway float64         `json:"percent_away"`
	Legs        []models.Symbol `json:"legs"`
}

//
// Init.
//
func init() {

	// Setup the broker
	broker = tradier.Api{
		DB:     nil,
		ApiKey: os.Getenv("TRADIER_ADMIN_ACCESS_TOKEN"),
	}

}

//
// Loop through all screeners and store them in cache.
//
func PrimeAllScreenerCaches(db models.Datastore) {

	for {

		screeners := []models.Screener{}
		db.New().Find(&screeners)

		for _, row := range screeners {

			// Get screen with screen items
			screen, err := db.GetScreenerByIdAndUserId(row.Id, row.UserId)

			if err != nil {
				services.BetterError(err)
				continue
			}

			// Run screen
			switch row.Strategy {

			// Put Credit Spread
			case "put-credit-spread":

				result, err := put_credit_spread.RunPutCreditSpread(screen, db)

				if err != nil {
					services.BetterError(err)
					continue
				}

				// Store result in cache.
				cache.SetExpire("oc-screener-result-"+strconv.Itoa(int(row.Id)), (time.Minute * 5), result)

				break

			}

		}

		// Sleep for 5 mins.
		time.Sleep(time.Second * 60 * 5)

	}

}

//
// Get last quote for a symbol we are screening.
//
func GetQuote(smb string) (types.Quote, error) {

	// Make call to get a current quote.
	q, err := broker.GetQuotes([]string{smb})

	if err != nil {
		return types.Quote{}, err
	}

	if len(q) <= 0 {
		return types.Quote{}, errors.New("No quote found.")
	}

	// return quote
	return q[0], nil
}

//
// Search filter items for a particular value.
//
func FindFilterItemValue(item string, screen models.Screener) (models.ScreenerItem, error) {

	for _, row := range screen.Items {

		if row.Key == item {
			return row, nil
		}

	}

	// Not found
	return models.ScreenerItem{}, errors.New("Item not found.")
}

//
// Find a strike price that is X number of strikes below.
//
func FindByStrike(chain []types.OptionsChainItem, strike float64) (types.OptionsChainItem, error) {

	for _, row := range chain {

		if strike == row.Strike {
			return row, nil
		}

	}

	return types.OptionsChainItem{}, errors.New("No leg found.")
}

/* End File */
