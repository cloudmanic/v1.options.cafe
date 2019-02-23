//
// Date: 2/22/2019
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"time"
)

// Backtest struct
type Backtest struct {
	Id              uint               `gorm:"primary_key" json:"id"`
	CreatedAt       time.Time          `json:"-"`
	UpdatedAt       time.Time          `json:"-"`
	UserId          uint               `sql:"not null;index:UserId" json:"user_id"`
	StartDate       time.Time          `sql:"not null" json:"start_date"`
	EndDate         time.Time          `sql:"not null" json:"end_date"`
	EndingBalance   float64            `sql:"not null" json:"ending_balance"`
	StartingBalance float64            `sql:"not null" json:"starting_balance"`
	TradeSelect     string             `sql:"not null;type:ENUM('highest-credit', 'median-credit', 'lowest-credit');default:'median-credit'" json:"trade_select"`
	Midpoint        bool               `sql:"not null" json:"midpoint"` // Open trade at the midpoint
	Screen          Screener           `json:"screen"`
	Positions       []BacktestPosition `json:"positions"`
}

/* End File */
