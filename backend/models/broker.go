package models

import (
  "time"
)

type Broker struct {
  Id uint `gorm:"primary_key"`
  CreatedAt time.Time `sql:"not null"`
  UpdatedAt time.Time `sql:"not null"`
  UserId uint `sql:"not null;index:UserId"` 
  Name string `sql:"not null;type:ENUM('Tradier', 'Tradeking', 'Etrade', 'Interactive Brokers'); default:'Tradier'"`
  AccessToken string `sql:"not null"`
  RefreshToken string `sql:"not null"`  
}     
      
/* End File */