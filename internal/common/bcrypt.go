package common

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func BcryptHash(password string) (string, error) {
	pwHash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}

	return string(pwHash), nil

}

func BcryptCompare(hashed string, plain string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}
