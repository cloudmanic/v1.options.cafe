//
// Date: 9/30/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
  "time"
  "errors"
  "app.options.cafe/backend/library/services"
)

type Watchlist struct {
  Id uint `gorm:"primary_key"`
  CreatedAt time.Time
  UpdatedAt time.Time
  UserId uint `sql:"not null;index:UserId"` 
  Name string `sql:"not null"`
} 

//
// Get a Watchlists by user id.
//
func (t * DB) GetWatchlistsByUserId(userId uint) ([]Watchlist, error) {
 
  var u []Watchlist
  
  if t.Connection.Where("user_id = ?", userId).Find(&u).RecordNotFound() {
    return u, errors.New("Records not found")
  }
  
  if len(u) <= 0 {
    return u, errors.New("Records not found")    
  }
  
  // Return the Watchlists.
  return u, nil
  
} 

//
// Create a new Watchlist entry.
//
func (t * DB) CreateNewWatchlist(user User, name string) (Watchlist, error) {
    
  // Create entry.
  wList := Watchlist{ Name: name, UserId: user.Id }
            
  t.Connection.Create(&wList)
  
  // Log broker creation.
  services.Log("CreateNewWatchlist - Created a new Watchlist entry - " + name + " " + user.Email)  
  
  // Return the user.
  return wList, nil  
   
} 
          
/* End File */