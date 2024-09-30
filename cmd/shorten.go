package cmd

import (
	"fmt"

	"github.com/skip2/go-qrcode"

	"github.com/abdealijaroli/jaro/cmd/utils"
	"github.com/abdealijaroli/jaro/store"
)

func ShortenURL(longURL string, store *store.PostgresStore) (string, error) {
	shortURL, shortCode := utils.GenerateShortCode(longURL)

	store.AddShortURLToDB(longURL, shortCode, false)

	qr, err := qrcode.New(shortURL, qrcode.Medium)
	if err != nil {
		return "", err
	}

	fmt.Println("QR Code for your short link: ")
	fmt.Println(qr.ToSmallString(false))

	return shortURL, nil
}
