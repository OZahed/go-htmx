package handlers

import "html/template"

const (
	contentType     = "Content-Type"
	textHtml        = "text/html"
	applicationJson = "application/json"
)

type PageMap map[string]any
type ExtraFunc template.FuncMap

type LayoutInfo struct {
	DataMap  PageMap
	PageName string
	Route    string
}
