//
// Date: 2018-10-27
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-10-27
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package screener

import (
	"testing"

	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/davecgh/go-spew/spew"
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
		Symbol:   "SPY",
		Strategy: "reverse-iron-condor",
		Items: []models.ScreenerItem{
			{Key: "put-leg-width", Operator: "=", ValueNumber: 2.00},
			{Key: "call-leg-width", Operator: "=", ValueNumber: 2.00},
			{Key: "put-leg-percent-away", Operator: ">", ValueNumber: 4.0},
			{Key: "call-leg-percent-away", Operator: "<", ValueNumber: 4.0},
			{Key: "open-credit", Operator: ">", ValueNumber: 0.18},
			{Key: "open-credit", Operator: "<", ValueNumber: 6.00},
			{Key: "days-to-expire", Operator: "<", ValueNumber: 46},
			{Key: "days-to-expire", Operator: ">", ValueNumber: 0},
		},
	}

	// Run back test
	result, err := RunReverseIronCondor(screen, db)

	spew.Dump(result)

	// Test result
	st.Expect(t, err, nil)

}

/* End File */
