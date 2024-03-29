//
// Date: 11/4/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"app.options.cafe/controllers/admin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/acme/autocert"
)

//
// Start the webserver
//
func (t *Controller) StartWebServer() {

	// Listen for data from our broker feeds.
	go t.WebsocketController.DoWsDispatch()

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

	// Start admin controllers
	a := &admin.Controller{DB: t.DB}

	// Do admin routes
	a.DoRoutes(router)

	// Are we in testing mode? If not give us some SSL
	if os.Getenv("APP_ENV") == "local" {

		s := &http.Server{
			Addr:         ":7080",
			Handler:      router,
			ReadTimeout:  120 * time.Second,
			WriteTimeout: 120 * time.Second,
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

		// Setup TLS config
		t := m.TLSConfig()
		t.ClientSessionCache = tls.NewLRUClientSessionCache(0)
		t.MinVersion = tls.VersionTLS11
		t.CurvePreferences = []tls.CurveID{
			tls.X25519,
			tls.CurveP521,
			tls.CurveP384,
			tls.CurveP256,
		}

		// Setup server
		s := &http.Server{
			Addr:         ":7043",
			ReadTimeout:  120 * time.Second,
			WriteTimeout: 120 * time.Second,
			IdleTimeout:  120 * time.Second,
			TLSConfig:    t,
			Handler:      router,
		}

		// Start http and https
		log.Printf("Starting secure server")
		go http.ListenAndServe(":7080", m.HTTPHandler(nil))
		log.Fatal(s.ListenAndServeTLS("", ""))
	}
}

/* End File */
