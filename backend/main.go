package main

import (
  "os"
  "runtime"  
  "github.com/joho/godotenv"
  "app.options.cafe/backend/users"
  "app.options.cafe/backend/models"
  "app.options.cafe/backend/websocket"  
  "app.options.cafe/backend/library/services"
)
      
//
// Main....
//
func main() {  
   
  // Setup CPU stuff.
  runtime.GOMAXPROCS(runtime.NumCPU())  
         
  // Load .env file 
  err := godotenv.Load()
  if err != nil {
    services.Fatal("Error loading .env file")
  }        
       
  // Lets get started
  services.MajorLog("App Started: " + os.Getenv("APP_ENV")) 
  
  // Setup channels
  websocket.WsWriteChannel = make(chan websocket.SendStruct, 1000)
  websocket.WsWriteQuoteChannel = make(chan websocket.SendStruct, 1000)
  
  // Connect to database and run Migrations.
  var DB = models.DB{}
  DB.Start()
  
  // Setup users object
  var Users = users.Base{ 
                      DB: &DB,
                      WsWriteChannel: websocket.WsWriteChannel,
                      WsWriteQuoteChannel: websocket.WsWriteQuoteChannel,
                    }
                    
  // Start users feeds
  Users.Start()

  // Start websockets
  websocket.Start()
  
} 

/* End File */
