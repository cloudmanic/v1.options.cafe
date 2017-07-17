package users

import (
  "app.options.cafe/backend/models"
  "app.options.cafe/backend/controllers"
  "app.options.cafe/backend/brokers"
  "app.options.cafe/backend/brokers/feed"
  "app.options.cafe/backend/brokers/tradier"  
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
    
    // Figure out which broker connection to setup.
    switch row.Name {
      
      case "Tradier":
        brokerApi = &tradier.Api{ ApiKey: row.AccessToken }
        
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

/* End File */