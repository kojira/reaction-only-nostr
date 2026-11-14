package main

import "testing"

func TestParseProfile(t *testing.T) {
	cases := []struct{ raw, name, picture string }{
		{`{"display_name":"Alice","name":"a","picture":"https://example.org/a.png"}`, "Alice", "https://example.org/a.png"},
		{`{"display_name":"","name":"Bob","picture":"javascript:alert(1)"}`, "Bob", ""},
		{`bad`, "0123456789ab", ""},
		{`{"name":"\u0000bad","picture":"https://user@example.org/a"}`, "0123456789ab", ""},
	}
	for _, c := range cases {
		p := parseProfile(c.raw, "0123456789abcdef")
		if p.Name != c.name || p.Picture != c.picture {
			t.Errorf("%q: %+v", c.raw, p)
		}
	}
}
