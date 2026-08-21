package config

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(plainText string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainText), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func IsPasswordMatch(input, digest string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(digest), []byte(input))
	return err == nil
}

func GenerateRandomShit() (string, error) {
	randomByte := make([]byte, 64)
	if _, err := rand.Read(randomByte); err != nil {
		return "", err
	}

	hash := sha256.Sum256(randomByte)

	return hex.EncodeToString(hash[:]), nil
}