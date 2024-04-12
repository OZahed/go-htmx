package handlers

import (
	"errors"
	"html/template"
	"path/filepath"
	"strings"
)

func LoadTemplates(rootDir string) *template.Template {
	return template.Must(template.New("").Funcs(template.FuncMap{
		"Contains": strings.Contains,
		"Dict":     Dict,
	}).ParseGlob(filepath.Join(rootDir, "*.html")))
}

func Dict(args ...string) (map[string]interface{}, error) {
	if len(args)%2 != 0 {
		return nil, errors.New("insufficient args number ")
	}

	mArgs := make(map[string]interface{})
	for i := 0; i < len(args); i += 2 {
		mArgs[args[i]] = args[i+1]
	}

	return mArgs, nil
}
