package analyzer

import (
	"strings"

	"github.com/andesdevroot/promptc/internal/models"
)

// Issue representa un problema encontrado en el prompt
type Issue struct {
	Type    string // "CRITICAL", "WARNING", "TIP"
	Message string
}

// Analyze realiza el análisis estático del prompt
func Analyze(p models.PromptSource) (int, []Issue) {
	score := 100
	var issues []Issue

	// 1. REGLA: Contexto Insuficiente (Riesgo Alto)
	if len(p.Context) < 20 {
		score -= 40
		issues = append(issues, Issue{
			Type:    "CRITICAL",
			Message: "El contexto es peligrosamente corto. El modelo tenderá a alucinar.",
		})
	}

	// 2. REGLA: Ambigüedad en la Tarea
	vagueWords := []string{"rápido", "breve", "mejor posible", "algo así", "creo"}
	taskLower := strings.ToLower(p.Task)
	for _, word := range vagueWords {
		if strings.Contains(taskLower, word) {
			score -= 10
			issues = append(issues, Issue{
				Type:    "WARNING",
				Message: "Palabra ambigua detectada: '" + word + "'. Sé más específico.",
			})
		}
	}

	// 3. REGLA: Falta de Constraints Negativos (Security)
	hasNegative := false
	for _, c := range p.Constraints {
		cLower := strings.ToLower(c)
		if strings.Contains(cLower, "no ") || strings.Contains(cLower, "evita") || strings.Contains(cLower, "nunca") {
			hasNegative = true
			break
		}
	}
	if !hasNegative {
		score -= 15
		issues = append(issues, Issue{
			Type:    "SECURITY",
			Message: "Faltan 'Negative Constraints'. Debes decir explícitamente qué NO hacer.",
		})
	}

	// Normalización del score
	if score < 0 {
		score = 0
	}

	return score, issues
}
