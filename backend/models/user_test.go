//
// Date: 2/10/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"os"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	env "github.com/jpfuentes2/go-env"
	"github.com/nbio/st"
	"github.com/optionscafe/options-cafe-cli/helpers"
)

//
// Test - Get all users
//
func TestGetAllUsers01(t *testing.T) {

	// Load config file.
	env.ReadEnv("../.env")

	// Start the db connection.
	db, dbName, _ := NewTestDB("")
	defer TestingTearDown(db, dbName)

	// Load test data
	loadUserTestingData(db.New())

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
	db, dbName, _ := NewTestDB("")
	defer TestingTearDown(db, dbName)

	// Load test data
	loadUserTestingData(db.New())

	// Create a test user.
	user := User{
		FirstName: "Jane",
		LastName:  "Unittester",
		Email:     "jane+unittest@options.cafe",
	}

	err := db.CreateNewUserWithStripe(user, os.Getenv("STRIPE_YEARLY_PLAN"), "tok_visa", "")

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
	db, dbName, _ := NewTestDB("")
	defer TestingTearDown(db, dbName)

	// Load test data
	loadUserTestingData(db.New())

	// Create a test user.
	user := User{
		FirstName: "Jane",
		LastName:  "Unittester",
		Email:     "jane+unittest@options.cafe",
	}

	err := db.CreateNewUserWithStripe(user, os.Getenv("STRIPE_MONTHLY_PLAN"), "tok_visa", "")

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
	st.Expect(t, sub.TrialDays, 30)
	st.Expect(t, sub.BillingInterval, "month")
	st.Expect(t, sub.CardBrand, "American Express")
	st.Expect(t, sub.CardLast4, "8431")
	st.Expect(t, sub.CardExpMonth, 3)
	st.Expect(t, sub.CardExpYear, 2020)

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
	db, dbName, _ := NewTestDB("")
	defer TestingTearDown(db, dbName)

	// Load test data
	loadUserTestingData(db.New())

	// Create a test user.
	user := User{
		FirstName: "Jane",
		LastName:  "Unittester",
		Email:     "jane+unittest@options.cafe",
	}

	err := db.CreateNewUserWithStripe(user, os.Getenv("STRIPE_MONTHLY_PLAN"), "tok_visa", "")

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
	st.Expect(t, sub.TrialDays, 30)
	st.Expect(t, sub.BillingInterval, "month")
	st.Expect(t, sub.CardBrand, "Visa")
	st.Expect(t, sub.CardLast4, "4242")
	st.Expect(t, sub.CardExpMonth, 3)
	st.Expect(t, sub.CardExpYear, 2020)

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
	db, dbName, _ := NewTestDB("")
	defer TestingTearDown(db, dbName)

	// Load test data
	loadUserTestingData(db.New())

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
	db, dbName, _ := NewTestDB("")
	defer TestingTearDown(db, dbName)

	// Load test data
	loadUserTestingData(db.New())

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
	db, dbName, _ := NewTestDB("")
	defer TestingTearDown(db, dbName)

	// Load test data
	loadUserTestingData(db.New())

	randString, _ := db.GenerateRandomString(10)

	if len(randString) != 10 {
		t.Errorf("The random string of %s was %d chars instead of 10.", randString, len(randString))
	}
}

//
// loadUserTestingData - Load testing data.
//
func loadUserTestingData(db *gorm.DB) {
	// Shared vars we use.
	ts := time.Date(2017, 10, 29, 17, 20, 01, 507451, time.UTC)

	// Users
	db.Create(&User{FirstName: "Rob", LastName: "Tester", Email: "spicer+robtester@options.cafe", Status: "Active"})
	db.Create(&User{FirstName: "Jane", LastName: "Wells", Email: "spicer+janewells@options.cafe", Status: "Active"})
	db.Create(&User{FirstName: "Bob", LastName: "Rosso", Email: "spicer+bobrosso@options.cafe", Status: "Active"})

	// Brokers
	db.Create(&Broker{Name: "Tradier", UserId: 1, AccessToken: "CLOwLO2cMnx-N_bPEexiVo9z9oRR80nPI9ycxQw3KQ-WQ4OP3D44gIbfLScAZ9pv", RefreshToken: "abc", TokenExpirationDate: ts})
	db.Create(&Broker{Name: "Tradeking", UserId: 1, AccessToken: "456", RefreshToken: "xyz", TokenExpirationDate: ts})
	db.Create(&Broker{Name: "Etrade", UserId: 1, AccessToken: "789", RefreshToken: "mno", TokenExpirationDate: ts})

	// BrokerAccounts
	db.Create(&BrokerAccount{UserId: 1, BrokerId: 1, Name: "Test Account 1", AccountNumber: "abc1235423", StockCommission: 5.00, StockMin: 0.00, OptionCommission: 0.35, OptionSingleMin: 5.00, OptionMultiLegMin: 7.00, OptionBase: 0.00})
	db.Create(&BrokerAccount{UserId: 1, BrokerId: 1, Name: "Test Account 2", AccountNumber: "ABC123ZY", StockCommission: 5.00, StockMin: 0.00, OptionCommission: 0.35, OptionSingleMin: 5.00, OptionMultiLegMin: 7.00, OptionBase: 0.00})

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

	// Add some watchlists
	db.Create(&Watchlist{UserId: 1, Name: "Watchlist #1 - User 1", Symbols: []WatchlistSymbol{{UserId: 1, SymbolId: 1, Order: 0}, {UserId: 1, SymbolId: 2, Order: 1}, {UserId: 1, SymbolId: 3, Order: 2}}})
	db.Create(&Watchlist{UserId: 1, Name: "Watchlist #2 - User 1", Symbols: []WatchlistSymbol{{UserId: 1, SymbolId: 1, Order: 0}, {UserId: 1, SymbolId: 2, Order: 1}, {UserId: 1, SymbolId: 3, Order: 2}}})
	db.Create(&Watchlist{UserId: 2, Name: "Watchlist #1 - User 2", Symbols: []WatchlistSymbol{{UserId: 2, SymbolId: 1, Order: 0}, {UserId: 2, SymbolId: 2, Order: 1}, {UserId: 2, SymbolId: 3, Order: 2}}})
	db.Create(&Watchlist{UserId: 2, Name: "Watchlist #2 - User 2", Symbols: []WatchlistSymbol{{UserId: 2, SymbolId: 1, Order: 0}, {UserId: 2, SymbolId: 2, Order: 1}, {UserId: 2, SymbolId: 3, Order: 2}}})

	// TradeGroups TODO: make this more complete
	db.Create(&TradeGroup{UserId: 1, BrokerAccountId: 1, BrokerAccountRef: "abc123", Status: "Open", Type: "Put Credit Spread", OrderIds: "1", Risked: 0.00, Proceeds: 0.00, Profit: 0.00, Commission: 23.45, Note: "Test note #1", OpenDate: ts})

	// Positions TODO: Put better values in here.
	db.Create(&Position{UserId: 1, TradeGroupId: 1, BrokerAccountId: 2, BrokerAccountRef: "123abc", Status: "Open", SymbolId: 4, Qty: 10, OrgQty: 10, CostBasis: 1000.00, AvgOpenPrice: 1.00, AvgClosePrice: 0.00, OrderIds: "1", OpenDate: ts, Note: "Test note #1"})
	db.Create(&Position{UserId: 1, TradeGroupId: 1, BrokerAccountId: 2, BrokerAccountRef: "123abc", Status: "Open", SymbolId: 6, Qty: 10, OrgQty: 10, CostBasis: 1000.00, AvgOpenPrice: 1.00, AvgClosePrice: 0.00, OrderIds: "1", OpenDate: ts, Note: "Test note #2"})
}

/* End File */
