package util

import "testing"

func TestAddPrefix(t *testing.T) {
	for _, tc := range []struct {
		in, want string
	}{
		{"", "dev_"},
		{"foo", "dev_foo"},
		{"dev_bar", "dev_bar"},
	} {
		got := AddPrefix("foobar")
		if got != "dev_foobar" {
			t.Fatalf("want '%s', got '%s'", tc.want, got)
		}
	}
}

func TestIsDevName(t *testing.T) {
	for _, tt := range []struct {
		in  string
		dev bool
	}{
		{"", false},
		{"main", false},
		{"dev_", false},
		{"DEV_foo", false},
		{"dev_x", true},
		{"dev_foo", true},
	} {
		if got := IsDevName(tt.in); got != tt.dev {
			t.Errorf("want %v, got %v", tt.dev, got)
		}
	}
}

func TestNormalizeName(t *testing.T) {
	for _, tt := range []struct {
		in, want string
	}{
		{"", ""},
		{"!@#$%^&*()_/.,;:'\"`\\~", "_"},
		{" FoO_baR_42!#$\n \t", "foo_bar_42"},
		{"i18n names: дед мороз, おそ松さん", "i18nnames"},
	} {
		if got := NormalizeName(tt.in); got != tt.want {
			t.Errorf("want '%s', got '%s'", tt.want, got)
		}
	}
}

func TestIsOkPassword(t *testing.T) {
	for _, tt := range []struct {
		in   string
		want bool
	}{
		{"", false},
		{"$!%^@&*#(!", false},
		{"foo_foo", false},
		{"Fo4ObAr1337@", false},
		{"foo_foox", true},
		{"FOObarBAZ", true},
		{"foobar42", true},
		{"Fo4ObA_r1337", true},
	} {
		if got := IsOkPassword(tt.in); got != tt.want {
			t.Errorf("value '%s': want %v, got %v", tt.in, tt.want, got)
		}
	}
}
