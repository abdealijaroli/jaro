package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/abdealijaroli/jaro/store"
)

var shortenURL string
var transferFile string

var rootCmd = &cobra.Command{
	Use:   "jaro",
	Short: "Jaro CLI",
	Long:  `Jaro CLI for shortening links and transferring files.`,
	Run: func(cmd *cobra.Command, args []string) {
		storage, err := store.NewPostgresStore()
		if err != nil {
			fmt.Printf("Error connecting to database: %v", err)
		}
		defer storage.Close()

		if shortenURL != "" {
			shortened, err := ShortenURL(shortenURL, storage)
			if err != nil {
				fmt.Println("Error shortening link: ", err)
				return
			}
			fmt.Printf("Your sweetened link: %s\n", shortened)
		} else if transferFile != "" {
			fmt.Printf("Transferring file: %s\n", transferFile)
		} else {
			fmt.Println("Please provide a valid option")
		}
	},
}

func Execute() {
	rootCmd.Flags().StringVarP(&shortenURL, "shorten", "s", "", "Link to shorten")
	rootCmd.Flags().StringVarP(&transferFile, "transfer", "t", "", "File to transfer")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
