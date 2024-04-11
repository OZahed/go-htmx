package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

const (
	contentType = "Content-Type"
	textHtml    = "text/html"
)

type LayputInfo struct {
	DataMap  map[string]any
	PageName string
	Route    string
}

var (
	tmpl          *template.Template
	templateFuncs = template.FuncMap{
		"Contains": strings.Contains,
		"Dict": func(args ...string) (map[string]any, error) {
			if len(args)%2 != 0 {
				return nil, errors.New("not sufficient number of inputs")
			}

			m := make(map[string]any)

			for i := 0; i < len(args); i += 2 {
				m[args[i]] = args[i+1]
			}

			return m, nil
		},
	}
)

func init() {
	tmpl = template.Must(template.New("").Funcs(templateFuncs).ParseGlob("./templates/*.html"))
}

func blogHandler(w http.ResponseWriter, r *http.Request) {
	if err := renderTemplate("OZ | Blog", w, r); err != nil {
		http.Error(w, fmt.Sprintf("could not render tamplate %s", err.Error()), http.StatusInternalServerError)
		return
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if err := renderTemplate("OZ", w, r); err != nil {
		http.Error(w, fmt.Sprintf("could not render tamplate %s", err.Error()), http.StatusInternalServerError)
		return
	}
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	if err := renderTemplate("OZ | About me", w, r); err != nil {
		http.Error(w, fmt.Sprintf("could not render tamplate %s", err.Error()), http.StatusInternalServerError)
		return
	}
}

func tagsHandler(w http.ResponseWriter, r *http.Request) {
	if err := renderTemplate("OZ | tags", w, r); err != nil {
		http.Error(w, fmt.Sprintf("could not render tamplate %s", err.Error()), http.StatusInternalServerError)
		return
	}
}

func renderTemplate(name string, w http.ResponseWriter, r *http.Request) error {
	w.Header().Add(contentType, textHtml)
	return tmpl.ExecuteTemplate(w, "Layout", LayputInfo{
		DataMap: map[string]any{
			"message": "this is the template message",
			"name":    name,
		},
		PageName: name,
		Route:    r.URL.Path,
	})
}
