//
// Date: 10/2/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package admin

import (
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/cloudmanic/app.options.cafe/backend/library/realip"
	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

//
// Delete a user's account. User will be deleted forever.
//
func (t *Controller) DeleteUser(c *gin.Context) {

	// Set as int
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Get User by id.
	user, err := t.DB.GetUserById(uint(id))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Delete a user
	err = t.DB.DeleteUser(&user)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Return happy
	c.JSON(204, nil)
}

//
// Login as a user (remember only admins can do this)
//
func (t *Controller) LoginAsUser(c *gin.Context) {

	// Parse json body
	body, err := ioutil.ReadAll(c.Request.Body)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// First we validate the grant type and client id. Make sure this is a known application.
	app, err := t.DB.ValidateClientIdGrantType(os.Getenv("CENTCOM_CLIENT_ID"), "password")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client_id or grant type."})
		return
	}

	// Get user Id
	userId := gjson.Get(string(body), "id").Int()

	user, err := t.DB.LoginUserById(uint(userId), app.Id, c.Request.UserAgent(), realip.RealIP(c.Request))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Return success json.
	c.JSON(200, gin.H{"access_token": user.Session.AccessToken, "user_id": user.Id})
}

//
// Return a list of all our users.
//
func (t *Controller) GetUsers(c *gin.Context) {

	users := []models.User{}

	// List the users
	err := t.DB.Query(&users, models.QueryParam{
		Order: "id",
		Sort:  "desc",
	})

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No users found."})
		return
	}

	// Return happy JSON
	c.JSON(200, users)
}

/* End File */
