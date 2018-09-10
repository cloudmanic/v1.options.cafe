//
// Date: 9/10/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package reports

import (
	"testing"
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/nbio/st"
)

//
// Load and return some test data.
//
func LoadAndGetTestData() models.BrokerAccount {

	// Start the db connection.
	db, _ := models.NewDB()
	defer db.Close()

	// Put test data into database
	brokerAccount := models.BrokerAccount{
		UserId:            1,
		BrokerId:          1,
		Name:              "Unit Test Account #1",
		AccountNumber:     "test12345",
		StockCommission:   5.00,
		StockMin:          0.00,
		OptionCommission:  0.35,
		OptionSingleMin:   5.00,
		OptionMultiLegMin: 7.00,
		OptionBase:        0.00,
	}

	// Insert into database.
	db.Create(&brokerAccount)

	// Load test Trade Group data
	ts := time.Date(2017, 10, 29, 17, 20, 01, 507451, time.UTC)
	db.Create(&models.TradeGroup{Risked: 1980.00, Credit: 220.00, Proceeds: -33.00, Profit: 171.60, Commission: 15.40, UserId: 1, PercentGain: 8.67, BrokerAccountId: 3, BrokerAccountRef: "test12345", Name: "Trade #1 - Put Credit Spread Trade", Status: "Closed", Type: "Put Credit Spread", OrderIds: "1", Note: "Test note", OpenDate: ts, ClosedDate: time.Date(2017, 01, 20, 17, 20, 01, 507451, time.UTC)})
	db.Create(&models.TradeGroup{Risked: 1980.00, Credit: 220.00, Proceeds: -33.00, Profit: 204.60, Commission: 15.40, UserId: 1, PercentGain: 8.67, BrokerAccountId: 3, BrokerAccountRef: "test12345", Name: "Trade #2 - Put Credit Spread Trade", Status: "Closed", Type: "Put Credit Spread", OrderIds: "1", Note: "Test note", OpenDate: ts, ClosedDate: time.Date(2017, 01, 29, 17, 20, 01, 507451, time.UTC)})
	db.Create(&models.TradeGroup{Risked: 1980.00, Credit: 220.00, Proceeds: -33.00, Profit: 171.60, Commission: 15.40, UserId: 1, PercentGain: 8.67, BrokerAccountId: 3, BrokerAccountRef: "test12345", Name: "Trade #3 - Put Credit Spread Trade", Status: "Closed", Type: "Put Credit Spread", OrderIds: "1", Note: "Test note", OpenDate: ts, ClosedDate: time.Date(2017, 02, 17, 17, 20, 01, 507451, time.UTC)})
	db.Create(&models.TradeGroup{Risked: 1980.00, Credit: 220.00, Proceeds: -33.00, Profit: 171.60, Commission: 15.40, UserId: 1, PercentGain: 8.67, BrokerAccountId: 3, BrokerAccountRef: "test12345", Name: "Trade #4 - Put Credit Spread Trade", Status: "Closed", Type: "Put Credit Spread", OrderIds: "1", Note: "Test note", OpenDate: ts, ClosedDate: time.Date(2017, 02, 13, 17, 20, 01, 507451, time.UTC)})
	db.Create(&models.TradeGroup{Risked: 1980.00, Credit: 220.00, Proceeds: -33.00, Profit: 171.60, Commission: 15.40, UserId: 1, PercentGain: 8.67, BrokerAccountId: 3, BrokerAccountRef: "test12345", Name: "Trade #5 - Put Credit Spread Trade", Status: "Closed", Type: "Put Credit Spread", OrderIds: "1", Note: "Test note", OpenDate: ts, ClosedDate: time.Date(2017, 03, 11, 17, 20, 01, 507451, time.UTC)})
	db.Create(&models.TradeGroup{Risked: 1980.00, Credit: 220.00, Proceeds: -33.00, Profit: 171.60, Commission: 15.40, UserId: 1, PercentGain: 8.67, BrokerAccountId: 3, BrokerAccountRef: "test12345", Name: "Trade #6 - Put Credit Spread Trade", Status: "Closed", Type: "Put Credit Spread", OrderIds: "1", Note: "Test note", OpenDate: ts, ClosedDate: time.Date(2017, 04, 04, 17, 20, 01, 507451, time.UTC)})
	db.Create(&models.TradeGroup{Risked: 1980.00, Credit: 220.00, Proceeds: -33.00, Profit: 171.60, Commission: 15.40, UserId: 1, PercentGain: 8.67, BrokerAccountId: 3, BrokerAccountRef: "test12345", Name: "Trade #7 - Put Credit Spread Trade", Status: "Closed", Type: "Put Credit Spread", OrderIds: "1", Note: "Test note", OpenDate: ts, ClosedDate: time.Date(2017, 04, 21, 17, 20, 01, 507451, time.UTC)})
	db.Create(&models.TradeGroup{Risked: 1947.00, Credit: 253.00, Proceeds: -33.00, Profit: 204.60, Commission: 15.40, UserId: 1, PercentGain: 10.51, BrokerAccountId: 3, BrokerAccountRef: "test12345", Name: "Trade #8 - Put Credit Spread Trade", Status: "Closed", Type: "Put Credit Spread", OrderIds: "1", Note: "Test note", OpenDate: ts, ClosedDate: time.Date(2017, 04, 15, 17, 20, 01, 507451, time.UTC)})
	db.Create(&models.TradeGroup{Risked: 1947.00, Credit: 253.00, Proceeds: -33.00, Profit: 204.60, Commission: 15.40, UserId: 1, PercentGain: 10.51, BrokerAccountId: 3, BrokerAccountRef: "test12345", Name: "Trade #9 - Put Credit Spread Trade", Status: "Closed", Type: "Put Credit Spread", OrderIds: "1", Note: "Test note", OpenDate: ts, ClosedDate: time.Date(2017, 05, 13, 17, 20, 01, 507451, time.UTC)})
	db.Create(&models.TradeGroup{Risked: 1947.00, Credit: 253.00, Proceeds: -33.00, Profit: 204.60, Commission: 15.40, UserId: 1, PercentGain: 10.51, BrokerAccountId: 3, BrokerAccountRef: "test12345", Name: "Trade #10 - Put Credit Spread Trade", Status: "Closed", Type: "Put Credit Spread", OrderIds: "1", Note: "Test note", OpenDate: ts, ClosedDate: time.Date(2017, 05, 14, 17, 20, 01, 507451, time.UTC)})
	db.Create(&models.TradeGroup{Risked: 1947.00, Credit: 253.00, Proceeds: -33.00, Profit: 204.60, Commission: 15.40, UserId: 1, PercentGain: 10.51, BrokerAccountId: 3, BrokerAccountRef: "test12345", Name: "Trade #11 - Put Credit Spread Trade", Status: "Closed", Type: "Put Credit Spread", OrderIds: "1", Note: "Test note", OpenDate: ts, ClosedDate: time.Date(2017, 06, 15, 17, 20, 01, 507451, time.UTC)})
	db.Create(&models.TradeGroup{Risked: 1947.00, Credit: 253.00, Proceeds: -33.00, Profit: 204.60, Commission: 15.40, UserId: 1, PercentGain: 10.51, BrokerAccountId: 3, BrokerAccountRef: "test12345", Name: "Trade #12 - Put Credit Spread Trade", Status: "Closed", Type: "Put Credit Spread", OrderIds: "1", Note: "Test note", OpenDate: ts, ClosedDate: time.Date(2017, 06, 16, 17, 20, 01, 507451, time.UTC)})
	db.Create(&models.TradeGroup{Risked: 1947.00, Credit: 253.00, Proceeds: -33.00, Profit: 204.60, Commission: 15.40, UserId: 1, PercentGain: 10.51, BrokerAccountId: 3, BrokerAccountRef: "test12345", Name: "Trade #13 - Put Credit Spread Trade", Status: "Closed", Type: "Put Credit Spread", OrderIds: "1", Note: "Test note", OpenDate: ts, ClosedDate: time.Date(2017, 07, 17, 17, 20, 01, 507451, time.UTC)})
	db.Create(&models.TradeGroup{Risked: 1947.00, Credit: 253.00, Proceeds: -33.00, Profit: 204.60, Commission: 15.40, UserId: 1, PercentGain: 10.51, BrokerAccountId: 3, BrokerAccountRef: "test12345", Name: "Trade #14 - Put Credit Spread Trade", Status: "Closed", Type: "Put Credit Spread", OrderIds: "1", Note: "Test note", OpenDate: ts, ClosedDate: time.Date(2017, 07, 18, 17, 20, 01, 507451, time.UTC)})
	db.Create(&models.TradeGroup{Risked: 1770.00, Credit: 230.00, Proceeds: -400.00, Profit: -184.000, Commission: 14.00, UserId: 1, PercentGain: -10.40, BrokerAccountId: 3, BrokerAccountRef: "test12345", Name: "Trade #15 - Put Credit Spread Trade", Status: "Closed", Type: "Put Credit Spread", OrderIds: "1", Note: "Test note", OpenDate: ts, ClosedDate: time.Date(2017, 02, 04, 17, 20, 01, 507451, time.UTC)})
	db.Create(&models.TradeGroup{Risked: 1770.00, Credit: 230.00, Proceeds: -400.00, Profit: -184.00, Commission: 14.00, UserId: 1, PercentGain: -10.40, BrokerAccountId: 3, BrokerAccountRef: "test12345", Name: "Trade #16 - Put Credit Spread Trade", Status: "Closed", Type: "Put Credit Spread", OrderIds: "1", Note: "Test note", OpenDate: ts, ClosedDate: time.Date(2017, 03, 21, 17, 20, 01, 507451, time.UTC)})
	db.Create(&models.TradeGroup{Risked: 3195.00, Credit: 405.00, Proceeds: -1242.00, Profit: -863.60, Commission: 26.60, UserId: 1, PercentGain: -27.03, BrokerAccountId: 3, BrokerAccountRef: "test12345", Name: "Trade #17 - Put Credit Spread Trade", Status: "Closed", Type: "Put Credit Spread", OrderIds: "1", Note: "Test note", OpenDate: ts, ClosedDate: time.Date(2017, 04, 15, 17, 20, 01, 507451, time.UTC)})
	db.Create(&models.TradeGroup{Risked: 3195.00, Credit: 405.00, Proceeds: -1242.00, Profit: -863.60, Commission: 26.60, UserId: 1, PercentGain: -27.03, BrokerAccountId: 3, BrokerAccountRef: "test12345", Name: "Trade #18 - Put Credit Spread Trade", Status: "Closed", Type: "Put Credit Spread", OrderIds: "1", Note: "Test note", OpenDate: ts, ClosedDate: time.Date(2017, 05, 13, 17, 20, 01, 507451, time.UTC)})
	db.Create(&models.TradeGroup{Risked: 1907.00, Credit: 253.00, Proceeds: -33.00, Profit: 244.60, Commission: 15.40, UserId: 1, PercentGain: 12.51, BrokerAccountId: 3, BrokerAccountRef: "test12345", Name: "Trade #19 - Put Credit Spread Trade", Status: "Closed", Type: "Put Credit Spread", OrderIds: "1", Note: "Test note", OpenDate: ts, ClosedDate: time.Date(2016, 05, 14, 17, 20, 01, 507451, time.UTC)})
	db.Create(&models.TradeGroup{Risked: 1907.00, Credit: 253.00, Proceeds: -33.00, Profit: 244.60, Commission: 15.40, UserId: 1, PercentGain: 12.51, BrokerAccountId: 3, BrokerAccountRef: "test12345", Name: "Trade #20 - Put Credit Spread Trade", Status: "Closed", Type: "Put Credit Spread", OrderIds: "1", Note: "Test note", OpenDate: ts, ClosedDate: time.Date(2016, 06, 15, 17, 20, 01, 507451, time.UTC)})
	db.Create(&models.TradeGroup{Risked: 1947.00, Credit: 253.00, Proceeds: -33.00, Profit: 204.60, Commission: 15.40, UserId: 1, PercentGain: 10.51, BrokerAccountId: 3, BrokerAccountRef: "test12345", Name: "Trade #21 - Put Credit Spread Trade", Status: "Closed", Type: "Put Credit Spread", OrderIds: "1", Note: "Test note", OpenDate: ts, ClosedDate: time.Date(2016, 06, 16, 17, 20, 01, 507451, time.UTC)})
	db.Create(&models.TradeGroup{Risked: 1947.00, Credit: 253.00, Proceeds: -33.00, Profit: 204.60, Commission: 15.40, UserId: 1, PercentGain: 10.51, BrokerAccountId: 3, BrokerAccountRef: "test12345", Name: "Trade #22 - Put Credit Spread Trade", Status: "Closed", Type: "Put Credit Spread", OrderIds: "1", Note: "Test note", OpenDate: ts, ClosedDate: time.Date(2016, 07, 17, 17, 20, 01, 507451, time.UTC)})
	db.Create(&models.TradeGroup{Risked: 1947.00, Credit: 253.00, Proceeds: -33.00, Profit: 204.60, Commission: 15.40, UserId: 1, PercentGain: 10.51, BrokerAccountId: 3, BrokerAccountRef: "test12345", Name: "Trade #23 - Put Credit Spread Trade", Status: "Closed", Type: "Put Credit Spread", OrderIds: "1", Note: "Test note", OpenDate: ts, ClosedDate: time.Date(2016, 07, 18, 17, 20, 01, 507451, time.UTC)})

	return brokerAccount
}

//
// Get a yearly summary based on account, year
//
func TestGetYearlySummaryByAccountYear01(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()
	defer db.Close()

	// Get test data
	brokerAccount := LoadAndGetTestData()

	// Run function we are testing. - 2016
	summary2016 := GetYearlySummaryByAccountYear(db, brokerAccount, 2016)

	// Verify the data was return as expected
	st.Expect(t, summary2016.Year, 2016)
	st.Expect(t, summary2016.TotalTrades, 5)
	st.Expect(t, summary2016.LossCount, 0)
	st.Expect(t, summary2016.WinCount, 5)
	st.Expect(t, summary2016.Profit, 1103.0)
	st.Expect(t, summary2016.Commission, 77.0)
	st.Expect(t, summary2016.WinPercent, 100.00)
	st.Expect(t, summary2016.LossPercent, 0.00)
	st.Expect(t, summary2016.ProfitStd, 19.60)
	st.Expect(t, summary2016.PercentGainStd, 0.98)
	st.Expect(t, summary2016.SharpeRatio, 9.47)
	st.Expect(t, summary2016.AvgRisked, 1931.0)
	st.Expect(t, summary2016.AvgPercentGain, 11.31)

	// Run function we are testing. - 2017
	summary2017 := GetYearlySummaryByAccountYear(db, brokerAccount, 2017)

	// Verify the data was return as expected
	st.Expect(t, summary2017.Year, 2017)
	st.Expect(t, summary2017.TotalTrades, 18)
	st.Expect(t, summary2017.LossCount, 4)
	st.Expect(t, summary2017.WinCount, 14)
	st.Expect(t, summary2017.Profit, 571.2)
	st.Expect(t, summary2017.Commission, 296.8)
	st.Expect(t, summary2017.WinPercent, 77.78)
	st.Expect(t, summary2017.LossPercent, 22.22)
	st.Expect(t, summary2017.ProfitStd, 337.7)
	st.Expect(t, summary2017.PercentGainStd, 12.43)
	st.Expect(t, summary2017.SharpeRatio, 0.1)
	st.Expect(t, summary2017.AvgRisked, 2078.83)
	st.Expect(t, summary2017.AvgPercentGain, 3.3)
}

/* End File */
