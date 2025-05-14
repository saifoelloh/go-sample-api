package utils

import (
	"crypto/sha512"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

func CryptoHash(inputString string) string {
	hasher := sha512.New()
	hasher.Write([]byte(inputString))
	return hex.EncodeToString(hasher.Sum(nil))
}

func HashBcrypt(password string) (string, error) {
	// Generate salt and hash in one step (bcrypt.GenerateFromPassword does both)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		// Wrap error into InternalServerError
		return "", err
	}

	return string(hashedPassword), nil
}
