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
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/acme/autocert"
)

//
// Start the webserver
//
func (t *Controller) StartWebServer() {

	// Listen for data from our broker feeds.
	go t.DoWsDispatch()

	// Set Router
	r := mux.NewRouter()

	// Register Routes
	t.DoRoutes(r)

	// Setup handler
	var handler = t.AuthMiddleware(r)

	if os.Getenv("HTTP_LOG_REQUESTS") == "true" {
		handler = handlers.CombinedLoggingHandler(os.Stdout, t.AuthMiddleware(r))
	}

	// Are we in testing mode? If not give us some SSL
	if os.Getenv("APP_ENV") == "local" {

		s := &http.Server{
			Addr:         ":7080",
			Handler:      handler,
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
		StartSecureServer(handler, m.GetCertificate)
	}
}

//
// Here we make sure we passed in a proper Bearer Access Token.
//
func (t *Controller) AuthMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// // Make sure we have a Bearer token.
		// auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

		// if len(auth) != 2 || auth[0] != "Bearer" {
		//   t.RespondError(w, http.StatusUnauthorized, "Authorization Failed (#001)")
		//   return
		// }

		// // Make sure we have a known access token.
		// if auth[1] != os.Getenv("ACCESS_TOKEN") {
		//   t.RespondError(w, http.StatusUnauthorized, "Authorization Failed (#002)")
		//   return
		// }

		// Manage OPTIONS requests
		if os.Getenv("APP_ENV") == "local" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "DNT,X-CustomHeader,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Content-Range,Range")
		}

		// On to next request in the Middleware chain.
		next.ServeHTTP(w, r)

	})
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
