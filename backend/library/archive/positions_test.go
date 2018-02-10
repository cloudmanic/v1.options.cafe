//
// Date: 2/9/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package archive

import (
	"testing"

	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/nbio/st"
)

//
// Test Commission for order.
//
func TestCalcCommissionForOrder01(t *testing.T) {

	// Build sample order.
	order := models.Order{
		Class: "multileg",
		Legs: []models.OrderLeg{
			{Qty: 10},
			{Qty: 10},
		},
	}

	// Build sample BrokerAccount
	brokerAccount := models.BrokerAccount{
		StockCommission:   5.00,
		StockMin:          0.00,
		OptionCommission:  0.35,
		OptionSingleMin:   5.00,
		OptionMultiLegMin: 7.00,
		OptionBase:        0.00,
	}

	// Figure out Commission
	commission := calcCommissionForOrder(&order, 1, &brokerAccount)

	// Verify the data was return as expected
	st.Expect(t, commission, 7.00)
}

//
// Test Commission for order.
//
func TestCalcCommissionForOrder02(t *testing.T) {

	// Build sample order.
	order := models.Order{
		Class: "multileg",
		Legs: []models.OrderLeg{
			{Qty: 15},
			{Qty: 10},
		},
	}

	// Build sample BrokerAccount
	brokerAccount := models.BrokerAccount{
		StockCommission:   5.00,
		StockMin:          0.00,
		OptionCommission:  0.35,
		OptionSingleMin:   5.00,
		OptionMultiLegMin: 7.00,
		OptionBase:        0.00,
	}

	// Figure out Commission
	commission := calcCommissionForOrder(&order, 1, &brokerAccount)

	// Verify the data was return as expected
	st.Expect(t, commission, 8.75)
}

/* End File */
