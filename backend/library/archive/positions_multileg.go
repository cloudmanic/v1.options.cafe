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

	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/cloudmanic/app.options.cafe/backend/models"
)

//
// Do multi leg orders - Just when you open a position.
//
func doMultiLegOrders(db models.Datastore, userId uint, brokerId uint) error {

	// Query and get all orders we have not reviewed before.
	orders, err := db.GetOrdersByUserClassStatusReviewed(userId, "multileg", "filled", "No")

	if err != nil {
		return nil
	}

	// Loop through the different orders and process.
	for _, row := range orders {

		var loopErr error
		var positions []models.Position

		// Loop through the legs and store
		for _, row2 := range row.Legs {

			// Deal with sides
			switch row2.Side {

			case "sell_to_open":
				row2.Qty = (row2.Qty * -1)
				pos, err := doOpenOneLegMultiLegOrder(row, row2, db, userId)
				positions = append(positions, pos)
				loopErr = err

			case "buy_to_open":
				pos, err := doOpenOneLegMultiLegOrder(row, row2, db, userId)
				positions = append(positions, pos)
				loopErr = err

			case "buy_to_close":
				pos, err := doCloseOneLegMultiLegOrder(row, row2, db, userId)
				positions = append(positions, pos)
				loopErr = err

			case "sell_to_close":
				row2.Qty = (row2.Qty * -1)
				pos, err := doCloseOneLegMultiLegOrder(row, row2, db, userId)
				positions = append(positions, pos)
				loopErr = err

			default:
				services.Critical("doMultiLegOrders() : Unknown Side")
				loopErr = errors.New("doMultiLegOrders() : Unknown Side")
			}

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
		err = doTradeGroupBuildFromPositions(row, &positions, db, userId, brokerId)

		if err != nil {
			services.Fatal(err)
			continue
		}

		// Mark the order as reviewed
		row.PositionReviewed = "Yes"
		err := db.UpdateOrder(&row)

		if err != nil {
			services.Fatal(err)
			continue
		}

	}

	// Return happy
	return nil
}

//
// Do one leg of a multi leg order - Open Order
//
func doOpenOneLegMultiLegOrder(order models.Order, leg models.OrderLeg, db models.Datastore, userId uint) (models.Position, error) {

	var qty int = 0
	var cost_basis float64 = 0.00

	// First we find out if we already have a position on for this.
	position, _ := db.GetPositionByUserSymbolStatusAccount(userId, leg.SymbolId, "Open", order.AccountId)

	// Is this a long trade closing?
	if leg.Side == "buy_to_open" {
		qty = int(leg.ExecQuantity)
		cost_basis = (float64(leg.ExecQuantity) * leg.AvgFillPrice * 100)
	}

	// Is this a short trade closing?
	if leg.Side == "sell_to_open" {
		qty = (int(leg.ExecQuantity) * -1)
		cost_basis = ((float64(leg.ExecQuantity) * leg.AvgFillPrice * 100) * -1)
	}

	// We found so we are just adding to a current position.
	if position.Id > 0 {

		// Update pos
		position.OrderIds = position.OrderIds + "," + strconv.Itoa(int(order.Id))
		position.UpdatedAt = time.Now()
		position.CostBasis = position.CostBasis + cost_basis
		position.Qty = qty + position.Qty
		position.OrgQty = qty + position.OrgQty
		position.AvgOpenPrice = (((leg.AvgFillPrice + position.AvgOpenPrice) / 2) * 100)
		position.Note = position.Note + "Updated - " + leg.TransactionDate.Format(time.RFC1123) + " :: "
		db.UpdatePosition(&position)

	} else {

		// Insert Position
		position = models.Position{
			UserId:        userId,
			TradeGroupId:  0,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			AccountId:     order.AccountId,
			SymbolId:      leg.SymbolId,
			Qty:           qty,
			OrgQty:        qty,
			CostBasis:     cost_basis,
			Proceeds:      0.00,
			AvgOpenPrice:  leg.AvgFillPrice,
			AvgClosePrice: 0.00,
			Note:          "",
			OpenDate:      leg.CreateDate,
			OrderIds:      strconv.Itoa(int(order.Id)),
			Status:        "Open",
		}

		// Insert into DB
		db.CreatePosition(&position)
	}

	// Return a list of position that we reviewed
	return position, nil
}

//
// Do one leg of a multi leg order - Close Order
//
func doCloseOneLegMultiLegOrder(order models.Order, leg models.OrderLeg, db models.Datastore, userId uint) (models.Position, error) {

	// First we find out if we already have a position on for this.
	position, _ := db.GetPositionByUserSymbolStatusAccount(userId, leg.SymbolId, "Open", order.AccountId)

	// We found so we are just removing to a current position.
	if position.Id > 0 {

		// Update pos
		position.OrderIds = position.OrderIds + "," + strconv.Itoa(int(order.Id))
		position.UpdatedAt = time.Now()

		// Is this a long trade closing?
		if leg.Side == "buy_to_close" {
			position.Qty = position.Qty + int(leg.ExecQuantity)
			position.Proceeds = position.Proceeds - (leg.ExecQuantity * leg.AvgFillPrice * 100)
			position.Profit = position.Profit + (position.Proceeds - position.CostBasis)
		}

		// Is this a short trade closing?
		if leg.Side == "sell_to_close" {
			position.Qty = position.Qty - int(leg.ExecQuantity)
			position.Proceeds = position.Proceeds + (leg.ExecQuantity * leg.AvgFillPrice * 100)
			position.Profit = position.Profit + (position.Proceeds - position.CostBasis)
		}

		// Average Close Price
		if position.AvgClosePrice != 0 {
			position.AvgClosePrice = ((leg.AvgFillPrice + position.AvgClosePrice) / 2)
		} else {
			position.AvgClosePrice = leg.AvgFillPrice
		}

		// Are we done with the trade?
		if position.Qty == 0 {
			position.ClosedDate = leg.TransactionDate
			position.Status = "Closed"
		} else {
			position.Note = position.Note + "Updated - " + leg.TransactionDate.Format(time.RFC1123) + " :: "
		}

		db.UpdatePosition(&position)

	} else {
		return models.Position{}, errors.New("Unable to find close position in our database. - " + strconv.Itoa(int(userId)) + " : " + strconv.Itoa(int(leg.SymbolId)) + " : " + order.AccountId)
	}

	// Return a list of position that we reviewed
	return position, nil

}

/* End File */
