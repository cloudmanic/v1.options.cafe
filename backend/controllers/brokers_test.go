//
// Date: 2/23/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudmanic/app.options.cafe/backend/brokers/types"
	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/nbio/st"
)

//
// Test - GetBrokers
//
func TestGetBrokers01(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()

	// Create controller
	c := &Controller{DB: db}

	// Make a mock request.
	req, _ := http.NewRequest("GET", "/api/v1/brokers", nil)
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(1)) })

	r.GET("/api/v1/brokers", c.GetBrokers)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	result := []models.Broker{}
	err := json.Unmarshal([]byte(w.Body.String()), &result)

	// Parse json that returned.
	st.Expect(t, err, nil)
	st.Expect(t, len(result), 3)
	st.Expect(t, result[0].Id, uint(1))
	st.Expect(t, result[1].Id, uint(2))
	st.Expect(t, result[2].Id, uint(3))
	st.Expect(t, result[0].Name, "Tradier")
	st.Expect(t, result[1].Name, "Tradeking")
	st.Expect(t, result[2].Name, "Etrade")
	st.Expect(t, len(result[0].BrokerAccounts), 2)
	st.Expect(t, len(result[1].BrokerAccounts), 0)
	st.Expect(t, len(result[2].BrokerAccounts), 0)
	st.Expect(t, result[0].BrokerAccounts[0].AccountNumber, "YYY123ZY")
	st.Expect(t, result[0].BrokerAccounts[1].AccountNumber, "ABC123ZY")
	st.Expect(t, result[0].BrokerAccounts[0].BrokerId, uint(1))
	st.Expect(t, result[0].BrokerAccounts[1].BrokerId, uint(1))
	st.Expect(t, result[0].BrokerAccounts[0].AccountNumber, "YYY123ZY")
	st.Expect(t, result[0].BrokerAccounts[1].AccountNumber, "ABC123ZY")
}

//
// Test - GetBalances
//
func TestGetBalances01(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()

	// Create controller
	c := &Controller{DB: db}

	// Make a mock request.
	req, _ := http.NewRequest("GET", "/api/v1/brokers/balances", nil)
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(1)) })

	r.GET("/api/v1/brokers/balances", c.GetBalances)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	result := []types.Balance{}
	err := json.Unmarshal([]byte(w.Body.String()), &result)

	// Parse json that returned.
	// TODO: WE can't test this any further as the Tradier API does not return in any particular order.
	// We could reorder the results and then test. So TODO here.
	st.Expect(t, err, nil)
	st.Expect(t, len(result), 3)
}

/* End File */
