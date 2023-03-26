package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mrkouhadi/go-booking-app/internal/config"
	"github.com/mrkouhadi/go-booking-app/internal/handlers"
)

func Routes(app *config.AppConfig) http.Handler {

	mux := chi.NewRouter()

	// midlwares
	mux.Use(middleware.Recoverer)
	mux.Use(Nosurf) //  ignore any POST request that doesn't have CSRF token
	mux.Use(LoadSession)

	// GET methods
	mux.Get("/", handlers.Repo.Home)
	mux.Get("/features", handlers.Repo.Features)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/contact", handlers.Repo.Contact) //
	mux.Get("/make-reservation", handlers.Repo.MakeReservation)
	mux.Get("/generals-quarters", handlers.Repo.GeneralRooms)
	mux.Get("/majors-suite", handlers.Repo.MajorSuite)
	mux.Get("/search-availability", handlers.Repo.SearchAvailability)

	// POST methods
	mux.Post("/search-availability", handlers.Repo.PostSearchAvailability)
	mux.Post("/search-availability-json", handlers.Repo.AvailabilityJSON)

	// render files in the template(html)
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
