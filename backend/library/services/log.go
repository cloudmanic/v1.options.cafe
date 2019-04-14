//
// Date: 5/31/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package services

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mgutz/ansi"
	"github.com/stvp/rollbar"
)

//
// Info Log. Used to log inportant information but a human does not need to review unless they are debugging something.
//
func Info(err error) {
	log.Println("[App:Info] " + MyCaller() + " : " + ansi.Color(err.Error(), "magenta"))
}

//
// Allow us to pass in just text instead of an error. Not always do we have a err to pass in.
//
func InfoMsg(msg string) {
	Info(errors.New(msg))
}

//
// Critical - We used this when we want to make splash. All Critical errors should be reviewed by a human.
//
func Critical(err error) {
	// Standard out
	log.Println(ansi.Color("[App:Critical] "+MyCaller()+" : "+err.Error(), "yellow"))

	// Rollbar
	RollbarError(err)
}

//
// Fatal Log. We use this wehn the app should die and not continue running.
//
func Fatal(err error) {
	// Standard out
	log.Fatal(ansi.Color("[App:Fatal] "+MyCaller()+" : "+err.Error(), "red"))

	// Rollbar
	RollbarError(err)
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
