//
// Date: 9/30/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"log"
	"os"

	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

//
// Start the DB connection.
//
func NewDB() (*DB, error) {

	var err error

	// Connect to Mysql
	db, err := gorm.Open("mysql", os.Getenv("DB_USERNAME")+":"+os.Getenv("DB_PASSWORD")+"@"+os.Getenv("DB_HOST")+"/"+os.Getenv("DB_DATABASE")+"?charset=utf8&parseTime=True&loc=Local")

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

	// Return db connection.
	return &DB{db}, nil
}

/* End File */
