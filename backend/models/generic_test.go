//
// Date: 2/10/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"testing"

	"github.com/nbio/st"
	"github.com/optionscafe/options-cafe-cli/helpers"
)

//
// Test - CreateNewRecord 01
//
func TestCreateNewRecord01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := NewTestDB("")
	defer TestingTearDown(db, dbName)

	// Add a new user test.
	user := &User{FirstName: "John", LastName: "Smith", Email: "js@cloudmanic.com", Password: "fake-password"}
	err := db.CreateNewRecord(&user, InsertParam{})
	st.Expect(t, err, nil)

	// Add a new user test.
	user2 := &User{FirstName: "John 2", LastName: "Smith", Email: "js2@cloudmanic.com", Password: "fake-password"}
	err = db.CreateNewRecord(&user2, InsertParam{})
	st.Expect(t, err, nil)

	// Test results
	st.Expect(t, user.Id, uint(1))
	st.Expect(t, user.FirstName, "John")
	st.Expect(t, user2.Id, uint(2))
	st.Expect(t, user2.FirstName, "John 2")
}

//
// Test - Get all users
//
func TestQuery01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := NewTestDB("")
	defer TestingTearDown(db, dbName)

	// Users
	db.Create(&User{FirstName: "Rob", LastName: "Tester", Email: "spicer+robtester@options.cafe", Status: "Active"})
	db.Create(&User{FirstName: "Jane", LastName: "Wells", Email: "spicer+janewells@options.cafe", Status: "Active"})
	db.Create(&User{FirstName: "Bob", LastName: "Rosso", Email: "spicer+bobrosso@options.cafe", Status: "Active"})

	// ---------  Test 1 -------- //

	// Place to store the results.
	var results = []User{}

	// Run the query
	err := db.Query(&results, QueryParam{})

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, len(results), 3)

	st.Expect(t, results[0].FirstName, "Rob")
	st.Expect(t, results[0].LastName, "Tester")
	st.Expect(t, results[0].Email, "spicer+robtester@options.cafe")

	st.Expect(t, results[1].FirstName, "Jane")
	st.Expect(t, results[1].LastName, "Wells")
	st.Expect(t, results[1].Email, "spicer+janewells@options.cafe")

	st.Expect(t, results[2].FirstName, "Bob")
	st.Expect(t, results[2].LastName, "Rosso")
	st.Expect(t, results[2].Email, "spicer+bobrosso@options.cafe")

	// ---------  Test 2 -------- //

	// Place to store the results.
	results = []User{}

	// Another test to see if search works
	err = db.Query(&results, QueryParam{
		SearchTerm: "wells",
		SearchCols: []string{"id", "first_name", "last_name", "email"},
	})

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, len(results), 1)

	st.Expect(t, results[0].FirstName, "Jane")
	st.Expect(t, results[0].LastName, "Wells")
	st.Expect(t, results[0].Email, "spicer+janewells@options.cafe")

	// ---------  Test 3 -------- //

	// Place to store the results.
	results = []User{}

	// Another test to see if search works
	err = db.Query(&results, QueryParam{
		Wheres: []KeyValue{{Key: "email", Value: "spicer+bobrosso@options.cafe"}},
	})

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, len(results), 1)

	st.Expect(t, results[0].Id, uint(3))
	st.Expect(t, results[0].FirstName, "Bob")
	st.Expect(t, results[0].LastName, "Rosso")
	st.Expect(t, results[0].Email, "spicer+bobrosso@options.cafe")

	// ---------  Test 4 -------- //

	// Place to store the results.
	results = []User{}

	// Another test to see if search works
	err = db.Query(&results, QueryParam{
		Order: "last_name",
		Sort:  "asc",
	})

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, len(results), 3)

	st.Expect(t, results[0].Id, uint(3))
	st.Expect(t, results[0].FirstName, "Bob")
	st.Expect(t, results[0].LastName, "Rosso")
	st.Expect(t, results[0].Email, "spicer+bobrosso@options.cafe")

	st.Expect(t, results[1].Id, uint(1))
	st.Expect(t, results[1].FirstName, "Rob")
	st.Expect(t, results[1].LastName, "Tester")
	st.Expect(t, results[1].Email, "spicer+robtester@options.cafe")

	st.Expect(t, results[2].Id, uint(2))
	st.Expect(t, results[2].FirstName, "Jane")
	st.Expect(t, results[2].LastName, "Wells")
	st.Expect(t, results[2].Email, "spicer+janewells@options.cafe")

	// ---------  Test 5 -------- //

	// Place to store the results.
	results = []User{}

	// Another test to see if search works
	err = db.Query(&results, QueryParam{
		Order: "last_name",
		Sort:  "desc",
	})

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, len(results), 3)

	st.Expect(t, results[0].Id, uint(2))
	st.Expect(t, results[0].FirstName, "Jane")
	st.Expect(t, results[0].LastName, "Wells")
	st.Expect(t, results[0].Email, "spicer+janewells@options.cafe")

	st.Expect(t, results[1].Id, uint(1))
	st.Expect(t, results[1].FirstName, "Rob")
	st.Expect(t, results[1].LastName, "Tester")
	st.Expect(t, results[1].Email, "spicer+robtester@options.cafe")

	st.Expect(t, results[2].Id, uint(3))
	st.Expect(t, results[2].FirstName, "Bob")
	st.Expect(t, results[2].LastName, "Rosso")
	st.Expect(t, results[2].Email, "spicer+bobrosso@options.cafe")

	// ---------  Test 6 -------- //

	// Place to store the results.
	results = []User{}

	// Another test to see if search works
	err = db.Query(&results, QueryParam{
		Order: "last_name",
		Sort:  "desc",
		Limit: 1,
	})

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, len(results), 1)

	st.Expect(t, results[0].Id, uint(2))
	st.Expect(t, results[0].FirstName, "Jane")
	st.Expect(t, results[0].LastName, "Wells")
	st.Expect(t, results[0].Email, "spicer+janewells@options.cafe")

	// ---------  Test 7 -------- //

	// Place to store the results.
	results = []User{}

	// Another test to see if search works
	err = db.Query(&results, QueryParam{
		Order:  "last_name",
		Sort:   "desc",
		Limit:  1,
		Offset: 2,
	})

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, len(results), 1)

	st.Expect(t, results[0].Id, uint(3))
	st.Expect(t, results[0].FirstName, "Bob")
	st.Expect(t, results[0].LastName, "Rosso")
	st.Expect(t, results[0].Email, "spicer+bobrosso@options.cafe")

	// ---------  Test 8 -------- //

	// Place to store the results.
	results = []User{}

	// Another test to see if search works
	err = db.Query(&results, QueryParam{
		Order: "last_name",
		Sort:  "desc",
		Limit: 1,
		Page:  2,
	})

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, len(results), 1)

	st.Expect(t, results[0].Id, uint(1))
	st.Expect(t, results[0].FirstName, "Rob")
	st.Expect(t, results[0].LastName, "Tester")
	st.Expect(t, results[0].Email, "spicer+robtester@options.cafe")
}

//
// Test 02 - Paging
//
func TestQuery02(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := NewTestDB("")
	defer TestingTearDown(db, dbName)

	// Users
	db.Create(&User{FirstName: "Rob", LastName: "Tester", Email: "spicer+robtester@options.cafe", Status: "Active"})
	db.Create(&User{FirstName: "Jane", LastName: "Wells", Email: "spicer+janewells@options.cafe", Status: "Active"})
	db.Create(&User{FirstName: "Bob", LastName: "Rosso", Email: "spicer+bobrosso@options.cafe", Status: "Active"})

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

	// ---------  Test 1 -------- //

	// Place to store the results.
	results := []Symbol{}

	// Another test to see if paging works
	q1 := QueryParam{Limit: 2, Page: 1, Debug: false}
	noFilterCount, err := db.QueryWithNoFilterCount(&results, q1)

	// Get the meta data related to this query.
	meta := db.GetQueryMetaData(len(results), noFilterCount, q1)

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, noFilterCount, 33)
	st.Expect(t, len(results), 2)
	st.Expect(t, results[0].Id, uint(1))
	st.Expect(t, results[1].Id, uint(2))
	st.Expect(t, meta.Page, 1)
	st.Expect(t, meta.Limit, 2)
	st.Expect(t, meta.Offset, 0)
	st.Expect(t, meta.PageCount, 17)
	st.Expect(t, meta.LimitCount, 2)
	st.Expect(t, meta.NoLimitCount, 33)

	// ---------  Test 2 -------- //

	// Another test to see if paging works
	err = db.Query(&results, QueryParam{
		Limit: 2,
		Page:  3,
		Debug: false,
	})

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, len(results), 2)
	st.Expect(t, results[0].Id, uint(5))
	st.Expect(t, results[1].Id, uint(6))

	// ---------  Test 3 -------- //

	// Setup query parm
	q3 := QueryParam{
		Limit:      2,
		Page:       1,
		Debug:      false,
		SearchCols: []string{"short_name"},
		SearchTerm: "SPY",
	}

	// Another test to see if paging works
	noFilterCount, err = db.QueryWithNoFilterCount(&results, q3)

	// Get the meta data related to this query.
	meta = db.GetQueryMetaData(len(results), noFilterCount, q3)

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, len(results), 2)
	st.Expect(t, noFilterCount, 23)
	st.Expect(t, results[0].Id, uint(1))
	st.Expect(t, results[1].Id, uint(4))
	st.Expect(t, meta.Page, 1)
	st.Expect(t, meta.Limit, 2)
	st.Expect(t, meta.Offset, 0)
	st.Expect(t, meta.PageCount, 12)
	st.Expect(t, meta.LimitCount, 2)
	st.Expect(t, meta.NoLimitCount, 23)
	st.Expect(t, meta.LastPage, false)

	// ---------  Test 4 -------- //

	// Setup query parm
	q4 := QueryParam{
		Limit:      2,
		Page:       2,
		Debug:      false,
		SearchCols: []string{"short_name"},
		SearchTerm: "SPY",
	}

	// Another test to see if paging works
	noFilterCount, err = db.QueryWithNoFilterCount(&results, q4)

	// Get the meta data related to this query.
	meta = db.GetQueryMetaData(len(results), noFilterCount, q4)

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, len(results), 2)
	st.Expect(t, noFilterCount, 23)
	st.Expect(t, results[0].Id, uint(6))
	st.Expect(t, meta.Page, 2)
	st.Expect(t, meta.Limit, 2)
	st.Expect(t, meta.Offset, 2)
	st.Expect(t, meta.PageCount, 12)
	st.Expect(t, meta.LimitCount, 2)
	st.Expect(t, meta.NoLimitCount, 23)
	st.Expect(t, meta.LastPage, false)

	// ---------  Test 5 -------- //

	// Place to store the results.
	results = []Symbol{}

	// Another test to see if paging works
	q5 := QueryParam{}
	noFilterCount, err = db.QueryWithNoFilterCount(&results, q5)

	// Get the meta data related to this query.
	meta = db.GetQueryMetaData(len(results), noFilterCount, q5)

	// Test results - Testing Offset when no parms are sent in.
	st.Expect(t, err, nil)
	st.Expect(t, meta.Page, 0)
	st.Expect(t, meta.Limit, 0)
	st.Expect(t, meta.Offset, 0)
}

//
// Test - Count
//
func TestCount01(t *testing.T) {
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

	// ---------  Test 1 -------- //

	// Run the query
	count, err := db.Count(&Symbol{}, QueryParam{})

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, count, uint(33))

	// ---------  Test 2 -------- //

	// Run the query
	count, err2 := db.Count(&Symbol{}, QueryParam{Wheres: []KeyValue{{Key: "type", Value: "Equity"}}})

	// Test results
	st.Expect(t, err2, nil)
	st.Expect(t, count, uint(7))
}

/* End File */
