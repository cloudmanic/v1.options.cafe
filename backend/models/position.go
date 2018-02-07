package models

import (
	"errors"
	"time"
)

type Position struct {
	Id            uint `gorm:"primary_key"`
	UserId        uint `sql:"not null;index:UserId"`
	TradeGroupId  uint `sql:"not null;index:TradeGroupId"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	AccountId     string `sql:"not null;index:AccountId"`
	Status        string `sql:"not null;type:ENUM('Open', 'Closed');default:'Open'"`
	Symbol        string
	Qty           int
	OrgQty        int
	CostBasis     float64 `sql:"type:DECIMAL(12,2)"`
	AvgOpenPrice  float64 `sql:"type:DECIMAL(12,2)"`
	AvgClosePrice float64 `sql:"type:DECIMAL(12,2)"`
	OrderIds      string
	Note          string `sql:"type:text"`
	OpenDate      time.Time
	ClosedDate    time.Time
}

//
// Store a new position.
//
func (t *DB) CreatePosition(position *Position) error {

	// Create position
	t.Create(position)

	// Return happy
	return nil
}

//
// Get positions by User and class and status and reviewed
//
func (t *DB) GetPositionByUserSymbolStatusAccount(userId uint, symbol string, status string, accountId string) (Position, error) {

	var position = Position{}

	// First we find out if we already have a position on for this.
	if t.Where("symbol = ? AND user_id = ? AND status = ? AND account_id = ?", symbol, userId, "Open", accountId).First(&position).RecordNotFound() {
		return position, errors.New("Record not found")
	}

	// Return happy
	return position, nil
}

/* End File */
