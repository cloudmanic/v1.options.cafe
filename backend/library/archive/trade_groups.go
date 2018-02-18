//
// Date: 2/9/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package archive

import (
	"strconv"
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/cloudmanic/app.options.cafe/backend/models"
)

//
// Build / Update a Tradegoup based on an array of positions
//
func DoTradeGroupBuildFromPositions(order models.Order, positions *[]models.Position, db models.Datastore, userId uint, brokerId uint) error {

	var tradeGroupId uint
	var proceeds float64 = 0.00
	var profit float64 = 0.00
	var percentGain float64 = 0.00
	var tradeGroupStatus = "Closed"

	// If we do not have at least 1 position we give up
	if len(*positions) == 0 {
		return nil
	}

	// Get broker account id
	brokerAccount, err := db.GetBrokerAccountByBrokerAccountNumber(brokerId, order.BrokerAccountRef)

	if err != nil {
		return nil
	}

	// See if we have a trade group of any of the positions
	tradeGroupId = 0

	for _, row := range *positions {

		// Mark if this trade group is open or closed.
		if row.Qty != 0 {
			tradeGroupStatus = "Open"
		}

		if row.TradeGroupId > 0 {
			tradeGroupId = row.TradeGroupId
		}

		// Figure out group profit
		profit = profit + row.Profit

		// Figure out proceeds
		proceeds = proceeds + row.Proceeds
	}

	// Figure out what type of trade group this is.
	tgType := ClassifyTradeGroup(positions)

	// Figure out Commission
	commission := calcCommissionForOrder(&order, brokerId, &brokerAccount)

	// Figure out max risked before commissions in this trade.
	risked, credit := GetAmountRiskedInTrade(positions)

	// Create or Update Trade Group
	if tradeGroupId == 0 {

		// Update profit to have commissions. -- this should almost never be hit
		if tradeGroupStatus == "Closed" {
			profit = profit - commission
		}

		// Figure out how many trade groups we have had thus far.
		count, err := db.Count(&models.TradeGroup{}, models.QueryParam{Wheres: []models.KeyValue{{Key: "user_id", Value: strconv.Itoa(int(userId))}}})

		if err != nil {
			services.BetterError(err)
		}

		countPlus := strconv.Itoa(int(count) + 1)

		// Build a new Trade Group
		var tradeGroup = &models.TradeGroup{
			UserId:           userId,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
			Name:             "Trade #" + countPlus + " - " + tgType + " Trade",
			BrokerAccountRef: order.BrokerAccountRef,
			BrokerAccountId:  brokerAccount.Id,
			Status:           tradeGroupStatus,
			OrderIds:         strconv.Itoa(int(order.Id)),
			Commission:       commission,
			Credit:           credit,
			Proceeds:         proceeds,
			Risked:           risked,
			Profit:           profit,
			PercentGain:      percentGain,
			Type:             tgType,
			Note:             "",
			OpenDate:         order.CreateDate,
			ClosedDate:       order.TransactionDate,
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

		// Update profit to have commissions.
		if tradeGroupStatus == "Closed" {
			profit = profit - commission - tradeGroup.Commission
			percentGain = (((risked + profit) - risked) / risked) * 100
		}

		tradeGroup.Type = tgType
		tradeGroup.Proceeds = proceeds
		tradeGroup.Risked = risked
		tradeGroup.Profit = profit
		tradeGroup.PercentGain = percentGain
		tradeGroup.Credit = credit
		tradeGroup.Commission += commission
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
