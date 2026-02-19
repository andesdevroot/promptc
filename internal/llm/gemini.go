package llm

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/andesdevroot/promptc/internal/analyzer"
	"github.com/andesdevroot/promptc/internal/models"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// AutoFix toma un prompt deficiente y los errores encontrados, y le pide a Gemini que lo reescriba.
func AutoFix(source models.PromptSource, issues []analyzer.Issue) (string, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("la variable de entorno GEMINI_API_KEY no está configurada")
	}

	ctx := context.Background()

	// Inicializamos el cliente oficial de Google
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return "", fmt.Errorf("error al crear el cliente de Gemini: %v", err)
	}
	defer client.Close()

	// ACTUALIZACIÓN CLAVE: Usamos gemini-2.5-flash (el 1.5 fue deprecado por Google)
	model := client.GenerativeModel("gemini-2.5-flash")

	// Construimos el "Meta-Prompt" (Un prompt para arreglar prompts)
	systemInstruction := `Eres un AI Prompt Engineer Senior.
Tu tarea es tomar un archivo de definición de prompt (YAML) deficiente y corregirlo basándote en los errores de un linter (análisis estático).
Debes devolver ÚNICAMENTE el código YAML corregido, sin bloques de código markdown ("` + "```yaml" + `") ni texto adicional. Solo el YAML puro.
Asegúrate de:
1. Expandir el contexto si es muy corto.
2. Reemplazar palabras ambiguas con instrucciones precisas.
3. Añadir restricciones negativas ("No hacer X") si faltan.`

	// Formateamos los errores para pasárselos a Gemini
	var issuesList strings.Builder
	for _, issue := range issues {
		issuesList.WriteString(fmt.Sprintf("- [%s]: %s\n", issue.Type, issue.Message))
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

	// Configuramos instrucciones del sistema (System Prompt para Gemini)
	model.SystemInstruction = genai.NewUserContent(genai.Text(systemInstruction))

	// Ejecutamos la petición
	resp, err := model.GenerateContent(ctx, genai.Text(userPrompt))
	if err != nil {
		return "", fmt.Errorf("error al generar respuesta de Gemini: %v", err)
	}

	// Extraemos el texto de la respuesta
	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		part := resp.Candidates[0].Content.Parts[0]
		if text, ok := part.(genai.Text); ok {
			// Limpiamos posibles backticks residuales que la IA a veces incluye
			cleanText := strings.TrimPrefix(string(text), "```yaml\n")
			cleanText = strings.TrimPrefix(cleanText, "```\n")
			cleanText = strings.TrimSuffix(cleanText, "\n```")
			return cleanText, nil
		}
	}

	return "", fmt.Errorf("respuesta vacía o formato desconocido de Gemini")
}

// Helper para imprimir las restricciones originales en el meta-prompt
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
