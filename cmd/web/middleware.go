package main

import (
	"net/http"

	"github.com/justinas/nosurf"
)

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
