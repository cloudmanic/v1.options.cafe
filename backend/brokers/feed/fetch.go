package feed

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cloudmanic/app.options.cafe/backend/controllers"
)

//
// Send cached data up the websocket pipeline.
// This is useful for when the application wants
// data between pooling of the broker.
//
func (t *Base) RefreshFromCached() error {

	// UserProfile - Send up websocket.
	err := t.WriteDataChannel("user/profile", t.UserProfile)

	if err != nil {
		return fmt.Errorf("RefreshFromCached() WriteDataChannel - user/profile : ", err)
	}

	// MarketStatus - Send up websocket.
	err = t.WriteDataChannel("market/status", t.MarketStatus)

	if err != nil {
		return fmt.Errorf("RefreshFromCached() WriteDataChannel - market/status : ", err)
	}

	// Orders - Send up websocket.
	err = t.WriteDataChannel("orders", t.Orders)

	if err != nil {
		return fmt.Errorf("RefreshFromCached() WriteDataChannel - orders : ", err)
	}

	// Balances - Send up websocket.
	err = t.WriteDataChannel("balances", t.Balances)

	if err != nil {
		return fmt.Errorf("RefreshFromCached() WriteDataChannel - balances : ", err)
	}

	// No error
	return nil
}

// ----------------- Helper Functions ---------------- //

//
// Return a json object ready to be sent up the websocket
//
func (t *Base) GetSendJson(uri string, data_json string) (string, error) {

	// Send Object
	send := SendStruct{
		Uri:  uri,
		Body: data_json,
	}

	send_json, err := json.Marshal(send)

	if err != nil {
		return "", err
	}

	return string(send_json), nil
}

//
// Send data up websocket.
//
func (t *Base) WriteDataChannel(send_type string, sendObject interface{}) error {

	// Convert to a json string.
	dataJson, err := json.Marshal(sendObject)

	if err != nil {
		return fmt.Errorf("WriteDataChannel() json.Marshal : ", err)
	}

	// Send data up websocket.
	sendJson, err := t.GetSendJson(send_type, string(dataJson))

	if err != nil {
		return fmt.Errorf("WriteDataChannel() GetSendJson Send Object : ", err)
	}

	// Write data out websocket
	t.DataChan <- controllers.SendStruct{UserId: t.User.Id, Body: sendJson}

	// Return happy.
	return nil
}

//
// Send data up quote websocket.
//
func (t *Base) WriteQuoteChannel(sendJson string) error {

	// Write data out websocket
	t.QuoteChan <- controllers.SendStruct{UserId: t.User.Id, Body: sendJson}

	// Return happy.
	return nil
}

//
// Remove duplicates from an array of strings.
//
func (t *Base) RemoveDupsStrings(list []string) []string {

	u := make([]string, 0, len(list))
	m := make(map[string]bool)

	for _, row := range list {

		if _, ok := m[row]; !ok {
			m[row] = true
			u = append(u, row)
		}

	}

	return u
}

//
// Take an array of strings and make them all upper case.
//
func (t *Base) ToUpperStrings(vs []string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = strings.ToUpper(v)
	}
	return vsm
}

/* End File */
