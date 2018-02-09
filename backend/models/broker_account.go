//
// Date: 2/9/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"strconv"
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/library/services"
)

type BrokerAccount struct {
	Id                uint      `gorm:"primary_key"`
	CreatedAt         time.Time `json:"-"`
	UpdatedAt         time.Time `json:"-"`
	UserId            uint      `gorm:"index" sql:"not null;index:UserId" json:"-"`
	BrokerId          uint      `gorm:"index" sql:"not null;index:BrokerId"`
	Name              string    `sql:"not null"`
	AccountNumber     string    `sql:"not null"`
	StockCommission   float64   `sql:"type:DECIMAL(12,2)"`
	StockMin          float64   `sql:"type:DECIMAL(12,2)"`
	OptionCommission  float64   `sql:"type:DECIMAL(12,2)"`
	OptionSingleMin   float64   `sql:"type:DECIMAL(12,2)"`
	OptionMultiLegMin float64   `sql:"type:DECIMAL(12,2)"`
	OptionBase        float64   `sql:"type:DECIMAL(12,2)"`
}

//
// Update the broker account object.
//
func (t *DB) UpdateBrokerAccount(brokerAccount *BrokerAccount) error {

	// Update entry.
	t.Save(&brokerAccount)

	// Return happy
	return nil
}

//
// Look for a broker account. If we can't find it create it.
// The bool return will be true if this was a new record created.
//
func (t *DB) FirstOrCreateBrokerAccount(brokerAccount *BrokerAccount) (bool, error) {

	// First lets see if this record is already in the DB.
	t.FirstOrInit(brokerAccount, brokerAccount)

	// Ok we did not find the record lets create it.
	if brokerAccount.Id == 0 {
		t.Create(brokerAccount)

		// Log user creation.
		services.Info("FirstOrCreateBrokerAccount - Created a new broker account - " + strconv.Itoa(int(brokerAccount.Id)))
	} else {
		return false, nil
	}

	// Return happy.
	return true, nil
}

/* End File */
