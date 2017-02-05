package models

import (
  "time"
  "github.com/jinzhu/gorm"
)

type User struct {
  Id uint `gorm:"primary_key"`
  CreatedAt time.Time
  UpdatedAt time.Time
  FirstName string `sql:"not null"`
  LastName string `sql:"not null"`
  Email string `sql:"not null"`
  Password string `sql:"not null"`
  AccessToken string `sql:"not null"`
  Status string `sql:"not null;type:ENUM('Active', 'Disable');default:'Active'"`
  Brokers []Broker
}     
      
// 
// Return an array of all users.
//
func (s *User) GetAllUsers(db *gorm.DB) ([]User) {
  
  var users []User
  
  db.Find(&users)
  
  // Add in our one to many look ups
  for i, _ := range users {
    db.Model(users[i]).Related(&users[i].Brokers)     
  }  
  
  return users
  
}
      
/* End File */