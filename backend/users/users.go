package users

import (
  "app.options.cafe/backend/models"
  "app.options.cafe/backend/websocket"
  "app.options.cafe/backend/brokers"
  "app.options.cafe/backend/brokers/feed"
  "app.options.cafe/backend/brokers/tradier"  
  "app.options.cafe/backend/library/services"
)

type Base struct {
  DB *models.DB
  Users map[uint]*User
  WsWriteChannel chan websocket.SendStruct
  WsWriteQuoteChannel chan websocket.SendStruct
}

type User struct {
  Profile models.User
  BrokerFeed map[uint]*feed.Base
  DataChannel chan string
  QuoteChannel chan string
}

//
// Start up our users.
//
func (t * Base) Start() {
  
  // Setup the map of users.
  t.Users = make(map[uint]*User)
  
  // Get all active users
  users := t.DB.GetAllActiveUsers()

  // Loop through the users
  for i, _ := range users {
    t.doUser(users[i])     
  }  
  
} 

//
// Start one user.
//
func (t * Base) doUser(user models.User) {
 
  var brokerApi brokers.Api   
 
  services.Log("Starting User Connection : " + user.Email)
  
  // This should not happen. But we double check this user is not already started.
  if _, ok := t.Users[user.Id]; ok {
    services.MajorLog("User Connection Is Already Going : " + user.Email)
    return
  }
  
  // Set the user to the object
  t.Users[user.Id] = &User{
                        Profile: user,
                        BrokerFeed: make(map[uint]*feed.Base),
                        DataChannel: make(chan string, 1000),
                        QuoteChannel: make(chan string, 1000),    
                      }
    
  // Loop through the different brokers for this user
  for _, row := range t.Users[user.Id].Profile.Brokers {
    
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
    t.Users[user.Id].BrokerFeed[row.Id] = &feed.Base{ 
                                            User: user, 
                                            Api: brokerApi,
                                            DataChannel: t.Users[user.Id].DataChannel,
                                            QuoteChannel: t.Users[user.Id].QuoteChannel,                       
                                          }
                      
    // Start fetching data for this user.
    go t.Users[user.Id].BrokerFeed[row.Id].Start()    
  }
  
  // Listen for data from our fetchers
  go t.doDataListen(t.Users[user.Id])
  
}

//
// Listen for message on the data channel.
//
func (t * Base) doDataListen(user *User) {
  
  // Listen for data from the fetcher. We then send it up the main channel after
  // adding a userId to the object. We do it this way because int he future we could do 
  // extra processing here.
  for {
  
    select {
    
      case message := <-user.DataChannel:
        
        // Send this message up the chain to the master channel
        t.WsWriteChannel <- websocket.SendStruct{ UserId: user.Profile.Id, Message: message }
        
      case message := <-user.QuoteChannel:

        // Send this message up the chain to the master channel
        t.WsWriteQuoteChannel <- websocket.SendStruct{ UserId: user.Profile.Id, Message: message }
       
    
    }
  
  }  
  
}

/* End File */