package archive

import (
	"strconv"
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/cloudmanic/app.options.cafe/backend/models"
)

//
// Here we loop through all the order data and create positions. We do this because
// brokers do not offer an api of past positions.
//
func StorePositions(db models.Datastore, userId uint) error {

	// Process different orders types.
	doMultiLegOrders(db, userId)

	// Return happy
	return nil

}

// //
// //
// //
// func ClassifyTradeGroup(positions *[]models.Position) {

// 	for _, row := range *positions {

// 		fmt.Println(row.Symbol)

// 		/*
// 		   s := []string{row.Symbol}

// 		   var q tradier.Quote.Api

// 		   quotes := q.GetQuotes(s)

// 		   fmt.Println(quotes)
// 		*/

// 	}

// }

//
// Do multi leg orders - Just when you open a position.
//
func doMultiLegOrders(db models.Datastore, userId uint) error {

	// Query and get all orders we have not reviewed before.
	orders, err := db.GetOrdersByUserClassStatusReviewed(userId, "multileg", "filled", "No")

	if err != nil {
		return err
	}

	// Loop through the different orders and process.
	for _, row := range orders {

		var positions []*models.Position

		// Loop through the legs and store
		for _, row2 := range row.Legs {

			// Deal with sides
			switch row2.Side {

			case "sell_to_open":
				row2.Qty = (row2.Qty * -1)

				pos, err := doOpenOneLegMultiLegOrder(row, row2, db, userId)

				if err != nil {
					services.Fatal(err)
					continue
				}

				positions = append(positions, pos)

			case "buy_to_open":
				pos, err := doOpenOneLegMultiLegOrder(row, row2, db, userId)

				if err != nil {
					services.Fatal(err)
					continue
				}

				positions = append(positions, pos)

			// case "buy_to_close":
			// 	pos, err := doCloseOneLegMultiLegOrder(row, row2, db, userId)

			// 	if err != nil {
			// 		services.Fatal(err)
			// 		continue
			// 	}

			// 	positions = append(positions, pos)

			// case "sell_to_close":
			// 	row2.Qty = (row2.Qty * -1)

			// 	pos, err := doCloseOneLegMultiLegOrder(row, row2, db, userId)

			// 	if err != nil {
			// 		services.Fatal(err)
			// 		continue
			// 	}

			// 	positions = append(positions, pos)

			default:
				services.Critical("doMultiLegOrders() : Unknown Side")
			}

		}

		// Build Trade Group
		err = doTradeGroupBuildFromPositions(row, positions, db, userId)

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
// Build / Update a Tradegoup based on an array of positions
//
func doTradeGroupBuildFromPositions(order models.Order, positions []*models.Position, db models.Datastore, userId uint) error {

	var tradeGroupId uint
	var tradeGroupStatus = "Closed"

	// If we do not have at least 1 position we give up
	if len(positions) == 0 {
		return nil
	}

	// See if we have a trade group of any of the positions
	tradeGroupId = 0

	for _, row := range positions {

		// Mark if this trade group is open or closed.
		if row.Qty != 0 {
			tradeGroupStatus = "Open"
		}

		if row.TradeGroupId > 0 {
			tradeGroupId = row.TradeGroupId
		}

	}

	// TODO: Figure out Risked, Commission, Gain, and Profit

	// Create or Update Trade Group
	if tradeGroupId == 0 {

		// Build a new Trade Group
		var tradeGroup = &models.TradeGroup{
			UserId:     userId,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
			AccountId:  order.AccountId,
			Status:     tradeGroupStatus,
			OrderIds:   strconv.Itoa(int(order.Id)),
			Note:       "",
			OpenDate:   order.CreateDate,
			ClosedDate: order.TransactionDate,
		}

		// Insert into DB
		db.CreateTradeGroup(tradeGroup)

		// Store tradegroup id
		tradeGroupId = tradeGroup.Id

	} else {

		// Update tradegroup with additional OrderIds
		tradeGroup, err := db.GetTradeGroupById(tradeGroupId)

		if err != nil {
			return err
		}

		tradeGroup.Status = tradeGroupStatus
		tradeGroup.ClosedDate = order.TransactionDate
		tradeGroup.OrderIds = tradeGroup.OrderIds + "," + strconv.Itoa(int(order.Id))
		db.UpdateTradeGroup(&tradeGroup)

	}

	// Loop through the positions and add the trade group id
	for _, row := range positions {
		row.TradeGroupId = tradeGroupId
		db.UpdatePosition(row)
	}

	// Log success
	services.Info("New TradeGroup created for user " + strconv.Itoa(int(userId)) + " TradeGroup Id: " + strconv.Itoa(int(tradeGroupId)))

	// Return happy.
	return nil
}

//
// Do one leg of a multi leg order - Open Order
//
func doOpenOneLegMultiLegOrder(order models.Order, leg models.OrderLeg, db models.Datastore, userId uint) (*models.Position, error) {

	// First we find out if we already have a position on for this.
	position, _ := db.GetPositionByUserSymbolStatusAccount(userId, leg.OptionSymbol, "Open", order.AccountId)

	// We found so we are just adding to a current position.
	if position.Id > 0 {

		// Update pos
		position.OrderIds = position.OrderIds + "," + strconv.Itoa(int(order.Id))
		position.UpdatedAt = time.Now()
		position.Qty = leg.Qty + position.Qty
		position.OrgQty = leg.Qty + position.OrgQty
		position.AvgOpenPrice = ((leg.AvgFillPrice + position.AvgOpenPrice) / 2)
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
			Symbol:        leg.OptionSymbol,
			Qty:           leg.Qty,
			OrgQty:        leg.Qty,
			CostBasis:     (float64(leg.Qty) * leg.AvgFillPrice * 100),
			AvgOpenPrice:  leg.AvgFillPrice,
			AvgClosePrice: 0.00,
			Note:          "",
			OpenDate:      leg.CreateDate,
			ClosedDate:    leg.TransactionDate,
			OrderIds:      strconv.Itoa(int(order.Id)),
			Status:        "Open",
		}

		// Insert into DB
		db.CreatePosition(&position)

	}

	// Return a list of position that we reviewed
	return &position, nil
}

// //
// // Do one leg of a multi leg order - Close Order
// //
// func doCloseOneLegMultiLegOrder(order models.Order, leg models.OrderLeg, db models.Datastore, userId uint) (*models.Position, error) {

// 	var position = &models.Position{}

// 	// First we find out if we already have a position on for this.
// 	db.Where("symbol = ? AND user_id = ? AND status = ? AND account_id = ?", leg.OptionSymbol, userId, "Open", order.AccountId).First(position)

// 	// We found so we are just removing to a current position.
// 	if position.Id > 0 {

// 		// Update pos
// 		position.OrderIds = position.OrderIds + "," + strconv.Itoa(int(order.Id))
// 		position.UpdatedAt = time.Now()
// 		position.Qty = leg.Qty + position.Qty

// 		if position.AvgClosePrice != 0 {
// 			position.AvgClosePrice = ((leg.AvgFillPrice + position.AvgClosePrice) / 2)
// 		} else {
// 			position.AvgClosePrice = leg.AvgFillPrice
// 		}

// 		// Are we done with the trade?
// 		if position.Qty == 0 {

// 			position.ClosedDate = leg.TransactionDate
// 			position.Status = "Closed"

// 		} else {

// 			position.Note = position.Note + "Updated - " + leg.TransactionDate.Format(time.RFC1123) + " :: "

// 		}

// 		db.Save(&position)

// 	} else {

// 		return nil, errors.New("Unable to find close position in our database.")

// 	}

// 	// Return a list of position that we reviewed
// 	return position, nil

// }

/* End File */
