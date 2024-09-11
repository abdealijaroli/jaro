package utils

import (
	"crypto/sha256"
	"encoding/base64"
)

func GenerateShortCode(filePath string) string {
	hash := sha256.Sum256([]byte(filePath))
	return base64.URLEncoding.EncodeToString(hash[:])[:6]
}
