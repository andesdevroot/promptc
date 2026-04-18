package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/andesdevroot/promptc/pkg/core"
)

const defaultOpenAIModel = "gpt-5.4-mini"

type OpenAIProvider struct {
	APIKey  string
	Model   string
	BaseURL string
	Client  *http.Client
}

func NewOpenAIProvider(apiKey string, model string) *OpenAIProvider {
	if strings.TrimSpace(model) == "" {
		model = defaultOpenAIModel
	}

	return &OpenAIProvider{
		APIKey:  apiKey,
		Model:   model,
		BaseURL: "https://api.openai.com/v1/responses",
		Client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (o *OpenAIProvider) Name() string {
	return fmt.Sprintf("OpenAI (%s)", o.Model)
}

func (o *OpenAIProvider) Optimize(ctx context.Context, p core.Prompt, issues []string) (string, error) {
	if strings.TrimSpace(o.APIKey) == "" {
		return "", fmt.Errorf("OPENAI_API_KEY no configurada")
	}

	systemMsg := `Eres el motor de compilación PROMPTC.
Devuelve exclusivamente el prompt final optimizado.

REGLAS:
1. IDIOMA: responde en español técnico.
2. FORMATO: entrega un prompt final estructurado y accionable.
3. SEGURIDAD: no inventes datos regulatorios ni hechos no presentes en el contexto.
4. ESTILO: prioriza claridad, determinismo y cumplimiento.`

	userMsg := fmt.Sprintf(
		"Optimiza este prompt industrial corrigiendo los siguientes hallazgos: %s\n\nROL: %s\nCONTEXTO: %s\nTAREA: %s\nRESTRICCIONES: %s",
		strings.Join(issues, ", "),
		p.Role,
		p.Context,
		p.Task,
		strings.Join(p.Constraints, " | "),
	)

	payload := map[string]interface{}{
		"model":        o.Model,
		"instructions": systemMsg,
		"input": []map[string]interface{}{
			{
				"role": "user",
				"content": []map[string]string{
					{
						"type": "input_text",
						"text": userMsg,
					},
				},
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, o.BaseURL, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+o.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := o.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var response struct {
		Error *struct {
			Message string `json:"message"`
		} `json:"error,omitempty"`
		Output []struct {
			Type    string `json:"type"`
			Role    string `json:"role"`
			Content []struct {
				Type string `json:"type"`
				Text string `json:"text,omitempty"`
			} `json:"content,omitempty"`
		} `json:"output,omitempty"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}

	if resp.StatusCode >= 400 {
		if response.Error != nil && response.Error.Message != "" {
			return "", fmt.Errorf("openai responses api: %s", response.Error.Message)
		}
		return "", fmt.Errorf("openai responses api devolvió status %d", resp.StatusCode)
	}

	var out strings.Builder
	for _, item := range response.Output {
		if item.Type != "message" || item.Role != "assistant" {
			continue
		}
		for _, content := range item.Content {
			if content.Type == "output_text" && strings.TrimSpace(content.Text) != "" {
				if out.Len() > 0 {
					out.WriteString("\n")
				}
				out.WriteString(content.Text)
			}
		}
	}

	text := strings.TrimSpace(out.String())
	if text == "" {
		return "", fmt.Errorf("openai no devolvió texto utilizable")
	}

	return text, nil
}
