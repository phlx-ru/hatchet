package secrets

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func MakeHash(secret string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func MustMakeHash(secret string) string {
	hash, err := MakeHash(secret)
	if err != nil {
		panic(err)
	}
	return hash
}

func CompareSourceAndHash(source, hash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(source))
	mismatched := errors.Is(err, bcrypt.ErrMismatchedHashAndPassword)
	if err != nil && !mismatched {
		return false, err
	}
	return !mismatched, nil
}

func MustCompareSourceAndHash(source, hash string) bool {
	matched, err := CompareSourceAndHash(source, hash)
	if err != nil {
		panic(err)
	}
	return matched
}
