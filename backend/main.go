package main

import (
	"os"
	"runtime"

	"github.com/app.options.cafe/backend/controllers"
	"github.com/app.options.cafe/backend/library/services"
	"github.com/app.options.cafe/backend/models"
	"github.com/app.options.cafe/backend/users"
	"github.com/gorilla/websocket"
)

//
// Main....
//
func main() {

	// Setup CPU stuff.
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Lets get started
	services.MajorLog("App Started: " + os.Getenv("APP_ENV"))

	// Start the db connection.
	db, err := models.NewDB()

	if err != nil {
		services.Fatal("Failed to connect database")
	}

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
