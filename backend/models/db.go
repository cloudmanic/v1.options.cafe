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

	// BrokerAccounts
	db.Exec("TRUNCATE TABLE broker_accounts;")
	db.Create(&BrokerAccount{UserId: 1, BrokerId: 2, Name: "Test Account 1", AccountNumber: "YYY123ZY", StockCommission: 5.00, StockMin: 0.00, OptionCommission: 0.35, OptionSingleMin: 5.00, OptionMultiLegMin: 7.00, OptionBase: 0.00})
	db.Create(&BrokerAccount{UserId: 1, BrokerId: 2, Name: "Test Account 2", AccountNumber: "ABC123ZY", StockCommission: 5.00, StockMin: 0.00, OptionCommission: 0.35, OptionSingleMin: 5.00, OptionMultiLegMin: 7.00, OptionBase: 0.00})

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
	db.Create(&Symbol{Name: "BARCLAYS BK PLC", ShortName: "VXX", Type: "Equity"})

	db.Create(&Symbol{Name: "SPY Feb 9, 2018 $276.00 Put", ShortName: "SPY180209P00276000", Type: "Option"})
	db.Create(&Symbol{Name: "VXX Mar 2, 2018 $46.00 Put", ShortName: "VXX180302P00046000", Type: "Option"})
	db.Create(&Symbol{Name: "VXX Feb 23, 2018 $50.00 Put", ShortName: "VXX180223P00050000", Type: "Option"})

	// TradeGroups TODO: make this more complete
	db.Exec("TRUNCATE TABLE trade_groups;")
	db.Create(&TradeGroup{UserId: 1, BrokerAccountId: 1, AccountId: "abc123", Status: "Open", Type: "Put Credit Spread", OrderIds: "1", Risked: 0.00, Proceeds: 0.00, Profit: 0.00, Commission: 23.45, Note: "Test note #1", OpenDate: ts})

	// Positions TODO: Put better values in here.
	db.Exec("TRUNCATE TABLE positions;")
	db.Create(&Position{UserId: 1, TradeGroupId: 1, AccountId: "123abc", Status: "Open", SymbolId: 4, Qty: 10, OrgQty: 10, CostBasis: 1000.00, AvgOpenPrice: 1.00, AvgClosePrice: 0.00, OrderIds: "1", OpenDate: ts, Note: "Test note #1"})
	db.Create(&Position{UserId: 1, TradeGroupId: 1, AccountId: "123abc", Status: "Open", SymbolId: 6, Qty: 10, OrgQty: 10, CostBasis: 1000.00, AvgOpenPrice: 1.00, AvgClosePrice: 0.00, OrderIds: "1", OpenDate: ts, Note: "Test note #2"})

	// Orders
	db.Exec("TRUNCATE TABLE orders;")
	db.Exec(`INSERT INTO orders (id, user_id, created_at, updated_at, broker_id, broker_ref, account_id, type, symbol_id, option_symbol_id, side, qty, status, duration, price, avg_fill_price, exec_quantity, last_fill_price, last_fill_quantity, remaining_quantity, create_date, transaction_date, class, num_legs, position_reviewed) VALUES 
	(1,1,'2018-02-17 00:13:53','2018-02-17 00:13:54',2,'734801','ABC123ZY','limit',1,11,'buy_to_open',1,'filled','day',2.40,2.39,1.00,2.39,1.00,0.00,'2018-01-16 11:54:50','2018-01-16 11:54:51','option',0,'No'),
	(2,1,'2018-02-17 00:13:53','2018-02-17 00:13:54',2,'735196','ABC123ZY','limit',1,11,'sell_to_close',-1,'filled','gtc',3.39,3.62,1.00,3.62,1.00,0.00,'2018-01-16 15:29:51','2018-02-05 06:30:03','option',0,'No'),
	(3,1,'2018-02-17 00:13:53','2018-02-17 00:13:53',2,'767256','ABC123ZY','credit',1,0,'buy',9,'filled','day',0.24,-0.24,9.00,0.00,1.00,0.00,'2018-01-30 09:20:25','2018-01-30 09:45:11','multileg',2,'No'),
	(4,1,'2018-02-17 00:13:53','2018-02-17 00:13:53',2,'770154','ABC123ZY','credit',1,0,'buy',9,'filled','day',0.24,-0.24,9.00,0.00,9.00,0.00,'2018-01-31 09:39:20','2018-01-31 11:21:38','multileg',2,'No'),
	(5,1,'2018-02-17 00:13:53','2018-02-17 00:13:53',2,'772977','ABC123ZY','debit',1,0,'buy',9,'filled','gtc',0.80,0.72,9.00,0.00,9.00,0.00,'2018-02-01 09:30:04','2018-02-05 11:57:31','multileg',2,'No'),
	(6,1,'2018-02-17 00:13:53','2018-02-17 00:13:53',2,'773020','ABC123ZY','credit',1,0,'buy',9,'filled','day',0.21,-0.21,9.00,0.00,9.00,0.00,'2018-02-01 09:40:44','2018-02-01 10:42:33','multileg',2,'No'),
	(7,1,'2018-02-17 00:13:53','2018-02-17 00:13:53',2,'775404','ABC123ZY','debit',1,0,'buy',18,'filled','gtc',0.80,0.69,18.00,0.00,18.00,0.00,'2018-02-02 07:27:02','2018-02-05 11:58:57','multileg',2,'No'),
	(8,1,'2018-02-17 00:13:53','2018-02-17 00:13:53',2,'775734','ABC123ZY','credit',1,0,'buy',9,'filled','day',0.23,-0.23,9.00,0.00,9.00,0.00,'2018-02-02 08:07:00','2018-02-02 08:19:15','multileg',2,'No'),
	(9,1,'2018-02-17 00:13:53','2018-02-17 00:13:53',2,'780197','ABC123ZY','credit',1,0,'buy',9,'filled','day',0.25,-0.25,9.00,0.00,9.00,0.00,'2018-02-05 08:01:06','2018-02-05 08:17:18','multileg',2,'No'),
	(10,1,'2018-02-17 00:13:53','2018-02-17 00:13:53',2,'781119','ABC123ZY','credit',1,0,'buy',9,'filled','day',0.22,-0.22,9.00,0.00,9.00,0.00,'2018-02-05 10:40:10','2018-02-05 11:26:06','multileg',2,'No'),
	(11,1,'2018-02-17 00:13:53','2018-02-17 00:13:54',2,'781816','ABC123ZY','credit',1,0,'buy',9,'filled','day',0.35,-0.35,9.00,0.00,1.00,0.00,'2018-02-05 12:11:47','2018-02-05 12:14:14','multileg',2,'No'),
	(12,1,'2018-02-17 00:13:53','2018-02-17 00:13:54',2,'784720','ABC123ZY','limit',10,12,'buy_to_open',2,'filled','day',6.60,6.45,2.00,6.45,1.00,0.00,'2018-02-06 09:36:01','2018-02-06 09:36:58','option',0,'No'),
	(13,1,'2018-02-17 00:13:53','2018-02-17 00:13:54',2,'784730','ABC123ZY','credit',10,0,'buy',2,'filled','day',1.20,-1.44,2.00,0.00,2.00,0.00,'2018-02-06 09:39:06','2018-02-06 09:39:06','multileg',2,'No'),
	(14,1,'2018-02-17 00:13:53','2018-02-17 00:13:54',2,'790726','ABC123ZY','credit',1,0,'buy',9,'filled','day',0.31,-0.31,9.00,0.00,9.00,0.00,'2018-02-08 09:53:10','2018-02-08 09:54:24','multileg',2,'No'),
	(15,1,'2018-02-17 00:13:53','2018-02-17 00:13:54',2,'790833','ABC123ZY','debit',1,0,'buy',9,'filled','gtc',0.03,0.03,9.00,0.00,9.00,0.00,'2018-02-08 10:09:37','2018-02-16 09:35:13','multileg',2,'No'),
	(16,1,'2018-02-17 00:13:53','2018-02-17 00:13:54',2,'791126','ABC123ZY','limit',10,13,'buy_to_open',2,'filled','day',5.00,5.00,2.00,5.00,2.00,0.00,'2018-02-08 11:05:07','2018-02-08 12:35:55','option',0,'No'),
	(17,1,'2018-02-17 00:13:53','2018-02-17 00:13:54',2,'793756','ABC123ZY','credit',1,0,'buy',9,'filled','day',0.28,-0.28,9.00,0.00,9.00,0.00,'2018-02-09 09:06:50','2018-02-09 09:08:48','multileg',2,'No'),
	(18,1,'2018-02-17 00:13:53','2018-02-17 00:13:54',2,'793804','ABC123ZY','limit',10,13,'sell_to_close',-2,'filled','gtc',9.00,9.00,2.00,9.00,2.00,0.00,'2018-02-09 09:11:52','2018-02-16 10:35:53','option',0,'No'),
	(19,1,'2018-02-17 00:13:53','2018-02-17 00:13:54',2,'797557','ABC123ZY','debit',1,0,'buy',9,'filled','gtc',0.03,0.03,9.00,0.00,9.00,0.00,'2018-02-12 09:09:44','2018-02-14 06:35:46','multileg',2,'No'),
	(20,1,'2018-02-17 00:13:53','2018-02-17 00:13:54',2,'802693','ABC123ZY','credit',1,0,'buy',9,'filled','day',0.29,-0.29,9.00,0.00,9.00,0.00,'2018-02-14 09:03:46','2018-02-14 09:10:39','multileg',2,'No');`)

	// OrderLegs
	db.Exec("TRUNCATE TABLE order_legs;")
	db.Exec(`INSERT INTO order_legs (id, user_id, created_at, updated_at, order_id, type, symbol_id, side, qty, status, duration, avg_fill_price, exec_quantity, last_fill_price, last_fill_quantity, remaining_quantity, create_date, transaction_date) VALUES 
	(1,1,'2018-02-17 00:13:53','2018-02-17 00:13:53',3,'credit',36007,'buy_to_open',9,'filled','day',1.62,9.00,1.65,1.00,0.00,'2018-01-30 09:20:25','2018-01-30 09:45:11'),
	(2,1,'2018-02-17 00:13:53','2018-02-17 00:13:53',3,'credit',36008,'sell_to_open',9,'filled','day',1.86,9.00,1.89,1.00,0.00,'2018-01-30 09:20:25','2018-01-30 09:45:11'),
	(3,1,'2018-02-17 00:13:53','2018-02-17 00:13:53',4,'credit',36009,'buy_to_open',9,'filled','day',1.49,9.00,1.49,3.00,0.00,'2018-01-31 09:39:20','2018-01-31 11:21:38'),
	(4,1,'2018-02-17 00:13:53','2018-02-17 00:13:53',4,'credit',36011,'sell_to_open',9,'filled','day',1.73,9.00,1.73,9.00,0.00,'2018-01-31 09:39:20','2018-01-31 11:21:38'),
	(5,1,'2018-02-17 00:13:53','2018-02-17 00:13:53',5,'debit',36007,'sell_to_close',9,'filled','gtc',6.27,9.00,6.27,9.00,0.00,'2018-02-01 09:30:04','2018-02-05 11:57:31'),
	(6,1,'2018-02-17 00:13:53','2018-02-17 00:13:53',5,'debit',36008,'buy_to_close',9,'filled','gtc',6.99,9.00,6.99,9.00,0.00,'2018-02-01 09:30:04','2018-02-05 11:57:31'),
	(7,1,'2018-02-17 00:13:53','2018-02-17 00:13:53',6,'credit',36009,'buy_to_open',9,'filled','day',1.22,9.00,1.22,1.00,0.00,'2018-02-01 09:40:44','2018-02-01 10:42:33'),
	(8,1,'2018-02-17 00:13:53','2018-02-17 00:13:53',6,'credit',36011,'sell_to_open',9,'filled','day',1.43,9.00,1.43,9.00,0.00,'2018-02-01 09:40:44','2018-02-01 10:42:33'),
	(9,1,'2018-02-17 00:13:53','2018-02-17 00:13:53',7,'debit',36009,'sell_to_close',18,'filled','gtc',6.82,18.00,6.82,18.00,0.00,'2018-02-02 07:27:02','2018-02-05 11:58:57'),
	(10,1,'2018-02-17 00:13:53','2018-02-17 00:13:53',7,'debit',36011,'buy_to_close',18,'filled','gtc',7.51,18.00,7.51,18.00,0.00,'2018-02-02 07:27:02','2018-02-05 11:58:57'),
	(11,1,'2018-02-17 00:13:53','2018-02-17 00:13:53',8,'credit',36013,'buy_to_open',9,'filled','day',1.48,9.00,1.48,9.00,0.00,'2018-02-02 08:07:00','2018-02-02 08:19:15'),
	(12,1,'2018-02-17 00:13:53','2018-02-17 00:13:53',8,'credit',36015,'sell_to_open',9,'filled','day',1.71,9.00,1.71,9.00,0.00,'2018-02-02 08:07:00','2018-02-02 08:19:15'),
	(13,1,'2018-02-17 00:13:53','2018-02-17 00:13:53',9,'credit',36017,'buy_to_open',9,'filled','day',1.56,9.00,1.56,9.00,0.00,'2018-02-05 08:01:06','2018-02-05 08:17:18'),
	(14,1,'2018-02-17 00:13:53','2018-02-17 00:13:53',9,'credit',36018,'sell_to_open',9,'filled','day',1.81,9.00,1.81,9.00,0.00,'2018-02-05 08:01:06','2018-02-05 08:17:18'),
	(15,1,'2018-02-17 00:13:53','2018-02-17 00:13:53',10,'credit',36021,'buy_to_open',9,'filled','day',2.33,9.00,2.33,9.00,0.00,'2018-02-05 10:40:10','2018-02-05 11:26:06'),
	(16,1,'2018-02-17 00:13:53','2018-02-17 00:13:53',10,'credit',36022,'sell_to_open',9,'filled','day',2.55,9.00,2.55,9.00,0.00,'2018-02-05 10:40:10','2018-02-05 11:26:06'),
	(17,1,'2018-02-17 00:13:53','2018-02-17 00:13:54',11,'credit',36024,'buy_to_open',9,'filled','day',2.56,9.00,4.55,1.00,0.00,'2018-02-05 12:11:47','2018-02-05 12:14:14'),
	(18,1,'2018-02-17 00:13:53','2018-02-17 00:13:54',11,'credit',36025,'sell_to_open',9,'filled','day',2.91,9.00,4.90,1.00,0.00,'2018-02-05 12:11:47','2018-02-05 12:14:14'),
	(19,1,'2018-02-17 00:13:53','2018-02-17 00:13:54',13,'credit',36026,'sell_to_open',2,'filled','day',6.35,2.00,6.35,2.00,0.00,'2018-02-06 09:39:06','2018-02-06 09:39:06'),
	(20,1,'2018-02-17 00:13:53','2018-02-17 00:13:54',13,'credit',36027,'buy_to_open',2,'filled','day',4.91,2.00,4.91,2.00,0.00,'2018-02-06 09:39:06','2018-02-06 09:39:06'),
	(21,1,'2018-02-17 00:13:53','2018-02-17 00:13:54',14,'credit',36028,'buy_to_open',9,'filled','day',3.28,9.00,3.28,9.00,0.00,'2018-02-08 09:53:10','2018-02-08 09:54:24'),
	(22,1,'2018-02-17 00:13:53','2018-02-17 00:13:54',14,'credit',36029,'sell_to_open',9,'filled','day',3.59,9.00,3.59,9.00,0.00,'2018-02-08 09:53:10','2018-02-08 09:54:24'),
	(23,1,'2018-02-17 00:13:53','2018-02-17 00:13:54',15,'debit',36028,'sell_to_close',9,'filled','gtc',0.14,9.00,0.14,4.00,0.00,'2018-02-08 10:09:37','2018-02-16 09:35:13'),
	(24,1,'2018-02-17 00:13:53','2018-02-17 00:13:54',15,'debit',36029,'buy_to_close',9,'filled','gtc',0.17,9.00,0.17,9.00,0.00,'2018-02-08 10:09:37','2018-02-16 09:35:13'),
	(25,1,'2018-02-17 00:13:53','2018-02-17 00:13:54',17,'credit',36030,'buy_to_open',9,'filled','day',2.26,9.00,2.26,9.00,0.00,'2018-02-09 09:06:50','2018-02-09 09:08:48'),
	(26,1,'2018-02-17 00:13:53','2018-02-17 00:13:54',17,'credit',36031,'sell_to_open',9,'filled','day',2.54,9.00,2.54,9.00,0.00,'2018-02-09 09:06:50','2018-02-09 09:08:48'),
	(27,1,'2018-02-17 00:13:53','2018-02-17 00:13:54',19,'debit',36030,'sell_to_close',9,'filled','gtc',0.10,9.00,0.10,9.00,0.00,'2018-02-12 09:09:44','2018-02-14 06:35:46'),
	(28,1,'2018-02-17 00:13:53','2018-02-17 00:13:54',19,'debit',36031,'buy_to_close',9,'filled','gtc',0.13,9.00,0.13,9.00,0.00,'2018-02-12 09:09:44','2018-02-14 06:35:46'),
	(29,1,'2018-02-17 00:13:53','2018-02-17 00:13:54',20,'credit',36032,'buy_to_open',9,'filled','day',2.17,9.00,2.17,9.00,0.00,'2018-02-14 09:03:46','2018-02-14 09:10:39'),
	(30,1,'2018-02-17 00:13:53','2018-02-17 00:13:54',20,'credit',36033,'sell_to_open',9,'filled','day',2.46,9.00,2.46,9.00,0.00,'2018-02-14 09:03:46','2018-02-14 09:10:39');`)
}

/* End File */
