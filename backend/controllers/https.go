//
// Date: 11/9/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"crypto/tls"
	"log"
	"net/http"
	"time"
)

//
// StartSecureServer starts an HTTPS server with a mux and a getCertificate
// function, which may be nil. Typically getCertificate will come from
// autocert.Manager, from package acme/autocert. The HTTPS server started
// enables HTST by default to ensure maximum protection (see
// https://www.owasp.org/index.php/HTTP_Strict_Transport_Security_Cheat_Sheet).
// StartSecureServer also starts an HTTP server that redirects all requests to
// their HTTPS counterpart and immediately terminates all connections.
//
// (we do not use 80 & 443 as we run the docker container as non-root)
//
func StartSecureServer(r http.Handler, getCertificate func(clientHello *tls.ClientHelloInfo) (*tls.Certificate, error)) {
	s := &http.Server{
		Addr:         ":7043",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      MiddlewareHSTS(r),
		TLSConfig: &tls.Config{
			GetCertificate: getCertificate,
			MinVersion:     tls.VersionTLS12,
			CurvePreferences: []tls.CurveID{
				tls.X25519, // requires go 1.8
				tls.CurveP521,
				tls.CurveP384,
				tls.CurveP256,
			},
			// Prefer this order of ciphers.
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				// required by HTTP-2.
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			},
		},
	}

	// Redirect regular HTTP requests to HTTPS.
	insecure := &http.Server{
		Addr:         ":7080",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Connection", "close")
			url := "https://" + req.Host + req.URL.String()
			http.Redirect(w, req, url, http.StatusMovedPermanently)
		}),
	}

	go func() { log.Fatal(insecure.ListenAndServe()) }()

	log.Fatal(s.ListenAndServeTLS("", ""))
}

//
// ServeHTTP implements http.Handler.
//
func MiddlewareHSTS(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Strict-Transport-Security", "max-age=86400; includeSubDomains")

		// On to next request in the Middleware chain.
		next.ServeHTTP(w, r)
	})
}

/* End File */
