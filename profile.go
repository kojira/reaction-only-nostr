package main

import (
	"encoding/json"
	"net/url"
	"strings"
	"unicode"
)

type profile struct{ Name, Picture string }

func safeHTTPS(raw string) bool {
	u, err := url.Parse(raw)
	return err == nil && u.Scheme == "https" && u.Hostname() != "" && u.User == nil
}
func parseProfile(raw, pubkey string) profile {
	fallback := pubkey
	if len(fallback) > 12 {
		fallback = fallback[:12]
	}
	p := profile{Name: fallback}
	var v struct {
		DisplayName string `json:"display_name"`
		Name        string `json:"name"`
		Picture     string `json:"picture"`
	}
	if json.Unmarshal([]byte(raw), &v) != nil {
		return p
	}
	for _, name := range []string{v.DisplayName, v.Name} {
		name = strings.TrimSpace(name)
		if name != "" && len([]rune(name)) <= 100 && strings.IndexFunc(name, unicode.IsControl) < 0 {
			p.Name = name
			break
		}
	}
	if safeHTTPS(v.Picture) {
		p.Picture = v.Picture
	}
	return p
}
