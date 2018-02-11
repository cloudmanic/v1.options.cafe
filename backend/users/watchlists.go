//
// Date: 9/30/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package users

import (
	"encoding/json"

	"github.com/cloudmanic/app.options.cafe/backend/controllers"
	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/cloudmanic/app.options.cafe/backend/models"
)

//
// Send a user's watchlist up the websocket channel
//
func (t *Base) WsSendWatchlists(user *UserFeed, request controllers.ReceivedStruct) {

	// Get the watchlists
	wLists, err := t.DB.GetWatchlistsByUserId(user.Profile.Id)

	if err != nil {
		return
	}

	// Loop through the different watchlists
	for _, row := range wLists {

		// Convert to a json string.
		dataJson, err := json.Marshal(row)

		if err != nil {
			services.Error(err, "WsSendWatchlists() json.Marshal (#1)")
			continue
		}

		// Build JSON we send
		jsonSend, err := t.WsSendJsonBuild("watchlists", dataJson)

		if err != nil {
			services.Error(err, "WsSendWatchlists() WsSendJsonBuild (#2)")
			continue
		}

		// Send up the websocket
		user.DataChan <- controllers.SendStruct{UserId: user.Profile.Id, Body: string(jsonSend)}

		// Example to send back to just the current connection. (this is useful for validation)
		//request.Connection.WriteChan <- jsonSend

	}

	// Return happy
	return
}

//
// Verify we have default watchlist in place.
//
func (t *Base) VerifyDefaultWatchList(user models.User) {

	// Setup defaults.
	type Y struct {
		SymShort string
		SymLong  string
	}

	var m []Y
	m = append(m, Y{SymShort: "SPY", SymLong: "SPDR S&P 500"})
	m = append(m, Y{SymShort: "IWM", SymLong: "Ishares Russell 2000 Etf"})
	m = append(m, Y{SymShort: "MCD", SymLong: "McDonald's Corp"})
	m = append(m, Y{SymShort: "XLF", SymLong: "SPDR Select Sector Fund - Financial"})
	m = append(m, Y{SymShort: "AMZN", SymLong: "Amazon.com Inc"})
	m = append(m, Y{SymShort: "AAPL", SymLong: "Apple Inc."})
	m = append(m, Y{SymShort: "SBUX", SymLong: "Starbucks Corp"})
	m = append(m, Y{SymShort: "BAC", SymLong: "Bank Of America Corporation"})
	m = append(m, Y{SymShort: "HD", SymLong: "The Home Depot Inc"})
	m = append(m, Y{SymShort: "CAT", SymLong: "Caterpillar Inc"})

	// See if this user already had a watchlist
	_, err := t.DB.GetWatchlistsByUserId(user.Id)

	// If no watchlists we create a default one with some default symbols.
	if err != nil {

		wList, err := t.DB.CreateNewWatchlist(user, "Default")

		if err != nil {
			services.Error(err, "(CreateNewWatchlist) Unable to create watchlist Default")
			return
		}

		for key, row := range m {

			// Add some default symbols - SPY
			symb, err := t.DB.CreateNewSymbol(row.SymShort, row.SymLong, "Equity")

			if err != nil {
				services.Error(err, "(VerifyDefaultWatchList) Unable to create symbol "+row.SymShort)
				return
			}

			// Add lookup
			_, err2 := t.DB.CreateNewWatchlistSymbol(wList, symb, user, uint(key))

			if err2 != nil {
				services.Error(err2, "(CreateNewWatchlistSymbol) Unable to create symbol "+row.SymShort+" lookup")
				return
			}

		}

	}

	return

}

/* End File */
