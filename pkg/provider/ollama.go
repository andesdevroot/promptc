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

type OllamaProvider struct {
	BaseURL string
	Model   string
	Client  *http.Client
}

func NewOllamaProvider(ip string) *OllamaProvider {
	return &OllamaProvider{
		BaseURL: fmt.Sprintf("http://%s:11434/api/generate", ip),
		Model:   "llama3",
		Client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (o *OllamaProvider) Name() string { return "Ollama Remote Node (Mac mini)" }

func (o *OllamaProvider) Optimize(ctx context.Context, p core.Prompt, issues []string) (string, error) {
	// Definimos el comportamiento esperado con un ejemplo claro (Few-Shot)
	// Esto obliga al modelo a seguir el patrón de idioma y formato.
	instruction := fmt.Sprintf(`Eres un Compilador de Prompts Técnico. Tu salida debe ser exclusivamente el prompt final optimizado.

### REGLAS DE ORO:
1. IDIOMA: Escribe TODO en ESPAÑOL DE CHILE/TÉCNICO.
2. PROHIBICIÓN: No hables con el usuario. No digas "Aquí está tu prompt". No uses inglés.
3. FORMATO: Devuelve un System Prompt estructurado.

### EJEMPLO DE COMPILACIÓN:
INPUT: {Role: "Dev", Context: "Web", Task: "Fix bug"}
OUTPUT: Actúa como un Desarrollador Senior. Tu contexto es un entorno web moderno. Tu tarea es identificar y corregir errores de lógica de forma eficiente.

### TAREA REAL A COMPILAR:
ROL: %s
CONTEXTO: %s
TAREA: %s
ERRORES A CORREGIR: %s

OUTPUT OPTIMIZADO EN ESPAÑOL:`, 
		p.Role, p.Context, p.Task, strings.Join(issues, ", "))

	payload := map[string]interface{}{
		"model":  o.Model,
		"prompt": instruction,
		"stream": false,
		"options": map[string]interface{}{
			"temperature": 0.3, // Bajamos la temperatura para que sea más determinista y menos "creativo" (evita alucinaciones de idioma)
		},
	}

	jsonData, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, "POST", o.BaseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := o.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error en enlace Tailscale: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Response string `json:"response"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	// Limpieza final por si el modelo ignora las instrucciones de no hablar
	finalPrompt := strings.TrimSpace(result.Response)
	finalPrompt = strings.TrimPrefix(finalPrompt, "Aquí está el prompt optimizado:")
	finalPrompt = strings.TrimPrefix(finalPrompt, "Optimized Prompt:")

	return finalPrompt, nil
}