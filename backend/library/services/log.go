//
// Date: 5/31/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package services

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mgutz/ansi"
	"github.com/stvp/rollbar"
)

// --------------- TODO: Get Rid Of These Functions ----------- //

func Error(err error, message string) {
	BetterError(err)
}

// -------------- END Get Rid Of -------------------- //

//
// Info Log.
//
func Info(message string) {
	if !strings.HasSuffix(os.Args[0], ".test") {
		log.Println("[App:Info] " + MyCaller() + " : " + ansi.Color(message, "magenta"))
	}
}

//
// Critical.
//
func Critical(message string) {

	caller := MyCaller()

	// Standard out
	log.Println(ansi.Color("[App:Critical] "+caller+" : "+message, "yellow"))

	// Rollbar
	RollbarInfo(caller + " : " + message)
}

//
// Fatal Log.
//
func Fatal(err error) {

	log.Fatal(ansi.Color("[App:Fatal] "+MyCaller()+" : "+err.Error(), "red"))

	// Rollbar
	RollbarError(err)
}

//
// Fatal with a message Log.
//
func FatalMsg(err error, message string) {

	log.Fatal(ansi.Color("[App:Fatal] "+MyCaller()+" : "+err.Error()+" - "+message, "red"))

	// Rollbar
	RollbarError(err)
}

//
// Warning Log. (error type only)
//
func Warning(err error) {
	if err != nil {
		log.Println(ansi.Color("[App:Warning] "+MyCaller()+" : "+err.Error(), "yellow+b"))
	}
}

//
// BetterError - Error Log. TODO: Convert this to normal Error Bettererror sucks as a name
//
func BetterError(err error) {
	caller := MyCaller()

	// Standard out
	log.Println(ansi.Color("[App:Error] "+caller+" : "+err.Error(), "red"))

	// Rollbar
	RollbarError(err)
}

//
// ErrorMsg - Error with a message Log.
//
func ErrorMsg(err error, message string) {
	caller := MyCaller()

	// Standard out
	log.Println(ansi.Color("[App:Error] "+caller+" : "+err.Error()+" : "+message, "red"))

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

//
// MyCaller returns the caller of the function that called the logger :)
//
func MyCaller() string {
	var filePath string
	var fnName string

	pc, file, line, ok := runtime.Caller(3)

	if !ok {
		file = "?"
		line = 0
	}

	fn := runtime.FuncForPC(pc)

	if fn == nil {
		fnName = "?()"
	} else {
		dotName := filepath.Ext(fn.Name())
		fnName = strings.TrimLeft(dotName, ".") + "()"
	}

	// Make the base of this code.
	parts := strings.Split(file, "app.options.cafe")

	if len(parts) == 2 {
		filePath = "app.options.cafe" + parts[1]
	} else {
		filePath = filepath.Base(file)
	}

	return fmt.Sprintf("%s:%d %s", filePath, line, fnName)
}

/* End File */
