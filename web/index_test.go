package web

import "testing"

func Test_getHostname(t *testing.T) {
	for _, tt := range []struct {
		in, want string
	}{
		{"localhost", "localhost"},
		{"localhost:1234", "localhost"},
		{"[::1]", "[::1]"},
		{"[::1]:1234", "[::1]"},
		{"100.1.2.3", "100.1.2.3"},
		{"10.3.15.6:5614", "10.3.15.6"},
		{"dev.db.example.com", "dev.db.example.com"},
		{"pg.example.com:9999", "pg.example.com"},
	} {
		if got := getHostname(tt.in); got != tt.want {
			t.Errorf("want '%s', got '%s'", tt.want, got)
		}
	}
}
