package utils

import "golang.org/x/crypto/bcrypt"

const hashCost = bcrypt.DefaultCost

// Generates bcrypt hash string for the given password.
func HashPassword(password string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), hashCost)
	return hash, err
}

// Verifies if the given password matches the stored hash.
func VerifyPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
