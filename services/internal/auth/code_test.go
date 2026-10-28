package auth

import (
	"regexp"
	"testing"
)

func TestNewCode(t *testing.T) {
	seen := map[string]bool{}
	for range 50 {
		code, err := NewCode()
		if err != nil {
			t.Fatal(err)
		}
		if !regexp.MustCompile(`^\d{6}$`).MatchString(code) {
			t.Fatalf("code %q is not six digits", code)
		}
		seen[code] = true
	}
	// 50 codes out of a million: a repeat or two at most
	if len(seen) < 48 {
		t.Errorf("only %d different codes of 50", len(seen))
	}
}

func TestCodesHash(t *testing.T) {
	c := Codes{Secret: []byte("a secret of at least thirty-two bytes")}
	h := c.Hash("u1", "verify", "123456")
	if h != c.Hash("u1", "verify", "123456") {
		t.Error("the same code hashed twice differs")
	}
	// Another user, purpose, code or secret: another hash
	for _, other := range []string{
		c.Hash("u2", "verify", "123456"),
		c.Hash("u1", "reset", "123456"),
		c.Hash("u1", "verify", "123457"),
		Codes{Secret: []byte("another secret, also thirty-two bytes")}.Hash("u1", "verify", "123456"),
	} {
		if other == h {
			t.Error("two different codes share a hash")
		}
	}
}
