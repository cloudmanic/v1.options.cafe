//
// Date: 9/30/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"errors"
	"go/build"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	env "github.com/jpfuentes2/go-env"

	"app.options.cafe/library/services"
)

//
// Init.
//
func init() {
	env.ReadEnv(build.Default.GOPATH + "/src/app.options.cafe/.env")
}

//
// NewDB - Start the DB connection.
//
func NewDB() (*DB, error) {
	// We should not be calling htis from testing.
	if strings.HasSuffix(os.Args[0], ".test") {
		log.Fatal(errors.New("We can not call NewDB() from testing."))
	}

	var err error

	// Build connection string
	conStr := os.Getenv("DB_USERNAME") + ":" + os.Getenv("DB_PASSWORD") + "@tcp(" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + ")/" + os.Getenv("DB_DATABASE")

	// Connect to Mysql
	db, err := gorm.Open("mysql", conStr+"?charset=utf8&parseTime=True&loc=Local")

	if err != nil {
		services.Fatal(errors.New(err.Error() + "Failed to connect database"))
	}

	// Enable
	// db.LogMode(true)
	// db.SetLogger(log.New(os.Stdout, "\r\n", 0))

	// Helpful settings
	//db.DB().SetMaxIdleConns(10)
	//db.DB().SetMaxOpenConns(100)

	// Run doMigrations
	doMigrations(db)

	// Ping every so often to keep the connection alive.
	go PingDbServer(db)

	// Return db connection.
	return &DB{db}, nil
}

//
// doMigrations - Run our migrations
//
func doMigrations(db *gorm.DB) {
	// Migrate the schemas (one per table).
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Gain{})
	db.AutoMigrate(&Broker{})
	db.AutoMigrate(&Order{})
	db.AutoMigrate(&Symbol{})
	db.AutoMigrate(&Session{})
	db.AutoMigrate(&Settings{})
	db.AutoMigrate(&OrderLeg{})
	db.AutoMigrate(&BrokerEvent{})
	db.AutoMigrate(&Watchlist{})
	db.AutoMigrate(&WatchlistSymbol{})
	db.AutoMigrate(&Position{})
	db.AutoMigrate(&Backtest{})
	db.AutoMigrate(&Screener{})
	db.AutoMigrate(&ScreenerItem{})
	db.AutoMigrate(&TradeGroup{})
	db.AutoMigrate(&ActiveSymbol{})
	db.AutoMigrate(&Application{})
	db.AutoMigrate(&HistoricalQuote{})
	db.AutoMigrate(&Notification{})
	db.AutoMigrate(&NotifyChannel{})
	db.AutoMigrate(&BrokerAccount{})
	db.AutoMigrate(&ForgotPassword{})
	db.AutoMigrate(&BalanceHistory{})
	db.AutoMigrate(&BacktestPosition{})
}

//
// PingDbServer - Just make a query on this connection every so often to keep it alive.
//
func PingDbServer(db *gorm.DB) {
	for {
		// Ping to keep server alive
		db.DB().Ping()

		// Sleep for X seconds
		time.Sleep(time.Second * 10)
	}
}

/* End File */
