// Copyright 2017 Eduardo Pinheiro (edpin@edpin.com). All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websocket

import (
	"log"
	"net/http"
	"crypto/tls"
	"time"
)

// StartSecureServer starts an HTTPS server with a mux and a getCertificate
// function, which may be nil. Typically getCertificate will come from
// autocert.Manager, from package acme/autocert. The HTTPS server started
// enables HTST by default to ensure maximum protection (see
// https://www.owasp.org/index.php/HTTP_Strict_Transport_Security_Cheat_Sheet).
// StartSecureServer also starts an HTTP server that redirects all requests to
// their HTTPS counterpart and immediately terminates all connections.
func StartSecureServer(mux *http.ServeMux, getCertificate func(clientHello *tls.ClientHelloInfo) (*tls.Certificate, error)) {
	s := &http.Server{
		Addr:         ":https",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
    IdleTimeout:  120 * time.Second,
		Handler:      NewHSTS(mux),
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

type htstMux struct {
	*http.ServeMux
}

// NewHSTS returns an HTTP handler that sets HSTS headers on all requests.
func NewHSTS(mux *http.ServeMux) http.Handler {
	return htstMux{
		ServeMux: mux,
	}
}

// ServeHTTP implements http.Handler.
func (h htstMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Strict-Transport-Security", "max-age=86400; includeSubDomains")
	h.ServeMux.ServeHTTP(w, r)
}
