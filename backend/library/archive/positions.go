//
// Date: 2/9/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package archive

import (
	"math"

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

	// Return happy
	return nil
}

//
// Review the order and calculate the commission for this order.
//
func calcCommissionForOrder(order *models.Order, brokerId uint, brokerAccount *models.BrokerAccount) float64 {

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
