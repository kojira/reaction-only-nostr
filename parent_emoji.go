package main

import "strings"

type parentEmojiPart struct{ Text, URL string }

// Only the displayed event's own NIP-30 tags may define its shortcodes.
func parentEmojiParts(ev Event) []parentEmojiPart {
	urls := map[string]string{}
	for _, tag := range ev.Tags {
		if len(tag) < 3 || tag[0] != "emoji" || !validParentShortcode(tag[1]) || !safeHTTPS(tag[2]) {
			continue
		}
		if _, exists := urls[tag[1]]; !exists {
			urls[tag[1]] = tag[2]
		}
	}
	var parts []parentEmojiPart
	text := ev.Content
	for len(text) > 0 {
		start := strings.IndexByte(text, ':')
		if start < 0 {
			parts = append(parts, parentEmojiPart{Text: text})
			break
		}
		end := strings.IndexByte(text[start+1:], ':')
		if end < 0 {
			parts = append(parts, parentEmojiPart{Text: text})
			break
		}
		end += start + 1
		code := text[start+1 : end]
		if url, ok := urls[code]; ok {
			if start > 0 {
				parts = append(parts, parentEmojiPart{Text: text[:start]})
			}
			parts = append(parts, parentEmojiPart{Text: text[start : end+1], URL: url})
			text = text[end+1:]
		} else {
			parts = append(parts, parentEmojiPart{Text: text[:start+1]})
			text = text[start+1:]
		}
	}
	return parts
}
func validParentShortcode(code string) bool {
	if code == "" {
		return false
	}
	for _, c := range code {
		if !(c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || c == '_') {
			return false
		}
	}
	return true
}
