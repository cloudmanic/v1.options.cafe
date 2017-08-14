package controllers

import (
  "os"
  "log"
  "time"
  "net/http"
  "golang.org/x/crypto/acme/autocert"
  "app.options.cafe/backend/models"
  "app.options.cafe/backend/library/services"
)

var (
  DB *models.DB
  WsReadChan chan SendStruct
  WsWriteChan chan SendStruct
  WsWriteQuoteChan chan SendStruct
)

//
// Start the webserver
//
func Start() {
  
  // Listen for data from our broker feeds.
  go DoWsDispatch()
  
  // Register some handlers:
  mux := http.NewServeMux()
  
  // Http Routes
  mux.HandleFunc("/", HtmlMainTemplate)    
  mux.HandleFunc("/register", DoRegister) 
    
  // Setup websocket
	mux.HandleFunc("/ws/core", DoWebsocketConnection)
	mux.HandleFunc("/ws/quotes", DoQuoteWebsocketConnection)

  // Are we in testing mode?
  if os.Getenv("APP_ENV") == "local" {
    
		s := &http.Server{
			Addr: ":7652",
			Handler: mux,
      ReadTimeout:  2 * time.Second,
      WriteTimeout: 2 * time.Second,			
		}
		
		log.Fatal(s.ListenAndServe())    
     
  } else {

    // Secure it with a TLS certificate using Let's  Encrypt:
    m := autocert.Manager{
  	  Prompt: autocert.AcceptTOS,
      Cache: autocert.DirCache("/etc/letsencrypt/"),
      Email: "support@options.cafe",
      HostPolicy: autocert.HostWhitelist("app.options.cafe"),
    }
  
    // Start a secure server:
    StartSecureServer(mux, m.GetCertificate)
    
  }
  
}

//
// Listen for data from our broker feeds.
// Take the data and then pass it up the websockets.
//
func DoWsDispatch() {
  
  for {
    
    select {

      // Core channel
      case send := <-WsWriteChan:
      
        for i, _ := range connections {
          
          // We only care about the user we passed in.
          if connections[i].userId == send.UserId {
            
            select {
              
              case connections[i].writeChan <-send.Message:
 	
              default:
                services.MajorLog("Channel full. Discarding value (Core channel)")   
                          
            }
            
          }
          
        }
      
      // Quotes channel
      case send := <-WsWriteQuoteChan:
      
        for i, _ := range quotesConnections {
          
          // We only care about the user we passed in.
          if quotesConnections[i].userId == send.UserId {
            
             select {
              
              case quotesConnections[i].writeChan <-send.Message:
 	
              default:
                services.MajorLog("Channel full. Discarding value (Quotes channel)")   
                          
            }           
            
          }
          
        }
         
    }
      
  }  
    
}