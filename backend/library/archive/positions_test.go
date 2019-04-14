//
// Date: 2/9/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package archive

import (
	"go/build"
	"testing"
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/library/helpers"
	"github.com/cloudmanic/app.options.cafe/backend/library/test"
	"github.com/cloudmanic/app.options.cafe/backend/models"
	env "github.com/jpfuentes2/go-env"
	"github.com/nbio/st"
)

//
// Test reviewing current positions for expired options.
//
func TestReviewCurrentPositionsForExpiredOptions01(t *testing.T) {
	if testing.Short() {
		t.Skipf("Skipping test since --short was requested")
	}

	// Load .env file
	env.ReadEnv(build.Default.GOPATH + "/src/github.com/cloudmanic/app.options.cafe/backend/.env")

	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Shared vars we use.
	ts := time.Date(2017, 10, 29, 17, 20, 01, 507451, time.UTC)

	// Users
	db.Create(&models.User{FirstName: "Rob", LastName: "Tester", Email: "spicer+robtester@options.cafe", Status: "Active"})
	db.Create(&models.User{FirstName: "Jane", LastName: "Wells", Email: "spicer+janewells@options.cafe", Status: "Active"})
	db.Create(&models.User{FirstName: "Bob", LastName: "Rosso", Email: "spicer+bobrosso@options.cafe", Status: "Active"})

	// BrokerAccounts
	db.Create(&models.BrokerAccount{UserId: 1, BrokerId: 1, Name: "Test Account 1", AccountNumber: "abc1235423", StockCommission: 5.00, StockMin: 0.00, OptionCommission: 0.35, OptionSingleMin: 5.00, OptionMultiLegMin: 7.00, OptionBase: 0.00})
	db.Create(&models.BrokerAccount{UserId: 1, BrokerId: 1, Name: "Test Account 2", AccountNumber: "ABC123ZY", StockCommission: 5.00, StockMin: 0.00, OptionCommission: 0.35, OptionSingleMin: 5.00, OptionMultiLegMin: 7.00, OptionBase: 0.00})

	// Symbols
	db.Create(&models.Symbol{Name: "SPDR S&P 500 ETF Trust", ShortName: "SPY", Type: "Equity"})
	db.Create(&models.Symbol{Name: "McDonald's Corp", ShortName: "MCD", Type: "Equity"})
	db.Create(&models.Symbol{Name: "Starbucks Corp", ShortName: "SBUX", Type: "Equity"})
	db.Create(&models.Symbol{Name: "SPY Mar 16, 2018 $253.00 Put", ShortName: "SPY180316P00253000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 253.00})
	db.Create(&models.Symbol{Name: "VXX Feb 23, 2018 $55.00 Call", ShortName: "VXX180223C00055000", Type: "Option", OptionUnderlying: "VXX", OptionExpire: models.Date{helpers.ParseDateNoError("2/23/2018").UTC()}, OptionType: "Call", OptionStrike: 55.00})
	db.Create(&models.Symbol{Name: "SPY Mar 16, 2018 $266.00 Put", ShortName: "SPY180316P00266000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 266.00})
	db.Create(&models.Symbol{Name: "Caterpillar Inc", ShortName: "CAT", Type: "Equity"})
	db.Create(&models.Symbol{Name: "Ascent Solar Technologies Inc", ShortName: "ASTI", Type: "Equity"})
	db.Create(&models.Symbol{Name: "Advanced Micro Devices Inc", ShortName: "AMD", Type: "Equity"})
	db.Create(&models.Symbol{Name: "BARCLAYS BK PLC", ShortName: "VXX", Type: "Equity"})
	db.Create(&models.Symbol{Name: "SPY Feb 9, 2018 $276.00 Put", ShortName: "SPY180209P00276000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("2/9/2018").UTC()}, OptionType: "Put", OptionStrike: 276.00})
	db.Create(&models.Symbol{Name: "VXX Mar 2, 2018 $46.00 Put", ShortName: "VXX180302P00046000", Type: "Option", OptionUnderlying: "VXX", OptionExpire: models.Date{helpers.ParseDateNoError("3/2/2018").UTC()}, OptionType: "Put", OptionStrike: 46.00})
	db.Create(&models.Symbol{Name: "VXX Feb 23, 2018 $50.00 Put", ShortName: "VXX180223P00050000", Type: "Option", OptionUnderlying: "VXX", OptionExpire: models.Date{helpers.ParseDateNoError("2/23/2018").UTC()}, OptionType: "Put", OptionStrike: 50.00})
	db.Create(&models.Symbol{Name: "SPY Mar 16, 2018 $267.00 Put", ShortName: "SPY180316P00267000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 267.00})
	db.Create(&models.Symbol{Name: "SPY Mar 16, 2018 $269.00 Put", ShortName: "SPY180316P00269000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 269.00})
	db.Create(&models.Symbol{Name: "SPY Mar 16, 2018 $268.00 Put", ShortName: "SPY180316P00268000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 268.00})
	db.Create(&models.Symbol{Name: "SPY Mar 16, 2018 $270.00 Put", ShortName: "SPY180316P00270000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 270.00})
	db.Create(&models.Symbol{Name: "SPY Mar 16, 2018 $264.00 Put", ShortName: "SPY180316P00264000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 264.00})
	db.Create(&models.Symbol{Name: "SPY Mar 16, 2018 $266.00 Put", ShortName: "SPY180316P00266000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 266.00})
	db.Create(&models.Symbol{Name: "SPY Mar 9, 2018 $262.00 Put", ShortName: "SPY180309P00262000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("3/9/2018").UTC()}, OptionType: "Put", OptionStrike: 262.00})
	db.Create(&models.Symbol{Name: "SPY Mar 9, 2018 $264.00 Put", ShortName: "SPY180309P00264000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("3/9/2018").UTC()}, OptionType: "Put", OptionStrike: 264.00})
	db.Create(&models.Symbol{Name: "SPY Mar 16, 2018 $251.00 Put", ShortName: "SPY180316P00251000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 251.00})
	db.Create(&models.Symbol{Name: "SPY Mar 16, 2018 $253.00 Put", ShortName: "SPY180316P00253000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 253.00})
	db.Create(&models.Symbol{Name: "SPY Mar 16, 2018 $250.00 Put", ShortName: "SPY180316P00250000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 250.00})
	db.Create(&models.Symbol{Name: "SPY Mar 16, 2018 $252.00 Put", ShortName: "SPY180316P00252000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 252.00})
	db.Create(&models.Symbol{Name: "VXX Feb 23, 2018 $50.00 Call", ShortName: "VXX180223C00050000", Type: "Option", OptionUnderlying: "VXX", OptionExpire: models.Date{helpers.ParseDateNoError("2/23/2018").UTC()}, OptionType: "Call", OptionStrike: 50.00})
	db.Create(&models.Symbol{Name: "SPY Feb 28, 2018 $249.00 Put", ShortName: "SPY180228P00249000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("2/28/2018").UTC()}, OptionType: "Put", OptionStrike: 249.00})
	db.Create(&models.Symbol{Name: "SPY Feb 28, 2018 $251.00 Put", ShortName: "SPY180228P00251000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("2/28/2018").UTC()}, OptionType: "Put", OptionStrike: 251.00})
	db.Create(&models.Symbol{Name: "SPY Feb 21, 2018 $241.00 Put", ShortName: "SPY180221P00241000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("2/21/2018").UTC()}, OptionType: "Put", OptionStrike: 241.00})
	db.Create(&models.Symbol{Name: "SPY Feb 21, 2018 $243.00 Put", ShortName: "SPY180221P00243000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("2/21/2018").UTC()}, OptionType: "Put", OptionStrike: 243.00})
	db.Create(&models.Symbol{Name: "SPY Mar 21, 2018 $254.00 Put", ShortName: "SPY180321P00254000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("3/21/2018").UTC()}, OptionType: "Put", OptionStrike: 254.00})
	db.Create(&models.Symbol{Name: "SPY Mar 21, 2018 $256.00 Put", ShortName: "SPY180321P00256000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("3/21/2018").UTC()}, OptionType: "Put", OptionStrike: 256.00})
	db.Create(&models.Symbol{Name: "SPY Jul 21, 2017 $234.00 Put", ShortName: "SPY170721P00234000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("7/21/2018").UTC()}, OptionType: "Put", OptionStrike: 234.00})

	// put test data into the DB.
	db.Exec("TRUNCATE TABLE positions;")
	db.Create(&models.Position{UserId: 1, TradeGroupId: 1, BrokerAccountId: 1, BrokerAccountRef: "123abc", Status: "Closed", SymbolId: 4, Qty: 10, OrgQty: 10, CostBasis: 1000.00, AvgOpenPrice: 1.00, AvgClosePrice: 0.00, OrderIds: "1", OpenDate: ts, Note: "Test note #1"})
	db.Create(&models.Position{UserId: 1, TradeGroupId: 2, BrokerAccountId: 1, BrokerAccountRef: "123abc", Status: "Open", SymbolId: 3, Qty: 10, OrgQty: 10, CostBasis: 1000.00, AvgOpenPrice: 1.00, AvgClosePrice: 0.00, OrderIds: "1", OpenDate: ts, Note: "Test note #2"})
	db.Create(&models.Position{UserId: 1, TradeGroupId: 3, BrokerAccountId: 1, BrokerAccountRef: "123abc", Status: "Open", SymbolId: 33, Qty: 10, OrgQty: 10, CostBasis: -1000.00, AvgOpenPrice: 1.00, AvgClosePrice: 0.00, OrderIds: "1", OpenDate: ts, Note: "Test note #3"})
	db.Create(&models.Position{UserId: 1, TradeGroupId: 4, BrokerAccountId: 1, BrokerAccountRef: "123abc", Status: "Open", SymbolId: 6, Qty: 10, OrgQty: 10, CostBasis: 1000.00, AvgOpenPrice: 1.00, AvgClosePrice: 0.00, OrderIds: "1", OpenDate: ts, Note: "Test note #4"})

	db.Exec("TRUNCATE TABLE trade_groups;")
	db.Create(&models.TradeGroup{UserId: 1, BrokerAccountId: 1, BrokerAccountRef: "abc123", Status: "Open", Type: "Option", OrderIds: "1", Risked: 0.00, Proceeds: 0.00, Profit: 0.00, Commission: 23.45, Note: "Test note #1", OpenDate: ts})
	db.Create(&models.TradeGroup{UserId: 1, BrokerAccountId: 1, BrokerAccountRef: "abc123", Status: "Open", Type: "Equity", OrderIds: "1", Risked: 0.00, Proceeds: 0.00, Profit: 0.00, Commission: 23.45, Note: "Test note #2", OpenDate: ts})
	db.Create(&models.TradeGroup{UserId: 1, BrokerAccountId: 1, BrokerAccountRef: "abc123", Status: "Open", Type: "Option", OrderIds: "1", Risked: 300.00, Proceeds: 0.00, Profit: 0.00, Commission: 23.45, Note: "Test note #3", OpenDate: ts})
	db.Create(&models.TradeGroup{UserId: 1, BrokerAccountId: 1, BrokerAccountRef: "abc123", Status: "Open", Type: "Option", OrderIds: "1", Risked: 0.00, Proceeds: 0.00, Profit: 0.00, Commission: 23.45, Note: "Test note #4", OpenDate: ts})

	// Run the test function we are testing.
	err := ReviewCurrentPositionsForExpiredOptions(db, 1, 1)

	// Verify the data was return as expected
	st.Expect(t, err, nil)

	// Get all open positions.
	var results = []models.Position{}

	err = db.Query(&results, models.QueryParam{Wheres: []models.KeyValue{
		{Key: "status", Value: "Closed"},
		{Key: "user_id", ValueInt: 1},
		{Key: "broker_account_id", ValueInt: 1},
	}})

	// Verify the data was return as expected
	st.Expect(t, err, nil)
	st.Expect(t, len(results), 3)
	st.Expect(t, results[1].Id, uint(3))
	st.Expect(t, results[1].CostBasis, -1000.00)
	st.Expect(t, results[1].Profit, 1000.00)

	// Get all open positions.
	var result = models.TradeGroup{}

	err = db.Query(&result, models.QueryParam{Wheres: []models.KeyValue{
		{Key: "id", ValueInt: int(results[1].TradeGroupId)},
		{Key: "user_id", ValueInt: 1},
	}})

	// Verify the data was return as expected
	st.Expect(t, err, nil)
	st.Expect(t, result.Id, uint(3))
	st.Expect(t, result.Status, "Closed")
	st.Expect(t, result.Profit, 976.55)
	st.Expect(t, result.PercentGain, 333.33)
}

//
// Test reviewing current positions for expired options.
//
func TestReviewCurrentPositionsForExpiredOptions02(t *testing.T) {
	if testing.Short() {
		t.Skipf("Skipping test since --short was requested")
	}

	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Users
	db.Create(&models.User{FirstName: "Rob", LastName: "Tester", Email: "spicer+robtester@options.cafe", Status: "Active"})
	db.Create(&models.User{FirstName: "Jane", LastName: "Wells", Email: "spicer+janewells@options.cafe", Status: "Active"})
	db.Create(&models.User{FirstName: "Bob", LastName: "Rosso", Email: "spicer+bobrosso@options.cafe", Status: "Active"})

	// put test data into the DB.
	openTs := time.Date(2018, 11, 8, 17, 20, 01, 507451, time.UTC)

	db.Create(&models.Symbol{Name: "SPY Dec 21 2018 $260.00 Put", ShortName: "SPY181221P00260000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("12/21/2018").UTC()}, OptionType: "Put", OptionStrike: 260.00})
	db.Create(&models.Symbol{Name: "SPY Dec 21 2018 $262.00 Put", ShortName: "SPY181221P00262000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("12/21/2018").UTC()}, OptionType: "Put", OptionStrike: 262.00})

	db.Create(&models.TradeGroup{UserId: 1, BrokerAccountId: 1, Name: "Put Credit Spread", BrokerAccountRef: "abc123", Status: "Open", Type: "Option", OrderIds: "1", Risked: 1969.00, Proceeds: 0.00, Profit: 0.00, Commission: 0.00, OpenDate: openTs})

	db.Create(&models.Position{UserId: 1, TradeGroupId: 1, BrokerAccountId: 1, BrokerAccountRef: "123abc", Status: "Open", SymbolId: 1, Qty: 11, OrgQty: 11, CostBasis: 1738.00, AvgOpenPrice: 1.58, AvgClosePrice: 0.00, OrderIds: "1", OpenDate: openTs})
	db.Create(&models.Position{UserId: 1, TradeGroupId: 1, BrokerAccountId: 1, BrokerAccountRef: "123abc", Status: "Open", SymbolId: 2, Qty: -11, OrgQty: -11, CostBasis: -1969.00, AvgOpenPrice: 1.79, AvgClosePrice: 0.00, OrderIds: "1", OpenDate: openTs})

	// Run the test function we are testing.
	err := ReviewCurrentPositionsForExpiredOptions(db, 1, 1)

	// Verify the data was return as expected
	st.Expect(t, err, nil)

	// Get all open positions.
	var results = []models.Position{}

	err = db.Query(&results, models.QueryParam{Wheres: []models.KeyValue{
		{Key: "status", Value: "Closed"},
		{Key: "user_id", ValueInt: 1},
		{Key: "broker_account_id", ValueInt: 1},
	}})

	// Verify the data was return as expected
	st.Expect(t, err, nil)
	st.Expect(t, len(results), 2)
	st.Expect(t, results[1].Id, uint(2))
	st.Expect(t, results[1].CostBasis, -1969.00)
	st.Expect(t, results[1].Profit, -21461.00)

	// Get all open positions.
	var result = models.TradeGroup{}

	err = db.Query(&result, models.QueryParam{Wheres: []models.KeyValue{
		{Key: "id", ValueInt: int(results[1].TradeGroupId)},
		{Key: "user_id", ValueInt: 1},
	}})

	// Verify the data was return as expected
	st.Expect(t, err, nil)
	st.Expect(t, result.Id, uint(1))
	st.Expect(t, result.Status, "Closed")
	st.Expect(t, result.Profit, -1969.00)
	st.Expect(t, result.PercentGain, -100.00)
}

//
// Test how we archive positions from orders.
//
func TestStorePositions01(t *testing.T) {

	// TODO: Revisit to make work

	// // Load config file.
	// env.ReadEnv("../../.env")
	//
	// // Start the db connection.
	// db, dbName, _ := models.NewTestDB("")
	// defer models.TestingTearDown(db, dbName)
	//
	// // Shared vars we use.
	// ts := time.Date(2017, 10, 29, 17, 20, 01, 507451, time.UTC)
	//
	// // Users
	// db.Create(&models.User{FirstName: "Rob", LastName: "Tester", Email: "spicer+robtester@options.cafe", Status: "Active"})
	// db.Create(&models.User{FirstName: "Jane", LastName: "Wells", Email: "spicer+janewells@options.cafe", Status: "Active"})
	// db.Create(&models.User{FirstName: "Bob", LastName: "Rosso", Email: "spicer+bobrosso@options.cafe", Status: "Active"})
	//
	// // BrokerAccounts
	// db.Create(&models.BrokerAccount{UserId: 1, BrokerId: 1, Name: "Test Account 1", AccountNumber: "abc1235423", StockCommission: 5.00, StockMin: 0.00, OptionCommission: 0.35, OptionSingleMin: 5.00, OptionMultiLegMin: 7.00, OptionBase: 0.00})
	// db.Create(&models.BrokerAccount{UserId: 1, BrokerId: 1, Name: "Test Account 2", AccountNumber: "ABC123ZY", StockCommission: 5.00, StockMin: 0.00, OptionCommission: 0.35, OptionSingleMin: 5.00, OptionMultiLegMin: 7.00, OptionBase: 0.00})
	//
	// // Symbols
	// db.Create(&models.Symbol{Name: "SPDR S&P 500 ETF Trust", ShortName: "SPY", Type: "Equity"})
	// db.Create(&models.Symbol{Name: "McDonald's Corp", ShortName: "MCD", Type: "Equity"})
	// db.Create(&models.Symbol{Name: "Starbucks Corp", ShortName: "SBUX", Type: "Equity"})
	// db.Create(&models.Symbol{Name: "SPY Mar 16, 2018 $253.00 Put", ShortName: "SPY180316P00253000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 253.00})
	// db.Create(&models.Symbol{Name: "VXX Feb 23, 2018 $55.00 Call", ShortName: "VXX180223C00055000", Type: "Option", OptionUnderlying: "VXX", OptionExpire: models.Date{helpers.ParseDateNoError("2/23/2018").UTC()}, OptionType: "Call", OptionStrike: 55.00})
	// db.Create(&models.Symbol{Name: "SPY Mar 16, 2018 $266.00 Put", ShortName: "SPY180316P00266000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 266.00})
	// db.Create(&models.Symbol{Name: "Caterpillar Inc", ShortName: "CAT", Type: "Equity"})
	// db.Create(&models.Symbol{Name: "Ascent Solar Technologies Inc", ShortName: "ASTI", Type: "Equity"})
	// db.Create(&models.Symbol{Name: "Advanced Micro Devices Inc", ShortName: "AMD", Type: "Equity"})
	// db.Create(&models.Symbol{Name: "BARCLAYS BK PLC", ShortName: "VXX", Type: "Equity"})
	// db.Create(&models.Symbol{Name: "SPY Feb 9, 2018 $276.00 Put", ShortName: "SPY180209P00276000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("2/9/2018").UTC()}, OptionType: "Put", OptionStrike: 276.00})
	// db.Create(&models.Symbol{Name: "VXX Mar 2, 2018 $46.00 Put", ShortName: "VXX180302P00046000", Type: "Option", OptionUnderlying: "VXX", OptionExpire: models.Date{helpers.ParseDateNoError("3/2/2018").UTC()}, OptionType: "Put", OptionStrike: 46.00})
	// db.Create(&models.Symbol{Name: "VXX Feb 23, 2018 $50.00 Put", ShortName: "VXX180223P00050000", Type: "Option", OptionUnderlying: "VXX", OptionExpire: models.Date{helpers.ParseDateNoError("2/23/2018").UTC()}, OptionType: "Put", OptionStrike: 50.00})
	// db.Create(&models.Symbol{Name: "SPY Mar 16, 2018 $267.00 Put", ShortName: "SPY180316P00267000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 267.00})
	// db.Create(&models.Symbol{Name: "SPY Mar 16, 2018 $269.00 Put", ShortName: "SPY180316P00269000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 269.00})
	// db.Create(&models.Symbol{Name: "SPY Mar 16, 2018 $268.00 Put", ShortName: "SPY180316P00268000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 268.00})
	// db.Create(&models.Symbol{Name: "SPY Mar 16, 2018 $270.00 Put", ShortName: "SPY180316P00270000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 270.00})
	// db.Create(&models.Symbol{Name: "SPY Mar 16, 2018 $264.00 Put", ShortName: "SPY180316P00264000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 264.00})
	// db.Create(&models.Symbol{Name: "SPY Mar 16, 2018 $266.00 Put", ShortName: "SPY180316P00266000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 266.00})
	// db.Create(&models.Symbol{Name: "SPY Mar 9, 2018 $262.00 Put", ShortName: "SPY180309P00262000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("3/9/2018").UTC()}, OptionType: "Put", OptionStrike: 262.00})
	// db.Create(&models.Symbol{Name: "SPY Mar 9, 2018 $264.00 Put", ShortName: "SPY180309P00264000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("3/9/2018").UTC()}, OptionType: "Put", OptionStrike: 264.00})
	// db.Create(&models.Symbol{Name: "SPY Mar 16, 2018 $251.00 Put", ShortName: "SPY180316P00251000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 251.00})
	// db.Create(&models.Symbol{Name: "SPY Mar 16, 2018 $253.00 Put", ShortName: "SPY180316P00253000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 253.00})
	// db.Create(&models.Symbol{Name: "SPY Mar 16, 2018 $250.00 Put", ShortName: "SPY180316P00250000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 250.00})
	// db.Create(&models.Symbol{Name: "SPY Mar 16, 2018 $252.00 Put", ShortName: "SPY180316P00252000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 252.00})
	// db.Create(&models.Symbol{Name: "VXX Feb 23, 2018 $50.00 Call", ShortName: "VXX180223C00050000", Type: "Option", OptionUnderlying: "VXX", OptionExpire: models.Date{helpers.ParseDateNoError("2/23/2018").UTC()}, OptionType: "Call", OptionStrike: 50.00})
	// db.Create(&models.Symbol{Name: "SPY Feb 28, 2018 $249.00 Put", ShortName: "SPY180228P00249000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("2/28/2018").UTC()}, OptionType: "Put", OptionStrike: 249.00})
	// db.Create(&models.Symbol{Name: "SPY Feb 28, 2018 $251.00 Put", ShortName: "SPY180228P00251000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("2/28/2018").UTC()}, OptionType: "Put", OptionStrike: 251.00})
	// db.Create(&models.Symbol{Name: "SPY Feb 21, 2018 $241.00 Put", ShortName: "SPY180221P00241000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("2/21/2018").UTC()}, OptionType: "Put", OptionStrike: 241.00})
	// db.Create(&models.Symbol{Name: "SPY Feb 21, 2018 $243.00 Put", ShortName: "SPY180221P00243000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("2/21/2018").UTC()}, OptionType: "Put", OptionStrike: 243.00})
	// db.Create(&models.Symbol{Name: "SPY Mar 21, 2018 $254.00 Put", ShortName: "SPY180321P00254000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("3/21/2018").UTC()}, OptionType: "Put", OptionStrike: 254.00})
	// db.Create(&models.Symbol{Name: "SPY Mar 21, 2018 $256.00 Put", ShortName: "SPY180321P00256000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("3/21/2018").UTC()}, OptionType: "Put", OptionStrike: 256.00})
	// db.Create(&models.Symbol{Name: "SPY Jul 21, 2017 $234.00 Put", ShortName: "SPY170721P00234000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("7/21/2018").UTC()}, OptionType: "Put", OptionStrike: 234.00})
	//
	// // TradeGroups TODO: make this more complete
	// db.Create(&models.TradeGroup{UserId: 1, BrokerAccountId: 1, BrokerAccountRef: "abc123", Status: "Open", Type: "Put Credit Spread", OrderIds: "1", Risked: 0.00, Proceeds: 0.00, Profit: 0.00, Commission: 23.45, Note: "Test note #1", OpenDate: ts})
	//
	// // Positions TODO: Put better values in here.
	// db.Create(&models.Position{UserId: 1, TradeGroupId: 1, BrokerAccountId: 2, BrokerAccountRef: "123abc", Status: "Open", SymbolId: 4, Qty: 10, OrgQty: 10, CostBasis: 1000.00, AvgOpenPrice: 1.00, AvgClosePrice: 0.00, OrderIds: "1", OpenDate: ts, Note: "Test note #1"})
	// db.Create(&models.Position{UserId: 1, TradeGroupId: 1, BrokerAccountId: 2, BrokerAccountRef: "123abc", Status: "Open", SymbolId: 6, Qty: 10, OrgQty: 10, CostBasis: 1000.00, AvgOpenPrice: 1.00, AvgClosePrice: 0.00, OrderIds: "1", OpenDate: ts, Note: "Test note #2"})
	//
	// // Orders
	// db.Exec(`INSERT INTO orders (id, user_id, created_at, updated_at, broker_account_id, broker_ref, broker_account_ref, type, symbol_id, option_symbol_id, side, qty, status, duration, price, avg_fill_price, exec_quantity, last_fill_price, last_fill_quantity, remaining_quantity, create_date, transaction_date, class, num_legs, position_reviewed) VALUES
	// (1,1,'2018-02-17 00:13:53','2018-02-17 00:13:54',2,'734801','ABC123ZY','limit',1,11,'buy_to_open',1,'filled','day',2.40,2.39,1.00,2.39,1.00,0.00,'2018-01-16 11:54:50','2018-01-16 11:54:51','option',0,'No'),
	// (2,1,'2018-02-17 00:13:53','2018-02-17 00:13:54',2,'735196','ABC123ZY','limit',1,11,'sell_to_close',-1,'filled','gtc',3.39,3.62,1.00,3.62,1.00,0.00,'2018-01-16 15:29:51','2018-02-05 06:30:03','option',0,'No'),
	// (3,1,'2018-02-17 00:13:53','2018-02-17 00:13:53',2,'767256','ABC123ZY','credit',1,0,'buy',9,'filled','day',0.24,-0.24,9.00,0.00,1.00,0.00,'2018-01-30 09:20:25','2018-01-30 09:45:11','multileg',2,'No'),
	// (4,1,'2018-02-17 00:13:53','2018-02-17 00:13:53',2,'770154','ABC123ZY','credit',1,0,'buy',9,'filled','day',0.24,-0.24,9.00,0.00,9.00,0.00,'2018-01-31 09:39:20','2018-01-31 11:21:38','multileg',2,'No'),
	// (5,1,'2018-02-17 00:13:53','2018-02-17 00:13:53',2,'772977','ABC123ZY','debit',1,0,'buy',9,'filled','gtc',0.80,0.72,9.00,0.00,9.00,0.00,'2018-02-01 09:30:04','2018-02-05 11:57:31','multileg',2,'No'),
	// (6,1,'2018-02-17 00:13:53','2018-02-17 00:13:53',2,'773020','ABC123ZY','credit',1,0,'buy',9,'filled','day',0.21,-0.21,9.00,0.00,9.00,0.00,'2018-02-01 09:40:44','2018-02-01 10:42:33','multileg',2,'No'),
	// (7,1,'2018-02-17 00:13:53','2018-02-17 00:13:53',2,'775404','ABC123ZY','debit',1,0,'buy',18,'filled','gtc',0.80,0.69,18.00,0.00,18.00,0.00,'2018-02-02 07:27:02','2018-02-05 11:58:57','multileg',2,'No'),
	// (8,1,'2018-02-17 00:13:53','2018-02-17 00:13:53',2,'775734','ABC123ZY','credit',1,0,'buy',9,'filled','day',0.23,-0.23,9.00,0.00,9.00,0.00,'2018-02-02 08:07:00','2018-02-02 08:19:15','multileg',2,'No'),
	// (9,1,'2018-02-17 00:13:53','2018-02-17 00:13:53',2,'780197','ABC123ZY','credit',1,0,'buy',9,'filled','day',0.25,-0.25,9.00,0.00,9.00,0.00,'2018-02-05 08:01:06','2018-02-05 08:17:18','multileg',2,'No'),
	// (10,1,'2018-02-17 00:13:53','2018-02-17 00:13:53',2,'781119','ABC123ZY','credit',1,0,'buy',9,'filled','day',0.22,-0.22,9.00,0.00,9.00,0.00,'2018-02-05 10:40:10','2018-02-05 11:26:06','multileg',2,'No'),
	// (11,1,'2018-02-17 00:13:53','2018-02-17 00:13:54',2,'781816','ABC123ZY','credit',1,0,'buy',9,'filled','day',0.35,-0.35,9.00,0.00,1.00,0.00,'2018-02-05 12:11:47','2018-02-05 12:14:14','multileg',2,'No'),
	// (12,1,'2018-02-17 00:13:53','2018-02-17 00:13:54',2,'784720','ABC123ZY','limit',10,12,'buy_to_open',2,'filled','day',6.60,6.45,2.00,6.45,1.00,0.00,'2018-02-06 09:36:01','2018-02-06 09:36:58','option',0,'No'),
	// (13,1,'2018-02-17 00:13:53','2018-02-17 00:13:54',2,'784730','ABC123ZY','credit',10,0,'buy',2,'filled','day',1.20,-1.44,2.00,0.00,2.00,0.00,'2018-02-06 09:39:06','2018-02-06 09:39:06','multileg',2,'No'),
	// (14,1,'2018-02-17 00:13:53','2018-02-17 00:13:54',2,'790726','ABC123ZY','credit',1,0,'buy',9,'filled','day',0.31,-0.31,9.00,0.00,9.00,0.00,'2018-02-08 09:53:10','2018-02-08 09:54:24','multileg',2,'No'),
	// (15,1,'2018-02-17 00:13:53','2018-02-17 00:13:54',2,'790833','ABC123ZY','debit',1,0,'buy',9,'filled','gtc',0.03,0.03,9.00,0.00,9.00,0.00,'2018-02-08 10:09:37','2018-02-16 09:35:13','multileg',2,'No'),
	// (16,1,'2018-02-17 00:13:53','2018-02-17 00:13:54',2,'791126','ABC123ZY','limit',10,13,'buy_to_open',2,'filled','day',5.00,5.00,2.00,5.00,2.00,0.00,'2018-02-08 11:05:07','2018-02-08 12:35:55','option',0,'No'),
	// (17,1,'2018-02-17 00:13:53','2018-02-17 00:13:54',2,'793756','ABC123ZY','credit',1,0,'buy',9,'filled','day',0.28,-0.28,9.00,0.00,9.00,0.00,'2018-02-09 09:06:50','2018-02-09 09:08:48','multileg',2,'No'),
	// (18,1,'2018-02-17 00:13:53','2018-02-17 00:13:54',2,'793804','ABC123ZY','limit',10,13,'sell_to_close',-2,'filled','gtc',9.00,9.00,2.00,9.00,2.00,0.00,'2018-02-09 09:11:52','2018-02-16 10:35:53','option',0,'No'),
	// (19,1,'2018-02-17 00:13:53','2018-02-17 00:13:54',2,'797557','ABC123ZY','debit',1,0,'buy',9,'filled','gtc',0.03,0.03,9.00,0.00,9.00,0.00,'2018-02-12 09:09:44','2018-02-14 06:35:46','multileg',2,'No'),
	// (20,1,'2018-02-17 00:13:53','2018-02-17 00:13:54',2,'802693','ABC123ZY','credit',1,0,'buy',9,'filled','day',0.29,-0.29,9.00,0.00,9.00,0.00,'2018-02-14 09:03:46','2018-02-14 09:10:39','multileg',2,'No');`)
	//
	// // OrderLegs
	// db.Exec(`INSERT INTO order_legs (id, user_id, created_at, updated_at, order_id, type, symbol_id, side, qty, status, duration, avg_fill_price, exec_quantity, last_fill_price, last_fill_quantity, remaining_quantity, create_date, transaction_date) VALUES
	// (1,1,'2018-02-17 00:13:53','2018-02-17 00:13:53',3,'credit',14,'buy_to_open',9,'filled','day',1.62,9.00,1.65,1.00,0.00,'2018-01-30 09:20:25','2018-01-30 09:45:11'),
	// (2,1,'2018-02-17 00:13:53','2018-02-17 00:13:53',3,'credit',15,'sell_to_open',9,'filled','day',1.86,9.00,1.89,1.00,0.00,'2018-01-30 09:20:25','2018-01-30 09:45:11'),
	// (3,1,'2018-02-17 00:13:53','2018-02-17 00:13:53',4,'credit',16,'buy_to_open',9,'filled','day',1.49,9.00,1.49,3.00,0.00,'2018-01-31 09:39:20','2018-01-31 11:21:38'),
	// (4,1,'2018-02-17 00:13:53','2018-02-17 00:13:53',4,'credit',17,'sell_to_open',9,'filled','day',1.73,9.00,1.73,9.00,0.00,'2018-01-31 09:39:20','2018-01-31 11:21:38'),
	// (5,1,'2018-02-17 00:13:53','2018-02-17 00:13:53',5,'debit',14,'sell_to_close',9,'filled','gtc',6.27,9.00,6.27,9.00,0.00,'2018-02-01 09:30:04','2018-02-05 11:57:31'),
	// (6,1,'2018-02-17 00:13:53','2018-02-17 00:13:53',5,'debit',15,'buy_to_close',9,'filled','gtc',6.99,9.00,6.99,9.00,0.00,'2018-02-01 09:30:04','2018-02-05 11:57:31'),
	// (7,1,'2018-02-17 00:13:53','2018-02-17 00:13:53',6,'credit',16,'buy_to_open',9,'filled','day',1.22,9.00,1.22,1.00,0.00,'2018-02-01 09:40:44','2018-02-01 10:42:33'),
	// (8,1,'2018-02-17 00:13:53','2018-02-17 00:13:53',6,'credit',17,'sell_to_open',9,'filled','day',1.43,9.00,1.43,9.00,0.00,'2018-02-01 09:40:44','2018-02-01 10:42:33'),
	// (9,1,'2018-02-17 00:13:53','2018-02-17 00:13:53',7,'debit',16,'sell_to_close',18,'filled','gtc',6.82,18.00,6.82,18.00,0.00,'2018-02-02 07:27:02','2018-02-05 11:58:57'),
	// (10,1,'2018-02-17 00:13:53','2018-02-17 00:13:53',7,'debit',17,'buy_to_close',18,'filled','gtc',7.51,18.00,7.51,18.00,0.00,'2018-02-02 07:27:02','2018-02-05 11:58:57'),
	// (11,1,'2018-02-17 00:13:53','2018-02-17 00:13:53',8,'credit',18,'buy_to_open',9,'filled','day',1.48,9.00,1.48,9.00,0.00,'2018-02-02 08:07:00','2018-02-02 08:19:15'),
	// (12,1,'2018-02-17 00:13:53','2018-02-17 00:13:53',8,'credit',19,'sell_to_open',9,'filled','day',1.71,9.00,1.71,9.00,0.00,'2018-02-02 08:07:00','2018-02-02 08:19:15'),
	// (13,1,'2018-02-17 00:13:53','2018-02-17 00:13:53',9,'credit',20,'buy_to_open',9,'filled','day',1.56,9.00,1.56,9.00,0.00,'2018-02-05 08:01:06','2018-02-05 08:17:18'),
	// (14,1,'2018-02-17 00:13:53','2018-02-17 00:13:53',9,'credit',21,'sell_to_open',9,'filled','day',1.81,9.00,1.81,9.00,0.00,'2018-02-05 08:01:06','2018-02-05 08:17:18'),
	// (15,1,'2018-02-17 00:13:53','2018-02-17 00:13:53',10,'credit',22,'buy_to_open',9,'filled','day',2.33,9.00,2.33,9.00,0.00,'2018-02-05 10:40:10','2018-02-05 11:26:06'),
	// (16,1,'2018-02-17 00:13:53','2018-02-17 00:13:53',10,'credit',23,'sell_to_open',9,'filled','day',2.55,9.00,2.55,9.00,0.00,'2018-02-05 10:40:10','2018-02-05 11:26:06'),
	// (17,1,'2018-02-17 00:13:53','2018-02-17 00:13:54',11,'credit',24,'buy_to_open',9,'filled','day',2.56,9.00,4.55,1.00,0.00,'2018-02-05 12:11:47','2018-02-05 12:14:14'),
	// (18,1,'2018-02-17 00:13:53','2018-02-17 00:13:54',11,'credit',25,'sell_to_open',9,'filled','day',2.91,9.00,4.90,1.00,0.00,'2018-02-05 12:11:47','2018-02-05 12:14:14'),
	// (19,1,'2018-02-17 00:13:53','2018-02-17 00:13:54',13,'credit',26,'sell_to_open',2,'filled','day',6.35,2.00,6.35,2.00,0.00,'2018-02-06 09:39:06','2018-02-06 09:39:06'),
	// (20,1,'2018-02-17 00:13:53','2018-02-17 00:13:54',13,'credit',5,'buy_to_open',2,'filled','day',4.91,2.00,4.91,2.00,0.00,'2018-02-06 09:39:06','2018-02-06 09:39:06'),
	// (21,1,'2018-02-17 00:13:53','2018-02-17 00:13:54',14,'credit',27,'buy_to_open',9,'filled','day',3.28,9.00,3.28,9.00,0.00,'2018-02-08 09:53:10','2018-02-08 09:54:24'),
	// (22,1,'2018-02-17 00:13:53','2018-02-17 00:13:54',14,'credit',28,'sell_to_open',9,'filled','day',3.59,9.00,3.59,9.00,0.00,'2018-02-08 09:53:10','2018-02-08 09:54:24'),
	// (23,1,'2018-02-17 00:13:53','2018-02-17 00:13:54',15,'debit',27,'sell_to_close',9,'filled','gtc',0.14,9.00,0.14,4.00,0.00,'2018-02-08 10:09:37','2018-02-16 09:35:13'),
	// (24,1,'2018-02-17 00:13:53','2018-02-17 00:13:54',15,'debit',28,'buy_to_close',9,'filled','gtc',0.17,9.00,0.17,9.00,0.00,'2018-02-08 10:09:37','2018-02-16 09:35:13'),
	// (25,1,'2018-02-17 00:13:53','2018-02-17 00:13:54',17,'credit',29,'buy_to_open',9,'filled','day',2.26,9.00,2.26,9.00,0.00,'2018-02-09 09:06:50','2018-02-09 09:08:48'),
	// (26,1,'2018-02-17 00:13:53','2018-02-17 00:13:54',17,'credit',30,'sell_to_open',9,'filled','day',2.54,9.00,2.54,9.00,0.00,'2018-02-09 09:06:50','2018-02-09 09:08:48'),
	// (27,1,'2018-02-17 00:13:53','2018-02-17 00:13:54',19,'debit',29,'sell_to_close',9,'filled','gtc',0.10,9.00,0.10,9.00,0.00,'2018-02-12 09:09:44','2018-02-14 06:35:46'),
	// (28,1,'2018-02-17 00:13:53','2018-02-17 00:13:54',19,'debit',30,'buy_to_close',9,'filled','gtc',0.13,9.00,0.13,9.00,0.00,'2018-02-12 09:09:44','2018-02-14 06:35:46'),
	// (29,1,'2018-02-17 00:13:53','2018-02-17 00:13:54',20,'credit',31,'buy_to_open',9,'filled','day',2.17,9.00,2.17,9.00,0.00,'2018-02-14 09:03:46','2018-02-14 09:10:39'),
	// (30,1,'2018-02-17 00:13:53','2018-02-17 00:13:54',20,'credit',32,'sell_to_open',9,'filled','day',2.46,9.00,2.46,9.00,0.00,'2018-02-14 09:03:46','2018-02-14 09:10:39');`)
	//
	// // Set known values from testing data
	// var userId uint = 1
	// var brokerId uint = 2
	//
	// // Run test.
	// err := StorePositions(db, userId, brokerId)
	//
	// // Verify the data was return as expected
	// st.Expect(t, err, nil)
	//
	// // Query and get the trade groups
	// var results = []models.TradeGroup{}
	//
	// // Run the query
	// err = db.Query(&results, models.QueryParam{
	// 	Order: "id",
	// 	Sort:  "asc",
	// })
	//
	// // Verify the data was return as expected
	// st.Expect(t, err, nil)
	// st.Expect(t, len(results), 13)
	//
	// // Check some of the single options
	// st.Expect(t, results[10].Status, "Closed")
	// st.Expect(t, results[10].Type, "Option")
	// st.Expect(t, results[10].OrderIds, "1,2")
	// st.Expect(t, results[10].Risked, 239.00)
	// st.Expect(t, results[10].Profit, 113.00)
	// st.Expect(t, results[10].Commission, 10.00)
	// st.Expect(t, results[10].PercentGain, 47.28)
	//
	// st.Expect(t, results[11].Status, "Closed")
	// st.Expect(t, results[11].Type, "Option")
	// st.Expect(t, results[11].OrderIds, "12")
	// st.Expect(t, results[11].Risked, 1290.00)
	// st.Expect(t, results[11].Profit, -1295.00)
	// st.Expect(t, results[11].Commission, 5.00)
	// st.Expect(t, results[11].PercentGain, -100.00)
	//
	// // Spot check some put credit spreads
	// st.Expect(t, results[1].Status, "Closed")
	// st.Expect(t, results[1].Type, "Put Credit Spread")
	// st.Expect(t, results[1].OrderIds, "4,6,7")
	// st.Expect(t, results[1].Risked, 3195.00)
	// st.Expect(t, results[1].Credit, 405.00)
	// st.Expect(t, results[1].Profit, -863.60)
	// st.Expect(t, results[1].Proceeds, -1242.00)
	// st.Expect(t, results[1].Commission, 26.60)
	// st.Expect(t, results[1].PercentGain, -27.03)
	//
	// // Spot check some call credit spreads
	// st.Expect(t, results[6].Status, "Closed")
	// st.Expect(t, results[6].Type, "Call Credit Spread")
	// st.Expect(t, results[6].OrderIds, "13")
	// st.Expect(t, results[6].Risked, 712.00)
	// st.Expect(t, results[6].Credit, 288.00)
	// st.Expect(t, results[6].Profit, 281.00)
	// st.Expect(t, results[6].Proceeds, 0.00)
	// st.Expect(t, results[6].Commission, 7.00)
	// st.Expect(t, results[6].PercentGain, 40.45)
}

//
// TestStorePositions02 - chadwick.mchugh@zoho.com Account - 4/13/19
//
func TestStorePositions02(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("testing_db")
	defer models.TestingTearDown(db, dbName)

	// Set vars
	var userId uint = 56
	var brokerId uint = 1

	// Users
	db.Create(&models.User{Id: userId, FirstName: "Rob", LastName: "Tester", Email: "spicer+robtester@options.cafe", Status: "Active"})

	// Brokers
	db.Create(&models.Broker{Id: brokerId, Name: "Tradier", UserId: userId})

	// Put test data into database
	brokerAccount := models.BrokerAccount{
		Id:                35,
		UserId:            userId,
		BrokerId:          brokerId,
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
	err := test.LoadSqlDump(db, "orders_1")
	st.Expect(t, err, nil)

	err = test.LoadSqlDump(db, "orders_legs_1")
	st.Expect(t, err, nil)

	err = test.LoadSqlDump(db, "symbols_1")
	st.Expect(t, err, nil)

	err = test.LoadSqlDump(db, "historical_quotes_1")
	st.Expect(t, err, nil)

	// Store orders
	err = StorePositions(db, userId, brokerId)
	st.Expect(t, err, nil)

	// Get all the trade groups and verify
	tgs := []models.TradeGroup{}
	db.New().Order("id ASC").Find(&tgs)

	// Verify Stored data (spot check)
	st.Expect(t, len(tgs), 15)
	st.Expect(t, tgs[5].Name, "Trade #6 - Long Put Butterfly Trade")
	st.Expect(t, tgs[6].Type, "Long Call Butterfly")
	st.Expect(t, tgs[9].OrderIds, "2298,2299")
	st.Expect(t, tgs[9].Note, " Trade Expired.")
	st.Expect(t, tgs[10].Note, " Trade Expired.")

	// Get all the positions and verify
	ps := []models.Position{}
	db.New().Order("id ASC").Find(&ps)

	// Verify Stored data (spot check)
	st.Expect(t, len(ps), 38)
	st.Expect(t, ps[6].CostBasis, 982.00)
	st.Expect(t, ps[6].TradeGroupId, uint(2))
}

//
// Test Commission for order.
//
func TestCalcCommissionForOrder01(t *testing.T) {
	// Build sample order.
	order := models.Order{
		Class: "multileg",
		Legs: []models.OrderLeg{
			{Qty: 10},
			{Qty: 10},
		},
	}

	// Build sample BrokerAccount
	brokerAccount := models.BrokerAccount{
		StockCommission:   5.00,
		StockMin:          0.00,
		OptionCommission:  0.35,
		OptionSingleMin:   5.00,
		OptionMultiLegMin: 7.00,
		OptionBase:        0.00,
	}

	// Figure out Commission
	commission := CalcCommissionForOrder(&order, 1, &brokerAccount)

	// Verify the data was return as expected
	st.Expect(t, commission, 7.00)
}

//
// Test Commission for order.
//
func TestCalcCommissionForOrder02(t *testing.T) {
	// Build sample order.
	order := models.Order{
		Class: "multileg",
		Legs: []models.OrderLeg{
			{Qty: 15},
			{Qty: 10},
		},
	}

	// Build sample BrokerAccount
	brokerAccount := models.BrokerAccount{
		StockCommission:   5.00,
		StockMin:          0.00,
		OptionCommission:  0.35,
		OptionSingleMin:   5.00,
		OptionMultiLegMin: 7.00,
		OptionBase:        0.00,
	}

	// Figure out Commission
	commission := CalcCommissionForOrder(&order, 1, &brokerAccount)

	// Verify the data was return as expected
	st.Expect(t, commission, 8.75)
}

/* End File */
