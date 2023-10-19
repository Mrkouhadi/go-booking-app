package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
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

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
	room, err := m.DB.GetRoomById(res.RoomId)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't find room")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
	res.Room.RoomName = room.RoomName
	m.App.Session.Put(r.Context(), "reservation", res)

	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")
	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]interface{})
	data["reservation"] = res
	render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
}

// POST
func (m *Repository) PostMakeReservation(w http.ResponseWriter, r *http.Request) {
	// it is recommend to do this part
	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse form!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")

	// 2020-01-01 -- 01/02 03:04:05PM '06 -0700

	layout := "2006-01-02"

	startDate, err := time.Parse(layout, sd)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse start date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	endDate, err := time.Parse(layout, ed)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't get parse end date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "invalid data!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	room, err := m.DB.GetRoomById(roomID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "invalid data!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Phone:     r.Form.Get("phone"),
		Email:     r.Form.Get("email"),
		StartDate: startDate,
		EndDate:   endDate,
		RoomId:    roomID,
		Room:      room,
	}

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email", "phone")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		http.Error(w, "my own error message", http.StatusSeeOther)
		render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	// store reservation data in db
	newReservationId, err := m.DB.InsertReservation(reservation)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't insert reservation into database")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	restriction := models.RoomRestrictions{
		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
		RoomId:        reservation.RoomId,
		ReservationId: newReservationId,
		RestrictionId: 1,
	}
	// store restriction in db
	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't insert room restriction")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	//  send a notification (email) to a customer
	mailContent := fmt.Sprintf(`
		<strong>Your reservation Confirmation</strong> <br>
		Dear %s, <br>
			Your Reservation for romm " %s" from %s to %s has been done succefully <br>
		Regards  <br>
	`, reservation.FirstName, room.RoomName, reservation.StartDate.Format("2006-01-02"), reservation.EndDate.Format("2006-01-02"))

	msg := models.MailData{
		To:       reservation.Email,
		From:     "kouhadibakr@gmail.com",
		Subject:  "Confirming your reservation",
		Content:  mailContent,
		Template: "basic.html",
	}
	m.App.MailChan <- msg

	//  send a notification (email)
	mailContent = fmt.Sprintf(`
			<strong> A Reservation Notification</strong> <br>
			Hi, <br>
				A Reservation for the room " %s " from %s to %s  has been made succefully by %s %s <br>
			Regards  <br>
		`, room.RoomName, reservation.StartDate.Format("2006-01-02"), reservation.EndDate.Format("2006-01-02"), reservation.FirstName, reservation.LastName)

	msg = models.MailData{
		To:       "kouhadibakr@gmail.com",
		From:     "contact@bookingapp.com",
		Subject:  "Reservation Notification",
		Content:  mailContent,
		Template: "basic.html",
	}
	m.App.MailChan <- msg

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
	OK         bool   `json:"ok"`
	MESSAGE    string `json:"message"`
	RoomId     string `json:"room_id"`
	Start_date string `json:"start_date"`
	End_date   string `json:"end_date"`
}

func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		resp := jsonResponse{
			OK:      false,
			MESSAGE: "Internal server error",
		}
		out, _ := json.MarshalIndent(resp, "", "        ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}
	sd := r.Form.Get("start")
	ed := r.Form.Get("end")
	// convert it from a string to a time.Time format
	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)
	roomId, _ := strconv.Atoi(r.Form.Get("room_id"))
	// get available room
	available, err := m.DB.SearchAvailabilityByDatesByRoomID(startDate, endDate, roomId)
	if err != nil {
		resp := jsonResponse{
			OK:      false,
			MESSAGE: "Error connecting to the database",
		}
		out, _ := json.MarshalIndent(resp, "", "        ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}
	res := jsonResponse{
		OK:         available,
		MESSAGE:    "", // not gonna use this field anyway
		Start_date: sd,
		End_date:   ed,
		RoomId:     strconv.Itoa(roomId),
	}

	out, _ := json.MarshalIndent(res, "", "		")

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

// //////// ReservationSummary shows the summary of the newly made reservation.
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
	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")
	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
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

// BookRoom get url params, builds sessional variable, and take user to make res screen
func (m *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {
	// id , s , e
	// get room's id and convert it from a string to an int
	ID, _ := strconv.Atoi(r.URL.Query().Get("id"))
	sd := r.URL.Query().Get("s")
	ed := r.URL.Query().Get("e")

	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)

	var res models.Reservation

	room, err := m.DB.GetRoomById(ID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	res.Room.RoomName = room.RoomName
	res.RoomId = ID
	res.StartDate = startDate
	res.EndDate = endDate

	m.App.Session.Put(r.Context(), "reservation", res)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

// ////////  AUTHENTICATION (login)
// GET
func (m *Repository) ShowLogin(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "login.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

// POST
func (m *Repository) PostShowLogin(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context())
	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "not able to parse the form")
		http.Redirect(w, r, "/admin", http.StatusTemporaryRedirect)
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	// check the validity of the form inputs
	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.IsEmail(email)
	form.MinLength("password", 4)
	///************************************* Here is the probelm. when this part is commented, it works fine
	// email: admin@kouhadi.com
	// password:pass12

	// if !form.Valid() {
	// 	render.Template(w, r, "login.page.tmpl", &models.TemplateData{
	// 		Form: form,
	// 	})
	// 	return
	// }
	//*************************************/
	// check the credentials
	id, _, err := m.DB.Authenticate(email, password)
	if err != nil {
		log.Println(err)
		m.App.Session.Put(r.Context(), "error", "Invalid Login Credentials !")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}
	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "flash", "Logged in successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ///// Logout
func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context()) // destroy all data in the session
	_ = m.App.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

// /////////////////////////// admindashboard (PROTECTED ROUTEs)
// Admin Dashboard
func (m *Repository) AminDashboard(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-dashboard.page.tmpl", &models.TemplateData{})
}

// shows al new reservations in admin tool
func (m *Repository) AdminNewReservations(w http.ResponseWriter, r *http.Request) {
	reservations, err := m.DB.AllNewReservations()

	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["reservations"] = reservations
	render.Template(w, r, "admin-new-reservations.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// all reservations
func (m *Repository) AdminAllReservations(w http.ResponseWriter, r *http.Request) {
	reservations, err := m.DB.AllReservations()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["reservations"] = reservations

	render.Template(w, r, "admin-all-reservations.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// admin show a single reservation by id
func (m *Repository) AdminShowReservation(w http.ResponseWriter, r *http.Request) {
	// "/admin/reservations/all/5"  so id is i position 4 and src in position 3
	exploded := strings.Split(r.RequestURI, "/")
	ID, err := strconv.Atoi(exploded[4])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	src := exploded[3]
	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")
	strmap := make(map[string]string)
	strmap["src"] = src
	strmap["month"] = month
	strmap["year"] = year

	res, err := m.DB.GetReservationByID(ID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	datamap := make(map[string]interface{})
	datamap["reservation"] = res
	render.Template(w, r, "admin-reservations-show.page.tmpl", &models.TemplateData{
		StringMap: strmap,
		Data:      datamap,
		Form:      forms.New(nil),
	})
}

// admin POST a single reservation by id
func (m *Repository) AdminPostShowReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		// because is just the admin that's wy we didn't use the error handling of parsing form like in other handlers
		helpers.ServerError(w, err)
		return
	}
	exploded := strings.Split(r.RequestURI, "/")
	ID, err := strconv.Atoi(exploded[4])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	src := exploded[3]
	strmap := make(map[string]string)
	strmap["src"] = src
	res, err := m.DB.GetReservationByID(ID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	// TODO: ERRORS CHECKING OF THE FORM HERE, AND ALSO ON THE CLIENT SIDE
	res.FirstName = r.Form.Get("first_name")
	res.LastName = r.Form.Get("last_name")
	res.Email = r.Form.Get("email")
	res.Phone = r.Form.Get("phone")
	error := m.DB.UpdateReservation(res)
	if error != nil {
		helpers.ServerError(w, err)
		return
	}
	month := r.Form.Get("month")
	year := r.Form.Get("year")

	m.App.Session.Put(r.Context(), "flash", "changes have been saved succesfully")
	if year == "" {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?y=%s&m=%s", year, month), http.StatusSeeOther)
	}
}

// AdminProcessReservation markes the reservation as processed
func (m *Repository) AdminProcessReservation(w http.ResponseWriter, r *http.Request) {
	ID, _ := strconv.Atoi(chi.URLParam(r, "id"))
	//src := chi.URLParam(r, "src")

	error := m.DB.UpdateProcessedForReservation(ID, 1)
	if error != nil {
		helpers.ServerError(w, error)
		return
	}
	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")

	m.App.Session.Put(r.Context(), "flash", "Reservation has been marked as processed succesfully")
	if year == "" {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-#{src}"), http.StatusSeeOther)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?y=%s&m=%s", year, month), http.StatusSeeOther)
	}
}

// AdminDeleteReservation deletes a reservation
func (m *Repository) AdminDeleteReservation(w http.ResponseWriter, r *http.Request) {
	ID, _ := strconv.Atoi(chi.URLParam(r, "id"))
	//src := chi.URLParam(r, "src")

	error := m.DB.DeleteReservation(ID)
	if error != nil {
		helpers.ServerError(w, error)
		return
	}

	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")

	m.App.Session.Put(r.Context(), "flash", "Reservation has been deleted succesfully")
	if year == "" {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-#{src}"), http.StatusSeeOther)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?y=%s&m=%s", year, month), http.StatusSeeOther)
	}
}

// reservations calendar
func (m *Repository) AdminReservationsCalendar(w http.ResponseWriter, r *http.Request) {
	// assume that there is no mont/year specified
	now := time.Now()
	if r.URL.Query().Get("y") != "" {
		year, _ := strconv.Atoi(r.URL.Query().Get("y"))
		month, _ := strconv.Atoi(r.URL.Query().Get("m"))
		now = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	}

	next := now.AddDate(0, 1, 0) // years,months,days
	last := now.AddDate(0, -1, 0)

	nextMonth := next.Format("01") // get 2 digits month
	nextMonthYear := next.Format("2006")

	lastMonth := last.Format("01")
	lastMonthYear := last.Format("2006")

	stringMap := make(map[string]string)

	stringMap["next_month"] = nextMonth
	stringMap["next_month_year"] = nextMonthYear

	stringMap["last_month"] = lastMonth
	stringMap["last_month_year"] = lastMonthYear

	stringMap["this_month"] = now.Format("01")
	stringMap["this_month_year"] = now.Format("2006")
	data := make(map[string]interface{})
	data["now"] = now
	currYear, currMonth, _ := now.Date()
	currLocation := now.Location()
	firstOfMonth := time.Date(currYear, currMonth, 1, 0, 0, 0, 0, currLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	intMap := make(map[string]int)
	intMap["days_in_month"] = lastOfMonth.Day()
	rooms, err := m.DB.AllRooms()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data["rooms"] = rooms
	for _, x := range rooms {
		// create maps
		reservationsMap := make(map[string]int)
		blockMap := make(map[string]int)
		for d := firstOfMonth; !d.After(lastOfMonth); d = d.AddDate(0, 0, 1) {
			reservationsMap[d.Format("2006-01-02")] = 0
			blockMap[d.Format("2006-01-02")] = 0
		}
		// get all restrictions
		restrictions, err := m.DB.GetRestrictionsFoorRoomByDate(x.ID, firstOfMonth, lastOfMonth)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
		for _, y := range restrictions {
			if y.ReservationId > 0 {
				// it's a reservation
				for d := y.StartDate; !d.After(y.EndDate); d = d.AddDate(0, 0, 1) {
					reservationsMap[d.Format("2006-01-02")] = y.ReservationId
				}
			} else {
				//  it's a block
				blockMap[y.StartDate.Format("2006-01-02")] = y.ID
			}
		}
		data[fmt.Sprintf("reservation_map_%d", x.ID)] = reservationsMap
		data[fmt.Sprintf("block_map_%d", x.ID)] = blockMap

		m.App.Session.Put(r.Context(), fmt.Sprintf("block_map_%d", x.ID), blockMap)
	}
	//  render temlates along with the data
	render.Template(w, r, "admin-reservations-calendar.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
		Data:      data,
		IntMap:    intMap,
	})
}
func (m *Repository) AdminPostReservationsCalendar(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	year, _ := strconv.Atoi(r.Form.Get("y"))
	month, _ := strconv.Atoi(r.Form.Get("m"))

	// get rooms
	rooms, err := m.DB.AllRooms()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	form := forms.New(r.PostForm)

	for _, rm := range rooms {
		// get the block map from the session. loop through the map and if we have the entry in the map
		// that doesn't exist in our posted data, and if the resteriction id > 0, then  it is a block
		// that we need to remove
		currMap := m.App.Session.Get(r.Context(), fmt.Sprintf("block_map_%d", rm.ID)).(map[string]int) // cast it from interface{} to map[string]int
		for name, value := range currMap {
			// ok will be false if the value is in the map
			if val, ok := currMap[name]; ok {
				// only pat atention to values > 0, and that are not in the form post
				// the rest are just placeholders for day without clocks
				if val > 0 {
					if !form.Has(fmt.Sprintf("remove_block_%d_%s", rm.ID, name), r) {
						// delete the restriction by id
						err := m.DB.DeleteBlockByID(value)
						if err != nil {
							log.Println(err)
						}
					}
				}
			}
		}
	}
	// now handle new blocks
	for name, _ := range r.PostForm {
		if strings.HasPrefix(name, "add_block") {
			exploded := strings.Split(name, "_")
			roomID, _ := strconv.Atoi(exploded[2])
			t, _ := time.Parse("2006-01-02", exploded[3])
			err := m.DB.InsertBlockForRoom(roomID, t)
			if err != nil {
				log.Println(err)
			}
		}
	}

	// store a success message in the session
	m.App.Session.Put(r.Context(), "flash", "Changes have been saved")
	// redirect the user after saving the changes
	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?y=%d&m=%d", year, month), http.StatusSeeOther)
}

// ////////////////////////// This is only for TESTing purposes
func NewTestRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewTestingRepo(a),
	}
}
