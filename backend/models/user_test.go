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
func TestGetAllUsers01(t *testing.T) {

	// Load config file.
	env.ReadEnv("../.env")

	// Start the db connection.
	db, _ := NewDB()
	defer db.Close()

	// Query and get test users
	users := db.GetAllUsers()

	// Verify data returned
	st.Expect(t, users[0].Id, uint(1))
	st.Expect(t, users[0].FirstName, "Rob")
	st.Expect(t, users[0].LastName, "Tester")
	st.Expect(t, users[0].Email, "spicer+robtester@options.cafe")

	st.Expect(t, users[1].Id, uint(2))
	st.Expect(t, users[1].FirstName, "Jane")
	st.Expect(t, users[1].LastName, "Wells")
	st.Expect(t, users[1].Email, "spicer+janewells@options.cafe")

	st.Expect(t, users[2].Id, uint(3))
	st.Expect(t, users[2].FirstName, "Bob")
	st.Expect(t, users[2].LastName, "Rosso")
	st.Expect(t, users[2].Email, "spicer+bobrosso@options.cafe")
}

//
// Validate an email address
//
func TestValidateEmailAddress(t *testing.T) {

	// Load config file.
	env.ReadEnv("../.env")

	// Start the db connection.
	db, _ := NewDB()
	defer db.Close()

	testData := map[string]bool{
		"spicer@options.cafe": true,
		"spicer matthews":     false,
		"@woot.com":           false,
		"me@example.com":      true,
	}

	for email, isReal := range testData {

		result := db.ValidateEmailAddress(email)

		if isReal && (result != nil) {
			t.Errorf("%s did not pass.", email)
		} else if (!isReal) && (result == nil) {
			t.Errorf("%s did not pass.", email)
		}

	}
}

//
// Test GenerateRandomBytes returns
//
func TestGenerateRandomBytes(t *testing.T) {

	// Load config file.
	env.ReadEnv("../.env")

	// Start the db connection.
	db, _ := NewDB()
	defer db.Close()

	randString, _ := db.GenerateRandomBytes(10)

	if len(randString) != 10 {
		t.Errorf("The random bytes was %d chars instead of 10.", len(randString))
	}

}

//
// Test GenerateRandomString returns
//
func TestGenerateRandomString(t *testing.T) {

	// Load config file.
	env.ReadEnv("../.env")

	// Start the db connection.
	db, _ := NewDB()
	defer db.Close()

	randString, _ := db.GenerateRandomString(10)

	if len(randString) != 10 {
		t.Errorf("The random string of %s was %d chars instead of 10.", randString, len(randString))
	}

}

/* End File */
