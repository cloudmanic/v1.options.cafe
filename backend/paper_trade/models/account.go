//
// Date: 2018-11-05
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-11-05
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import "time"

type Account struct {
	Id            uint      `gorm:"primary_key" json:"id"`
	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
	FirstName     string    `sql:"not null" json:"first_name"`
	LastName      string    `sql:"not null" json:"last_name"`
	Email         string    `sql:"not null" json:"email"`
	AccountNumber string    `sql:"not null" json:"account_number"`
	AccessToken   string    `sql:"not null" json:"-"`
}

/* End File */
