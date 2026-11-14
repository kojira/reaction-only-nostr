package main

import (
	"reflect"
)

// Accept one emoji sequence, including modifiers, flags and joined emoji.
func validEmoji(s string) bool {
	r := []rune(s)
	if len(r) == 0 || len(r) > 32 {
		return false
	}
	bases := 0
	joined := false
	for i, c := range r {
		switch {
		case c == 0x200d:
			if i == 0 || i == len(r)-1 || joined {
				return false
			}
			joined = true
		case c == 0xfe0f || c == 0xfe0e || c >= 0x1f3fb && c <= 0x1f3ff:
			if bases == 0 || joined {
				return false
			}
		case c >= 0x1f1e6 && c <= 0x1f1ff:
			if joined || bases >= 2 {
				return false
			}
			bases++
		case c >= 0x1f000 && c <= 0x1faff || c >= 0x2600 && c <= 0x27bf || c == 0x00a9 || c == 0x00ae || c == 0x203c || c == 0x2049 || c == 0x2122 || c == 0x2139 || c == 0x3030 || c == 0x303d || c == 0x3297 || c == 0x3299:
			if bases > 0 && !joined {
				return false
			}
			bases++
			joined = false
		default:
			return false
		}
	}
	return bases > 0 && !joined
}
func sameTags(a, b [][]string) bool { return reflect.DeepEqual(a, b) }
