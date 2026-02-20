package analyzer

import (
	"strings"

	"github.com/andesdevroot/promptc/internal/core"
)

// Heurísticas de Anti-Alucinación y Calidad
func Analyze(p core.PromptSource) (int, []string) {
	score := 100
	var issues []string

	// 1. Validación de Rol (Evita alucinaciones de identidad)
	if len(p.Role) < 10 {
		score -= 20
		issues = append(issues, "Rol demasiado vago. Define expertise (ej: 'Senior Cloud Architect' vs 'Arquitecto').")
	}

	// 2. Validación de Restricciones (Crucial para Anti-Hallucination)
	hasNegativeConstraint := false
	for _, c := range p.Constraints {
		lowC := strings.ToLower(c)
		if strings.Contains(lowC, "no") || strings.Contains(lowC, "evita") || strings.Contains(lowC, "sin") {
			hasNegativeConstraint = true
		}
	}

	if !hasNegativeConstraint {
		score -= 30
		issues = append(issues, "Faltan restricciones negativas. Sin límites, el LLM es propenso a alucinar.")
	}

	// 3. Diferenciador: Análisis de Tono en Español
	if !detectSpanishClarity(p.Task) {
		score -= 15
		issues = append(issues, "La tarea principal carece de verbos de acción fuertes en español (ej: 'Analiza', 'Genera', 'Valida').")
	}

	return score, issues
}

func detectSpanishClarity(task string) bool {
	// Verbos de alto impacto para LLMs en español
	strongVerbs := []string{"Analiza", "Genera", "Escribe", "Valida", "Calcula", "Resume", "Traduce"}
	for _, v := range strongVerbs {
		if strings.Contains(strings.ToLower(task), strings.ToLower(v)) {
			return true
		}
	}
	return false
}
