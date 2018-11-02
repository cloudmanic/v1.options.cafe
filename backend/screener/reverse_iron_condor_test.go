//
// Date: 2018-10-27
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-11-01
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package screener

import (
	"fmt"
	"os"
	"testing"

	"github.com/cloudmanic/app.options.cafe/backend/brokers/tradier"
	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/nbio/st"
)

//
// RunReverseIronCondor - 01
//
// Note: This is too hard to mock as it makes many requests to Tradier.
// This unit test is just for development.
//
func TestRunReverseIronCondor01(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()
	defer db.Close()

	// Build screener object
	screen := models.Screener{
		Symbol:   "VXX",
		Strategy: "reverse-iron-condor",
		Items: []models.ScreenerItem{
			{Key: "put-leg-width", Operator: "=", ValueNumber: 2.00},
			{Key: "call-leg-width", Operator: "=", ValueNumber: 2.00},
			{Key: "put-leg-percent-away", Operator: ">", ValueNumber: 4.0},
			{Key: "call-leg-percent-away", Operator: ">", ValueNumber: 4.0},
			{Key: "open-debit", Operator: ">", ValueNumber: 0.50},
			{Key: "open-debit", Operator: "<", ValueNumber: 3.00},
			{Key: "days-to-expire", Operator: "<", ValueNumber: 46},
			{Key: "days-to-expire", Operator: ">", ValueNumber: 0},
		},
	}

	// Setup the broker
	broker := tradier.Api{
		DB:     nil,
		ApiKey: os.Getenv("TRADIER_ADMIN_ACCESS_TOKEN"),
	}

	// New screener instance
	s := NewScreen(db, &broker)

	// Run back test
	result, err := s.RunReverseIronCondor(screen)

	for _, row := range result {
		fmt.Println(row.Debit)
	}

	// Test result
	st.Expect(t, err, nil)

}

/* End File */
