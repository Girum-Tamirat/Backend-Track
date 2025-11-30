package Infrastructure

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword returns bcrypt hash
func HashPassword(pw string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(b), err
}

// ComparePassword compares hash with plaintext
func ComparePassword(hash, pw string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw))
}
