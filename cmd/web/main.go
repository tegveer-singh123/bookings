package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/tegveer-singh123/bookings/pkg/config"
	"github.com/tegveer-singh123/bookings/pkg/handlers"
	"github.com/tegveer-singh123/bookings/pkg/render"
)

const PortNumber = ":8080"

func main() {

	var app config.AppConfig
     
    //change this to true when inproduction
	app.InProduction = false

	session := scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.Secure = app.InProduction

	tc, err := render.CreateNewTemplateCache()
	if err != nil {
		log.Fatal("Error creating template cache ")
	}

	app.TemplateCache = tc

	app.UseCache = false

	render.NewTemplates(&app)

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

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
