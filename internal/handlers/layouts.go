package handlers

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
)

type LayoutHandlers struct {
	tmpl     *template.Template
	lg       *slog.Logger
	siteName string
	root     string
}

func NewLayoutHandler(tmp *template.Template, sn, rootTempName string, lg *slog.Logger) *LayoutHandlers {
	return &LayoutHandlers{
		tmpl:     tmp,
		lg:       lg,
		siteName: sn,
		root:     rootTempName,
	}
}

func (lh *LayoutHandlers) BlogHandler(w http.ResponseWriter, r *http.Request) {
	data := LayoutInfo{
		SubTmplName: "Blog",
		PageName:    fmt.Sprintf("%s > Blog", lh.siteName),
		Route:       r.URL.Path,
	}
	lh.renderTemplate(data, w, r)
}

func (lh *LayoutHandlers) IndexHandler(w http.ResponseWriter, r *http.Request) {
	data := LayoutInfo{
		SubTmplName: "",
		PageName:    lh.siteName,
		Route:       r.URL.Path,
	}
	lh.renderTemplate(data, w, r)
}

func (lh *LayoutHandlers) AboutHandler(w http.ResponseWriter, r *http.Request) {
	data := LayoutInfo{
		SubTmplName: "About",
		PageName:    fmt.Sprintf("%s > About", lh.siteName),
		Route:       r.URL.Path,
	}
	lh.renderTemplate(data, w, r)
}

func (lh *LayoutHandlers) TagsHandler(w http.ResponseWriter, r *http.Request) {
	data := LayoutInfo{
		SubTmplName: "Tags",
		PageName:    fmt.Sprintf("%s > Tags", lh.siteName),
		Route:       r.URL.Path,
	}
	lh.renderTemplate(data, w, r)
}

func (lh *LayoutHandlers) renderTemplate(data LayoutInfo, w http.ResponseWriter, r *http.Request) {
	w.Header().Add(contentType, textHtml)
	err := lh.tmpl.ExecuteTemplate(w, lh.root, data)
	if err != nil {
		lh.lg.Error("failed to execute template", "path", r.URL.Path, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
