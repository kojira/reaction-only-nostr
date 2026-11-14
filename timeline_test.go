package main

import "testing"

func TestTimeline(t *testing.T) {
	e := func(id, parent string, ts int64) Event {
		return Event{ID: id, CreatedAt: ts, Tags: [][]string{{"e", parent}}}
	}
	got := buildTimeline([]Event{e("child", "root", 3), e("root", "missing", 2), e("child", "root", 3), e("a", "b", 5), e("b", "a", 4), e("new", "", 9)})
	if len(got) != 5 {
		t.Fatalf("count %d", len(got))
	}
	if got[0].Event.ID != "new" || got[1].Event.ID != "root" || !got[1].MissingParent || got[2].Event.ID != "child" || got[2].Depth != 1 || !got[3].Cycle {
		t.Fatalf("unexpected tree %+v", got)
	}
}
