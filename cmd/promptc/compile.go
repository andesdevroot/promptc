package main

import (
	"fmt"
	"os"

	"github.com/andesdevroot/promptc/internal/parser"
	"github.com/andesdevroot/promptc/pkg/engine"
	"github.com/spf13/cobra"
	// Ya no necesitamos internal/core aquí
)

var compileCmd = &cobra.Command{
	Use:   "compile [archivo.yaml]",
	Short: "Compila un prompt YAML a formato plano optimizado",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Ahora p es automáticamente de tipo pkg/core.Prompt
		p, err := parser.ParseFile(args[0])
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		compiler := engine.New()
		output, err := compiler.Compile(p)
		if err != nil {
			fmt.Printf("Error de compilación: %v\n", err)
			os.Exit(1)
		}

		fmt.Print(output)
	},
}

func init() {
	rootCmd.AddCommand(compileCmd)
}
