package engine

import (
	"testing"
	"time"

	"github.com/andesdevroot/promptc/pkg/core"
)

func TestCompilerEngine_Compile(t *testing.T) {
	engine := New()

	tests := []struct {
		name     string
		prompt   core.Prompt
		expected string
	}{
		{
			name: "Compilación Base",
			prompt: core.Prompt{
				Role:    "Analista Senior",
				Context: "Prevención de Fraudes",
				Task:    "Eres un {{role}} experto en {{context}}. Evalúa este caso.",
			},
			expected: "### ROLE\nAnalista Senior\n\n### CONTEXT\nPrevención de Fraudes\n\n### TASK\nEres un Analista Senior experto en Prevención de Fraudes. Evalúa este caso.",
		},
		{
			name: "Escenario Completo Industrial",
			prompt: core.Prompt{
				Role:    "Auditor CMF",
				Context: "Riesgo Operacional",
				Task:    "Actúa como {{role}} en {{context}}. Revisa el proceso de {{area}}.",
				Variables: map[string]string{
					"area": "Tesorería",
				},
				Constraints: []string{"Citar norma 461"},
				CreatedAt:   time.Now(),
			},
			expected: "### ROLE\nAuditor CMF\n\n### CONTEXT\nRiesgo Operacional\n\n### TASK\nActúa como Auditor CMF en Riesgo Operacional. Revisa el proceso de Tesorería.\n\n### CONSTRAINTS\n- Citar norma 461\n\n### VARIABLES\n- area: Tesorería",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := engine.Compile(tt.prompt)
			if err != nil {
				t.Fatalf("Compile() error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("\n=== FALLO EN COMPILACIÓN: %s ===\nEXPECTED:\n%q\nGOT:\n%q", tt.name, tt.expected, result)
			}
		})
	}
}

func TestCompilerEngine_Analyze(t *testing.T) {
	eng := New()

	tests := []struct {
		name           string
		prompt         core.Prompt
		minScore       int
		shouldReliable bool
	}{
		{
			name: "Prompt Débil (Sin estructura)",
			prompt: core.Prompt{
				Task: "Haz un resumen.",
			},
			minScore:       20,
			shouldReliable: false,
		},
		{
			name: "Prompt Industrial (Completo)",
			prompt: core.Prompt{
				Role:    "Cloud Architect",
				Context: "Migración de microservicios a AWS",
				Task:    "Diseña un diagrama de arquitectura para una aplicación bancaria que debe soportar 10k TPS con latencia menor a 50ms.",
				Constraints: []string{
					"Usar AWS EKS",
					"Incluir esquema de seguridad IAM",
				},
			},
			minScore:       90,
			shouldReliable: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := eng.Analyze(tt.prompt)
			if result.Score < tt.minScore {
				t.Errorf("%s: Score demasiado bajo. Esperado min %d, obtenido %d", tt.name, tt.minScore, result.Score)
			}
			if result.IsReliable != tt.shouldReliable {
				t.Errorf("%s: Confiabilidad incorrecta. Esperada %v, obtenida %v", tt.name, tt.shouldReliable, result.IsReliable)
			}
		})
	}
}
