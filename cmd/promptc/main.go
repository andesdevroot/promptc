package main

import (
	"fmt"
	"os"

	"github.com/andesdevroot/promptc/internal/cli"
	"github.com/andesdevroot/promptc/internal/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "promptc",
	Short: "PromptC: The Prompt Compiler for Engineering Excellence",
	Long: `PromptC es una herramienta de ingeniería de software para LLMs.
Valida, analiza y compila prompts deterministas reduciendo alucinaciones.`,
	Run: func(cmd *cobra.Command, args []string) {
		cli.PrintBanner()
		cfg, _ := config.Load()

		if cfg.APIKey == "" {
			fmt.Println(cli.ColorYellow + "👋 ¡Bienvenido a PromptC!" + cli.ColorReset)
			fmt.Println("Parece que aún no has configurado tu motor de IA.")
			fmt.Println("Para comenzar, inicia el asistente:")
			fmt.Println("\n    " + cli.ColorCyan + "promptc config" + cli.ColorReset + "\n")
		} else {
			cmd.Help()
		}
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
