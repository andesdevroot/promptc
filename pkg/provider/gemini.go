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
}

func NewGeminiProvider(ctx context.Context, apiKey string) (*GeminiProvider, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}
	return &GeminiProvider{client: client}, nil
}

func (g *GeminiProvider) Name() string {
	return "Gemini 1.5 Pro Cloud"
}

func (g *GeminiProvider) Optimize(ctx context.Context, p core.Prompt, issues []string) (string, error) {
	model := g.client.GenerativeModel("gemini-1.5-pro")

	// Construcción del Meta-Prompt de Optimización
	var metaPrompt strings.Builder
	metaPrompt.WriteString("Eres un experto en Ingeniería de Prompts Industriales.\n")
	metaPrompt.WriteString("Tu misión es reconstruir un prompt débil basándote en los ISSUES detectados.\n\n")

	metaPrompt.WriteString("ISSUES A CORREGIR:\n")
	for _, issue := range issues {
		metaPrompt.WriteString("- " + issue + "\n")
	}

	metaPrompt.WriteString("\nESTRUCTURA OBLIGATORIA DEL OUTPUT:\n")
	metaPrompt.WriteString("### ROLE\n[Definición técnica del experto]\n\n")
	metaPrompt.WriteString("### CONTEXT\n[Escenario detallado de la industria]\n\n")
	metaPrompt.WriteString("### TASK\n[La instrucción principal optimizada]\n\n")
	metaPrompt.WriteString("### CONSTRAINTS\n- [Restricción 1]\n- [Restricción 2]\n\n")

	metaPrompt.WriteString("PROMPT ORIGINAL DEL USUARIO:\n")
	metaPrompt.WriteString(fmt.Sprintf("Rol: %s\nContexto: %s\nTarea: %s\n", p.Role, p.Context, p.Task))

	resp, err := model.GenerateContent(ctx, genai.Text(metaPrompt.String()))
	if err != nil {
		return "", err
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("Gemini devolvió una respuesta vacía")
	}

	// Extraemos la primera parte de la respuesta
	part := resp.Candidates[0].Content.Parts[0]
	return fmt.Sprintf("%v", part), nil
}
