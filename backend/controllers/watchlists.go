//
// Date: 11/10/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/tidwall/gjson"
)

//
// Return watchlists in our database.
//
func (t *Controller) GetWatchlists(w http.ResponseWriter, r *http.Request) {

	// Get the watchlists
	wLists, err := t.DB.GetWatchlistsByUserId(1)

	if t.DoRespondError(w, err, httpGenericErrMsg) {
		return
	}

	// Return happy JSON
	t.RespondJSON(w, http.StatusOK, wLists)
}

//
// Return watchlist in our database.
//
func (t *Controller) GetWatchlist(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Set as int
	id, err := strconv.ParseInt(vars["id"], 10, 32)

	if t.DoRespondError(w, err, httpGenericErrMsg) {
		return
	}

	// Get the watchlist by id.
	wLists, err := t.DB.GetWatchlistsById(uint(id))

	if t.DoRespondError(w, err, httpNoRecordFound) {
		return
	}

	// Return happy JSON
	t.RespondJSON(w, http.StatusOK, wLists)
}

//
// Watchlist - Create
//
// curl -H "Content-Type: application/json" -X POST -d '{"name":"Super Cool Watchlist"}' -H "Authorization: Bearer XXXXXX" http://localhost:7080/api/v1/watchlists
//
func (t *Controller) CreateWatchlist(w http.ResponseWriter, r *http.Request) {

	// Parse json body
	body, err := ioutil.ReadAll(r.Body)

	if t.DoRespondError(w, err, httpGenericErrMsg) {
		return
	}

	name := gjson.Get(string(body), "name").String()

	// Get the watchlists
	wLists, err := t.DB.CreateWatchlist(1, name)

	if t.DoRespondError(w, err, httpGenericErrMsg) {
		return
	}

	// Return happy JSON
	t.RespondJSON(w, http.StatusOK, wLists)
}

/* End File */
