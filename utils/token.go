package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
)

// GenerateSecureToken returns a random token string and its base64-encoded SHA256 hash
func GenerateSecureToken() (string, string, error) {
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", "", err
	}

	rawToken := base64.URLEncoding.EncodeToString(tokenBytes)
	encodedToken := GetEncodedHashedToken(rawToken)

	return rawToken, encodedToken, nil
}

// VerifySecureToken checks if the tokenHash is the SHA256 hash of rawToken
func VerifySecureToken(tokenHash string, rawToken string) bool {
	encodedToken := GetEncodedHashedToken(rawToken)
	return encodedToken == tokenHash
}

// GetEncodedHashedToken returns a base64-encoded SHA256 hash
func GetEncodedHashedToken(rawToken string) string {
	hash := sha256.Sum256([]byte(rawToken))
	b64Hash := base64.URLEncoding.EncodeToString(hash[:])

	return b64Hash
}
