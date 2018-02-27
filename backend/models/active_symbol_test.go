//
// Date: 2/26/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"testing"

	env "github.com/jpfuentes2/go-env"
	"github.com/nbio/st"
)

//
// Test - CreateActiveSymbol
//
func TestCreateActiveSymbol01(t *testing.T) {

	// Load config file.
	env.ReadEnv("../.env")

	// Start the db connection.
	db, _ := NewDB()
	defer db.Close()

	// Make query
	syb1, _ := db.CreateActiveSymbol(1, "SPY")
	syb2, _ := db.CreateActiveSymbol(1, "HD")
	syb3, _ := db.CreateActiveSymbol(1, "SPY")

	// Test results
	st.Expect(t, syb1.Id, uint(1))
	st.Expect(t, syb2.Id, uint(2))
	st.Expect(t, syb3.Id, uint(1))

}

//
// Test - GetActiveSymbolsByUser
//
func TestGetActiveSymbolsByUser01(t *testing.T) {

	// Load config file.
	env.ReadEnv("../.env")

	// Start the db connection.
	db, _ := NewDB()
	defer db.Close()

	// Test data into place.
	db.CreateActiveSymbol(1, "SPY")
	db.CreateActiveSymbol(1, "HD")
	db.CreateActiveSymbol(1, "SPY")
	db.CreateActiveSymbol(2, "SPY")
	db.CreateActiveSymbol(3, "HD")
	db.CreateActiveSymbol(3, "SPY")
	db.CreateActiveSymbol(5, "SPY")
	db.CreateActiveSymbol(6, "HD")
	db.CreateActiveSymbol(7, "SPY")

	// Make query.
	results, err := db.GetActiveSymbolsByUser(3)

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, results[0].UserId, uint(3))
	st.Expect(t, results[0].Symbol, "HD")
	st.Expect(t, results[1].UserId, uint(3))
	st.Expect(t, results[1].Symbol, "SPY")
}

/* End File */
