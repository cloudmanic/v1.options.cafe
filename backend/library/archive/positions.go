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

	"github.com/cloudmanic/app.options.cafe/backend/library/helpers"
	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/cloudmanic/app.options.cafe/backend/models"
)

//
// Here we loop through all the order data and create positions. We do this because
// brokers do not offer an api of past positions.
//
func StorePositions(db models.Datastore, userId uint, brokerId uint) error {

	// Just easier to do this since we often comment stuff out for testing
	var err error

	// Process multi leg orders
	err = doMultiLegOrders(db, userId, brokerId)

	if err != nil {
		return err
	}

	// Process single option order
	err = doSingleOptionOrder(db, userId, brokerId)

	if err != nil {
		return err
	}

	// Get the different broker accounts for this user
	var results = []models.BrokerAccount{}

	err = db.Query(&results, models.QueryParam{
		Wheres: []models.KeyValue{
			{Key: "user_id", ValueInt: int(userId)},
		}})

	if err != nil {
		return err
	}

	// Just double check we do not have any expired positions. (loop through the different broker accounts)
	for _, row := range results {
		err = ReviewCurrentPositionsForExpiredOptions(db, userId, row.Id)

		if err != nil {
			return err
		}
	}

	// Return happy
	return nil
}

//
// Review current Positions for options that have expired.
//
func ReviewCurrentPositionsForExpiredOptions(db models.Datastore, userId uint, brokerId uint) error {

	// Set results var
	var results = []models.Position{}

	// Get all open positions.
	err := db.Query(&results, models.QueryParam{
		PreLoads: []string{"Symbol"},
		Wheres: []models.KeyValue{
			{Key: "status", Value: "Open"},
			{Key: "user_id", ValueInt: int(userId)},
			{Key: "broker_account_id", ValueInt: int(brokerId)},
		}})

	if err != nil {
		return err
	}

	// Set current date / time.
	now := time.Now()
	nowDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	// Loop through the positions looking for expired options.
	for _, row := range results {

		// Get the options parts.
		parts, err := helpers.OptionParse(row.Symbol.ShortName)

		// If we have an error it is not an option position.
		if err != nil {
			continue
		}

		// This value will be positive if option has expired.
		if nowDate.Sub(parts.Expire) > 0 {

			// Figure out profit.
			row.Profit = row.CostBasis * -1

			// Lets expire the position
			row.Qty = 0
			row.Proceeds = 0.00
			row.Status = "Closed"
			row.ClosedDate = parts.Expire
			db.UpdatePosition(&row)

			// Update the trade group.
			tradeGroup, err := db.GetTradeGroupById(row.TradeGroupId)

			if err != nil {
				return err
			}

			// See if this tradegroup is to close.
			var profit = 0.00
			var proceeds = 0.00
			var tradeGroupStatus = "Closed"

			// Loop through the positions
			for _, row2 := range tradeGroup.Positions {
				// Figure out profit
				profit = profit + row2.Profit

				// Figure out proceeds
				proceeds = proceeds + row.Proceeds

				// Is this trade group still open?
				if row2.Status == "Open" {
					tradeGroupStatus = "Open"
				}
			}

			// Update profit to have commissions.
			if tradeGroupStatus == "Closed" {
				tradeGroup.Note = tradeGroup.Note + " Trade Expired."
				tradeGroup.Proceeds = proceeds
				tradeGroup.Profit = profit - tradeGroup.Commission
				tradeGroup.PercentGain = (((tradeGroup.Risked + profit) - tradeGroup.Risked) / tradeGroup.Risked) * 100
				tradeGroup.Status = tradeGroupStatus
				db.UpdateTradeGroup(&tradeGroup)

				// Log success
				services.Info("New TradeGroup updated (ReviewCurrentPositionsForExpiredOptions) for user " + strconv.Itoa(int(userId)) + " TradeGroup Id: " + strconv.Itoa(int(tradeGroup.Id)))
			}
		}

	}

	// Return happy.
	return nil
}

//
// Review the order and calculate the commission for this order.
//
func CalcCommissionForOrder(order *models.Order, brokerId uint, brokerAccount *models.BrokerAccount) float64 {

	var qty = 0.00
	var commission = 0.00

	// TODO: Deal with combo orders

	// Go through the legs
	if order.Class == "multileg" {
		for _, row := range order.Legs {
			qty = qty + math.Abs(float64(row.Qty))
		}

		commission = qty * brokerAccount.OptionCommission
	} else if order.Class == "option" {
		commission = order.ExecQuantity * brokerAccount.OptionCommission
	} else {
		commission = brokerAccount.StockCommission
	}

	// See if we hit the min for multileg?
	if (order.Class == "multileg") && (brokerAccount.OptionMultiLegMin > commission) {
		commission = brokerAccount.OptionMultiLegMin
	}

	// See if we hit the min for options?
	if (order.Class == "option") && (brokerAccount.OptionSingleMin > commission) {
		commission = brokerAccount.OptionSingleMin
	}

	// Return commission value
	return commission
}

/* End File */
