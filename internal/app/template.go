package app

import (
	"bytes"
	"fmt"
	"forum/internal/models"
	"net/http"
	"path/filepath"
	"text/template"
	"time"
)

func NewTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}

func (app *Application) Render(w http.ResponseWriter, status int, page string, data *models.TemplateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.ServerError(w, err)
		return
	}

	buf := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.ServerError(w, err)
		return
	}
	w.WriteHeader(status)
	buf.WriteTo(w)
}

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}
