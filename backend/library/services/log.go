//
// Date: 5/31/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package services

import (
  "os"
  "log"
  "github.com/stvp/rollbar"  
)

//
// Normal Log.
//
func Log(message string) {
  
  log.Println(message)
  
}

//
// Fatal Log.
//
func Fatal(message string) {
  
  log.Fatal(message)
  
}

//
// Error Log.
//
func Error(err error, message string) {
  
  Log(message + " (" + err.Error()  + ")")

  go func() {
    rollbar.Token = os.Getenv("ROLLBAR_TOKEN")
    rollbar.Environment = os.Getenv("ROLLBAR_ENV")
    rollbar.Error(rollbar.ERR, err)
    rollbar.Wait()
  }()
  
}

//
// Major Log - Log to every place.
//
func MajorLog(message string) {
  
  Log(message)
  
  go func() {
    rollbar.Token = os.Getenv("ROLLBAR_TOKEN")
    rollbar.Environment = os.Getenv("ROLLBAR_ENV")
    rollbar.Message("info", "App started.")
    rollbar.Wait()
  }()
}

/* End File */