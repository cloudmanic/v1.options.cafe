package users

import (
  //"fmt"
  "app.options.cafe/backend/models"
  "app.options.cafe/backend/brokers"
  "app.options.cafe/backend/brokers/feed"
  "app.options.cafe/backend/brokers/tradier"  
  "app.options.cafe/backend/library/services"
)

type Base struct {
  DB *models.DB
  Users map[uint]*User
}

type User struct {
  Profile models.User
  BrokerFeed map[uint]*feed.Base
  WsWriteChannel chan string
  WsWriteQuoteChannel chan string
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
    go t.doUser(users[i])     
  }  
  
} 

//
// Start one user.
//
func (t * Base) doUser(user models.User) {
 
  var brokerApi brokers.Api
  var wsWriteChannel chan string
  var wsWriteQuoteChannel chan string     
 
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
                        WsWriteChannel: make(chan string),
                        WsWriteQuoteChannel: make(chan string),    
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
                                            WsWriteChannel: wsWriteChannel,
                                            WsWriteQuoteChannel: wsWriteQuoteChannel,                       
                                          }
                      
    // Start fetching data for this user.
    t.Users[user.Id].BrokerFeed[row.Id].Start()      
    
/*    
    // Set the broker we are going to use.
    var broker = tradier.Api{ ApiKey: row.AccessToken }
    var fetch = Fetch{ broker: broker, user: userConnections[user.Id] }
    
    // Set Broker hash and lets get going.
    userConnections[user.Id].BrokerConnections[row.Id] = &BrokerFeed{ fetch: fetch, userId: user.Id }
  
    // Start the broker feed.
    userConnections[user.Id].BrokerConnections[row.Id].Start()
*/     
    
  }  
  
}

/*
type Connection struct {
  UserId uint
  BrokerConnections map[uint]*feed.Broker
  WsWriteChannel chan string
  WsWriteQuoteChannel chan string
}

//
// Start up our user connections.
//
func StartUserConnections() {
  
  var user models.User
  usrs := user.GetAllUsers(db)

  for i, _ := range usrs {
    go StartUserConnection(usrs[i])     
  }  
  
}

//
// Start one user connection.
//
func StartUserConnection(user models.User) {
  
  fmt.Println("Starting User Connection : " + user.Email)
  
  // Lock the memory
	mu.Lock()
	defer mu.Unlock() 

  // Make sure we do not already have this licenseKey going.
  if _, ok := userConnections[user.Id]; ok {
    rollbar.Message("info", "User Connection Is Already Going : " + user.Email)
    fmt.Println("User Connection Is Already Going : " + user.Email)
    return
  }
  
  // Set the user connection.
  userConnections[user.Id] = &http_ws.Connection{
    UserId: user.Id,
    BrokerConnections: make(map[uint]*feed.Broker),
    WsWriteChannel: make(chan string),
    WsWriteQuoteChannel: make(chan string),
  }
  
  // Start the websocket write dispatcher for this user.
  go ws.DoWsDispatch(userConnections[user.Id])
  
  // Loop through the different brokers for this user
  for _, row := range user.Brokers {
    
    // Need an access token to continue
    if len(row.AccessToken) <= 0 {
      rollbar.Message("info", "User Connection (Brokers) No Access Token Found : " + user.Email)
      fmt.Println("User Connection (Brokers) No Access Token Found : " + user.Email)
      continue
    }
    
    // Set the broker we are going to use.
    var broker = tradier.Api{ ApiKey: row.AccessToken }
    var fetch = Fetch{ broker: broker, user: userConnections[user.Id] }
    
    // Set Broker hash and lets get going.
    userConnections[user.Id].BrokerConnections[row.Id] = &BrokerFeed{ fetch: fetch, userId: user.Id }
  
    // Start the broker feed.
    userConnections[user.Id].BrokerConnections[row.Id].Start()     
    
  }
    
}
*/