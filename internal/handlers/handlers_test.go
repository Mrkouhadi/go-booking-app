package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/mrkouhadi/go-booking-app/internal/models"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"about", "/about", "GET", []postData{}, http.StatusOK},
	{"generals-quarters", "/generals-quarters", "GET", []postData{}, http.StatusOK},
	{"majors-suite", "/majors-suite", "GET", []postData{}, http.StatusOK},
	{"search-availability", "/search-availability", "GET", []postData{}, http.StatusOK},
	{"contact", "/contact", "GET", []postData{}, http.StatusOK},
	// {"make-res", "/make-reservation", "GET", []postData{}, http.StatusOK},
	// {"search-availability", "/search-availability", "POST", []postData{
	// 	{key: "start", value: "2020-03-09"},
	// 	{key: "end", value: "2020-03-23"},
	// }, http.StatusOK},
	// {"search-availability-json", "/search-availability-json", "POST", []postData{
	// 	{key: "start", value: "2020-03-09"},
	// 	{key: "end", value: "2020-03-23"},
	// }, http.StatusOK},
	{"make-reservation", "/make-reservation", "POST", []postData{
		{key: "first_name", value: "Bakr"},
		{key: "last_name", value: "Kouhadi"},
		{key: "email", value: "Kouhadi@me.com"},
		{key: "phone", value: "15592081881"},
	}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()

	// a test server to run our server test
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		resp, err := ts.Client().Get(ts.URL + e.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}
		if resp.StatusCode != e.expectedStatusCode {
			t.Errorf("for %s expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
		}
	}
}

func TestRepository_Resrvation(t *testing.T) {
	reservation := models.Reservation{
		RoomId: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}

	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)
	handler := http.HandlerFunc(Repo.MakeReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("MakeReservation handler returns wrong response code. got %d, wanted %d", rr.Code, http.StatusOK)
	}

	// test case where reservation is not in the session (reset everything)
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("MakeReservation handler returns wrong response code. got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test with nonexistent room
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	reservation.RoomId = 100
	session.Put(ctx, "reservation", reservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("MakeReservation handler returns wrong response code. got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

}
func TestRepository_PostReservation(t *testing.T) {

	// reqBody := "start_date=2050-01-01"
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-10")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=Bryan")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "lasst_name=Kouhadi")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "email=Bryan@kouhadi.com")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=15592810818")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	postedData := url.Values{}
	postedData.Add("start_date", "2050-01-01")
	postedData.Add("end_date", "2050-01-10")
	postedData.Add("first_name", "Bryan")
	postedData.Add("last_name", "Kouhadi")
	postedData.Add("email", "Kouhadi@bryan.com")
	postedData.Add("phone", "15598789198")
	postedData.Add("room_id", "1")

	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	// tell the server that the request that you about to receive is of form POST
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.PostMakeReservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostMakeReservation handler returns wrong response code. got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test for missing post body

	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostMakeReservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostMakeReservation handler returns wrong response code for missing post body. got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// testing invalid START date
	postedData = url.Values{}
	postedData.Add("start_date", "invalid")
	postedData.Add("end_date", "2050-01-10")
	postedData.Add("first_name", "Bryan")
	postedData.Add("last_name", "Kouhadi")
	postedData.Add("email", "Kouhadi@bryan.com")
	postedData.Add("phone", "15598789198")
	postedData.Add("room_id", "1")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostMakeReservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostMakeReservation handler returns wrong response code for invalid start date. got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
	// testing invalid END date
	postedData = url.Values{}
	postedData.Add("start_date", "2050-01-01")
	postedData.Add("end_date", "invalid")
	postedData.Add("first_name", "Bryan")
	postedData.Add("last_name", "Kouhadi")
	postedData.Add("email", "Kouhadi@bryan.com")
	postedData.Add("phone", "15598789198")
	postedData.Add("room_id", "1")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostMakeReservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostMakeReservation handler returns wrong response code for invalid end date. got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// testing invalid room_id
	postedData = url.Values{}
	postedData.Add("start_date", "2050-01-01")
	postedData.Add("end_date", "2050-01-10")
	postedData.Add("first_name", "Bryan")
	postedData.Add("last_name", "Kouhadi")
	postedData.Add("email", "Kouhadi@bryan.com")
	postedData.Add("phone", "15598789198")
	postedData.Add("room_id", "invalid")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostMakeReservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostMakeReservation handler returns wrong response code for invalid room ID. got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
	// testing wrong data
	postedData = url.Values{}
	postedData.Add("start_date", "invalid")
	postedData.Add("end_date", "2050-01-10")
	postedData.Add("first_name", "B")
	postedData.Add("last_name", "Kouhadi")
	postedData.Add("email", "Kouhadi@bryan.com")
	postedData.Add("phone", "15598789198")
	postedData.Add("room_id", "1")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostMakeReservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostMakeReservation handler returns wrong response code for invalid DATA. got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// testing failure of inserting a reservation to the database
	postedData = url.Values{}
	postedData.Add("start_date", "invalid")
	postedData.Add("end_date", "2050-01-10")
	postedData.Add("first_name", "Bryan")
	postedData.Add("last_name", "Kouhadi")
	postedData.Add("email", "Kouhadi@bryan.com")
	postedData.Add("phone", "15598789198")
	postedData.Add("room_id", "2") // giving it 2 to generate the failure

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostMakeReservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostMakeReservation handler failed when trying to fail inserting reservation to DB. got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// testing failure of inserting a restriction to the database
	postedData = url.Values{}
	postedData.Add("start_date", "invalid")
	postedData.Add("end_date", "2050-01-10")
	postedData.Add("first_name", "Bryan")
	postedData.Add("last_name", "Kouhadi")
	postedData.Add("email", "Kouhadi@bryan.com")
	postedData.Add("phone", "15598789198")
	postedData.Add("room_id", "1000000") // giving an imposible number to generate failure

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostMakeReservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostMakeReservation handler failed when trying to fail inserting restriction to DB. got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}
func TestRepository_AvailabilityJSON(t *testing.T) {
	// first case: rooms are not available
	reqBody := "start=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2050-01-10")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")
	// create a new request
	req, _ := http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody))

	//  get context with session
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// make a handlerfunc
	handler := http.HandlerFunc(Repo.AvailabilityJSON)

	// make a response recorder
	rr := httptest.NewRecorder()
	// get request to our handler
	handler.ServeHTTP(rr, req)

	var js jsonResponse
	err := json.Unmarshal(rr.Body.Bytes(), &js) // rr.Body.Bytes() instead of []byte(rr.Body.String())
	if err != nil {
		t.Error("failed to parse json")
	}
}
func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("x-session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}
