package models

import (
  "time"
)

type User struct {
  Id uint `gorm:"primary_key"`
  CreatedAt time.Time
  UpdatedAt time.Time
  FirstName string
  LastName string
  Email string
  Password string
  Token string
}     
      
/* End File */