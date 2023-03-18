package handlers

import (
	"net/http"

	"github.com/mrkouhadi/go-booking-app/pkg/config"
	"github.com/mrkouhadi/go-booking-app/pkg/models"
	"github.com/mrkouhadi/go-booking-app/pkg/render"
)

// the repository used by the handlers
var Repo *Repository

// repository type
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates the new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	// get IP address of the visitor
	remoteIP := r.RemoteAddr
	// write data in a session
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)
	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{})
}

// About page
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	strMap := make(map[string]string)
	strMap["test"] = "I am a piece of Data passed to the about page from about handler."
	remoteIp := m.App.Session.GetString(r.Context(), "remote_ip")
	strMap["remote_ip"] = remoteIp
	render.RenderTemplate(w, "about.page.tmpl", &models.TemplateData{StringMap: strMap})
}

// features page
func (m *Repository) Features(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "features.page.tmpl", &models.TemplateData{})
}

// contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "contact.page.tmpl", &models.TemplateData{})
}

// make a reservation page
func (m *Repository) MakeReservation(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "make-reservation.page.tmpl", &models.TemplateData{})
}

// general rooms page handler
func (m *Repository) GeneralRooms(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "generals.page.tmpl", &models.TemplateData{})
}

// major suite rooms page handler
func (m *Repository) MajorSuite(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "major-suite.page.tmpl", &models.TemplateData{})
}
