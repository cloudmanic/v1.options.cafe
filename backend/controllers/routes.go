//
// Date: 11/4/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"net/http"

	"app.options.cafe/backend/brokers/tradier"
	"github.com/gorilla/mux"
)

var noAuthRoutes map[string]bool

//
// Do Routes
//
func (t *Controller) DoRoutes(r *mux.Router) {

	// Set the routes we do not want to authenticate against.
	setRoutesWeSkipAuthOn([]string{
		"/",
		"/login",
		"/register",
		"/forgot-password",
		"/ws/core",
		"/ws/quotes",
		"/webhooks/stripe",
	})

	// --------- API V1 sub-routes ----------- //

	apiV1 := r.PathPrefix("/api/v1").Subrouter()

	// Symbols
	apiV1.HandleFunc("/symbols", t.GetSymbols).Methods("GET")

	// Watchlists
	apiV1.HandleFunc("/watchlists", t.GetWatchlists).Methods("GET")
	apiV1.HandleFunc("/watchlists", t.CreateWatchlist).Methods("POST", "OPTIONS")
	apiV1.HandleFunc("/watchlists/{id:[0-9]+}", t.GetWatchlist).Methods("GET")

	// ------- End API V1 sub-routes --------- //

	// Auth Routes
	r.HandleFunc("/login", t.DoLogin).Methods("POST", "OPTIONS")
	r.HandleFunc("/register", t.DoRegister).Methods("POST", "OPTIONS")
	r.HandleFunc("/reset-password", t.DoResetPassword).Methods("POST", "OPTIONS")
	r.HandleFunc("/forgot-password", t.DoForgotPassword).Methods("POST", "OPTIONS")

	// Webhooks
	r.HandleFunc("/webhooks/stripe", t.DoStripeWebhook).Methods("POST", "OPTIONS")

	// Tradier Oauth
	tr := &tradier.TradierAuth{DB: t.DB}
	r.HandleFunc("/tradier/authorize", tr.DoAuthCode).Methods("GET", "OPTIONS")
	r.HandleFunc("/tradier/callback", tr.DoAuthCallback).Methods("GET", "OPTIONS")

	// Setup websocket
	r.HandleFunc("/ws/core", t.DoWebsocketConnection)
	r.HandleFunc("/ws/quotes", t.DoQuoteWebsocketConnection)

	// Static files.
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("/frontend")))
}

//
// Setup routes we skip auth on.
//
func setRoutesWeSkipAuthOn(skips []string) {

	noAuthRoutes = make(map[string]bool)

	for _, row := range skips {
		noAuthRoutes[row] = true
	}
}

/* End File */
