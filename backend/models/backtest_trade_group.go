//
// Date: 2/22/2019
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"time"
)

// BacktestTradeGroup struct
type BacktestTradeGroup struct {
	Id               uint               `gorm:"primary_key" json:"id"`
	CreatedAt        time.Time          `sql:"not null" json:"-"`
	UpdatedAt        time.Time          `sql:"not null" json:"-"`
	UserId           uint               `sql:"not null;index:UserId" json:"user_id"`
	BacktestId       uint               `sql:"not null;index:BacktestId" json:"backtest_id"`
	Strategy         string             `sql:"not null" json:"strategy"`
	Status           string             `sql:"not null;type:ENUM('Open', 'Closed');default:'Open'" json:"status"`
	OpenDate         Date               `gorm:"type:date" sql:"not null" json:"open_date"`
	CloseDate        Date               `gorm:"type:date" sql:"not null" json:"close_date"`
	OpenPrice        float64            `sql:"not null" json:"open_price"`
	ClosePrice       float64            `sql:"not null" json:"close_price"`
	Credit           float64            `sql:"not null" json:"credit"`
	ReturnPercent    float64            `sql:"not null" json:"return_percent"`
	ReturnFromStart  float64            `sql:"not null" json:"return_from_start"`
	Margin           float64            `sql:"not null" json:"margin"`
	Balance          float64            `sql:"not null" json:"balance"`
	PutPrecentAway   float64            `sql:"not null" json:"put_percent_away"`
	CallPrecentAway  float64            `sql:"not null" json:"call_percent_away"`
	BenchmarkLast    float64            `sql:"not null" json:"benchmark_last"`
	BenchmarkBalance float64            `sql:"not null" json:"benchmark_balance"`
	BenchmarkReturn  float64            `sql:"not null" json:"benchmark_return"`
	Lots             int                `sql:"not null" json:"lots"`
	Positions        []BacktestPosition `json:"positions"`
	Note             string             `sql:"not null" json:"note"`
}

/* End File */
