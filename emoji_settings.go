//go:build js && wasm

package main

import (
	"encoding/json"
	"strings"
	"syscall/js"
)

const reactionStorageKey = "reaction-only-default-reaction-v1"
const emojiStorageKey = "reaction-only-default-emojis-v1"

var selectedReaction = defaultReaction{Content: "❤️"}

func loadDefaultReaction() {
	defer func() {
		if recover() != nil {
			text("emojiSettingsStatus", "保存データを読み込めませんでした")
		}
		updateDefaultFields()
	}()
	storage := js.Global().Get("localStorage")
	raw, old := "", ""
	if v := storage.Call("getItem", reactionStorageKey); !v.IsNull() {
		raw = v.String()
	}
	if v := storage.Call("getItem", emojiStorageKey); !v.IsNull() {
		old = v.String()
	}
	selectedReaction = chooseDefaultReaction(raw, old)
}
func updateDefaultFields() {
	el("defaultReaction").Set("value", selectedReaction.Content)
	el("defaultReactionURL").Set("value", selectedReaction.URL)
}
func saveDefaultReaction() {
	value := strings.TrimSpace(el("defaultReaction").Get("value").String())
	url := strings.TrimSpace(el("defaultReactionURL").Get("value").String())
	r := defaultReaction{Content: value, URL: url}
	if strings.HasPrefix(value, ":") && strings.HasSuffix(value, ":") && len(value) > 2 {
		r.Code = value[1 : len(value)-1]
	}
	if !validDefaultReaction(r) {
		text("emojiSettingsStatus", "通常絵文字を1つ、または :ショートコード: と画像のHTTPS URLを入力してください")
		return
	}
	raw, _ := json.Marshal(r)
	saved := false
	func() {
		defer func() { _ = recover() }()
		js.Global().Get("localStorage").Call("setItem", reactionStorageKey, string(raw))
		saved = true
	}()
	if !saved {
		text("emojiSettingsStatus", "保存できませんでした。ブラウザの設定を確認してください")
		return
	}
	selectedReaction = r
	updateDefaultFields()
	text("emojiSettingsStatus", "保存しました")
	renderTimeline()
}
