package main

import (
	"encoding/hex"
	"strings"
)

func emojiAddress(tags [][]string) string {
	for _, t := range tags {
		if len(t) > 1 && t[0] == "d" {
			return t[1]
		}
	}
	return ""
}
func emojiReferences(tags [][]string) []string {
	var out []string
	seen := map[string]bool{}
	for _, t := range tags {
		if len(t) < 2 || t[0] != "a" {
			continue
		}
		p := strings.SplitN(t[1], ":", 3)
		if len(p) != 3 || p[0] != "30030" || !emojiValidPubkey(p[1]) || seen[t[1]] {
			continue
		}
		seen[t[1]] = true
		out = append(out, t[1])
	}
	return out
}

func emojiValidPubkey(s string) bool {
	if len(s) != 64 || strings.ToLower(s) != s {
		return false
	}
	_, err := hex.DecodeString(s)
	return err == nil
}
