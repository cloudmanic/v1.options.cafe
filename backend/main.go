package main

import (
	"os"
	"runtime"

	"app.options.cafe/backend/controllers"
	"app.options.cafe/backend/library/services"
	"app.options.cafe/backend/models"
	"app.options.cafe/backend/users"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

//
// Main....
//
func main() {

	// Setup CPU stuff.
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		services.Fatal("Error loading .env file")
	}

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
	WsReadChan := make(chan controllers.SendStruct, 1000)
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
