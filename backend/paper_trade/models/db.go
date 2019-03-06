//
// Date: 2018-11-05
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-11-05
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"errors"
	"go/build"
	"log"
	"os"
	"strings"

	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/jinzhu/gorm"
	env "github.com/jpfuentes2/go-env"
)

//
// Init.
//
func init() {
	env.ReadEnv(build.Default.GOPATH + "/src/github.com/cloudmanic/app.options.cafe/backend/paper_trade/.env")
}

//
// Start the DB connection.
//
func NewDB() (*DB, error) {
	// We should not be calling htis from testing.
	if strings.HasSuffix(os.Args[0], ".test") {
		log.Fatal(errors.New("We can not call NewDB() from testing."))
	}

	var err error

	dbName := os.Getenv("DB_DATABASE")

	// Connect to Mysql
	db, err := gorm.Open("mysql", os.Getenv("DB_USERNAME")+":"+os.Getenv("DB_PASSWORD")+"@"+os.Getenv("DB_HOST")+"/"+dbName+"?charset=utf8&parseTime=True&loc=Local")

	if err != nil {
		services.Error(err, "Failed to connect database")
		log.Fatal(err)
	}

	// Migrate the schemas (one per table).
	db.AutoMigrate(&Symbol{})
	db.AutoMigrate(&Account{})
	db.AutoMigrate(&Position{})

	// Is this a testing run? If so load testing data.
	if strings.HasSuffix(os.Args[0], ".test") {
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
