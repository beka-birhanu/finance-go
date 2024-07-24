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

type HashService struct{}

var _ hash.IHashService = &HashService{}

var (
	instance *HashService
	once     sync.Once
)

func GetHashService() *HashService {
	once.Do(func() {
		instance = &HashService{}
	})
	return instance
}

func (hs *HashService) Hash(word string) (string, error) {
	salt := make([]byte, saltSize)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := pbkdf2.Key([]byte(word), salt, iterations, keySize, sha256.New)

	result := append(salt, hash...)

	return base64.StdEncoding.EncodeToString(result), nil
}

func (hs *HashService) Match(hashedWord, plainWord string) (bool, error) {
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
