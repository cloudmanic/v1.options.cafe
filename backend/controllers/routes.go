//
// Date: 11/4/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"github.com/cloudmanic/app.options.cafe/backend/brokers/tradier"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

//
// Do Routes
//
func (t *Controller) DoRoutes(r *gin.Engine) {

	// --------- API V1 sub-routes ----------- //

	apiV1 := r.Group("/api/v1")

	apiV1.Use(t.AuthMiddleware())
	{
		// Brokers
		apiV1.GET("/brokers", t.GetBrokers)

		// Symbols
		apiV1.GET("/symbols", t.GetSymbols)

		// Watchlists
		apiV1.GET("/watchlists", t.GetWatchlists)
		apiV1.POST("/watchlists", t.CreateWatchlist)
		apiV1.GET("/watchlists/:id", t.GetWatchlist)

		// Trade Groups
		apiV1.GET("/tradegroups", t.GetTradeGroups)

	}

	// ---------- Websockets -------------- //

	r.GET("/ws/core", t.DoWebsocketConnection)
	r.GET("/ws/quotes", t.DoQuoteWebsocketConnection)

	// ------------ Non-Auth Routes ------ //

	// // Auth Routes
	r.POST("/login", t.DoLogin)
	r.POST("/register", t.DoRegister)
	r.POST("/reset-password", t.DoResetPassword)
	r.POST("/forgot-password", t.DoForgotPassword)

	// Webhooks
	r.GET("/webhooks/stripe", t.DoStripeWebhook)

	// // Tradier Oauth
	tr := &tradier.TradierAuth{DB: t.DB}
	r.GET("/tradier/authorize", tr.DoAuthCode)
	r.GET("/tradier/callback", tr.DoAuthCallback)

	// -------- Static Files ------------ //

	r.Use(static.Serve("/", static.LocalFile("/frontend", true)))
	r.NoRoute(func(c *gin.Context) { c.File("/frontend/index.html") })
}

/* End File */
