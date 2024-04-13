package config

import (
	"net/url"
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestGetEnv(t *testing.T) {
	const (
		urlString = "https://test.com"
		unixDate  = "Sat Apr 13 17:42:36 +0330 2024"
		rfcDate   = "2024-04-13 14:12:56+00:00"
		port      = 3000
		intVal    = 123456
		appName   = "TEST"
		badKey    = "BAD_KEY"
		stringVal = "abcdefg"
		timeStr   = "2024-01-01"
	)

	strings := []string{"item1", "item2", "item3"}

	testEnvs := map[string]string{
		"TEST_PORT":        strconv.Itoa(port),
		"TEST_DURATION":    "2s",
		"TEST_URL":         urlString,
		"TEST_STRING_VAL":  stringVal,
		"TEST_INT_VAL":     strconv.Itoa(intVal),
		"TEST_STRINGS":     "item1 item2 item3",
		"TEST_STRINGS1":    "item1,item2,item3",
		"TEST_STRINGS2":    "item1;item2;item3",
		"TEST_STRINGS_ONE": "item1",
		"TEST_DATE":        "2024-01-01",
		"TEST_BOOL1":       "t",
		"TEST_BOOL2":       "1",
		"TEST_BOOL3":       "True",
		"TEST_BOOL4":       "False",
	}

	for k, v := range testEnvs {
		os.Setenv(k, v)
	}

	u, _ := url.Parse(urlString)
	defUrl, _ := url.Parse("http://localhost:8080")
	date, _ := time.Parse(time.DateOnly, timeStr)

	keyProvider := makeKeyProviderPrefix(appName)

	type args struct {
		def  any
		name string
	}
	tests := []struct {
		want any
		name string
		args args
	}{
		{
			name: "int value",
			want: any(port),
			args: args{
				name: keyProvider("PORT"),
				def:  8080,
			},
		},
		{
			name: "default int",
			want: any(400),
			args: args{
				name: keyProvider(badKey),
				def:  400,
			},
		},
		{
			name: "test duration",
			want: any(time.Second * 2),
			args: args{
				name: keyProvider("DURATION"),
				def:  time.Millisecond,
			},
		},
		{
			name: "test duration default",
			want: any(time.Hour),
			args: args{
				name: keyProvider(badKey),
				def:  time.Hour,
			},
		},
		{
			name: "test url",
			want: any(u),
			args: args{
				name: keyProvider("URL"),
				def:  defUrl,
			},
		},
		{
			name: "test url default",
			want: any(defUrl),
			args: args{
				name: keyProvider(badKey),
				def:  defUrl,
			},
		},
		{
			name: "test string",
			want: any(stringVal),
			args: args{
				name: keyProvider("STRING_VAL"),
				def:  "",
			},
		},
		{
			name: "test string bad key",
			want: any(stringVal),
			args: args{
				name: keyProvider(badKey),
				def:  stringVal,
			},
		},
		{
			name: "test string slice",
			want: any(strings),
			args: args{
				name: keyProvider("STRINGS"),
				def:  []string{},
			},
		},
		{
			name: "test string slice1",
			want: any(strings),
			args: args{
				name: keyProvider("STRINGS1"),
				def:  []string{},
			},
		},
		{
			name: "test string slice2",
			want: any(strings),
			args: args{
				name: keyProvider("STRINGS2"),
				def:  []string{},
			},
		},
		{
			name: "should return array with one element",
			want: any([]string{"item1"}),
			args: args{
				name: keyProvider("STRINGS_ONE"),
				def:  []string{},
			},
		},
		{
			name: "test time",
			want: any(date),
			args: args{
				name: keyProvider("DATE"),
				def:  time.Time{},
			},
		},
		{
			name: "test bool1",
			want: any(true),
			args: args{
				name: keyProvider("BOOL1"),
				def:  false,
			},
		},
		{
			name: "test bool2",
			want: any(true),
			args: args{
				name: keyProvider("BOOL2"),
				def:  false,
			},
		},
		{
			name: "test bool3",
			want: any(true),
			args: args{
				name: keyProvider("BOOL3"),
				def:  false,
			},
		},
		{
			name: "test bool4",
			want: any(false),
			args: args{
				name: keyProvider("BOOL4"),
				def:  true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getEnv(tt.args.name, tt.args.def)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}
