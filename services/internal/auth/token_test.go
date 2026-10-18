package auth

import (
	"testing"
	"time"
)

func TestTokens(t *testing.T) {
	tokens := Tokens{Secret: []byte("a secret of at least thirty-two bytes"), TTL: 15 * time.Minute}

	token, err := tokens.Issue("u1", time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if id, err := tokens.UserID(token); err != nil || id != "u1" {
		t.Errorf("fresh token: %q, %v", id, err)
	}

	// Issued half an hour ago, counted for 15 minutes
	old, _ := tokens.Issue("u1", time.Now().Add(-30*time.Minute))
	if _, err := tokens.UserID(old); err == nil {
		t.Error("an expired token checked out")
	}

	// Signed with another secret: someone made it up
	other := Tokens{Secret: []byte("another secret, also thirty-two bytes"), TTL: 15 * time.Minute}
	forged, _ := other.Issue("u1", time.Now())
	if _, err := tokens.UserID(forged); err == nil {
		t.Error("a forged token checked out")
	}
}
