package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/mrkouhadi/go-booking-app/internal/config"
	"github.com/mrkouhadi/go-booking-app/internal/driver"
	"github.com/mrkouhadi/go-booking-app/internal/handlers"
	"github.com/mrkouhadi/go-booking-app/internal/helpers"
	"github.com/mrkouhadi/go-booking-app/internal/models"
	"github.com/mrkouhadi/go-booking-app/internal/render"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}

	defer db.SQL.Close()
	defer close(app.MailChan)
	fmt.Println("Started mail listener ")
	listenForMail()

	fmt.Println("Listening to Port:8080")
	srv := &http.Server{
		Addr:    portNumber,
		Handler: Routes(&app),
	}
	err = srv.ListenAndServe()

	log.Fatal(err)
}

func run() (*driver.DB, error) {
	// what am I going to store in the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(map[string]int{})

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime) // \t means tab (bunch of spaces)
	app.InfoLog = infoLog
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session
	// connect to the database
	log.Println("CNNECTING TO A DATABASE...")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=kouhadi password=")
	if err != nil {
		log.Fatal("Cannot connect to the database ! Dying...")
	}
	log.Println("SUCCESSFULLY CONNECTed TO A DATABASE !")

	// create templates cache
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot Create Template Cache。")
		return nil, err
	}
	app.TemplateCache = tc

	app.UseCache = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.Newhelpers(&app)
	return db, nil
}
