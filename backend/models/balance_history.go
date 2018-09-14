//
// Date: 9/14/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"time"
)

type BalanceHistory struct {
	Id                uint      `gorm:"primary_key" json:"id"`
	CreatedAt         time.Time `json:"-"`
	UpdatedAt         time.Time `json:"-"`
	UserId            uint      `gorm:"index" sql:"not null;index:UserId" json:"-"`
	BrokerAccountId   uint      `gorm:"index" sql:"not null;index:BrokerId" json:"broker_id"`
	Date              Date      `gorm:"type:date" json:"date"`
	AccountNumber     string    `json:"account_number"`
	AccountValue      float64   `json:"account_value"`
	TotalCash         float64   `json:"total_cash"`
	OptionBuyingPower float64   `json:"option_buying_power"`
	StockBuyingPower  float64   `json:"stock_buying_power"`
}

/* End File */
