// internal/parser/parser.go
package parser

import (
	"os"

	"github.com/andesdevroot/promptc/internal/models"
	"gopkg.in/yaml.v3"
)

func ParseFile(filePath string) (models.PromptSource, error) {
	var source models.PromptSource
	data, err := os.ReadFile(filePath)
	if err != nil {
		return source, err
	}
	err = yaml.Unmarshal(data, &source)
	return source, err
}
