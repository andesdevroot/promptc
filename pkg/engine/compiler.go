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
		MinScoreThreshold: 85,
	}
}

// ResolveVariables sustituye todos los {{placeholders}} en el contenido
// usando primero los campos core del Prompt, luego el mapa Variables.
// Si un placeholder no tiene valor, lo marca como [MISSING:key] para
// que el operador lo vea en el dashboard y en el audit.log.
func (e *CompilerEngine) ResolveVariables(content string, p core.Prompt) string {
	// 1. Campos core del struct — siempre disponibles
	coreVars := map[string]string{
		"role":    strings.TrimSpace(p.Role),
		"context": strings.TrimSpace(p.Context),
		"task":    strings.TrimSpace(p.Task),
	}

	// 2. Constraints como bloque de texto para el placeholder {{constraints}}
	if len(p.Constraints) > 0 {
		parts := make([]string, 0, len(p.Constraints))
		for _, c := range p.Constraints {
			if strings.TrimSpace(c) != "" {
				parts = append(parts, "- "+strings.TrimSpace(c))
			}
		}
		coreVars["constraints"] = strings.Join(parts, "\n")
	} else {
		coreVars["constraints"] = ""
	}

	// 3. Merge con Variables map — permite override y variables de negocio custom
	merged := make(map[string]string)
	for k, v := range coreVars {
		merged[k] = v
	}
	for k, v := range p.Variables {
		merged[k] = v
	}

	// 4. Resolución de placeholders
	result := content
	for key, value := range merged {
		placeholder := "{{" + key + "}}"
		result = strings.ReplaceAll(result, placeholder, value)
	}

	// 5. Detectar y marcar placeholders no resueltos
	// Esto aparece en el audit log como [MISSING:nombre_variable]
	// permitiendo al operador identificar qué variables faltan
	for {
		start := strings.Index(result, "{{")
		if start == -1 {
			break
		}
		end := strings.Index(result[start:], "}}")
		if end == -1 {
			break
		}
		placeholder := result[start : start+end+2]
		key := result[start+2 : start+end]
		result = strings.ReplaceAll(result, placeholder, "[MISSING:"+strings.TrimSpace(key)+"]")
	}

	return result
}

// Compile construye el prompt final resolviendo variables antes de
// armar la estructura por secciones. El Task puede contener un template
// completo con placeholders — ResolveVariables lo resuelve primero.
func (e *CompilerEngine) Compile(p core.Prompt) (string, error) {
	var sb strings.Builder

	// Resolver el Task si contiene placeholders de template
	resolvedTask := e.ResolveVariables(p.Task, p)

	sb.WriteString(fmt.Sprintf("### ROLE\n%s\n\n", strings.TrimSpace(p.Role)))
	sb.WriteString(fmt.Sprintf("### CONTEXT\n%s\n\n", strings.TrimSpace(p.Context)))
	sb.WriteString(fmt.Sprintf("### TASK\n%s\n\n", resolvedTask))

	if len(p.Constraints) > 0 {
		sb.WriteString("### CONSTRAINTS\n")
		for _, c := range p.Constraints {
			if strings.TrimSpace(c) != "" {
				sb.WriteString(fmt.Sprintf("- %s\n", strings.TrimSpace(c)))
			}
		}
		sb.WriteString("\n")
	}

	// Variables de negocio adicionales expuestas al modelo como sección propia
	// Esto permite al nodo de inferencia ver el contexto completo de variables
	if len(p.Variables) > 0 {
		sb.WriteString("### VARIABLES\n")
		for k, v := range p.Variables {
			sb.WriteString(fmt.Sprintf("- %s: %s\n", k, v))
		}
		sb.WriteString("\n")
	}

	return sb.String(), nil
}

func (e *CompilerEngine) Analyze(p core.Prompt) core.Result {
	score := 100
	var issues []string
	var suggestions []string

	// 1. Rigor del Rol
	cleanRole := strings.TrimSpace(p.Role)
	if len(cleanRole) < 25 {
		score -= 25
		issues = append(issues, "Identidad de agente débil.")
		suggestions = append(suggestions, "En español, usa roles con contexto de industria (ej: 'Ingeniero de Minas experto en Seguridad' en lugar de 'Experto').")
	}

	// 2. Anti-Hallucination: subjuntivo
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

	// 3. Negative Constraints — obligatorio para industrias críticas
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

	// 4. Penalizar placeholders sin resolver en el Task
	// Si llegan {{variables}} sin resolver al Analyze, significa que
	// el operador no proveyó el mapa Variables antes de compilar
	if strings.Contains(p.Task, "{{") && strings.Contains(p.Task, "}}") {
		score -= 20
		issues = append(issues, "El Task contiene placeholders sin resolver.")
		suggestions = append(suggestions, "Provee el mapa Variables con los valores correspondientes antes de compilar.")
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
