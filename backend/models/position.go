//
// Date: 2/9/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"errors"
	"time"
)

type Position struct {
	Id            uint      `gorm:"primary_key" json:"id"`
	UserId        uint      `sql:"not null;index:UserId" json:"_"`
	CreatedAt     time.Time `json:"_"`
	UpdatedAt     time.Time `json:"_"`
	TradeGroupId  uint      `sql:"not null;index:TradeGroupId" json:"trade_group_id"`
	AccountId     string    `sql:"not null;index:AccountId" json:"account_id"`
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

//
// Store a new position.
//
func (t *DB) CreatePosition(position *Position) error {

	// Create position
	t.Create(position)

	// Add in symbol
	t.Model(position).Related(&position.Symbol)

	// Return happy
	return nil
}

//
// Update a position.
//
func (t *DB) UpdatePosition(position *Position) error {

	// Update entry.
	t.Save(&position)

	// Add in symbol
	t.Model(position).Related(&position.Symbol)

	// Return happy
	return nil
}

//
// Get positions by User and class and status and reviewed
//
func (t *DB) GetPositionByUserSymbolStatusAccount(userId uint, symbolId uint, status string, accountId string) (Position, error) {

	var position = Position{}

	// First we find out if we already have a position on for this.
	if t.Where("symbol_id = ? AND user_id = ? AND status = ? AND account_id = ?", symbolId, userId, status, accountId).First(&position).RecordNotFound() {
		return position, errors.New("Record not found")
	}

	// Return happy
	return position, nil
}

/* End File */
