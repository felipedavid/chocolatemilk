package chocolatemilk

import (
	"bytes"
	"errors"
	"html/template"
	"net/http"
	"path/filepath"
)

var ErrUnknownTemplate = errors.New("unknown template")

type templateData struct {
	Data any
}

var functions = template.FuncMap{}

func newTemplateCache() (map[string]*template.Template, error) {
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

func (app *App) render(w http.ResponseWriter, status int, templateName string, data *templateData) error {
	ts, ok := app.templateCache[templateName]
	if !ok {
		return ErrUnknownTemplate
	}

	buf := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		return err
	}

	w.WriteHeader(status)
	buf.WriteTo(w)
	return nil
}
