package main

import (
	"net/http"

	"github.com/justinas/nosurf"
	"github.com/mrkouhadi/go-booking-app/internal/helpers"
)

// csrf : ignore any POST request that doesn't have CSRF token
func Nosurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

// sessions loader
func LoadSession(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

// auth
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !helpers.IsAuthenticated(r) {
			session.Put(r.Context(), "error", "Please Log in first !")
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}
