//go:build js && wasm

package main

import (
	"fmt"
	"syscall/js"
)

var parentEvents = map[string]Event{}
var parentRequests = map[string]string{} // subscription ID -> exact event ID
var parentSubscriptions = map[string]string{}

func requestParent(id string) {
	if !validHex(id, 64) {
		return
	}
	if _, ok := parentEvents[id]; ok {
		return
	}
	if _, ok := parentSubscriptions[id]; ok {
		return
	}
	sub := fmt.Sprintf("parent-%d", len(parentRequests)+1)
	parentRequests[sub] = id
	parentSubscriptions[id] = sub
	for _, ws := range sockets {
		subscribeParent(ws, sub, id)
	}
}
func subscribeParent(ws js.Value, sub, id string) {
	send(ws, []any{"REQ", sub, map[string]any{"ids": []string{id}, "limit": 1}})
}
func subscribeParents(ws js.Value) {
	for sub, id := range parentRequests {
		if _, ok := parentEvents[id]; !ok {
			subscribeParent(ws, sub, id)
		}
	}
}
func receiveParent(sub string, ev Event) bool {
	id, ok := parentRequests[sub]
	if !ok || ev.ID != id {
		return false
	}
	if _, ok := parentEvents[id]; ok {
		return false
	}
	parentEvents[id] = ev
	requestProfile(ev.Pubkey)
	for _, ws := range sockets {
		send(ws, []any{"CLOSE", sub})
	}
	return true
}
