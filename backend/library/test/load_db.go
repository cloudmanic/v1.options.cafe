//
// Date: 3/6/2019
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

// Package test used to load data data into our database.
package test

import (
	"fmt"
	"go/build"
	"os/exec"
	"strings"

	"github.com/cloudmanic/app.options.cafe/backend/models"
)

//
// LoadSqlDump - Take a dataset string and load content from
// a data file into our database. These dataset strings are mysql dumps
//
func LoadSqlDump(db models.Datastore, dataSet string) error {
	// Get DB name - yes, very hacky
	s := fmt.Sprintf("%s", db.New().CommonDB())
	s1 := strings.Split(s, "(127.0.0.1:9906)/")
	s2 := strings.Split(s1[1], "?charset=")
	dbName := s2[0]

	// Set file
	dataFile := build.Default.GOPATH + "/src/github.com/cloudmanic/app.options.cafe/backend/library/test/data/" + dataSet + ".sql"

	// Run CLI command to import data
	_, err := exec.Command("bash", "-c", "mysql --host=127.0.0.1 --port=9906 -u root -pfoobar "+dbName+" < "+dataFile).Output()

	if err != nil {
		return err
	}

	// Return happy
	return nil
}

/* End File */
