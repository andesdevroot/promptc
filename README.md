# 🚀 PROMPTC: The Industrial Prompt Compiler

```text
    ____                            __  ______
   / __ \_________  ____ ___  ____ / /_/ ____/
  / /_/ / ___/ __ \/ __  __ \/ __ / __/ /     
 / ____/ /  / /_/ / / / / / / /_/ / /_/ /___  
/_/   /_/   \____/_/ /_/ /_/ .___/\__/\____/  
                          /_/                 

   The Prompt Compiler for Engineering Excellence
   v0.3.1 • by Cesar Rivas
```

![Version](https://img.shields.io/badge/version-v0.3.1-blue.svg)
![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8.svg?logo=go)
![Architecture](https://img.shields.io/badge/Architecture-MCP_Dual--Tier-00ff41.svg)
![Build](https://img.shields.io/badge/build-passing-brightgreen.svg)
![License](https://img.shields.io/badge/license-MIT-green.svg)

> **"Prompt Engineering is Software Engineering."**

**PROMPTC** es un compilador nativo y orquestador L7 desarrollado en Go. Resuelve el vacío de soberanía, latencia y determinismo en la adopción de IA Generativa, permitiendo a los desarrolladores y equipos corporativos inyectar contexto, resolver variables y aplicar restricciones de negocio estables antes de que el LLM procese la solicitud.

---

## ⚡ Instalación Rápida (Community Edition)

¿Quieres probar PROMPTC en tu entorno local en menos de 60 segundos? Hemos diseñado un auto-instalador *Plug & Play* que configura el motor y lo conecta con tu **Claude Desktop** automáticamente.

**Requisitos previos:**
* macOS (M1/M2/M3 o Intel) o Linux.
* Claude Desktop instalado.
* Una [API Key de Google AI Studio](https://aistudio.google.com/app/apikey) (gratuita).

Abre tu terminal y ejecuta:

```bash
curl -sSL https://raw.githubusercontent.com/andesdevroot/promptc/main/install.sh | bash
```

Una vez finalizado, **reinicia Claude Desktop** (Cmd + Q) y pídele:
*"Usa la herramienta optimize_prompt de PROMPTC para crear un protocolo usando el template PROMPTC_MINERIA_BASE"*.

---

## 🎛️ Core Dashboard (Local Observability)

PROMPTC incluye un panel de control local de grado industrial (`http://localhost:8080`) que te permite:
- Monitorear el consumo de tokens y latencia en tiempo real.
- Auditar el flujo de datos y decisiones de ruteo (*Audit Log Stream*).
- **Hot-Reload de Templates:** Edita tus plantillas industriales en formato JSON y aplícalas en caliente sin reiniciar el servidor MCP ni Claude.

---

## 🏗️ Arquitectura: Private-First AI (Dual-Tier)

En sectores regulados como la **Minería (Sernageomin), Banca (CMF) y el Sector Legal**, la lógica de negocio es un activo crítico que no puede ser expuesto a nubes públicas de forma descontrolada. PROMPTC actúa como un **L7 Gateway para LLMs**, adaptándose a tus requerimientos:

### 1. Community Mode (Rápida Adopción)
* **Orquestación:** Ejecuta el binario localmente de forma ultra-ligera.
* **Inferencia:** Rutea la optimización hacia **Gemini 1.5 Pro** de forma transparente.
* **Uso:** Ideal para startups, desarrolladores y flujos ágiles.

### 2. Enterprise Mode (Air-Gapped / Soberanía Total)
* **Orquestación:** Ejecuta el binario localmente interceptando el prompt.
* **Inferencia:** Rutea estrictamente hacia nodos locales (ej. Mac Mini corriendo Ollama / Llama 3) vía redes privadas virtuales como **Tailscale** (`100.x.x.x`).
* **Seguridad:** Cero telemetría externa. Los datos de la compañía nunca tocan la internet pública.

---

## ✨ Características Principales

* **Servidor MCP Nativo**: Implementación completa del *Model Context Protocol* sobre JSON-RPC 2.0. Integración transparente con Claude Desktop y Cursor.
* **Binario Estático en Go**: Cero dependencias (Runtime-free). Rendimiento de alto nivel con consumo mínimo de recursos (RAM < 15MB).
* **Prompt-as-Code (PaC)**: Gestión de plantillas mediante componentes versionables pre-certificados (`PROMPTC_MINERIA_BASE`, `PROMPTC_BANCA_RIESGO`).
* **Compilación Determinista**: Transforma lenguaje ambiguo en estructuras Markdown blindadas (Role, Context, Task, Constraints) inyectando variables dinámicas (`{{variable}}`).

---

## 🛠️ Herramientas MCP Expuestas

PROMPTC expone las siguientes funciones al agente de tu IDE/Chat:

1. `get_template`: Extrae una plantilla industrial pre-aprobada desde tu almacén local (`~/.promptc/templates.json`).
2. `optimize_prompt`: El motor central. Inyecta contexto, resuelve variables en tiempo de compilación y emite un prompt determinista listo para inferencia de alta precisión.

---

## 💻 Para Contribuidores (TDD & Build)

PROMPTC se construye bajo una estricta política de **Test-Driven Development (TDD) y Cero Regresiones**. Si deseas compilar el código fuente o aportar heurísticas:

```bash
# 1. Clonar el repositorio
git clone https://github.com/andesdevroot/promptc.git
cd promptc

# 2. Ejecutar la suite de tests (Obligatorio antes de compilar)
go test -v ./...

# 3. Compilar el motor optimizado
go build -ldflags="-s -w" -o build/promptc ./cmd/promptc/main.go
```

---

## 🤝 Filosofía y Licencia

Este es un proyecto Open Source diseñado para fortalecer el desarrollo de IA determinista y soberana en la región. Las contribuciones son bienvenidas vía Pull Requests, siempre que incluyan su respectiva cobertura de tests.

**Licencia**: MIT  
**Autor**: Cesar Rivas - Senior Software Engineer.  
*Desarrollado en La Serena, Chile.* 🇨🇱

