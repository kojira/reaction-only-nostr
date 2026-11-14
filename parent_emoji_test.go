package main

import (
	"reflect"
	"testing"
)

func TestParentEmojiParts(t *testing.T) {
	tags := [][]string{{"emoji", "ok", "https://example.org/a.png"}, {"emoji", "bad", "javascript:alert(1)"}, {"emoji", "x-y", "https://example.org/x"}, {"emoji", "short"}}
	for _, kind := range []int{1, 7} {
		ev := Event{Kind: kind, Content: "hello :ok: :bad: :unknown: :x-y: :ok:", Tags: tags}
		got := parentEmojiParts(ev)
		var text string
		var urls []string
		for _, p := range got {
			text += p.Text
			if p.URL != "" {
				urls = append(urls, p.URL)
			}
		}
		if text != ev.Content || !reflect.DeepEqual(urls, []string{"https://example.org/a.png", "https://example.org/a.png"}) {
			t.Fatalf("kind %d: %#v", kind, got)
		}
	}
	if got := parentEmojiParts(Event{Content: ":ok:", Tags: [][]string{{"emoji", "ok", "https://user@example.org/a"}}}); len(got) == 0 || got[0].URL != "" {
		t.Fatal(got)
	}
}
