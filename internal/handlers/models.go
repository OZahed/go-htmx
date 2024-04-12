package handlers

import "html/template"

const (
	contentType = "Content-Type"
	textHtml    = "text/html"
)

type PageMap map[string]any
type ExtraFunc template.FuncMap

type LayputInfo struct {
	DataMap  PageMap
	PageName string
	Route    string
}
