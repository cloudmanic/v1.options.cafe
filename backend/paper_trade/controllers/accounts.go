//
// Date: 2018-11-05
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-11-05
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import "github.com/gin-gonic/gin"

//
// Create a new paper trading account.
//
func (t *Controller) CreateAccount(c *gin.Context) {

	// Return happy JSON
	c.JSON(200, gin.H{"status": "ok"})
}

/* End File */
