//
// Date: 2/9/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package archive

import (
	"testing"
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/models"
	env "github.com/jpfuentes2/go-env"
	"github.com/nbio/st"
)

//
// Test reviewing current positions for expired options.
//
func TestReviewCurrentPositionsForExpiredOptions01(t *testing.T) {

	// Load config file.
	env.ReadEnv("../../.env")

	// Start the db connection.
	db, _ := models.NewDB()
	defer db.Close()

	// put test data into the DB.
	db.Exec("TRUNCATE TABLE positions;")
	ts := time.Date(2017, 10, 29, 17, 20, 01, 507451, time.UTC)
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
	st.Expect(t, len(results), 2)
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
// Test how we archive positions from orders.
//
func TestStorePositions01(t *testing.T) {

	// Load config file.
	env.ReadEnv("../../.env")

	// Start the db connection.
	db, _ := models.NewDB()
	defer db.Close()

	// Best to clear this data so we can test against empty tables.
	// This test is responsible for populating these tables.
	db.Exec("TRUNCATE TABLE trade_groups;")
	db.Exec("TRUNCATE TABLE positions;")

	// Set known values from testing data
	var userId uint = 1
	var brokerId uint = 2

	// Run test.
	err := StorePositions(db, userId, brokerId)

	// Verify the data was return as expected
	st.Expect(t, err, nil)

	// Query and get the trade groups
	var results = []models.TradeGroup{}

	// Run the query
	err = db.Query(&results, models.QueryParam{
		Order: "id",
		Sort:  "asc",
	})

	// Verify the data was return as expected
	st.Expect(t, err, nil)
	st.Expect(t, len(results), 13)

	//spew.Dump(results)

	// Check some of the single options
	st.Expect(t, results[10].Status, "Closed")
	st.Expect(t, results[10].Type, "Option")
	st.Expect(t, results[10].OrderIds, "1,2")
	st.Expect(t, results[10].Risked, 239.00)
	st.Expect(t, results[10].Profit, 113.00)
	st.Expect(t, results[10].Commission, 10.00)
	st.Expect(t, results[10].PercentGain, 47.28)

	st.Expect(t, results[11].Status, "Open")
	st.Expect(t, results[11].Type, "Option")
	st.Expect(t, results[11].OrderIds, "12")
	st.Expect(t, results[11].Risked, 1290.00)
	st.Expect(t, results[11].Profit, 0.00)
	st.Expect(t, results[11].Commission, 5.00)
	st.Expect(t, results[11].PercentGain, 0.00)

	// Spot check some put credit spreads
	st.Expect(t, results[1].Status, "Closed")
	st.Expect(t, results[1].Type, "Put Credit Spread")
	st.Expect(t, results[1].OrderIds, "4,6,7")
	st.Expect(t, results[1].Risked, 3195.00)
	st.Expect(t, results[1].Credit, 405.00)
	st.Expect(t, results[1].Profit, -863.60)
	st.Expect(t, results[1].Proceeds, -1242.00)
	st.Expect(t, results[1].Commission, 26.60)
	st.Expect(t, results[1].PercentGain, -27.03)

	// Spot check some call credit spreads
	st.Expect(t, results[6].Status, "Open")
	st.Expect(t, results[6].Type, "Call Credit Spread")
	st.Expect(t, results[6].OrderIds, "13")
	st.Expect(t, results[6].Risked, 712.00)
	st.Expect(t, results[6].Credit, 288.00)
	st.Expect(t, results[6].Profit, 0.00)
	st.Expect(t, results[6].Proceeds, 0.00)
	st.Expect(t, results[6].Commission, 7.00)
	st.Expect(t, results[6].PercentGain, 0.00)
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
