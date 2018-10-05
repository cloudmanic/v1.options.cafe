//
// Date: 10/4/2018
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

	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/nbio/st"
)

//
// CreateNotifyChannel - 01
//
func TestCreateNotifyChannel01(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()

	// Create controller
	c := &Controller{DB: db}

	// Post data
	chanPost := models.NotifyChannel{
		Type:      "Web Push",
		ChannelId: "abc123-123-324afasf-asdf",
	}

	postStr, _ := json.Marshal(chanPost)

	// Make a mock request.
	req, _ := http.NewRequest("POST", "/api/v1/notifications/add-channel", bytes.NewBuffer(postStr))
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(2)) })

	r.POST("/api/v1/notifications/add-channel", c.CreateNotifyChannel)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Make sure we have a UserId
	testObj := models.NotifyChannel{}
	db.Query(&testObj, models.QueryParam{
		Wheres: []models.KeyValue{
			{Key: "id", ValueInt: 1},
		},
	})

	// Validate result
	st.Expect(t, w.Code, 201)
	st.Expect(t, testObj.UserId, uint(2))
	st.Expect(t, w.Body.String(), `{"id":1,"type":"Web Push","channel_id":"abc123-123-324afasf-asdf"}`)
}

//
// CreateNotifyChannel - 02
//
func TestCreateNotifyChannel02(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()

	// Create controller
	c := &Controller{DB: db}

	// Post data
	chanPost := models.NotifyChannel{
		Type:      "Blah Woots",
		ChannelId: "abc123-123-324afasf-asdf",
	}

	postStr, _ := json.Marshal(chanPost)

	// Make a mock request.
	req, _ := http.NewRequest("POST", "/api/v1/notifications/add-channel", bytes.NewBuffer(postStr))
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(2)) })

	r.POST("/api/v1/notifications/add-channel", c.CreateNotifyChannel)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Validate result
	st.Expect(t, w.Code, 400)
	st.Expect(t, w.Body.String(), `{"errors":{"type":"Field channel_id must be set to Web Push."}}`)
}

//
// CreateNotifyChannel - 03
//
func TestCreateNotifyChannel03(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()

	// Create controller
	c := &Controller{DB: db}

	// Post data
	chanPost := models.NotifyChannel{
		Type:      "Web Push",
		ChannelId: "",
	}

	postStr, _ := json.Marshal(chanPost)

	// Make a mock request.
	req, _ := http.NewRequest("POST", "/api/v1/notifications/add-channel", bytes.NewBuffer(postStr))
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(2)) })

	r.POST("/api/v1/notifications/add-channel", c.CreateNotifyChannel)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Validate result
	st.Expect(t, w.Code, 400)
	st.Expect(t, w.Body.String(), `{"errors":{"channel_id":"The channel_id field is required."}}`)
}

/* End File */
