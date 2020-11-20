//
// Date: 2/9/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package archive

import (
	"errors"
	"math"
	"strconv"
	"time"

	"app.options.cafe/brokers/types"
	"app.options.cafe/library/helpers"
	"app.options.cafe/library/market"
	"app.options.cafe/library/services"
	"app.options.cafe/models"
)

//
// StorePositions - Here we loop through all the order data and create positions.
// We do this because brokers do not offer an api of past positions.
//
func StorePositions(db models.Datastore, userId uint, brokerId uint) error {
	// Just easier to do this since we often comment stuff out for testing
	var err error

	// Process Equity Orders
	err = DoEquityOrder(db, userId)

	if err != nil {
		return err
	}

	// Process multi leg orders
	err = DoMultiLegOrders(db, userId, brokerId)

	if err != nil {
		return err
	}

	// Process single option order
	err = DoSingleOptionOrder(db, userId, brokerId)

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
			services.InfoMsg("archive.ReviewCurrentPositionsForExpiredOptions returned error " + strconv.Itoa(int(userId)) + " Broker Id: " + strconv.Itoa(int(brokerId)) + " Broker Account Id: " + strconv.Itoa(int(row.Id)))
			return err
		}
	}

	// Return happy
	return nil
}

//
// Create a trade group from a position.
// Since we only import orders post 2017 from Tradier
// we could have a position that get added before 2017 and not closed
// this is how we catch that and create a trade group for it.
// We are not doing any trade classifying.
//
func PastCreateTradeGroupFromPosition(db models.Datastore, userId uint, brokerId uint, pos types.Position) error {

	// Get the symbol we are after (or create if it is not in the system)
	symbol, err := db.GetSymbolByShortName(pos.Symbol)

	if err != nil {
		return err
	}

	// Get the broker account id.
	brokerAccount, err := db.GetBrokerAccountByBrokerAccountNumber(brokerId, pos.AccountId)

	if err != nil {
		return err
	}

	// Check to see if we already have an open position
	_, err2 := db.GetPositionByUserSymbolStatusAccount(userId, symbol.Id, "Open", brokerAccount.Id)

	if (err2 != nil) && (err2.Error() == "Record not found") {

		// Create trade group
		tg := models.TradeGroup{
			UserId:           userId,
			Name:             "Past " + symbol.Type + " Trade",
			BrokerAccountId:  brokerAccount.Id,
			BrokerAccountRef: pos.AccountId,
			Status:           "Open",
			Type:             symbol.Type,
			Risked:           pos.CostBasis,
			Note:             "Order opened before 2018, so no order history imported. Please verify details of this trade.",
			Commission:       brokerAccount.StockCommission,
			OpenDate:         helpers.ParseDateNoError(pos.OpenDate),
		}

		err = db.CreateTradeGroup(&tg)

		if err != nil {
			return err
		}

		// Create position
		newPos := models.Position{
			UserId:           userId,
			TradeGroupId:     tg.Id,
			BrokerAccountId:  brokerAccount.Id,
			BrokerAccountRef: pos.AccountId,
			Status:           "Open",
			SymbolId:         symbol.Id,
			Qty:              int(pos.Quantity),
			OrgQty:           int(pos.Quantity),
			CostBasis:        pos.CostBasis,
			AvgOpenPrice:     (pos.CostBasis / pos.Quantity),
			OpenDate:         helpers.ParseDateNoError(pos.OpenDate),
		}

		db.CreatePosition(&newPos)

		// Log success
		services.InfoMsg("New Past TradeGroup created for user " + strconv.Itoa(int(userId)) + " TradeGroup Id: " + strconv.Itoa(int(tg.Id)))
	}

	// Return Happy
	return nil
}

//
// ReviewCurrentPositionsForExpiredOptions - Review current Positions for options that have expired.
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

			// Get underlying quote
			stockQuote, err := market.GetUnderlayingQuoteByDate(db, userId, row.Symbol.OptionUnderlying, parts.Expire)

			if err != nil {
				services.Critical(errors.New(err.Error() + "market.GetUnderlayingQuoteByDate returned error UserId: " + strconv.Itoa(int(userId)) + " Broker Id: " + strconv.Itoa(int(brokerId)) + " OptionUnderlying: " + row.Symbol.OptionUnderlying + " Today: " + parts.Expire.Format("2006-01-02")))

				// This is hacky but if we can't get a quote we need to do something.
				// Common issue is Tradier does not give historical quotes on futures.
				stockQuote = 0.00

				// TODO(spicer): Notify user or something. (we put a note in the tradegroup below)
			}

			// Figure out profit.
			if row.OrgQty > 0 {

				// Figure out based on Call or Put - Long
				if (row.Symbol.OptionType == "Call") && (row.Symbol.OptionStrike < stockQuote) {
					row.Profit = (stockQuote - row.Symbol.OptionStrike) * float64(row.OrgQty) * 100.00
				} else if (row.Symbol.OptionType == "Put") && (row.Symbol.OptionStrike > stockQuote) {
					row.Profit = ((row.Symbol.OptionStrike - stockQuote) * float64(row.OrgQty) * 100.00) - row.CostBasis
				} else {
					row.Profit = row.CostBasis * -1
				}

			} else {

				// Figure out based on Call or Put - Short
				if (row.Symbol.OptionType == "Call") && (row.Symbol.OptionStrike < stockQuote) {
					row.Profit = ((stockQuote - row.Symbol.OptionStrike) * float64(row.OrgQty) * 100.00) + (row.CostBasis * -1)
				} else if (row.Symbol.OptionType == "Put") && (row.Symbol.OptionStrike > stockQuote) {
					row.Profit = ((row.Symbol.OptionStrike - stockQuote) * float64(row.OrgQty) * 100.00) + (row.CostBasis * -1)
				} else {
					row.Profit = row.CostBasis * -1
				}

			}

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
				tradeGroup.ClosedDate = tradeGroup.Positions[0].ClosedDate

				// Poor Tradier note
				if stockQuote == 0 {
					tradeGroup.Note = tradeGroup.Note + " Due to Tradier's poor historical data it is expected the trade data here is not correct. Please contact help@options.cafe to update."
				}

				db.UpdateTradeGroup(&tradeGroup)

				// Log success
				services.InfoMsg("New TradeGroup updated (ReviewCurrentPositionsForExpiredOptions) for user " + strconv.Itoa(int(userId)) + " TradeGroup Id: " + strconv.Itoa(int(tradeGroup.Id)))
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
