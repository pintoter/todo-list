package hash

import (
	"crypto/sha1"
	"fmt"
)

type Hasher interface {
	Hash(password string) (string, error)
}

type Config interface {
	GetSalt() string
}

type SHA1Hash struct {
	salt string
}

func New(cfg Config) *SHA1Hash {
	return &SHA1Hash{
		salt: cfg.GetSalt(),
	}
}

func (h *SHA1Hash) Hash(password string) (string, error) {
	hasher := sha1.New()

	if _, err := hasher.Write([]byte(password)); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hasher.Sum([]byte(h.salt))), nil
}
