//
// Date: 2/9/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package archive

import (
	"math"
	"strconv"
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/cloudmanic/app.options.cafe/backend/models"
)

//
// Here we loop through all the order data and create positions. We do this because
// brokers do not offer an api of past positions.
//
func StorePositions(db models.Datastore, userId uint, brokerId uint) error {

	// Process multi leg orders
	err := doMultiLegOrders(db, userId, brokerId)

	if err != nil {
		return err
	}

	// Process single option order
	err = doSingleOptionOrder(db, userId, brokerId)

	if err != nil {
		return err
	}

	// Return happy
	return nil
}

//
// Build / Update a Tradegoup based on an array of positions
//
func doTradeGroupBuildFromPositions(order models.Order, positions *[]models.Position, db models.Datastore, userId uint, brokerId uint) error {

	var totalQty float64
	var tradeGroupId uint
	var tradeGroupStatus = "Closed"

	// If we do not have at least 1 position we give up
	if len(*positions) == 0 {
		return nil
	}

	// Get broker account id
	brokerAccount, err := db.GetBrokerAccountByBrokerAccountNumber(brokerId, order.AccountId)

	if err != nil {
		return err
	}

	// See if we have a trade group of any of the positions
	totalQty = 0
	tradeGroupId = 0

	for _, row := range *positions {

		// Mark if this trade group is open or closed.
		if row.Qty != 0 {
			tradeGroupStatus = "Open"
		}

		if row.TradeGroupId > 0 {
			tradeGroupId = row.TradeGroupId
		}

		// Get the total qty
		// TODO: There is a bug here. If you start the position with more than one order where your paying the
		// min commission more than once this number will not be correct.
		totalQty = totalQty + math.Abs(float64(row.OrgQty))
	}

	// Figure out what type of trade group this is.
	tgType := ClassifyTradeGroup(positions)

	// Figure out Commission
	commission := (totalQty * .35)

	// See if we hit the min for multileg?
	if (order.Class == "multileg") && (brokerAccount.OptionMultiLegMin > commission) {
		commission = brokerAccount.OptionMultiLegMin
	}

	// See if we hit the min for options?
	if (order.Class == "option") && (brokerAccount.OptionSingleMin > commission) {
		commission = brokerAccount.OptionSingleMin
	}

	// Closing order?
	if tradeGroupStatus == "Closed" {
		commission = commission * 2
	}

	// TODO: Figure out Risked, Gain, and Profit (if this is closed)

	// Create or Update Trade Group
	if tradeGroupId == 0 {

		// Build a new Trade Group
		var tradeGroup = &models.TradeGroup{
			UserId:          userId,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
			AccountId:       order.AccountId,
			BrokerAccountId: brokerAccount.Id,
			Status:          tradeGroupStatus,
			OrderIds:        strconv.Itoa(int(order.Id)),
			Commission:      commission,
			Type:            tgType,
			Note:            "",
			OpenDate:        order.CreateDate,
			ClosedDate:      order.TransactionDate,
		}

		// Insert into DB
		db.CreateTradeGroup(tradeGroup)

		// Store tradegroup id
		tradeGroupId = tradeGroup.Id

		// Log success
		services.Info("New TradeGroup created for user " + strconv.Itoa(int(userId)) + " TradeGroup Id: " + strconv.Itoa(int(tradeGroupId)))
	} else {

		// Update tradegroup with additional OrderIds
		tradeGroup, err := db.GetTradeGroupById(tradeGroupId)

		if err != nil {
			return err
		}

		tradeGroup.Type = tgType
		tradeGroup.Commission = commission
		tradeGroup.Status = tradeGroupStatus
		tradeGroup.ClosedDate = order.TransactionDate
		tradeGroup.OrderIds = tradeGroup.OrderIds + "," + strconv.Itoa(int(order.Id))
		db.UpdateTradeGroup(&tradeGroup)

		// Log success
		services.Info("New TradeGroup updated for user " + strconv.Itoa(int(userId)) + " TradeGroup Id: " + strconv.Itoa(int(tradeGroupId)))
	}

	// Loop through the positions and add the trade group id
	for _, row := range *positions {
		row.TradeGroupId = tradeGroupId
		db.UpdatePosition(&row)
	}

	// Return happy.
	return nil
}

/* End File */
