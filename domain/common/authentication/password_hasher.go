package hash

type IHashService interface {
	Hash(word string) (string, error)
	Match(hashedWord, plainWord string) (bool, error)
}

