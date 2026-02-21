# 🚀 PROMPTC: The Prompt Compiler (v0.1.0-alpha)

```text
    ____                            __  ______
   / __ \_________  ____ ___  ____ / /_/ ____/
  / /_/ / ___/ __ \/ __  __ \/ __ / __/ /     
 / ____/ /  / /_/ / / / / / / /_/ / /_/ /___  
/_/   /_/   \____/_/ /_/ /_/ .___/\__/\____/  
                          /_/                 

   The Prompt Compiler for Engineering Excellence
   v0.1.0-alpha • by Cesar Rivas
```

> **"Prompt Engineering is Software Engineering."**
> PROMPTC es un compilador nativo desarrollado en Go, diseñado para resolver el vacío de soberanía y determinismo en la adopción de IA Generativa para industrias críticas en LATAM.

---

## 🏗️ Visión de Arquitectura: Private-First AI

En sectores regulados como la **Minería (Sernageomin), Banca (CMF) y el Sector Legal**, la lógica de negocio es un activo crítico que no puede ser expuesto a nubes públicas. **PROMPTC** actúa como un **L7 Gateway para LLMs**, permitiendo que herramientas como Claude Desktop o Cursor consuman contextos privados sin que la data sensible abandone la infraestructura corporativa.

### Diferenciadores Core

1.  **Soberanía de Datos**: Orquestación de inferencia local mediante túneles **Tailscale** hacia nodos privados (Mac mini, Ollama, vLLM).
2.  **Abstracción de Vendor**: Compila una vez, despliega en cualquier modelo. Control total sobre el flujo de tokens y el presupuesto de inferencia.
3.  **Compilación Determinista**: Transforma lenguaje ambiguo en estructuras técnicas blindadas, eliminando el "Spanglish" y asegurando el cumplimiento de normativas regionales.

---

## ✨ Características Principales

* **Servidor MCP Nativo**: Implementación completa del *Model Context Protocol* sobre JSON-RPC 2.0 para integración directa con el ecosistema Anthropic y Cursor.
* **Binario Estático en Go**: Cero dependencias en tiempo de ejecución. Rendimiento de alto nivel con consumo mínimo de recursos en workstations y servidores de borde.
* **Prompt-as-Code (PaC)**: Gestión de plantillas mediante componentes versionables y pre-certificados:
    * `PROMPTC_MINERIA_BASE`: Protocolos EPP, seguridad de faena y normativa minera local.
    * `PROMPTC_BANCA_RIESGO`: Alineado con normativas CMF y prevención de fraude (AML).
    * `PROMPTC_LEGAL_CONTRATOS`: Revisión de cláusulas críticas y derecho corporativo.
* **Motor Anti-Spanglish**: Validación semántica estricta que fuerza el uso de terminología técnica precisa en español nativo, eliminando alucinaciones culturales.

---

## 🛠️ Instalación y Configuración

### 1. Prerrequisitos
* Go 1.22+
* Tailscale (Opcional, para modo de inferencia híbrida)
* Ollama o vLLM (Para soberanía total del dato)

### 2. Compilación del Sistema
Para generar un binario de producción optimizado:
```bash
# Limpiar dependencias y compilar
go mod tidy
go build -ldflags="-s -w" -o bin/promptc ./cmd/promptc/main.go
```

### 3. Integración con Claude Desktop / Cursor
Añade el servidor a tu archivo de configuración `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "PROMPTC_CORE": {
      "command": "/Users/cesar_rivas/promptc/bin/promptc",
      "args": ["-mode", "hybrid"],
      "env": {
        "PROMPTC_REMOTE_NODE": "100.x.y.z",
        "PROMPTC_ENV": "production",
        "PROMPTC_LOG_LEVEL": "debug"
      }
    }
  }
}
```

---

## 🎮 Caso de Uso: Orquestación Industrial

Una vez activo, puedes delegar tareas complejas directamente desde tu IDE:

**Usuario**: *"Analiza este reporte de incidente en faena usando PROMPTC_MINERIA_BASE y genera el XML de cumplimiento para el regulador."*

**PROMPTC Workflow**:
1.  **Intercept**: El servidor MCP recibe la solicitud localmente antes de que llegue a la nube pública.
2.  **Route**: Enruta el contexto por el túnel seguro al nodo de inferencia privado (Mac mini).
3.  **Compile**: Inyecta las reglas de la plantilla industrial y valida el output idiomático.
4.  **Deliver**: Devuelve una respuesta determinista y segura a tu editor como un System Prompt.

---

## 🏗️ Estructura del Proyecto (Clean Architecture)

* `cmd/promptc`: Punto de entrada del Servidor JSON-RPC 2.0 (Stdio Bridge).
* `pkg/sdk`: Orquestador de inferencia y lógica de compilación de prompts.
* `pkg/core`: Definiciones de dominio, contratos y esquemas de validación.
* `pkg/provider`: Adaptadores para Ollama (Local) y OpenRouter (Cloud fallback).

---

## 📅 Roadmap v0.2.0

- [x] Implementación Core MCP (Stdio).
- [x] Integración nativa con Ollama local.
- [ ] **PROMPTC Dashboard**: Visualización de observabilidad en tiempo real (Next.js + Go Fiber).
- [ ] **Schema Enforcement**: Validación estricta de estructuras de salida mediante JSON Schema.
- [ ] **Multi-node Load Balancing**: Soporte para clústeres de inferencia distribuida.

---

## 🤝 Contribuciones y Licencia

Este es un proyecto Open Source nacido en Chile para fortalecer el desarrollo de IA soberana en la región. Las contribuciones son bienvenidas vía Pull Requests.

**Licencia**: MIT | **Autor**: Cesar Rivas - Senior Backend Engineer & Cloud Architect.
Desarrollado en La Serena, Chile.