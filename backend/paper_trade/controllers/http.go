//
// Date: 2018-11-05
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-11-05
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/cloudmanic/app.options.cafe/backend/library/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

//
// Start the webserver
//
func (t *Controller) StartWebServer() {

	// Set GIN Settings
	gin.SetMode("release")
	gin.DisableConsoleColor()

	// Set Router
	router := gin.New()

	// Logger - Global middleware
	if os.Getenv("HTTP_LOG_REQUESTS") == "true" {
		router.Use(gin.Logger())
	}

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.Recovery())

	// CORS Middleware - Global middleware
	router.Use(cors.New(cors.Config{
		AllowMethods:     []string{"POST", "GET", "PUT", "OPTIONS", "DELETE"},
		AllowHeaders:     []string{"Authorization", "DNT", "X-CustomHeader", "Keep-Alive", "User-Agent", "X-Requested-With", "If-Modified-Since", "Cache-Control", "Content-Type", "Content-Range,Range"},
		ExposeHeaders:    []string{"Content-Length", "X-Last-Page", "X-Offset", "X-Limit", "X-No-Limit-Count"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return (origin == os.Getenv("SITE_URL")) || strings.Contains(origin, "localhost")
		},
		MaxAge: 12 * time.Hour,
	}))

	// Register Routes
	t.DoRoutes(router)

	// Setup HTTP Server
	s := &http.Server{
		Addr:         ":7081",
		Handler:      router,
		ReadTimeout:  120 * time.Second,
		WriteTimeout: 120 * time.Second,
	}

	// Lets get started
	services.InfoMsg("HTTP Server Started : " + os.Getenv("SITE_URL"))

	log.Fatal(s.ListenAndServe())

}

/* End File */
