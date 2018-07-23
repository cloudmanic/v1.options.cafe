//
// Date: 7/18/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/nbio/st"
)

//
// TestGetSettings01
//
func TestGetSettings01(t *testing.T) {

	// Json response
	jsonResponse := `{"id":1,"user_id":1,"strategy_pcs_close_price":0.03,"strategy_pcs_open_price":"mid-point","strategy_pcs_lots":10,"strategy_ccs_close_price":0.03,"strategy_ccs_open_price":"mid-point","strategy_ccs_lots":10,"strategy_pds_close_price":0.03,"strategy_pds_open_price":"mid-point","strategy_pds_lots":10,"strategy_cds_close_price":0.03,"strategy_cds_open_price":"mid-point","strategy_cds_lots":10}`

	// Start the db connection.
	db, _ := models.NewDB()

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

/* End File */
