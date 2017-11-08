//
// Date: 11/8/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package data_import

import (
	"os"
	"strconv"

	"app.options.cafe/backend/brokers/tradier"
	"app.options.cafe/backend/library/services"
)

var chars = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

//
// Connect to Tradier using our shared admin account and download all possible symbols.
//
func (t *Base) DoSymbolImport() {

	knownSymbols := make(map[string]string)

	// Map known symbols
	s := t.DB.GetAllSymbols()

	for _, row := range s {
		knownSymbols[row.ShortName] = row.Name
	}

	// Log...
	services.Log("[DoSymbolImport] Found " + strconv.Itoa(len(knownSymbols)) + " known symbols.")

	// Loop through each char and import into db.
	for _, row := range chars {
		t.ProcessLetter(row, knownSymbols)
	}
}

//
// Process a letter.
//
func (t *Base) ProcessLetter(letter string, knownSymbols map[string]string) {

	// Create new tradier instance
	tr := &tradier.Api{ApiKey: os.Getenv("TRADIER_ADMIN_ACCESS_TOKEN")}

	symbols, err := tr.SearchBySymbolName(letter)

	if err != nil {
		services.Error(err, "[DoSymbolImport] SearchBySymbolOrCompanyName failed.")
		return
	}

	// Log...
	services.Log("[DoSymbolImport] Processing letter " + letter + " found " + strconv.Itoa(len(symbols)) + " symbols from Tradier.")

	// Loop through each result and add to db.
	for _, row := range symbols {

		// Make sure we don't already have this symbol
		if _, ok := knownSymbols[row.Name]; ok {

			// TODO: See if the company name updated.....

			// Continue nothing to do.
			continue
		}

		// Add Symbol to our database.
		t.DB.CreateNewSymbol(row.Name, row.Description)
	}
}

/* End File */
