//
// Date: 2/11/2018
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
// Test - Get all symbols
//
func TestGetAllSymbols01(t *testing.T) {

	// Load config file.
	env.ReadEnv("../.env")

	// Start the db connection.
	db, _ := NewDB()
	defer db.Close()

	// Query and get test users
	syms := db.GetAllSymbols()

	// Verify data returned
	st.Expect(t, len(syms), 9)
	st.Expect(t, syms[0].Id, uint(1))
	st.Expect(t, syms[0].ShortName, "SPY")
	st.Expect(t, syms[0].Name, "SPDR S&P 500 ETF Trust")
	st.Expect(t, syms[0].Type, "Equity")

	st.Expect(t, syms[1].Id, uint(2))
	st.Expect(t, syms[1].ShortName, "MCD")
	st.Expect(t, syms[1].Name, "McDonald's Corp")
	st.Expect(t, syms[1].Type, "Equity")

	st.Expect(t, syms[2].Id, uint(3))
	st.Expect(t, syms[2].ShortName, "SBUX")
	st.Expect(t, syms[2].Name, "Starbucks Corp")
	st.Expect(t, syms[2].Type, "Equity")
}

//
// Test - Create a symbol
//
func TestCreateNewSymbol01(t *testing.T) {

	// Load config file.
	env.ReadEnv("../.env")

	// Start the db connection.
	db, _ := NewDB()
	defer db.Close()

	// Query and get test users
	sym, err := db.CreateNewSymbol("hd", "Home Depot", "Equity")

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, sym.Id, uint(10))
	st.Expect(t, sym.ShortName, "HD") // verify it turns to caps
	st.Expect(t, sym.Name, "Home Depot")
	st.Expect(t, sym.Type, "Equity")

	// ---- Now we test again. This function should not add the same symbol twice.

	// Query and get test users
	sym, err2 := db.CreateNewSymbol("hd", "Home Depot", "Equity")

	// Test results
	st.Expect(t, err2, nil)
	st.Expect(t, sym.Id, uint(10))
	st.Expect(t, sym.ShortName, "HD") // verify it turns to caps
	st.Expect(t, sym.Name, "Home Depot")
	st.Expect(t, sym.Type, "Equity")
}

//
// Test - Searching symbols
//
func TestSearchSymbols01(t *testing.T) {

	// Load config file.
	env.ReadEnv("../.env")

	// Start the db connection.
	db, _ := NewDB()
	defer db.Close()

	// Query and get test users
	syms, err := db.SearchSymbols("py", "Equity")

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, len(syms), 1)
	st.Expect(t, syms[0].ShortName, "SPY")
	st.Expect(t, syms[0].Name, "SPDR S&P 500 ETF Trust")
	st.Expect(t, syms[0].Type, "Equity")

	// Query and get test users
	syms, err2 := db.SearchSymbols("", "Equity")

	// Test results
	st.Expect(t, err2, nil)
	st.Expect(t, len(syms), 6)
	st.Expect(t, syms[0].Id, uint(1))
	st.Expect(t, syms[0].ShortName, "SPY")
	st.Expect(t, syms[0].Name, "SPDR S&P 500 ETF Trust")
	st.Expect(t, syms[0].Type, "Equity")

	st.Expect(t, syms[1].Id, uint(2))
	st.Expect(t, syms[1].ShortName, "MCD")
	st.Expect(t, syms[1].Name, "McDonald's Corp")
	st.Expect(t, syms[1].Type, "Equity")

	st.Expect(t, syms[2].Id, uint(3))
	st.Expect(t, syms[2].ShortName, "SBUX")
	st.Expect(t, syms[2].Name, "Starbucks Corp")
	st.Expect(t, syms[2].Type, "Equity")
}

//
// Test - Create New Option Symbol
//
func TestCreateNewOptionSymbol01(t *testing.T) {

	// Load config file.
	env.ReadEnv("../.env")

	// Start the db connection.
	db, _ := NewDB()
	defer db.Close()

	// Query and get test users
	sym, err := db.CreateNewOptionSymbol("SPY180209P00276000")

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, sym.Id, uint(10))
	st.Expect(t, sym.ShortName, "SPY180209P00276000")
	st.Expect(t, sym.Name, "SPY Feb 9, 2018 $276.00 Put")
	st.Expect(t, sym.Type, "Option")
}

/* End File */
