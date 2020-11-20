//
// Date: 2018-12-23
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-12-23
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package market

import (
	"testing"
	"time"

	"app.options.cafe/library/helpers"
	"app.options.cafe/models"
	"github.com/nbio/st"
	gock "gopkg.in/h2non/gock.v1"
)

//
// Test - Return a underlying stock quote based on date in time.
//
func TestGetUnderlayingQuoteByDate01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Shared vars we use.
	ts := time.Date(2017, 10, 29, 17, 20, 01, 507451, time.UTC)

	// Users
	db.Create(&models.User{FirstName: "Rob", LastName: "Tester", Email: "spicer+robtester@options.cafe", Status: "Active"})
	db.Create(&models.User{FirstName: "Jane", LastName: "Wells", Email: "spicer+janewells@options.cafe", Status: "Active"})
	db.Create(&models.User{FirstName: "Bob", LastName: "Rosso", Email: "spicer+bobrosso@options.cafe", Status: "Active"})

	// Brokers
	db.Create(&models.Broker{Name: "Tradier", UserId: 1, AccessToken: "CLOwLO2cMnx-N_bPEexiVo9z9oRR80nPI9ycxQw3KQ-WQ4OP3D44gIbfLScAZ9pv", RefreshToken: "abc", TokenExpirationDate: ts})
	db.Create(&models.Broker{Name: "Tradeking", UserId: 1, AccessToken: "456", RefreshToken: "xyz", TokenExpirationDate: ts})
	db.Create(&models.Broker{Name: "Etrade", UserId: 1, AccessToken: "789", RefreshToken: "mno", TokenExpirationDate: ts})

	// BrokerAccounts
	db.Create(&models.BrokerAccount{UserId: 1, BrokerId: 1, Name: "Test Account 1", AccountNumber: "abc1235423", StockCommission: 5.00, StockMin: 0.00, OptionCommission: 0.35, OptionSingleMin: 5.00, OptionMultiLegMin: 7.00, OptionBase: 0.00})
	db.Create(&models.BrokerAccount{UserId: 1, BrokerId: 1, Name: "Test Account 2", AccountNumber: "ABC123ZY", StockCommission: 5.00, StockMin: 0.00, OptionCommission: 0.35, OptionSingleMin: 5.00, OptionMultiLegMin: 7.00, OptionBase: 0.00})

	// Symbols
	db.Create(&models.Symbol{Name: "SPDR S&P 500 ETF Trust", ShortName: "SPY", Type: "Equity"})
	db.Create(&models.Symbol{Name: "McDonald's Corp", ShortName: "MCD", Type: "Equity"})
	db.Create(&models.Symbol{Name: "Starbucks Corp", ShortName: "SBUX", Type: "Equity"})
	db.Create(&models.Symbol{Name: "SPY Mar 16, 2018 $253.00 Put", ShortName: "SPY180316P00253000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 253.00})
	db.Create(&models.Symbol{Name: "VXX Feb 23, 2018 $55.00 Call", ShortName: "VXX180223C00055000", Type: "Option", OptionUnderlying: "VXX", OptionExpire: models.Date{helpers.ParseDateNoError("2/23/2018").UTC()}, OptionType: "Call", OptionStrike: 55.00})
	db.Create(&models.Symbol{Name: "SPY Mar 16, 2018 $266.00 Put", ShortName: "SPY180316P00266000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 266.00})
	db.Create(&models.Symbol{Name: "Caterpillar Inc", ShortName: "CAT", Type: "Equity"})
	db.Create(&models.Symbol{Name: "Ascent Solar Technologies Inc", ShortName: "ASTI", Type: "Equity"})
	db.Create(&models.Symbol{Name: "Advanced Micro Devices Inc", ShortName: "AMD", Type: "Equity"})
	db.Create(&models.Symbol{Name: "BARCLAYS BK PLC", ShortName: "VXX", Type: "Equity"})
	db.Create(&models.Symbol{Name: "SPY Feb 9, 2018 $276.00 Put", ShortName: "SPY180209P00276000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("2/9/2018").UTC()}, OptionType: "Put", OptionStrike: 276.00})
	db.Create(&models.Symbol{Name: "VXX Mar 2, 2018 $46.00 Put", ShortName: "VXX180302P00046000", Type: "Option", OptionUnderlying: "VXX", OptionExpire: models.Date{helpers.ParseDateNoError("3/2/2018").UTC()}, OptionType: "Put", OptionStrike: 46.00})
	db.Create(&models.Symbol{Name: "VXX Feb 23, 2018 $50.00 Put", ShortName: "VXX180223P00050000", Type: "Option", OptionUnderlying: "VXX", OptionExpire: models.Date{helpers.ParseDateNoError("2/23/2018").UTC()}, OptionType: "Put", OptionStrike: 50.00})
	db.Create(&models.Symbol{Name: "SPY Mar 16, 2018 $267.00 Put", ShortName: "SPY180316P00267000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 267.00})
	db.Create(&models.Symbol{Name: "SPY Mar 16, 2018 $269.00 Put", ShortName: "SPY180316P00269000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 269.00})
	db.Create(&models.Symbol{Name: "SPY Mar 16, 2018 $268.00 Put", ShortName: "SPY180316P00268000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 268.00})
	db.Create(&models.Symbol{Name: "SPY Mar 16, 2018 $270.00 Put", ShortName: "SPY180316P00270000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 270.00})
	db.Create(&models.Symbol{Name: "SPY Mar 16, 2018 $264.00 Put", ShortName: "SPY180316P00264000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 264.00})
	db.Create(&models.Symbol{Name: "SPY Mar 16, 2018 $266.00 Put", ShortName: "SPY180316P00266000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 266.00})
	db.Create(&models.Symbol{Name: "SPY Mar 9, 2018 $262.00 Put", ShortName: "SPY180309P00262000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("3/9/2018").UTC()}, OptionType: "Put", OptionStrike: 262.00})
	db.Create(&models.Symbol{Name: "SPY Mar 9, 2018 $264.00 Put", ShortName: "SPY180309P00264000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("3/9/2018").UTC()}, OptionType: "Put", OptionStrike: 264.00})
	db.Create(&models.Symbol{Name: "SPY Mar 16, 2018 $251.00 Put", ShortName: "SPY180316P00251000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 251.00})
	db.Create(&models.Symbol{Name: "SPY Mar 16, 2018 $253.00 Put", ShortName: "SPY180316P00253000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 253.00})
	db.Create(&models.Symbol{Name: "SPY Mar 16, 2018 $250.00 Put", ShortName: "SPY180316P00250000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 250.00})
	db.Create(&models.Symbol{Name: "SPY Mar 16, 2018 $252.00 Put", ShortName: "SPY180316P00252000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("3/16/2018").UTC()}, OptionType: "Put", OptionStrike: 252.00})
	db.Create(&models.Symbol{Name: "VXX Feb 23, 2018 $50.00 Call", ShortName: "VXX180223C00050000", Type: "Option", OptionUnderlying: "VXX", OptionExpire: models.Date{helpers.ParseDateNoError("2/23/2018").UTC()}, OptionType: "Call", OptionStrike: 50.00})
	db.Create(&models.Symbol{Name: "SPY Feb 28, 2018 $249.00 Put", ShortName: "SPY180228P00249000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("2/28/2018").UTC()}, OptionType: "Put", OptionStrike: 249.00})
	db.Create(&models.Symbol{Name: "SPY Feb 28, 2018 $251.00 Put", ShortName: "SPY180228P00251000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("2/28/2018").UTC()}, OptionType: "Put", OptionStrike: 251.00})
	db.Create(&models.Symbol{Name: "SPY Feb 21, 2018 $241.00 Put", ShortName: "SPY180221P00241000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("2/21/2018").UTC()}, OptionType: "Put", OptionStrike: 241.00})
	db.Create(&models.Symbol{Name: "SPY Feb 21, 2018 $243.00 Put", ShortName: "SPY180221P00243000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("2/21/2018").UTC()}, OptionType: "Put", OptionStrike: 243.00})
	db.Create(&models.Symbol{Name: "SPY Mar 21, 2018 $254.00 Put", ShortName: "SPY180321P00254000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("3/21/2018").UTC()}, OptionType: "Put", OptionStrike: 254.00})
	db.Create(&models.Symbol{Name: "SPY Mar 21, 2018 $256.00 Put", ShortName: "SPY180321P00256000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("3/21/2018").UTC()}, OptionType: "Put", OptionStrike: 256.00})
	db.Create(&models.Symbol{Name: "SPY Jul 21, 2017 $234.00 Put", ShortName: "SPY170721P00234000", Type: "Option", OptionUnderlying: "SPY", OptionExpire: models.Date{helpers.ParseDateNoError("7/21/2018").UTC()}, OptionType: "Put", OptionStrike: 234.00})

	// Clear tables - Default test values throws API Token to Tradier off.
	db.Exec("TRUNCATE TABLE brokers;")
	db.Exec("TRUNCATE TABLE broker_accounts;")

	// Flush pending mocks after test execution
	defer gock.Off()

	// Setup mock request.
	gock.New("https://api.tradier.com/v1").
		Get("/markets/history").
		Reply(200).
		BodyString(`{"history":{"day":{"date":"2018-12-21","open":246.74,"high":249.71,"low":239.98,"close":240.7,"volume":255320360}}}`)

	// Make API call
	quote, err := GetUnderlayingQuoteByDate(db, uint(1), "spy", helpers.ParseDateNoError("12/21/18"))

	if err != nil {
		panic(err)
	}

	// Verify the data was return as expected
	st.Expect(t, quote, 240.70)

}

/* End File */
