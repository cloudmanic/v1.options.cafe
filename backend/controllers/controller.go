package controllers


import (
  "fmt"
  "github.com/tidwall/gjson"
  "app.options.cafe/backend/models"
  "app.options.cafe/backend/library/services"    
)

//
// Process a read request that was sent in from the client
//
func ProcessRead(conn *WebsocketConnection, message string, data map[string]interface{})  {

  switch data["type"] {
    
    // Ping to make sure we are alive.
    case "ping":
      conn.writeChan <- "{\"type\":\"pong\"}"
    break;
  
    // Refresh all cached data.
    case "refresh-all-data":
      RefreshAllData(conn)
    break;
  
    // The user authenticates.
    case "set-access-token":
      device_id := gjson.Get(message, "data.device_id").String()
      access_token := gjson.Get(message, "data.access_token").String()
      AuthenticateConnection(conn, access_token, device_id)
    break; 
         
  }
    
}

//
// Authenticate Connection
//
func AuthenticateConnection(conn *WebsocketConnection, access_token string, device_id string) {

  var user models.User    
  var session models.Session
  
  services.Log("Connected Device Id : " + device_id)
  
  // Store the device id
  conn.muDeviceId.Lock()
  conn.deviceId = device_id
  conn.muDeviceId.Unlock()  
  
  // See if this session is in our db.
  if DB.Connection.First(&session, "access_token = ?", access_token).RecordNotFound() {
    services.MajorLog("Access Token Not Found - Unable to Authenticate")
    return
  }
  
  // Get this user is in our db.  
  if DB.Connection.First(&user, session.UserId).RecordNotFound() {
    services.MajorLog("User Not Found - Unable to Authenticate - UserId : " + fmt.Sprint(session.UserId) + " - Session Id : " + fmt.Sprint(session.Id))     
    return
  }  
  
  services.Log("Authenticated : " + user.Email)
  
  // Store the user id from this connection because the auth was successful
  conn.muUserId.Lock()
  conn.userId = user.Id
  conn.muUserId.Unlock()
  
  // Do the writing. 
  go DoWsWriting(conn)
  
  // Send cached data so they do not have to wait for polling.
  RefreshAllData(conn)

}

//
// A request from the client to send cached data up the 
// websocket. This is often used when the page changes
// or the state of a page changes and they need to 
// refresh the data on the client.
//
func RefreshAllData(conn *WebsocketConnection)  {
  
  WsReadChan <- SendStruct{ UserId: conn.userId, Message: "FromCache:refresh" }  
 
}

/* End File */