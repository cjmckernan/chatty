package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func GenerateSessionId() (string, error) {

	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("failed to generate session ID: %w", err)
	}
	sessionID := base64.URLEncoding.EncodeToString(b)
	return sessionID, nil
}

func HashPassword(password string) (string, error) {
	handlePassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to has password %w", err)
	}

	return string(handlePassword), nil
}
