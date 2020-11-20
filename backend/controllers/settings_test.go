//
// Date: 7/18/2018
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

	"app.options.cafe/models"
	"github.com/gin-gonic/gin"
	"github.com/nbio/st"
)

//
// TestGetSettings01
//
func TestGetSettings01(t *testing.T) {

	// Json response
	jsonResponse := `{"id":1,"user_id":1,"strategy_pcs_close_price":0.03,"strategy_pcs_open_price":"mid-point","strategy_pcs_lots":10,"strategy_ccs_close_price":0.03,"strategy_ccs_open_price":"mid-point","strategy_ccs_lots":10,"strategy_pds_close_price":0.03,"strategy_pds_open_price":"mid-point","strategy_pds_lots":10,"strategy_cds_close_price":0.03,"strategy_cds_open_price":"mid-point","strategy_cds_lots":10,"notice_trade_filled_email":"No","notice_trade_filled_sms":"No","notice_trade_filled_push":"No","notice_market_open_email":"No","notice_market_open_sms":"No","notice_market_open_push":"No","notice_market_closed_email":"No","notice_market_closed_sms":"No","notice_market_closed_push":"No"}`

	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{DB: db}

	// Make a mock request.
	req, _ := http.NewRequest("GET", "/api/v1/settings", nil)
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(1)) })

	r.GET("/api/v1/settings", c.GetSettings)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Parse json that returned.
	st.Expect(t, w.Code, 200)
	st.Expect(t, w.Body.String(), jsonResponse)

	// ------- We call it again to make sure a second DB entry is not created. ------------ //

	// Make a mock request.
	req2, _ := http.NewRequest("GET", "/api/v1/settings", nil)
	req2.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r2 := gin.New()

	r2.Use(func(c *gin.Context) { c.Set("userId", uint(1)) })

	r2.GET("/api/v1/settings", c.GetSettings)

	// Setup writer.
	w2 := httptest.NewRecorder()
	r2.ServeHTTP(w2, req)

	// Parse json that returned.
	st.Expect(t, w2.Code, 200)
	st.Expect(t, w2.Body.String(), jsonResponse)
}

//
// TestGetSettings01
//
func TestUpdateSettings01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{DB: db}

	// Create a temp setting
	settings := db.SettingsGetOrCreateByUserId(1)

	// Change some values to test.
	settings.StrategyCcsClosePrice = 0.40
	settings.NoticeTradeFilledEmail = "Yes"
	settings.NoticeTradeFilledSms = "Yes"
	json, _ := json.Marshal(settings)

	// Make a mock request.
	req, _ := http.NewRequest("PUT", "/api/v1/settings", bytes.NewBuffer(json))
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(1)) })

	r.PUT("/api/v1/settings", c.UpdateSettings)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Parse json that returned.
	st.Expect(t, w.Code, 204)
	st.Expect(t, w.Body.String(), "")

	// Get settings again to verify
	s := db.SettingsGetOrCreateByUserId(1)

	// Verify
	st.Expect(t, s.StrategyCcsClosePrice, 0.40)
	st.Expect(t, s.NoticeTradeFilledEmail, "Yes")
	st.Expect(t, s.NoticeTradeFilledSms, "Yes")
	st.Expect(t, s.NoticeTradeFilledPush, "No")
	st.Expect(t, s.StrategyCcsOpenPrice, "mid-point")
	st.Expect(t, s.StrategyPcsClosePrice, 0.03)
}

//
// TestGetSettings02
//
func TestUpdateSettings02(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{DB: db}

	// Create a temp setting
	settings := db.SettingsGetOrCreateByUserId(1)

	// Change some values to test.
	settings.UserId = 12 // Making sure this is ignored
	settings.StrategyCcsClosePrice = 0.40
	settings.NoticeTradeFilledEmail = "Yes"
	settings.NoticeTradeFilledSms = "Yes"
	json, _ := json.Marshal(settings)

	// Make a mock request.
	req, _ := http.NewRequest("PUT", "/api/v1/settings", bytes.NewBuffer(json))
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(1)) })

	r.PUT("/api/v1/settings", c.UpdateSettings)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Parse json that returned.
	st.Expect(t, w.Code, 204)
	st.Expect(t, w.Body.String(), "")

	// Get settings again to verify
	s := db.SettingsGetOrCreateByUserId(1)

	// Verify
	st.Expect(t, s.UserId, uint(1))
	st.Expect(t, s.StrategyCcsClosePrice, 0.40)
	st.Expect(t, s.NoticeTradeFilledEmail, "Yes")
	st.Expect(t, s.NoticeTradeFilledSms, "Yes")
	st.Expect(t, s.NoticeTradeFilledPush, "No")
	st.Expect(t, s.StrategyCcsOpenPrice, "mid-point")
	st.Expect(t, s.StrategyPcsClosePrice, 0.03)
}

//
// TestGetSettings03
//
func TestUpdateSettings03(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{DB: db}

	// Create a temp setting
	settings := db.SettingsGetOrCreateByUserId(1)

	// Change some values to test.
	settings.StrategyCcsClosePrice = 0.40
	settings.NoticeTradeFilledEmail = "break"
	settings.NoticeTradeFilledSms = "Yes"
	json, _ := json.Marshal(settings)

	// Make a mock request.
	req, _ := http.NewRequest("PUT", "/api/v1/settings", bytes.NewBuffer(json))
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(1)) })

	r.PUT("/api/v1/settings", c.UpdateSettings)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Parse json that returned.
	st.Expect(t, w.Code, 400)
	st.Expect(t, w.Body.String(), `{"errors":{"notice_trade_filled_email":"The notice_trade_filled_email must be Yes or No."}}`)
}

/* End File */
