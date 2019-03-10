//
// Date: 3/7/2019
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package reports

import (
	"testing"

	"github.com/nbio/st"

	"github.com/cloudmanic/app.options.cafe/backend/library/helpers"
	"github.com/cloudmanic/app.options.cafe/backend/library/test"
	"github.com/cloudmanic/app.options.cafe/backend/models"
)

//
// TestGetBalances01 - Load and return some test data. - Day
//
func TestGetBalances01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)
	// db, _, _ := models.NewTestDB("testing_db")

	// BrokerAccounts
	db.Create(&models.BrokerAccount{UserId: 1, BrokerId: 1, Name: "Test Account 1", AccountNumber: "abc1235423", StockCommission: 5.00, StockMin: 0.00, OptionCommission: 0.35, OptionSingleMin: 5.00, OptionMultiLegMin: 7.00, OptionBase: 0.00})
	db.Create(&models.BrokerAccount{UserId: 1, BrokerId: 1, Name: "Test Account 1", AccountNumber: "abc1235423", StockCommission: 5.00, StockMin: 0.00, OptionCommission: 0.35, OptionSingleMin: 5.00, OptionMultiLegMin: 7.00, OptionBase: 0.00})

	// Put test data into database - We want broker account id = 3 to match with the sql dump data
	brokerAccount := models.BrokerAccount{
		UserId:            1,
		BrokerId:          1,
		Name:              "Unit Test Account #1",
		AccountNumber:     "test12345",
		StockCommission:   5.00,
		StockMin:          0.00,
		OptionCommission:  0.35,
		OptionSingleMin:   5.00,
		OptionMultiLegMin: 7.00,
		OptionBase:        0.00,
	}

	// Insert into database.
	db.Create(&brokerAccount)

	// Load testing data.
	err := test.LoadSqlDump(db, "balance_histories_1")
	st.Expect(t, err, nil)

	// Call function we are testing.
	b := GetBalances(db, brokerAccount, BalancesParams{
		StartDate: helpers.ParseDateNoError("2019-01-01"),
		EndDate:   helpers.ParseDateNoError("2019-12-31"),
		GroupBy:   "day",
		Sort:      "asc",
	})

	// Test results
	st.Expect(t, len(b), 67)
	st.Expect(t, b[0].TotalCash, 2855.63)
	st.Expect(t, b[0].AccountValue, 1769.63)
	st.Expect(t, b[0].Date.Format("2006-01-02"), "2019-01-01")

	st.Expect(t, b[11].TotalCash, 1616.47)
	st.Expect(t, b[11].AccountValue, 2618.47)
	st.Expect(t, b[11].Date.Format("2006-01-02"), "2019-01-12")

	st.Expect(t, b[50].TotalCash, 4566.77)
	st.Expect(t, b[50].AccountValue, 4482.77)
	st.Expect(t, b[50].Date.Format("2006-01-02"), "2019-02-20")
}

//
// TestGetBalances02 - Load and return some test data. - Month
//
func TestGetBalances02(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)
	// db, _, _ := models.NewTestDB("testing_db")

	// BrokerAccounts
	db.Create(&models.BrokerAccount{UserId: 1, BrokerId: 1, Name: "Test Account 1", AccountNumber: "abc1235423", StockCommission: 5.00, StockMin: 0.00, OptionCommission: 0.35, OptionSingleMin: 5.00, OptionMultiLegMin: 7.00, OptionBase: 0.00})
	db.Create(&models.BrokerAccount{UserId: 1, BrokerId: 1, Name: "Test Account 1", AccountNumber: "abc1235423", StockCommission: 5.00, StockMin: 0.00, OptionCommission: 0.35, OptionSingleMin: 5.00, OptionMultiLegMin: 7.00, OptionBase: 0.00})

	// Put test data into database - We want broker account id = 3 to match with the sql dump data
	brokerAccount := models.BrokerAccount{
		UserId:            1,
		BrokerId:          1,
		Name:              "Unit Test Account #1",
		AccountNumber:     "test12345",
		StockCommission:   5.00,
		StockMin:          0.00,
		OptionCommission:  0.35,
		OptionSingleMin:   5.00,
		OptionMultiLegMin: 7.00,
		OptionBase:        0.00,
	}

	// Insert into database.
	db.Create(&brokerAccount)

	// Load testing data.
	err := test.LoadSqlDump(db, "balance_histories_1")
	st.Expect(t, err, nil)

	// Call function we are testing.
	b := GetBalances(db, brokerAccount, BalancesParams{
		StartDate: helpers.ParseDateNoError("2019-01-01"),
		EndDate:   helpers.ParseDateNoError("2019-12-31"),
		GroupBy:   "month",
		Sort:      "asc",
	})

	// Test results
	st.Expect(t, len(b), 3)
	st.Expect(t, b[0].TotalCash, 3711.97)
	st.Expect(t, b[0].AccountValue, 3656.97)
	st.Expect(t, b[0].Date.Format("2006-01-02"), "2019-01-31")

	st.Expect(t, b[1].TotalCash, 4732.85)
	st.Expect(t, b[1].AccountValue, 4528.35)
	st.Expect(t, b[1].Date.Format("2006-01-02"), "2019-02-28")

	st.Expect(t, b[2].TotalCash, 4907.36)
	st.Expect(t, b[2].AccountValue, 4284.36)
	st.Expect(t, b[2].Date.Format("2006-01-02"), "2019-03-31")
}

//
// TestGetBalances03 - Load and return some test data. - Year
//
func TestGetBalances03(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)
	// db, _, _ := models.NewTestDB("testing_db")

	// BrokerAccounts
	db.Create(&models.BrokerAccount{UserId: 1, BrokerId: 1, Name: "Test Account 1", AccountNumber: "abc1235423", StockCommission: 5.00, StockMin: 0.00, OptionCommission: 0.35, OptionSingleMin: 5.00, OptionMultiLegMin: 7.00, OptionBase: 0.00})
	db.Create(&models.BrokerAccount{UserId: 1, BrokerId: 1, Name: "Test Account 1", AccountNumber: "abc1235423", StockCommission: 5.00, StockMin: 0.00, OptionCommission: 0.35, OptionSingleMin: 5.00, OptionMultiLegMin: 7.00, OptionBase: 0.00})

	// Put test data into database - We want broker account id = 3 to match with the sql dump data
	brokerAccount := models.BrokerAccount{
		UserId:            1,
		BrokerId:          1,
		Name:              "Unit Test Account #1",
		AccountNumber:     "test12345",
		StockCommission:   5.00,
		StockMin:          0.00,
		OptionCommission:  0.35,
		OptionSingleMin:   5.00,
		OptionMultiLegMin: 7.00,
		OptionBase:        0.00,
	}

	// Insert into database.
	db.Create(&brokerAccount)

	// Load testing data.
	err := test.LoadSqlDump(db, "balance_histories_1")
	st.Expect(t, err, nil)

	// Call function we are testing.
	b := GetBalances(db, brokerAccount, BalancesParams{
		StartDate: helpers.ParseDateNoError("2017-01-01"),
		EndDate:   helpers.ParseDateNoError("2019-12-31"),
		GroupBy:   "year",
		Sort:      "asc",
	})

	// Test results
	st.Expect(t, len(b), 2)
	st.Expect(t, b[0].TotalCash, 2856.07)
	st.Expect(t, b[0].AccountValue, 1770.07)
	st.Expect(t, b[0].Date.Format("2006-01-02"), "2018-12-31")

	st.Expect(t, b[1].TotalCash, 4907.36)
	st.Expect(t, b[1].AccountValue, 4284.36)
	st.Expect(t, b[1].Date.Format("2006-01-02"), "2019-12-31")
}

//
// TestGetAccountReturns01 - Load and return some test data.
//
func TestGetAccountReturns01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// BrokerAccounts
	db.Create(&models.BrokerAccount{UserId: 1, BrokerId: 1, Name: "Test Account 1", AccountNumber: "abc1235423", StockCommission: 5.00, StockMin: 0.00, OptionCommission: 0.35, OptionSingleMin: 5.00, OptionMultiLegMin: 7.00, OptionBase: 0.00})
	db.Create(&models.BrokerAccount{UserId: 1, BrokerId: 1, Name: "Test Account 1", AccountNumber: "abc1235423", StockCommission: 5.00, StockMin: 0.00, OptionCommission: 0.35, OptionSingleMin: 5.00, OptionMultiLegMin: 7.00, OptionBase: 0.00})

	// Put test data into database - We want broker account id = 3 to match with the sql dump data
	brokerAccount := models.BrokerAccount{
		UserId:            1,
		BrokerId:          1,
		Name:              "Unit Test Account #1",
		AccountNumber:     "test12345",
		StockCommission:   5.00,
		StockMin:          0.00,
		OptionCommission:  0.35,
		OptionSingleMin:   5.00,
		OptionMultiLegMin: 7.00,
		OptionBase:        0.00,
	}

	// Insert into database.
	db.Create(&brokerAccount)

	// Load testing data.
	err := test.LoadSqlDump(db, "balance_histories_1")
	st.Expect(t, err, nil)
	err = test.LoadSqlDump(db, "broker_events_1")
	st.Expect(t, err, nil)

	// Call function we are testing.
	r := GetAccountReturns(db, brokerAccount, BalancesParams{
		StartDate: helpers.ParseDateNoError("2018-01-01"),
		EndDate:   helpers.ParseDateNoError("2018-12-31"),
	})

	// Test results
	st.Expect(t, len(r), 109)
	st.Expect(t, r[0].TotalCash, 4948.86)
	st.Expect(t, r[0].PricePer, 1.00)
	st.Expect(t, r[0].Percent, 0.00)
	st.Expect(t, r[0].Units, 4844.36)
	st.Expect(t, r[0].AccountValue, 4844.36)
	st.Expect(t, r[0].Date.Format("2006-01-02"), "2018-09-14")

	st.Expect(t, r[20].TotalCash, 5128.32)
	st.Expect(t, r[20].PricePer, 1.0166090051110983)
	st.Expect(t, r[20].Percent, 0.016609005111098307)
	st.Expect(t, r[20].Units, 4844.36)
	st.Expect(t, r[20].AccountValue, 4924.82)
	st.Expect(t, r[20].Date.Format("2006-01-02"), "2018-10-04")

	st.Expect(t, r[30].TotalCash, 7092.32)
	st.Expect(t, r[30].PricePer, 0.8793378557103995)
	st.Expect(t, r[30].Percent, -0.12066214428960054)
	st.Expect(t, r[30].Units, 6257.913228987943)
	st.Expect(t, r[30].AccountValue, 5502.82)
	st.Expect(t, r[30].Date.Format("2006-01-02"), "2018-10-14")

	st.Expect(t, r[108].TotalCash, 2856.07)
	st.Expect(t, r[108].PricePer, 0.16260125625509866)
	st.Expect(t, r[108].Percent, -0.8373987437449013)
	st.Expect(t, r[108].Units, 10885.955255001274)
	st.Expect(t, r[108].AccountValue, 1770.07)
	st.Expect(t, r[108].Date.Format("2006-01-02"), "2018-12-31")
}

/* End File */
