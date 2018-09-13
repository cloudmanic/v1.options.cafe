//
// Date: 7/18/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/nbio/st"
)

//
// Load and return some test data.
//
func LoadAndGetBrokerEventsTestData() {

	// Start the db connection.
	db, _ := models.NewDB()
	defer db.Close()

	// Put test data into database
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

	// Load test Event data,
	ts := time.Date(2017, 10, 29, 17, 20, 01, 507451, time.UTC)
	db.Create(&models.BrokerEvent{UserId: 1, CreatedAt: ts, UpdatedAt: ts, BrokerAccountId: 3, BrokerId: "46dcdd5817b760661a353308b81c5aac", Type: "Trade", Date: ts, Symbol: "spy", Commission: 5.00, Description: "test 123", Price: 13.33, Quantity: 122, TradeType: "Equity"})
	db.Create(&models.BrokerEvent{UserId: 1, CreatedAt: ts, UpdatedAt: ts, BrokerAccountId: 3, BrokerId: "46dcdd5817b76066ff35d3308b81c5aac", Type: "Trade", Date: ts, Symbol: "spy", Commission: 5.00, Description: "test abc", Price: 32.33, Quantity: 1, TradeType: "Interest"})
	db.Create(&models.BrokerEvent{UserId: 1, CreatedAt: ts, UpdatedAt: ts, BrokerAccountId: 3, BrokerId: "46dcdd5817b76066234353308b81c5aac", Type: "Trade", Date: ts, Symbol: "spy", Commission: 5.00, Description: "test ttt", Price: 23.33, Quantity: 17, TradeType: "Option"})
	db.Create(&models.BrokerEvent{UserId: 1, CreatedAt: ts, UpdatedAt: ts, BrokerAccountId: 3, BrokerId: "46dcdd5817b76066113153308b81c5aac", Type: "Trade", Date: ts, Symbol: "spy", Commission: 5.00, Description: "test 888", Price: 223.33, Quantity: 18, TradeType: "Equity"})
}

//
// TestGetBrokerEvents01
//
func TestGetBrokerEvents01(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()

	// Load test data
	LoadAndGetBrokerEventsTestData()

	// Create controller
	c := &Controller{DB: db}

	// Make a mock request.
	req, _ := http.NewRequest("GET", "/api/v1/broker-events/3", nil)
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(1)) })

	r.GET("/api/v1/broker-events/:brokerAccount", c.GetBrokerEvents)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	result := []models.BrokerEvent{}
	err := json.Unmarshal([]byte(w.Body.String()), &result)

	// Parse json that returned.
	st.Expect(t, err, nil)
	st.Expect(t, w.Code, 200)
	st.Expect(t, w.Body.String(), `[{"id":1,"type":"Trade","date":"2017-10-29T10:20:01-07:00","amount":0,"symbol":"spy","commission":5,"description":"test 123","price":13.33,"quantity":122,"trade_type":"Equity"},{"id":2,"type":"Trade","date":"2017-10-29T10:20:01-07:00","amount":0,"symbol":"spy","commission":5,"description":"test ttt","price":23.33,"quantity":17,"trade_type":"Option"},{"id":3,"type":"Trade","date":"2017-10-29T10:20:01-07:00","amount":0,"symbol":"spy","commission":5,"description":"test 888","price":223.33,"quantity":18,"trade_type":"Equity"}]`)
}

/* End File */
