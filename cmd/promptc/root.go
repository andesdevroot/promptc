package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "promptc",
	Short: "PromptC es un compilador y optimizador de prompts",
	Long: `PromptC es una herramienta CLI para gestionar, compilar y optimizar 
prompts de IA estructurados en archivos YAML.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
