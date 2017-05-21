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
  
  // Connect to database and run Migrations.
  var DB = models.DB{}
  DB.Start()
  
  // Loop through and startup user feeds
  var Users = users.Base{ DB: &DB }
  Users.Start()

  // Start websockets
  websocket.Start()
  
} 

/* End File */
