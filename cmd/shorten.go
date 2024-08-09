package cmd

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)

func ShortenURL(url string) string {
	h := sha1.New()
	h.Write([]byte(url))
	shortURL := hex.EncodeToString(h.Sum(nil))[:5]
	return fmt.Sprintf("https://jaro.li/%s", shortURL)
}
