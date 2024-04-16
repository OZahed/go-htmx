package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	r "reflect"
	"strconv"
	"strings"
	"time"
)

var (
	timeFormats = []string{time.DateOnly, time.TimeOnly, time.DateTime, "2006-01-02 15:04:05-07:00",
		time.Kitchen, time.RFC3339, time.RFC1123, time.RFC1123Z, time.ANSIC,
		"2006/01/02", "2006/01/02 15:04:05", time.UnixDate, time.RubyDate}

	stringSeparators = []string{",", ";", ";", "-", " "}

	timeType     = r.TypeOf(time.Time{})
	durationType = r.TypeOf(time.Duration(0))
	urlType      = r.TypeOf(&url.URL{})
)

var (
	// later on you can read string value over a config server
	// or read it from yaml or read value form json file
	// just replace the Getter function with any other kind of
	DefaultEnvGetter ValueFunc = func(key, def string) string {
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
		g = DefaultEnvGetter
	}

	if k == nil {
		k = DefaultKeyBuilder
	}

	return &Marshaler{BuildKey: k, Get: g}
}

// TODO: append key names
// TODO: Add Support for

//nolint:funlen
func (m *Marshaler) Marshal(dest interface{}, prefix string) (err error) {
	dst := r.ValueOf(dest)
	valueType := dst.Type()

	if valueType.Kind() != r.Pointer {
		return fmt.Errorf("Kind %s is not a pointer", valueType.Kind())
	}

	if dst.IsNil() {
		dst.Set(r.New(dst.Type().Elem()))
	}

	elm := valueType.Elem()
	if elm.Kind() != r.Struct {
		return fmt.Errorf("destination is of type %s and not struct", elm.Kind())
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

		if strValues == "" && fieldType.Type.Kind() != r.Struct {
			continue
		}

		err = m.MarshalValue(fieldValue, strValues, prefix, key)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Marshaler) MarshalValue(reflectValue r.Value, strValue, prefix, key string) error {
	if !reflectValue.CanSet() {
		return nil
	}

	if reflectValue.Kind() == r.Func {
		return nil
	}

	// Checking for non-builtin types
	switch reflectValue.Type() {
	case timeType:
		t, err := parseTime(strValue)
		if err != nil {
			return err
		}
		reflectValue.Set(r.ValueOf(t))
		return nil
	case urlType:
		u, err := url.Parse(strValue)
		if err != nil {
			return err
		}

		reflectValue.Set(r.ValueOf(u))
		return nil
	case durationType:
		d, err := time.ParseDuration(strValue)
		if err != nil {
			return err
		}

		reflectValue.Set(r.ValueOf(d))
		return nil
	}

	// Checking for built int types
	switch reflectValue.Kind() {
	case r.String:
		reflectValue.SetString(strValue)
	case r.Int, r.Int8, r.Int32, r.Int16, r.Int64:
		n, err := strconv.ParseInt(strValue, 10, 64)
		if err != nil {
			return err
		}
		reflectValue.SetInt(n)
	case r.Uint, r.Uint8, r.Uint16, r.Uint32, r.Uint64, r.Uintptr:
		n, err := strconv.ParseUint(strValue, 10, 64)
		if err != nil {
			return err
		}

		reflectValue.SetUint(n)
	case r.Float32, r.Float64:
		f, err := strconv.ParseFloat(strValue, 64)
		if err != nil {
			return err
		}
		reflectValue.SetFloat(f)
	case r.Bool:
		b, err := strconv.ParseBool(strValue)
		if err != nil {
			return err
		}

		reflectValue.SetBool(b)
	case r.Map:
		return m.parseMap(reflectValue, strValue)
	case r.Slice:
		return m.parseArray(strValue, reflectValue, key)
	case r.Struct:
		// The ParseEnv should be on pointer
		ptr := reflectValue.Addr()
		// checking for ParseEnv() error method first
		parser := ptr.MethodByName("ParseEnv")
		if parser.IsValid() {
			callResult := parser.Call([]r.Value{r.ValueOf(key)})

			e := callResult[0].Interface()

			if e == nil {
				return nil
			}

			return e.(error)
		}

		if !reflectValue.CanAddr() || reflectValue.Type() == r.TypeOf(struct{}{}) {
			return nil
		}

		return m.Marshal(reflectValue.Addr().Interface(), key)
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
func (m *Marshaler) parseMap(value r.Value, str string) (err error) {
	if value.Type().Kind() != r.Map {
		return fmt.Errorf("%s is not a map", value.Type().Name())
	}

	fmt.Println(value.Len())

	keyType := value.Type().Key()
	valueType := value.Type().Elem()
	value.Set(r.MakeMap(value.Type()))

	kv := splitStr(str)
	for _, pair := range kv {
		splt := strings.Split(pair, ":")
		fmt.Println(value.Len())
		if len(splt) < 2 {
			return fmt.Errorf("%s can not is in wrong format as key value pair", pair)
		}

		keyStr := strings.TrimSpace(splt[0])
		valStr := strings.TrimSpace(splt[1])
		k := r.New(keyType).Elem()
		v := r.New(valueType).Elem()
		fmt.Printf("key = %s value = %s\n", k.Kind().String(), v.Kind().String())

		if err = m.MarshalValue(k, keyStr, "", ""); err != nil {
			return fmt.Errorf("%s can not be parsed as %s", keyStr, k.Kind())
		}

		if err = m.MarshalValue(v, valStr, "", ""); err != nil {
			return fmt.Errorf("%s can not be parsed as %s", valStr, v.Kind())
		}

		value.SetMapIndex(k, v)
	}
	return nil
}

func (m *Marshaler) parseArray(value string, fieldValue r.Value, currentKey string) error {
	splits := splitStr(value)

	if len(splits) > fieldValue.Len() {
		fieldValue.Grow(len(splits) - fieldValue.Len())
	}

	fieldValue.SetLen(len(splits))

	for i, split := range splits {
		split = strings.TrimSpace(split)
		// for slice values prefix should become key and there should be no keys
		err := m.MarshalValue(fieldValue.Index(i), split, currentKey, "")
		if err != nil {
			return err
		}
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
