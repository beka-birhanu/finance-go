package hash

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"sync"

	"github.com/beka-birhanu/finance-go/domain/common/hash"
	"golang.org/x/crypto/pbkdf2"
)

const (
	iterations = 10000
	saltSize   = 16
	keySize    = 32
)

// Service is an implementation of hash.IService for hashing and matching.
type Service struct{}

var _ hash.IService = &Service{}

var (
	instance *Service
	once     sync.Once
)

// SingletonService returns a singleton instance of the hash Service.
// It ensures that only one instance of the Service is created.
func SingletonService() *Service {
	once.Do(func() {
		instance = &Service{}
	})
	return instance
}

// Hash generates a hashed representation of the given word.
// It creates a random salt, combines it with the word, and hashes the result
// using PBKDF2 with SHA-256. The final result is the base64-encoded combination
// of the salt and the hash.
func (hs *Service) Hash(word string) (string, error) {
	salt := make([]byte, saltSize)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := pbkdf2.Key([]byte(word), salt, iterations, keySize, sha256.New)

	result := append(salt, hash...)

	return base64.StdEncoding.EncodeToString(result), nil
}

// Match compares a plain text word to a hashed word to determine if they match.
// It extracts the salt from the hashed word, re-hashes the plain text word with the same salt,
// and compares the result to the hash part of the hashed word.
//
// Parameters:
//   - hashedWord: The base64-encoded combination of salt and hash to be compared against.
//   - plainWord: The plain text word to be hashed and compared.
//
// Returns:
//   - A boolean indicating whether the plain text word matches the hashed word.
//   - An error if the hashed word is not in the expected format.
func (hs *Service) Match(hashedWord, plainWord string) (bool, error) {
	hashedWordBytes, err := base64.StdEncoding.DecodeString(hashedWord)
	if err != nil {
		return false, err
	}

	if len(hashedWordBytes) != saltSize+keySize {
		return false, errors.New("invalid hashed word length")
	}

	salt := hashedWordBytes[:saltSize]
	expectedHash := hashedWordBytes[saltSize:]

	hash := pbkdf2.Key([]byte(plainWord), salt, iterations, keySize, sha256.New)
	for i := 0; i < keySize; i++ {
		if expectedHash[i] != hash[i] {
			return false, nil
		}
	}

	return true, nil
}
