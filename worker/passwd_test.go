package worker

import "testing"

func Test_validPassword(t *testing.T) {
	const name = "Ol8fVmKqCrQjPh4lxD"

	if !validPassword(name, "foo", "foo") {
		t.Fatalf("valid password must have matched admin password")
	}
	if !validPassword(name, "", "bar") {
		t.Fatalf("missing passwd file must allow empty password")
	}
	if validPassword(name, "foo", "bar") {
		t.Fatalf("missing passwd file must NOT allow non-empty password")
	}
}
