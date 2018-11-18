//
// Date: 2018-11-05
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-11-18
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type DB struct {
	*gorm.DB
}

// Database interface
type Datastore interface {

	// Gorm Functions
	New() *gorm.DB
}

/* End File */
