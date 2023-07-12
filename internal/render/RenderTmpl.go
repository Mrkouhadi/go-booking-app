package render

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"text/template"
	"time"

	"github.com/justinas/nosurf"
	"github.com/mrkouhadi/go-booking-app/internal/config"
	"github.com/mrkouhadi/go-booking-app/internal/models"
)

var functions = template.FuncMap{
	"HumanDate": HumanDate, // we just passed a function to our templates
}
var pathToTemplates = "./templates"

var app *config.AppConfig

// newTemplates sets the config for the template package
func NewRenderer(a *config.AppConfig) {
	app = a
}

// format dates to YYYY-MMM-DD
func HumanDate(t time.Time) string {
	layout := "2006-01-02"
	return t.Format(layout)
}

// add default data for the first rendered pages
func AddDefaultData(tmplData *models.TemplateData, r *http.Request) *models.TemplateData {
	// PopString puts somethign in session untill next page displayed and then taken out automatically
	tmplData.Flash = app.Session.PopString(r.Context(), "flash")
	tmplData.Warning = app.Session.PopString(r.Context(), "warning")
	tmplData.Error = app.Session.PopString(r.Context(), "error")
	tmplData.CSRFToken = nosurf.Token(r)
	// check the user if is authenticated
	if app.Session.Exists(r.Context(), "user_id") {
		tmplData.IsAuthenticated = 1 // default value is 0 BUT if user authenticated then it'll be 1
	}
	return tmplData
}

// RenderTemplate renders the requested template
func Template(w http.ResponseWriter, r *http.Request, tmpl string, tmplData *models.TemplateData) error {
	// Get the template cache from the AppConfig
	var tmplCache map[string]*template.Template
	if app.UseCache {
		tmplCache = app.TemplateCache
	} else {
		tmplCache, _ = CreateTemplateCache()
	}

	// get requested template from cached templates
	t, ok := tmplCache[tmpl]
	if !ok {
		return errors.New("could not get the template from cached templates")
	}
	buf := new(bytes.Buffer)
	tmplData = AddDefaultData(tmplData, r)
	err := t.Execute(buf, tmplData)
	if err != nil {
		log.Fatal(err)
	}

	// render template
	_, err = buf.WriteTo(w)
	if err != nil {
		fmt.Println("error wrinting template to the browser")
		return err
	}
	return nil
}

// CreateTemplateCache create cache for templates
func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// get all files named *.page.tmpl from ./templates
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		return myCache, err
	}

	// range through all files ending with *.page.tmpl
	for _, page := range pages {
		fileName := filepath.Base(page) // filepath.Base returns the last element of the path
		templSet, err := template.New(fileName).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}
		// look for any layout that exist in that directory
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return myCache, err
		}
		if len(matches) > 0 {
			templSet, err = templSet.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}

		myCache[fileName] = templSet
	}
	return myCache, nil
}
