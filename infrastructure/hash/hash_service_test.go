package hash

import "testing"

var hashService = GetHashService()

func TestHash(t *testing.T) {
	hashedWord, err := hashService.Hash("plainWord")
	if err != nil {
		t.Errorf("error hashing plainWord: %v", err)
	}

	if hashedWord == "" {
		t.Error("expected hash to be not empty")
	}

	if hashedWord == "plainWord" {
		t.Error("expected hash to be diffrent from plainWord")
	}
}

func TestMatchHash(t *testing.T) {
	hashedWrod, err := hashService.Hash("plainWord")
	if err != nil {
		t.Errorf("error hashing plainWord: %v", err)
	}

	matchResult, err := hashService.Match(hashedWrod, "plainWord")
	if err != nil {
		t.Errorf("error matching plainWord: %v", err)
	}

	if matchResult != true {
		t.Errorf("expected password to match hash")
	}

	matchResult, err = hashService.Match(hashedWrod, "incorrect plainWord")
	if err != nil {
		t.Errorf("error matching incorrect plainWord: %v", err)
	}

	if matchResult == true {
		t.Errorf("expected incorrect plainWord to not match hash")
	}
}
