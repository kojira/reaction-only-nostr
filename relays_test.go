package main

import (
	"reflect"
	"testing"
)

func TestRelayURLs(t *testing.T) {
	want := [...]string{"wss://r.kojira.io", "wss://x.kojira.io", "wss://n.kojira.io"}
	if !reflect.DeepEqual(relayURLs, want) {
		t.Fatalf("relay URLs = %v, want %v", relayURLs, want)
	}
}
