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
		// Ping
		apiV1.GET("/ping", t.PingFromServer)

		// Settings
		apiV1.GET("/settings", t.GetSettings)

		// Orders
		apiV1.POST("/orders", t.SubmitOrder)
		apiV1.POST("/orders/preview", t.PreviewOrder)

		// Brokers
		apiV1.GET("/brokers", t.GetBrokers)
		apiV1.POST("/brokers", t.CreateBroker)
		apiV1.GET("/brokers/:id/balances", t.GetBalances)
		apiV1.GET("/brokers/:id/orders", t.GetBrokerActiveOrders)
		apiV1.PUT("/brokers/:id/accounts/:acctId", t.UpdateBrokerAccount)
		apiV1.GET("/brokers/:id/accounts/:acctId/balance", t.BrokerAccountGetBalance)

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

		// Broker events.
		apiV1.GET("/broker-events/:brokerAccount", t.GetBrokerEvents)

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
		apiV1.DELETE("/screeners/:id", t.DeleteScreener)
		apiV1.GET("/screeners/:id", t.GetScreener)
		apiV1.GET("/screeners/:id/results", t.GetScreenerResults)
		apiV1.POST("/screeners/results", t.GetScreenerResultsFromFilters)

		// Reports
		apiV1.GET("/reports/:brokerAccount/tradegroup/years", t.ReportsGetTradeGroupYears)
		apiV1.GET("/reports/:brokerAccount/summary/yearly/:year", t.ReportsGetAccountYearlySummary)

		// Status
		apiV1.GET("/status/market", t.GetMarketStatus)

		// Notifications
		apiV1.POST("/notifications/add-channel", t.CreateNotifyChannel)
	}

	// ---------- Websockets -------------- //

	r.GET("/ws", t.WebsocketController.DoWebsocketConnection)

	// ------------ Non-Auth Routes ------ //

	// Google login Routes
	r.GET("/oauth/google", t.DoGoogleAuthLogin)
	r.GET("/oauth/google/callback", t.DoGoogleCallback)
	r.POST("/oauth/google/session", t.DoStartGoogleLoginSession)
	r.POST("/oauth/google/token", t.DoGetAccessTokenAfterGoogleAuth)

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

	// Redirect after oauth to start the broker feed.
	r.GET("/broker-feed/start", t.DoStartBrokerFeed)

	// -------- Static Files ------------ //

	r.Use(static.Serve("/", static.LocalFile("/frontend", true)))
	r.NoRoute(func(c *gin.Context) { c.File("/frontend/index.html") })
}

/* End File */
