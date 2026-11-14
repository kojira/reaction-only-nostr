package main

import "testing"

func TestReactionParent(t *testing.T) {
	ev := Event{Kind: 7, Tags: [][]string{{"e", "root"}, {"p", "author"}, {"e", "parent", "wss://example"}}}
	if got := reactionParent(ev); got != "parent" {
		t.Fatalf("got %q", got)
	}
	ev.Kind = 1
	if got := reactionParent(ev); got != "" {
		t.Fatalf("kind 1: %q", got)
	}
}
