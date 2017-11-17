//
// Date: 9/8/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/app.options.cafe/backend/library/realip"
	"github.com/app.options.cafe/backend/library/services"
	"github.com/gin-gonic/gin"
)

//
// Register a new account.
//
func (t *Controller) DoRegister(c *gin.Context) {

	// Set response
	if os.Getenv("APP_ENV") == "local" {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	}

	// Decode json passed in
	decoder := json.NewDecoder(c.Request.Body)

	type RegisterPost struct {
		First    string
		Last     string
		Email    string
		Password string
	}

	var post RegisterPost

	err := decoder.Decode(&post)

	if err != nil {
		services.Error(err, "DoRegisterPost - Failed to decode JSON posted in")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Something went wrong while logging into your account. Please try again or contact help@options.cafe. Sorry for the trouble."})
		return
	}

	defer c.Request.Body.Close()

	// Validate user.
	if err := t.DB.ValidateCreateUser(post.First, post.Last, post.Email, post.Password); err != nil {

		// Respond with error
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Install new user.
	user, err := t.DB.CreateUser(post.First, post.Last, post.Email, post.Password, c.Request.UserAgent(), realip.RealIP(c.Request))

	if err != nil {
		services.Error(err, "DoRegisterPost - Unable to register new user. (CreateUser)")

		// Respond with error
		c.JSON(http.StatusBadRequest, gin.H{"error": "Something went wrong while logging into your account. Please try again or contact help@options.cafe. Sorry for the trouble."})
		return
	}

	type Response struct {
		Status      uint   `json:"status"`
		UserId      uint   `json:"user_id"`
		AccessToken string `json:"access_token"`
		BrokerCount int    `json:"broker_count"`
	}

	resObj := &Response{
		Status:      1,
		UserId:      user.Id,
		AccessToken: user.Session.AccessToken,
		BrokerCount: 0,
	}

	// Return success json.
	c.JSON(200, resObj)
}

/* End File */
