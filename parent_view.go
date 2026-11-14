//go:build js && wasm

package main

import (
	"fmt"
	"syscall/js"
)

var expandedParents = map[string]bool{}

// Show one parent at a time; expanding a reaction lets users follow the chain.
func renderParent(card js.Value, ev Event) {
	id := reactionParent(ev)
	if id == "" {
		return
	}
	box := doc.Call("createElement", "div")
	box.Set("className", "parent")
	if !validHex(id, 64) {
		box.Set("textContent", "参照先が不正です")
		card.Call("appendChild", box)
		return
	}
	parent, ok := parentEvents[id]
	if !ok {
		box.Set("textContent", "参照先を取得中 / 未取得: "+id[:12])
	} else {
		name := parseProfile(profiles[parent.Pubkey].Content, parent.Pubkey).Name
		label := fmt.Sprintf("%s · kind %d", name, parent.Kind)
		meta := doc.Call("createElement", "small")
		meta.Set("textContent", label)
		box.Call("appendChild", meta)
		content := doc.Call("createElement", "div")
		content.Set("className", "parent-content")
		content.Set("textContent", parent.Content)
		box.Call("appendChild", content)
		if parent.Kind == 7 && reactionParent(parent) != "" {
			if expandedParents[id] {
				renderParent(box, parent)
			} else {
				button := doc.Call("createElement", "button")
				button.Set("textContent", "さらに参照先を見る")
				button.Call("addEventListener", "click", callback(func(_ js.Value, _ []js.Value) any {
					expandedParents[id] = true
					requestParent(reactionParent(parent))
					renderTimeline()
					return nil
				}))
				box.Call("appendChild", button)
			}
		}
	}
	card.Call("appendChild", box)
}
