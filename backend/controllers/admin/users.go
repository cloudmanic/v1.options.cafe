//
// Date: 10/2/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package admin

import (
	"net/http"

	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/gin-gonic/gin"
)

//
// Return a list of all our users.
//
func (t *Controller) GetUsers(c *gin.Context) {

	users := []models.User{}

	// List the users
	err := t.DB.Query(&users, models.QueryParam{
		Order: "id",
		Sort:  "desc",
	})

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No users found."})
		return
	}

	// Return happy JSON
	c.JSON(200, users)
}

/* End File */
