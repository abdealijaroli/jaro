package cmd

import (
	// "encoding/json"
	"fmt"
	// "io"
	"log"
	// "net/http"
	"os"
	// "sync"

	"github.com/skip2/go-qrcode"
	"github.com/spf13/cobra"
	// "github.com/abdealijaroli/jaro/store"
)

func transferFile(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	shortCode := GenerateShortCode(filePath)
	shortURL := fmt.Sprintf("https://jaroli.me/%s", shortCode)

	qr, err := qrcode.New(shortURL, qrcode.Medium)
	if err != nil {
		log.Printf("Error generating QR code: %v\n", err)
		return
	}

	fmt.Println("QR Code for your shareable file link: ")
	fmt.Println(qr.ToSmallString(false))

	fmt.Printf("Your shareable file link is: %s\n", shortURL)
}

var transferCmd = &cobra.Command{
	Use:   "transfer",
	Short: "Transfer a file",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Println("Please provide a file to transfer. Run 'jaro --help' for more information.")
			return
		}
		filePath := args[0]

		transferFile(filePath)
	},
}

func init() {
	rootCmd.AddCommand(transferCmd)
}
