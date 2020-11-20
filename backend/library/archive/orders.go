package archive

import (
	"time"

	"app.options.cafe/brokers/types"
	"app.options.cafe/library/services"
	"app.options.cafe/models"
)

//
// Pass in all orders and archive them by putting them into our database.
// We only archive orders that are filled.
//
// TODO: Review partially_filled orders.
//
func StoreOrders(db models.Datastore, orders []types.Order, userId uint, brokerId uint) error {

	// Loop through the orders and process
	for _, row := range orders {

		// This is how we handle parcel fills. A little hack.
		if (row.Status == "canceled") && (row.ExecQuantity > 0) {
			row.Status = "filled"
		}

		// We only care about filled orders
		if row.Status != "filled" {
			continue
		}

		// See if we already have this record in our database
		if db.HasOrderByBrokerRefUserId(row.Id, userId) {
			continue
		}

		// Get broker account id
		brokerAccount, err := db.GetBrokerAccountByBrokerAccountNumber(brokerId, row.AccountId)

		// Timestamp Layout
		layout := "2006-01-02T15:04:05.000Z"

		// Convert Create Date
		createDate, err := time.Parse(layout, row.CreateDate)

		if err != nil {
			services.Fatal(err)
			continue
		}

		// Convert TransactionDate
		transactionDate, err := time.Parse(layout, row.TransactionDate)

		if err != nil {
			services.Fatal(err)
			continue
		}

		//
		// Because Tradier data is bad we only import from 2018 forward.
		//
		if transactionDate.Year() < 2017 {
			continue
		}

		//
		// Because Tradier data is bad we only import from 2018 forward.
		//
		if createDate.Year() < 2017 {
			continue
		}

		// Get the symbol id.
		// TODO: notice the symbol name is blank. We should always have this symbol in the DB. But maybe we should prepare for it not being there just in case
		sym, err := db.CreateNewSymbol(row.Symbol, "", "Equity")

		if err != nil {
			services.Fatal(err)
			continue
		}

		// Do we have an option symbol?
		var oSym uint = 0

		if len(row.OptionSymbol) > 0 {
			sO, err := db.CreateNewOptionSymbol(row.OptionSymbol)

			if err != nil {
				services.Info(err)
				continue
			}

			oSym = sO.Id
		}

		// Insert into DB
		order := &models.Order{
			UserId:            userId,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			BrokerAccountId:   brokerAccount.Id,
			BrokerRef:         row.Id,
			BrokerAccountRef:  row.AccountId,
			Type:              row.Type,
			SymbolId:          sym.Id,
			OptionSymbolId:    oSym,
			Side:              row.Side,
			Qty:               int(row.Quantity),
			Status:            row.Status,
			Duration:          row.Duration,
			Price:             row.Price,
			AvgFillPrice:      row.AvgFillPrice,
			ExecQuantity:      row.ExecQuantity,
			LastFillPrice:     row.LastFillPrice,
			LastFillQuantity:  row.LastFillQuantity,
			RemainingQuantity: row.RemainingQuantity,
			CreateDate:        createDate,
			TransactionDate:   transactionDate,
			Class:             row.Class,
			PositionReviewed:  "No",
			NumLegs:           row.NumLegs,
		}

		err = db.CreateOrder(order)

		if err != nil {
			services.Info(err)
			continue
		}

		// Build the legs
		legs := []models.OrderLeg{}

		for _, row2 := range row.Legs {

			// Get the symbol id.
			sym, err := db.CreateNewOptionSymbol(row2.OptionSymbol)

			if err != nil {
				services.Info(err)
				continue
			}

			legs = append(legs, models.OrderLeg{
				UserId:            userId,
				OrderId:           order.Id,
				Type:              row2.Type,
				SymbolId:          sym.Id,
				Side:              row2.Side,
				Qty:               int(row2.Quantity),
				Status:            row2.Status,
				Duration:          row2.Duration,
				AvgFillPrice:      row2.AvgFillPrice,
				ExecQuantity:      row2.ExecQuantity,
				LastFillPrice:     row2.LastFillPrice,
				LastFillQuantity:  row2.LastFillQuantity,
				RemainingQuantity: row2.RemainingQuantity,
				CreateDate:        createDate,
				TransactionDate:   transactionDate,
			})
		}

		// Add the legs to the order and update.
		order.Legs = legs

		err = db.UpdateOrder(order)

		if err != nil {
			services.Fatal(err)
			continue
		}
	}

	// Now build out our positions database table based on past orders.
	StorePositions(db, userId, brokerId)

	// Return Happy
	return nil

}

/* End File */
