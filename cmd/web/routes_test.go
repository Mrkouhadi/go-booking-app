package main

import (
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/mrkouhadi/go-booking-app/internal/config"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	mux := Routes(&app)

	switch v := mux.(type) {
	case *chi.Mux:
		// do nothing
	default:
		t.Errorf("type is not *chi.Mux, but is %T", v)
	}
}
