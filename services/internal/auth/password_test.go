package auth

import "testing"

func TestCheckPassword(t *testing.T) {
	for p, want := range map[string]error{
		"flowscore1":             nil,
		"короткий1":              nil, // 9 letters: runes count, not bytes
		"short1":                 ErrWeakPassword,
		"no digits here":         ErrWeakPassword,
		string(make([]byte, 73)): ErrLongPassword,
	} {
		if got := CheckPassword(p); got != want {
			t.Errorf("%q: %v, want %v", p, got, want)
		}
	}
}

func TestHashPassword(t *testing.T) {
	a, err := HashPassword("flowscore1")
	if err != nil {
		t.Fatal(err)
	}
	b, _ := HashPassword("flowscore1")
	// Salted: one password, two hashes, and both check out
	if a == b || !SamePassword(a, "flowscore1") || !SamePassword(b, "flowscore1") {
		t.Errorf("hashes %q and %q", a, b)
	}
	if SamePassword(a, "flowscore2") {
		t.Error("a wrong password checked out")
	}
}

func TestNormalizeEmail(t *testing.T) {
	if got := NormalizeEmail("  Ayan@Mail.COM "); got != "ayan@mail.com" {
		t.Errorf("got %q", got)
	}
}
