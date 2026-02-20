package parser

import (
	"fmt"
	"os"

	"github.com/andesdevroot/promptc/pkg/core"
	"gopkg.in/yaml.v3"
)

// LoadPrompt lee un archivo YAML del disco y lo convierte en la estructura core.Prompt
func LoadPrompt(filename string) (core.Prompt, error) {
	var p core.Prompt

	// 1. Leer el archivo físico
	data, err := os.ReadFile(filename)
	if err != nil {
		return p, fmt.Errorf("no se pudo leer el archivo %s: %w", filename, err)
	}

	// 2. Deserializar (Unmarshal) el YAML a la estructura de Go
	if err := yaml.Unmarshal(data, &p); err != nil {
		return p, fmt.Errorf("error de sintaxis en el YAML: %w", err)
	}

	return p, nil
}
