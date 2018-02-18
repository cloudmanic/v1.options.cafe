//
// Date: 2/9/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package archive

import (
	"testing"

	"github.com/cloudmanic/app.options.cafe/backend/models"
	env "github.com/jpfuentes2/go-env"
	"github.com/nbio/st"
)

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
	err = db.Query(&results, models.QueryParam{})

	// Verify the data was return as expected
	st.Expect(t, err, nil)
	st.Expect(t, len(results), 13)

	// Check some of the single options
	st.Expect(t, results[10].Status, "Closed")
	st.Expect(t, results[10].Type, "Option")
	st.Expect(t, results[10].OrderIds, "1,2")
	st.Expect(t, results[10].Risked, 239.00)
	st.Expect(t, results[10].Profit, 113.00)
	st.Expect(t, results[10].Commission, 10.00)

	st.Expect(t, results[11].Status, "Open")
	st.Expect(t, results[11].Type, "Option")
	st.Expect(t, results[11].OrderIds, "12")
	st.Expect(t, results[11].Risked, 1290.00)
	st.Expect(t, results[11].Profit, 0.00)
	st.Expect(t, results[11].Commission, 5.00)

	// Spot check some put credit spreads
	st.Expect(t, results[1].Status, "Closed")
	st.Expect(t, results[1].Type, "Put Credit Spread")
	st.Expect(t, results[1].OrderIds, "4,6,7")
	st.Expect(t, results[1].Risked, 3195.00)
	st.Expect(t, results[1].Credit, 405.00)
	st.Expect(t, results[1].Profit, -863.60)
	st.Expect(t, results[1].Proceeds, -1242.00)
	st.Expect(t, results[1].Commission, 26.60)

	// Spot check some call credit spreads
	st.Expect(t, results[6].Status, "Open")
	st.Expect(t, results[6].Type, "Call Credit Spread")
	st.Expect(t, results[6].OrderIds, "13")
	st.Expect(t, results[6].Risked, 712.00)
	st.Expect(t, results[6].Credit, 288.00)
	st.Expect(t, results[6].Profit, 0.00)
	st.Expect(t, results[6].Proceeds, 0.00)
	st.Expect(t, results[6].Commission, 7.00)
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
	commission := calcCommissionForOrder(&order, 1, &brokerAccount)

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
	commission := calcCommissionForOrder(&order, 1, &brokerAccount)

	// Verify the data was return as expected
	st.Expect(t, commission, 8.75)
}

/* End File */
