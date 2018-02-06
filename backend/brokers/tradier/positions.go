//
// Date: 2/6/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package tradier

import (
	"encoding/json"
	"fmt"

	"github.com/cloudmanic/app.options.cafe/backend/brokers/types"
	"github.com/tidwall/gjson"
)

//
// Get Positions
//
func (t *Api) GetPositions() ([]types.Position, error) {

	var positions []types.Position

	// Get the JSON
	jsonRt, err := t.SendGetRequest("/user/positions")

	if err != nil {
		return positions, err
	}

	// Process JSON returned
	return t.parsePositionsJson(jsonRt)

}

//
// Process JSON returned from an positions api call.
//
func (t *Api) parsePositionsJson(body string) ([]types.Position, error) {

	var positions []types.Position

	// Make sure we have at least one account (this should never happen)
	vo := gjson.Get(body, "accounts.account")

	if !vo.Exists() {
		// Return happy (just no positions found)
		return positions, nil
	}

	// Do we have only one account?
	vo = gjson.Get(body, "accounts.account.account_number")

	// Only one account
	if vo.Exists() {

		// if t.positionsParseOneAccount(body, &t_orders) != nil {
		// 	return positions, nil
		// }

	} else // More than one accounts
	{

		if t.positionsParseMoreThanOneAccount(body, &positions) != nil {
			return positions, nil
		}

	}

	// Convert to an formal order array
	//t.tempOrderArray2OrderArray(&t_orders, &orders)

	// Return happy
	return positions, nil

}

//
// Parse the case where the user has more than one account.
//
func (t *Api) positionsParseMoreThanOneAccount(body string, positions *[]types.Position) error {

	vo := gjson.Get(body, "accounts.account")

	type TempPosition struct {
		Id           int
		AccountId    string
		Symbol       string
		DateAcquired string  `json:"date_acquired"`
		CostBasis    float64 `json:"cost_basis"`
		Quantity     float64
	}

	// Loop through the different accounts.
	vo.ForEach(func(key, value gjson.Result) bool {

		// Do we have any positions?
		vo2 := gjson.Get(value.String(), "positions")

		if !vo2.Exists() {
			return true
		}

		// Set the account id
		account_number := gjson.Get(value.String(), "account_number").String()

		// Do we have more than one position
		vo2 = gjson.Get(value.String(), "positions.position.id")

		// More than one position??
		if !vo2.Exists() {

			var ws []TempPosition

			// Get just the orders part
			vo3 := gjson.Get(value.String(), "positions.position")

			// Unmarshal json
			if err := json.Unmarshal([]byte(vo3.String()), &ws); err != nil {
				return true
			}

			// Set the orders to our return
			for _, row := range ws {
				*positions = append(*positions, types.Position{
					Id:           row.Id,
					AccountId:    account_number,
					Symbol:       row.Symbol,
					DateAcquired: row.DateAcquired,
					CostBasis:    row.CostBasis,
					Quantity:     row.Quantity,
				})
			}

		} else {

			fmt.Println("Only one")

			var ws TempPosition

			// Get just the orders part
			vo3 := gjson.Get(value.String(), "positions.position")

			// Unmarshal json
			if err := json.Unmarshal([]byte(vo3.String()), &ws); err != nil {
				return true
			}

			// Set the orders to our return
			*positions = append(*positions, types.Position{
				Id:           ws.Id,
				AccountId:    account_number,
				Symbol:       ws.Symbol,
				DateAcquired: ws.DateAcquired,
				CostBasis:    ws.CostBasis,
				Quantity:     ws.Quantity,
			})

		}

		// keep iterating
		return true

	})

	return nil
}

// //
// // Parse the case where the user just has one account.
// //
// func (t *Api) positionsParseOneAccount(body string, t_orders *[]types.TradierOrder) error {

//   // Do we have any orders.
//   vo2 := gjson.Get(body, "accounts.account.orders")

//   // if !vo2.Exists() {
//   //   return errors.New("No Orders Found")
//   // }

//   // // Set the account id
//   // account_number := gjson.Get(body, "accounts.account.account_number").String()

//   // // Do we have more than one order
//   // vo2 = gjson.Get(body, "accounts.account.orders.order.id")

//   // // More than one order??
//   // if !vo2.Exists() {

//   //   var ws []types.TradierOrder

//   //   // Get just the orders part
//   //   vo3 := gjson.Get(body, "accounts.account.orders.order")

//   //   // Unmarshal json
//   //   if err := json.Unmarshal([]byte(vo3.String()), &ws); err != nil {
//   //     return err
//   //   }

//   //   // Set the orders to our return
//   //   for _, row := range ws {
//   //     row.AccountId = account_number
//   //     *t_orders = append(*t_orders, row)
//   //   }

//   // } else {
//   //   var ws types.TradierOrder

//   //   // Get just the orders part
//   //   vo3 := gjson.Get(body, "accounts.account.orders.order")

//   //   // Unmarshal json
//   //   if err := json.Unmarshal([]byte(vo3.String()), &ws); err != nil {
//   //     return err
//   //   }

//   //   // Set the orders we return.
//   //   ws.AccountId = account_number
//   //   *t_orders = append(*t_orders, ws)

//   // }

//   // Success
//   return nil
// }

/* End File */
