package main

import "testing"

func TestFollowing(t *testing.T) {
	p := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	f := following([][]string{{"p", p}, {"p", p}, {"p", "bad"}, {"e", p}})
	if len(f) != 1 || !f[p] {
		t.Fatal(f)
	}
}
func TestNewer(t *testing.T) {
	if !newer(Event{CreatedAt: 2}, Event{CreatedAt: 1}) || !newer(Event{CreatedAt: 1, ID: "a"}, Event{CreatedAt: 1, ID: "b"}) {
		t.Fatal("order")
	}
}
