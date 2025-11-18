package encrypt

import (
	"crypto/sha256"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

func HashWithSHA256String(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func HashWithSHA256Bytes(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

func CompareWithSHA256String(data string, hash string) bool {
	return sha256.Sum256([]byte(data)) == sha256.Sum256([]byte(hash))
}

func CompareWithSHA256Bytes(data []byte, hash []byte) bool {
	return sha256.Sum256(data) == sha256.Sum256(hash)
}

func HashWithBcryptString(data string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(data), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hash[:]), nil
}

func HashWithBcryptBytes(data []byte) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword(data, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func CompareWithBcryptString(data string, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(data)) == nil
}

func CompareWithBcryptBytes(data []byte, hash []byte) bool {
	return bcrypt.CompareHashAndPassword(hash, data) == nil
}
