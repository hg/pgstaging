package util

import (
	"testing"
)

func TestGetUserIdNobody(t *testing.T) {
	id, err := GetUserId("nobody")
	if err != nil {
		return
	}
	if id.UID == 0 || id.GID == 0 {
		t.Fatalf("unexpected root ID")
	}
	id2, err := GetUserId("nobody")
	if err != nil {
		t.Fatalf("got error on second try: %v", err)
	}
	if id.UID != id2.UID {
		t.Fatalf("UIDs do not match: %d, %d", id.UID, id2.UID)
	}
	if id.GID != id2.GID {
		t.Fatalf("GIDs do not match: %d, %d", id.GID, id2.GID)
	}
}

func TestGetUserIdRoot(t *testing.T) {
	id, err := GetUserId("root")
	if err != nil {
		t.Fatalf("got error: %v", err)
	}
	if id.UID != 0 || id.GID != 0 {
		t.Fatalf("unexpected uid or gid: %v", id)
	}
}

func TestGetUserFake(t *testing.T) {
	id, err := GetUserId("3e18JXHNqHlutC2t9j")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if id.UID != 0 || id.GID != 0 {
		t.Fatalf("expected zero result, got %v", id)
	}
}
