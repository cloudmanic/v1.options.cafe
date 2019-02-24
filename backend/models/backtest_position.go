//
// Date: 2/22/2019
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"time"
)

// BacktestPosition struct
type BacktestPosition struct {
	Id         uint      `gorm:"primary_key" json:"id"`
	CreatedAt  time.Time `sql:"not null" json:"-"`
	UpdatedAt  time.Time `sql:"not null" json:"-"`
	UserId     uint      `sql:"not null" index:UserId" json:"user_id"`
	BacktestId uint      `sql:"not null" index:BacktestId" json:"backtest_id"`
	Status     string    `sql:"not null;type:ENUM('Open', 'Closed');default:'Open'" json:"status"`
	OpenDate   Date      `gorm:"type:date" sql:"not null" json:"open_date"`
	CloseDate  Date      `gorm:"type:date" sql:"not null" json:"close_date"`
	OpenPrice  float64   `sql:"not null" json:"open_price"`
	ClosePrice float64   `sql:"not null" json:"close_price"`
	Margin     float64   `sql:"not null" json:"margin"`
	Balance    float64   `sql:"not null" json:"balance"`
	Lots       int       `sql:"not null" json:"lots"`
	Legs       []Symbol  `sql:"not null" json:"legs" gorm:"many2many:backtest_positions_symbols;"`
	Note       string    `sql:"not null" json:"note"`
}

/* End File */
