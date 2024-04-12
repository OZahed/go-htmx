package logger

import "testing"

func TestHumanReadableBytes(t *testing.T) {
	tests := []struct {
		name string
		want string
		n    int
	}{
		{
			name: "should return 100 B",
			n:    100,
			want: "100 B",
		},
		{
			name: "should return 1000 B",
			n:    1000,
			want: "1000 B",
		},
		{
			name: "1 Killo byte",
			n:    1 << 10,
			want: "1.00 KB",
		},
		{
			name: "1.20 Killo byte",
			n:    1<<10 + 213,
			want: "1.21 KB",
		},
		{
			name: "1.03 Killo byte",
			n:    1<<10 + 3,
			want: "1.00 KB",
		},
		{
			name: "1 Giga byte ",
			n:    1 << 20,
			want: "1.00 GB",
		},
		{
			name: "1.5 Giga byte (it is common to use one digit after truncing the float64 to 2 digits)",
			n:    1<<20 + (500 * (1 << 10)),
			want: "1.49 GB",
		},
		{
			name: "1.3 Giga byte (it is common to use one digit after truncing the float64 to 2 digits)",
			n:    1<<20 + (300 * (1 << 10)),
			want: "1.29 GB",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HumanReadableBytes(tt.n)
			if got != tt.want {
				t.Errorf("HumanReadableBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
