package main

import (
  "fmt" 
  "sync"
  "time"  
  "net/http"
  "encoding/json"
  "github.com/tidwall/gjson"  
  "github.com/gorilla/websocket" 
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
	
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second
	
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

type Websockets struct {   
  connections map[*websocket.Conn]*WebsocketConnection
  quotesConnections map[*websocket.Conn]*WebsocketConnection 
}

type WebsocketConnection struct {
  muWrite sync.Mutex 
  connection *websocket.Conn

  muLicenseKey sync.Mutex
  licenseKey string
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

  println("New Websocket Connection - Standard")
  
  // Close connection when this function ends
  defer conn.Close()  

  // Add in the ping (rom server) and the pong (back from browser)
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error { conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	
  // Add the connection to our connection array
  r_con := WebsocketConnection{
    connection: conn,
    licenseKey: "",
  }
    
  t.connections[conn] = &r_con 

  // Send a message that we are connected.
  r_con.muWrite.Lock()
  r_con.connection.WriteMessage(websocket.TextMessage, []byte("{\"type\":\"Status:connected\"}"))
  r_con.muWrite.Unlock()
  
  // Start handling reading messages from the client.
  t.DoWebsocketReading(&r_con)
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
        conn.muWrite.Lock();
        conn.connection.WriteMessage(websocket.TextMessage, []byte("{\"type\":\"pong\"}"))
        conn.muWrite.Unlock();
      break;
      
      // Get the tradier API sent from the frontend
      case "tradier-api-key":
        t.DoTradierApiKey(conn, message)
      break;
      
      // Set the account id we should be using when talking to the broker.
      case "set-account-id":
        id := gjson.Get(string(message), "data.id")
        SetBrokerAccountId(conn.licenseKey, id.String())
      break;
      
    }
        
  }  
}

//
// Send data out the websocket
//
func (t *Websockets) DoWebsocketSending() {
  
  // Setup the ping ticker
  ticker := time.NewTicker(pingPeriod)
  defer ticker.Stop()
  
  for {

    select {
      
      // Websocket Channel
      case message := <-websocketSendChannel:
      
        //fmt.Println(string(message))
        
        // TODO: Filter this by lisc. keys
        for _, row := range t.connections {
          
          //fmt.Println("Here1")
          
          row.muWrite.Lock()
          row.connection.WriteMessage(websocket.TextMessage, []byte(message))
          row.muWrite.Unlock()
          
          //fmt.Println("Here2")
        }
        
      break      
     
      // Ticker - Send ping messages
      case <-ticker.C:
              
        for _, row := range t.connections {
          row.muWrite.Lock()
          
          row.connection.SetWriteDeadline(time.Now().Add(writeWait))
          
          if err := row.connection.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
            fmt.Println("Client Disconnected... (Error Ticker)")
			    }
			    
			    row.muWrite.Unlock()
			  }
			  
      break
    }

  }
	  
}

//
// Grab the tradier API key sent over from the front-end
//
func (t *Websockets) DoTradierApiKey(conn *WebsocketConnection, message []byte) (string, error) {
  
  var obj TradierApiKeyStruct
          
  if err := json.Unmarshal(message, &obj); err != nil {
    fmt.Println("json: (do_tradier_api_key)", err)
    return "", err      
  } 
  
  // Start collecting data for this connection if we have not already started
  StartBrokerFeed("fa9a93242ds234kasdf", string(obj.Data.Key))
  
  // TODO: change this to the real lisc key
  conn.muLicenseKey.Lock()
  conn.licenseKey = "fa9a93242ds234kasdf"
  conn.muLicenseKey.Unlock()
  
  return obj.Data.Key, nil
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
  
  // Add in the ping (rom server) and the pong (back from browser)
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error { conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	  
  // Add the connection to our connection array
  r_con := WebsocketConnection{connection: conn}
 
  t.quotesConnections[conn] = &r_con
  
  // Do reading
  t.DoQuoteWebsocketRead(&r_con); 
}

//
// Send data out quote data websocket 
//
func (t *Websockets) DoWebsocketQuoteSending() {
  
  // Setup the ping ticker
  ticker := time.NewTicker(pingPeriod)
  defer ticker.Stop()
  
  for {

    select {
      
      // Websocket Channel
      case message := <-websocketSendQuoteChannel:
      
        // TODO: Filter this by lisc. keys
        for _, row := range t.quotesConnections {
          row.muWrite.Lock();
          row.connection.WriteMessage(websocket.TextMessage, []byte(message));
          row.muWrite.Unlock();
        }
        
      break      
     
      // Ticker - Send ping messages
      case <-ticker.C:
              
        for _, row := range t.quotesConnections {
          row.muWrite.Lock()
          
          row.connection.SetWriteDeadline(time.Now().Add(writeWait))
          
          if err := row.connection.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
            fmt.Println("Client Quote Disconnected... (Error Ticker)")
			    }
			    
			    row.muWrite.Unlock()
			  }
			  
      break
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
      
      println("Quote Client Disconnected...")
      
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
        conn.muWrite.Lock();
        conn.connection.WriteMessage(websocket.TextMessage, []byte("{\"type\":\"pong\"}"))
        conn.muWrite.Unlock();        
      break;
      
    }
  
  }
      
}

/* End File */
