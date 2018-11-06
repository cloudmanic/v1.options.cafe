//
// Date: 2018-11-05
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-11-05
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

	apiV1.Use()
	{
		// accounts
		apiV1.POST("/accounts", t.CreateAccount)

	}

}

/* End File */
