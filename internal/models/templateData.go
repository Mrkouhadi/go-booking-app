package models

import "github.com/mrkouhadi/go-booking-app/internal/forms"

// TempleData holds data to be sent from handlers to templates
type TemplateData struct {
	StringMap       map[string]string
	IntMap          map[string]int
	FloatMap        map[string]float32
	Data            map[string]interface{}
	CSRFToken       string
	Flash           string
	Warning         string
	Error           string
	Form            *forms.Form
	IsAuthenticated int // logged in => > 0, Not logged in == 0
}
