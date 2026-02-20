package core

import "time"

// Prompt es la unidad fundamental del SDK.
type Prompt struct {
	ID          string            `json:"id"`
	Version     string            `json:"version"`
	Role        string            `json:"role"`
	Context     string            `json:"context"`
	Task        string            `json:"task"`
	Constraints []string          `json:"constraints"`
	Variables   map[string]string `json:"variables"`
	CreatedAt   time.Time         `json:"created_at"`
}

// Result es lo que devuelve el compilador tras analizar un prompt.
type Result struct {
	Score       int      `json:"score"`       // 0-100
	IsReliable  bool     `json:"is_reliable"` // ¿Pasa los filtros anti-alucinación?
	Issues      []string `json:"issues"`      // Problemas encontrados
	Suggestions []string `json:"suggestions"` // Cómo mejorarlo (especialmente en español)
}
