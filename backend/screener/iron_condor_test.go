//
// Date: 2018-10-27
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-11-01
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package screener

import (
	"testing"

	"github.com/cloudmanic/app.options.cafe/backend/brokers/eod"
	"github.com/cloudmanic/app.options.cafe/backend/library/helpers"
	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/davecgh/go-spew/spew"
	"github.com/nbio/st"
)

//
// RunIronCondor - 01
//
func TestRunIronCondor01(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()
	defer db.Close()

	// Build screener object
	screen := models.Screener{
		Symbol:   "SPY",
		Strategy: "iron-condor",
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

	// // Setup the broker - Tradier Test
	// broker := tradier.Api{
	// 	DB:     nil,
	// 	ApiKey: os.Getenv("TRADIER_ADMIN_ACCESS_TOKEN"),
	// }

	// Setup the broker - EOD Test
	broker := eod.Api{
		DB:  db,
		Day: helpers.ParseDateNoError("2018-10-18").UTC(),
	}

	// New screener instance
	s := NewScreen(db, &broker)

	// Run back test
	result, err := s.RunIronCondor(screen)

	spew.Dump(result)

	// for _, row := range result {
	// 	fmt.Println(row.Debit)
	// }

	// Test result
	st.Expect(t, err, nil)

}

/* End File */
