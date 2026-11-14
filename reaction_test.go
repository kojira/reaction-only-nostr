package main

import "testing"

func TestValidEmoji(t *testing.T) {
	for _, s := range []string{"❤️", "👍🏽", "👨‍👩‍👧", "🇯🇵"} {
		if !validEmoji(s) {
			t.Errorf("rejected %q", s)
		}
	}
	for _, s := range []string{"", "hello", "❤️ hello", "👍👍", "‍👍", "👍‍"} {
		if validEmoji(s) {
			t.Errorf("accepted %q", s)
		}
	}
}
