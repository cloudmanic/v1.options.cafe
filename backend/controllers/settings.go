//
// Date: 7/18/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"github.com/gin-gonic/gin"
)

//
// Get a settings by UserId.
//
func (t *Controller) GetSettings(c *gin.Context) {

	// Get the user id.
	userId := c.MustGet("userId").(uint)

	// Get the settings by id.
	settings := t.DB.SettingsGetOrCreateByUserId(userId)

	// Return happy JSON
	c.JSON(200, settings)
}

/* End File */
