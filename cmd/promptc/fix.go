package main

import (
	"fmt"

	"github.com/andesdevroot/promptc/internal/analyzer"
	"github.com/andesdevroot/promptc/internal/cli"
	"github.com/andesdevroot/promptc/internal/llm"
	"github.com/andesdevroot/promptc/internal/parser"
	"github.com/spf13/cobra"
)

var fixCmd = &cobra.Command{
	Use:   "fix [archivo.yaml]",
	Short: "Auto-corrige un prompt deficiente",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cli.PrintBanner()
		source, _ := parser.ParseFile(args[0])
		score, issues := analyzer.Analyze(source)

		cli.PrintWarning(fmt.Sprintf("Score: %d/100. Corrigiendo...", score))

		fixed, err := llm.AutoFix(source, issues)
		if err != nil {
			cli.PrintError(fmt.Sprintf("Error: %v", err))
			return
		}

		cli.PrintSuccess("¡Prompt corregido!")
		fmt.Println(fixed)
	},
}

func init() {
	rootCmd.AddCommand(fixCmd)
}
