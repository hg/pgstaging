package sessions

import (
	"fmt"
	"testing"
)

func TestSession_AddEvent(t *testing.T) {
	var s Session

	s.AddEvent(StatusQueued, "foobar")
	if len(s.events) != 1 {
		t.Fatalf("expected one event, got %d", len(s.events))
	}

	for i := 0; i < keepLastEvents*2; i++ {
		s.AddEvent(StatusSuccess, fmt.Sprintf("event #%d", i))
	}

	if len(s.events) != keepLastEvents {
		t.Fatalf("expected %d events, got %d", keepLastEvents, len(s.events))
	}

	for _, cs := range []struct {
		idx int
		msg string
	}{
		{0, fmt.Sprintf("event #%d", keepLastEvents*2-1)},
		{1, fmt.Sprintf("event #%d", keepLastEvents*2-2)},
		{2, fmt.Sprintf("event #%d", keepLastEvents*2-3)},
		{keepLastEvents - 1, fmt.Sprintf("event #%d", keepLastEvents)},
	} {
		if got := s.events[cs.idx].Message; got != cs.msg {
			t.Fatalf("expected message '%s', got '%s'", cs.msg, got)
		}
		if got := s.events[cs.idx].Status; got != StatusSuccess {
			t.Fatalf("expected status %s, got %s", StatusSuccess, got)
		}
	}
}

func TestSession_Events(t *testing.T) {
	var s Session
	s.AddEvent(StatusSuccess, "success")
	s.AddEvent(StatusError, "error")

	evt := s.Events()

	checkEvt := func() {
		if len(evt) != 2 {
			t.Fatalf("expected length 2, got %d", len(evt))
		}
		if evt[0].Status != StatusError || evt[0].Message != "error" {
			t.Fatalf("invalid first record %v", evt[0])
		}
		if evt[1].Status != StatusSuccess || evt[1].Message != "success" {
			t.Fatalf("invalid second record %v", evt[1])
		}
	}

	checkEvt()

	s.AddEvent(StatusQueued, "queued")
	if len(s.events) != 3 {
		t.Fatalf("expected new event in %v", s.events)
	}

	if len(s.events) != 3 {
		t.Fatalf("did not expect new event in %v", evt)
	}

	checkEvt()
}
