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

func TestKeycapEmoji(t *testing.T) {
	for _, s := range []string{"1️⃣", "#️⃣", "*⃣"} {
		if !validEmoji(s) {
			t.Errorf("keycap %q rejected", s)
		}
	}
	for _, s := range []string{"1", "12️⃣", "1️"} {
		if validEmoji(s) {
			t.Errorf("invalid keycap %q accepted", s)
		}
	}
}
