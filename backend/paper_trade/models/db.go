//
// Date: 2018-11-05
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-11-05
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"flag"
	"go/build"
	"log"
	"os"

	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/jinzhu/gorm"
	env "github.com/jpfuentes2/go-env"
)

//
// Init.
//
func init() {
	env.ReadEnv(build.Default.GOPATH + "/src/github.com/cloudmanic/app.options.cafe/backend/.env")
}

//
// Start the DB connection.
//
func NewDB() (*DB, error) {

	var err error

	dbName := os.Getenv("PT_DB_DATABASE")

	// Is this a testing run?
	if flag.Lookup("test.v") != nil {
		dbName = os.Getenv("PT_DB_DATABASE_TESTING")
	}

	// Connect to Mysql
	db, err := gorm.Open("mysql", os.Getenv("PT_DB_USERNAME")+":"+os.Getenv("PT_DB_PASSWORD")+"@"+os.Getenv("PT_DB_HOST")+"/"+dbName+"?charset=utf8&parseTime=True&loc=Local")

	if err != nil {
		services.Error(err, "Failed to connect database")
		log.Fatal(err)
	}

	// Migrate the schemas (one per table).
	db.AutoMigrate(&Symbol{})
	db.AutoMigrate(&Account{})
	db.AutoMigrate(&Position{})

	// Is this a testing run? If so load testing data.
	if flag.Lookup("test.v") != nil {
		ClearTestingData(db)
	}

	// Return db connection.
	return &DB{db}, nil
}

//
// Clear Testing Data
//
func ClearTestingData(db *gorm.DB) {
	db.Exec("TRUNCATE TABLE symbols;")
	db.Exec("TRUNCATE TABLE accounts;")
	db.Exec("TRUNCATE TABLE positions;")
}
