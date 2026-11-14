//go:build js && wasm

package main

import (
	"encoding/json"
	"strings"
	"syscall/js"
)

const emojiStorageKey = "reaction-only-default-emojis-v1"

var defaultEmojis = initialDefaultEmojis()

func loadDefaultEmojis() {
	// Storage access may be blocked by browser privacy settings.
	defer func() {
		if recover() != nil {
			text("emojiSettingsStatus", "保存データを読み込めませんでした")
		}
		el("defaultEmojis").Set("value", strings.Join(defaultEmojis, " "))
	}()
	stored := js.Global().Get("localStorage").Call("getItem", emojiStorageKey)
	if stored.IsNull() {
		return
	}
	var values []string
	if json.Unmarshal([]byte(stored.String()), &values) != nil || !validDefaultEmojis(values) {
		text("emojiSettingsStatus", "保存データを読み込めませんでした")
		return
	}
	defaultEmojis = values
}

func validDefaultEmojis(values []string) bool {
	if len(values) == 0 || len(values) > 50 {
		return false
	}
	for _, value := range values {
		if !validEmoji(value) {
			return false
		}
	}
	return true
}

func saveDefaultEmojis() {
	values := strings.Fields(el("defaultEmojis").Get("value").String())
	if !validDefaultEmojis(values) {
		text("emojiSettingsStatus", "絵文字を空白で区切って1〜50個入力してください")
		return
	}
	raw, _ := json.Marshal(values)
	saved := false
	func() {
		defer func() { _ = recover() }()
		js.Global().Get("localStorage").Call("setItem", emojiStorageKey, string(raw))
		saved = true
	}()
	if !saved {
		text("emojiSettingsStatus", "保存できませんでした。ブラウザの設定を確認してください")
		return
	}
	defaultEmojis = values
	el("defaultEmojis").Set("value", strings.Join(values, " "))
	text("emojiSettingsStatus", "保存しました")
	renderTimeline()
}
