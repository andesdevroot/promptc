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

// OllamaProvider implementa la interfaz Optimizer para nodos locales.
type OllamaProvider struct {
	BaseURL string
	Model   string
	Client  *http.Client
}

// OllamaRequest es la estructura para el endpoint /api/generate
type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

// OllamaResponse es la respuesta simplificada del servidor
type OllamaResponse struct {
	Response string `json:"response"`
}

// NewOllamaProvider inicializa el nodo de inferencia privada.
// Se espera que remoteIP sea la IP de Tailscale de tu Mac Mini.
func NewOllamaProvider(remoteIP string) *OllamaProvider {
	return &OllamaProvider{
		BaseURL: fmt.Sprintf("http://%s:11434", remoteIP),
		Model:   "llama3", // O el modelo que tengas corriendo (mistral, phi3, etc)
		Client: &http.Client{
			Timeout: 10 * time.Second, // Timeout corto para asegurar fallback rápido
		},
	}
}

func (o *OllamaProvider) Name() string {
	return "Ollama Enterprise Node (Private)"
}

func (o *OllamaProvider) Optimize(ctx context.Context, p core.Prompt, issues []string) (string, error) {
	// Construimos el Meta-Prompt idéntico al de Gemini para consistencia industrial
	var sb strings.Builder
	sb.WriteString("### SYSTEM INSTRUCTION\n")
	sb.WriteString("Eres un experto en Ingeniería de Prompts Industriales. Tu misión es corregir los ISSUES del prompt.\n\n")

	sb.WriteString("ISSUES DETECTADOS:\n")
	for _, issue := range issues {
		sb.WriteString("- " + issue + "\n")
	}

	sb.WriteString("\nFORMATO DE SALIDA:\n")
	sb.WriteString("Solo entrega el prompt final estructurado en Markdown (ROLE, CONTEXT, TASK, CONSTRAINTS).\n\n")

	sb.WriteString("ENTRADA CRUDA:\n")
	sb.WriteString(fmt.Sprintf("Rol: %s | Contexto: %s | Tarea: %s\n", p.Role, p.Context, p.Task))

	// Preparamos el payload JSON
	reqBody := OllamaRequest{
		Model:  o.Model,
		Prompt: sb.String(),
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling ollama request: %w", err)
	}

	// Ejecutamos la llamada HTTP al nodo local
	req, err := http.NewRequestWithContext(ctx, "POST", o.BaseURL+"/api/generate", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := o.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("nodo Ollama inaccesible: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ollama retorno status: %d", resp.StatusCode)
	}

	var ollamaResp OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return "", fmt.Errorf("error decoding ollama response: %w", err)
	}

	return ollamaResp.Response, nil
}
