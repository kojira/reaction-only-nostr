package main

func following(tags [][]string) map[string]bool {
	result := map[string]bool{}
	for _, t := range tags {
		if len(t) > 1 && t[0] == "p" && validPubkey(t[1]) {
			result[t[1]] = true
		}
	}
	return result
}
func validPubkey(s string) bool {
	if len(s) != 64 {
		return false
	}
	for _, r := range s {
		if !(r >= '0' && r <= '9' || r >= 'a' && r <= 'f') {
			return false
		}
	}
	return true
}
func newer(a, b Event) bool {
	return a.CreatedAt > b.CreatedAt || a.CreatedAt == b.CreatedAt && a.ID < b.ID
}
