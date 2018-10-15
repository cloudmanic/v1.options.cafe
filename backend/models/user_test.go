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
// Test - CreateNewUserWithStripe
//
func TestCreateNewUserWithStripe01(t *testing.T) {

	// Load config file.
	env.ReadEnv("../.env")

	// Start the db connection.
	db, _ := NewDB()
	defer db.Close()

	// Create a test user.
	user := User{
		FirstName: "Jane",
		LastName:  "Unittester",
		Email:     "jane+unittest@options.cafe",
	}

	err := db.CreateNewUserWithStripe(user, "", "", "")

	// Since we are testing we know the user is id 4
	dbUser, _ := db.GetUserById(4)

	// Verify data returned
	st.Expect(t, err, nil)
	st.Expect(t, dbUser.Id, uint(4))
	st.Expect(t, dbUser.FirstName, "Jane")
	st.Expect(t, dbUser.LastName, "Unittester")
	st.Expect(t, dbUser.Email, "jane+unittest@options.cafe")
	st.Expect(t, len(dbUser.StripeCustomer), 18)
	st.Expect(t, len(dbUser.StripeSubscription), 18)

	// Clean things up at stripe
	db.DeleteUserWithStripe(dbUser)
}

//
// Test - GetSubscriptionWithStripe
//
func TestGetSubscriptionWithStripe01(t *testing.T) {

	// Load config file.
	env.ReadEnv("../.env")

	// Start the db connection.
	db, _ := NewDB()
	defer db.Close()

	// Create a test user.
	user := User{
		FirstName: "Jane",
		LastName:  "Unittester",
		Email:     "jane+unittest@options.cafe",
	}

	err := db.CreateNewUserWithStripe(user, "", "", "")

	// Since we are testing we know the user is id 4
	dbUser, _ := db.GetUserById(4)

	// Verify data returned
	st.Expect(t, err, nil)
	st.Expect(t, dbUser.Id, uint(4))
	st.Expect(t, dbUser.FirstName, "Jane")
	st.Expect(t, dbUser.LastName, "Unittester")
	st.Expect(t, dbUser.Email, "jane+unittest@options.cafe")
	st.Expect(t, len(dbUser.StripeCustomer), 18)
	st.Expect(t, len(dbUser.StripeSubscription), 18)

	// Add a credit card.
	err = db.UpdateCreditCard(dbUser, "tok_visa")

	// Verify data returned
	st.Expect(t, err, nil)

	// Add a second credit card for testing we want to verify stripe only has one card at a time.
	err = db.UpdateCreditCard(dbUser, "tok_amex")

	// Verify data returned
	st.Expect(t, err, nil)

	// Get subscription with stripe
	sub, err := db.GetSubscriptionWithStripe(dbUser)

	// Verify data returned
	st.Expect(t, err, nil)
	st.Expect(t, sub.Name, "Monthly - $20")
	st.Expect(t, sub.Amount, 20.00)
	st.Expect(t, sub.TrialDays, 7)
	st.Expect(t, sub.BillingInterval, "month")
	st.Expect(t, sub.CardBrand, "American Express")
	st.Expect(t, sub.CardLast4, "8431")
	st.Expect(t, sub.CardExpMonth, 10)
	st.Expect(t, sub.CardExpYear, 2019)

	// Clean things up at stripe
	db.DeleteUserWithStripe(dbUser)
}

//
// Test - GetSubscriptionWithStripe - No card on file
//
func TestGetSubscriptionWithStripe02(t *testing.T) {

	// Load config file.
	env.ReadEnv("../.env")

	// Start the db connection.
	db, _ := NewDB()
	defer db.Close()

	// Create a test user.
	user := User{
		FirstName: "Jane",
		LastName:  "Unittester",
		Email:     "jane+unittest@options.cafe",
	}

	err := db.CreateNewUserWithStripe(user, "", "", "")

	// Since we are testing we know the user is id 4
	dbUser, _ := db.GetUserById(4)

	// Verify data returned
	st.Expect(t, err, nil)
	st.Expect(t, dbUser.Id, uint(4))
	st.Expect(t, dbUser.FirstName, "Jane")
	st.Expect(t, dbUser.LastName, "Unittester")
	st.Expect(t, dbUser.Email, "jane+unittest@options.cafe")
	st.Expect(t, len(dbUser.StripeCustomer), 18)
	st.Expect(t, len(dbUser.StripeSubscription), 18)

	// Get subscription with stripe
	sub, err := db.GetSubscriptionWithStripe(dbUser)

	// Verify data returned
	st.Expect(t, err, nil)
	st.Expect(t, sub.Name, "Monthly - $20")
	st.Expect(t, sub.Amount, 20.00)
	st.Expect(t, sub.TrialDays, 7)
	st.Expect(t, sub.BillingInterval, "month")
	st.Expect(t, sub.CardBrand, "")
	st.Expect(t, sub.CardLast4, "")
	st.Expect(t, sub.CardExpMonth, 0)
	st.Expect(t, sub.CardExpYear, 0)

	// Clean things up at stripe
	db.DeleteUserWithStripe(dbUser)
}

//
// Test - GetInvoiceHistoryWithStripe
//
func TestGetInvoiceHistoryWithStripe01(t *testing.T) {

	// // Load config file.
	// env.ReadEnv("../.env")

	// // Start the db connection.
	// db, _ := NewDB()
	// defer db.Close()

	// // Create a test user.
	// user := User{
	// 	FirstName:      "Jane",
	// 	LastName:       "Unittester",
	// 	Email:          "jane+unittest@options.cafe",
	// 	StripeCustomer: "cus_Djqq5q9mnW0lm6",
	// }

	// // Get invoices from stripe.
	// inv, err := db.GetInvoiceHistoryWithStripe(user)
	// st.Expect(t, err, nil)

	// spew.Dump(inv)
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
