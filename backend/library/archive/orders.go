package archive

import (
	"fmt"
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/brokers/types"
	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/stvp/rollbar"
)

//
// Pass in all orders and archive them by putting them into our database.
// We only archive orders that are filled.
//
// TODO: Review partially_filled orders.
//
func StoreOrders(db models.Datastore, orders []types.Order, userId uint) error {

	// Loop through the orders and process
	for _, row := range orders {

		// We only care about filled orders
		if row.Status != "filled" {
			continue
		}

		// See if we already have this record in our database
		if db.HasOrderByBrokerIdUserId(uint(row.Id), userId) {
			continue
		}

		// Timestamp Layout
		layout := "2006-01-02T15:04:05.000Z"

		// Convert Create Date
		createDate, err := time.Parse(layout, row.CreateDate)

		if err != nil {
			fmt.Println(err)
			rollbar.Error(rollbar.ERR, err)
			continue
		}

		// Convert TransactionDate
		transactionDate, err := time.Parse(layout, row.TransactionDate)

		if err != nil {
			services.Fatal(err)
			continue
		}

		// Insert into DB
		order := &models.Order{
			UserId:            userId,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			BrokerId:          row.Id,
			AccountId:         row.AccountId,
			Type:              row.Type,
			Symbol:            row.Symbol,
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
			services.Fatal(err)
			continue
		}

		// Build the legs
		legs := []models.OrderLeg{}

		for _, row2 := range row.Legs {
			legs = append(legs, models.OrderLeg{
				UserId:            userId,
				OrderId:           order.Id,
				Type:              row2.Type,
				Symbol:            row2.Symbol,
				OptionSymbol:      row2.OptionSymbol,
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

	// Return Happy
	return nil

}

/* End File */
