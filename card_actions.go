//go:build js && wasm

package main

import "syscall/js"

func buildReactionActions(actions js.Value, ev Event) {
	button := doc.Call("createElement", "button")
	button.Set("textContent", defaultEmojis[0])
	button.Set("ariaLabel", "デフォルトのリアクションを送る")
	button.Call("addEventListener", "click", callback(func(_ js.Value, _ []js.Value) any {
		publishReaction(ev.ID, ev.Pubkey, defaultEmojis[0], "", "")
		return nil
	}))
	actions.Call("appendChild", button)
	toggle := doc.Call("createElement", "button")
	toggle.Set("textContent", "+")
	toggle.Set("ariaLabel", "絵文字を選ぶ")
	toggle.Call("setAttribute", "aria-expanded", "false")
	actions.Call("appendChild", toggle)
	choices := doc.Call("createElement", "div")
	choices.Set("className", "emoji-choices")
	choices.Set("hidden", true)
	toggle.Call("addEventListener", "click", callback(func(_ js.Value, _ []js.Value) any {
		hidden := !choices.Get("hidden").Bool()
		choices.Set("hidden", hidden)
		toggle.Call("setAttribute", "aria-expanded", !hidden)
		return nil
	}))
	for _, emoji := range defaultEmojis {
		emoji := emoji
		b := doc.Call("createElement", "button")
		b.Set("textContent", emoji)
		b.Set("title", "このリアクションに "+emoji+" を送る")
		b.Call("addEventListener", "click", callback(func(_ js.Value, _ []js.Value) any { publishReaction(ev.ID, ev.Pubkey, emoji, "", ""); return nil }))
		choices.Call("appendChild", b)
	}
	for _, t := range setTags {
		if len(t) >= 3 && t[0] == "emoji" && safeHTTPS(t[2]) {
			tag := append([]string(nil), t...)
			b := doc.Call("createElement", "button")
			b.Set("textContent", ":"+tag[1]+":")
			b.Set("title", "このリアクションに :"+tag[1]+": を送る")
			b.Call("addEventListener", "click", callback(func(_ js.Value, _ []js.Value) any { publishReaction(ev.ID, ev.Pubkey, "", tag[1], tag[2]); return nil }))
			choices.Call("appendChild", b)
		}
	}
	actions.Call("appendChild", choices)
}
