//
// Date: 2018-11-05
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-11-18
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"github.com/gin-gonic/gin"
)

//
// Do Routes
//
func (t *Controller) DoRoutes(r *gin.Engine) {

	// --------- API V1 sub-routes ----------- //

	apiV1 := r.Group("/api/v1")

	apiV1.Use(t.AuthMiddleware())
	{
		// accounts
		apiV1.GET("/accounts/me", t.GetAccountMe)
		apiV1.POST("/accounts", t.CreateAccount)

		// quotes
		apiV1.GET("/quotes", t.GetQuotes)
	}

}

/* End File */
