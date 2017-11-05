//
// Date: 11/4/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"log"
	"net/http"
	"time"

	"app.options.cafe/backend/library/services"
)

//
// Start the webserver
//
func (t *Controller) StartWebServer() {

	// Listen for data from our broker feeds.
	go t.DoWsDispatch()

	// Register some handlers:
	mux := http.NewServeMux()

	// Register Routes
	t.DoRoutes(mux)

	// Are we in testing mode?
	s := &http.Server{
		Addr:         ":7080",
		Handler:      mux,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	}

	log.Fatal(s.ListenAndServe())
}

//
// Listen for data from our broker feeds.
// Take the data and then pass it up the websockets.
//
func (t *Controller) DoWsDispatch() {

	for {

		select {

		// Core channel
		case send := <-t.WsWriteChan:

			for i := range t.Connections {

				// We only care about the user we passed in.
				if t.Connections[i].userId == send.UserId {

					select {

					case t.Connections[i].writeChan <- send.Message:

					default:
						services.MajorLog("Channel full. Discarding value (Core channel)")

					}

				}

			}

		// Quotes channel
		case send := <-t.WsWriteQuoteChan:

			for i := range t.QuotesConnections {

				// We only care about the user we passed in.
				if t.QuotesConnections[i].userId == send.UserId {

					select {

					case t.QuotesConnections[i].writeChan <- send.Message:

					default:
						services.MajorLog("Channel full. Discarding value (Quotes channel)")

					}

				}

			}

		}

	}

}

/* End File */
