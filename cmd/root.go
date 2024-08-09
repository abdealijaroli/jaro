package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var shortenURL string
var transferFile string

var rootCmd = &cobra.Command{
	Use:   "jaro",
	Short: "Jaro CLI",
	Long:  `Jaro CLI for shortening links and transferring files.`,
	Run: func(cmd *cobra.Command, args []string) {
		if shortenURL != "" {
			shortened := ShortenURL(shortenURL)
			fmt.Printf("Shortened URL: %s\n", shortened)
		} else if transferFile != "" {
			fmt.Printf("Transferring file: %s\n", transferFile)
		} else {
			fmt.Println("Please provide a valid option")
		}
	},
}

func Execute() {
	rootCmd.Flags().StringVarP(&shortenURL, "shorten", "s", "", "URL to shorten")
	rootCmd.Flags().StringVarP(&transferFile, "transfer", "t", "", "File to transfer")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}