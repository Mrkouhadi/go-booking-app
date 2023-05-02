package render

import (
	"encoding/gob"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/mrkouhadi/go-booking-app/internal/config"
	"github.com/mrkouhadi/go-booking-app/internal/models"
)

var session *scs.SessionManager
var testApp config.AppConfig

func TestMain(m *testing.M) {
	// what am I going to store in the session
	gob.Register(models.Reservation{})

	testApp.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false

	testApp.Session = session
	app = &testApp   // app declared in renderTmpl.go
	os.Exit(m.Run()) // before it exits, it runs M
}

// create a response writer that we need in TestRenderTemplate()
type mywriter struct{}

func (tw *mywriter) Header() http.Header {
	var h http.Header
	return h
}
func (tw *mywriter) Write(dt []byte) (int, error) {
	length := len(dt)
	return length, nil
}
func (tw *mywriter) WriteHeader(i int) {

}
