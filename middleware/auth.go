package middleware

import (
	"net/http"

	"github.com/robyzzz/csl-backend/config"
	"github.com/robyzzz/csl-backend/controller"
	"github.com/robyzzz/csl-backend/utils"
)

// Used to let client access only routes that need him to be logged in
func IsAuthenticated(h func(w http.ResponseWriter, r *http.Request)) http.Handler {
	next := http.HandlerFunc(h)
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Headers","Content-Type,access-control-allow-origin, access-control-allow-headers")

			if config.SessionAlreadyExists(r) {
				next.ServeHTTP(w, r)
			} else {
				utils.APIErrorRespond(w, utils.NewAPIError(http.StatusUnauthorized, "Not authenticated."))
			}
		})
}

//! not needed
// // Used to update steam user data when acessing /login
// // If user is already logged in, we update, else we redirect to login page
// func BeforeLogin(h func(w http.ResponseWriter, r *http.Request)) http.Handler {
// 	next := http.HandlerFunc(h)
// 	return http.HandlerFunc(
// 		func(w http.ResponseWriter, r *http.Request) {
// 			if config.SessionAlreadyExists(r) {
// 				controller.UpdateSteamUser(w, r)
// 				next.ServeHTTP(w, r)
// 			} else {
// 				next.ServeHTTP(w, r)
// 			}
// 		})
// }

func BeforeAuth(h func(w http.ResponseWriter, r *http.Request)) http.Handler {
	next := http.HandlerFunc(h)
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if config.SessionAlreadyExists(r) {
				controller.UpdateSteamUser(w, r)
				next.ServeHTTP(w, r)
			} else {
				utils.APIErrorRespond(w, utils.NewAPIError(http.StatusUnauthorized, "Unauthorized"))
			}
		})
}
