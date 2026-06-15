package hash

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Hasher abstracts password hashing so it can be swapped in tests.
type Hasher interface {
	Hash(password string) (string, error)
	Verify(password, hashed string) bool
}

type bcryptHasher struct {
	cost int
}

// NewBcrypt returns a bcrypt hasher with the recommended cost (12).
// Cost 12 is ~250ms on modern hardware — suitable for production.
func NewBcrypt() Hasher {
	return &bcryptHasher{cost: 12}
}

func (h *bcryptHasher) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", fmt.Errorf("hash: %w", err)
	}
	return string(bytes), nil
}

func (h *bcryptHasher) Verify(password, hashed string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)) == nil
}
