package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
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
	{"make-res", "/make-reservation", "GET", []postData{}, http.StatusOK},
	{"search-availability", "/search-availability", "POST", []postData{
		{key: "start", value: "2020-03-09"},
		{key: "end", value: "2020-03-23"},
	}, http.StatusOK},
	{"search-availability-json", "/search-availability-json", "POST", []postData{
		{key: "start", value: "2020-03-09"},
		{key: "end", value: "2020-03-23"},
	}, http.StatusOK},
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
		if e.method == "GET" {
			resp, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		} else {
			vals := url.Values{}
			for _, x := range e.params {
				vals.Add(x.key, x.value)
			}
			resp, err := ts.Client().PostForm(ts.URL+e.url, vals)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		}
	}
}
