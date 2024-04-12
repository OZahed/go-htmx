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

func NewLayoutHanler(tmp *template.Template, sn, rootTempName string, lg *slog.Logger) *LayoutHandlers {
	return &LayoutHandlers{
		tmpl:     tmp,
		lg:       lg,
		siteName: sn,
		root:     rootTempName,
	}
}

func (lh *LayoutHandlers) BlogHandler(w http.ResponseWriter, r *http.Request) {
	lh.renderTemplate(fmt.Sprintf("%s > Blog", lh.siteName), lh.root, w, r)
}

func (lh *LayoutHandlers) IndexHandler(w http.ResponseWriter, r *http.Request) {
	lh.renderTemplate(lh.siteName, lh.root, w, r)
}

func (lh *LayoutHandlers) AboutHandler(w http.ResponseWriter, r *http.Request) {
	lh.renderTemplate(fmt.Sprintf("%s > about", lh.siteName), lh.root, w, r)
}

func (lh *LayoutHandlers) TagsHandler(w http.ResponseWriter, r *http.Request) {
	lh.renderTemplate(fmt.Sprintf("%s > Tags", lh.siteName), lh.root, w, r)
}

func (lh *LayoutHandlers) renderTemplate(name, root string, w http.ResponseWriter, r *http.Request) {
	w.Header().Add(contentType, textHtml)
	err := lh.tmpl.ExecuteTemplate(w, root, LayputInfo{
		DataMap: map[string]any{
			"message": "this is the template message",
			"name":    name,
		},
		PageName: name,
		Route:    r.URL.Path,
	})

	if err != nil {
		lh.lg.Error("failed to execute template", "name", name, "path", r.URL.Path, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
