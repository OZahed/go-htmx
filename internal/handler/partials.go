package handler

import (
	"errors"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"
)

// this is the most stupid form of count handling but go on with it for now
var count int

type Partials struct {
	tmpl *template.Template
	lg   *slog.Logger
}

func NewPartials(tmp *template.Template, lg *slog.Logger) *Partials {
	return &Partials{
		tmpl: tmp,
		lg:   lg,
	}
}

func (p *Partials) Counter(w http.ResponseWriter, r *http.Request) {
	countStr := r.PathValue("count")
	count, err := strconv.ParseInt(countStr, 10, 64)
	if err != nil {
		p.handleErr(errors.New("count value is not a number"), w, r)
		return
	}

	p.renderTemplate(PartialInfo{Content: map[string]interface{}{
		"Count": count + 1,
	},
		Name:   "Layout_Counter",
		Caller: r.Referer()}, w, r)
}

func (p *Partials) renderTemplate(data PartialInfo, w http.ResponseWriter, r *http.Request) {
	w.Header().Add(contentType, textHtml)
	err := p.tmpl.ExecuteTemplate(w, data.Name, data)

	if err != nil {
		p.lg.Error("failed to execute template", "name", data.Name, "path", r.URL.Path, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (p *Partials) handleErr(err error, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	p.lg.Error("error initiated", "path", r.URL.Path, "error", err.Error())

	w.Header().Add(contentType, textPlain)
	w.Write([]byte("bad request"))
}
