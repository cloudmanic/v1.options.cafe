//
// Date: 11/09/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"net/http"

	"app.options.cafe/backend/library/services"
)

//
// Return symbols in our database.
//
func (t *Controller) GetSymbols(w http.ResponseWriter, r *http.Request) {

	// Get the GET parms
	search := r.URL.Query().Get("search")

	// Search for symbol
	if search != "" {
		t.DoSymbolSearch(w, r)
	}
}

//
// Do Symbol Search
//
func (t *Controller) DoSymbolSearch(w http.ResponseWriter, r *http.Request) {

	// Get the query.
	search := r.URL.Query().Get("search")

	// Run DB query
	symbols, err := t.DB.SearchSymbols(search)

	if err != nil {
		services.Error(err, "Controller:SearchBySymbolOrCompanyName() Mysql Call.")
		t.RespondError(w, http.StatusBadRequest, err.Error())
	}

	// Return happy JSON
	t.RespondJSON(w, http.StatusOK, symbols)
}

/* End File */
