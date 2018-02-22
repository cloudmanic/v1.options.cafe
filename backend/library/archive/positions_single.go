//
// Date: 2/9/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package archive

import (
	"errors"
	"strconv"
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/library/notify"
	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/cloudmanic/app.options.cafe/backend/models"
)

//
// Do single option order
//
func doSingleOptionOrder(db models.Datastore, userId uint, brokerId uint) error {

	// Query and get all orders we have not reviewed before.
	orders, err := db.GetOrdersByUserClassStatusReviewed(userId, "option", "filled", "No")

	if err != nil {
		return nil
	}

	// Loop through the different orders and process.
	for _, row := range orders {

		var positions []models.Position
		var loopErr error

		// Deal with sides
		switch row.Side {

		case "sell_to_open":
			row.Qty = (row.Qty * -1)
			pos, err := doOpenSingleOptionOrder(row, db, userId)
			positions = append(positions, pos)
			loopErr = err

		case "buy_to_open":
			pos, err := doOpenSingleOptionOrder(row, db, userId)
			positions = append(positions, pos)
			loopErr = err

		case "buy_to_close":
			pos, err := doCloseSingleOptionOrder(row, db, userId)
			positions = append(positions, pos)
			loopErr = err

		case "sell_to_close":
			row.Qty = (row.Qty * -1)
			pos, err := doCloseSingleOptionOrder(row, db, userId)
			positions = append(positions, pos)
			loopErr = err

		default:
			services.Critical("doSingleOptionOrder() : Unknown Side")
			loopErr = errors.New("doSingleOptionOrder() : Unknown Side")
		}

		// Did we have an err
		if loopErr != nil {
			services.BetterError(loopErr)

			// Mark the order as reviewed
			row.PositionReviewed = "Error"
			err := db.UpdateOrder(&row)

			if err != nil {
				services.Fatal(err)
			}

			continue
		}

		// Build Trade Group
		err = DoTradeGroupBuildFromPositions(row, &positions, db, userId, brokerId)

		// Mark the order as reviewed
		row.PositionReviewed = "Yes"
		err := db.UpdateOrder(&row)

		if err != nil {
			services.Fatal(err)
			continue
		}

		// Notify
		notify.PushWebsocket(userId, "change-detected", `{ "type": "order-filled", "id": `+strconv.Itoa(int(row.Id))+` }`)
	}

	// Return happy
	return nil
}

//
// Do a single option order - Open Order
//
func doOpenSingleOptionOrder(order models.Order, db models.Datastore, userId uint) (models.Position, error) {

	var qty int = 0
	var cost_basis float64 = 0.00

	// First we find out if we already have a position on for this.
	position, _ := db.GetPositionByUserSymbolStatusAccount(userId, order.OptionSymbolId, "Open", order.BrokerAccountId)

	// Is this a long trade closing?
	if order.Side == "buy_to_open" {
		qty = int(order.ExecQuantity)
		cost_basis = (float64(order.ExecQuantity) * order.AvgFillPrice * 100)
	}

	// Is this a short trade closing?
	if order.Side == "sell_to_open" {
		qty = (int(order.ExecQuantity) * -1)
		cost_basis = ((float64(order.ExecQuantity) * order.AvgFillPrice * 100) * -1)
	}

	// We found so we are just adding to a current position.
	if position.Id > 0 {

		// Update pos
		position.CostBasis = position.CostBasis + cost_basis
		position.OrderIds = position.OrderIds + "," + strconv.Itoa(int(order.Id))
		position.UpdatedAt = time.Now()
		position.Qty = qty + position.Qty
		position.OrgQty = qty + position.OrgQty
		position.AvgOpenPrice = ((order.AvgFillPrice + position.AvgOpenPrice) / 2)
		position.Note = position.Note + "Updated - " + order.TransactionDate.Format(time.RFC1123) + " :: "
		db.UpdatePosition(&position)

	} else {

		// Insert Position
		position = models.Position{
			UserId:           userId,
			TradeGroupId:     0,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
			BrokerAccountId:  uint(order.BrokerAccountId),
			BrokerAccountRef: order.BrokerAccountRef,
			SymbolId:         order.OptionSymbolId,
			Qty:              qty,
			OrgQty:           qty,
			CostBasis:        cost_basis,
			AvgOpenPrice:     order.AvgFillPrice,
			AvgClosePrice:    0.00,
			Note:             "",
			OpenDate:         order.CreateDate,
			OrderIds:         strconv.Itoa(int(order.Id)),
			Status:           "Open",
		}

		// Insert into DB
		db.CreatePosition(&position)
	}

	// Return a list of position that we reviewed
	return position, nil
}

//
// Do a single option order - Close Order
//
func doCloseSingleOptionOrder(order models.Order, db models.Datastore, userId uint) (models.Position, error) {

	// First we find out if we already have a position on for this.
	position, _ := db.GetPositionByUserSymbolStatusAccount(userId, order.OptionSymbolId, "Open", order.BrokerAccountId)

	// We found so we are just removing to a current position.
	if position.Id > 0 {

		// Update pos
		position.OrderIds = position.OrderIds + "," + strconv.Itoa(int(order.Id))
		position.UpdatedAt = time.Now()

		// Is this a long trade closing?
		if order.Side == "buy_to_close" {
			position.Qty = position.Qty + int(order.ExecQuantity)
			position.Proceeds = position.Proceeds - (order.ExecQuantity * order.AvgFillPrice * 100)
			position.Profit = position.Profit + (position.Proceeds - position.CostBasis)
		}

		// Is this a short trade closing?
		if order.Side == "sell_to_close" {
			position.Qty = position.Qty - int(order.ExecQuantity)
			position.Proceeds = position.Proceeds + (order.ExecQuantity * order.AvgFillPrice * 100)
			position.Profit = position.Profit + (position.Proceeds - position.CostBasis)
		}

		// Average Close Price
		if position.AvgClosePrice != 0 {
			position.AvgClosePrice = ((order.AvgFillPrice + position.AvgClosePrice) / 2)
		} else {
			position.AvgClosePrice = order.AvgFillPrice
		}

		// Are we done with the trade?
		if position.Qty == 0 {
			position.ClosedDate = order.TransactionDate
			position.Status = "Closed"
		} else {
			position.Note = position.Note + "Updated - " + order.TransactionDate.Format(time.RFC1123) + " :: "
		}

		db.UpdatePosition(&position)

	} else {
		return models.Position{}, errors.New("Unable to find close position in our database. - " + strconv.Itoa(int(userId)) + " : " + order.BrokerAccountRef)
	}

	// Return a list of position that we reviewed
	return position, nil
}

/* End File */
