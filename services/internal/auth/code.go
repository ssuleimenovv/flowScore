package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
)

// NewCode is a random 6-digit code. It comes from crypto/rand: math/rand
// would let someone who saw a few codes work out the next ones
func NewCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1_000_000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

// Codes hashes the codes the database keeps. It is an HMAC with the server's
// secret: a plain hash of six digits is undone by trying all million, so a
// copy of the table alone would give the codes away
type Codes struct {
	Secret []byte
}

// Hash is the code as the database keeps it, tied to the user and the purpose
func (c Codes) Hash(userID, purpose, code string) string {
	m := hmac.New(sha256.New, c.Secret)
	m.Write([]byte(userID + ":" + purpose + ":" + code))
	return hex.EncodeToString(m.Sum(nil))
}
