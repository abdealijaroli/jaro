package cmd

import (
	"log"

	"github.com/abdealijaroli/jaro/cmd/utils"
	"github.com/abdealijaroli/jaro/store"
)

func TransferFile(filePath string, store *store.PostgresStore) {
	shortURL, shortCode := utils.GenerateShortCode(filePath)
	utils.GenerateQRCode(shortURL)

	// add shortURL to database
	err := store.AddShortURLToDB(filePath, shortCode, true)
	if err != nil {
		log.Printf("Error adding shortURL to database: %v\n", err)
	}
}
