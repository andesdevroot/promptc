package analyzer

import (
	"strings"

	"github.com/andesdevroot/promptc/internal/models"
)

// Listado de palabras que causan alucinación (Vaguedad)
var vagueWords = []string{
	"hazlo lo mejor posible",
	"hazlo lo mejor que puedas",
	"hazlo lo mejor posible",
	"hazlo lo mejor que puedas",
	"rápido",
	"brevemente",
	"algo así",
	"importante", // Subjetivo: ¿que es importante para el usuario?
	"creo que",
	"espero que",
	"ojalá",
	"ojalá que",
}

// Analyze revisa el prompt y devuelve warnings y un score calculado.
func Analyze(p models.PromptSource) ([]string, int) {
	var warnings []string
	score := 100

	// REGLA 1: Longitud de Contexto
	// Un contexto vacío o muy corto es garantía de invención.
	if len(p.Context) < 20 {
		warnings = append(warnings, "CRITICAL: El contexto es demasiado corto. El modelo inventará datos.")
		score -= 40
	}

	// REGLA 2: Detección de Ambigüedad
	// Buscamos palabras prohibidas en la Tarea.
	taskLower := strings.ToLower(p.Task)
	for _, word := range vagueWords {
		if strings.Contains(taskLower, word) {
			warnings = append(warnings, "AMBIGUITY: Evita la palabra '"+word+"'. Sé específico con métricas o formatos.")
			score -= 10
		}
	}

	// REGLA 3: Restricciones Negativas
	// Es vital decir qué NO hacer.
	hasNegative := false
	for _, c := range p.Constraints {
		cLower := strings.ToLower(c)
		if strings.Contains(cLower, "no ") || strings.Contains(cLower, "evita") || strings.Contains(cLower, "nunca") {
			hasNegative = true
			break
		}
	}
	if !hasNegative {
		warnings = append(warnings, "SECURITY: Faltan restricciones negativas (Negative Constraints).")
		score -= 15
	}

	// Normalizar score (que no baje de 0)
	if score < 0 {
		score = 0
	}

	return warnings, score
}
