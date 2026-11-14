//go:build js && wasm

package main

import "syscall/js"

const loginStorageKey = "reaction-only-login-pubkey-v1"

var loginGeneration int

func storedLogin() (result string) {
	defer func() { _ = recover() }()
	v := js.Global().Get("localStorage").Call("getItem", loginStorageKey)
	if !v.IsNull() && validHex(v.String(), 64) {
		return v.String()
	}
	return ""
}
func saveLogin(p string) {
	defer func() { _ = recover() }()
	js.Global().Get("localStorage").Call("setItem", loginStorageKey, p)
}
func clearStoredLogin() {
	defer func() { _ = recover() }()
	js.Global().Get("localStorage").Call("removeItem", loginStorageKey)
}
