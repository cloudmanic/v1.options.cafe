//
// Date: 2018-07-17
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package screener

import (
	"errors"
	"os"

	"github.com/cloudmanic/app.options.cafe/backend/brokers/tradier"
	"github.com/cloudmanic/app.options.cafe/backend/brokers/types"
	"github.com/cloudmanic/app.options.cafe/backend/models"
)

var (
	broker tradier.Api
)

type Result struct {
	Credit      float64                  `json:"credit"`
	MidPoint    float64                  `json:"midpoint"`
	PrecentAway float64                  `json:"percent_away"`
	Legs        []types.OptionsChainItem `json:"legs"`
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
