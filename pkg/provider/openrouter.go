package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/andesdevroot/promptc/pkg/core"
)

type OpenRouterProvider struct {
	APIKey string
	Model  string // "anthropic/claude-3.5-sonnet"
}

func NewOpenRouter(apiKey string) *OpenRouterProvider {
	return &OpenRouterProvider{
		APIKey: apiKey,
		Model:  "anthropic/claude-3.5-sonnet",
	}
}

func (o *OpenRouterProvider) Name() string { return "OpenRouter (Claude 3.5 Sonnet)" }

func (o *OpenRouterProvider) Optimize(ctx context.Context, p core.Prompt, issues []string) (string, error) {
	url := "https://openrouter.ai/api/v1/chat/completions"

	// Instrucción nivel Senior para Claude
	systemMsg := `Eres el motor de compilación PROMPTC. Tu misión es transformar borradores YAML en prompts de sistema deterministas y profesionales.
	REGLAS:
	1. IDIOMA: Español de Chile técnico/minero.
	2. FORMATO: Devuelve SOLO el prompt final, sin introducciones.
	3. CALIDAD: Usa ingeniería de prompts avanzada (Chain of Thought implícito).`

	userMsg := fmt.Sprintf("Optimiza este prompt eliminando: %s\n\nDatos:\nRole: %s\nContext: %s\nTask: %s",
		strings.Join(issues, ", "), p.Role, p.Context, p.Task)

	payload := map[string]interface{}{
		"model": o.Model,
		"messages": []map[string]string{
			{"role": "system", "content": systemMsg},
			{"role": "user", "content": userMsg},
		},
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+o.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var res struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	json.NewDecoder(resp.Body).Decode(&res)

	if len(res.Choices) == 0 {
		return "", fmt.Errorf("OpenRouter no devolvió opciones")
	}

	return strings.TrimSpace(res.Choices[0].Message.Content), nil
}
