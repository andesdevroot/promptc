package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Muestra la versión de PROMPTC",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("PROMPTC v0.3.1 (Codex-Ready Industrial MCP)")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
