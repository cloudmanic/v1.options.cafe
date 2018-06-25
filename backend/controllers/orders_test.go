//
// Date: 6/23/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudmanic/app.options.cafe/backend/brokers/types"
	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/nbio/st"
	gock "gopkg.in/h2non/gock.v1"
)

//
// Test - PreviewOrder 01
//
func TestPreviewOrder01(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()

	// Create controller
	c := &Controller{DB: db}

	// Flush pending mocks after test execution
	defer gock.Off()

	// Setup mock request.
	gock.New("https://api.tradier.com/v1").
		Post("/accounts/abc1235423/orders").
		Reply(200).
		BodyString(`{"order":{"status":"ok","commission":7.00000000,"cost":538.0000000000,"fees":0,"symbol":"SPY","type":"credit","duration":"day","result":true,"price":0.23,"order_cost":-69.0000000000,"margin_change":600.00000000,"option_requirement":600.00000000,"request_date":"2018-06-24T07:10:14.14","extended_hours":false,"class":"multileg","strategy":"spread"}}`)

	// Post data
	order := types.Order{
		BrokerAccountId: 1,
		Class:           "multileg",
		Symbol:          "SPY",
		Duration:        "day",
		Type:            "credit",
		Price:           0.23,
		Legs: []types.OrderLeg{
			{Side: "buy_to_open", Quantity: 3, OptionSymbol: "SPY180709P00268000"},
			{Side: "sell_to_open", Quantity: 3, OptionSymbol: "SPY180709P00270000"},
		},
	}

	// Encode json
	jsonText, _ := json.Marshal(order)

	// Make a mock request.
	req, _ := http.NewRequest("POST", "/api/v1/orders/preview", bytes.NewBuffer(jsonText))
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(1)) })

	r.POST("/api/v1/orders/preview", c.PreviewOrder)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	result := types.OrderPreview{}
	err := json.Unmarshal([]byte(w.Body.String()), &result)

	// Test result
	st.Expect(t, err, nil)
	st.Expect(t, result.Status, "ok")
	st.Expect(t, result.Cost, 538.00)
	st.Expect(t, result.Type, "credit")
	st.Expect(t, result.Strategy, "spread")
	st.Expect(t, result.Commission, 7.00)
}

//
// Test - PreviewOrder 02
//
func TestPreviewOrder02(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()

	// Create controller
	c := &Controller{DB: db}

	// Flush pending mocks after test execution
	defer gock.Off()

	// Setup mock request.
	gock.New("https://api.tradier.com/v1").
		Post("/accounts/abc1235423/orders").
		Reply(200).
		BodyString(`{"errors":{"error":"You do not have enough buying power for this trade."}}`)

	// Post data
	order := types.Order{
		BrokerAccountId: 1,
		Class:           "multileg",
		Symbol:          "SPY",
		Duration:        "day",
		Type:            "credit",
		Price:           0.23,
		Legs: []types.OrderLeg{
			{Side: "buy_to_open", Quantity: 300, OptionSymbol: "SPY180709P00268000"},
			{Side: "sell_to_open", Quantity: 300, OptionSymbol: "SPY180709P00270000"},
		},
	}

	// Encode json
	jsonText, _ := json.Marshal(order)

	// Make a mock request.
	req, _ := http.NewRequest("POST", "/api/v1/orders/preview", bytes.NewBuffer(jsonText))
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(1)) })

	r.POST("/api/v1/orders/preview", c.PreviewOrder)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	result := types.OrderPreview{}
	err := json.Unmarshal([]byte(w.Body.String()), &result)

	// Test result
	st.Expect(t, err, nil)
	st.Expect(t, result.Status, "Error")
	st.Expect(t, result.Cost, 0.00)
	st.Expect(t, result.Type, "")
	st.Expect(t, result.Strategy, "")
	st.Expect(t, result.Commission, 0.00)
	st.Expect(t, result.Error, "You do not have enough buying power for this trade.")
}

//
// Test - SubmitOrder 01
//
func TestSubmitOrder01(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()

	// Create controller
	c := &Controller{DB: db}

	// Flush pending mocks after test execution
	defer gock.Off()

	// Setup mock request.
	gock.New("https://api.tradier.com/v1").
		Post("/accounts/abc1235423/orders").
		Reply(200).
		BodyString(`{"order":{"id":1086137,"status":"ok","partner_id":"1186d9ee-eff8-452c-8cd9-e5d8aaa1157a"}}`)

	// Post data
	order := types.Order{
		BrokerAccountId: 1,
		Class:           "multileg",
		Symbol:          "SPY",
		Duration:        "day",
		Type:            "credit",
		Price:           0.23,
		Legs: []types.OrderLeg{
			{Side: "buy_to_open", Quantity: 3, OptionSymbol: "SPY180709P00268000"},
			{Side: "sell_to_open", Quantity: 3, OptionSymbol: "SPY180709P00270000"},
		},
	}

	// Encode json
	jsonText, _ := json.Marshal(order)

	// Make a mock request.
	req, _ := http.NewRequest("POST", "/api/v1/orders", bytes.NewBuffer(jsonText))
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(1)) })

	r.POST("/api/v1/orders", c.SubmitOrder)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	result := types.OrderSubmit{}
	err := json.Unmarshal([]byte(w.Body.String()), &result)

	// Test result
	st.Expect(t, err, nil)
	st.Expect(t, result.Status, "ok")
	st.Expect(t, result.Id, uint(1086137))
}

/* End File */
