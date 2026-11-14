package main

import (
	"reflect"
	"testing"
)

func TestInitialDefaultEmojis(t *testing.T) {
	want := []string{"❤️", "👍", "🎉", "😮"}
	if got := initialDefaultEmojis(); !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}
