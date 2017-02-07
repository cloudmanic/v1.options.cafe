package main

import (
  "fmt" 
  "sync"
  "time"
  "./models"  
  "net/http"
  "encoding/json"
  "github.com/tidwall/gjson"  
  "github.com/gorilla/websocket" 
)

const writeWait = 3 * time.Second

type Websockets struct {
  WsPongChannel chan *websocket.Conn
  WsQuotePongChannel chan *websocket.Conn     
  connections map[*websocket.Conn]*WebsocketConnection
  quotesConnections map[*websocket.Conn]*WebsocketConnection 
}

type WebsocketConnection struct {
  connection *websocket.Conn
  
  muUserId sync.Mutex
  userId uint  
}

type TradierApiKeyStruct struct {
  Type string `json:"type"`
  Data struct {
    Key string `json:"key"` 
  }    
}

type WebsocketSendStruct struct {
  Type string `json:"type"`
  Data string `json:"data"`
}

//
// Check Origin
//
func (t *Websockets) CheckOrigin(r *http.Request) bool {
  
  origin := r.Header.Get("Origin")
  
  if origin == "file://" {
    return true;
  }

  if origin == "http://localhost:4200" {
    return true;
  } 

  if origin == "http://localhost:7652" {
    return true;
  } 
  
  if origin == "https://s3.amazonaws.com" {
    return true;
  }  

  if origin == "https://app.options.cafe" {
    return true;
  } 

  if origin == "https://cdn.options.cafe" {
    return true;
  } 

  return false;
}

//
// Return a json object ready to be sent up the websocket
//
func (t *Websockets) GetSendJson(send_type string, data_json string) (string, error) {
  
  // Send Object
  send := WebsocketSendStruct{
    Type: send_type,
    Data: data_json} 
  
  send_json, err := json.Marshal(send)

  if err != nil {
    return "", err
  } 

  return string(send_json), nil
}

//
// Handle new connections to the app.
//
func (t *Websockets) DoWebsocketConnection(w http.ResponseWriter, r *http.Request) {

  // setup upgrader
  var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: t.CheckOrigin,
  }

  // Upgrade connection
  conn, err := upgrader.Upgrade(w, r, nil)

  if err != nil {
    fmt.Println(err)
    return
  }

  fmt.Println("New Websocket Connection - Standard")
  
  // Close connection when this function ends
  defer conn.Close()  

  // Add the connection to our connection array
  r_con := WebsocketConnection{ connection: conn }
    
  t.connections[conn] = &r_con 
  
  // Start handling reading messages from the client.
  t.DoWebsocketReading(&r_con)
}

//
// Authenticate Connection
//
func (t *Websockets) AuthenticateConnection(conn *WebsocketConnection, access_token string) {
    
  var user models.User
  
  if db.First(&user, "access_token = ?", access_token).RecordNotFound() {
    fmt.Println("Access Token Not Found - Unable to Authenticate")
    return
  }
  
  fmt.Println("Authenticated : " + user.Email)
  
  // Store the user id from this connection because the auth was successful
  conn.muUserId.Lock()
  conn.userId = user.Id
  conn.muUserId.Unlock()
  
} 

//
// Do reading message.
//
func (t *Websockets) DoWebsocketReading(conn *WebsocketConnection) {  
  
  for {
    
    // Block waiting for a message to arrive
    mt, message, err := conn.connection.ReadMessage()
				
		// Connection closed.
    if mt < 0 {
      
			_, ok := t.connections[conn.connection]
			
			if ok {
				delete(t.connections, conn.connection)
			}    
      
      fmt.Println("Client Disconnected...")
      
      break
    }
		 
    // this should come after the mt test. 
    if err != nil {
			fmt.Println(err)
      break
    }		  
		    
    // Json decode message.
    var data map[string]interface{}
    if err := json.Unmarshal(message, &data); err != nil {
			println("json:", err)
      break      
    }
    
    // Switch on the type of requests.
    switch data["type"] {
      
      // Ping to make sure we are alive.
      case "ping":
        t.WsPongChannel <- conn.connection
      break;

      // The user authenticates.
      case "set-access-token":
        access_token := gjson.Get(string(message), "data.access_token").String()
        t.AuthenticateConnection(conn, access_token)
      break;      
    }
        
  }  
}

//
// Send data out the websocket
//
func (t *Websockets) DoWebsocketWriting(user *UsersConnection) {
  
  for {
    
    select {
      
      // Standard write message.
      case message := <-user.WebsocketWriteChannel:
      
        for i, _ := range t.connections {
          
          // We only care about the user we passed in.
          if t.connections[i].userId != user.UserId {
            continue
          }
          
          t.connections[i].connection.WriteMessage(websocket.TextMessage, []byte(message))
          
        }      
    
      // Manage pongs back to the browser
      case conn := <-t.WsPongChannel:
      
        conn.WriteMessage(websocket.TextMessage, []byte("{\"type\":\"pong\"}"))    
      
    }
    
  }
  	  
}

// ---------------------------- Special Handling For Quotes ----------------------- //

//
// Handle new quote connections to the app.
//
func (t *Websockets) DoQuoteWebsocketConnection(w http.ResponseWriter, r *http.Request) {

  // setup upgrader
  var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: t.CheckOrigin,
  }

  // Upgrade connection
  conn, err := upgrader.Upgrade(w, r, nil)

  if err != nil {
    fmt.Println(err)
    return
  }

  fmt.Println("New Websocket Connection - Quote")
  
  // Close connection when this function ends
  defer conn.Close()
	  
  // Add the connection to our connection array
  r_con := WebsocketConnection{connection: conn}
  
  t.quotesConnections[conn] = &r_con
  
  // Do reading
  t.DoQuoteWebsocketRead(&r_con); 
}

//
// Send data out quote data websocket 
//
func (t *Websockets) DoWebsocketQuoteWriting(user *UsersConnection) {
  
  for {
    
    select {

      // Standard write message.
      case message := <-user.WebsocketWriteQuoteChannel:
    
        for i, _ := range t.quotesConnections {

          // We only care about the user we passed in.
          if t.quotesConnections[i].userId != user.UserId {
            continue
          }
          
          t.quotesConnections[i].connection.WriteMessage(websocket.TextMessage, []byte(message))
          
        }
     
      // Manage pongs back to the browser
      case conn := <-t.WsQuotePongChannel:
      
        conn.WriteMessage(websocket.TextMessage, []byte("{\"type\":\"pong\"}"))
      
    }
    
  }  
	  
}

//
// Handle a message from the web client with Quotes.
//
func (t *Websockets) DoQuoteWebsocketRead(conn *WebsocketConnection) {
  
  for {
    
    // Block waiting for a message to arrive
    mt, message, err := conn.connection.ReadMessage()    
  		
		// Connection closed.
    if mt < 0 {
      
			_, ok := t.quotesConnections[conn.connection]
			
			if ok {
				delete(t.quotesConnections, conn.connection)
			}     
      
      fmt.Println("Quote Client Disconnected...")
      
      break
    } 
    
    // this should come after the mt test. 
    if err != nil {
			fmt.Println(err)
      break
    }	     
    
    // Json decode message.
    var data map[string]interface{}
    
    if err := json.Unmarshal(message, &data); err != nil {
      println("json:", err)
      break     
    }
    
    // Switch on the type of requests.
    switch data["type"] {
      
      // Ping to make sure we are alive.
      case "ping":
        t.WsQuotePongChannel <- conn.connection       
      break;
      
      // The user authenticates.
      case "set-access-token":
        access_token := gjson.Get(string(message), "data.access_token").String()
        t.AuthenticateConnection(conn, access_token)
      break;     
    }
  
  }
      
}

/* End File */