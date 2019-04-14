//
// Date: 4/13/2019
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"time"
)

// HistoricalQuote struct - Main use of this struct is to make unit testing better.
// So we do not need to call broker APIs
type HistoricalQuote struct {
	Id        uint      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	ShortName string    `sql:"not null" json:"short_name"`
	Date      Date      `gorm:"type:date" json:"date"`
	Price     float64   `json:"price"`
}

/* End File */
