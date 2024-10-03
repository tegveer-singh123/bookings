package render

import (
	"bytes"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/tegveer-singh123/bookings/pkg/config"
	"github.com/tegveer-singh123/bookings/pkg/models"
)

var functions = template.FuncMap{}

var app *config.AppConfig

func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

// renderTemplate renders the template
func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	//create template cache

	var tc map[string]*template.Template
	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateNewTemplateCache()
	}

	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Error ")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td)

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

	pages, err := filepath.Glob("../../templates/*page.tmpl")
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("../../templates/*layout.tmpl")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("../../templates/*layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts

	}

	return myCache, nil
}
