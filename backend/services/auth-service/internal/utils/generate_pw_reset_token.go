package utils

import (
	"crypto/rand"
	"encoding/base64"
)

func GeneratePasswordResetToken() (string, error) {
	const tokenLength = 32 
	b := make([]byte, tokenLength)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	token := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(b)
	return token, nil
}


