//
// Date: 9/14/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package admin

import (
	"github.com/gin-gonic/gin"
)

//
// Collect a ping from the server to know we are alive.
// Mainly used as a way to see if a user is logged in an allowed
// to access the admin side of the application.
//
func (t *Controller) PingFromServer(c *gin.Context) {

	// Make sure the UserId is correct.
	//userId := c.MustGet("userId").(uint)

	// Return happy JSON
	c.JSON(200, gin.H{"status": "ok"})
}

/* End File */
