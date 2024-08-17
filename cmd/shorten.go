package cmd

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/skip2/go-qrcode"

	"github.com/abdealijaroli/jaro/store"
)

func ShortenURL(url string, storage *store.PostgresStore) (string, error) {
	hash := sha256.Sum256([]byte(url))
	shortCode := base64.URLEncoding.EncodeToString(hash[:])[:6]
	shortURL := fmt.Sprintf("https://jaroli.me/%s", shortCode)

	err := storage.CreateShortURL(url, shortCode)
	if err != nil {
		return "", err
	}

	qr, err := qrcode.New(shortURL, qrcode.Medium)
	if err != nil {
		return "", err
	}

	fmt.Println("QR Code:")
	fmt.Println(qr.ToSmallString(false))

	return shortURL, nil
}
