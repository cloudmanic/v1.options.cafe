//
// Date: 2018-07-17
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package screener

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
)

//
// RunPutCreditSpread - 01
//
func TestRunPutCreditSpread01(t *testing.T) {

	filter := Filter{
		Symbol:   "SPY",
		Strategy: "put-credit-spread",
		FilterItems: []FilterItems{
			{Key: "short-strike-percent-away", Operator: "=", ValueNumber: 2.5},
			{Key: "spread-width", Operator: "=", ValueNumber: 2.00},
			{Key: "min-credit", Operator: "=", ValueNumber: 0.18},
			{Key: "max-days-to-expire", Operator: "=", ValueNumber: 45},
			{Key: "min-days-to-expire", Operator: "=", ValueNumber: 0},
		},
	}

	result, _ := RunPutCreditSpread(filter)

	spew.Dump(result)

}

/* End File */
