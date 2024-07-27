package hash

import "testing"

func setup() *Service {
	return SingletonService()
}

func TestHashService(t *testing.T) {
	hashService := setup()

	// Testing Hash function
	t.Run("HashWord", func(t *testing.T) {
		hashedWord, err := hashService.Hash("plainWord")
		if err != nil {
			t.Fatalf("error hashing plainWord: %v", err)
		}

		if hashedWord == "" {
			t.Error("expected hash to be not empty")
		}

		if hashedWord == "plainWord" {
			t.Error("expected hash to be different from plainWord")
		}
	})

	// Testing Match with correct plain text
	t.Run("MatchHashedWithCorrectPlain", func(t *testing.T) {
		hashedWord, err := hashService.Hash("plainWord")
		if err != nil {
			t.Fatalf("error hashing plainWord: %v", err)
		}

		matchResult, err := hashService.Match(hashedWord, "plainWord")
		if err != nil {
			t.Fatalf("error matching plainWord: %v", err)
		}

		if !matchResult {
			t.Errorf("expected plainWord to match hash")
		}
	})

	// Testing Match with incorrect plain text
	t.Run("MatchHashedWithIncorrectPlain", func(t *testing.T) {
		hashedWord, err := hashService.Hash("plainWord")
		if err != nil {
			t.Fatalf("error hashing plainWord: %v", err)
		}

		matchResult, err := hashService.Match(hashedWord, "incorrectPlainWord")
		if err != nil {
			t.Fatalf("error matching incorrectPlainWord: %v", err)
		}

		if matchResult {
			t.Errorf("expected incorrectPlainWord to not match hash")
		}
	})

	// Testing Match with invalid hashedWord
	t.Run("MatchWithInvalidHashedWord", func(t *testing.T) {
		_, err := hashService.Match("invalidHashedWord", "plainWord")
		if err == nil {
			t.Error("expected an error with invalid hashed word")
		}
	})
}
