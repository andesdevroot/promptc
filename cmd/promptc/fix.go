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
	Short: "Analiza y auto-optimiza un prompt usando el motor de IA de PromptC",
	Long:  `Analiza la estructura de un prompt y utiliza Gemini Pro para corregir deficiencias semánticas y de determinismo.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cli.PrintBanner()

		// 1. Cargar configuración de usuario
		cfg, err := config.Load()
		if err != nil || cfg.APIKey == "" {
			cli.PrintError("Error: API Key no configurada. Ejecuta 'promptc config' primero.")
			os.Exit(1)
		}

		// 2. Parsear el archivo YAML usando el motor del SDK (pkg/core.Prompt)
		p, err := parser.ParseFile(args[0])
		if err != nil {
			cli.PrintError(fmt.Sprintf("Error al leer el prompt: %v", err))
			os.Exit(1)
		}

		// 3. Inicializar el SDK
		// Usamos context.Background() para la gestión de la conexión con la API de Google
		ctx := context.Background()
		promptcSDK, err := sdk.NewSDK(ctx, cfg.APIKey)
		if err != nil {
			cli.PrintError(fmt.Sprintf("Error al inicializar el SDK: %v", err))
			os.Exit(1)
		}

		// 4. Ejecutar el análisis técnico (Reglas del SDK)
		cli.PrintSection("📋 Análisis de Calidad del SDK")
		analysis := promptcSDK.Analyze(p)

		// Mostrar el Score con color según su valor
		renderScore(analysis.Score)

		// 5. Lógica de Optimización si el Score es insuficiente
		if !analysis.IsReliable {
			cli.PrintWarning("⚠️  Calidad insuficiente para producción. Iniciando optimización...")

			// Llamada al motor de IA del SDK para reparar el prompt
			optimized, err := promptcSDK.Optimize(ctx, p)
			if err != nil {
				cli.PrintError(fmt.Sprintf("Error durante la optimización: %v", err))
				os.Exit(1)
			}

			cli.PrintSuccess("✨ Prompt Optimizado por PromptC:")
			fmt.Println("\n" + optimized)
		} else {
			cli.PrintSuccess("✅ El prompt cumple con los estándares de ingeniería de PromptC.")
		}
	},
}

// renderScore ayuda a visualizar la calidad en la terminal
func renderScore(score int) {
	color := cli.ColorGreen
	if score < 80 {
		color = cli.ColorRed
	} else if score < 95 {
		color = cli.ColorYellow
	}
	fmt.Printf("Score de Ingeniería: %s%d/100%s\n\n", color, score, cli.ColorReset)
}

func init() {
	rootCmd.AddCommand(fixCmd)
}
