package main

import "encoding/json"

type defaultReaction struct {
	Content string `json:"content"`
	Code    string `json:"code,omitempty"`
	URL     string `json:"url,omitempty"`
}

func validDefaultReaction(r defaultReaction) bool {
	if r.Code == "" {
		return r.URL == "" && validEmoji(r.Content)
	}
	for _, c := range r.Code {
		if !(c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || c == '_') {
			return false
		}
	}
	return r.Content == ":"+r.Code+":" && safeHTTPS(r.URL)
}
func chooseDefaultReaction(raw, legacy string) defaultReaction {
	var r defaultReaction
	if json.Unmarshal([]byte(raw), &r) == nil && validDefaultReaction(r) {
		return r
	}
	var old []string
	if json.Unmarshal([]byte(legacy), &old) == nil && len(old) > 0 && validEmoji(old[0]) {
		return defaultReaction{Content: old[0]}
	}
	return defaultReaction{Content: "❤️"}
}
