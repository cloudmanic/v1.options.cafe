package models

import (
	"errors"
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/library/helpers"
	"github.com/cloudmanic/app.options.cafe/backend/library/services"
)

type Broker struct {
	Id                  uint            `gorm:"primary_key" json:"id"`
	CreatedAt           time.Time       `json:"-"`
	UpdatedAt           time.Time       `json:"-"`
	UserId              uint            `gorm:"index" sql:"not null;index:UserId" json:"-"`
	Name                string          `sql:"not null;type:ENUM('Tradier', 'Tradier Sandbox', 'Tradeking', 'Etrade', 'Interactive Brokers'); default:'Tradier'" json:"name"`
	AccessToken         string          `sql:"not null" json:"-"`
	RefreshToken        string          `sql:"not null" json:"-"`
	TokenExpirationDate time.Time       `json:"-"`
	BrokerAccounts      []BrokerAccount `json:"broker_accounts"`
}

//
// Get a broker by Id.
//
func (t *DB) GetBrokerById(id uint) (Broker, error) {

	var u Broker

	if t.Where("Id = ?", id).First(&u).RecordNotFound() {
		return u, errors.New("Record not found")
	}

	// Return the user.
	return u, nil

}

//
// Get a brokers by type and user id.
//
func (t *DB) GetBrokerTypeAndUserId(userId uint, brokerType string) ([]Broker, error) {

	var u []Broker

	if t.Where("user_id = ? AND name = ?", userId, brokerType).Find(&u).RecordNotFound() {
		return u, errors.New("Records not found")
	}

	// Return the user.
	return u, nil

}

//
// Create a new broker entry.
//
func (t *DB) CreateNewBroker(name string, user User, accessToken string, refreshToken string, tokenExpirationDate time.Time) (Broker, error) {

	// Encrypt the access token
	encryptAccessToken, err := helpers.Encrypt(accessToken)

	if err != nil {
		services.Error(err, "[Models:CreateNewBroker] Unable to encrypt message (#1)")
		return Broker{}, errors.New("[CreateNewBroker] Unable to encrypt message (#1)")
	}

	// Encrypt the refresh token
	encryptRefreshToken, err := helpers.Encrypt(refreshToken)

	if err != nil {
		services.Error(err, "[Models:CreateNewBroker] Unable to encrypt message (#2)")
		return Broker{}, errors.New("[Models:CreateNewBroker] Unable to encrypt message (#2)")
	}

	// Create entry.
	broker := Broker{
		Name:                name,
		UserId:              user.Id,
		AccessToken:         encryptAccessToken,
		RefreshToken:        encryptRefreshToken,
		TokenExpirationDate: tokenExpirationDate,
	}

	t.Create(&broker)

	// Log broker creation.
	services.Info("[Models:CreateNewBroker] - Created a new broker entry - " + name + " " + user.Email)

	// Return the user.
	return broker, nil

}

//
// Update a new broker entry.
//
func (t *DB) UpdateBroker(broker Broker) error {

	// Encrypt the access token
	encryptAccessToken, err := helpers.Encrypt(broker.AccessToken)

	if err != nil {
		services.Error(err, "(UpdateBroker) Unable to encrypt message (#1)")
		return errors.New("(UpdateBroker) Unable to encrypt message (#1)")
	}

	broker.AccessToken = encryptAccessToken

	// Encrypt the refresh token
	encryptRefreshToken, err := helpers.Encrypt(broker.RefreshToken)

	if err != nil {
		services.Error(err, "(UpdateBroker) Unable to encrypt message (#2)")
		return errors.New("(UpdateBroker) Unable to encrypt message (#2)")
	}

	broker.RefreshToken = encryptRefreshToken

	// Update entry.
	t.Save(&broker)

	// Return the user.
	return nil

}

/* End File */
