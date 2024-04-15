package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var (
	timeFormats = []string{time.DateOnly, time.TimeOnly, time.DateTime, "2006-01-02 15:04:05-07:00",
		time.Kitchen, time.RFC3339, time.RFC1123, time.RFC1123Z, time.ANSIC,
		"2006/01/02", "2006/01/02 15:04:05", time.UnixDate, time.RubyDate}

	stringSeparators = []string{",", ";", ";", "-"}
)

var (
	// later on you can read string value over a config server
	// or read it from yaml or read value form json file
	// just replace the Getter function with any other kind of
	DefaultGetter ValueFunc = func(key, def string) string {
		val := os.Getenv(key)
		if val == "" {
			return def
		}

		return val
	}

	DefaultKeyBuilder KeyFunc = func(key string) string {
		return strings.ReplaceAll(strings.TrimSpace(key), ".", "_")
	}
)

// ValueFunc is the function is required because sometimes we need to read values sources other than os.getEnv
type ValueFunc func(key, def string) string

// KeyFunc is a function that returns altered keys, for example some times you need
// to replace some characters or you need to add a prefix or suffix
type KeyFunc func(string) string

type Marshaler struct {
	BuildKey func(string) string
	Get      func(name, def string) string
}

func NewMarshaler(k KeyFunc, g ValueFunc) *Marshaler {
	if g == nil {
		g = DefaultGetter
	}

	if k == nil {
		k = DefaultKeyBuilder
	}

	return &Marshaler{BuildKey: k, Get: g}
}

func KeyFuncWithPrefix(prefix string) KeyFunc {
	if !strings.HasSuffix(prefix, "_") {
		prefix = prefix + "_"
	}

	return func(key string) string {
		return prefix + key
	}
}

func (m *Marshaler) Marshal(dest any, prefix string) error {
	dst := reflect.ValueOf(dest)

	return m.marshal(dst, prefix)
}

// TODO: append key names
// TODO: Add Support for

//nolint:funlen
func (m *Marshaler) marshal(dst reflect.Value, prefix string) (err error) {
	valueType := dst.Type()

	if valueType.Kind() != reflect.Pointer {
		return errors.New("kind is not a pointer")
	}

	elm := valueType.Elem()
	if elm.Kind() != reflect.Struct {
		return errors.New("destination is not a struct")
	}

	valueType = valueType.Elem()
	dst = dst.Elem()

	for i := 0; i < valueType.NumField(); i++ {
		fieldValue := dst.Field(i)
		fieldType := valueType.Field(i)

		if !fieldType.IsExported() {
			continue
		}

		tagVal, hasKey := fieldType.Tag.Lookup("env")
		if !hasKey {
			continue
		}

		// set string up
		key, def := parseStructTags(tagVal)
		if prefix != "" {
			key = fmt.Sprintf("%s.%s", prefix, key)
		}

		// KeyBuilder removes
		strValues := m.Get(m.BuildKey(key), def)

		if strValues == "" {
			continue
		}

		err = m.marshalValue(fieldValue, strValues, prefix, key)
		if err != nil {
			return err
		}

	}

	return nil
}

func (m *Marshaler) marshalValue(fieldValue reflect.Value, value, prefix, key string) error {
	if !fieldValue.CanSet() {
		return nil
	}

	switch fieldValue.Kind() {
	case reflect.String:
		fieldValue.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int16, reflect.Int64:
		n, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		fieldValue.SetInt(n)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		n, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}

		fieldValue.SetUint(n)
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		fieldValue.SetFloat(f)
	case reflect.Map:
		err := makeMap(value, fieldValue)
		if err != nil {
			return err
		}

	case reflect.Slice:
		// TODO: parse Array
		splits := splitStr(value)

		if len(splits) > fieldValue.Len() {
			fieldValue.Grow(len(splits) - fieldValue.Len())
		}

		fieldValue.SetLen(len(splits))

		for i, split := range splits {
			if !fieldValue.Index(i).CanAddr() {
				continue
			}

			err := m.marshalValue(fieldValue.Index(i).Addr(), split, prefix, key)
			if err != nil {
				return err
			}
		}
	case reflect.Struct:
		if fieldValue.Type() == reflect.TypeOf(time.Time{}) {
			t, err := parseTime(value)
			if err != nil {
				return err
			}

			fieldValue.Set(reflect.ValueOf(t))
			return nil
		}

		// Url is always a pointer
		if fieldValue.Type() == reflect.TypeOf(new(url.URL)) {
			u, err := url.Parse(value)
			if err != nil {
				return err
			}

			fieldValue.Set(reflect.ValueOf(u))
			return nil
		}

		err := m.marshal(fieldValue.Addr(), fmt.Sprintf("%s.%s", prefix, key))
		if err != nil {
			return err
		}
	}

	if fieldValue.Type() == reflect.TypeOf(time.Duration(0)) {
		d, err := time.ParseDuration(value)
		if err != nil {
			return err
		}

		fieldValue.Set(reflect.ValueOf(d))
	}

	return nil
}

// I does not worth an scanner
func parseStructTags(tagVal string) (key, def string) {
	tagVal = strings.TrimSpace(tagVal)
	if tagVal == "-" || tagVal == "" {
		return "", ""
	}

	parts := strings.Split(tagVal, ",")
	key = parts[0]
	if len(parts) < 2 {
		return key, ""
	}

	parts[1] = strings.ReplaceAll(parts[1], "default=", "")

	def = parts[1]
	if len(parts) > 2 {
		def = strings.Join(parts[1:], ",")
	}

	return key, def
}

// Turns strings like: key1:val1,key2:val2 into map[K]V
// Only string and int are supported for now
//
// TODO: add map[string]any
func makeMap(str string, value reflect.Value) error {
	if value.Type().Kind() != reflect.Map {
		return fmt.Errorf("%s is not a map", value.Type().Name())
	}

	keyType := value.Type().Key()
	valueType := value.Type().Elem()
	value.Set(reflect.MakeMap(value.Type()))

	if keyType.Kind() != reflect.String && keyType.Kind() != reflect.Int {
		return fmt.Errorf("only int and string keys are acceptable, got key of type %s", keyType.Kind())
	}

	if valueType.Kind() != reflect.String && valueType.Kind() != reflect.Int {
		return fmt.Errorf("only key and value data values are acceptable, got key of type %s", valueType.Kind())
	}

	kv := splitStr(str)
	for _, pair := range kv {
		splt := strings.Split(pair, ":")
		if len(splt) < 2 {
			return fmt.Errorf("%s can not is in wrong format as key value pair", pair)
		}

		k := reflect.ValueOf(splt[0])
		v := reflect.ValueOf(splt[1])

		if keyType.Kind() == reflect.Int {
			n, err := strconv.ParseInt(splt[1], 10, 64)
			if err != nil {
				return fmt.Errorf("key %s for map is can not be turned into int", splt[1])
			}
			k = reflect.ValueOf(n)
		}

		if valueType.Kind() == reflect.Int {
			n, err := strconv.ParseInt(splt[1], 10, 64)
			if err != nil {
				return fmt.Errorf("value %s for map is can not be turned into int", splt[1])
			}
			v = reflect.ValueOf(n)
		}

		value.SetMapIndex(k, v)
	}

	return nil
}

func splitStr(value string) (split []string) {
	for _, sep := range stringSeparators {
		split = strings.Split(value, sep)
		if split[0] != value {
			return
		}
	}

	return
}

func parseTime(value string) (time.Time, error) {
	var err []error
	for _, format := range timeFormats {
		t, e := time.Parse(format, value)
		if e == nil {
			return t, nil
		}
		err = append(err, e)
	}
	return time.Time{}, errors.Join(err...)
}
