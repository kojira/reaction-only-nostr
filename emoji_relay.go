//go:build js && wasm

package main

import (
	"strings"
	"syscall/js"
)

var emojiList Event
var emojiSets = map[string]Event{}

func resetEmojiLists() { emojiList = Event{}; emojiSets = map[string]Event{} }
func subscribeEmoji(ws js.Value) {
	send(ws, []any{"REQ", "emoji-list", map[string]any{"kinds": []int{10030}, "authors": []string{pubkey}, "limit": 1}})
	send(ws, []any{"REQ", "emoji-set", map[string]any{"kinds": []int{30030}, "authors": []string{pubkey}, "#d": []string{"reaction-only"}}})
	for _, a := range emojiReferences(emojiList.Tags) {
		subscribeEmojiRef(ws, a)
	}
}
func subscribeEmojiRef(ws js.Value, a string) {
	p := strings.SplitN(a, ":", 3)
	send(ws, []any{"REQ", "emoji-ref-" + a, map[string]any{"kinds": []int{30030}, "authors": []string{p[1]}, "#d": []string{p[2]}, "limit": 1}})
}
func handleEmojiEvent(ev Event, sub string) {
	if sub == "emoji-list" && ev.Kind == 10030 && ev.Pubkey == pubkey && (emojiList.ID == "" || newer(ev, emojiList)) {
		emojiList = ev
		emojiSets = map[string]Event{}
		for _, a := range emojiReferences(ev.Tags) {
			for _, ws := range sockets {
				subscribeEmojiRef(ws, a)
			}
		}
		return
	}
	if ev.Kind != 30030 || !strings.HasPrefix(sub, "emoji-ref-") {
		return
	}
	a := "30030:" + ev.Pubkey + ":" + emojiAddress(ev.Tags)
	if sub != "emoji-ref-"+a {
		return
	}
	found := false
	for _, ref := range emojiReferences(emojiList.Tags) {
		if ref == a {
			found = true
			break
		}
	}
	if !found {
		return
	}
	old, ok := emojiSets[a]
	if !ok || newer(ev, old) {
		emojiSets[a] = ev
	}
}
func availableEmojiTags() [][]string {
	tags := append([][]string{}, setTags...)
	tags = append(tags, emojiList.Tags...)
	for _, a := range emojiReferences(emojiList.Tags) {
		if e, ok := emojiSets[a]; ok {
			tags = append(tags, e.Tags...)
		}
	}
	return tags
}
