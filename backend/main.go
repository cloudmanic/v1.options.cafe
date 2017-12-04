package main

import (
	"os"
	"runtime"

	"github.com/app.options.cafe/backend/cmd"
	"github.com/app.options.cafe/backend/controllers"
	"github.com/app.options.cafe/backend/library/services"
	"github.com/app.options.cafe/backend/models"
	"github.com/app.options.cafe/backend/users"
	"github.com/gorilla/websocket"
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
	status := cmd.Run()

	if status == true {
		return
	}

	// -------------- If we made it this far it is time to start the http server -------------- //

	// Lets get started
	services.Critical("App Started: " + os.Getenv("APP_ENV"))

	// Close db when this app dies. (This might be useless)
	defer db.Close()

	// Setup shared channels
	WsReadChan := make(chan controllers.ReceivedStruct, 1000)
	WsWriteChan := make(chan controllers.SendStruct, 1000)
	WsWriteQuoteChan := make(chan controllers.SendStruct, 1000)

	// Setup users object & Start users feeds
	u := &users.Base{
		DB:              db,
		Users:           make(map[uint]*users.UserFeed),
		DataChan:        WsWriteChan,
		QuoteChan:       WsWriteQuoteChan,
		FeedRequestChan: WsReadChan,
	}

	// Start user feed
	u.StartFeeds()

	// Startup controller & websockets
	c := &controllers.Controller{
		DB:                db,
		WsReadChan:        WsReadChan,
		WsWriteChan:       WsWriteChan,
		WsWriteQuoteChan:  WsWriteQuoteChan,
		Connections:       make(map[*websocket.Conn]*controllers.WebsocketConnection),
		QuotesConnections: make(map[*websocket.Conn]*controllers.WebsocketConnection),
	}

	// Start websockets & controllers
	c.StartWebServer()

}

/* End File */
