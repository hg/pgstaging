package web

import (
	"github.com/hg/pgstaging/util"
	"testing"
)

func Test_toIsoCodes(t *testing.T) {
	got := toIsoCodes("fr-CH, en;q=0.8, FR;q=0.7, de;q=0.7, *;q=0.5")
	want := []string{"fr", "en", "de"}

	if !util.SlicesEqual(want, got) {
		t.Fatalf("want %v, got %v", want, got)
	}
}
