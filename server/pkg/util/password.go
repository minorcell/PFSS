package util

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const (
	// Cost determines the computational cost of hashing passwords
	defaultBcryptCost = 12
	// Length of the salt in bytes
	saltLength = 16
)

// GenerateSalt generates a random salt
func GenerateSalt() (string, error) {
	salt := make([]byte, saltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}
	return base64.StdEncoding.EncodeToString(salt), nil
}

// HashPassword combines password with salt and creates a hash
func HashPassword(password string) (string, error) {
	// Generate hash
	hash, err := bcrypt.GenerateFromPassword([]byte(password), defaultBcryptCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hash), nil
}

// ValidatePassword checks if the provided password matches the hash
func ValidatePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
