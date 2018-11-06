//
// Date: 2/9/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"time"
)

type Position struct {
	Id            uint      `gorm:"primary_key" json:"id"`
	AccountId     uint      `sql:"not null;index:AccountId" json:"_"`
	CreatedAt     time.Time `json:"_"`
	UpdatedAt     time.Time `json:"_"`
	Status        string    `sql:"not null;type:ENUM('Open', 'Closed');default:'Open'" json:"status"`
	SymbolId      uint      `json:"_"`
	Qty           int       `json:"qty"`
	OrgQty        int       `json:"org_qty"`
	CostBasis     float64   `sql:"type:DECIMAL(12,2)" json:"cost_basis"`
	AvgOpenPrice  float64   `sql:"type:DECIMAL(12,2)" json:"avg_open_price"`
	AvgClosePrice float64   `sql:"type:DECIMAL(12,2)" json:"avg_close_price"`
	OrderIds      string    `json:"_"`
	Note          string    `sql:"type:text" json:"note"`
	OpenDate      time.Time `json:"open_date"`
	ClosedDate    time.Time `json:"close_date"`
	Symbol        Symbol    `json:"symbol"`
}

/* End File */
