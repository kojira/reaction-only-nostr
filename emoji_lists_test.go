package main

import "testing"

func TestEmojiReferences(t *testing.T) {
	p := "2fc29418a41a18753e56eac7953a8c2ffbd0c7118a38d6fa3ae41b9e8bce6b0a"
	good := "30030:" + p + ":pack:variant"
	refs := emojiReferences([][]string{{"a", good, "wss://example.com"}, {"a", good}, {"a", "30030:bad:x"}, {"a", "30031:" + p + ":x"}, {"emoji", "foo", "https://example.com/foo.png"}})
	if len(refs) != 1 || refs[0] != good {
		t.Fatalf("unexpected refs: %v", refs)
	}
}
