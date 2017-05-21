package services

import (
  "os"
  "fmt"
  "log"
  "github.com/stvp/rollbar"  
)

//
// Normal Log.
//
func Log(message string) {
  
  fmt.Println(message)
  
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
  
  Log(message)

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