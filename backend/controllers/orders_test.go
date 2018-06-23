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
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/nbio/st"
)

//
// Test - PreviewOrder
//
func TestPreviewOrder01(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()

	// Create controller
	c := &Controller{DB: db}

	// // Flush pending mocks after test execution
	// defer gock.Off()

	// // Setup mock request.
	// gock.New("https://api.tradier.com/v1").
	// 	Post("/accounts/abc1235423/orders").
	// 	Reply(200).
	// 	BodyString(`{"jjj": "jj"}`)

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

	spew.Dump(result)

	// Test result
	st.Expect(t, err, nil)
}

/* End File */
