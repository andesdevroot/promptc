package llm

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/andesdevroot/promptc/internal/config"
	"github.com/andesdevroot/promptc/internal/core"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// AutoFix toma un prompt deficiente y le pide a Gemini que lo reescriba.
func AutoFix(source core.PromptSource, issues []string) (string, error) {
	// 1. Intentamos cargar la configuración local
	cfg, err := config.Load()
	if err != nil {
		return "", fmt.Errorf("error leyendo configuración local: %v", err)
	}

	apiKey := cfg.APIKey

	// 2. Fallback: Si no hay llave en la configuración, buscamos en las variables de entorno (Ideal para CI/CD)
	if apiKey == "" {
		apiKey = os.Getenv("GEMINI_API_KEY")
	}

	// 3. Si definitivamente no hay llave, guiamos al usuario a la solución
	if apiKey == "" {
		return "", fmt.Errorf("no se encontró una API Key.\nEjecuta 'promptc config' para configurarla interactivamente")
	}

	ctx := context.Background()

	// Inicializamos el cliente oficial de Google
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return "", fmt.Errorf("error al crear el cliente de Gemini: %v", err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-2.5-flash")

	systemInstruction := `Eres un AI Prompt Engineer Senior.
Tu tarea es tomar un archivo de definición de prompt (YAML) deficiente y corregirlo basándote en los errores de un linter (análisis estático).
Debes devolver ÚNICAMENTE el código YAML corregido, sin bloques de código markdown ("` + "```yaml" + `") ni texto adicional. Solo el YAML puro.
Asegúrate de:
1. Expandir el contexto si es muy corto.
2. Reemplazar palabras ambiguas con instrucciones precisas.
3. Añadir restricciones negativas ("No hacer X") si faltan.`

	var issuesList strings.Builder
	for _, issue := range issues {
		issuesList.WriteString(fmt.Sprintf("- %s\n", issue))
	}

	userPrompt := fmt.Sprintf(`
ERRORES DETECTADOS POR EL LINTER:
%s

PROMPT ORIGINAL (YAML):
role: "%s"
context: "%s"
task: "%s"
constraints:
%s

Por favor, devuelve la versión corregida en formato YAML puro:`,
		issuesList.String(),
		source.Role,
		source.Context,
		source.Task,
		formatConstraintsList(source.Constraints))

	model.SystemInstruction = genai.NewUserContent(genai.Text(systemInstruction))

	resp, err := model.GenerateContent(ctx, genai.Text(userPrompt))
	if err != nil {
		return "", fmt.Errorf("error al generar respuesta de Gemini: %v", err)
	}

	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		part := resp.Candidates[0].Content.Parts[0]
		if text, ok := part.(genai.Text); ok {
			cleanText := strings.TrimPrefix(string(text), "```yaml\n")
			cleanText = strings.TrimPrefix(cleanText, "```\n")
			cleanText = strings.TrimSuffix(cleanText, "\n```")
			return cleanText, nil
		}
	}

	return "", fmt.Errorf("respuesta vacía o formato desconocido de Gemini")
}

func formatConstraintsList(constraints []string) string {
	if len(constraints) == 0 {
		return "  []"
	}
	var sb strings.Builder
	for _, c := range constraints {
		sb.WriteString(fmt.Sprintf("  - \"%s\"\n", c))
	}
	return sb.String()
}
