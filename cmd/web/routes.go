package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mrkouhadi/go-booking-app/pkg/config"
	"github.com/mrkouhadi/go-booking-app/pkg/handlers"
)

func Routes(app *config.AppConfig) http.Handler {

	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)

	mux.Use(Nosurf)
	mux.Use(LoadSession)
	mux.Get("/", handlers.Repo.Home)
	mux.Get("/features", handlers.Repo.Features)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/contact", handlers.Repo.Contact) //
	mux.Get("/make-reservation", handlers.Repo.MakeReservation)
	mux.Get("/general-rooms", handlers.Repo.GeneralRooms)
	mux.Get("/major-suite", handlers.Repo.MajorSuite)

	// render files in the template(html)
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
