package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/andesdevroot/promptc/pkg/core"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type GeminiProvider struct {
	client      *genai.Client
	activeModel string
}

func NewGeminiProvider(ctx context.Context, apiKey string) (*GeminiProvider, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("error al inicializar cliente genai: %w", err)
	}

	provider := &GeminiProvider{client: client}
	err = provider.discoverBestModel(ctx)
	if err != nil {
		return nil, err
	}

	return provider, nil
}

func (g *GeminiProvider) discoverBestModel(ctx context.Context) error {
	iter := g.client.ListModels(ctx)

	// Prioridad: 1.5 Flash (estable), 1.5 Pro (estable), luego el resto.
	// PRAGMATISMO: Evitamos modelos '-exp' que reportan quota limit 0.
	available := make(map[string]bool)

	for {
		m, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}

		name := strings.TrimPrefix(m.Name, "models/")

		// Filtro: Solo modelos que soporten generación y NO sean experimentales
		supportsGenerate := false
		for _, op := range m.SupportedGenerationMethods {
			if op == "generateContent" {
				supportsGenerate = true
				break
			}
		}

		if supportsGenerate && !strings.Contains(name, "-exp") {
			available[name] = true
		}
	}

	// Selección por ranking de estabilidad
	ranking := []string{"gemini-1.5-flash", "gemini-1.5-pro", "gemini-1.0-pro"}
	for _, r := range ranking {
		if available[r] {
			g.activeModel = r
			return nil
		}
	}

	// Fallback a cualquier modelo no experimental
	for name := range available {
		g.activeModel = name
		return nil
	}

	return fmt.Errorf("no se encontraron modelos estables con cuota disponible")
}

func (g *GeminiProvider) Name() string {
	if g.activeModel != "" {
		return fmt.Sprintf("Google Gemini (%s)", g.activeModel)
	}
	return "Google Gemini"
}

func (g *GeminiProvider) Optimize(ctx context.Context, p core.Prompt, issues []string) (string, error) {
	model := g.client.GenerativeModel(g.activeModel)
	model.SetTemperature(0.2)

	var sb strings.Builder
	sb.WriteString("Optimiza este prompt profesionalmente en Español:\n")
	sb.WriteString(fmt.Sprintf("ROLE: %s\nTASK: %s\n", p.Role, p.Task))
	sb.WriteString("Usa headers ### ROLE, ### CONTEXT, ### TASK, ### CONSTRAINTS.")

	resp, err := model.GenerateContent(ctx, genai.Text(sb.String()))
	if err != nil {
		// Retornamos el error original para que el SDK decida qué hacer
		return "", err
	}

	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil {
		return "", fmt.Errorf("respuesta vacía de la IA")
	}

	var res strings.Builder
	for _, part := range resp.Candidates[0].Content.Parts {
		res.WriteString(fmt.Sprint(part))
	}
	return res.String(), nil
}
