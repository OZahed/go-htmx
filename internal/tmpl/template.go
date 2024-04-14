package tmpl

import (
	"errors"
	"html/template"
	"path/filepath"
	"reflect"
	"strings"
)

func LoadTemplates(rootDir string) *template.Template {
	return template.Must(template.New("").Funcs(template.FuncMap{
		"Contains": strings.Contains,
		"Dict":     Dict,
		"Attach":   Attach,
	}).ParseGlob(filepath.Join(rootDir, "*.html")))
}

func Dict(keyVal ...string) (map[string]interface{}, error) {
	if len(keyVal)%2 != 0 {
		return nil, errors.New("insufficient number of arguments")
	}

	mArgs := make(map[string]interface{})
	for i := 0; i < len(keyVal); i += 2 {
		mArgs[keyVal[i]] = keyVal[i+1]
	}

	return mArgs, nil
}

func Attach(s interface{}, keyVal ...string) (map[string]interface{}, error) {
	sType := reflect.TypeOf(s)
	acceptable := false
	switch sType.Kind() {
	case reflect.Struct:
		acceptable = true
	case reflect.Map:
		_, acceptable = s.(map[string]interface{})
	case reflect.Pointer:
		if sType.Elem().Kind() == reflect.Struct {
			sType = sType.Elem()
			acceptable = true
		}
	}

	if !acceptable {
		return nil, errors.New("unacceptable base type")
	}

	if len(keyVal)%2 != 0 {
		return nil, errors.New("insufficient number of arguments")
	}

	res := make(map[string]interface{})
	val := reflect.ValueOf(s)

	// we already checked the if the map is map[string]interface{} using acceptable
	if sType.Kind() == reflect.Map {
		return attachMap(s.(map[string]interface{}), keyVal...)
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		inter := field.Interface()
		name := sType.Field(i).Name
		res[name] = inter
	}

	for i := 0; i < len(keyVal); i += 2 {
		res[keyVal[i]] = keyVal[i+1]
	}

	return res, nil
}

func attachMap(m map[string]interface{}, args ...string) (map[string]interface{}, error) {
	if len(args) == 0 || len(m) == 0 {
		return m, nil
	}

	for i := 0; i < len(args); i += 2 {
		m[args[i]] = args[i+1]
	}

	return m, nil
}
