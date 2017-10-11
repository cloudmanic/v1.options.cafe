//
// Date: 5/31/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package main

import(
  "os"
  "runtime"
  "github.com/jasonlvhit/gocron"
  "app.options.cafe/backend/cron/data_import" 
  "app.options.cafe/backend/library/services"  
)

//
// Main....
//
func main() {
  
  // Setup CPU stuff.
  runtime.GOMAXPROCS(runtime.NumCPU())  
           
  // Lets get started
  services.MajorLog("Cron Started: " + os.Getenv("APP_ENV"))

  // Setup jobs we need to run 
  gocron.Every(1).Day().At("22:00").Do(data_import.DoEodOptionsImport) 

  // function Start start all the pending jobs
  <- gocron.Start()
}

/* End File */