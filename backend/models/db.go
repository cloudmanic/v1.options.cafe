//
// Date: 9/30/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

//
// Start the DB connection.
//
func NewDB() (*DB, error) {

	var err error

	dbName := os.Getenv("DB_DATABASE")

	// Is this a testing run?
	if flag.Lookup("test.v") != nil {
		dbName = os.Getenv("DB_DATABASE_TESTING")
	}

	// Connect to Mysql
	db, err := gorm.Open("mysql", os.Getenv("DB_USERNAME")+":"+os.Getenv("DB_PASSWORD")+"@"+os.Getenv("DB_HOST")+"/"+dbName+"?charset=utf8&parseTime=True&loc=Local")

	if err != nil {
		services.Error(err, "Failed to connect database")
		log.Fatal(err)
	}

	// Enable
	//db.LogMode(true)
	//db.SetLogger(log.New(os.Stdout, "\r\n", 0))

	// Migrate the schemas (one per table).
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Broker{})
	db.AutoMigrate(&Order{})
	db.AutoMigrate(&Symbol{})
	db.AutoMigrate(&Session{})
	db.AutoMigrate(&OrderLeg{})
	db.AutoMigrate(&Watchlist{})
	db.AutoMigrate(&WatchlistSymbol{})
	db.AutoMigrate(&Position{})
	db.AutoMigrate(&TradeGroup{})
	db.AutoMigrate(&BrokerAccount{})
	db.AutoMigrate(&ForgotPassword{})

	// Is this a testing run? If so load testing data.
	if flag.Lookup("test.v") != nil {
		LoadTestingData(db)
	}

	// Return db connection.
	return &DB{db}, nil
}

//
// Load testing data.
//
func LoadTestingData(db *gorm.DB) {

	// Shared vars we use.
	ts := time.Date(2017, 10, 29, 17, 20, 01, 507451, time.UTC)

	// Users
	db.Exec("TRUNCATE TABLE users;")
	db.Create(&User{FirstName: "Rob", LastName: "Tester", Email: "spicer+robtester@options.cafe", Status: "Active"})
	db.Create(&User{FirstName: "Jane", LastName: "Wells", Email: "spicer+janewells@options.cafe", Status: "Active"})
	db.Create(&User{FirstName: "Bob", LastName: "Rosso", Email: "spicer+bobrosso@options.cafe", Status: "Active"})

	// Brokers
	db.Exec("TRUNCATE TABLE brokers;")
	db.Create(&Broker{Name: "Tradier", UserId: 1, AccessToken: "123", RefreshToken: "abc", TokenExpirationDate: ts})
	db.Create(&Broker{Name: "Tradeking", UserId: 1, AccessToken: "456", RefreshToken: "xyz", TokenExpirationDate: ts})
	db.Create(&Broker{Name: "Etrade", UserId: 1, AccessToken: "789", RefreshToken: "mno", TokenExpirationDate: ts})

	// Symbols
	db.Exec("TRUNCATE TABLE symbols;")
	db.Create(&Symbol{Name: "SPDR S&P 500 ETF Trust", ShortName: "SPY", Type: "Equity"})
	db.Create(&Symbol{Name: "McDonald's Corp", ShortName: "MCD", Type: "Equity"})
	db.Create(&Symbol{Name: "Starbucks Corp", ShortName: "SBUX", Type: "Equity"})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $253.00 Put", ShortName: "SPY180316P00253000", Type: "Option"})
	db.Create(&Symbol{Name: "VXX Feb 23, 2018 $55.00 Call", ShortName: "VXX180223C00055000", Type: "Option"})
	db.Create(&Symbol{Name: "SPY Mar 16, 2018 $266.00 Put", ShortName: "SPY180316P00266000", Type: "Option"})
	db.Create(&Symbol{Name: "Caterpillar Inc", ShortName: "CAT", Type: "Equity"})
	db.Create(&Symbol{Name: "Ascent Solar Technologies Inc", ShortName: "ASTI", Type: "Equity"})
	db.Create(&Symbol{Name: "Advanced Micro Devices Inc", ShortName: "AMD", Type: "Equity"})

	// Orders TODO: make this more complete
	db.Exec("TRUNCATE TABLE orders;")
	db.Create(&Order{UserId: 1, BrokerId: 1, BrokerRef: "734801", AccountId: "123abc", Type: "limit"})

	// TradeGroups TODO: make this more complete
	db.Exec("TRUNCATE TABLE trade_groups;")
	db.Create(&TradeGroup{UserId: 1, BrokerAccountId: 1, AccountId: "abc123", Status: "Open", Type: "Put Credit Spread", OrderIds: "1", Risked: 0.00, Gain: 0.00, Profit: 0.00, Commission: 23.45, Note: "Test note #1", OpenDate: ts})

	// Positions TODO: Put better values in here.
	db.Exec("TRUNCATE TABLE positions;")
	db.Create(&Position{UserId: 1, TradeGroupId: 1, AccountId: "123abc", Status: "Open", SymbolId: 4, Qty: 10, OrgQty: 10, CostBasis: 1000.00, AvgOpenPrice: 1.00, AvgClosePrice: 0.00, OrderIds: "1", OpenDate: ts, Note: "Test note #1"})
	db.Create(&Position{UserId: 1, TradeGroupId: 1, AccountId: "123abc", Status: "Open", SymbolId: 6, Qty: 10, OrgQty: 10, CostBasis: 1000.00, AvgOpenPrice: 1.00, AvgClosePrice: 0.00, OrderIds: "1", OpenDate: ts, Note: "Test note #2"})
}

/* End File */
