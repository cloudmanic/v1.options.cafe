package tradier

import (
	"encoding/json"
	"errors"
	"net/url"
	"strconv"

	"github.com/cloudmanic/app.options.cafe/backend/brokers/types"
	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/tidwall/gjson"
)

//
// Cancel order
//
func (t *Api) CancelOrder(accountId string, orderId string) error {

	// Log order cancel
	services.Info("Canceling order for account - " + accountId + " : " + orderId)

	// Get the JSON
	_, err := t.SendDeleteRequest("/accounts/" + accountId + "/orders/" + orderId)

	if err != nil {
		return err
	}

	return nil
}

//
// Submit order
//
func (t *Api) SubmitOrder(accountId string, order types.Order) (types.OrderSubmit, error) {

	// Log order
	orderJson, _ := json.Marshal(order)
	services.Info("Placing order for account - " + accountId + " : " + string(orderJson))

	// Prep Order
	params := prepOrder(order)

	// Get the JSON
	jsonRt, err := t.SendPostRequest("/accounts/"+accountId+"/orders", params)

	if err != nil {
		return types.OrderSubmit{}, err
	}

	// Check for errors
	errMsg := gjson.Get(jsonRt, "errors.error").String()

	if len(errMsg) > 0 {
		return types.OrderSubmit{Status: "Error", Error: errMsg}, nil
	}

	// Grab result and convert to strut
	result := types.OrderSubmit{}
	vo := gjson.Get(jsonRt, "order").String()
	err2 := json.Unmarshal([]byte(vo), &result)

	if err2 != nil {
		return types.OrderSubmit{}, err2
	}

	return result, nil
}

//
// Preview order
//
func (t *Api) PreviewOrder(accountId string, order types.Order) (types.OrderPreview, error) {

	// Prep Order
	params := prepOrder(order)

	// Set to Preview
	params.Set("preview", "true")

	// Get the JSON
	jsonRt, err := t.SendPostRequest("/accounts/"+accountId+"/orders", params)

	if err != nil {
		return types.OrderPreview{}, err
	}

	// Check for errors
	errMsg := gjson.Get(jsonRt, "errors.error").String()

	if len(errMsg) > 0 {
		return types.OrderPreview{Status: "Error", Error: errMsg}, nil
	}

	// Grab result and convert to strut
	result := types.OrderPreview{}
	vo := gjson.Get(jsonRt, "order").String()
	err2 := json.Unmarshal([]byte(vo), &result)

	if err2 != nil {
		return types.OrderPreview{}, err2
	}

	return result, nil
}

//
// Get All Orders Ever
//
func (t *Api) GetAllOrders() ([]types.Order, error) {

	page := 1
	orders := []types.Order{}

	// Get the user profile
	user, err := t.GetUserProfile()

	if err != nil {
		return orders, err
	}

	// Loop through the different accounts and get all orders.
	for _, row := range user.Accounts {

		page = 1

		// Here we loop through the orders getting the orders 1000 orders at a time.
		for {

			var tmp []types.Order

			err := t.processOrderDataForGetAllOrders(row.AccountNumber, &tmp, page, 1000)

			if err != nil {
				services.BetterError(err)
			}

			// Add orders into master var.
			for _, row := range tmp {
				orders = append(orders, row)
			}

			// If we got no results break
			if len(tmp) <= 0 {
				break
			}

			// Increase page.
			page++

		}

	}

	// Return the data.
	return orders, nil

}

//
// Get Orders
//
func (t *Api) GetOrders() ([]types.Order, error) {

	var orders []types.Order

	// Get the JSON
	jsonRt, err := t.SendGetRequest("/user/orders")

	if err != nil {
		return orders, err
	}

	// Process JSON returned
	return t.parseOrdersJson(jsonRt)

}

//
// Process JSON returned from an ORDERS api call.
//
func (t *Api) parseOrdersJson(body string) ([]types.Order, error) {

	var orders []types.Order
	var t_orders []types.TradierOrder

	// Make sure we have at least one account (this should never happen)
	vo := gjson.Get(body, "accounts.account")

	if !vo.Exists() {
		// Return happy (just no orders found)
		return orders, nil
	}

	// Do we have only one account?
	vo = gjson.Get(body, "accounts.account.account_number")

	// Only one account
	if vo.Exists() {

		if t.orderParseOneAccount(body, &t_orders) != nil {
			return orders, nil
		}

	} else // More than one accounts
	{

		if t.orderParseMoreThanOneAccount(body, &t_orders) != nil {
			return orders, nil
		}

	}

	// Convert to an formal order array
	t.tempOrderArray2OrderArray(&t_orders, &orders)

	// Return happy
	return orders, nil

}

//
// Parse the case where the user just has more than one account.
// This function just sets the t_orders
//
func (t *Api) orderParseMoreThanOneAccount(body string, t_orders *[]types.TradierOrder) error {

	vo := gjson.Get(body, "accounts.account")

	// Loop through the different accounts.
	vo.ForEach(func(key, value gjson.Result) bool {

		// Do we have any orders.
		vo2 := gjson.Get(value.String(), "orders")

		if !vo2.Exists() {
			return true
		}

		// Set the account id
		account_number := gjson.Get(value.String(), "account_number").String()

		// Do we have more than one order
		vo2 = gjson.Get(value.String(), "orders.order.id")

		// More than one order??
		if !vo2.Exists() {

			var ws []types.TradierOrder

			// Get just the orders part
			vo3 := gjson.Get(value.String(), "orders.order")

			// Unmarshal json
			if err := json.Unmarshal([]byte(vo3.String()), &ws); err != nil {
				return true
			}

			// Set the orders to our return
			for _, row := range ws {
				row.AccountId = account_number
				*t_orders = append(*t_orders, row)
			}

		} else {
			var ws types.TradierOrder

			// Get just the orders part
			vo3 := gjson.Get(value.String(), "orders.order")

			// Unmarshal json
			if err := json.Unmarshal([]byte(vo3.String()), &ws); err != nil {
				return true
			}

			// Set the orders we return.
			ws.AccountId = account_number
			*t_orders = append(*t_orders, ws)

		}

		// keep iterating
		return true

	})

	// Happy
	return nil
}

//
// Parse the case where the user just has one account.
// This function just sets the t_orders
//
func (t *Api) orderParseOneAccount(body string, t_orders *[]types.TradierOrder) error {

	// Do we have any orders.
	vo2 := gjson.Get(body, "accounts.account.orders")

	if !vo2.Exists() {
		return errors.New("No Orders Found")
	}

	// Set the account id
	account_number := gjson.Get(body, "accounts.account.account_number").String()

	// Do we have more than one order
	vo2 = gjson.Get(body, "accounts.account.orders.order.id")

	// More than one order??
	if !vo2.Exists() {

		var ws []types.TradierOrder

		// Get just the orders part
		vo3 := gjson.Get(body, "accounts.account.orders.order")

		// Unmarshal json
		if err := json.Unmarshal([]byte(vo3.String()), &ws); err != nil {
			return err
		}

		// Set the orders to our return
		for _, row := range ws {
			row.AccountId = account_number
			*t_orders = append(*t_orders, row)
		}

	} else {
		var ws types.TradierOrder

		// Get just the orders part
		vo3 := gjson.Get(body, "accounts.account.orders.order")

		// Unmarshal json
		if err := json.Unmarshal([]byte(vo3.String()), &ws); err != nil {
			return err
		}

		// Set the orders we return.
		ws.AccountId = account_number
		*t_orders = append(*t_orders, ws)

	}

	// Success
	return nil
}

//
// Process order data for GetAllOrders
//
func (t *Api) processOrderDataForGetAllOrders(accountNumber string, orders *[]types.Order, page int, limit int) error {

	var t_orders []types.TradierOrder

	// Get the JSON
	jsonRt, err := t.SendGetRequest("/accounts/" + accountNumber + "/orders?filter=all&limit=" + strconv.Itoa(limit) + "&page=" + strconv.Itoa(page))

	if err != nil {
		return err
	}

	// See if we got zero results.
	if jsonRt == `{"orders":"null"}` {
		return nil
	}

	// Do we have more than one order
	vo1 := gjson.Get(jsonRt, "orders.order.Id")

	// More than one order??
	if !vo1.Exists() {

		var ws []types.TradierOrder

		// Get just the orders part
		vo2 := gjson.Get(jsonRt, "orders.order")

		// Unmarshal json
		if err := json.Unmarshal([]byte(vo2.String()), &ws); err != nil {
			return err
		}

		// Set the orders to our return
		for _, row := range ws {
			row.AccountId = accountNumber
			t_orders = append(t_orders, row)
		}

	} else {
		var ws types.TradierOrder

		// Get just the orders part
		vo2 := gjson.Get(jsonRt, "orders.order")

		// Unmarshal json
		if err := json.Unmarshal([]byte(vo2.String()), &ws); err != nil {
			return err
		}

		// Set the orders we return.
		ws.AccountId = accountNumber
		t_orders = append(t_orders, ws)

	}

	// Convert to an formal order array
	t.tempOrderArray2OrderArray(&t_orders, orders)

	// Return the data.
	return nil

}

//
// Go from a temp order array to a formal order array.
//
func (t *Api) tempOrderArray2OrderArray(t_orders *[]types.TradierOrder, orders *[]types.Order) error {

	// Clean up the array and make it more generic
	for _, row := range *t_orders {

		var legs []types.OrderLeg

		// Merge in the legs
		if len(row.Legs) > 0 {

			for _, row2 := range row.Legs {

				legs = append(legs, types.OrderLeg{
					Type:              row2.Type,
					Symbol:            row2.Symbol,
					OptionSymbol:      row2.OptionSymbol,
					Side:              row2.Side,
					Quantity:          row2.Quantity,
					Status:            row2.Status,
					Duration:          row2.Duration,
					AvgFillPrice:      row2.AvgFillPrice,
					ExecQuantity:      row2.ExecQuantity,
					LastFillPrice:     row2.LastFillPrice,
					LastFillQuantity:  row2.LastFillQuantity,
					RemainingQuantity: row2.RemainingQuantity,
					CreateDate:        row2.CreateDate,
					TransactionDate:   row2.TransactionDate,
				})

			}

		}

		// Append the orders
		*orders = append(*orders, types.Order{
			Id:                strconv.Itoa(row.Id),
			AccountId:         row.AccountId,
			Type:              row.Type,
			Symbol:            row.Symbol,
			Side:              row.Side,
			Quantity:          row.Quantity,
			Status:            row.Status,
			Duration:          row.Duration,
			Price:             row.Price,
			AvgFillPrice:      row.AvgFillPrice,
			ExecQuantity:      row.ExecQuantity,
			LastFillPrice:     row.LastFillPrice,
			LastFillQuantity:  row.LastFillQuantity,
			RemainingQuantity: row.RemainingQuantity,
			CreateDate:        row.CreateDate,
			TransactionDate:   row.TransactionDate,
			Class:             row.Class,
			OptionSymbol:      row.OptionSymbol,
			NumLegs:           row.NumLegs,
			Legs:              legs,
		})

	}

	// Success
	return nil

}

// ----------------- Helper Functions -------------------- //

//
// Prep order
//
func prepOrder(order types.Order) url.Values {

	// Build Form
	params := url.Values{}
	params.Set("price", strconv.FormatFloat(order.Price, 'f', 2, 64))
	params.Set("type", order.Type)
	params.Set("symbol", order.Symbol)
	params.Set("duration", order.Duration)
	params.Set("class", order.Class)
	params.Set("side", order.Side)
	params.Set("stop", strconv.FormatFloat(order.Stop, 'f', 2, 64))

	if order.Class != "multileg" {
		params.Set("quantity", strconv.FormatFloat(order.Quantity, 'f', 2, 64))
	}

	// Multi Leg?
	for key, row := range order.Legs {
		params.Set("side["+strconv.Itoa(key)+"]", row.Side)
		params.Set("quantity["+strconv.Itoa(key)+"]", strconv.FormatFloat(row.Quantity, 'f', 2, 64))
		params.Set("option_symbol["+strconv.Itoa(key)+"]", row.OptionSymbol)
	}

	// Return Happy
	return params
}

/* End File */
