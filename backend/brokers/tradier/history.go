package tradier

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/cloudmanic/app.options.cafe/backend/brokers/types"
	"github.com/tidwall/gjson"
)

//
// Get all history events for a broker.
//
func (t *Api) GetAllHistory() ([]types.History, error) {

	var history []types.History

	// Get the user profile
	user, err := t.GetUserProfile()

	if err != nil {
		return history, err
	}

	// Loop through the different accounts and get all orders.
	for _, row := range user.Accounts {

		tmp, _ := t.GetHistoryByAccountId(row.AccountNumber)

		for _, row2 := range tmp {
			history = append(history, row2)
		}

	}

	// Return happy
	return history, nil
}

//
// Get History by Account Id
//
func (t *Api) GetHistoryByAccountId(accountId string) ([]types.History, error) {

	var history []types.History

	// Get the JSON - Set big limit to get all data
	jsonRt, err := t.SendGetRequest("/accounts/" + accountId + "/history?limit=1000000")

	if err != nil {
		return history, err
	}

	// Make sure we have at least one event
	vo := gjson.Get(jsonRt, "history.event")

	if !vo.Exists() {
		return history, nil
	}

	// Do we have only one event or more?
	vo = gjson.Get(jsonRt, "history.event.date")

	// Only one event
	if vo.Exists() {

		// Add to balances array
		history = append(history, t.addJsonToEvent(gjson.Get(jsonRt, "history.event").String(), accountId))

	} else // More than one historical event
	{

		vo := gjson.Get(jsonRt, "history.event")

		// Loop through the different accounts.
		vo.ForEach(func(key, value gjson.Result) bool {

			// Add to balances array
			history = append(history, t.addJsonToEvent(value.String(), accountId))

			// keep iterating
			return true

		})

	}

	// Return happy
	return history, nil

}

//
// Take a json string and turn it into a history object.
//
func (t *Api) addJsonToEvent(eventJson string, accountId string) types.History {

	hasher := md5.New()
	hasher.Write([]byte(eventJson))

	// Get description
	eventType := gjson.Get(eventJson, "type").String()
	qty := gjson.Get(eventJson, "trade.quantity").Int()
	description := gjson.Get(eventJson, "trade.description").String()

	// See if this is a journal
	if gjson.Get(eventJson, "type").String() == "journal" {
		qty = 0
		description = gjson.Get(eventJson, "journal.description").String()
	}

	// See if this is a dividend
	if gjson.Get(eventJson, "type").String() == "dividend" {
		qty = gjson.Get(eventJson, "dividend.quantity").Int()
		description = gjson.Get(eventJson, "dividend.description").String()
	}

	// See if this is a interest
	if gjson.Get(eventJson, "type").String() == "interest" {
		qty = gjson.Get(eventJson, "interest.quantity").Int()
		description = gjson.Get(eventJson, "interest.description").String()
	}

	// See if this is a interest
	if gjson.Get(eventJson, "type").String() == "ach" {
		qty = gjson.Get(eventJson, "ach.quantity").Int()
		description = gjson.Get(eventJson, "ach.description").String()
	}

	// See if this is a DIVADJ
	if gjson.Get(eventJson, "type").String() == "DIVADJ" {
		eventType = "dividend"
		qty = gjson.Get(eventJson, "adjustment.quantity").Int()
		description = gjson.Get(eventJson, "adjustment.description").String()
	}

	// See if this is a position adjustments
	if gjson.Get(eventJson, "type").String() == "position adjustments" {
		qty = gjson.Get(eventJson, "adjustment.quantity").Int()
		description = gjson.Get(eventJson, "adjustment.description").String()
	}

	// See if this is a DIVPAY
	if gjson.Get(eventJson, "type").String() == "DIVPAY" {
		eventType = "dividend"
		qty = gjson.Get(eventJson, "adjustment.quantity").Int()
		description = gjson.Get(eventJson, "adjustment.description").String()
	}

	// See if this is a adjustment
	if gjson.Get(eventJson, "type").String() == "adjustment" {
		qty = gjson.Get(eventJson, "adjustment.quantity").Int()
		description = gjson.Get(eventJson, "adjustment.description").String()
	}

	// Return History event.
	return types.History{
		Id:          hex.EncodeToString(hasher.Sum(nil)),
		BrokerId:    accountId,
		Type:        eventType,
		Date:        gjson.Get(eventJson, "date").String(),
		Amount:      gjson.Get(eventJson, "amount").Float(),
		Symbol:      gjson.Get(eventJson, "trade.symbol").String(),
		Commission:  gjson.Get(eventJson, "trade.commission").Float(),
		Description: description,
		Price:       gjson.Get(eventJson, "trade.price").Float(),
		Quantity:    qty,
		TradeType:   gjson.Get(eventJson, "trade.trade_type").String(),
	}

}

/* End File */
