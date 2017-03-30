package main

import (
  "fmt" 
  "sync"
  "time" 
  "net/http"
  "encoding/json"
  "github.com/stvp/rollbar"
  "github.com/gorilla/websocket" 
)

const writeWait = 5 * time.Second

type Websockets struct { 
  connections map[*websocket.Conn]*WebsocketConnection
  quotesConnections map[*websocket.Conn]*WebsocketConnection
  controller WsController 
}

type WebsocketConnection struct {
  writeChan chan string
  connection *websocket.Conn
  
  muUserId sync.Mutex
  userId uint 
  
  muDeviceId sync.Mutex
  deviceId string    
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

  if origin == "https://app.options.cafe" {
    return true;
  } 

  return false;
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
  r_con := WebsocketConnection{ connection: conn, writeChan: make(chan string, 100) }
    
  t.connections[conn] = &r_con 
  
  // Start handling reading messages from the client.
  t.DoWsReading(&r_con)
}

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
    rollbar.Error(rollbar.ERR, err)
    fmt.Println(err)
    return
  }

  fmt.Println("New Websocket Connection - Quote")
  
  // Close connection when this function ends
  defer conn.Close()
	  
  // Add the connection to our connection array
  r_con := WebsocketConnection{connection: conn, writeChan: make(chan string, 1000)}
  
  t.quotesConnections[conn] = &r_con
  
  // Do reading
  t.DoWsReading(&r_con) 
}

//
// Start a writer for the websocket connection.
//
func (t *Websockets) DoWsWriting(conn *WebsocketConnection) {
  
  conn.connection.SetWriteDeadline(time.Now().Add(writeWait))
  
  for {
    
    message := <-conn.writeChan
    conn.connection.WriteMessage(websocket.TextMessage, []byte(message))
    conn.connection.SetWriteDeadline(time.Now().Add(writeWait))
    
  }
  
}

//
// Start a reader for this websocket connection.
//
func (t *Websockets) DoWsReading(conn *WebsocketConnection) {  
  
  for {
    
    // Block waiting for a message to arrive
    mt, message, err := conn.connection.ReadMessage()
				
		// Connection closed.
    if mt < 0 {
      
			_, ok := t.connections[conn.connection]
			
			if ok {
				delete(t.connections, conn.connection)
			}    
      
      fmt.Println("Client Disconnected (" + conn.deviceId + ") ...")
      
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
      rollbar.Error(rollbar.ERR, err)
			fmt.Println("json:", err)
      break      
    }
    
    // Switch on the type of requests.
    t.controller.ProcessRead(conn, string(message), data)
        
  }  
}

//
// Grab data and dispatch it to the different websocket connections.
//
func (t *Websockets) DoWsDispatch(user *UsersConnection) {

  for {
    
    select {

      // Core channel
      case message := <-user.WsWriteChannel:
      
        for i, _ := range t.connections {
          
          // We only care about the user we passed in.
          if t.connections[i].userId == user.UserId {
            
            select {
              
              case t.connections[i].writeChan <- message:
 	
              default:
                rollbar.Message("error", "Channel full. Discarding value (Core channel)")
                fmt.Println("Channel full. Discarding value (Core channel)")   
                          
            }
            
          }
          
        }
      
      // Quotes channel
      case message := <-user.WsWriteQuoteChannel:
      
        for i, _ := range t.quotesConnections {
          
          // We only care about the user we passed in.
          if t.quotesConnections[i].userId == user.UserId {
            
             select {
              
              case t.quotesConnections[i].writeChan <- message:
 	
              default:
                rollbar.Message("error", "Channel full. Discarding value (Quotes channel)")
                fmt.Println("Channel full. Discarding value (Quotes channel)")   
                          
            }           
            
          }
          
        }
              
    }
      
  }
  	  
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

/* End File */