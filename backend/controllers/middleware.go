//
// Date: 11/11/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/library/realip"
	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/gin-gonic/gin"
)

//
// Here we make sure we passed in a proper Bearer Access Token.
//
func (t *Controller) AuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		// Set access token and start the auth process
		var access_token = ""

		// Make sure we have a Bearer token.
		auth := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)

		if len(auth) != 2 || auth[0] != "Bearer" {

			// We allow access token from the command line
			if os.Getenv("APP_ENV") == "local" {

				access_token = c.Query("access_token")

				if len(access_token) <= 0 {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization Failed (#101)"})
					c.AbortWithStatus(401)
					return
				}

			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization Failed (#001)"})
				c.AbortWithStatus(401)
				return
			}

		} else {
			access_token = auth[1]
		}

		// If access_token is empty do nothing
		if len(access_token) <= 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization Failed (#004)"})
			c.AbortWithStatus(401)
			return
		}

		// See if this session is in our db.
		session, err := t.DB.GetByAccessToken(access_token)

		if err != nil {
			services.InfoMsg("Access Token Not Found - Unable to Authenticate via HTTP (#002) - " + access_token)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization Failed (#002)"})
			c.AbortWithStatus(401)
			return
		}

		// Get this user is in our db.
		user, err := t.DB.GetUserById(session.UserId)

		if err != nil {
			services.InfoMsg("User Not Found - Unable to Authenticate - UserId (HTTP) : " + fmt.Sprint(session.UserId) + " - Session Id : " + fmt.Sprint(session.Id))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization Failed (#003)"})
			c.AbortWithStatus(401)
			return
		}

		// Make sure we really have a user.
		if user.Id <= 0 {
			services.InfoMsg("User Not Found - Unable to Authenticate - UserId <= 0 (HTTP) : " + fmt.Sprint(session.UserId) + " - Session Id : " + fmt.Sprint(session.Id))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization Failed (#004)"})
			c.AbortWithStatus(401)
			return
		}

		// Add last_activity to users model as well.
		user.LastActivity = time.Now()
		t.DB.UpdateUser(&user)

		// Log this request into the last_activity col.
		session.LastActivity = time.Now()
		session.LastIpAddress = realip.RealIP(c.Request)
		t.DB.UpdateSession(&session)

		// Add this user to the context
		c.Set("userId", user.Id)

		// CORS for local development.
		if os.Getenv("APP_ENV") == "local" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		}

		// On to next request in the Middleware chain.
		c.Next()
	}
}

/* End File */
