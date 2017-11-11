//
// Date: 11/4/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
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

		// Setup routes we skip. (TODO: this is a little hacky we should revisit.)
		skip := make(map[string]bool)
		skip["/"] = true
		skip["/login"] = true
		skip["/register"] = true
		skip["/forgot-password"] = true
		skip["/ws/core"] = true
		skip["/ws/quotes"] = true
		skip["/webhooks/stripe"] = true

		if _, ok := skip[r.URL.Path]; ok {
			fmt.Println(r.URL.Path)
			next.ServeHTTP(w, r)
		}

		// Set access token and start the auth proccess
		var access_token = ""

		// Make sure we have a Bearer token.
		auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

		if len(auth) != 2 || auth[0] != "Bearer" {

			// We allow access token from the command line
			if os.Getenv("APP_ENV") == "local" {

				access_token = r.URL.Query().Get("access_token")

				if len(access_token) <= 0 {
					t.RespondError(w, http.StatusUnauthorized, "Authorization Failed (#101)")
					return
				}

			} else {
				t.RespondError(w, http.StatusUnauthorized, "Authorization Failed (#001)")
				return
			}

		} else {
			access_token = auth[1]
		}

		// See if this session is in our db.
		session, err := t.DB.GetByAccessToken(access_token)

		if err != nil {
			services.MajorLog("Access Token Not Found - Unable to Authenticate via HTTP")
			t.RespondError(w, http.StatusUnauthorized, "Authorization Failed (#002)")
			return
		}

		// Get this user is in our db.
		user, err := t.DB.GetUserById(session.UserId)

		if err != nil {
			services.MajorLog("User Not Found - Unable to Authenticate - UserId (HTTP) : " + fmt.Sprint(session.UserId) + " - Session Id : " + fmt.Sprint(session.Id))
			t.RespondError(w, http.StatusUnauthorized, "Authorization Failed (#003)")
			return
		}

		// Manage OPTIONS requests
		if os.Getenv("APP_ENV") == "local" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "DNT,X-CustomHeader,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Content-Range,Range")
		}

		fmt.Println(user.Id)

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
