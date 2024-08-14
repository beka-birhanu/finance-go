// Package hash provides secure hashing and matching functionalities using PBKDF2 with SHA-256.
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

// Service implements hash.IService for hashing and matching.
type Service struct{}

var _ hash.IService = &Service{}

var (
	instance *Service
	once     sync.Once
)

// SingletonService returns a singleton instance of Service.
func SingletonService() *Service {
	once.Do(func() {
		instance = &Service{}
	})
	return instance
}

// Hash generates a hashed representation of the given word.
func (hs *Service) Hash(word string) (string, error) {
	salt := make([]byte, saltSize)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := pbkdf2.Key([]byte(word), salt, iterations, keySize, sha256.New)

	result := append(salt, hash...)
	return base64.StdEncoding.EncodeToString(result), nil
}

// Match compares a plain word to a hashed word.
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
