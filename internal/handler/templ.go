package handler

import (
	"log/slog"
	"net/http"

	"github.com/OZahed/go-htmx/internal/templ-files"
)

type TemplHandler struct {
	lg *slog.Logger
}

func NewTempleHandler(lg *slog.Logger) *TemplHandler {
	return &TemplHandler{
		lg: lg,
	}
}

func (t *TemplHandler) Greet(w http.ResponseWriter, r *http.Request) {
	age := r.PathValue("age")
	name := r.PathValue("name")

	comp := templ.Hello(name, age)

	_ = comp.Render(r.Context(), w)
}
