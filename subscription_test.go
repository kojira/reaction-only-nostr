package main

import (
	"fmt"
	"testing"
)

func TestProfileSubscriptionIDs(t *testing.T) {
	ids := map[string]bool{}
	for i := 0; i < 10000; i++ {
		key := fmt.Sprintf("%064x", i)
		id := profileSubscriptionID(key)
		if len(id) > 64 || ids[id] || id != profileSubscriptionID(key) {
			t.Fatalf("invalid or duplicate subscription ID %q", id)
		}
		ids[id] = true
	}
}
