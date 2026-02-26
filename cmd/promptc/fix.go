package main

import (
	"context"
	"fmt"
	"os"

	"github.com/andesdevroot/promptc/internal/cli"
	"github.com/andesdevroot/promptc/internal/config"
	"github.com/andesdevroot/promptc/internal/parser"
	"github.com/andesdevroot/promptc/pkg/sdk"
	"github.com/spf13/cobra"
)

var fixCmd = &cobra.Command{
	Use:   "fix [archivo.yaml]",
	Short: "Analiza y repara un prompt con redundancia de IA",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cli.PrintBanner()
		cfg, _ := config.Load()
		p, err := parser.ParseFile(args[0])
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		ctx := context.Background()
		promptcSDK, _ := sdk.NewSDK(ctx, cfg.APIKey, "")

		analysis := promptcSDK.Engine.Analyze(p)
		
		score := analysis.Score
		fmt.Printf("Score: %d/100\n", score)

		isReliable := analysis.IsReliable
		if !isReliable {
			optimized, err := promptcSDK.CompileAndOptimize(ctx, p)
			if err != nil {
				fmt.Printf("\n❌ Error Crítico: %v\n", err)
				os.Exit(1)
			}
			cli.PrintSuccess("\n✨ Prompt Optimizado:")
			fmt.Println("\n" + fmt.Sprint(optimized))
		} else {
			output, _ := promptcSDK.Engine.Compile(p)
			fmt.Println("\n" + output)
		}
	},
}

func init() {
	rootCmd.AddCommand(fixCmd)
}
