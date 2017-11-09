//
// Date: 11/4/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"net/http"

	"app.options.cafe/backend/brokers/tradier"
)

//
// Do Routes
//
func (t *Controller) DoRoutes(mux *http.ServeMux) {

	// Http Routes
	mux.HandleFunc("/login", t.DoLogin)
	mux.HandleFunc("/register", t.DoRegister)
	mux.HandleFunc("/reset-password", t.DoResetPassword)
	mux.HandleFunc("/forgot-password", t.DoForgotPassword)

	// Webhooks
	mux.HandleFunc("/webhooks/stripe", t.DoStripeWebhook)

	// Tradier Oauth
	tr := &tradier.TradierAuth{DB: t.DB}
	mux.HandleFunc("/tradier/authorize", tr.DoAuthCode)
	mux.HandleFunc("/tradier/callback", tr.DoAuthCallback)

	// Setup websocket
	mux.HandleFunc("/ws/core", t.DoWebsocketConnection)
	mux.HandleFunc("/ws/quotes", t.DoQuoteWebsocketConnection)

	// Static files.
	mux.Handle("/", http.FileServer(http.Dir("/frontend")))
}

/* End File */
