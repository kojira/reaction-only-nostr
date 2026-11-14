//go:build js && wasm

package main

import (
	"syscall/js"
	"time"
)

var retryPending bool
var hiddenAt time.Time

func scheduleReconnect() {
	if retryPending {
		return
	}
	retryPending = true
	js.Global().Call("setTimeout", callback(func(_ js.Value, _ []js.Value) any {
		retryPending = false
		if doc.Get("visibilityState").String() != "visible" {
			return nil
		}
		for _, ws := range sockets {
			if ws.Get("readyState").Int() > 1 {
				status("再接続中…")
				connect()
				break
			}
		}
		return nil
	}), 3000)
}

func setupRecovery() {
	doc.Call("addEventListener", "visibilitychange", callback(func(_ js.Value, _ []js.Value) any {
		if doc.Get("visibilityState").String() == "hidden" {
			hiddenAt = time.Now()
			return nil
		}
		// A suspended socket can still report OPEN: rebuild after a background interval.
		if !hiddenAt.IsZero() && time.Since(hiddenAt) > 5*time.Second {
			hiddenAt = time.Time{}
			status("再接続中…")
			connect()
			return nil
		}
		for _, ws := range sockets {
			if ws.Get("readyState").Int() > 1 {
				scheduleReconnect()
				break
			}
		}
		return nil
	}))
	js.Global().Call("addEventListener", "pageshow", callback(func(_ js.Value, a []js.Value) any {
		if len(a) > 0 && a[0].Get("persisted").Bool() {
			status("再接続中…")
			connect()
		}
		return nil
	}))
}
