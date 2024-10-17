package render

import (
	"bytes"
	"log"
	"net/http"
	"path/filepath"
	"text/template"
	"time"

	"github.com/justinas/nosurf"
	"github.com/tegveer-singh123/bookings/internal/config"
	"github.com/tegveer-singh123/bookings/internal/models"
)

var functions = template.FuncMap{
	"humanDate":  HumanDate,
	"formatDate": FormatDate,
	"iterate":    Iterate,
	"add":        Add,
}

var app *config.AppConfig

func NewRenderer(a *config.AppConfig) {
	app = a
}

func Iterate(count int) []int {
	var i int
	var items []int
	for i = 1; i <= count; i++ {
		items = append(items, i)
	}
	//log.Println("Iterating with count:", count, " Result slice size:", len(items))
	return items

}

// it returns date in YYYY-MM-DD format
func HumanDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func FormatDate(t time.Time, f string) string {
	return t.Format(f)
}

func Add(a, b int) int {
	return a + b
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	if app.Session.Exists(r.Context(), "user_id") {
		td.IsAuthenticated = 1
	}

	//sets the csrf token
	td.CSRFToken = nosurf.Token(r)

	//log.Println("CSRF token generated :", td.CSRFToken)
	return td
}

// renderTemplate renders the template
func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) {
	//create template cache

	var tc map[string]*template.Template
	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateNewTemplateCache()
	}

	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Error Parsing Template")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)

	err := t.Execute(buf, td)
	if err != nil {
		log.Fatal(err)
	}

	// parse the template
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Fatal(err)
	}

}

// it creates the template cache
func CreateNewTemplateCache() (map[string]*template.Template, error) {

	myCache := map[string]*template.Template{}

	//get all files named *page.tmpl from templates

	pages, err := filepath.Glob("../../templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("../../templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("../../templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts

	}

	return myCache, nil
}
