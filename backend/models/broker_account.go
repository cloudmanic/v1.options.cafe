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
	validation "github.com/go-ozzo/ozzo-validation"
)

type BrokerAccount struct {
	Id                uint      `gorm:"primary_key" json:"id"`
	CreatedAt         time.Time `json:"-"`
	UpdatedAt         time.Time `json:"-"`
	UserId            uint      `gorm:"index" sql:"not null;index:UserId" json:"-"`
	BrokerId          uint      `gorm:"index" sql:"not null;index:BrokerId" json:"broker_id"`
	Name              string    `sql:"not null" json:"name"`
	AccountNumber     string    `sql:"not null" json:"account_number"`
	StockCommission   float64   `sql:"type:DECIMAL(12,2)" json:"stock_commission"`
	StockMin          float64   `sql:"type:DECIMAL(12,2)" json:"stock_min"`
	OptionCommission  float64   `sql:"type:DECIMAL(12,2)" json:"option_commission"`
	OptionSingleMin   float64   `sql:"type:DECIMAL(12,2)" json:"option_single_min"`
	OptionMultiLegMin float64   `sql:"type:DECIMAL(12,2)" json:"option_multi_leg_min"`
	OptionBase        float64   `sql:"type:DECIMAL(12,2)" json:"option_base"`
}

//
// Validate for this model.
//
func (a BrokerAccount) Validate(db Datastore) error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Name, validation.Required.Error("The name field is required.")),
		// validation.Field(&a.Strategy, validation.Required.Error("The strategy field is required.")),
		// validation.Field(&a.Symbol, validation.Required.Error("The symbol field is required.")),
	)
}

//
// Get the broker account by id and user id
//
func (t *DB) GetBrokerAccountByIdUserId(id uint, userId uint) (BrokerAccount, error) {

	ba := BrokerAccount{}

	// Query and get broker account
	t.Where("id = ? AND user_id = ?", id, userId).First(&ba)

	// Return happy
	return ba, nil
}

//
// Get broker from broker account.
//
func (t *DB) GetBrokerFromBrokerAccountAndUserId(id uint, userId uint) (Broker, error) {

	// Get broker account.
	brokerAccount, err := t.GetBrokerAccountByIdUserId(id, userId)

	if err != nil {
		return Broker{}, err
	}

	// Get the broker
	broker, err := t.GetBrokerById(brokerAccount.BrokerId)

	if err != nil {
		return Broker{}, err
	}

	return broker, nil
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
// Get the broker account by broker and account number
//
func (t *DB) GetBrokerAccountByBrokerAccountNumber(brokerId uint, accountNumber string) (BrokerAccount, error) {

	ba := BrokerAccount{}

	// Query and get broker account
	t.Where("broker_id = ? AND account_number = ?", brokerId, accountNumber).First(&ba)

	// Return happy
	return ba, nil
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
