package parser

import (
	"os"

	"gopkg.in/yaml.v3"

	// Importamos el core del SDK
	"github.com/andesdevroot/promptc/pkg/core"
)

// ParseFile ahora devuelve el tipo oficial del SDK: core.Prompt
func ParseFile(path string) (core.Prompt, error) {
	var p core.Prompt

	data, err := os.ReadFile(path)
	if err != nil {
		return p, err
	}

	err = yaml.Unmarshal(data, &p)
	if err != nil {
		return p, err
	}

	return p, nil
}
