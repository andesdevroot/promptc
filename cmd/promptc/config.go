package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/andesdevroot/promptc/internal/cli"
	"github.com/andesdevroot/promptc/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configura de forma interactiva las credenciales",
	Run: func(cmd *cobra.Command, args []string) {
		cli.PrintBanner()
		cli.PrintSection("⚙️  Configuración de PromptC")

		reader := bufio.NewReader(os.Stdin)

		fmt.Println(cli.ColorCyan + "¿Qué proveedor de IA deseas usar?" + cli.ColorReset)
		fmt.Println("1) Google Gemini")
		fmt.Print(cli.ColorYellow + "> " + cli.ColorReset)

		providerOption, _ := reader.ReadString('\n')
		providerOption = strings.TrimSpace(providerOption)

		var provider string
		switch providerOption {
		case "1":
			provider = "gemini"
		default:
			provider = "gemini"
		}

		fmt.Printf("\n🔑 Ingresa tu API Key:\n")
		fmt.Print(cli.ColorYellow + "> " + cli.ColorReset)

		apiKey, _ := reader.ReadString('\n')
		apiKey = strings.TrimSpace(apiKey)

		cfg := config.AppConfig{
			Provider: provider,
			APIKey:   apiKey,
		}

		config.Save(cfg)
		cli.PrintSuccess("¡Configuración guardada en ~/.promptc/config.yaml!")
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
