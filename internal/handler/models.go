package handler

import "html/template"

const (
	contentType     = "Content-Type"
	textHtml        = "text/html"
	textPlain       = "text/plain"
	applicationJson = "application/json"
)

type ExtraFunc template.FuncMap

type LayoutInfo struct {
	Meta        map[string]interface{}
	SubTmplName string
	PageName    string
	Route       string
}

type PartialInfo struct {
	Content map[string]interface{}
	Name    string
	Caller  string
}
