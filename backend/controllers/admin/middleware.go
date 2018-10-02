//
// Date: 12/2/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package admin

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

// IPs allowed to access these admin routes
var allowedIps = map[string]bool{
	"127.0.0.1":     true,
	"71.238.46.175": true, // Spicer home
	"96.239.59.69":  true, // VMG VPN
}

//
// Here we make sure we passed in a proper Bearer Access Token.
//
func (t *Controller) AuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		// We only allow this request from a few IP addresses
		if _, ok := allowedIps[realip.RealIP(c.Request)]; !ok {
			services.Critical("UnAuthorization IP address. - " + realip.RealIP(c.Request))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization Failed (#005)"})
			c.AbortWithStatus(401)
			return
		}

		// Set access token and start the auth process
		var access_token = ""

		// Make sure we have a Bearer token.
		auth := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)

		if len(auth) != 2 || auth[0] != "Bearer" {

			// We allow access token from the url
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

		// Make sure this is an admin user (Options Cafe Employee)
		if user.Admin != "Yes" {
			services.Critical("User Not Found - Unable to Authenticate - UserId (HTTP) : " + fmt.Sprint(session.UserId) + " - Session Id : " + fmt.Sprint(session.Id))
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
