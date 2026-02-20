package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/andesdevroot/promptc/pkg/core"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GeminiProvider struct {
	client *genai.Client
	apiKey string
}

func NewGeminiProvider(ctx context.Context, apiKey string) (*GeminiProvider, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}
	return &GeminiProvider{client: client, apiKey: apiKey}, nil
}

func (g *GeminiProvider) Optimize(ctx context.Context, p core.Prompt, issues []string) (string, error) {
	model := g.client.GenerativeModel("gemini-1.5-flash") // Usamos flash para velocidad

	// El "System Prompt" que garantiza la calidad del SDK
	instructions := fmt.Sprintf(`Actúa como un experto en Ingeniería de Prompts y Arquitectura de Software.
Tu misión es mejorar el siguiente prompt que ha fallado en las pruebas de calidad.

ERRORES DETECTADOS:
%s

PROMPT ORIGINAL (YAML):
Role: %s
Context: %s
Task: %s
Constraints: %s

INSTRUCCIONES DE REPARACIÓN:
1. Resuelve todos los errores detectados.
2. Mantén el idioma original (español).
3. Devuelve únicamente el contenido del prompt MEJORADO en formato texto plano, estructurado con ### ROLE, ### CONTEXT, ### TASK y ### CONSTRAINTS.
4. Asegura que el resultado sea determinista y reduzca alucinaciones.`,
		strings.Join(issues, "\n"), p.Role, p.Context, p.Task, strings.Join(p.Constraints, ", "))

	resp, err := model.GenerateContent(ctx, genai.Text(instructions))
	if err != nil {
		return "", err
	}

	if len(resp.Candidates) == 0 {
		return "", fmt.Errorf("no se recibió respuesta de Gemini")
	}

	// Extraemos el texto de la respuesta
	var result strings.Builder
	for _, part := range resp.Candidates[0].Content.Parts {
		result.WriteString(fmt.Sprintf("%v", part))
	}

	return result.String(), nil
}
