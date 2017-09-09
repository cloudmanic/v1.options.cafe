package models

import (
  "time"
  "errors"
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
  
  // Create entry.
  broker := Broker{ 
              Name: name, 
              UserId: user.Id,
              AccessToken: accessToken,
              RefreshToken: refreshToken,
              TokenExpirationDate: tokenExpirationDate,
            }
            
  t.Connection.Create(&broker)
  
  // Log broker creation.
  services.Log("CreateNewBroker - Created a new broker entry - " + name + " " + user.Email)  
  
  // Return the user.
  return broker, nil  
   
}      
      
/* End File */