package feed

import (
	"encoding/json"
	"fmt"
	"sort"
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

//
// Return active symbols. This is handy because we
// sort and filter before returning.
//
func (t *Base) GetActiveSymbols() []string {

	var activeSymbols []string

	// Symbols we always want
	activeSymbols = append(activeSymbols, "$DJI")
	activeSymbols = append(activeSymbols, "SPX")
	activeSymbols = append(activeSymbols, "COMP")
	activeSymbols = append(activeSymbols, "VIX")

	// Get the watch lists for this user.
	watchList, err := t.DB.GetWatchlistsByUserId(t.User.Id)

	// Loop through the watchlists and add the symbols
	if err == nil {

		for _, row := range watchList {

			for _, row2 := range row.Symbols {

				activeSymbols = append(activeSymbols, row2.Symbol.ShortName)

			}

		}

	}

	// Add in the orders we want.
	t.muOrders.Lock()

	for _, row := range t.Orders {

		// Stock symbol
		activeSymbols = append(activeSymbols, row.Symbol)

		// Multi leg trade
		if row.NumLegs > 0 {

			for _, row2 := range row.Legs {

				activeSymbols = append(activeSymbols, row2.OptionSymbol)

			}

		}

		// Single option order
		if len(row.OptionSymbol) > 0 {
			activeSymbols = append(activeSymbols, row.OptionSymbol)
		}

	}

	t.muOrders.Unlock()

	// Clean up the list.
	activeSymbols = t.ToUpperStrings(activeSymbols)
	activeSymbols = t.RemoveDupsStrings(activeSymbols)

	// Sort the list.
	sort.Strings(activeSymbols)

	// Return the cleaned up list.
	return activeSymbols

}

// ----------------- Market Status ------------------- //

//
// Do get market status
//
func (t *Base) GetMarketStatus() error {

	// Make API call
	marketStatus, err := t.Api.GetMarketStatus()

	if err != nil {
		return err
	}

	// Save the market status in the fetch object
	t.muMarketStatus.Lock()
	t.MarketStatus = marketStatus
	t.muMarketStatus.Unlock()

	// Send up websocket.
	err = t.WriteDataChannel("market/status", marketStatus)

	if err != nil {
		return fmt.Errorf("GetMarketStatus() WriteDataChannel : ", err)
	}

	// Return Happy
	return nil

}

// ----------------- Positions ------------------- //

//
// Do get positions
//
func (t *Base) GetPositions() error {

	// Make API call
	positions, err := t.Api.GetPositions()

	if err != nil {
		return err
	}

	// Save the positions in the fetch object
	t.muPositions.Lock()
	t.Positions = positions
	t.muPositions.Unlock()

	// Send up websocket.
	err = t.WriteDataChannel("positions", positions)

	if err != nil {
		return fmt.Errorf("Fetch.GetPositions() : ", err)
	}

	// Return Happy
	return nil

}

// ----------------- Quotes ------------------- //

//
// Do get quotes - more details from the streaming - activeSymbols
//
func (t *Base) GetActiveSymbolsDetailedQuotes() error {

	symbols := t.GetActiveSymbols()
	detailedQuotes, err := t.Api.GetQuotes(symbols)

	if err != nil {
		return err
	}

	// Loop through the quotes sending them up the websocket channel
	for _, row := range detailedQuotes {

		// Convert to a json string.
		data_json, err := json.Marshal(row)

		if err != nil {
			return err
		}

		// Send data up websocket.
		send_json, err := t.GetSendJson("quote", string(data_json))

		if err != nil {
			return err
		}

		// Send up websocket.
		err = t.WriteQuoteChannel(send_json)

		if err != nil {
			return fmt.Errorf("GetActiveSymbolsDetailedQuotes() WriteDataChannel : ", err)
		}

	}

	// Return happy
	return nil

}

// ----------------- Balances ------------------- //

//
// Do get Balances
//
func (t *Base) GetBalances() error {

	balances, err := t.Api.GetBalances()

	if err != nil {
		return err
	}

	// Save the balances in the fetch object
	t.muBalances.Lock()
	t.Balances = balances
	t.muBalances.Unlock()

	// Send up websocket.
	err = t.WriteDataChannel("balances", balances)

	if err != nil {
		return fmt.Errorf("GetBalances() WriteDataChannel : ", err)
	}

	// Return Happy
	return nil

}

// ----------------- Access Tokens ------------------- //

//
// Do update access token from refresh
//
func (t *Base) AccessTokenRefresh() error {

	err := t.Api.DoRefreshAccessTokenIfNeeded(t.User)

	if err != nil {
		return err
	}

	// Return Happy
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
