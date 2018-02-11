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
		Page:  "2",
	})

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, len(results), 1)

	st.Expect(t, results[0].Id, uint(1))
	st.Expect(t, results[0].FirstName, "Rob")
	st.Expect(t, results[0].LastName, "Tester")
	st.Expect(t, results[0].Email, "spicer+robtester@options.cafe")
}

/* End File */
