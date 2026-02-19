package parser

import (
	"os"

	"github.com/andesdevroot/promptc/internal/models"
	"gopkg.in/yaml.v3"
)

// ParseFile lee un archivo YAML y lo convierte en un PromptSource
func ParseFile(filePath string) (models.PromptSource, error) {
	var source models.PromptSource

	data, err := os.ReadFile(filePath)
	if err != nil {
		return source, err
	}

	err = yaml.Unmarshal(data, &source)
	if err != nil {
		return source, err
	}

	return source, nil
}
