//
// Date: 11/4/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"log"
	"net/http"
	"os"
	"time"

	"app.options.cafe/backend/library/services"
	"golang.org/x/crypto/acme/autocert"
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

	// Are we in testing mode? If not give us some SSL
	if os.Getenv("APP_ENV") == "local" {

		s := &http.Server{
			Addr:         ":7080",
			Handler:      mux,
			ReadTimeout:  2 * time.Second,
			WriteTimeout: 2 * time.Second,
		}

		log.Fatal(s.ListenAndServe())

	} else {

		// Secure it with a TLS certificate using Let's  Encrypt:
		m := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			Cache:      autocert.DirCache("/letsencrypt/"),
			Email:      "help@options.cafe",
			HostPolicy: autocert.HostWhitelist("app.options.cafe"),
		}

		// Start a secure server:
		StartSecureServer(mux, m.GetCertificate)
	}
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

					case t.Connections[i].WriteChan <- send.Body:

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

					case t.QuotesConnections[i].WriteChan <- send.Body:

					default:
						services.MajorLog("Channel full. Discarding value (Quotes channel)")

					}

				}

			}

		}

	}

}

/* End File */
