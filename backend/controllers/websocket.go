package controllers

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

type SendStruct struct {
  UserId uint
  Message string
}

var (
  connections = make(map[*websocket.Conn]*WebsocketConnection)
  quotesConnections = make(map[*websocket.Conn]*WebsocketConnection)
)
  
//
// Check Origin
//
func CheckOrigin(r *http.Request) bool {
  
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
func DoWebsocketConnection(w http.ResponseWriter, r *http.Request) {

  // setup upgrader
  var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: CheckOrigin,
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
    
  connections[conn] = &r_con 
  
  // Start handling reading messages from the client.
  DoWsReading(&r_con)
}

//
// Handle new quote connections to the app.
//
func DoQuoteWebsocketConnection(w http.ResponseWriter, r *http.Request) {

  // setup upgrader
  var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: CheckOrigin,
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
  
  quotesConnections[conn] = &r_con
  
  // Do reading
  DoWsReading(&r_con) 
}

//
// Start a writer for the websocket connection.
//
func DoWsWriting(conn *WebsocketConnection) {
  
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
func DoWsReading(conn *WebsocketConnection) {  
  
  for {
    
    // Block waiting for a message to arrive
    mt, message, err := conn.connection.ReadMessage()
				
		// Connection closed.
    if mt < 0 {
      
			_, ok := connections[conn.connection]
			
			if ok {
				delete(connections, conn.connection)
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
    ProcessRead(conn, string(message), data)
        
  }  
}

/* End File */