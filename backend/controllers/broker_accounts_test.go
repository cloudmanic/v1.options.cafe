//
// Date: 9/21/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/nbio/st"
)

//
// Test - UpdateBrokerAccount - 01 - Success
//
func TestUpdateBrokerAccount01(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()

	// Create controller
	c := &Controller{DB: db}

	// Shared vars we use.
	ts := time.Date(2017, 10, 29, 17, 20, 01, 507451, time.UTC)

	// Install test data.
	db.Exec("TRUNCATE TABLE brokers;")
	db.Create(&models.Broker{Name: "Tradier", UserId: 1, AccessToken: "CLOwLO2cMnx-N_bPEexiVo9z9oRR80nPI9ycxQw3KQ-WQ4OP3D44gIbfLScAZ9pv", RefreshToken: "abc", TokenExpirationDate: ts, Status: "Active"})

	db.Exec("TRUNCATE TABLE broker_accounts;")
	db.Create(&models.BrokerAccount{UserId: 1, BrokerId: 1, Name: "Test Account 1", AccountNumber: "YYY123ZY", StockCommission: 1.00, StockMin: 1.00, OptionCommission: 1.35, OptionSingleMin: 1.00, OptionMultiLegMin: 1.00, OptionBase: 1.00})
	db.Create(&models.BrokerAccount{UserId: 1, BrokerId: 1, Name: "Test Account 2", AccountNumber: "ABC123ZY", StockCommission: 1.00, StockMin: 1.00, OptionCommission: 1.35, OptionSingleMin: 1.00, OptionMultiLegMin: 1.00, OptionBase: 1.00})

	// Body data // We send in limited data as all other broker data is added through the auth process.
	var bodyStr = []byte(`{"name":"unit_test_account_#1", "account_number":"abc123", "stock_commission": 3.95, "stock_min": 5.00, "option_commission": 0.35, "option_single_min": 5.00, "option_multi_leg_min": 7.00, "option_base": 3.00}`)

	// Make a mock request.
	req, _ := http.NewRequest("PUT", "/api/v1/brokers/1/accounts/1", bytes.NewBuffer(bodyStr))
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(1)) })

	r.PUT("/api/v1/brokers/:id/accounts/:acctId", c.UpdateBrokerAccount)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Validate result
	st.Expect(t, w.Code, 204)

	// Get the broker account we just updated
	orgObj, _ := db.GetBrokerAccountByIdUserId(uint(1), uint(1))
	st.Expect(t, orgObj.Id, uint(1))
	st.Expect(t, orgObj.AccountNumber, "YYY123ZY")
	st.Expect(t, orgObj.StockCommission, 3.95)
	st.Expect(t, orgObj.StockMin, 5.00)
	st.Expect(t, orgObj.OptionCommission, 0.35)
	st.Expect(t, orgObj.OptionSingleMin, 5.00)
	st.Expect(t, orgObj.OptionMultiLegMin, 7.00)
	st.Expect(t, orgObj.OptionBase, 3.00)
}

//
// Test - UpdateBrokerAccount - 02 - Failed
//
func TestUpdateBrokerAccount02(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()

	// Create controller
	c := &Controller{DB: db}

	// Shared vars we use.
	ts := time.Date(2017, 10, 29, 17, 20, 01, 507451, time.UTC)

	// Install test data.
	db.Exec("TRUNCATE TABLE brokers;")
	db.Create(&models.Broker{Name: "Tradier", UserId: 1, AccessToken: "CLOwLO2cMnx-N_bPEexiVo9z9oRR80nPI9ycxQw3KQ-WQ4OP3D44gIbfLScAZ9pv", RefreshToken: "abc", TokenExpirationDate: ts, Status: "Active"})

	db.Exec("TRUNCATE TABLE broker_accounts;")
	db.Create(&models.BrokerAccount{UserId: 1, BrokerId: 1, Name: "Test Account 1", AccountNumber: "YYY123ZY", StockCommission: 1.00, StockMin: 1.00, OptionCommission: 1.35, OptionSingleMin: 1.00, OptionMultiLegMin: 1.00, OptionBase: 1.00})
	db.Create(&models.BrokerAccount{UserId: 1, BrokerId: 1, Name: "Test Account 2", AccountNumber: "ABC123ZY", StockCommission: 1.00, StockMin: 1.00, OptionCommission: 1.35, OptionSingleMin: 1.00, OptionMultiLegMin: 1.00, OptionBase: 1.00})

	// Body data // We send in limited data as all other broker data is added through the auth process.
	var bodyStr = []byte(`{"name":"unit_test_account_#1", "account_number":"abc123", "stock_commission": 3.95, "stock_min": 5.00, "option_commission": 0.35, "option_single_min": 5.00, "option_multi_leg_min": 7.00, "option_base": 3.00}`)

	// Make a mock request.
	req, _ := http.NewRequest("PUT", "/api/v1/brokers/1/accounts/5", bytes.NewBuffer(bodyStr))
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(1)) })

	r.PUT("/api/v1/brokers/:id/accounts/:acctId", c.UpdateBrokerAccount)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Validate result
	st.Expect(t, w.Code, 400)
	st.Expect(t, w.Body.String(), `{"error":"No Record Found."}`)
}

//
// Test - UpdateBrokerAccount - 03 - Failed
//
func TestUpdateBrokerAccount03(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()

	// Create controller
	c := &Controller{DB: db}

	// Shared vars we use.
	ts := time.Date(2017, 10, 29, 17, 20, 01, 507451, time.UTC)

	// Install test data.
	db.Exec("TRUNCATE TABLE brokers;")
	db.Create(&models.Broker{Name: "Tradier", UserId: 1, AccessToken: "CLOwLO2cMnx-N_bPEexiVo9z9oRR80nPI9ycxQw3KQ-WQ4OP3D44gIbfLScAZ9pv", RefreshToken: "abc", TokenExpirationDate: ts, Status: "Active"})

	db.Exec("TRUNCATE TABLE broker_accounts;")
	db.Create(&models.BrokerAccount{UserId: 1, BrokerId: 1, Name: "Test Account 1", AccountNumber: "YYY123ZY", StockCommission: 1.00, StockMin: 1.00, OptionCommission: 1.35, OptionSingleMin: 1.00, OptionMultiLegMin: 1.00, OptionBase: 1.00})
	db.Create(&models.BrokerAccount{UserId: 1, BrokerId: 1, Name: "Test Account 2", AccountNumber: "ABC123ZY", StockCommission: 1.00, StockMin: 1.00, OptionCommission: 1.35, OptionSingleMin: 1.00, OptionMultiLegMin: 1.00, OptionBase: 1.00})

	// Body data // We send in limited data as all other broker data is added through the auth process.
	var bodyStr = []byte(`{"name":"unit_test_account_#1", "account_number":"abc123", "stock_commission": 3.95, "stock_min": 5.00, "option_commission": 0.35, "option_single_min": 5.00, "option_multi_leg_min": 7.00, "option_base": 3.00}`)

	// Make a mock request.
	req, _ := http.NewRequest("PUT", "/api/v1/brokers/2/accounts/1", bytes.NewBuffer(bodyStr))
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(1)) })

	r.PUT("/api/v1/brokers/:id/accounts/:acctId", c.UpdateBrokerAccount)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Validate result
	st.Expect(t, w.Code, 400)
	st.Expect(t, w.Body.String(), `{"error":"No Record Found."}`)
}

//
// Test - UpdateBrokerAccount - 04 - Success
//
func TestUpdateBrokerAccount04(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()

	// Create controller
	c := &Controller{DB: db}

	// Shared vars we use.
	ts := time.Date(2017, 10, 29, 17, 20, 01, 507451, time.UTC)

	// Install test data.
	db.Exec("TRUNCATE TABLE brokers;")
	db.Create(&models.Broker{Name: "Tradier", UserId: 1, AccessToken: "CLOwLO2cMnx-N_bPEexiVo9z9oRR80nPI9ycxQw3KQ-WQ4OP3D44gIbfLScAZ9pv", RefreshToken: "abc", TokenExpirationDate: ts, Status: "Active"})

	db.Exec("TRUNCATE TABLE broker_accounts;")
	db.Create(&models.BrokerAccount{UserId: 1, BrokerId: 1, Name: "Test Account 1", AccountNumber: "YYY123ZY", StockCommission: 1.00, StockMin: 1.00, OptionCommission: 1.35, OptionSingleMin: 1.00, OptionMultiLegMin: 1.00, OptionBase: 1.00})
	db.Create(&models.BrokerAccount{UserId: 1, BrokerId: 1, Name: "Test Account 2", AccountNumber: "ABC123ZY", StockCommission: 1.00, StockMin: 1.00, OptionCommission: 1.35, OptionSingleMin: 1.00, OptionMultiLegMin: 1.00, OptionBase: 1.00})

	// Body data // We send in limited data as all other broker data is added through the auth process.
	var bodyStr = []byte(`{"name":"unit_test_account_#1", "account_number":"abc123", "stock_commission": 3.95, "stock_min": 5.00, "option_commission": 0.35, "option_single_min": 5.00, "option_multi_leg_min": 7.00, "option_base": 3.00}`)

	// Make a mock request.
	req, _ := http.NewRequest("PUT", "/api/v1/brokers/1/accounts/2", bytes.NewBuffer(bodyStr))
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(1)) })

	r.PUT("/api/v1/brokers/:id/accounts/:acctId", c.UpdateBrokerAccount)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Validate result
	st.Expect(t, w.Code, 204)

	// Get the broker account we just updated
	orgObj, _ := db.GetBrokerAccountByIdUserId(uint(2), uint(1))
	st.Expect(t, orgObj.Id, uint(2))
	st.Expect(t, orgObj.AccountNumber, "ABC123ZY")
	st.Expect(t, orgObj.StockCommission, 3.95)
	st.Expect(t, orgObj.StockMin, 5.00)
	st.Expect(t, orgObj.OptionCommission, 0.35)
	st.Expect(t, orgObj.OptionSingleMin, 5.00)
	st.Expect(t, orgObj.OptionMultiLegMin, 7.00)
	st.Expect(t, orgObj.OptionBase, 3.00)
}

//
// Test - UpdateBrokerAccount - 05 - Failed
//
func TestUpdateBrokerAccount05(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()

	// Create controller
	c := &Controller{DB: db}

	// Shared vars we use.
	ts := time.Date(2017, 10, 29, 17, 20, 01, 507451, time.UTC)

	// Install test data.
	db.Exec("TRUNCATE TABLE brokers;")
	db.Create(&models.Broker{Name: "Tradier", UserId: 1, AccessToken: "CLOwLO2cMnx-N_bPEexiVo9z9oRR80nPI9ycxQw3KQ-WQ4OP3D44gIbfLScAZ9pv", RefreshToken: "abc", TokenExpirationDate: ts, Status: "Active"})

	db.Exec("TRUNCATE TABLE broker_accounts;")
	db.Create(&models.BrokerAccount{UserId: 1, BrokerId: 1, Name: "Test Account 1", AccountNumber: "YYY123ZY", StockCommission: 1.00, StockMin: 1.00, OptionCommission: 1.35, OptionSingleMin: 1.00, OptionMultiLegMin: 1.00, OptionBase: 1.00})
	db.Create(&models.BrokerAccount{UserId: 1, BrokerId: 1, Name: "Test Account 2", AccountNumber: "ABC123ZY", StockCommission: 1.00, StockMin: 1.00, OptionCommission: 1.35, OptionSingleMin: 1.00, OptionMultiLegMin: 1.00, OptionBase: 1.00})

	// Body data // We send in limited data as all other broker data is added through the auth process.
	var bodyStr = []byte(`{"name":"unit_test_account_#1", "account_number":"abc123", "stock_commission": "3.95", "stock_min": 5.00, "option_commission": 0.35, "option_single_min": 5.00, "option_multi_leg_min": 7.00, "option_base": 3.00}`)

	// Make a mock request.
	req, _ := http.NewRequest("PUT", "/api/v1/brokers/1/accounts/5", bytes.NewBuffer(bodyStr))
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(1)) })

	r.PUT("/api/v1/brokers/:id/accounts/:acctId", c.UpdateBrokerAccount)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Validate result
	st.Expect(t, w.Code, 400)
	st.Expect(t, w.Body.String(), `{"error":"Invalid JSON in body. There is a chance the JSON maybe valid but does not match the data type requirements. For example maybe you passed a string in for an integer."}`)
}

/* End File */
