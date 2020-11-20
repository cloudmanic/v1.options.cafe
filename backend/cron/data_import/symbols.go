//
// Date: 11/8/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package data_import

import (
	"errors"
	"net/http"
	"os"
	"strconv"

	"app.options.cafe/brokers/tradier"
	"app.options.cafe/library/services"
	"app.options.cafe/models"
)

var chars = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

//
// Connect to Tradier using our shared admin account and download all possible symbols.
//
func DoSymbolImport(db *models.DB) {

	knownSymbols := make(map[string]models.Symbol)

	// Map known symbols
	s := db.GetAllSymbols()

	for _, row := range s {
		knownSymbols[row.ShortName] = row
	}

	// Log...
	services.InfoMsg("[DoSymbolImport] Found " + strconv.Itoa(len(knownSymbols)) + " known symbols.")

	// Loop through each char and import into db.
	for _, row := range chars {
		ProcessLetter(db, row, knownSymbols)
	}

	// Send health check notice.
	if len(os.Getenv("HEALTH_CHECK_SYMBOLS_IMPORT_URL")) > 0 {

		resp, err := http.Get(os.Getenv("HEALTH_CHECK_SYMBOLS_IMPORT_URL"))

		if err != nil {
			services.Critical(errors.New(err.Error() + "Could send health check - " + os.Getenv("HEALTH_CHECK_SYMBOLS_IMPORT_URL")))
		}

		defer resp.Body.Close()

	}
}

//
// Process a letter.
//
func ProcessLetter(db *models.DB, letter string, knownSymbols map[string]models.Symbol) {

	// Create new tradier instance
	tr := &tradier.Api{ApiKey: os.Getenv("TRADIER_ADMIN_ACCESS_TOKEN")}

	symbols, err := tr.SearchBySymbolName(letter)

	if err != nil {
		services.Info(errors.New(err.Error() + "[DoSymbolImport] SearchBySymbolOrCompanyName failed."))
		return
	}

	// Log...
	services.InfoMsg("[DoSymbolImport] Processing letter " + letter + " found " + strconv.Itoa(len(symbols)) + " symbols from Tradier.")

	// Loop through each result and add to db.
	for _, row := range symbols {

		// Make sure we don't already have this symbol
		if _, ok := knownSymbols[row.Name]; ok {

			// TODO: See if the company name updated.....
			if row.Description != knownSymbols[row.Name].Name {
				db.UpdateSymbol(knownSymbols[row.Name].Id, row.Name, row.Description, "Equity")
			}

			// Continue nothing to do.
			continue
		}

		// Add Symbol to our database.
		db.CreateNewSymbol(row.Name, row.Description, "Equity")
	}
}

/* End File */
