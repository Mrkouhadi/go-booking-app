package handlers

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
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
		t.Errorf("Reservation handler returns wrong response code. got %d, wanted %d", rr.Code, http.StatusOK)
	}

	// test case where reservation is not in the session (reset everything)
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returns wrong response code. got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
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
		t.Errorf("Reservation handler returns wrong response code. got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

}

func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("x-session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}
