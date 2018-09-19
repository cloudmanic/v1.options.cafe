//
// Date: 9/18/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package archive

import (
	"strconv"
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/library/notify"
	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/cloudmanic/app.options.cafe/backend/models"
)

//
// Do equity order
//
func DoEquityOrder(db models.Datastore, userId uint) error {

	// Query and get all orders we have not reviewed before.
	orders, err := db.GetOrdersByUserClassStatusReviewed(userId, "equity", "filled", "No")

	if err != nil {
		return nil
	}

	// Loop through the different orders and process.
	for _, row := range orders {

		var positions []models.Position
		pos, err := doOpenEquityOrder(row, db, userId)

		if err != nil {
			return nil
		}

		// Build positions
		positions = append(positions, pos)

		// Build Trade Group
		err = DoTradeGroupBuildFromPositions(row, &positions, db, userId, row.BrokerAccountId)

		if err != nil {
			services.Fatal(err)
			continue
		}

		// Mark the order as reviewed
		row.PositionReviewed = "Yes"
		err = db.UpdateOrder(&row)

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
// Do a single equity order - Open Order
//
func doOpenEquityOrder(order models.Order, db models.Datastore, userId uint) (models.Position, error) {

	var qty int = 0
	var cost_basis float64 = 0.00

	// First we find out if we already have a position on for this.
	position, _ := db.GetPositionByUserSymbolStatusAccount(userId, order.SymbolId, "Open", order.BrokerAccountId)

	// Is this a long trade closing?
	if order.Side == "buy" {
		qty = int(order.ExecQuantity)
		cost_basis = (float64(order.ExecQuantity) * order.AvgFillPrice)
	}

	// Is this a short trade closing?
	if order.Side == "sell" {
		qty = (int(order.ExecQuantity) * -1)
	}

	// We found so we are just adding to a current position.
	if position.Id > 0 {

		// Update pos
		position.CostBasis = position.CostBasis + cost_basis
		position.OrderIds = position.OrderIds + "," + strconv.Itoa(int(order.Id))
		position.UpdatedAt = time.Now()
		position.Qty = qty + position.Qty

		if order.Side == "buy" {
			position.OrgQty = qty + position.OrgQty
			position.AvgOpenPrice = ((order.AvgFillPrice + position.AvgOpenPrice) / 2)
		}

		position.Note = position.Note + "Updated - " + order.TransactionDate.Format(time.RFC1123) + " :: "
		db.UpdatePosition(&position)

		// Are we done with the trade?
		if position.Qty == 0 {
			position.Proceeds = position.Proceeds + (order.ExecQuantity * order.AvgFillPrice)
			position.Profit = position.Profit + (position.Proceeds - position.CostBasis)
			position.ClosedDate = order.TransactionDate
			position.Status = "Closed"
		}

	} else {

		// Insert Position
		position = models.Position{
			UserId:           userId,
			TradeGroupId:     0,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
			BrokerAccountId:  uint(order.BrokerAccountId),
			BrokerAccountRef: order.BrokerAccountRef,
			SymbolId:         order.SymbolId,
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

/* End File */
