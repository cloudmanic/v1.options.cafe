//
// Date: 2/11/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"testing"

	"app.options.cafe/library/helpers"
	"github.com/nbio/st"
)

//
// Test - Get all symbols
//
func TestGetAllSymbols01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := NewTestDB("")
	defer TestingTearDown(db, dbName)

	// Symbols
	db.Create(&Symbol{Name: "SPDR S&P 500 ETF Trust", ShortName: "SPY", Type: "Equity"})
	db.Create(&Symbol{Name: "McDonald's Corp", ShortName: "MCD", Type: "Equity"})
	db.Create(&Symbol{Name: "Starbucks Corp", ShortName: "SBUX", Type: "Equity"})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $253.00 Put", ShortName: "SPY180316P00253000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 253.00})
	db.Create(&Symbol{Name: "VXX Feb 23, 2018 $55.00 Call", ShortName: "VXX180223C00055000", Type: "Option", OptionUnderlying: "VXX", OptionExpire: Date{helpers.ParseDateNoError("2/23/2018").UTC()}, OptionType: "Call", OptionStrike: 55.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $266.00 Put", ShortName: "SPY180316P00266000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 266.00})
	db.Create(&Symbol{Name: "Caterpillar Inc", ShortName: "CAT", Type: "Equity"})
	db.Create(&Symbol{Name: "Ascent Solar Technologies Inc", ShortName: "ASTI", Type: "Equity"})
	db.Create(&Symbol{Name: "Advanced Micro Devices Inc", ShortName: "AMD", Type: "Equity"})
	db.Create(&Symbol{Name: "BARCLAYS BK PLC", ShortName: "VXX", Type: "Equity"})
	db.Create(&Symbol{Name: "SPY Feb 9, 2018 $276.00 Put", ShortName: "SPY180209P00276000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("2/9/2018").UTC()}, OptionType: "Put", OptionStrike: 276.00})
	db.Create(&Symbol{Name: "VXX Mar 2, 2018 $46.00 Put", ShortName: "VXX180302P00046000", Type: "Option", OptionUnderlying: "VXX", OptionExpire: Date{helpers.ParseDateNoError("3/2/2018").UTC()}, OptionType: "Put", OptionStrike: 46.00})
	db.Create(&Symbol{Name: "VXX Feb 23, 2018 $50.00 Put", ShortName: "VXX180223P00050000", Type: "Option", OptionUnderlying: "VXX", OptionExpire: Date{helpers.ParseDateNoError("2/23/2018").UTC()}, OptionType: "Put", OptionStrike: 50.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $267.00 Put", ShortName: "SPY180316P00267000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 267.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $269.00 Put", ShortName: "SPY180316P00269000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 269.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $268.00 Put", ShortName: "SPY180316P00268000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 268.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $270.00 Put", ShortName: "SPY180316P00270000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 270.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $264.00 Put", ShortName: "SPY180316P00264000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 264.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $266.00 Put", ShortName: "SPY180316P00266000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 266.00})
	db.Create(&Symbol{Name: "SPY Mar 9, 2018 $262.00 Put", ShortName: "SPY180309P00262000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/9/2018").UTC()}, OptionType: "Put", OptionStrike: 262.00})
	db.Create(&Symbol{Name: "SPY Mar 9, 2018 $264.00 Put", ShortName: "SPY180309P00264000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/9/2018").UTC()}, OptionType: "Put", OptionStrike: 264.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $251.00 Put", ShortName: "SPY180316P00251000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 251.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $253.00 Put", ShortName: "SPY180316P00253000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 253.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $250.00 Put", ShortName: "SPY180316P00250000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 250.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $252.00 Put", ShortName: "SPY180316P00252000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 252.00})
	db.Create(&Symbol{Name: "VXX Feb 23, 2018 $50.00 Call", ShortName: "VXX180223C00050000", Type: "Option", OptionUnderlying: "VXX", OptionExpire: Date{helpers.ParseDateNoError("2/23/2018").UTC()}, OptionType: "Call", OptionStrike: 50.00})
	db.Create(&Symbol{Name: "SPY Feb 28, 2018 $249.00 Put", ShortName: "SPY180228P00249000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("2/28/2018").UTC()}, OptionType: "Put", OptionStrike: 249.00})
	db.Create(&Symbol{Name: "SPY Feb 28, 2018 $251.00 Put", ShortName: "SPY180228P00251000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("2/28/2018").UTC()}, OptionType: "Put", OptionStrike: 251.00})
	db.Create(&Symbol{Name: "SPY Feb 21, 2018 $241.00 Put", ShortName: "SPY180221P00241000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("2/21/2018").UTC()}, OptionType: "Put", OptionStrike: 241.00})
	db.Create(&Symbol{Name: "SPY Feb 21, 2018 $243.00 Put", ShortName: "SPY180221P00243000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("2/21/2018").UTC()}, OptionType: "Put", OptionStrike: 243.00})
	db.Create(&Symbol{Name: "SPY Mar 21, 2018 $254.00 Put", ShortName: "SPY180321P00254000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/21/2018").UTC()}, OptionType: "Put", OptionStrike: 254.00})
	db.Create(&Symbol{Name: "SPY Mar 21, 2018 $256.00 Put", ShortName: "SPY180321P00256000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/21/2018").UTC()}, OptionType: "Put", OptionStrike: 256.00})
	db.Create(&Symbol{Name: "SPY Jul 21, 2017 $234.00 Put", ShortName: "SPY170721P00234000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("7/21/2018").UTC()}, OptionType: "Put", OptionStrike: 234.00})

	// Query and get test users
	syms := db.GetAllSymbols()

	// Verify data returned
	st.Expect(t, len(syms), 33)
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

	st.Expect(t, syms[3].Id, uint(4))
	st.Expect(t, syms[3].ShortName, "SPY180316P00253000")
	st.Expect(t, syms[3].Name, "SPY Mar 16, 2018 $253.00 Put")
	st.Expect(t, syms[3].Type, "Option")
	st.Expect(t, syms[3].OptionType, "Put")
	st.Expect(t, syms[3].OptionStrike, 253.00)
	st.Expect(t, syms[3].OptionExpire.Format("01/02/2006"), "03/16/2018")
}

//
// Test - Create a symbol
//
func TestCreateNewSymbol01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := NewTestDB("")
	defer TestingTearDown(db, dbName)

	// Symbols
	db.Create(&Symbol{Name: "SPDR S&P 500 ETF Trust", ShortName: "SPY", Type: "Equity"})
	db.Create(&Symbol{Name: "McDonald's Corp", ShortName: "MCD", Type: "Equity"})
	db.Create(&Symbol{Name: "Starbucks Corp", ShortName: "SBUX", Type: "Equity"})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $253.00 Put", ShortName: "SPY180316P00253000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 253.00})
	db.Create(&Symbol{Name: "VXX Feb 23, 2018 $55.00 Call", ShortName: "VXX180223C00055000", Type: "Option", OptionUnderlying: "VXX", OptionExpire: Date{helpers.ParseDateNoError("2/23/2018").UTC()}, OptionType: "Call", OptionStrike: 55.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $266.00 Put", ShortName: "SPY180316P00266000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 266.00})
	db.Create(&Symbol{Name: "Caterpillar Inc", ShortName: "CAT", Type: "Equity"})
	db.Create(&Symbol{Name: "Ascent Solar Technologies Inc", ShortName: "ASTI", Type: "Equity"})
	db.Create(&Symbol{Name: "Advanced Micro Devices Inc", ShortName: "AMD", Type: "Equity"})
	db.Create(&Symbol{Name: "BARCLAYS BK PLC", ShortName: "VXX", Type: "Equity"})
	db.Create(&Symbol{Name: "SPY Feb 9, 2018 $276.00 Put", ShortName: "SPY180209P00276000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("2/9/2018").UTC()}, OptionType: "Put", OptionStrike: 276.00})
	db.Create(&Symbol{Name: "VXX Mar 2, 2018 $46.00 Put", ShortName: "VXX180302P00046000", Type: "Option", OptionUnderlying: "VXX", OptionExpire: Date{helpers.ParseDateNoError("3/2/2018").UTC()}, OptionType: "Put", OptionStrike: 46.00})
	db.Create(&Symbol{Name: "VXX Feb 23, 2018 $50.00 Put", ShortName: "VXX180223P00050000", Type: "Option", OptionUnderlying: "VXX", OptionExpire: Date{helpers.ParseDateNoError("2/23/2018").UTC()}, OptionType: "Put", OptionStrike: 50.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $267.00 Put", ShortName: "SPY180316P00267000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 267.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $269.00 Put", ShortName: "SPY180316P00269000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 269.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $268.00 Put", ShortName: "SPY180316P00268000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 268.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $270.00 Put", ShortName: "SPY180316P00270000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 270.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $264.00 Put", ShortName: "SPY180316P00264000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 264.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $266.00 Put", ShortName: "SPY180316P00266000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 266.00})
	db.Create(&Symbol{Name: "SPY Mar 9, 2018 $262.00 Put", ShortName: "SPY180309P00262000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/9/2018").UTC()}, OptionType: "Put", OptionStrike: 262.00})
	db.Create(&Symbol{Name: "SPY Mar 9, 2018 $264.00 Put", ShortName: "SPY180309P00264000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/9/2018").UTC()}, OptionType: "Put", OptionStrike: 264.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $251.00 Put", ShortName: "SPY180316P00251000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 251.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $253.00 Put", ShortName: "SPY180316P00253000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 253.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $250.00 Put", ShortName: "SPY180316P00250000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 250.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $252.00 Put", ShortName: "SPY180316P00252000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 252.00})
	db.Create(&Symbol{Name: "VXX Feb 23, 2018 $50.00 Call", ShortName: "VXX180223C00050000", Type: "Option", OptionUnderlying: "VXX", OptionExpire: Date{helpers.ParseDateNoError("2/23/2018").UTC()}, OptionType: "Call", OptionStrike: 50.00})
	db.Create(&Symbol{Name: "SPY Feb 28, 2018 $249.00 Put", ShortName: "SPY180228P00249000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("2/28/2018").UTC()}, OptionType: "Put", OptionStrike: 249.00})
	db.Create(&Symbol{Name: "SPY Feb 28, 2018 $251.00 Put", ShortName: "SPY180228P00251000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("2/28/2018").UTC()}, OptionType: "Put", OptionStrike: 251.00})
	db.Create(&Symbol{Name: "SPY Feb 21, 2018 $241.00 Put", ShortName: "SPY180221P00241000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("2/21/2018").UTC()}, OptionType: "Put", OptionStrike: 241.00})
	db.Create(&Symbol{Name: "SPY Feb 21, 2018 $243.00 Put", ShortName: "SPY180221P00243000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("2/21/2018").UTC()}, OptionType: "Put", OptionStrike: 243.00})
	db.Create(&Symbol{Name: "SPY Mar 21, 2018 $254.00 Put", ShortName: "SPY180321P00254000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/21/2018").UTC()}, OptionType: "Put", OptionStrike: 254.00})
	db.Create(&Symbol{Name: "SPY Mar 21, 2018 $256.00 Put", ShortName: "SPY180321P00256000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/21/2018").UTC()}, OptionType: "Put", OptionStrike: 256.00})
	db.Create(&Symbol{Name: "SPY Jul 21, 2017 $234.00 Put", ShortName: "SPY170721P00234000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("7/21/2018").UTC()}, OptionType: "Put", OptionStrike: 234.00})

	// Query and get test users
	sym, err := db.CreateNewSymbol("hd", "Home Depot", "Equity")

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, sym.Id, uint(34))
	st.Expect(t, sym.ShortName, "HD") // verify it turns to caps
	st.Expect(t, sym.Name, "Home Depot")
	st.Expect(t, sym.Type, "Equity")

	// ---- Now we test again. This function should not add the same symbol twice.

	// Query and get test users
	sym, err2 := db.CreateNewSymbol("hd", "Home Depot", "Equity")

	// Test results
	st.Expect(t, err2, nil)
	st.Expect(t, sym.Id, uint(34))
	st.Expect(t, sym.ShortName, "HD") // verify it turns to caps
	st.Expect(t, sym.Name, "Home Depot")
	st.Expect(t, sym.Type, "Equity")
}

//
// Test - Searching symbols
//
func TestSearchSymbols01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := NewTestDB("")
	defer TestingTearDown(db, dbName)

	// Symbols
	db.Create(&Symbol{Name: "SPDR S&P 500 ETF Trust", ShortName: "SPY", Type: "Equity"})
	db.Create(&Symbol{Name: "McDonald's Corp", ShortName: "MCD", Type: "Equity"})
	db.Create(&Symbol{Name: "Starbucks Corp", ShortName: "SBUX", Type: "Equity"})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $253.00 Put", ShortName: "SPY180316P00253000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 253.00})
	db.Create(&Symbol{Name: "VXX Feb 23, 2018 $55.00 Call", ShortName: "VXX180223C00055000", Type: "Option", OptionUnderlying: "VXX", OptionExpire: Date{helpers.ParseDateNoError("2/23/2018").UTC()}, OptionType: "Call", OptionStrike: 55.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $266.00 Put", ShortName: "SPY180316P00266000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 266.00})
	db.Create(&Symbol{Name: "Caterpillar Inc", ShortName: "CAT", Type: "Equity"})
	db.Create(&Symbol{Name: "Ascent Solar Technologies Inc", ShortName: "ASTI", Type: "Equity"})
	db.Create(&Symbol{Name: "Advanced Micro Devices Inc", ShortName: "AMD", Type: "Equity"})
	db.Create(&Symbol{Name: "BARCLAYS BK PLC", ShortName: "VXX", Type: "Equity"})
	db.Create(&Symbol{Name: "SPY Feb 9, 2018 $276.00 Put", ShortName: "SPY180209P00276000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("2/9/2018").UTC()}, OptionType: "Put", OptionStrike: 276.00})
	db.Create(&Symbol{Name: "VXX Mar 2, 2018 $46.00 Put", ShortName: "VXX180302P00046000", Type: "Option", OptionUnderlying: "VXX", OptionExpire: Date{helpers.ParseDateNoError("3/2/2018").UTC()}, OptionType: "Put", OptionStrike: 46.00})
	db.Create(&Symbol{Name: "VXX Feb 23, 2018 $50.00 Put", ShortName: "VXX180223P00050000", Type: "Option", OptionUnderlying: "VXX", OptionExpire: Date{helpers.ParseDateNoError("2/23/2018").UTC()}, OptionType: "Put", OptionStrike: 50.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $267.00 Put", ShortName: "SPY180316P00267000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 267.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $269.00 Put", ShortName: "SPY180316P00269000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 269.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $268.00 Put", ShortName: "SPY180316P00268000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 268.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $270.00 Put", ShortName: "SPY180316P00270000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 270.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $264.00 Put", ShortName: "SPY180316P00264000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 264.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $266.00 Put", ShortName: "SPY180316P00266000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 266.00})
	db.Create(&Symbol{Name: "SPY Mar 9, 2018 $262.00 Put", ShortName: "SPY180309P00262000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/9/2018").UTC()}, OptionType: "Put", OptionStrike: 262.00})
	db.Create(&Symbol{Name: "SPY Mar 9, 2018 $264.00 Put", ShortName: "SPY180309P00264000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/9/2018").UTC()}, OptionType: "Put", OptionStrike: 264.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $251.00 Put", ShortName: "SPY180316P00251000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 251.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $253.00 Put", ShortName: "SPY180316P00253000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 253.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $250.00 Put", ShortName: "SPY180316P00250000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 250.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $252.00 Put", ShortName: "SPY180316P00252000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 252.00})
	db.Create(&Symbol{Name: "VXX Feb 23, 2018 $50.00 Call", ShortName: "VXX180223C00050000", Type: "Option", OptionUnderlying: "VXX", OptionExpire: Date{helpers.ParseDateNoError("2/23/2018").UTC()}, OptionType: "Call", OptionStrike: 50.00})
	db.Create(&Symbol{Name: "SPY Feb 28, 2018 $249.00 Put", ShortName: "SPY180228P00249000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("2/28/2018").UTC()}, OptionType: "Put", OptionStrike: 249.00})
	db.Create(&Symbol{Name: "SPY Feb 28, 2018 $251.00 Put", ShortName: "SPY180228P00251000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("2/28/2018").UTC()}, OptionType: "Put", OptionStrike: 251.00})
	db.Create(&Symbol{Name: "SPY Feb 21, 2018 $241.00 Put", ShortName: "SPY180221P00241000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("2/21/2018").UTC()}, OptionType: "Put", OptionStrike: 241.00})
	db.Create(&Symbol{Name: "SPY Feb 21, 2018 $243.00 Put", ShortName: "SPY180221P00243000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("2/21/2018").UTC()}, OptionType: "Put", OptionStrike: 243.00})
	db.Create(&Symbol{Name: "SPY Mar 21, 2018 $254.00 Put", ShortName: "SPY180321P00254000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/21/2018").UTC()}, OptionType: "Put", OptionStrike: 254.00})
	db.Create(&Symbol{Name: "SPY Mar 21, 2018 $256.00 Put", ShortName: "SPY180321P00256000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/21/2018").UTC()}, OptionType: "Put", OptionStrike: 256.00})
	db.Create(&Symbol{Name: "SPY Jul 21, 2017 $234.00 Put", ShortName: "SPY170721P00234000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("7/21/2018").UTC()}, OptionType: "Put", OptionStrike: 234.00})

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
	st.Expect(t, len(syms), 7)
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
	// Start the db connection.
	db, dbName, _ := NewTestDB("")
	defer TestingTearDown(db, dbName)

	// Symbols
	db.Create(&Symbol{Name: "SPDR S&P 500 ETF Trust", ShortName: "SPY", Type: "Equity"})
	db.Create(&Symbol{Name: "McDonald's Corp", ShortName: "MCD", Type: "Equity"})
	db.Create(&Symbol{Name: "Starbucks Corp", ShortName: "SBUX", Type: "Equity"})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $253.00 Put", ShortName: "SPY180316P00253000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 253.00})
	db.Create(&Symbol{Name: "VXX Feb 23, 2018 $55.00 Call", ShortName: "VXX180223C00055000", Type: "Option", OptionUnderlying: "VXX", OptionExpire: Date{helpers.ParseDateNoError("2/23/2018").UTC()}, OptionType: "Call", OptionStrike: 55.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $266.00 Put", ShortName: "SPY180316P00266000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 266.00})
	db.Create(&Symbol{Name: "Caterpillar Inc", ShortName: "CAT", Type: "Equity"})
	db.Create(&Symbol{Name: "Ascent Solar Technologies Inc", ShortName: "ASTI", Type: "Equity"})
	db.Create(&Symbol{Name: "Advanced Micro Devices Inc", ShortName: "AMD", Type: "Equity"})
	db.Create(&Symbol{Name: "BARCLAYS BK PLC", ShortName: "VXX", Type: "Equity"})
	db.Create(&Symbol{Name: "SPY Feb 9, 2018 $276.00 Put", ShortName: "SPY180209P00276000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("2/9/2018").UTC()}, OptionType: "Put", OptionStrike: 276.00})
	db.Create(&Symbol{Name: "VXX Mar 2, 2018 $46.00 Put", ShortName: "VXX180302P00046000", Type: "Option", OptionUnderlying: "VXX", OptionExpire: Date{helpers.ParseDateNoError("3/2/2018").UTC()}, OptionType: "Put", OptionStrike: 46.00})
	db.Create(&Symbol{Name: "VXX Feb 23, 2018 $50.00 Put", ShortName: "VXX180223P00050000", Type: "Option", OptionUnderlying: "VXX", OptionExpire: Date{helpers.ParseDateNoError("2/23/2018").UTC()}, OptionType: "Put", OptionStrike: 50.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $267.00 Put", ShortName: "SPY180316P00267000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 267.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $269.00 Put", ShortName: "SPY180316P00269000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 269.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $268.00 Put", ShortName: "SPY180316P00268000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 268.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $270.00 Put", ShortName: "SPY180316P00270000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 270.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $264.00 Put", ShortName: "SPY180316P00264000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 264.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $266.00 Put", ShortName: "SPY180316P00266000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 266.00})
	db.Create(&Symbol{Name: "SPY Mar 9, 2018 $262.00 Put", ShortName: "SPY180309P00262000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/9/2018").UTC()}, OptionType: "Put", OptionStrike: 262.00})
	db.Create(&Symbol{Name: "SPY Mar 9, 2018 $264.00 Put", ShortName: "SPY180309P00264000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/9/2018").UTC()}, OptionType: "Put", OptionStrike: 264.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $251.00 Put", ShortName: "SPY180316P00251000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 251.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $253.00 Put", ShortName: "SPY180316P00253000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 253.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $250.00 Put", ShortName: "SPY180316P00250000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 250.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $252.00 Put", ShortName: "SPY180316P00252000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 252.00})
	db.Create(&Symbol{Name: "VXX Feb 23, 2018 $50.00 Call", ShortName: "VXX180223C00050000", Type: "Option", OptionUnderlying: "VXX", OptionExpire: Date{helpers.ParseDateNoError("2/23/2018").UTC()}, OptionType: "Call", OptionStrike: 50.00})
	db.Create(&Symbol{Name: "SPY Feb 28, 2018 $249.00 Put", ShortName: "SPY180228P00249000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("2/28/2018").UTC()}, OptionType: "Put", OptionStrike: 249.00})
	db.Create(&Symbol{Name: "SPY Feb 28, 2018 $251.00 Put", ShortName: "SPY180228P00251000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("2/28/2018").UTC()}, OptionType: "Put", OptionStrike: 251.00})
	db.Create(&Symbol{Name: "SPY Feb 21, 2018 $241.00 Put", ShortName: "SPY180221P00241000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("2/21/2018").UTC()}, OptionType: "Put", OptionStrike: 241.00})
	db.Create(&Symbol{Name: "SPY Feb 21, 2018 $243.00 Put", ShortName: "SPY180221P00243000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("2/21/2018").UTC()}, OptionType: "Put", OptionStrike: 243.00})
	db.Create(&Symbol{Name: "SPY Mar 21, 2018 $254.00 Put", ShortName: "SPY180321P00254000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/21/2018").UTC()}, OptionType: "Put", OptionStrike: 254.00})
	db.Create(&Symbol{Name: "SPY Mar 21, 2018 $256.00 Put", ShortName: "SPY180321P00256000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/21/2018").UTC()}, OptionType: "Put", OptionStrike: 256.00})
	db.Create(&Symbol{Name: "SPY Jul 21, 2017 $234.00 Put", ShortName: "SPY170721P00234000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("7/21/2018").UTC()}, OptionType: "Put", OptionStrike: 234.00})

	// Query and get test users
	sym, err := db.CreateNewOptionSymbol("SPY180209P00276000")

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, sym.Id, uint(11))
	st.Expect(t, sym.ShortName, "SPY180209P00276000")
	st.Expect(t, sym.Name, "SPY Feb 9, 2018 $276.00 Put")
	st.Expect(t, sym.Type, "Option")
}

//
// Test - GetOptionByParts
//
func TestGetOptionByParts01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := NewTestDB("")
	defer TestingTearDown(db, dbName)

	// Symbols
	db.Create(&Symbol{Name: "SPDR S&P 500 ETF Trust", ShortName: "SPY", Type: "Equity"})
	db.Create(&Symbol{Name: "McDonald's Corp", ShortName: "MCD", Type: "Equity"})
	db.Create(&Symbol{Name: "Starbucks Corp", ShortName: "SBUX", Type: "Equity"})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $253.00 Put", ShortName: "SPY180316P00253000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 253.00})
	db.Create(&Symbol{Name: "VXX Feb 23, 2018 $55.00 Call", ShortName: "VXX180223C00055000", Type: "Option", OptionUnderlying: "VXX", OptionExpire: Date{helpers.ParseDateNoError("2/23/2018").UTC()}, OptionType: "Call", OptionStrike: 55.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $266.00 Put", ShortName: "SPY180316P00266000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 266.00})
	db.Create(&Symbol{Name: "Caterpillar Inc", ShortName: "CAT", Type: "Equity"})
	db.Create(&Symbol{Name: "Ascent Solar Technologies Inc", ShortName: "ASTI", Type: "Equity"})
	db.Create(&Symbol{Name: "Advanced Micro Devices Inc", ShortName: "AMD", Type: "Equity"})
	db.Create(&Symbol{Name: "BARCLAYS BK PLC", ShortName: "VXX", Type: "Equity"})
	db.Create(&Symbol{Name: "SPY Feb 9, 2018 $276.00 Put", ShortName: "SPY180209P00276000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("2/9/2018").UTC()}, OptionType: "Put", OptionStrike: 276.00})
	db.Create(&Symbol{Name: "VXX Mar 2, 2018 $46.00 Put", ShortName: "VXX180302P00046000", Type: "Option", OptionUnderlying: "VXX", OptionExpire: Date{helpers.ParseDateNoError("3/2/2018").UTC()}, OptionType: "Put", OptionStrike: 46.00})
	db.Create(&Symbol{Name: "VXX Feb 23, 2018 $50.00 Put", ShortName: "VXX180223P00050000", Type: "Option", OptionUnderlying: "VXX", OptionExpire: Date{helpers.ParseDateNoError("2/23/2018").UTC()}, OptionType: "Put", OptionStrike: 50.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $267.00 Put", ShortName: "SPY180316P00267000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 267.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $269.00 Put", ShortName: "SPY180316P00269000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 269.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $268.00 Put", ShortName: "SPY180316P00268000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 268.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $270.00 Put", ShortName: "SPY180316P00270000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 270.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $264.00 Put", ShortName: "SPY180316P00264000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 264.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $266.00 Put", ShortName: "SPY180316P00266000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 266.00})
	db.Create(&Symbol{Name: "SPY Mar 9, 2018 $262.00 Put", ShortName: "SPY180309P00262000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/9/2018").UTC()}, OptionType: "Put", OptionStrike: 262.00})
	db.Create(&Symbol{Name: "SPY Mar 9, 2018 $264.00 Put", ShortName: "SPY180309P00264000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/9/2018").UTC()}, OptionType: "Put", OptionStrike: 264.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $251.00 Put", ShortName: "SPY180316P00251000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 251.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $253.00 Put", ShortName: "SPY180316P00253000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 253.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $250.00 Put", ShortName: "SPY180316P00250000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 250.00})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $252.00 Put", ShortName: "SPY180316P00252000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 252.00})
	db.Create(&Symbol{Name: "VXX Feb 23, 2018 $50.00 Call", ShortName: "VXX180223C00050000", Type: "Option", OptionUnderlying: "VXX", OptionExpire: Date{helpers.ParseDateNoError("2/23/2018").UTC()}, OptionType: "Call", OptionStrike: 50.00})
	db.Create(&Symbol{Name: "SPY Feb 28, 2018 $249.00 Put", ShortName: "SPY180228P00249000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("2/28/2018").UTC()}, OptionType: "Put", OptionStrike: 249.00})
	db.Create(&Symbol{Name: "SPY Feb 28, 2018 $251.00 Put", ShortName: "SPY180228P00251000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("2/28/2018").UTC()}, OptionType: "Put", OptionStrike: 251.00})
	db.Create(&Symbol{Name: "SPY Feb 21, 2018 $241.00 Put", ShortName: "SPY180221P00241000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("2/21/2018").UTC()}, OptionType: "Put", OptionStrike: 241.00})
	db.Create(&Symbol{Name: "SPY Feb 21, 2018 $243.00 Put", ShortName: "SPY180221P00243000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("2/21/2018").UTC()}, OptionType: "Put", OptionStrike: 243.00})
	db.Create(&Symbol{Name: "SPY Mar 21, 2018 $254.00 Put", ShortName: "SPY180321P00254000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/21/2018").UTC()}, OptionType: "Put", OptionStrike: 254.00})
	db.Create(&Symbol{Name: "SPY Mar 21, 2018 $256.00 Put", ShortName: "SPY180321P00256000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("3/21/2018").UTC()}, OptionType: "Put", OptionStrike: 256.00})
	db.Create(&Symbol{Name: "SPY Jul 21, 2017 $234.00 Put", ShortName: "SPY170721P00234000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: Date{helpers.ParseDateNoError("7/21/2018").UTC()}, OptionType: "Put", OptionStrike: 234.00})

	// Set date.
	expireDate := helpers.ParseDateNoError("2018-02-09")

	// Query database
	sym, err := db.GetOptionByParts("SPY", "Put", expireDate, 276)

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, sym.Id, uint(11))
	st.Expect(t, sym.ShortName, "SPY180209P00276000")
	st.Expect(t, sym.Name, "SPY Feb 9, 2018 $276.00 Put")
	st.Expect(t, sym.OptionUnderlying, "SPY")
	st.Expect(t, sym.OptionType, "Put")
	st.Expect(t, sym.OptionExpire.Format("2006-01-02"), expireDate.Format("2006-01-02"))
	st.Expect(t, sym.OptionStrike, 276.00)
}

/* End File */
