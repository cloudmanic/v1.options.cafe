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
		// Orders
		apiV1.POST("/orders", t.SubmitOrder)
		apiV1.POST("/orders/preview", t.PreviewOrder)

		// Brokers
		apiV1.GET("/brokers", t.GetBrokers)
		apiV1.GET("/brokers/:id/balances", t.GetBalances)
		apiV1.GET("/brokers/:id/orders", t.GetBrokerActiveOrders)

		// Symbols
		apiV1.GET("/symbols", t.GetSymbols)
		apiV1.GET("/symbols/:symb", t.GetSymbol)
		apiV1.POST("/symbols/add-active-symbol", t.AddActiveSymbol)
		apiV1.POST("/symbols/get-option-symbol-from-parts", t.GetOptionSymbolFromParts)

		// Watchlists
		apiV1.GET("/watchlists", t.GetWatchlists)
		apiV1.POST("/watchlists", t.CreateWatchlist)
		apiV1.GET("/watchlists/:id", t.GetWatchlist)
		apiV1.PUT("/watchlists/:id", t.UpdateWatchlist)
		apiV1.DELETE("/watchlists/:id", t.DeleteWatchlist)
		apiV1.POST("/watchlists/:id/symbol", t.WatchlistAddSymbol)
		apiV1.DELETE("/watchlists/:id/symbol/:symb", t.WatchlistDeleteSymbol)
		apiV1.PUT("/watchlists/:id/reorder", t.WatchlistReorder)

		// Trade Groups
		apiV1.GET("/tradegroups", t.GetTradeGroups)

		// Quotes
		apiV1.GET("/quotes/historical", t.GetHistoricalQuotes)
		apiV1.GET("/quotes/rank/:symb", t.GetRank)
		apiV1.GET("/quotes/options/expirations/:symb", t.GetOptionsExpirations)
		apiV1.GET("/quotes/options/chain/:symb/:expire", t.GetOptionsChainByExpiration)
		apiV1.GET("/quotes/options/strikes/:symb/:expire", t.GetOptionsStikesBySymbolExpiration)

		// Screeners
		apiV1.GET("/screeners", t.GetScreeners)
		apiV1.POST("/screeners", t.CreateScreener)
		apiV1.PUT("/screeners/:id", t.UpdateScreener)
		apiV1.GET("/screeners/:id", t.GetScreener)
		apiV1.GET("/screeners/:id/results", t.GetScreenerResults)

		// Status
		apiV1.GET("/status/market", t.GetMarketStatus)
	}

	// ---------- Websockets -------------- //

	r.GET("/ws", t.WebsocketController.DoWebsocketConnection)

	// ------------ Non-Auth Routes ------ //

	// oAuth Routes
	r.POST("/oauth/token", t.DoOauthToken)
	r.GET("/oauth/logout", t.DoLogOut)

	// Other Auth Routes
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
