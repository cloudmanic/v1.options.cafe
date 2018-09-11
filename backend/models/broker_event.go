//
// Date: 2018-09-11
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-09-11
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"time"
)

type BrokerEvent struct {
	Id              uint `gorm:"primary_key"`
	UserId          uint `sql:"not null;index:UserId" json:"_"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	BrokerAccountId uint      `sql:"not null;index:BrokerAccountId" json:"_"`
	BrokerId        string    `json:"_"`
	Type            string    `sql:"not null;type:ENUM('Ach', 'Trade', 'Option', 'Interest', 'Journal', 'Dividend', 'Adjustment', 'Other');default:'Other'" json:"type"`
	Date            time.Time `json:"date"`
	Amount          float64   `json:"amount"`
	Symbol          string    `json:"symbol"`
	Commission      float64   `json:"commission"`
	Description     string    `json:"description"`
	Price           float64   `json:"price"`
	Quantity        int64     `json:"quantity"`
	TradeType       string    `sql:"not null;type:ENUM('Equity', 'Option', 'ETF', 'Preferred Stock', 'Other');default:'Other'" json:"trade_type"`
}

/* End File */
