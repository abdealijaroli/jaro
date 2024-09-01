package cmd

import (
	"fmt"
	"github.com/skip2/go-qrcode"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"os"
)

func transferFile(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	shortCode := generateShortCode(filePath)
	shortURL := fmt.Sprintf("https://jaroli.me/%s", shortCode)

	qr, err := qrcode.New(shortURL, qrcode.Medium)
	if err != nil {
		fmt.Printf("Error generating QR code: %v\n", err)
		return
	}

	fmt.Printf("Your shareable file link is: %s\n", shortURL)
	fmt.Println("QR Code for your link:")
	fmt.Println(qr.ToSmallString(false))

	go startFileServer(shortCode, file)
}

func startFileServer(shortCode string, file *os.File) {
	http.HandleFunc("/"+shortCode, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "<h1>File Transfer</h1>")
		fmt.Fprintf(w, "<p>Sender: %s</p>", "SenderName")
		fmt.Fprintf(w, "<p>File: %s</p>", file.Name())
		fmt.Fprintf(w, "<form method='POST' action='/accept/%s'><button type='submit'>Accept</button></form>", shortCode)
	})

	http.HandleFunc("/accept/"+shortCode, func(w http.ResponseWriter, r *http.Request) {
		transferToRecipient(w, file)
	})

	http.ListenAndServe(":8080", nil)
}

func transferToRecipient(w http.ResponseWriter, file *os.File) {
	w.Header().Set("Content-Disposition", "attachment; filename="+file.Name())
	w.Header().Set("Content-Type", "application/octet-stream")

	_, err := io.Copy(w, file)
	if err != nil {
		http.Error(w, "Error transferring file", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "File delivered successfully.")
}

func generateShortCode(filePath string) string {
	return ""
}

var transferCmd = &cobra.Command{
	Use:   "transfer",
	Short: "Transfer a file",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Please provide a file to transfer. Run 'jaro --help' for more information.")
			return
		}
		filePath := args[0]

		transferFile(filePath)
	},
}

func init() {
	rootCmd.AddCommand(transferCmd)
}
