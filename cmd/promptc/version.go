package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Muestra la versión",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("PromptC v0.1.0-alpha")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
