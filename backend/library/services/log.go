//
// Date: 5/31/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package services

import (
	"log"
	"os"

	"github.com/stvp/rollbar"
)

//
// Normal Log.
//
func Log(message string) {

	log.Println("[App Log] " + message)
}

//
// Fatal Log.
//
func Fatal(message string) {

	log.Fatal("[App Log] " + message)

	// Rollbar
	RollbarInfo(message)
}

//
// Error Log.
//
func Error(err error, message string) {

	// Standard out
	log.Println("[App Log] " + message + " (" + err.Error() + ")")

	// Rollbar
	RollbarError(err)
}

//
// Error Log. Not a major error. But should log.
//
func LogErrorOnly(err error) {

	// Standard out
	log.Println("[App Log] " + err.Error())
}

//
// Major Log - Log to every place.
//
func MajorLog(message string) {

	// Standard out
	log.Println("[App Log] " + message)

	// Rollbar
	RollbarInfo(message)
}

//
// Major Error - Log to every place.
//
func MajorError(err error, message string) {

	// Standard out
	log.Println("[App Log] " + message + " (" + err.Error() + ")")

	// Rollbar
	RollbarError(err)
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

/* End File */
