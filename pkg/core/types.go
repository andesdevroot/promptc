package core

import (
	"time"
)

// Prompt representa la unidad fundamental de datos para el SDK.
type Prompt struct {
	ID          string            `yaml:"id" json:"id"`
	Version     string            `yaml:"version" json:"version"`
	Role        string            `yaml:"role" json:"role"`
	Context     string            `yaml:"context" json:"context"`
	Task        string            `yaml:"task" json:"task"`
	Constraints []string          `yaml:"constraints" json:"constraints"`
	Variables   map[string]string `yaml:"variables" json:"variables"`
	CreatedAt   time.Time         `yaml:"created_at" json:"created_at"`
}

// Result contiene el veredicto técnico del análisis de un prompt.
type Result struct {
	Score       int      `json:"score"`
	IsReliable  bool     `json:"is_reliable"`
	Issues      []string `json:"issues"`
	Suggestions []string `json:"suggestions"`
}
