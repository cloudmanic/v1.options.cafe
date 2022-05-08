//
// Date: 5/2/2022
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2022 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"time"
)

// BacktestPosition struct
type BacktestPosition struct {
	Id                   uint      `gorm:"primary_key" json:"id"`
	UserId               uint      `sql:"not null;index:UserId" json:"_"`
	CreatedAt            time.Time `json:"_"`
	UpdatedAt            time.Time `json:"_"`
	BacktestTradeGroupId uint      `sql:"not null;index:BacktestTradeGroupId" json:"backtest_trade_group_id"`
	Status               string    `sql:"not null;type:ENUM('Open', 'Closed');default:'Open'" json:"status"`
	SymbolId             uint      `json:"_"`
	Symbol               Symbol    `json:"symbol"`
	Qty                  int       `json:"qty"`
	OrgQty               int       `json:"org_qty"`
	CostBasis            float64   `sql:"type:DECIMAL(12,2)" json:"cost_basis"`
	Proceeds             float64   `sql:"type:DECIMAL(12,2)" json:"proceeds"`
	Profit               float64   `sql:"type:DECIMAL(12,2)" json:"profit"`
	AvgOpenPrice         float64   `sql:"type:DECIMAL(12,2)" json:"avg_open_price"`
	AvgClosePrice        float64   `sql:"type:DECIMAL(12,2)" json:"avg_close_price"`
	Note                 string    `sql:"type:text" json:"note"`
	OpenDate             time.Time `json:"open_date"`
	ClosedDate           time.Time `json:"close_date"`
}

/* End File */
