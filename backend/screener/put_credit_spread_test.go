//
// Date: 2018-07-17
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package screener

import (
	"testing"

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

	screen := models.Screener{
		Symbol:   "SPY",
		Strategy: "put-credit-spread",
		Items: []models.ScreenerItem{
			{Key: "short-strike-percent-away", Operator: "=", ValueNumber: 2.5},
			{Key: "spread-width", Operator: "=", ValueNumber: 2.00},
			{Key: "min-credit", Operator: "=", ValueNumber: 0.18},
			{Key: "max-days-to-expire", Operator: "=", ValueNumber: 45},
			{Key: "min-days-to-expire", Operator: "=", ValueNumber: 0},
		},
	}

	// Run back test
	result, err := RunPutCreditSpread(screen)

	spew.Dump(result)

	// Test result
	st.Expect(t, err, nil)

}

/* End File */
