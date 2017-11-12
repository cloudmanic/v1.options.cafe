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

	// --------- API V1 sub-routes ----------- //

	apiV1 := r.PathPrefix("/api/v1").Subrouter()

	// Symbols
	apiV1.Handle("/symbols", t.AuthMiddleware(http.HandlerFunc(t.GetSymbols))).Methods("GET", "OPTIONS")

	// Watchlists
	apiV1.Handle("/watchlists", t.AuthMiddleware(http.HandlerFunc(t.GetWatchlists))).Methods("GET", "OPTIONS")
	apiV1.Handle("/watchlists", t.AuthMiddleware(http.HandlerFunc(t.CreateWatchlist))).Methods("POST", "OPTIONS")
	apiV1.Handle("/watchlists/{id:[0-9]+}", t.AuthMiddleware(http.HandlerFunc(t.GetWatchlist))).Methods("GET", "OPTIONS")

	// ------- End API V1 sub-routes --------- //

	// Auth Routes
	r.Handle("/login", t.CorsMiddleware(http.HandlerFunc(t.DoLogin))).Methods("POST", "OPTIONS")
	r.Handle("/register", t.CorsMiddleware(http.HandlerFunc(t.DoRegister))).Methods("POST", "OPTIONS")
	r.Handle("/reset-password", t.CorsMiddleware(http.HandlerFunc(t.DoResetPassword))).Methods("POST", "OPTIONS")
	r.Handle("/forgot-password", t.CorsMiddleware(http.HandlerFunc(t.DoForgotPassword))).Methods("POST", "OPTIONS")

	// Webhooks
	r.HandleFunc("/webhooks/stripe", t.DoStripeWebhook).Methods("POST")

	// Tradier Oauth
	tr := &tradier.TradierAuth{DB: t.DB}
	r.HandleFunc("/tradier/authorize", tr.DoAuthCode).Methods("GET")
	r.HandleFunc("/tradier/callback", tr.DoAuthCallback).Methods("GET")

	// Setup websocket
	r.HandleFunc("/ws/core", t.DoWebsocketConnection)
	r.HandleFunc("/ws/quotes", t.DoQuoteWebsocketConnection)

	// Static files.
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("/frontend")))
}

/* End File */
