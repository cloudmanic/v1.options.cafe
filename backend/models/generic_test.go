//
// Date: 2/10/2018
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
// Test - Get all users
//
func TestQuery01(t *testing.T) {

	// Load config file.
	env.ReadEnv("../.env")

	// Start the db connection.
	db, _ := NewDB()
	defer db.Close()

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

	// Load config file.
	env.ReadEnv("../.env")

	// Start the db connection.
	db, _ := NewDB()
	defer db.Close()

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
	st.Expect(t, noFilterCount, 9)
	st.Expect(t, len(results), 2)
	st.Expect(t, results[0].Id, uint(1))
	st.Expect(t, results[1].Id, uint(2))
	st.Expect(t, meta.Page, 1)
	st.Expect(t, meta.Limit, 2)
	st.Expect(t, meta.Offset, 0)
	st.Expect(t, meta.PageCount, 5)
	st.Expect(t, meta.LimitCount, 2)
	st.Expect(t, meta.NoLimitCount, 9)

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
	st.Expect(t, noFilterCount, 3)
	st.Expect(t, results[0].Id, uint(1))
	st.Expect(t, results[1].Id, uint(4))
	st.Expect(t, meta.Page, 1)
	st.Expect(t, meta.Limit, 2)
	st.Expect(t, meta.Offset, 0)
	st.Expect(t, meta.PageCount, 2)
	st.Expect(t, meta.LimitCount, 2)
	st.Expect(t, meta.NoLimitCount, 3)
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
	st.Expect(t, len(results), 1)
	st.Expect(t, noFilterCount, 3)
	st.Expect(t, results[0].Id, uint(6))
	st.Expect(t, meta.Page, 2)
	st.Expect(t, meta.Limit, 2)
	st.Expect(t, meta.Offset, 2)
	st.Expect(t, meta.PageCount, 2)
	st.Expect(t, meta.LimitCount, 1)
	st.Expect(t, meta.NoLimitCount, 3)
	st.Expect(t, meta.LastPage, true)

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

	// Load config file.
	env.ReadEnv("../.env")

	// Start the db connection.
	db, _ := NewDB()
	defer db.Close()

	// ---------  Test 1 -------- //

	// Run the query
	count, err := db.Count(&Symbol{}, QueryParam{})

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, count, uint(9))

	// ---------  Test 2 -------- //

	// Run the query
	count, err2 := db.Count(&Symbol{}, QueryParam{Wheres: []KeyValue{{Key: "type", Value: "Equity"}}})

	// Test results
	st.Expect(t, err2, nil)
	st.Expect(t, count, uint(6))
}

/* End File */
