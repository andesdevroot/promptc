package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/andesdevroot/promptc/pkg/core"
	"github.com/andesdevroot/promptc/pkg/sdk"
)

// Estructuras nativas para JSON-RPC 2.0 (Estándar MCP)
type JSONRPCMessage struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id,omitempty"`
	Method  string          `json:"method,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type JSONRPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// Catálogo de Plantillas convertido a Librería Interna (Soberanía de Datos)
var templateLibrary = map[string]string{
	"PROMPTC_MINERIA_BASE": `ROLE: Ingeniero Experto en Seguridad y Salud Ocupacional (Sernageomin).
CONTEXTO: Operación de faena minera en Chile. Altas exigencias de seguridad, uso de EPP, control de fatiga y somnolencia.
TASK: Evaluar protocolos de prevención de riesgos y emitir directrices de mitigación.`,

	"PROMPTC_BANCA_RIESGO": `ROLE: Analista Senior de Riesgo Crediticio y Cumplimiento CMF.
CONTEXTO: Institución financiera chilena bancaria. Evaluación de carteras comerciales y prevención de lavado de activos (UAF).
TASK: Estructurar reglas de decisión para evaluación de créditos corporativos.`,

	"PROMPTC_LEGAL_CONTRATOS": `ROLE: Abogado Corporativo Senior (Corporate Law).
CONTEXTO: Legislación chilena. Redacción y revisión de contratos comerciales, NDAs y acuerdos de confidencialidad.
TASK: Detectar vacíos legales, cláusulas abusivas y asegurar protección de propiedad intelectual.`,
}

func main() {
	// Todos los logs de auditoría van a Stderr obligatoriamente
	fmt.Fprintf(os.Stderr, "[PROMPTC] Iniciando PROMPTC Engine (MCP Server)...\n")

	geminiKey := os.Getenv("GEMINI_API_KEY")
	macMiniIP := os.Getenv("PROMPTC_MACMINI_IP")
	if macMiniIP == "" {
		macMiniIP = "100.90.6.101"
	}

	ctx := context.Background()
	app, err := sdk.NewSDK(ctx, geminiKey, macMiniIP)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[PROMPTC] Error crítico inicializando SDK: %v\n", err)
		os.Exit(1)
	}

	// Scanner para leer las peticiones de Claude Desktop/Cursor vía Stdio
	scanner := bufio.NewScanner(os.Stdin)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	for scanner.Scan() {
		line := scanner.Bytes()
		var req JSONRPCMessage

		if err := json.Unmarshal(line, &req); err != nil {
			fmt.Fprintf(os.Stderr, "[PROMPTC] Ignorando entrada no-JSON.\n")
			continue
		}

		// Enrutador del protocolo MCP
		switch req.Method {
		case "initialize":
			sendResponse(req.ID, map[string]interface{}{
				"protocolVersion": "2024-11-05",
				// BRANDING CLAVE: Esto le dice a Claude cómo llamarnos en la UI
				"serverInfo": map[string]string{"name": "PROMPTC-Engine", "version": "1.3.0"},
				"capabilities": map[string]interface{}{
					"tools": map[string]interface{}{},
				},
			})

		case "notifications/initialized":
			fmt.Fprintf(os.Stderr, "[PROMPTC] Handshake completado. Herramientas activas en el IDE.\n")
			continue

		case "tools/list":
			sendResponse(req.ID, map[string]interface{}{
				"tools": []map[string]interface{}{
					{
						"name":        "get_template",
						"description": "Obtiene una plantilla base de PROMPTC para una industria específica. Úsala SIEMPRE que el usuario mencione una plantilla.",
						"inputSchema": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"template_name": map[string]string{
									"type":        "string",
									"description": "Nombre de la plantilla. Opciones válidas: PROMPTC_MINERIA_BASE, PROMPTC_BANCA_RIESGO, PROMPTC_LEGAL_CONTRATOS",
								},
							},
							"required": []string{"template_name"},
						},
					},
					{
						"name":        "optimize_prompt",
						"description": "Compila y optimiza un borrador de prompt usando el motor PROMPTC. Elimina Spanglish y fuerza terminología técnica en español.",
						"inputSchema": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"role":    map[string]string{"type": "string", "description": "Rol del agente. DEBE estar en español."},
								"context": map[string]string{"type": "string", "description": "Contexto del proyecto. Especificar industria y país."},
								"task":    map[string]string{"type": "string", "description": "Tarea específica a ejecutar por el agente."},
							},
							"required": []string{"role", "context", "task"},
						},
					},
				},
			})

		case "tools/call":
			handleToolCall(req, app)
		}
	}
}

func handleToolCall(req JSONRPCMessage, app *sdk.PromptC) {
	// Parseo dinámico para soportar múltiples herramientas
	var call struct {
		Name      string          `json:"name"`
		Arguments json.RawMessage `json:"arguments"`
	}

	if err := json.Unmarshal(req.Params, &call); err != nil {
		sendError(req.ID, -32602, "Invalid params parsing tool call")
		return
	}

	// ENRUTADOR DE HERRAMIENTAS
	switch call.Name {

	case "get_template":
		var args struct {
			TemplateName string `json:"template_name"`
		}
		json.Unmarshal(call.Arguments, &args)
		fmt.Fprintf(os.Stderr, "[PROMPTC] Ejecutando get_template: %s\n", args.TemplateName)

		if content, exists := templateLibrary[args.TemplateName]; exists {
			sendResponse(req.ID, map[string]interface{}{
				"content": []map[string]interface{}{
					{"type": "text", "text": fmt.Sprintf("Contenido de la plantilla %s:\n\n%s", args.TemplateName, content)},
				},
			})
		} else {
			sendError(req.ID, -32602, fmt.Sprintf("Plantilla no encontrada: %s", args.TemplateName))
		}

	case "optimize_prompt":
		var args struct {
			Role    string `json:"role"`
			Context string `json:"context"`
			Task    string `json:"task"`
		}
		json.Unmarshal(call.Arguments, &args)
		fmt.Fprintf(os.Stderr, "[PROMPTC] Ejecutando compilación en nodo remoto...\n")

		ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
		defer cancel()

		p := core.Prompt{
			Role:    args.Role,
			Context: args.Context,
			Task:    args.Task,
		}

		result, err := app.CompileAndOptimize(ctx, p)
		if err != nil {
			sendError(req.ID, -32000, fmt.Sprintf("Error de compilación: %v", err))
			return
		}

		// BRANDING DE LANZAMIENTO ESTRICTO
		finalOutput := fmt.Sprintf("### ⚡ PROMPTC - COMPILACIÓN EXITOSA ⚡\n\n%s\n\n---\n*Compilado de forma determinista y segura por PROMPTC Engine.*", result)

		sendResponse(req.ID, map[string]interface{}{
			"content": []map[string]interface{}{
				{"type": "text", "text": finalOutput},
			},
		})

	default:
		sendError(req.ID, -32601, "Tool not found")
	}
}

// Helpers para JSON-RPC (Estos son los ÚNICOS lugares donde se usa fmt.Println hacia Stdout)
func sendResponse(id interface{}, result interface{}) {
	resp := JSONRPCResponse{JSONRPC: "2.0", ID: id, Result: result}
	out, _ := json.Marshal(resp)
	fmt.Println(string(out))
}

func sendError(id interface{}, code int, message string) {
	resp := JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error:   map[string]interface{}{"code": code, "message": message},
	}
	out, _ := json.Marshal(resp)
	fmt.Println(string(out))
}
