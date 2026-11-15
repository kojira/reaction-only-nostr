package main

import "testing"

func TestChooseDefaultReaction(t *testing.T) {
	tests := []struct {
		raw, old string
		want     defaultReaction
	}{
		{"", "", defaultReaction{Content: "❤️"}},
		{"", `["👍","🎉"]`, defaultReaction{Content: "👍"}},
		{`{"content":"🎉"}`, `["👍"]`, defaultReaction{Content: "🎉"}},
		{`{"content":":cat:","code":"cat","url":"https://example.com/cat.png"}`, "", defaultReaction{Content: ":cat:", Code: "cat", URL: "https://example.com/cat.png"}},
		{`{"content":":cat:","code":"cat","url":"http://example.com/cat.png"}`, `["👍"]`, defaultReaction{Content: "👍"}},
		{`{"content":"x"}`, "broken", defaultReaction{Content: "❤️"}},
	}
	for _, tt := range tests {
		if got := chooseDefaultReaction(tt.raw, tt.old); got != tt.want {
			t.Errorf("got %#v want %#v", got, tt.want)
		}
	}
}
func TestInvalidDefaultReaction(t *testing.T) {
	for _, r := range []defaultReaction{{Content: "👍", URL: "https://example.com/a"}, {Content: ":bad-code:", Code: "bad-code", URL: "https://example.com/a"}, {Content: ":other:", Code: "cat", URL: "https://example.com/a"}} {
		if validDefaultReaction(r) {
			t.Errorf("accepted %#v", r)
		}
	}
}
