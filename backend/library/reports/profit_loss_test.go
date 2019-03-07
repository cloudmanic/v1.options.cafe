//
// Date: 3/6/2019
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package reports

import (
	"testing"

	"github.com/nbio/st"
	"github.com/optionscafe/options-cafe-cli/helpers"

	"github.com/cloudmanic/app.options.cafe/backend/library/test"
	"github.com/cloudmanic/app.options.cafe/backend/models"
)

//
// TestGetProfitLoss01 - Load and return some test data.
// Monthly test
//
func TestGetProfitLoss01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)
	// db, _, _ := models.NewTestDB("testing_db")

	// BrokerAccounts
	db.Create(&models.BrokerAccount{UserId: 1, BrokerId: 1, Name: "Test Account 1", AccountNumber: "abc1235423", StockCommission: 5.00, StockMin: 0.00, OptionCommission: 0.35, OptionSingleMin: 5.00, OptionMultiLegMin: 7.00, OptionBase: 0.00})

	// Put test data into database - We want broker account id = 2 to match with the sql dump data
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
	err := test.LoadSqlDump(db, "trade_groups_1")
	st.Expect(t, err, nil)

	// Call function we are testing.
	pl := GetProfitLoss(db, brokerAccount, ProfitLossParams{
		StartDate:  helpers.ParseDateNoError("2019-01-01"),
		EndDate:    helpers.ParseDateNoError("2019-12-31"),
		GroupBy:    "month",
		Sort:       "asc",
		Cumulative: false,
	})

	// Test results
	st.Expect(t, pl[0].DateStr, "2019-01")
	st.Expect(t, pl[0].Profit, -763.00)
	st.Expect(t, pl[0].TradeCount, 3)
	st.Expect(t, pl[0].Commissions, 41.00)
	st.Expect(t, pl[0].ProfitPerTrade, -254.33333333333334)
	st.Expect(t, pl[0].WinRatio, 66.66666666666666)
	st.Expect(t, pl[0].LossCount, 1)
	st.Expect(t, pl[0].WinCount, 2)

	st.Expect(t, pl[1].DateStr, "2019-02")
	st.Expect(t, pl[1].Profit, 144.00)
	st.Expect(t, pl[1].TradeCount, 5)
	st.Expect(t, pl[1].Commissions, 70.00)
	st.Expect(t, pl[1].ProfitPerTrade, 28.8)
	st.Expect(t, pl[1].WinRatio, 100.00)
	st.Expect(t, pl[1].LossCount, 0)
	st.Expect(t, pl[1].WinCount, 5)
}

//
// TestGetProfitLoss02 - Load and return some test data.
// Monthly test
//
func TestGetProfitLoss02(t *testing.T) {
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
	err := test.LoadSqlDump(db, "trade_groups_2")
	st.Expect(t, err, nil)

	// Call function we are testing.
	pl := GetProfitLoss(db, brokerAccount, ProfitLossParams{
		StartDate:  helpers.ParseDateNoError("2019-01-01"),
		EndDate:    helpers.ParseDateNoError("2019-12-31"),
		GroupBy:    "month",
		Sort:       "asc",
		Cumulative: false,
	})

	// Test results
	st.Expect(t, pl[0].DateStr, "2019-01")
	st.Expect(t, pl[0].Profit, -722.00)
	st.Expect(t, pl[0].TradeCount, 3)
	st.Expect(t, pl[0].Commissions, 0.00)
	st.Expect(t, pl[0].ProfitPerTrade, -240.66666666666666)
	st.Expect(t, pl[0].WinRatio, 66.66666666666666)
	st.Expect(t, pl[0].LossCount, 1)
	st.Expect(t, pl[0].WinCount, 2)

	st.Expect(t, pl[1].DateStr, "2019-02")
	st.Expect(t, pl[1].Profit, 214.00)
	st.Expect(t, pl[1].TradeCount, 5)
	st.Expect(t, pl[1].Commissions, 0.00)
	st.Expect(t, pl[1].ProfitPerTrade, 42.8)
	st.Expect(t, pl[1].WinRatio, 100.00)
	st.Expect(t, pl[1].LossCount, 0)
	st.Expect(t, pl[1].WinCount, 5)
}

//
// TestGetProfitLoss03 - Load and return some test data.
// Yearly test
//
func TestGetProfitLoss03(t *testing.T) {
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
	err := test.LoadSqlDump(db, "trade_groups_2")
	st.Expect(t, err, nil)

	// Call function we are testing.
	pl := GetProfitLoss(db, brokerAccount, ProfitLossParams{
		StartDate:  helpers.ParseDateNoError("2017-03-06"),
		EndDate:    helpers.ParseDateNoError("2019-03-06"),
		GroupBy:    "year",
		Sort:       "asc",
		Cumulative: false,
	})

	// Test results
	st.Expect(t, pl[0].DateStr, "2017")
	st.Expect(t, pl[0].Profit, 2780.00)
	st.Expect(t, pl[0].TradeCount, 23)
	st.Expect(t, pl[0].Commissions, 0.00)
	st.Expect(t, pl[0].ProfitPerTrade, 120.8695652173913)
	st.Expect(t, pl[0].WinRatio, 86.95652173913044)
	st.Expect(t, pl[0].LossCount, 3)
	st.Expect(t, pl[0].WinCount, 20)

	st.Expect(t, pl[1].DateStr, "2018")
	st.Expect(t, pl[1].Profit, -1126.80)
	st.Expect(t, pl[1].TradeCount, 60)
	st.Expect(t, pl[1].Commissions, 0.00)
	st.Expect(t, pl[1].ProfitPerTrade, -18.779999999999998)
	st.Expect(t, pl[1].WinRatio, 78.33333333333333)
	st.Expect(t, pl[1].LossCount, 13)
	st.Expect(t, pl[1].WinCount, 47)

	st.Expect(t, pl[2].DateStr, "2019")
	st.Expect(t, pl[2].Profit, -508.00)
	st.Expect(t, pl[2].TradeCount, 8)
	st.Expect(t, pl[2].Commissions, 0.00)
	st.Expect(t, pl[2].ProfitPerTrade, -63.50)
	st.Expect(t, pl[2].WinRatio, 87.5)
	st.Expect(t, pl[2].LossCount, 1)
	st.Expect(t, pl[2].WinCount, 7)
}

/* End File */
