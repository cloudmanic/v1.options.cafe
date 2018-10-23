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

				result, err := RunPutCreditSpread(screen, db)

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
func FindFilterItemsByKey(item string, screen models.Screener) []models.ScreenerItem {

	rt := []models.ScreenerItem{}

	for _, row := range screen.Items {

		if row.Key == item {
			rt = append(rt, row)
		}

	}

	// Return finds
	return rt
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

// ---------------- Shared Filters --------------------- //

//
// Review trade for open credit
//
func FilterOpenCredit(screen models.Screener, credit float64) bool {

	// Find keys related to this filter.
	items := FindFilterItemsByKey("open-credit", screen)

	// Loop through the keys. Return false if something does not match up.
	for _, row := range items {

		// Switch based on the operator
		switch row.Operator {

		// Is the credit > than this value.
		case "<":
			if credit > row.ValueNumber {
				return false
			}
			break

		// Is the credit < this value
		case ">":
			if credit < row.ValueNumber {
				return false
			}
			break

		// Is the credit = this value
		case "=":
			if credit != row.ValueNumber {
				return false
			}
			break

		}

	}

	// If we made it this far it is true.
	return true
}

//
// Review trade for days to expire.
//
func FilterDaysToExpireDaysToExpire(screen models.Screener, expire time.Time) bool {

	// Get days to expire.
	today := time.Now()
	daysToExpire := int(today.Sub(expire).Hours()/24) * -1

	// Find keys related to this filter.
	items := FindFilterItemsByKey("days-to-expire", screen)

	// Loop through the keys. Return false if something does not match up.
	for _, row := range items {

		// Switch based on the operator
		switch row.Operator {

		// Is this spread to far out?
		case "<":
			if float64(daysToExpire) > row.ValueNumber {
				return false
			}
			break

		// Is this spread to far in?
		case ">":
			if float64(daysToExpire) < row.ValueNumber {
				return false
			}
			break

		// Days to expire must be equal to the value we pass in.
		case "=":
			if float64(daysToExpire) != row.ValueNumber {
				return false
			}
			break

		}

	}

	// If we made it this far it is true.
	return true
}

/* End File */
