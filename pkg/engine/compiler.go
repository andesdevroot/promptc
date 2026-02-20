package engine

import (
	"fmt"
	"strings"

	"github.com/andesdevroot/promptc/pkg/core"
)

type CompilerEngine struct {
	MinScoreThreshold int
}

func New() *CompilerEngine {
	return &CompilerEngine{
		MinScoreThreshold: 85, // Subimos el estándar para el lanzamiento
	}
}

func (e *CompilerEngine) Compile(p core.Prompt) (string, error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("### ROLE\n%s\n\n", strings.TrimSpace(p.Role)))
	sb.WriteString(fmt.Sprintf("### CONTEXT\n%s\n\n", strings.TrimSpace(p.Context)))
	sb.WriteString(fmt.Sprintf("### TASK\n%s\n\n", strings.TrimSpace(p.Task)))

	if len(p.Constraints) > 0 {
		sb.WriteString("### CONSTRAINTS\n")
		for _, c := range p.Constraints {
			if strings.TrimSpace(c) != "" {
				sb.WriteString(fmt.Sprintf("- %s\n", strings.TrimSpace(c)))
			}
		}
	}
	return sb.String(), nil
}

func (e *CompilerEngine) Analyze(p core.Prompt) core.Result {
	score := 100
	var issues []string
	var suggestions []string

	// 1. Rigor del Rol (Lógica Senior)
	cleanRole := strings.TrimSpace(p.Role)
	if len(cleanRole) < 25 {
		score -= 25
		issues = append(issues, "Identidad de agente débil.")
		suggestions = append(suggestions, "En español, usa roles con contexto de industria (ej: 'Ingeniero de Minas experto en Seguridad' en lugar de 'Experto').")
	}

	// 2. Anti-Hallucination: El 'Fuego Amigo' del Subjuntivo
	// En español, usar el modo subjuntivo ("quisiera que hicieras") confunde el determinismo.
	subjunctivePatterns := []string{"quisiera", "me gustaría", "tal vez", "podrías"}
	cleanTask := strings.ToLower(p.Task)
	for _, pattern := range subjunctivePatterns {
		if strings.Contains(cleanTask, pattern) {
			score -= 15
			issues = append(issues, "Uso de lenguaje condicional o ambiguo.")
			suggestions = append(suggestions, "Cambia el condicional por imperativos directos: 'Analiza', 'Genera', 'Calcula'.")
			break
		}
	}

	// 3. Verificación de Restricciones Críticas
	hasNegative := false
	keywords := []string{"no ", "evita", "nunca", "prohibido", "sin inventar", "excluye"}
	for _, c := range p.Constraints {
		lowC := strings.ToLower(c)
		for _, kw := range keywords {
			if strings.Contains(lowC, kw) {
				hasNegative = true
				break
			}
		}
	}

	if !hasNegative {
		score -= 40
		issues = append(issues, "Ausencia de Negative Constraints.")
		suggestions = append(suggestions, "Añade: 'No utilices información fuera del contexto proporcionado' para blindar el prompt.")
	}

	return core.Result{
		Score:       e.clamp(score),
		IsReliable:  score >= e.MinScoreThreshold,
		Issues:      issues,
		Suggestions: suggestions,
	}
}

func (e *CompilerEngine) clamp(score int) int {
	if score < 0 {
		return 0
	}
	return score
}
