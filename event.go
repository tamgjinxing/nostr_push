package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
)

// Tag is a tag of an event.
type Tag []string
type EventKind int64

// Event is an Nostr event.
type Event struct {
	ID        string    `json:"id"`
	PubKey    string    `json:"pubkey"`
	CreatedAt int64     `json:"created_at"`
	Kind      EventKind `json:"kind"`
	Tags      []Tag     `json:"tags"`
	Content   string    `json:"content"`
	Sig       string    `json:"sig"`
}

type EventMessage struct {
	SubscriptionID string // optional
	Event          *Event
}

// Sing signs the event with the given private key.
// It sets the ID, PubKey, and Sig fields.
func (e *Event) Sign(privKey string) error {
	s, err := hex.DecodeString(privKey)
	if err != nil {
		return fmt.Errorf("invalid private key: %w", err)
	}
	sk, pk := btcec.PrivKeyFromBytes(s)

	// public key
	e.PubKey = hex.EncodeToString(pk.SerializeCompressed()[1:])

	serial, err := e.serialize()
	if err != nil {
		return fmt.Errorf("invalid event: %w", err)
	}
	serialHash := sha256.Sum256(serial)

	// id
	e.ID = hex.EncodeToString(serialHash[:])

	// signature
	sig, err := schnorr.Sign(sk, serialHash[:])
	if err != nil {
		return err
	}
	e.Sig = hex.EncodeToString(sig.Serialize())
	return nil
}

func (e *Event) serialize() ([]byte, error) {
	b, err := json.Marshal([]any{
		0,
		e.PubKey,
		e.CreatedAt,
		e.Kind,
		e.Tags,
		e.Content,
	})
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (e *Event) ToEventString() (string, error) {
	result := []interface{}{"EVENT", e}
	resultJSON, err := json.Marshal(result)
	if err != nil {
		log.Println("Error marshalling result JSON:", err)
		return "", err
	}

	return string(resultJSON), nil
}

func (e *Event) ToAuthString() (string, error) {
	result := []interface{}{"AUTH", e}
	resultJSON, err := json.Marshal(result)
	if err != nil {
		log.Println("Error marshalling result JSON:", err)
		return "", err
	}

	return string(resultJSON), nil
}
