//
// Date: 7/18/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"time"
)

type ScreenerItem struct {
	Id          uint      `gorm:"primary_key" json:"id"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
	ScreenerId  uint      `sql:"not null;index:ScreenerId" json:"screener_id"`
	UserId      uint      `sql:"not null;index:UserId" json:"user_id"`
	Key         string    `sql:"not null" json:"key"`
	Operator    string    `json:"operator"`
	ValueString string    `json:"value_string"`
	ValueNumber float64   `json:"value_number"`
}

/* End File */
