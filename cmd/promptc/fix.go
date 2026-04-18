package main

import (
	"context"
	"fmt"
	"os"
	"strings"

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
		cfg, err := config.Load()
		if err != nil {
			fmt.Println("Error cargando configuración:", err)
			os.Exit(1)
		}

		p, err := parser.ParseFile(args[0])
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		ctx := context.Background()
		openAIKey := strings.TrimSpace(cfg.OpenAIAPIKey)
		geminiKey := strings.TrimSpace(cfg.GeminiAPIKey)

		switch cfg.Provider {
		case "openai":
			if openAIKey == "" {
				openAIKey = strings.TrimSpace(cfg.APIKey)
			}
		case "gemini":
			if geminiKey == "" {
				geminiKey = strings.TrimSpace(cfg.APIKey)
			}
		}

		if envKey := strings.TrimSpace(os.Getenv("OPENAI_API_KEY")); envKey != "" {
			openAIKey = envKey
		}
		if envKey := strings.TrimSpace(os.Getenv("GEMINI_API_KEY")); envKey != "" {
			geminiKey = envKey
		}

		openAIModel := strings.TrimSpace(cfg.OpenAIModel)
		if envModel := strings.TrimSpace(os.Getenv("OPENAI_MODEL")); envModel != "" {
			openAIModel = envModel
		}

		promptcSDK, err := sdk.NewSDK(ctx, openAIKey, openAIModel, geminiKey, os.Getenv("PROMPTC_MACMINI_IP"))
		if err != nil {
			fmt.Println("Error inicializando SDK:", err)
			os.Exit(1)
		}

		analysis := promptcSDK.Analyze(p)
		fmt.Printf("Score: %d/100\n", analysis.Score)

		if !analysis.IsReliable {
			optimized, err := promptcSDK.Optimize(ctx, p)
			if err != nil {
				fmt.Printf("\n❌ Error Crítico: %v\n", err)
				os.Exit(1)
			}
			cli.PrintSuccess("\n✨ Prompt Optimizado:")
			fmt.Println("\n" + optimized)
		} else {
			output, _ := promptcSDK.Engine.Compile(p)
			fmt.Println("\n" + output)
		}
	},
}

func init() {
	rootCmd.AddCommand(fixCmd)
}
