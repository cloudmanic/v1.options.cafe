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

	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/gin-gonic/gin"
)

//
// Cors middleware for local development.
//
func (t *Controller) CorsMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		// Manage OPTIONS requests
		if (os.Getenv("APP_ENV") == "local") && (c.Request.Method == http.MethodOptions) {
			c.Writer.Header().Set("Content-Type", "application/json")
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, OPTIONS")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization,DNT,X-CustomHeader,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Content-Range,Range")
			c.AbortWithStatus(200)
			return
		}

		// Set useful headers
		c.Writer.Header().Set("Access-Control-Expose-Headers", "X-Last-Page, X-Offset, X-Limit, X-No-Limit-Count")

		// On to next request in the Middleware chain.
		c.Next()
	}
}

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

		// See if this session is in our db.
		session, err := t.DB.GetByAccessToken(access_token)

		if err != nil {
			services.Critical("Access Token Not Found - Unable to Authenticate via HTTP (#002)")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization Failed (#002)"})
			c.AbortWithStatus(401)
			return
		}

		// Get this user is in our db.
		user, err := t.DB.GetUserById(session.UserId)

		if err != nil {
			services.Critical("User Not Found - Unable to Authenticate - UserId (HTTP) : " + fmt.Sprint(session.UserId) + " - Session Id : " + fmt.Sprint(session.Id))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization Failed (#003)"})
			c.AbortWithStatus(401)
			return
		}

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
