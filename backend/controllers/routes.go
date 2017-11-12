//
// Date: 11/4/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"net/http"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

//
// Do Routes
//
func (t *Controller) DoRoutes(router *gin.Engine) {

	// --------- API V1 sub-routes ----------- //

	apiV1 := router.Group("/api/v1")

	apiV1.Use(t.AuthMiddleware())
	{
		// Symbols
		apiV1.GET("/symbols", t.GetSymbols)

		// // Watchlists
		apiV1.GET("/watchlists", t.GetWatchlists)
		apiV1.POST("/watchlists", t.CreateWatchlist)
		apiV1.GET("/watchlists/:id", t.GetWatchlist)
	}

	// ------- Websockets --------- //

	// Setup websocket - Core
	router.GET("/ws/core", func(c *gin.Context) {
		handler := http.HandlerFunc(t.DoWebsocketConnection)
		handler.ServeHTTP(c.Writer, c.Request)
	})

	// Setup websocket - Quotes
	router.GET("/ws/quotes", func(c *gin.Context) {
		handler := http.HandlerFunc(t.DoQuoteWebsocketConnection)
		handler.ServeHTTP(c.Writer, c.Request)
	})

	// -------- Static Files ------------ //

	router.Use(static.Serve("/", static.LocalFile("/frontend", true)))
	router.Use(static.Serve("/login", static.LocalFile("/frontend", true)))
	router.Use(static.Serve("/screener", static.LocalFile("/frontend", true)))

	// router.StaticFile("/", "/frontend/index.html")
	// router.Static("/assets", "/frontend/assets")
	// router.Static("*.css", "/frontend")

	// // Auth Routes
	// r.Handle("/login", t.CorsMiddleware(http.HandlerFunc(t.DoLogin))).Methods("POST", "OPTIONS")
	// r.Handle("/register", t.CorsMiddleware(http.HandlerFunc(t.DoRegister))).Methods("POST", "OPTIONS")
	// r.Handle("/reset-password", t.CorsMiddleware(http.HandlerFunc(t.DoResetPassword))).Methods("POST", "OPTIONS")
	// r.Handle("/forgot-password", t.CorsMiddleware(http.HandlerFunc(t.DoForgotPassword))).Methods("POST", "OPTIONS")

	// // Webhooks
	// r.HandleFunc("/webhooks/stripe", t.DoStripeWebhook).Methods("POST")

	// // Tradier Oauth
	// tr := &tradier.TradierAuth{DB: t.DB}
	// r.HandleFunc("/tradier/authorize", tr.DoAuthCode).Methods("GET")
	// r.HandleFunc("/tradier/callback", tr.DoAuthCallback).Methods("GET")

	// // Setup websocket
	// r.HandleFunc("/ws/core", t.DoWebsocketConnection)
	// r.HandleFunc("/ws/quotes", t.DoQuoteWebsocketConnection)

	// // Static files.
	// fs := http.FileServer(http.Dir("/frontend/"))
	// r.PathPrefix("/").Handler(fs)
	// r.Handle("/screener/", http.StripPrefix("/screener/", fs))

	// // ------- End API V1 sub-routes --------- //

	// // Auth Routes
	// r.Handle("/login", t.CorsMiddleware(http.HandlerFunc(t.DoLogin))).Methods("POST", "OPTIONS")
	// r.Handle("/register", t.CorsMiddleware(http.HandlerFunc(t.DoRegister))).Methods("POST", "OPTIONS")
	// r.Handle("/reset-password", t.CorsMiddleware(http.HandlerFunc(t.DoResetPassword))).Methods("POST", "OPTIONS")
	// r.Handle("/forgot-password", t.CorsMiddleware(http.HandlerFunc(t.DoForgotPassword))).Methods("POST", "OPTIONS")

	// // Webhooks
	// r.HandleFunc("/webhooks/stripe", t.DoStripeWebhook).Methods("POST")

	// // Tradier Oauth
	// tr := &tradier.TradierAuth{DB: t.DB}
	// r.HandleFunc("/tradier/authorize", tr.DoAuthCode).Methods("GET")
	// r.HandleFunc("/tradier/callback", tr.DoAuthCallback).Methods("GET")

	// // Setup websocket
	// r.HandleFunc("/ws/core", t.DoWebsocketConnection)
	// r.HandleFunc("/ws/quotes", t.DoQuoteWebsocketConnection)

	// // Static files.
	// fs := http.FileServer(http.Dir("/frontend/"))
	// r.PathPrefix("/").Handler(fs)
	// r.Handle("/screener/", http.StripPrefix("/screener/", fs))

}

/* End File */
