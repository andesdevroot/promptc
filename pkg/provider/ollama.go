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

// OllamaProvider implementa la interfaz Optimizer para modelos locales.
type OllamaProvider struct {
	BaseURL string
	Model   string
	Client  *http.Client
}

// NewOllamaProvider inicializa el proveedor con la IP de la Mac mini.
func NewOllamaProvider(ip string) *OllamaProvider {
	return &OllamaProvider{
		BaseURL: fmt.Sprintf("http://%s:11434/api/generate", ip),
		Model:   "llama3",
		Client: &http.Client{
			Timeout: 60 * time.Second, // Timeout extendido para inferencia local
		},
	}
}

func (o *OllamaProvider) Name() string {
	return "Ollama Remote Node (Mac mini)"
}

// Optimize envía el prompt al modelo local para su mejora técnica.
func (o *OllamaProvider) Optimize(ctx context.Context, p core.Prompt, issues []string) (string, error) {
	// Construcción del meta-prompt para la optimización
	instruction := fmt.Sprintf(`### TAREA
Eres un Senior Prompt Engineer. Optimiza el siguiente prompt corrigiendo: %s.

### PROMPT ORIGINAL
ROL: %s
CONTEXTO: %s
TAREA: %s

### REGLA
Responde exclusivamente con el nuevo prompt optimizado. Sin explicaciones.`,
		strings.Join(issues, ", "), p.Role, p.Context, p.Task)

	payload := map[string]interface{}{
		"model":  o.Model,
		"prompt": instruction,
		"stream": false,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("error marshalling payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", o.BaseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := o.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("nodo Mac mini no alcanzable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error del servidor Ollama: status %d", resp.StatusCode)
	}

	var result struct {
		Response string `json:"response"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("error decoding ollama response: %w", err)
	}

	return strings.TrimSpace(result.Response), nil
}
