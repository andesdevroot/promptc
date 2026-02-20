package engine

import (
	"fmt"
	"strings"

	"github.com/andesdevroot/promptc/pkg/core"
)

// CompilerEngine es el motor que procesa y valida los prompts.
type CompilerEngine struct {
	MinScoreThreshold int
}

// New genera una nueva instancia del motor con valores por defecto.
func New() *CompilerEngine {
	return &CompilerEngine{
		MinScoreThreshold: 80, // Un prompt de menos de 80 es "peligroso"
	}
}

// Compile transforma el objeto Prompt en un string optimizado para LLMs.
func (e *CompilerEngine) Compile(p core.Prompt) (string, error) {
	var sb strings.Builder

	// Aplicamos estructura de ingeniería
	sb.WriteString(fmt.Sprintf("### ROLE: %s\n\n", p.Role))
	sb.WriteString(fmt.Sprintf("### CONTEXT: %s\n\n", p.Context))
	sb.WriteString(fmt.Sprintf("### TASK: %s\n\n", p.Task))

	if len(p.Constraints) > 0 {
		sb.WriteString("### CONSTRAINTS:\n")
		for _, c := range p.Constraints {
			sb.WriteString(fmt.Sprintf("- %s\n", c))
		}
	}

	return sb.String(), nil
}

// Analyze verifica la calidad del prompt y detecta riesgos de alucinación.
func (e *CompilerEngine) Analyze(p core.Prompt) core.Result {
	score := 100
	var issues []string
	var suggestions []string

	// 1. Validación de Restricciones Negativas (Anti-alucinación)
	hasNegative := false
	for _, c := range p.Constraints {
		cLow := strings.ToLower(c)
		// Buscamos patrones de control en español
		if strings.Contains(cLow, "no ") || strings.Contains(cLow, "evita") || strings.Contains(cLow, "prohibido") {
			hasNegative = true
		}
	}

	if !hasNegative {
		score -= 30
		issues = append(issues, "Faltan restricciones negativas.")
		suggestions = append(suggestions, "Agrega: 'No inventes información si no estás seguro de la respuesta'.")
	}

	// 2. Validación de especificidad del Rol
	if len(p.Role) < 15 {
		score -= 20
		issues = append(issues, "Rol demasiado genérico.")
		suggestions = append(suggestions, "Define el expertise técnico, ej: 'Experto en Seguridad Cloud en AWS'.")
	}

	return core.Result{
		Score:       score,
		IsReliable:  score >= e.MinScoreThreshold,
		Issues:      issues,
		Suggestions: suggestions,
	}
}
