//
// Date: 9/30/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package users

import (
  "app.options.cafe/backend/models"
  "app.options.cafe/backend/controllers"
  "app.options.cafe/backend/brokers"
  "app.options.cafe/backend/brokers/feed"
  "app.options.cafe/backend/brokers/tradier"
  "app.options.cafe/backend/library/helpers"  
  "app.options.cafe/backend/library/services"
)

var (
  DB *models.DB
  Users map[uint]*User
  DataChan chan controllers.SendStruct
  QuoteChan chan controllers.SendStruct
  FeedRequestChan chan controllers.SendStruct   
)
 
type User struct {
  Profile models.User
  BrokerFeed map[uint]*feed.Base
}

//
// Start up our user feeds.
//
func StartFeeds() {
  
  // Setup the map of users.
  Users = make(map[uint]*User)
  
  // Get all active users
  users := DB.GetAllActiveUsers()

  // Loop through the users
  for i, _ := range users {
    DoUserFeed(users[i])     
  }
  
  // Listen of income Feed Requests.
  go DoFeedRequestListen()  
  
} 

//
// Start one user.
//
func DoUserFeed(user models.User) {
 
  var brokerApi brokers.Api   
 
  services.Log("Starting User Connection : " + user.Email)
  
  // This should not happen. But we double check this user is not already started.
  if _, ok := Users[user.Id]; ok {
    services.MajorLog("User Connection Is Already Going : " + user.Email)
    return
  }
  
  // Verify some default data.
  VerifyDefaultWatchList(user)
  
  // Set the user to the object
  Users[user.Id] = &User{
                      Profile: user,
                      BrokerFeed: make(map[uint]*feed.Base),   
                    }
    
  // Loop through the different brokers for this user
  for _, row := range Users[user.Id].Profile.Brokers {
    
    // Need an access token to continue
    if len(row.AccessToken) <= 0 {
      services.MajorLog("User Connection (Brokers) No Access Token Found : " + user.Email + " (" + row.Name + ")")
      continue
    }    
    
    // Decrypt the access token
    decryptAccessToken, err := helpers.Decrypt(row.AccessToken)
    
    if err != nil {
      services.Error(err, "(DoUserFeed) Unable to decrypt message (#1)")
      continue
    }     
    
    // Figure out which broker connection to setup.
    switch row.Name {
      
      case "Tradier":
        brokerApi = &tradier.Api{ ApiKey: decryptAccessToken }
        
      default:
        services.MajorLog("Unknown Broker : " + row.Name + " (" + user.Email + ")")
        continue
        
    }
    
    // Log magic
    services.Log("Setting up to use " + row.Name + " as the broker for " + user.Email)  
    
    // Set the library we use to fetching data from our broker's API
    Users[user.Id].BrokerFeed[row.Id] = &feed.Base{ 
                                          User: user, 
                                          Api: brokerApi,
                                          DataChan: DataChan,
                                          QuoteChan: QuoteChan,                       
                                        }
                      
    // Start fetching data for this user.
    go Users[user.Id].BrokerFeed[row.Id].Start()    
  }
  
}

//
// Listen for incomeing feed requests.
//
func DoFeedRequestListen() {

  for {
  
    send := <-FeedRequestChan
    
    switch send.Message {
      
      // Refresh all data from cache - FromCache:refresh
      case "FromCache:refresh":
        
        // Loop through each broker and refresh the data.
        for _, row := range Users[send.UserId].BrokerFeed {
          
          row.RefreshFromCached()
          
        }
        
    }

  } 
  
}

// ------------- Helper Functions ------------------ //

//
// Verify we have default watchlist in place.
//
func VerifyDefaultWatchList(user models.User) {
  
  // Setup defaults.
  type Y struct {
    SymShort string
    SymLong string
  }
  
  var m []Y
  m = append(m, Y{ SymShort: "SPY", SymLong: "SPDR S&P 500" })
  m = append(m, Y{ SymShort: "IWM", SymLong: "Ishares Russell 2000 Etf" })
  m = append(m, Y{ SymShort: "VIX", SymLong: "CBOE Volatility S&P 500 Index" })
  m = append(m, Y{ SymShort: "AMZN", SymLong: "Amazon.com Inc" })
  m = append(m, Y{ SymShort: "AAPL", SymLong: "Apple Inc." })      
  m = append(m, Y{ SymShort: "SBUX", SymLong: "Starbucks Corp" })
  m = append(m, Y{ SymShort: "BAC", SymLong: "Bank Of America Corporation" })

  // See if this user already had a watchlist
  _, err := DB.GetWatchlistsByUserId(user.Id)
  
  // If no watchlists we create a default one with some default symbols.  
  if err != nil {

    wList, err := DB.CreateNewWatchlist(user, "Default")

    if err != nil {
      services.Error(err, "(CreateNewWatchlist) Unable to create watchlist Default")
      return
    }

    for key, row := range m {

      // Add some default symbols - SPY
      symb, err := DB.CreateNewSymbol(row.SymShort, row.SymLong)
      
      if err != nil {
        services.Error(err, "(VerifyDefaultWatchList) Unable to create symbol " + row.SymShort)
        return
      }
      
      // Add lookup
      _, err2 := DB.CreateNewWatchlistSymbol(wList, symb, user, uint(key))      
  
      if err2 != nil {
        services.Error(err2, "(CreateNewWatchlistSymbol) Unable to create symbol " + row.SymShort + " lookup")
        return
      }
    
    }
    
  }
  
  return
  
}

/* End File */