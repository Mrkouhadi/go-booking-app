package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/mrkouhadi/go-booking-app/internal/config"
	"github.com/mrkouhadi/go-booking-app/internal/driver"
	"github.com/mrkouhadi/go-booking-app/internal/forms"
	"github.com/mrkouhadi/go-booking-app/internal/helpers"
	"github.com/mrkouhadi/go-booking-app/internal/models"
	"github.com/mrkouhadi/go-booking-app/internal/render"
	"github.com/mrkouhadi/go-booking-app/internal/repository"
	"github.com/mrkouhadi/go-booking-app/internal/repository/dbrepo"
)

// the repository used by the handlers
var Repo *Repository

// repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepo creates the new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// /////// Home page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

// /////// About page
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "about.page.tmpl", &models.TemplateData{})
}

// /////// features page
func (m *Repository) Features(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "features.page.tmpl", &models.TemplateData{})
}

// /////// contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	m.DB.AllUsers()
	render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
}

// /////// make a reservation page
// GET
func (m *Repository) MakeReservation(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation
	render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil), // include an empty,
		Data: data,
	})
}

// POST
func (m *Repository) PostMakeReservation(w http.ResponseWriter, r *http.Request) {
	// it is recommende to do this part
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	roomId, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
		RoomId:    roomId,
		StartDate: startDate,
		EndDate:   endDate,
	}
	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email", "phone")
	form.MinLength("first_name", 3)
	form.IsEmail("email")
	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	// store reservation data in db
	newReservationId, err := m.DB.InsertReservation(reservation)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	restriction := models.RoomRestrictions{
		StartDate:     startDate,
		EndDate:       endDate,
		RoomId:        roomId,
		ReservationId: newReservationId,
		RestrictionId: 1,
	}
	// store restriction in db
	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	// store reservation details into session
	m.App.Session.Put(r.Context(), "reservation", reservation)
	// redirect the user to a different url after submitting the form
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

// /////// make a search-availability page
// GET
func (m *Repository) SearchAvailability(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// Post
func (m *Repository) PostSearchAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, start)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(layout, end)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	rooms, err := m.DB.SearchAvailablityForAllRooms(startDate, endDate)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	if len(rooms) == 0 {
		m.App.Session.Put(r.Context(), "error", " No available rooms for your Dates ")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}
	data := make(map[string]interface{})
	data["rooms"] = rooms

	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}

	m.App.Session.Put(r.Context(), "reservation", res)

	render.Template(w, r, "choose-room.page.tmpl", &models.TemplateData{
		Data: data,
	})
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
		helpers.ServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.Write(out)
}

// /////// general rooms page handler
func (m *Repository) GeneralRooms(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "generals.page.tmpl", &models.TemplateData{})
}

// /////// major suite rooms page handler
func (m *Repository) MajorSuite(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "majors.page.tmpl", &models.TemplateData{})
}

// ////////
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("cannot get reservation data from the session")
		m.App.Session.Put(r.Context(), "error", "Cannot get reservation data from the Session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	// remove data from reservation
	m.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// /// ChooseRoom method
func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	roomId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, err)
		return
	}
	res.RoomId = roomId
	m.App.Session.Put(r.Context(), "reservation", res)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}
