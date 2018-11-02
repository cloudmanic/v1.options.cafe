//
// Date: 10/2/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package admin

import (
	"github.com/gin-gonic/gin"
)

//
// Do Routes
//
func (t *Controller) DoRoutes(r *gin.Engine) {

	// ------------- Admin API --------------- //

	adminApi := r.Group("/api/admin")

	adminApi.Use(t.AuthMiddleware())
	{
		// Ping
		adminApi.GET("/ping", t.PingFromServer)

		// Users
		adminApi.GET("/users", t.GetUsers)
		adminApi.DELETE("/users/:id", t.DeleteUser)
		adminApi.POST("/users/login-as-user", t.LoginAsUser)
	}

}

/* End File */
