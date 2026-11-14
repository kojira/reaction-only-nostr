package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"github.com/btcsuite/btcd/btcec/v2/schnorr"
)

type Event struct {
	ID        string     `json:"id"`
	Pubkey    string     `json:"pubkey"`
	CreatedAt int64      `json:"created_at"`
	Kind      int        `json:"kind"`
	Tags      [][]string `json:"tags"`
	Content   string     `json:"content"`
	Sig       string     `json:"sig"`
}

func verifyEvent(ev Event) bool {
	if ev.Tags == nil {
		return false
	}
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	enc.SetEscapeHTML(false)
	if enc.Encode([]any{0, ev.Pubkey, ev.CreatedAt, ev.Kind, ev.Tags, ev.Content}) != nil {
		return false
	}
	sum := sha256.Sum256(bytes.TrimSuffix(b.Bytes(), []byte("\n")))
	if hex.EncodeToString(sum[:]) != ev.ID {
		return false
	}
	pk, err := hex.DecodeString(ev.Pubkey)
	if err != nil || len(pk) != 32 {
		return false
	}
	key, err := schnorr.ParsePubKey(pk)
	if err != nil {
		return false
	}
	sb, err := hex.DecodeString(ev.Sig)
	if err != nil || len(sb) != 64 {
		return false
	}
	sig, err := schnorr.ParseSignature(sb)
	return err == nil && sig.Verify(sum[:], key)
}
