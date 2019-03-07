//
// Date: 3/6/2019
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

// Package test used to load data data into our database.
package test

import (
	"go/build"
	"io/ioutil"

	"github.com/cloudmanic/app.options.cafe/backend/models"
)

//
// LoadSqlDump - Take a dataset string and load content from
// a data file into our database. These dataset strings are mysql dumps
//
func LoadSqlDump(db models.Datastore, dataSet string) error {
	// Set file
	dataFile := build.Default.GOPATH + "/src/github.com/cloudmanic/app.options.cafe/backend/library/test/data/" + dataSet + ".sql"

	// Read JSON file
	dat, err := ioutil.ReadFile(dataFile)

	if err != nil {
		return err
	}

	// Run SQL dump we got from file.
	db.New().Exec(string(dat))

	// Return happy
	return nil
}

/* End File */
