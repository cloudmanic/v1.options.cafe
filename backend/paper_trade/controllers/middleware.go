//
// Date: 11/11/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"net/http"
	"os"
	"strings"

	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/cloudmanic/app.options.cafe/backend/paper_trade/models"
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

		// See if this is a valid access token
		account := models.Account{}
		t.DB.New().Where("access_token = ?", access_token).Find(&account)

		if account.Id <= 0 {
			services.Critical("Access Token Not Found - Unable to Authenticate via HTTP (#002)")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization Failed (#002)"})
			c.AbortWithStatus(401)
			return
		}

		// Add this account to the context
		c.Set("accountId", account.Id)

		// On to next request in the Middleware chain.
		c.Next()
	}
}

/* End File */
