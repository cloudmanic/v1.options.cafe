package models

import (
  "time"
  "errors"
  "app.options.cafe/backend/library/helpers"
  "app.options.cafe/backend/library/services"
)

type Broker struct {
  Id uint `gorm:"primary_key"`
  CreatedAt time.Time
  UpdatedAt time.Time
  UserId uint `sql:"not null;index:UserId"` 
  Name string `sql:"not null;type:ENUM('Tradier', 'Tradeking', 'Etrade', 'Interactive Brokers'); default:'Tradier'"`
  AccessToken string `sql:"not null"`
  RefreshToken string `sql:"not null"`
  TokenExpirationDate time.Time 
}      

//
// Get a broker by Id.
//
func (t * DB) GetBrokerById(id uint) (Broker, error) {
 
  var u Broker
  
  if t.Connection.Where("Id = ?", id).First(&u).RecordNotFound() {
    return u, errors.New("Record not found")
  }
  
  // Return the user.
  return u, nil
  
} 

//
// Get a brokers by type and user id.
//
func (t * DB) GetBrokerTypeAndUserId(userId uint, brokerType string) ([]Broker, error) {
 
  var u []Broker
  
  if t.Connection.Where("user_id = ? AND name = ?", userId, brokerType).Find(&u).RecordNotFound() {
    return u, errors.New("Records not found")
  }
  
  // Return the user.
  return u, nil
  
}

//
// Create a new broker entry.
//
func (t * DB) CreateNewBroker(name string, user User, accessToken string, refreshToken string, tokenExpirationDate time.Time) (Broker, error) {
  
  // Encrypt the access token
	encryptAccessToken, err := helpers.Encrypt(accessToken)
	
  if err != nil {
    services.Error(err, "(CreateNewBroker) Unable to encrypt message (#1)")
    return Broker{}, errors.New("(CreateNewBroker) Unable to encrypt message (#1)")
  }  

  // Encrypt the refresh token
	encryptRefreshToken, err := helpers.Encrypt(refreshToken)
	
  if err != nil {
    services.Error(err, "(CreateNewBroker) Unable to encrypt message (#2)")
    return Broker{}, errors.New("(CreateNewBroker) Unable to encrypt message (#2)")
  }
  
  // Create entry.
  broker := Broker{ 
              Name: name, 
              UserId: user.Id,
              AccessToken: encryptAccessToken,
              RefreshToken: encryptRefreshToken,
              TokenExpirationDate: tokenExpirationDate,
            }
            
  t.Connection.Create(&broker)
  
  // Log broker creation.
  services.Log("CreateNewBroker - Created a new broker entry - " + name + " " + user.Email)  
  
  // Return the user.
  return broker, nil  
   
}

//
// Update a new broker entry.
//
func (t * DB) UpdateBroker(broker Broker) error {
  
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
  t.Connection.Save(&broker)
   
  // Return the user.
  return nil  
   
}      
      
/* End File */