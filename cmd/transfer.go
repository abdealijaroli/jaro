package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/skip2/go-qrcode"

	"github.com/abdealijaroli/jaro/cmd/stream"
	"github.com/abdealijaroli/jaro/store"
)

func TransferFile(filePath string, store *store.PostgresStore) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	shortURL := stream.InitiateTransfer(filePath, store)

	qr, err := qrcode.New(shortURL, qrcode.Medium)
	if err != nil {
		log.Printf("Error generating QR code: %v\n", err)
		return
	}

	fmt.Println("QR Code for your shareable file link: ")
	fmt.Println(qr.ToSmallString(false))

	fmt.Printf("Your shareable file link is: %s\n", shortURL)
}
