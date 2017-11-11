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

//
// Do Routes
//
func (t *Controller) DoRoutes(r *mux.Router) {

	// Auth Routes
	r.HandleFunc("/login", t.DoLogin).Methods("POST", "OPTIONS")
	r.HandleFunc("/register", t.DoRegister).Methods("POST", "OPTIONS")
	r.HandleFunc("/reset-password", t.DoResetPassword).Methods("POST", "OPTIONS")
	r.HandleFunc("/forgot-password", t.DoForgotPassword).Methods("POST", "OPTIONS")

	// Symbols
	r.HandleFunc("/api/v1/symbols", t.GetSymbols).Methods("GET", "OPTIONS")

	// Watchlists
	r.HandleFunc("/api/v1/watchlists", t.GetWatchlists).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/v1/watchlists", t.CreateWatchlist).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/v1/watchlists/{id:[0-9]+}", t.GetWatchlist).Methods("GET", "OPTIONS")

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

/* End File */
