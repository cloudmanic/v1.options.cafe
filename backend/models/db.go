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

}

/* End File */
