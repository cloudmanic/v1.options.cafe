//
// Date: 2018-12-23
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-12-23
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package actions

import (
	"app.options.cafe/models"
)

//
// Rebuild trade groups.
//
// go run main.go --cmd="user-rebuild-trade-groups" --user_id=1
//
func UserRebuildTradeGroups(db *models.DB, userId int) {
	// Delete all the orders / trade group stuff
	db.New().Where("user_id = ?", userId).Delete(models.TradeGroup{})
	db.New().Where("user_id = ?", userId).Delete(models.Position{})
	db.New().Where("user_id = ?", userId).Delete(models.OrderLeg{})
	db.New().Where("user_id = ?", userId).Delete(models.Order{}) // Best to do this last so things do not get out of wack
}

/* End File */
