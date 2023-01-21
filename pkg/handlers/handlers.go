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

// Home is the Home page handler
func (m *Repository) Home(w http.ResponseWriter, req *http.Request) {
	// get IP address of the visitor
	remoteIP := req.RemoteAddr
	// write data in a session
	m.App.Session.Put(req.Context(), "remote_ip", remoteIP)
	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{})
}

// About is the About page handler
func (m *Repository) About(w http.ResponseWriter, req *http.Request) {
	strMap := make(map[string]string)
	strMap["test"] = "I am a piece of Data passed to the about page from about handler."
	remoteIp := m.App.Session.GetString(req.Context(), "remote_ip")
	strMap["remote_ip"] = remoteIp
	render.RenderTemplate(w, "about.page.tmpl", &models.TemplateData{StringMap: strMap})
}

func (m *Repository) Features(w http.ResponseWriter, req *http.Request) {
	render.RenderTemplate(w, "features.page.tmpl", &models.TemplateData{})
}
