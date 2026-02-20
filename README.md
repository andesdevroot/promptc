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

> **Ingeniería de Prompts Determinista para la Era de la Confiabilidad.**

`PROMPTC` es una herramienta de sistema profesional, escrita en **Go**, diseñada para transformar la intención humana vaga en instrucciones blindadas para Modelos de Lenguaje Extensos (LLMs). 

A diferencia de los "templates" tradicionales, `PROMPTC` aplica principios de **compiladores** (análisis, optimización semántica y enrutamiento híbrido) para reducir alucinaciones, asegurar la estructura técnica y estandarizar la calidad del output en entornos de alta criticidad.

---

## 💡 El Problema: El Abismo del Compliance en IA

En industrias reguladas como la **Minería, Banca y Sector Legal**, el uso de LLMs comerciales presenta un riesgo inaceptable: la fuga de propiedad intelectual y datos sensibles hacia nubes públicas. Los ingenieros necesitan el poder de la IA en sus IDEs, pero las normativas (CMF, Sernageomin, GDPR) exigen soberanía sobre los datos.

## 🛡️ La Solución: Arquitectura Soberana PROMPTC

`PROMPTC` no es solo un optimizador; es un **Servidor MCP (Model Context Protocol)** que actúa como un puente seguro entre tu IDE y tu propia infraestructura:

1.  **Inferencia Híbrida:** Enruta las solicitudes de optimización a través de un túnel **Tailscale** hacia nodos de computación privados (ej. un Mac mini local con Llama 3).
2.  **Orquestación Autónoma:** Expone herramientas inteligentes que los LLMs (Claude, Cursor) pueden encadenar para construir soluciones complejas.
3.  **Determinismo Regional:** Fuerza el uso de terminología técnica en español chileno/latino, eliminando el "Spanglish" y las alucinaciones culturales de los modelos base.

---

## ✨ Características Principales

* **Tool Chaining:** Permite al LLM buscar plantillas industriales y optimizarlas en un solo flujo de pensamiento.
* **Librería de Componentes (Resources):**
    * `PROMPTC_MINERIA_BASE`: Foco en seguridad de faena, protocolos EPP y normativa Sernageomin.
    * `PROMPTC_BANCA_RIESGO`: Alineado con normativas CMF, prevención de fraudes y lavado de activos.
    * `PROMPTC_LEGAL_CONTRATOS`: Estructuras de derecho corporativo y revisión de cláusulas críticas.
* **Seguridad por Diseño:** Comunicación vía `stdio` (entrada/salida estándar), garantizando que el servidor MCP solo responda a procesos autorizados localmente.

---

## 🛠️ Instalación y Configuración

### 1. Compilación del Binario
Requiere Go 1.21 o superior.
```bash
go mod tidy
go build -o promptc ./cmd/promptc/main.go
```

### 2. Integración con Claude Desktop / Cursor
Añade el servidor a tu configuración de MCP (`claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "PROMPTC": {
      "command": "/Users/TU_USUARIO/Desktop/GO/promptc/promptc",
      "env": {
        "PROMPTC_TU_MAQUINA": "TU_IP"
      }
    }
  }
}
```

---

## 🎮 Caso de Uso: Orquestación en Acción

Una vez configurado, puedes interactuar con **PROMPTC** de forma natural en tu chat de IA:

**Usuario:** *"Usa la plantilla PROMPTC_BANCA_RIESGO y compílala para un agente que analice fraudes en transferencias Swift."*

**PROMPTC Engine:**
1.  Llamada a `get_template("PROMPTC_BANCA_RIESGO")` -> Extrae reglas de cumplimiento local.
2.  Llamada a `optimize_prompt(...)` -> Cruza los datos hacia el nodo privado (Mac mini).
3.  **Resultado:** Un prompt de sistema blindado, listo para producción.

---

## 🏗️ Estructura del Proyecto

* `cmd/promptc`: Servidor JSON-RPC 2.0 nativo (Stdio Bridge).
* `pkg/sdk`: Orquestador de inferencia y lógica de compilación.
* `pkg/core`: Contratos y modelos de dominio para la ingeniería de prompts.
* `pkg/provider`: Implementaciones para Ollama (Local) y OpenRouter (Cloud fallback).

---

## 📝 Licencia
Este proyecto es Open Source bajo la licencia **MIT**.

---
**PROMPTC: Elevando la ingeniería de prompts al estándar de la ingeniería de software.** Desarrollado con ❤️ en Chile / La Serena por **Cesar Rivas**.