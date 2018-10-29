//
// Date: 2018-07-17
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package screener

import (
	"os"
	"testing"

	"github.com/cloudmanic/app.options.cafe/backend/brokers/tradier"
	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/davecgh/go-spew/spew"
	"github.com/nbio/st"
)

//
// RunPutCreditSpread - 01
//
// Note: This is too hard to mock as it makes many requests to Tradier.
// This unit test is just for development.
//
func TestRunPutCreditSpread01(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()
	defer db.Close()

	// Build screener object
	screen := models.Screener{
		Symbol:   "SPY",
		Strategy: "put-credit-spread",
		Items: []models.ScreenerItem{
			{Key: "short-strike-percent-away", Operator: "<", ValueNumber: 4.0},
			{Key: "spread-width", Operator: "=", ValueNumber: 2.00},
			{Key: "open-credit", Operator: ">", ValueNumber: 0.18},
			{Key: "open-credit", Operator: "<", ValueNumber: 0.20},
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
	result, err := s.RunPutCreditSpread(screen)

	spew.Dump(result)

	// Test result
	st.Expect(t, err, nil)

}

/* End File */
