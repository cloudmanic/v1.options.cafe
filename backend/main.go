package main

import (
  "os"
  "runtime"  
  "github.com/joho/godotenv"
  "app.options.cafe/backend/users"
  "app.options.cafe/backend/models"
  "app.options.cafe/backend/controllers"  
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
  
  // Connect to database and run Migrations.
  var DB = models.DB{}
  DB.Start()
  defer DB.Connection.Close()  
  
  // Setup websockets & controllers
  controllers.DB = &DB
  controllers.WsReadChan = make(chan controllers.SendStruct, 1000)
  controllers.WsWriteChan = make(chan controllers.SendStruct, 1000)
  controllers.WsWriteQuoteChan = make(chan controllers.SendStruct, 1000)
    
  // Setup users object & Start users feeds
  users.DB = &DB
  users.DataChan = controllers.WsWriteChan
  users.QuoteChan = controllers.WsWriteQuoteChan
  users.FeedRequestChan = controllers.WsReadChan    
  users.StartFeeds()

  // Start websockets & controllers
  controllers.Start()
  
} 

/* End File */
