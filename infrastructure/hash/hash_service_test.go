package hash

import "testing"

func TestHashService(t *testing.T) {
	hashService := GetHashService()

	t.Run("HashWord", func(t *testing.T) {
		hashedWord, err := hashService.Hash("plainWord")
		if err != nil {
			t.Errorf("error hashing plainWord: %v", err)
		}

		if hashedWord == "" {
			t.Error("expected hash to be not empty")
		}

		if hashedWord == "plainWord" {
			t.Error("expected hash to be different from plainWord")
		}
	})

	t.Run("MatchHashedWithCorrectPlain", func(t *testing.T) {
		hashedWord, err := hashService.Hash("plainWord")
		if err != nil {
			t.Errorf("error hashing plainWord: %v", err)
		}

		matchResult, err := hashService.Match(hashedWord, "plainWord")
		if err != nil {
			t.Errorf("error matching plainWord: %v", err)
		}

		if !matchResult {
			t.Errorf("expected password to match hash")
		}
	})

	t.Run("MatchHashedWithIncorrectPlain", func(t *testing.T) {
		hashedWord, err := hashService.Hash("plainWord")
		if err != nil {
			t.Errorf("error hashing plainWord: %v", err)
		}

		matchResult, err := hashService.Match(hashedWord, "incorrectPlainWord")
		if err != nil {
			t.Errorf("error matching incorrectPlainWord: %v", err)
		}

		if matchResult {
			t.Errorf("expected incorrectPlainWord to not match hash")
		}
	})
}

