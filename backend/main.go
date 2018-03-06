package main

import (
	"os"
	"runtime"

	"github.com/cloudmanic/app.options.cafe/backend/cmd"
	"github.com/cloudmanic/app.options.cafe/backend/controllers"
	"github.com/cloudmanic/app.options.cafe/backend/library/notify"
	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/cloudmanic/app.options.cafe/backend/models"
	"github.com/cloudmanic/app.options.cafe/backend/users"
	"github.com/cloudmanic/app.options.cafe/backend/websocket"
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

	// See if this a command. If so run the command and do not start the app.
	status := cmd.Run(db)

	if status == true {
		return
	}

	// -------------- If we made it this far it is time to start the http server -------------- //

	// Lets get started
	services.Critical("App Started: " + os.Getenv("APP_ENV"))

	// Close db when this app dies. (This might be useless)
	defer db.Close()

	// Setup shared channels
	WsWriteChan := make(chan websocket.SendStruct, 1000)

	// Setup the notification channel
	notify.SetWebsocketChannel(WsWriteChan)

	// Setup users object & Start users feeds
	u := &users.Base{
		DB:          db,
		Users:       make(map[uint]*users.UserFeed),
		WsWriteChan: WsWriteChan,
	}

	// Start user feed
	u.StartFeeds()

	// Create new websocket
	w := websocket.NewController(db, WsWriteChan)

	// Startup controller & websockets
	c := &controllers.Controller{DB: db, WebsocketController: w}

	// Start market status feed
	go w.StartMarketStatusFeed()

	// Start websockets & controllers
	c.StartWebServer()

}

/* End File */
