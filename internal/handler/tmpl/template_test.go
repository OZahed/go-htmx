package tmpl_test

import (
	"reflect"
	"testing"

	"github.com/OZahed/go-htmx/internal/handler/tmpl"
)

func TestAttach(t *testing.T) {
	wants := map[string]interface{}{
		"Key1":    "1",
		"Key2":    2,
		"string1": "string2",
	}

	type input struct {
		Key1 string
		Key2 int
	}

	testVal := input{"1", 2}
	inString := []string{"string1", "string2"}
	type args struct {
		s      interface{}
		keyVal []string
	}
	tests := []struct {
		want    map[string]interface{}
		name    string
		args    args
		wantErr bool
	}{
		{
			want: wants,
			name: "struct",
			args: args{
				s:      testVal,
				keyVal: inString,
			},
			wantErr: false,
		},
		{
			want: nil,
			name: "map error due to map type",
			args: args{
				s: any(map[string]int{}),
			},
			wantErr: true,
		},
		{
			want: wants,
			name: "map attach",
			args: args{
				s: any(map[string]interface{}{
					"Key1": "1",
					"Key2": 2,
				}),
				keyVal: inString,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tmpl.Attach(tt.args.s, tt.args.keyVal...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Attach() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Attach() = %v, want %v", got, tt.want)
			}
		})
	}
}
