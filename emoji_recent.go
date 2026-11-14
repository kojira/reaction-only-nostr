//go:build js && wasm

package main

import (
	"encoding/json"
	"syscall/js"
)

type recentEmoji struct {
	Content string `json:"content"`
	Code    string `json:"code,omitempty"`
	URL     string `json:"url,omitempty"`
}

const recentKey = "reaction-only-recent-emojis-v1"

func loadRecent() (items []recentEmoji) {
	defer func() { _ = recover() }()
	v := js.Global().Get("localStorage").Call("getItem", recentKey)
	if !v.IsNull() {
		_ = json.Unmarshal([]byte(v.String()), &items)
	}
	return
}
func rememberEmoji(content, code, url string) {
	defer func() { _ = recover() }()
	item := recentEmoji{content, code, url}
	items := []recentEmoji{item}
	for _, old := range loadRecent() {
		if old != item {
			items = append(items, old)
		}
		if len(items) >= 40 {
			break
		}
	}
	b, _ := json.Marshal(items)
	js.Global().Get("localStorage").Call("setItem", recentKey, string(b))
}
