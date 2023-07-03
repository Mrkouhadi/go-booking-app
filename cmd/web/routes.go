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
	mux.Get("/reservation-summary", handlers.Repo.ReservationSummary)
	mux.Get("/choose-room/{id}", handlers.Repo.ChooseRoom)
	mux.Get("/book-room", handlers.Repo.BookRoom)
	mux.Get("/user/login", handlers.Repo.ShowLogin)
	mux.Get("/user/logout", handlers.Repo.Logout)

	// POST methods
	mux.Post("/search-availability", handlers.Repo.PostSearchAvailability)
	mux.Post("/search-availability-json", handlers.Repo.AvailabilityJSON)
	mux.Post("/make-reservation", handlers.Repo.PostMakeReservation)
	mux.Post("/reservation-summary", handlers.Repo.ReservationSummary)
	mux.Post("/user/login", handlers.Repo.PostShowLogin)

	// render files in the template(html)
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	// protected routes
	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(Auth)

		mux.Get("/dashboard", handlers.Repo.AminDashboard)                         // route will be : admin/dashboard
		mux.Get("/reservations-new", handlers.Repo.AdminNewReservations)           // admin/reservations-new
		mux.Get("/reservations-all", handlers.Repo.AdminAllReservations)           // admin/reservations-all
		mux.Get("/reservations-calendar", handlers.Repo.AdminReservationsCalendar) // admin/reservations-calendar

	})
	return mux
}
