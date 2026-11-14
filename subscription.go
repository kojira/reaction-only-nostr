package main

import "strconv"

// Profile subscriptions need a unique, stable ID of at most 64 characters.
// A monotonically allocated integer avoids collisions between public keys.
var profileSubscriptionIDs = map[string]string{}

func profileSubscriptionID(key string) string {
	if id, ok := profileSubscriptionIDs[key]; ok {
		return id
	}
	id := "p" + strconv.Itoa(len(profileSubscriptionIDs)+1)
	profileSubscriptionIDs[key] = id
	return id
}
