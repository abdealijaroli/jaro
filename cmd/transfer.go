package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/skip2/go-qrcode"
	"github.com/spf13/cobra"

	"github.com/abdealijaroli/jaro/store"
)

func InitiateTransfer(w http.ResponseWriter, r *http.Request, storage *store.PostgresStore) {
	// Parse the request body
	var req struct {
		FilePath string `json:"filePath"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Handle WebSocket connections
}

var serverPort = ":8080"
var wg sync.WaitGroup

func transferFile(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	shortCode := generateShortCode(filePath)
	shortURL := fmt.Sprintf("https://jaroli.me/%s", shortCode)

	qr, err := qrcode.New(shortURL, qrcode.Medium)
	if err != nil {
		log.Printf("Error generating QR code: %v\n", err)
		return
	}

	fmt.Printf("Your shareable file link is: %s\n", shortURL)
	fmt.Println("QR Code for your link:")
	fmt.Println(qr.ToSmallString(false))

	wg.Add(1)
	go startFileServer(shortCode, file)
}

func startFileServer(shortCode string, file *os.File) {
	defer wg.Done()

	http.HandleFunc("/"+shortCode, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, file.Name())
	})

	http.HandleFunc("/accept/"+shortCode, func(w http.ResponseWriter, r *http.Request) {
		transferToRecipient(w, file)
	})

	log.Printf("Starting file server on port %s\n", serverPort)
	if err := http.ListenAndServe(serverPort, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func transferToRecipient(w http.ResponseWriter, file *os.File) {
	w.Header().Set("Content-Disposition", "attachment; filename="+file.Name())
	w.Header().Set("Content-Type", "application/octet-stream")

	if _, err := io.Copy(w, file); err != nil {
		http.Error(w, "Error transferring file", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "File delivered successfully.")
}

func generateShortCode(filePath string) string {
	// Simple short code generation logic (e.g., using the file name)
	return filePath[len(filePath)-6:] // Example: last 6 characters of the file name
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
