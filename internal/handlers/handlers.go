package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/mrkouhadi/go-booking-app/internal/config"
	"github.com/mrkouhadi/go-booking-app/internal/forms"
	"github.com/mrkouhadi/go-booking-app/internal/models"
	"github.com/mrkouhadi/go-booking-app/internal/render"
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

// /////// Home page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	// get IP address of the visitor
	remoteIP := r.RemoteAddr
	// write data in a session
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)
	render.RenderTemplate(w, r, "home.page.tmpl", &models.TemplateData{})
}

// /////// About page
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	strMap := make(map[string]string)
	strMap["test"] = "I am a piece of Data passed to the about page from about handler."
	remoteIp := m.App.Session.GetString(r.Context(), "remote_ip")
	strMap["remote_ip"] = remoteIp
	render.RenderTemplate(w, r, "about.page.tmpl", &models.TemplateData{StringMap: strMap})
}

// /////// features page
func (m *Repository) Features(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "features.page.tmpl", &models.TemplateData{})
}

// /////// contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "contact.page.tmpl", &models.TemplateData{})
}

// /////// make a reservation page
// GET
func (m *Repository) MakeReservation(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation
	render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil), // include an empty,
		Data: data,
	})
}

// POST
func (m *Repository) PostMakeReservation(w http.ResponseWriter, r *http.Request) {
	// it is recommende to do this part
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}
	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}
	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email", "phone")
	form.MinLength("first_name", 3, r)

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
}

// /////// make a search-availability page
// GET
func (m *Repository) SearchAvailability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// Post
func (m *Repository) PostSearchAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")
	w.Write([]byte(fmt.Sprintf("start date is %s and end date is %s", start, end)))
}

// send data JSON back
type jsonResponse struct {
	OK      bool   `json:"ok"`
	MESSAGE string `json:"message"`
}

func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	res := jsonResponse{
		OK:      true,
		MESSAGE: "There is an available room!",
	}
	out, err := json.MarshalIndent(res, "", "		")
	if err != nil {
		log.Println(err)
	}

	// log.Println(r.Body)

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// /////// general rooms page handler
func (m *Repository) GeneralRooms(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "generals.page.tmpl", &models.TemplateData{})
}

// /////// major suite rooms page handler
func (m *Repository) MajorSuite(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "majors.page.tmpl", &models.TemplateData{})
}
