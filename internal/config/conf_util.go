/*
config should be turned into a dedicated package for itself
- [ ] Add struct tags parsing
- [ ] Add cache layer on os.Getenv
- [ ] Add Struct Load using Struct Tags
- [ ] Add Generate sample config file
*/
package config

import (
	"fmt"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var (
	separators  = []string{" ", ",", ";", "\n"}
	timeLayouts = []string{time.RFC3339, time.RFC3339Nano, time.DateTime, time.Stamp, time.DateOnly,
		time.TimeOnly, time.ANSIC, time.RFC822, time.UnixDate,
	}
)

type Getter struct {
	keyProvider func(name string) string
}

func (a *Getter) GetString(name, def string) string {
	return parseEnv(a.keyProvider(name), def).(string)
}

func (a *Getter) GetStringSlice(name string) []string {
	return parseEnv(name, []string{}).([]string)
}

func (a *Getter) GetInt(name string, def int) int {
	return parseEnv(a.keyProvider(name), def).(int)
}

func (a *Getter) GetInt64(name string, def int64) int64 {
	return parseEnv(a.keyProvider(name), def).(int64)
}

func (a *Getter) GetInt32(name string, def int32) int32 {
	return parseEnv(a.keyProvider(name), def).(int32)
}

func (a *Getter) GetFloat64(name string, def float64) float64 {
	return parseEnv(a.keyProvider(name), def).(float64)
}

func (a *Getter) GetFloat32(name string, def float32) float32 {
	return parseEnv(a.keyProvider(name), def).(float32)
}

func (a *Getter) GetBool(name string) bool {
	return parseEnv(a.keyProvider(name), false).(bool)
}

func (a *Getter) GetTime(name string) time.Time {
	return parseEnv(a.keyProvider(name), time.Time{}).(time.Time)
}

func (a *Getter) GetDuration(name string, def time.Duration) time.Duration {
	return parseEnv(a.keyProvider(name), def).(time.Duration)
}

func (a *Getter) GetUrl(name string) *url.URL {
	return parseEnv(a.keyProvider(name), (*url.URL)(nil)).(*url.URL)
}

// Name Without prefix
func parseEnv(name string, def any) any {
	val := os.Getenv(name)
	if val == "" {
		return def
	}

	switch def.(type) {
	case string:
		if val == "" {
			return def
		}
		return val
	case []string:
		if len(val) == 0 {
			return def
		}
		for _, sep := range separators {
			split := strings.Split(val, sep)
			if split[0] != val {
				return split
			}
		}
		return []string{val}
	case int:
		n := int(parseInt64(val))
		if n == 0 {
			return def
		}
		return n
	case int32:
		n := parseInt64(val)
		if n == 0 {
			return def
		}
		return int32(n)

	case int64:
		n := parseInt64(val)
		if n == 0 {
			return def
		}
		return n
	case float64:
		res, _ := strconv.ParseFloat(val, 64)
		if res == 0 {
			return def
		}

		return res
	case float32:
		res, _ := strconv.ParseFloat(val, 32)
		if res == 0 {
			return def
		}
		return res
	case bool:
		res, err := strconv.ParseBool(val)
		if err != nil {
			return def
		}

		return res

	case time.Time:
		for _, layout := range timeLayouts {
			t, err := time.Parse(layout, val)
			if err == nil && !t.IsZero() {
				return t
			}
		}
		return def
	case time.Duration:
		d, err := time.ParseDuration(val)
		if err != nil || d == 0 {
			return def
		}
		return d

	case *url.URL:
		u, err := url.Parse(val)
		if err != nil {
			return def
		}
		return u
	default:
		return def
	}
}

func GetDefault[T any](name string, def T) T {
	val := Get[T](name)

	if reflect.ValueOf(val).IsZero() {
		return def
	}

	return val
}

func Get[T any](name string) T {
	tp := reflect.TypeFor[T]()
	var res any

	val := os.Getenv(name)
	if val == "" {
		return reflect.New(tp).Elem().Interface().(T)
	}

	switch tp.Kind() {
	case reflect.String:
		res = val
	case reflect.Slice:
		if tp.Elem().Kind() == reflect.String {
			for _, sep := range separators {
				split := strings.Split(val, sep)
				if split[0] != val {
					res = split
					break
				}
			}
		}

		if tp.Elem().Kind() == reflect.Int {
			split := strings.Split(val, ",")
			arr := make([]int, 0)
			for _, str := range split {
				arr = append(arr, int(parseInt64(str)))
			}
			res = arr
		}
	case reflect.Int:
		res = int(parseInt64(val))
	case reflect.Int32:
		res = int32(parseInt64(val))
	case reflect.Int64:
		res = parseInt64(val)
	case reflect.Float64:
		res, _ = strconv.ParseFloat(val, 64)
	case reflect.Float32:
		res, _ = strconv.ParseFloat(val, 32)
	case reflect.Bool:
		res, _ = strconv.ParseBool(val)
	}

	if tp == reflect.TypeOf(time.Duration(0)) {
		res, _ = time.ParseDuration(val)
	}

	if tp == reflect.TypeOf(time.Time{}) {
		for _, layout := range timeLayouts {
			t, err := time.Parse(layout, val)
			if err == nil && !t.IsZero() {
				res = t
				break
			}
		}
	}

	if res == nil {
		fmt.Println("nil")
	}

	if reflect.TypeOf(res) != reflect.TypeFor[T]() {
		return reflect.New(tp).Elem().Interface().(T)
	}

	return reflect.ValueOf(res).Interface().(T)
}

func parseInt64(val string) int64 {
	n, _ := strconv.ParseInt(val, 10, 64)
	return n
}

func makeKeyProviderPrefix(prefix string) func(name string) string {
	return func(name string) string {
		if prefix == "" {
			return name
		}

		return prefix + "_" + name
	}
}
