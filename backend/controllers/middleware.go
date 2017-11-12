//
// Date: 11/11/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"app.options.cafe/backend/library/services"
)

//
// Cors middleware for local development.
//
func (t *Controller) CorsMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Manage OPTIONS requests
		if (os.Getenv("APP_ENV") == "local") && (r.Method == http.MethodOptions) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization,DNT,X-CustomHeader,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Content-Range,Range")
			return
		}

		// On to next request in the Middleware chain.
		next.ServeHTTP(w, r)
	})
}

//
// Here we make sure we passed in a proper Bearer Access Token.
//
func (t *Controller) AuthMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Manage OPTIONS requests
		if (os.Getenv("APP_ENV") == "local") && (r.Method == http.MethodOptions) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization,DNT,X-CustomHeader,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Content-Range,Range")
			return
		}

		// Set access token and start the auth process
		var access_token = ""

		// Make sure we have a Bearer token.
		auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

		if len(auth) != 2 || auth[0] != "Bearer" {

			// We allow access token from the command line
			if os.Getenv("APP_ENV") == "local" {

				access_token = r.URL.Query().Get("access_token")

				if len(access_token) <= 0 {
					t.RespondError(w, http.StatusUnauthorized, "Authorization Failed (#101)")
					return
				}

			} else {
				t.RespondError(w, http.StatusUnauthorized, "Authorization Failed (#001)")
				return
			}

		} else {
			access_token = auth[1]
		}

		// See if this session is in our db.
		session, err := t.DB.GetByAccessToken(access_token)

		if err != nil {
			services.MajorLog("Access Token Not Found - Unable to Authenticate via HTTP (#002)")
			t.RespondError(w, http.StatusUnauthorized, "Authorization Failed (#002)")
			return
		}

		// Get this user is in our db.
		user, err := t.DB.GetUserById(session.UserId)

		if err != nil {
			services.MajorLog("User Not Found - Unable to Authenticate - UserId (HTTP) : " + fmt.Sprint(session.UserId) + " - Session Id : " + fmt.Sprint(session.Id))
			t.RespondError(w, http.StatusUnauthorized, "Authorization Failed (#003)")
			return
		}

		// Add this user to the context
		ctx := context.WithValue(r.Context(), "userId", user.Id)

		// CORS for local development.
		if os.Getenv("APP_ENV") == "local" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization,DNT,X-CustomHeader,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Content-Range,Range")
		}

		// On to next request in the Middleware chain.
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

//
// Get User Id from context
//
func (t *Controller) GetUserIdFromContext(r *http.Request) uint {
	return r.Context().Value("userId").(uint)
}

/* End File */
