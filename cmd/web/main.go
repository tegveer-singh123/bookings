package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/tegveer-singh123/bookings/internal/config"
	"github.com/tegveer-singh123/bookings/internal/drivers"
	"github.com/tegveer-singh123/bookings/internal/handlers"
	"github.com/tegveer-singh123/bookings/internal/helpers"
	"github.com/tegveer-singh123/bookings/internal/models"
	"github.com/tegveer-singh123/bookings/internal/render"
)

const PortNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {

	db, err := Run()
	if err != nil {
		log.Fatal(err)
	}

	defer db.SQL.Close()

	defer close(app.Mailchan)

	listenForMail()

	
	fmt.Println(fmt.Sprintf("Starting Application on port %s", PortNumber))

	srv := &http.Server{
		Addr:    PortNumber,
		Handler: Routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}

func Run() (*drivers.DB, error) {
	//session storing data
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.RoomRestriction{})
	gob.Register(map[string]int{})

	mailChan := make(chan models.MailData)
	app.Mailchan = mailChan 

	//change this to true when inproduction
	app.InProduction = false

    session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.HttpOnly = true
	session.Cookie.Secure = app.InProduction

	app.Session = session

	app.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.ErrorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	log.Println("Connecting to database...")
	db, err := drivers.ConnectSQL("host=localhost port=5432 dbname=bookings user=tegi password=")
	if err != nil {
		log.Fatal("Cannot connect to the database")
	}

	log.Println("Connected to the database")

	tc, err := render.CreateNewTemplateCache()
	if err != nil {
		log.Fatal("Error creating template cache", err)
		return nil, err
	}

	app.TemplateCache = tc

	app.UseCache = false

	render.NewRenderer(&app)

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	helpers.NewHelpers(&app)

	return db, nil
}
