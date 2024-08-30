package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var transferCmd = &cobra.Command{
	Use:   "transfer",
	Short: "Transfer a file",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Please provide a file to transfer. Run 'jaro --help' for more information.")
			return
		}
		file := args[0]

		// file transfer logic here
		fmt.Printf("Transferring file: %s\n", file)
	},
}

func init() {
	rootCmd.AddCommand(transferCmd)
}
