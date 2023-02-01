package models

import (
	"time"
)

type Gain struct {
	Id              uint      `sql:"not null" gorm:"primary_key"`
	CreatedAt       time.Time `sql:"not null"`
	UpdatedAt       time.Time `sql:"not null"`
	UserId          uint      `sql:"not null;index:UserId"`
	TradeId         uint      `sql:"not null;index:TradeId" json:"trade_id"`
	BrokerAccountId uint      `sql:"not null;index:BrokerAccountId"`
	SymbolId        uint      `sql:"not null"`
	Symbol          Symbol    `sql:"not null" json:"symbol"`
	OpenDate        Date      `gorm:"type:date" sql:"not null" json:"open_date"`
	CloseDate       Date      `gorm:"type:date" sql:"not null" json:"close_date"`
	Cost            float64   `sql:"not null" json:"cost"`
	GainLoss        float64   `sql:"not null" json:"gain_loss"`
	GainLossPercent float64   `sql:"not null" json:"gain_loss_percent"`
	Proceeds        float64   `sql:"not null" json:"proceeds"`
	Quantity        float64   `sql:"not null" json:"quantity"`
	Term            int64     `sql:"not null" json:"term"`
	Reviewed        string    `sql:"not null;type:ENUM('No', 'Yes', 'Error');default:'No'"`
}
