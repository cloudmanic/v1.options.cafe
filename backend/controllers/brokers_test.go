//
// Date: 2/23/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/brokers/types"
	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/nbio/st"
	gock "gopkg.in/h2non/gock.v1"
)

//
// Test - GetBrokers
//
func TestGetBrokers01(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()

	// Create controller
	c := &Controller{DB: db}

	// Shared vars we use.
	ts := time.Date(2017, 10, 29, 17, 20, 01, 507451, time.UTC)

	// Install test data.
	db.Exec("TRUNCATE TABLE brokers;")
	db.Create(&models.Broker{Name: "Tradier", UserId: 1, AccessToken: "CLOwLO2cMnx-N_bPEexiVo9z9oRR80nPI9ycxQw3KQ-WQ4OP3D44gIbfLScAZ9pv", RefreshToken: "abc", TokenExpirationDate: ts, Status: "Active"})
	db.Create(&models.Broker{Name: "Tradeking", UserId: 1, AccessToken: "456", RefreshToken: "xyz", TokenExpirationDate: ts, Status: "Active"})
	db.Create(&models.Broker{Name: "Etrade", UserId: 1, AccessToken: "789", RefreshToken: "mno", TokenExpirationDate: ts, Status: "Active"})

	db.Exec("TRUNCATE TABLE broker_accounts;")
	db.Create(&models.BrokerAccount{UserId: 1, BrokerId: 1, Name: "Test Account 1", AccountNumber: "YYY123ZY", StockCommission: 5.00, StockMin: 0.00, OptionCommission: 0.35, OptionSingleMin: 5.00, OptionMultiLegMin: 7.00, OptionBase: 0.00})
	db.Create(&models.BrokerAccount{UserId: 1, BrokerId: 1, Name: "Test Account 2", AccountNumber: "ABC123ZY", StockCommission: 5.00, StockMin: 0.00, OptionCommission: 0.35, OptionSingleMin: 5.00, OptionMultiLegMin: 7.00, OptionBase: 0.00})

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
// Test - CreateBroker - 01 - Success
//
func TestCreateBroker01(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()

	// Create controller
	c := &Controller{DB: db}

	// Shared vars we use.
	ts := time.Date(2017, 10, 29, 17, 20, 01, 507451, time.UTC)

	// Install test data.
	db.Exec("TRUNCATE TABLE brokers;")
	db.Create(&models.Broker{Name: "Tradier", UserId: 1, AccessToken: "CLOwLO2cMnx-N_bPEexiVo9z9oRR80nPI9ycxQw3KQ-WQ4OP3D44gIbfLScAZ9pv", RefreshToken: "abc", TokenExpirationDate: ts, Status: "Active"})

	// Body data // We send in limited data as all other broker data is added through the auth process.
	var bodyStr = []byte(`{"name":"Tradier", "display_name":"Unit Test Broker#1"}`)

	// Make a mock request.
	req, _ := http.NewRequest("POST", "/api/v1/brokers", bytes.NewBuffer(bodyStr))
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(1)) })

	r.POST("/api/v1/brokers", c.CreateBroker)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	result := models.Broker{}
	err := json.Unmarshal([]byte(w.Body.String()), &result)

	// Parse json that returned.
	st.Expect(t, err, nil)
	st.Expect(t, w.Body.String(), `{"id":2,"name":"Tradier","display_name":"Unit Test Broker#1","broker_accounts":null,"status":"Disabled"}`)
}

//
// Test - CreateBroker 02 - Fail
//
func TestCreateBroker02(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()

	// Create controller
	c := &Controller{DB: db}

	// Shared vars we use.
	ts := time.Date(2017, 10, 29, 17, 20, 01, 507451, time.UTC)

	// Install test data.
	db.Exec("TRUNCATE TABLE brokers;")
	db.Create(&models.Broker{Name: "Tradier", UserId: 1, AccessToken: "CLOwLO2cMnx-N_bPEexiVo9z9oRR80nPI9ycxQw3KQ-WQ4OP3D44gIbfLScAZ9pv", RefreshToken: "abc", TokenExpirationDate: ts, Status: "Active"})

	// Body data // We send in limited data as all other broker data is added through the auth process.
	var bodyStr = []byte(`{"name":"Bad Broker", "display_name":"Unit Test Broker#1"}`)

	// Make a mock request.
	req, _ := http.NewRequest("POST", "/api/v1/brokers", bytes.NewBuffer(bodyStr))
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(1)) })

	r.POST("/api/v1/brokers", c.CreateBroker)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	result := models.Broker{}
	err := json.Unmarshal([]byte(w.Body.String()), &result)

	// Parse json that returned.
	st.Expect(t, err, nil)
	st.Expect(t, w.Body.String(), `{"error":"Broker name not valid."}`)
}

//
// Test - CreateBroker 03 - Fail
//
func TestCreateBroker03(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()

	// Create controller
	c := &Controller{DB: db}

	// Shared vars we use.
	ts := time.Date(2017, 10, 29, 17, 20, 01, 507451, time.UTC)

	// Install test data.
	db.Exec("TRUNCATE TABLE brokers;")
	db.Create(&models.Broker{Name: "Tradier", UserId: 1, AccessToken: "CLOwLO2cMnx-N_bPEexiVo9z9oRR80nPI9ycxQw3KQ-WQ4OP3D44gIbfLScAZ9pv", RefreshToken: "abc", TokenExpirationDate: ts, Status: "Active"})

	// Body data // We send in limited data as all other broker data is added through the auth process.
	var bodyStr = []byte(`{"name":"", "display_name":""}`)

	// Make a mock request.
	req, _ := http.NewRequest("POST", "/api/v1/brokers", bytes.NewBuffer(bodyStr))
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(1)) })

	r.POST("/api/v1/brokers", c.CreateBroker)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	result := models.Broker{}
	err := json.Unmarshal([]byte(w.Body.String()), &result)

	// Parse json that returned.
	st.Expect(t, err, nil)
	st.Expect(t, w.Body.String(), `{"error":"Name field can not be empty."}`)
}

//
// Test - CreateBroker 04 - Fail
//
func TestCreateBroker04(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()

	// Create controller
	c := &Controller{DB: db}

	// Shared vars we use.
	ts := time.Date(2017, 10, 29, 17, 20, 01, 507451, time.UTC)

	// Install test data.
	db.Exec("TRUNCATE TABLE brokers;")
	db.Create(&models.Broker{Name: "Tradier", UserId: 1, AccessToken: "CLOwLO2cMnx-N_bPEexiVo9z9oRR80nPI9ycxQw3KQ-WQ4OP3D44gIbfLScAZ9pv", RefreshToken: "abc", TokenExpirationDate: ts, Status: "Active"})

	// Body data // We send in limited data as all other broker data is added through the auth process.
	var bodyStr = []byte(`{"name":"Tradier", "display_name":""}`)

	// Make a mock request.
	req, _ := http.NewRequest("POST", "/api/v1/brokers", bytes.NewBuffer(bodyStr))
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(1)) })

	r.POST("/api/v1/brokers", c.CreateBroker)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	result := models.Broker{}
	err := json.Unmarshal([]byte(w.Body.String()), &result)

	// Parse json that returned.
	st.Expect(t, err, nil)
	st.Expect(t, w.Body.String(), `{"error":"Display Name field can not be empty."}`)
}

//
// Test - GetBalances
//
func TestGetBalances01(t *testing.T) {

	// Flush pending mocks after test execution
	defer gock.Off()

	// Setup mock request.
	gock.New("https://api.tradier.com/v1").
		Get("/user/balances").
		Reply(200).
		BodyString(`{"accounts":{"account":[{"balances":{"option_short_value":0,"total_equity":0.00000000,"account_number":"6Y111184","account_type":"cash","close_pl":0,"current_requirement":0,"equity":0,"long_market_value":0,"market_value":0,"open_pl":0,"option_long_value":0,"option_requirement":0,"pending_orders_count":0,"short_market_value":0,"stock_long_value":0,"total_cash":0.00000000,"uncleared_funds":0,"pending_cash":0,"cash":{"cash_available":0.00000000,"sweep":0,"unsettled_funds":0}},"account_number":"6Y111184"},{"balances":{"option_short_value":-5166.0000000000000000000,"total_equity":115978.2300000000000000000,"account_number":"6Y777785","account_type":"margin","close_pl":0.00000000,"current_requirement":12600.0000000000000000,"equity":0,"long_market_value":0,"market_value":-751.5000000000000000000,"open_pl":850.5000000000000000000,"option_long_value":4414.5000000000000000000,"option_requirement":12600.0000000000000000,"pending_orders_count":14,"short_market_value":0,"stock_long_value":0,"total_cash":116729.73000000,"uncleared_funds":0,"pending_cash":0,"margin":{"fed_call":0,"maintenance_call":0,"option_buying_power":4129.73000000,"stock_buying_power":8259.46,"stock_short_value":0,"sweep":0}},"account_number":"6Y777785"},{"balances":{"option_short_value":0,"total_equity":3165.660000000000000000,"account_number":"6YA88882","account_type":"cash","close_pl":0.00000000,"current_requirement":0,"equity":0,"long_market_value":0,"market_value":2903.260000000000000000,"open_pl":946.243200000000000000,"option_long_value":0,"option_requirement":0,"pending_orders_count":0,"short_market_value":0,"stock_long_value":2903.260000000000000000,"total_cash":262.40000000,"uncleared_funds":0,"pending_cash":0,"cash":{"cash_available":0.00000000,"sweep":0,"unsettled_funds":0}},"account_number":"6YA88882"}]}}`)

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
	st.Expect(t, err, nil)
	st.Expect(t, len(result), 3)
	st.Expect(t, result[0].AccountNumber, "6Y111184")
	st.Expect(t, result[0].AccountValue, 0.00)
	st.Expect(t, result[0].TotalCash, 0.00)
	st.Expect(t, result[0].OptionBuyingPower, 0.00)
	st.Expect(t, result[0].StockBuyingPower, 0.00)
	st.Expect(t, result[1].AccountNumber, "6Y777785")
	st.Expect(t, result[1].AccountValue, 115978.23)
	st.Expect(t, result[1].TotalCash, 116729.73)
	st.Expect(t, result[1].OptionBuyingPower, 4129.73)
	st.Expect(t, result[1].StockBuyingPower, 8259.46)
	st.Expect(t, result[2].AccountNumber, "6YA88882")
	st.Expect(t, result[2].AccountValue, 3165.66)
	st.Expect(t, result[2].TotalCash, 262.4)
	st.Expect(t, result[2].OptionBuyingPower, 0.00)
	st.Expect(t, result[2].StockBuyingPower, 0.00)
	st.Expect(t, gock.IsDone(), true)
}

/* End File */
