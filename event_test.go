package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
)

func TestVerifyEvent(t *testing.T) {
	key, _ := btcec.PrivKeyFromBytes(bytes.Repeat([]byte{1}, 32))
	ev := Event{Pubkey: hex.EncodeToString(schnorr.SerializePubKey(key.PubKey())), CreatedAt: 123, Kind: 7, Tags: [][]string{{"e", "abc"}}, Content: "🌸 <x>"}
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	enc.SetEscapeHTML(false)
	enc.Encode([]any{0, ev.Pubkey, ev.CreatedAt, ev.Kind, ev.Tags, ev.Content})
	hash := sha256.Sum256(bytes.TrimSuffix(b.Bytes(), []byte("\n")))
	ev.ID = hex.EncodeToString(hash[:])
	sig, err := schnorr.Sign(key, hash[:])
	if err != nil {
		t.Fatal(err)
	}
	ev.Sig = hex.EncodeToString(sig.Serialize())
	if !verifyEvent(ev) {
		t.Fatal("valid event rejected")
	}
	ev.Content = "modified"
	if verifyEvent(ev) {
		t.Fatal("tampered event accepted")
	}
}
