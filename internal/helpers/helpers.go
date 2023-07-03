package helpers

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/mrkouhadi/go-booking-app/internal/config"
)

var app *config.AppConfig

// Newhelpers sets up appconfig for helpers
func Newhelpers(a *config.AppConfig) {
	app = a
}

// ClientError sends errors to the client when something goes wrong on the clientside
func ClientError(w http.ResponseWriter, status int) {
	app.InfoLog.Println("Client error with status of ", status)
	http.Error(w, http.StatusText(status), status)
}

// ServerError sends errors to the client when something goes wrong on the serverside
func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func IsAuthenticated(r *http.Request) bool {
	exists := app.Session.Exists(r.Context(), "user_id")
	return exists
}
