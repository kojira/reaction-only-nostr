package main

// NIP-25 identifies the reacted-to event with the last e tag.
func reactionParent(ev Event) string {
	if ev.Kind != 7 {
		return ""
	}
	return parentID(ev)
}
