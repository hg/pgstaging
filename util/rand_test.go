package util

import (
	"regexp"
	"testing"
)

var validRandomStr = regexp.MustCompile(`^[a-zA-Z0-9]+$`)

func TestRandomString(t *testing.T) {
	seen := make(map[string]bool)

	for i := 0; i < 16; i++ {
		for _, length := range []uint{0, 8, 16, 25, 37} {
			str := RandomString(length)

			if len(str) != int(length) {
				t.Fatalf("wanted length %d, got %d (string %s)", length, len(str), str)
			}
			if length == 0 {
				continue
			}
			if !validRandomStr.MatchString(str) {
				t.Fatalf("expected '%s' to match alnum regex'", str)
			}
			if seen[str] {
				t.Fatalf("already seen string '%s'", str)
			}
			seen[str] = true
		}
	}
}
