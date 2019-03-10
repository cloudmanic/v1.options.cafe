//
// Date: 7/18/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/nbio/st"

	"github.com/cloudmanic/app.options.cafe/backend/library/reports"
	"github.com/cloudmanic/app.options.cafe/backend/library/test"
	"github.com/cloudmanic/app.options.cafe/backend/models"
)

//
// Load and return some test data.
//
func LoadAndGetTestData(db *gorm.DB) {
	// BrokerAccounts
	db.Create(&models.BrokerAccount{UserId: 1, BrokerId: 1, Name: "Test Account 1", AccountNumber: "abc1235423", StockCommission: 5.00, StockMin: 0.00, OptionCommission: 0.35, OptionSingleMin: 5.00, OptionMultiLegMin: 7.00, OptionBase: 0.00})
	db.Create(&models.BrokerAccount{UserId: 1, BrokerId: 1, Name: "Test Account 1", AccountNumber: "abc1235423", StockCommission: 5.00, StockMin: 0.00, OptionCommission: 0.35, OptionSingleMin: 5.00, OptionMultiLegMin: 7.00, OptionBase: 0.00})

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
}

//
// TestReportsGetTradeGroupYears01
//
func TestReportsGetTradeGroupYears01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)
	//db, _, _ := models.NewTestDB("testing_db")

	// Load test data
	LoadAndGetTestData(db.New())

	// Create controller
	c := &Controller{DB: db}

	// Make a mock request.
	req, _ := http.NewRequest("GET", "/api/v1/reports/3/tradegroup/years", nil)
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(1)) })

	r.GET("/api/v1/reports/:brokerAccount/tradegroup/years", c.ReportsGetTradeGroupYears)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	var result []int
	err := json.Unmarshal([]byte(w.Body.String()), &result)

	// Parse json that returned.
	st.Expect(t, err, nil)
	st.Expect(t, w.Code, 200)
	st.Expect(t, w.Body.String(), `[2017,2016]`)
}

//
// TestReportsGetAccountYearlySummary01
//
func TestReportsGetAccountYearlySummary01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Load test data
	LoadAndGetTestData(db.New())

	// Create controller
	c := &Controller{DB: db}

	// Make a mock request.
	req, _ := http.NewRequest("GET", "/api/v1/reports/3/summary/yearly/2017", nil)
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(1)) })

	r.GET("/api/v1/reports/:brokerAccount/summary/yearly/:year", c.ReportsGetAccountYearlySummary)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	result := reports.YearlySummary{}
	err := json.Unmarshal([]byte(w.Body.String()), &result)

	// Parse json that returned.
	st.Expect(t, err, nil)
	st.Expect(t, w.Code, 200)
	st.Expect(t, w.Body.String(), `{"year":2017,"total_trades":18,"loss_count":4,"win_count":14,"profit":571.2,"commission":296.8,"win_percent":77.78,"loss_percent":22.22,"profit_std":337.7,"precent_gain_std":12.43,"sharpe_ratio":0.1,"avg_risked":2078.83,"avg_percent_gain":3.3}`)
}

//
// TestReportsGetAccountYearlySummary02
//
func TestReportsGetAccountYearlySummary02(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Load test data
	LoadAndGetTestData(db.New())

	// Create controller
	c := &Controller{DB: db}

	// Make a mock request.
	req, _ := http.NewRequest("GET", "/api/v1/reports/3/summary/yearly/2017", nil)
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(1)) })

	r.GET("/api/v1/reports/:brokerAccount/summary/yearly/:year", c.ReportsGetAccountYearlySummary)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	result := reports.YearlySummary{}
	err := json.Unmarshal([]byte(w.Body.String()), &result)

	// Parse json that returned.
	st.Expect(t, err, nil)
	st.Expect(t, w.Code, 200)
	st.Expect(t, w.Body.String(), `{"year":2017,"total_trades":18,"loss_count":4,"win_count":14,"profit":571.2,"commission":296.8,"win_percent":77.78,"loss_percent":22.22,"profit_std":337.7,"precent_gain_std":12.43,"sharpe_ratio":0.1,"avg_risked":2078.83,"avg_percent_gain":3.3}`)
}

//
// TestReportsGetAccountReturns01
//
func TestReportsGetAccountReturns01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// BrokerAccounts
	db.Create(&models.BrokerAccount{UserId: 1, BrokerId: 1, Name: "Test Account 1", AccountNumber: "abc1235423", StockCommission: 5.00, StockMin: 0.00, OptionCommission: 0.35, OptionSingleMin: 5.00, OptionMultiLegMin: 7.00, OptionBase: 0.00})
	db.Create(&models.BrokerAccount{UserId: 1, BrokerId: 1, Name: "Test Account 1", AccountNumber: "abc1235423", StockCommission: 5.00, StockMin: 0.00, OptionCommission: 0.35, OptionSingleMin: 5.00, OptionMultiLegMin: 7.00, OptionBase: 0.00})

	// Put test data into database - We want broker account id = 3 to match with the sql dump data
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

	// Load testing data.
	err := test.LoadSqlDump(db, "balance_histories_1")
	st.Expect(t, err, nil)
	err = test.LoadSqlDump(db, "broker_events_1")
	st.Expect(t, err, nil)

	// Create controller
	c := &Controller{DB: db}

	// Make a mock request.
	req, _ := http.NewRequest("GET", "/api/v1/reports/3/account-returns?start=2018-12-20&end=2018-12-31", nil)
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(1)) })

	r.GET("/api/v1/reports/:brokerAccount/account-returns", c.ReportsGetAccountReturns)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	results := []reports.Returns{}
	err = json.Unmarshal([]byte(w.Body.String()), &results)

	// Parse json that returned.
	st.Expect(t, err, nil)
	st.Expect(t, w.Code, 200)
	st.Expect(t, len(results), 12)
	st.Expect(t, w.Body.String(), `[{"date":"2018-12-20","percent":0,"account_value":3007.22,"total_cash":8096.72,"price_per":1,"units":3007.22},{"date":"2018-12-21","percent":-0.29339389868383414,"account_value":2124.92,"total_cash":8056.42,"price_per":0.7066061013161659,"units":3714.8278158236544},{"date":"2018-12-22","percent":-0.5330416143082024,"account_value":1734.67,"total_cash":3612.67,"price_per":0.4669583856917976,"units":3714.8278158236544},{"date":"2018-12-23","percent":-0.5330416143082024,"account_value":1734.67,"total_cash":3612.67,"price_per":0.4669583856917976,"units":3714.8278158236544},{"date":"2018-12-24","percent":-0.05334239583853195,"account_value":3516.67,"total_cash":3612.67,"price_per":0.946657604161468,"units":3714.8278158236544},{"date":"2018-12-25","percent":-0.05334239583853195,"account_value":3516.67,"total_cash":3612.67,"price_per":0.946657604161468,"units":3714.8278158236544},{"date":"2018-12-26","percent":-0.2374693685844661,"account_value":2832.67,"total_cash":3612.67,"price_per":0.7625306314155339,"units":3714.8278158236544},{"date":"2018-12-27","percent":-0.5160825510289716,"account_value":1797.67,"total_cash":3612.67,"price_per":0.4839174489710284,"units":3714.8278158236544},{"date":"2018-12-28","percent":-0.5233507210057848,"account_value":1770.67,"total_cash":3612.67,"price_per":0.47664927899421516,"units":3714.8278158236544},{"date":"2018-12-29","percent":-0.5233507210057848,"account_value":1770.67,"total_cash":3612.67,"price_per":0.47664927899421516,"units":3714.8278158236544},{"date":"2018-12-30","percent":-0.5233507210057848,"account_value":1770.67,"total_cash":3612.67,"price_per":0.47664927899421516,"units":3714.8278158236544},{"date":"2018-12-31","percent":-0.5235122358941584,"account_value":1770.07,"total_cash":2856.07,"price_per":0.47648776410584154,"units":3714.8278158236544}]`)
}

//
// TestReportsGetAccountReturns02
//
func TestReportsGetAccountReturns02(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// BrokerAccounts
	db.Create(&models.BrokerAccount{UserId: 1, BrokerId: 1, Name: "Test Account 1", AccountNumber: "abc1235423", StockCommission: 5.00, StockMin: 0.00, OptionCommission: 0.35, OptionSingleMin: 5.00, OptionMultiLegMin: 7.00, OptionBase: 0.00})
	db.Create(&models.BrokerAccount{UserId: 1, BrokerId: 1, Name: "Test Account 1", AccountNumber: "abc1235423", StockCommission: 5.00, StockMin: 0.00, OptionCommission: 0.35, OptionSingleMin: 5.00, OptionMultiLegMin: 7.00, OptionBase: 0.00})

	// Put test data into database - We want broker account id = 3 to match with the sql dump data
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

	// Load testing data.
	err := test.LoadSqlDump(db, "balance_histories_1")
	st.Expect(t, err, nil)
	err = test.LoadSqlDump(db, "broker_events_1")
	st.Expect(t, err, nil)

	// Create controller
	c := &Controller{DB: db}

	// Make a mock request.
	req, _ := http.NewRequest("GET", "/api/v1/reports/3/account-returns?start=2018-01-01&end=2018-12-31", nil)
	req.Header.Set("Accept", "application/json")

	// Setup GIN Router
	gin.SetMode("release")
	gin.DisableConsoleColor()
	r := gin.New()

	r.Use(func(c *gin.Context) { c.Set("userId", uint(1)) })

	r.GET("/api/v1/reports/:brokerAccount/account-returns", c.ReportsGetAccountReturns)

	// Setup writer.
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	results := []reports.Returns{}
	err = json.Unmarshal([]byte(w.Body.String()), &results)

	// Review Results
	st.Expect(t, len(results), 109)
	st.Expect(t, err, nil)
	st.Expect(t, w.Code, 200)

	st.Expect(t, results[0].TotalCash, 4948.86)
	st.Expect(t, results[0].PricePer, 1.00)
	st.Expect(t, results[0].Percent, 0.00)
	st.Expect(t, results[0].Units, 4844.36)
	st.Expect(t, results[0].AccountValue, 4844.36)
	st.Expect(t, results[0].Date.Format("2006-01-02"), "2018-09-14")

	st.Expect(t, results[20].TotalCash, 5128.32)
	st.Expect(t, results[20].PricePer, 1.0166090051110983)
	st.Expect(t, results[20].Percent, 0.016609005111098307)
	st.Expect(t, results[20].Units, 4844.36)
	st.Expect(t, results[20].AccountValue, 4924.82)
	st.Expect(t, results[20].Date.Format("2006-01-02"), "2018-10-04")

	st.Expect(t, results[30].TotalCash, 7092.32)
	st.Expect(t, results[30].PricePer, 0.8793378557103995)
	st.Expect(t, results[30].Percent, -0.12066214428960054)
	st.Expect(t, results[30].Units, 6257.913228987943)
	st.Expect(t, results[30].AccountValue, 5502.82)
	st.Expect(t, results[30].Date.Format("2006-01-02"), "2018-10-14")

	st.Expect(t, results[108].TotalCash, 2856.07)
	st.Expect(t, results[108].PricePer, 0.16260125625509866)
	st.Expect(t, results[108].Percent, -0.8373987437449013)
	st.Expect(t, results[108].Units, 10885.955255001274)
	st.Expect(t, results[108].AccountValue, 1770.07)
	st.Expect(t, results[108].Date.Format("2006-01-02"), "2018-12-31")
}

/* End File */
