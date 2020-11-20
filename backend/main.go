package main

import (
	"os"
	"runtime"

	"app.options.cafe/cmd"
	"app.options.cafe/controllers"
	"app.options.cafe/cron"
	"app.options.cafe/library/polling"
	"app.options.cafe/library/queue"
	"app.options.cafe/library/seed"
	"app.options.cafe/library/services"
	"app.options.cafe/library/worker/jobs"
	"app.options.cafe/models"
	"app.options.cafe/websocket"
	_ "github.com/jpfuentes2/go-env/autoload"
)

//
// Main....
//
func main() {
	// Setup CPU stuff.
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Start the db connection.
	db, err := models.NewDB()

	if err != nil {
		services.Fatal(err)
	}

	// See local database (if we are in local dev ENV)
	seed.LocalDatabase(db)

	// Fire up the queue connection
	queue.Start()

	// See if this a command. If so run the command and do not start the app.
	status := cmd.Run(db)

	if status == true {
		return
	}

	// -------------- If we made it this far it is time to start the http server -------------- //

	// Lets get started
	services.InfoMsg("App Started: " + os.Getenv("APP_ENV"))

	// Close db when this app dies. (This might be useless)
	defer db.Close()

	// Setup shared channels
	WsWriteChan := make(chan websocket.SendStruct, 1000)

	// Create new websocket
	w := websocket.NewController(db, WsWriteChan)

	// Startup controller & websockets
	c := &controllers.Controller{DB: db, WebsocketController: w}

	// Stuff we start a as a different process in production using --cmd. In local dev we start here
	// "-cmd=cron"
	// --cmd=worker
	// "-cmd=broker-feed-poller"
	if os.Getenv("APP_ENV") == "local" {
		polling.Start(db)
		go cron.Start(db)
		go jobs.Start(db)
	}

	// Start websockets & controllers
	c.StartWebServer()
}

/* End File */
