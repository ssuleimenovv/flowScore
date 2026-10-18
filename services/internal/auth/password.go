package auth

import (
	"errors"
	"strings"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrWeakPassword: the rule the sign-up form tells, "Не короче 8
	// символов и хотя бы одна цифра"
	ErrWeakPassword = errors.New("password: at least 8 characters, one of them a digit")
	// ErrLongPassword: bcrypt reads only the first 72 bytes, so a longer
	// password would match anything that starts the same
	ErrLongPassword = errors.New("password: at most 72 bytes")
)

// CheckPassword tells whether a new password follows the rule.
func CheckPassword(p string) error {
	if len(p) > 72 {
		return ErrLongPassword
	}
	if utf8.RuneCountInString(p) < 8 || !strings.ContainsAny(p, "0123456789") {
		return ErrWeakPassword
	}
	return nil
}

// HashPassword is what the database keeps instead of the password. bcrypt
// salts it, so two users with one password get different hashes.
func HashPassword(p string) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	return string(h), err
}

// SamePassword tells whether p is the password the hash was made from.
func SamePassword(hash, p string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(p)) == nil
}

// NormalizeEmail is the email as the database keeps it.
func NormalizeEmail(e string) string {
	return strings.ToLower(strings.TrimSpace(e))
}
