//
// Date: 9/30/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"time"
)

type Symbol struct {
	Id               uint      `gorm:"primary_key" json:"id"`
	CreatedAt        time.Time `json:"-"`
	UpdatedAt        time.Time `json:"-"`
	ShortName        string    `sql:"not null" json:"short_name"`
	Name             string    `sql:"not null" json:"name"`
	Type             string    `sql:"not null;type:ENUM('Equity', 'Option', 'Other');default:'Equity'" json:"type"`
	OptionUnderlying string    `json:"option_underlying"`
	OptionType       string    `json:"option_type"`
	OptionExpire     Date      `gorm:"type:date" json:"option_expire"`
	OptionStrike     float64   `json:"option_strike"`
}

/* End File */
