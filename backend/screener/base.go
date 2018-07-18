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
)

var (
	broker tradier.Api
)

type Filter struct {
	Symbol      string
	Strategy    string // put-credit-spread
	FilterItems []FilterItems
}

type FilterItems struct {
	Key         string
	Operator    string
	ValueNumber float64
	ValueString string
}

type Result struct {
	Credit      float64
	MidPoint    float64
	PrecentAway float64
	Legs        []types.OptionsChainItem
	//Wing1SellLegStrike string
	//   'sell_leg' => $row2['strike'],
	//   'buy_leg' => $buy_leg['strike'],
	//   'expire' => $row,
	//   'expire_df1' => date('n/j/y', strtotime($row)),
	//   'credit' => number_format($credit, 2),
	//   'midpoint' => number_format($mid_point, 2),
	//   'precent_away' => number_format((1 - $row2['strike'] / $stock['last']) * 100, 2),
	//   'occ_sell' => $row2['symbol'],
	//   'occ_buy' => $buy_leg['symbol']
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
func FindFilterItemValue(item string, filter Filter) (FilterItems, error) {

	for _, row := range filter.FilterItems {

		if row.Key == item {
			return row, nil
		}

	}

	// Not found
	return FilterItems{}, errors.New("Item not found.")
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
