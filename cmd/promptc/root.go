package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "promptc",
	Short: "PROMPTC CLI para compilacion de prompts industriales",
	Long:  "PROMPTC expone un servidor MCP y utilidades CLI para compilacion, analisis y operacion auditable de prompts industriales.",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	},
}

func execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
