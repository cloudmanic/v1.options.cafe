//
// Date: 5/31/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package services

import (
  "os"
  "log"
  "log/syslog"
  "github.com/stvp/rollbar"  
)

//
// Normal Log.
//
func Log(message string) {
  
  log.Println(message)

  // Papertrail log.
  PaperTrailLog(message, "info")
}

//
// Fatal Log.
//
func Fatal(message string) {
  
  log.Fatal(message)

  // Papertrail log.
  PaperTrailLog(message, "error")  
  
  // Rollbar
  RollbarInfo(message)  
}

//
// Error Log.
//
func Error(err error, message string) {
  
  // Standard out
  log.Println(message + " (" + err.Error()  + ")")

  // Papertrail log.
  PaperTrailLog(message, "error")  

  // Rollbar
  RollbarError(err)
}

//
// Major Log - Log to every place.
//
func MajorLog(message string) {
  
  // Standard out
  log.Println(message)

  // Papertrail log.
  PaperTrailLog(message, "info")   
  
  // Rollbar
  RollbarInfo(message)
}

//
// Send log to rollbar
//
func RollbarInfo(message string) {

  if len(os.Getenv("ROLLBAR_TOKEN")) > 0 {

    go func() {
      rollbar.Token = os.Getenv("ROLLBAR_TOKEN")
      rollbar.Environment = os.Getenv("ROLLBAR_ENV")
      rollbar.Message("info", message)
      rollbar.Wait()
    }()

  }
}

//
// Send log to rollbar
//
func RollbarError(err error) {

  if len(os.Getenv("ROLLBAR_TOKEN")) > 0 {

    go func() {
      rollbar.Token = os.Getenv("ROLLBAR_TOKEN")
      rollbar.Environment = os.Getenv("ROLLBAR_ENV")
      rollbar.Error(rollbar.ERR, err) 
      rollbar.Wait()
    }()

  }
}

//
// Send log to PaperTrail 
//
// We might not need this as we run within docker and we just setup docker 
// to send logs to PaperTrail. 
//
func PaperTrailLog(message string, msgType string) {

  if len(os.Getenv("PAPERTRAIL_URL")) > 0 {

    go func() {
      w, err := syslog.Dial("udp", os.Getenv("PAPERTRAIL_URL"), syslog.LOG_EMERG | syslog.LOG_KERN, os.Getenv("APP_NAME"))
    
      if err != nil {
        Error(err, "Failed to dial syslog (PaperTrailLog)")
      } 

      // Info log
      if msgType == "info" {
        w.Info(message)
      } 

      // Error log
      if msgType == "error" {
        w.Err(message)
      } 

    }()

  }
}

/* End File */