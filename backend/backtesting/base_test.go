//
// Date: 2019-02-22
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package backtesting

import (
	"os"
	"testing"

	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/nbio/st"
	"github.com/optionscafe/options-cafe-cli/helpers"
)

//
// TestDoBacktestDays01 - Run a backtest looping through each day
//
func TestDoBacktestDays01(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()
	defer db.Close()

	// Setup a new backtesting
	bt := New(db)

	// Run blank backtest
	err := bt.DoBacktestDays(Action{
		Type:      "blank",
		Symbol:    "SPY",
		StartDate: helpers.ParseDateNoError("2018-01-01"),
		EndDate:   helpers.ParseDateNoError("2018-01-03"),
	})
	st.Expect(t, err, nil)

	// Verify our cache files got set.
	cacheDir := os.Getenv("CACHE_DIR") + "/" + cacheDirBase
	cacheFile := cacheDir + "/SPY-2018-01-02.json"
	_, err2 := os.Stat(cacheFile)
	st.Expect(t, err2, nil)
}

/* End File */
