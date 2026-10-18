package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Tokens issues and reads the access tokens: a JWT that names the user and
// the moment it stops counting, signed with the server's secret. Anyone can
// read what it says; only the server can make one that checks out.
type Tokens struct {
	Secret []byte
	TTL    time.Duration // how long a token counts
}

// Issue makes a token for the user, counting from now.
func (t Tokens) Issue(userID string, now time.Time) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(t.TTL)),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(t.Secret)
}

// UserID reads a token back: the user it was issued to, or an error for a
// token that is forged, broken or expired.
func (t Tokens) UserID(token string) (string, error) {
	var claims jwt.RegisteredClaims
	_, err := jwt.ParseWithClaims(token, &claims,
		func(*jwt.Token) (any, error) { return t.Secret, nil },
		// Only our own algorithm: a token that says "none" or another one
		// is turned away before its signature is looked at
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return "", err
	}
	return claims.Subject, nil
}
