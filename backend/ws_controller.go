package main

import (
  "fmt" 
  "./models"   
  "github.com/tidwall/gjson"  
)

type WsController struct {}

//
// Process a read request that was sent in from the client
//
func (t *WsController) ProcessRead(conn *WebsocketConnection, message string, data map[string]interface{})  {

  switch data["type"] {
    
    // Ping to make sure we are alive.
    case "ping":
      conn.writeChan <- "{\"type\":\"pong\"}"
    break;
  
    // Refresh all cached data.
    case "refresh-all-data":
      t.RefreshAllData(conn)
    break;
  
    // The user authenticates.
    case "set-access-token":
      device_id := gjson.Get(message, "data.device_id").String()
      access_token := gjson.Get(message, "data.access_token").String()
      t.AuthenticateConnection(conn, access_token, device_id)
    break; 
         
  }
    
}

//
// Authenticate Connection
//
func (t *WsController) AuthenticateConnection(conn *WebsocketConnection, access_token string, device_id string) {
    
  var user models.User
  
  fmt.Println("Device Id : " + device_id)
  
  // Store the device id
  conn.muDeviceId.Lock()
  conn.deviceId = device_id
  conn.muDeviceId.Unlock()  
  
  // See if this user is in our db.
  if db.First(&user, "access_token = ?", access_token).RecordNotFound() {
    fmt.Println("Access Token Not Found - Unable to Authenticate")
    return
  }
  
  fmt.Println("Authenticated : " + user.Email)
  
  // Store the user id from this connection because the auth was successful
  conn.muUserId.Lock()
  conn.userId = user.Id
  conn.muUserId.Unlock()
  
  // Do the writing. 
  go ws.DoWsWriting(conn)
  
  // Send cached data so they do not have to wait for polling.
  for key, _ := range userConnections[conn.userId].BrokerConnections {
    
    userConnections[conn.userId].BrokerConnections[key].fetch.RefreshFromCached()
    
  }

}

//
// A request from the client to send cached data up the 
// websocket. This is often used when the page changes
// or the state of a page changes and they need to 
// refresh the data on the client.
//
func (t *WsController) RefreshAllData(conn *WebsocketConnection)  {
  
  // Send cached data so they do not have to wait for polling.
  for key, _ := range userConnections[conn.userId].BrokerConnections {
    
    userConnections[conn.userId].BrokerConnections[key].fetch.RefreshFromCached()
    
  }
   
}

/* End File */