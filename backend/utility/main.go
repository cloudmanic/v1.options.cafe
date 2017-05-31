//
// Date: 5/31/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package main

import(
  "os"
  "runtime"  
  "github.com/joho/godotenv"
  "app.options.cafe/backend/library/services"
  "app.options.cafe/backend/utility/data_import"  
)

//
// Main....
//
func main() {
  
  // Setup CPU stuff.
  runtime.GOMAXPROCS(runtime.NumCPU())  
         
  // Load .env file 
  err := godotenv.Load("../.env")
  if err != nil {
    services.Fatal("Error loading .env file")
  }    
  
  // Lets get started
  services.MajorLog("Utility Started: " + os.Getenv("APP_ENV"))
  
  // Start up End of Day Import Services
  data_import.DoEodOptionsImport()  
  
}

/* End File */