//go:build js && wasm

package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"syscall/js"
	"time"
)

type Event struct {
	ID        string     `json:"id"`
	Pubkey    string     `json:"pubkey"`
	CreatedAt int64      `json:"created_at"`
	Kind      int        `json:"kind"`
	Tags      [][]string `json:"tags"`
	Content   string     `json:"content"`
	Sig       string     `json:"sig"`
}

var doc = js.Global().Get("document")
var callbacks []js.Func
var sockets []js.Value
var pubkey string
var seen = map[string]bool{}
var setTags = [][]string{{"d", "reaction-only"}}
var setTime int64

func el(id string) js.Value  { return doc.Call("getElementById", id) }
func value(id string) string { return strings.TrimSpace(el(id).Get("value").String()) }
func text(id, s string)      { el(id).Set("textContent", s) }
func callback(f func(js.Value, []js.Value) any) js.Func {
	x := js.FuncOf(f)
	callbacks = append(callbacks, x)
	return x
}
func listen(id, event string, f func(js.Value, []js.Value) any) {
	el(id).Call("addEventListener", event, callback(f))
}
func status(s string) { text("status", s) }
func promise(p js.Value, ok func(js.Value)) {
	p.Call("then", callback(func(_ js.Value, a []js.Value) any { ok(a[0]); return nil })).Call("catch", callback(func(_ js.Value, a []js.Value) any { status("エラー: " + a[0].String()); return nil }))
}
func jsonJS(v any) js.Value {
	b, _ := json.Marshal(v)
	return js.Global().Get("JSON").Call("parse", string(b))
}
func send(ws js.Value, v any) {
	b, e := json.Marshal(v)
	if e == nil && ws.Get("readyState").Int() == 1 {
		ws.Call("send", string(b))
	}
}
func connect() {
	for _, s := range sockets {
		s.Call("close")
	}
	sockets = nil
	for _, u := range relayURLs {
		func(url string) {
			ws := js.Global().Get("WebSocket").New(url)
			sockets = append(sockets, ws)
			ws.Set("onopen", callback(func(_ js.Value, _ []js.Value) any {
				send(ws, []any{"REQ", "reactions", map[string]any{"kinds": []int{7}, "limit": 100}})
				status("接続: " + url)
				return nil
			}))
			ws.Set("onmessage", callback(func(_ js.Value, a []js.Value) any {
				var msg []json.RawMessage
				if json.Unmarshal([]byte(a[0].Get("data").String()), &msg) != nil || len(msg) < 3 {
					return nil
				}
				var typ string
				json.Unmarshal(msg[0], &typ)
				if typ == "EVENT" {
					var ev Event
					if json.Unmarshal(msg[2], &ev) == nil && ev.Kind == 30030 && ev.Pubkey == pubkey && ev.CreatedAt > setTime {
						for _, t := range ev.Tags {
							if len(t) > 1 && t[0] == "d" && t[1] == "reaction-only" {
								setTime = ev.CreatedAt
								setTags = ev.Tags
								refreshSet()
								break
							}
						}
					}
					if ev.Kind == 7 && validHex(ev.ID, 64) && validHex(ev.Pubkey, 64) && !seen[ev.ID] {
						seen[ev.ID] = true
						render(ev)
					}
				}
				if typ == "OK" {
					var accepted bool
					json.Unmarshal(msg[2], &accepted)
					status(fmt.Sprintf("リレー応答 %s: 受理=%v", url, accepted))
				}
				return nil
			}))
			ws.Set("onerror", callback(func(_ js.Value, _ []js.Value) any { status("接続エラー: " + url); return nil }))
		}(u)
	}
}
func validHex(s string, n int) bool {
	if len(s) != n {
		return false
	}
	for _, r := range s {
		if !(r >= '0' && r <= '9' || r >= 'a' && r <= 'f') {
			return false
		}
	}
	return true
}
func target(ev Event) (string, string) {
	var e, p string
	for _, t := range ev.Tags {
		if len(t) >= 2 {
			if t[0] == "e" {
				e = t[1]
			}
			if t[0] == "p" {
				p = t[1]
			}
		}
	}
	return e, p
}
func render(ev Event) {
	card := doc.Call("createElement", "article")
	line := doc.Call("createElement", "div")
	line.Set("className", "reaction")
	custom := false
	for _, t := range ev.Tags {
		if len(t) >= 3 && t[0] == "emoji" && ev.Content == ":"+t[1]+":" && strings.HasPrefix(t[2], "https://") {
			img := doc.Call("createElement", "img")
			img.Set("src", t[2])
			img.Set("alt", ev.Content)
			img.Set("referrerPolicy", "no-referrer")
			line.Call("appendChild", img)
			custom = true
			break
		}
	}
	if !custom {
		line.Set("textContent", ev.Content)
	}
	card.Call("appendChild", line)
	meta := doc.Call("createElement", "small")
	meta.Set("textContent", ev.Pubkey[:12]+" · "+time.Unix(ev.CreatedAt, 0).Format("2006-01-02 15:04"))
	card.Call("appendChild", meta)
	e, p := target(ev)
	if validHex(e, 64) && validHex(p, 64) {
		btn := doc.Call("createElement", "button")
		btn.Set("textContent", "この対象にリアクション")
		btn.Call("addEventListener", "click", callback(func(_ js.Value, _ []js.Value) any {
			el("targetEvent").Set("value", e)
			el("targetPubkey").Set("value", p)
			return nil
		}))
		card.Call("appendChild", btn)
	}
	btn := doc.Call("createElement", "button")
	btn.Set("textContent", "このリアクションに返信")
	btn.Call("addEventListener", "click", callback(func(_ js.Value, _ []js.Value) any {
		el("targetEvent").Set("value", ev.ID)
		el("targetPubkey").Set("value", ev.Pubkey)
		return nil
	}))
	card.Call("appendChild", btn)
	el("feed").Call("prepend", card)
}
func login() {
	n := js.Global().Get("nostr")
	if n.IsUndefined() || n.IsNull() {
		status("NIP-07拡張機能が必要です")
		return
	}
	promise(n.Call("getPublicKey"), func(v js.Value) {
		p := v.String()
		if !validHex(p, 64) {
			status("公開鍵が不正です")
			return
		}
		pubkey = p
		text("account", "ログイン: "+p)
		for _, ws := range sockets {
			send(ws, []any{"REQ", "emoji-set", map[string]any{"kinds": []int{30030}, "authors": []string{pubkey}, "#d": []string{"reaction-only"}}})
		}
		status("ログインしました")
	})
}
func publish() {
	if pubkey == "" {
		status("先にログインしてください")
		return
	}
	id, p := value("targetEvent"), value("targetPubkey")
	if !validHex(id, 64) || !validHex(p, 64) {
		status("対象イベントIDと公開鍵は64桁の小文字hexで指定してください")
		return
	}
	content := value("emoji")
	shortcode, url := value("shortcode"), value("emojiURL")
	tags := [][]string{{"e", id}, {"p", p}}
	if shortcode != "" || url != "" {
		if shortcode == "" || !strings.HasPrefix(url, "https://") {
			status("カスタム絵文字はショートコードとhttps URLが必要です")
			return
		}
		for _, r := range shortcode {
			if !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '_') {
				status("ショートコードは英数字と_のみです")
				return
			}
		}
		content = ":" + shortcode + ":"
		tags = append(tags, []string{"emoji", shortcode, url})
	}
	if content == "" {
		status("絵文字を入力してください")
		return
	}
	count := 0
	for _, ws := range sockets {
		if ws.Get("readyState").Int() == 1 {
			count++
		}
	}
	if count == 0 {
		status("接続中のリレーがありません")
		return
	}
	unsigned := map[string]any{"kind": 7, "pubkey": pubkey, "created_at": time.Now().Unix(), "tags": tags, "content": content}
	promise(js.Global().Get("nostr").Call("signEvent", jsonJS(unsigned)), func(signed js.Value) {
		var ev Event
		raw := js.Global().Get("JSON").Call("stringify", signed).String()
		if json.Unmarshal([]byte(raw), &ev) != nil || ev.Pubkey != pubkey || ev.Kind != 7 || !validHex(ev.ID, 64) || !validHex(ev.Sig, 128) {
			status("署名結果が不正です")
			return
		}
		for _, ws := range sockets {
			send(ws, []any{"EVENT", signedToMap(raw)})
		}
		status("送信しました。リレーの受理応答を待っています")
	})
}
func signedToMap(s string) any { var v any; json.Unmarshal([]byte(s), &v); return v }
func main() {
	listen("connect", "click", func(_ js.Value, _ []js.Value) any { connect(); return nil })
	listen("login", "click", func(_ js.Value, _ []js.Value) any { login(); return nil })
	listen("register", "click", func(_ js.Value, _ []js.Value) any { registerSet(); return nil })
	listen("publish", "click", func(_ js.Value, _ []js.Value) any { publish(); return nil })
	connect()
	status("起動しました")
	select {}
}

func refreshSet() {
	box := el("emojiSet")
	box.Set("textContent", "")
	for _, t := range setTags {
		if len(t) >= 3 && t[0] == "emoji" {
			tag := t
			b := doc.Call("createElement", "button")
			b.Set("textContent", ":"+tag[1]+":")
			b.Call("addEventListener", "click", callback(func(_ js.Value, _ []js.Value) any {
				el("shortcode").Set("value", tag[1])
				el("emojiURL").Set("value", tag[2])
				return nil
			}))
			box.Call("appendChild", b)
		}
	}
}
func registerSet() {
	code, url := value("shortcode"), value("emojiURL")
	if pubkey == "" || code == "" || !strings.HasPrefix(url, "https://") {
		status("ログイン、ショートコード、https画像URLが必要です")
		return
	}
	for _, r := range code {
		if !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '_') {
			status("ショートコードは英数字と_のみです")
			return
		}
	}
	tags := [][]string{}
	for _, t := range setTags {
		if len(t) >= 2 && t[0] == "emoji" && t[1] == code {
			continue
		}
		tags = append(tags, t)
	}
	tags = append(tags, []string{"emoji", code, url})
	event := map[string]any{"kind": 30030, "pubkey": pubkey, "created_at": time.Now().Unix(), "tags": tags, "content": ""}
	promise(js.Global().Get("nostr").Call("signEvent", jsonJS(event)), func(signed js.Value) {
		for _, ws := range sockets {
			send(ws, []any{"EVENT", signedToMap(js.Global().Get("JSON").Call("stringify", signed).String())})
		}
		status("絵文字セットの公開を試みました。リレー応答を確認してください")
	})
}
