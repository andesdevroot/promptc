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
![Security](https://img.shields.io/badge/Security-PII_Masking-red.svg)
![Build](https://img.shields.io/badge/build-passing-brightgreen.svg)
![License](https://img.shields.io/badge/license-MIT-green.svg)

> **"Prompt Engineering is Software Engineering."**

🌐 *Read this in other languages:* [Español](README.es.md)

**PROMPTC** is a native compiler and Layer 7 (L7) orchestrator developed in Go. It bridges the gap between sovereignty, latency, and determinism in Generative AI adoption. It allows engineering teams to inject context, resolve variables, and enforce stable business constraints before an LLM processes the request.

---

## ⚡ Quick Start (Community Edition)

Get PROMPTC running in your local environment in under 60 seconds. Our *Plug & Play* installer automatically configures the engine and connects it with **Claude Desktop**.

**Prerequisites:**
* macOS (M1/M2/M3) or Linux.
* Claude Desktop installed.
* A [Google AI Studio API Key](https://aistudio.google.com/app/apikey).

Run in your terminal:

```bash
curl -sSL https://raw.githubusercontent.com/andesdevroot/promptc/master/install.sh | bash
```

---

## 🏗️ Architecture: Private-First AI (Dual-Tier)

In regulated sectors like **Mining, Banking, and Legal**, business logic is a critical asset. PROMPTC acts as an **L7 Gateway for LLMs**:

### 1. Community Mode (Agile Development)
* **Orchestration:** Ultra-lightweight local binary execution.
* **Inference:** Routes optimization to **Gemini 1.5 Pro** transparently.

### 2. Enterprise Mode (Air-Gapped / Full Sovereignty)
* **Orchestration:** Local interceptor for all outgoing prompts.
* **Inference:** Routes strictly to local nodes (e.g., Mac Mini running **Ollama / Llama 3**) via VPNs like **Tailscale**.
* **Security:** Zero external telemetry. Corporate data never touches the public internet.

---

## 🛡️ Industrial Security Layer

* **PII Masking**: Automatic detection and masking of sensitive data (RUT, Emails, API Keys).
* **Resource Guard**: Strict validation of prompt size to prevent infrastructure exhaustion (DoS protection).
* **Audit Logging**: Real-time stream of all compilation decisions for compliance monitoring.

---

## ✨ Key Features

* **Native MCP Server**: Full *Model Context Protocol* implementation over JSON-RPC 2.0.
* **Static Go Binary**: Zero dependencies, runtime-free, and high performance (RAM < 15MB).
* **Prompt-as-Code (PaC)**: Manage templates through versioned, pre-certified components.
* **Deterministic Compilation**: Transforms ambiguous language into structured Markdown (Role, Context, Task, Constraints).

---

## 💻 For Contributors (TDD & Build)

PROMPTC is built under a strict **Test-Driven Development (TDD)** and **Zero Regression** policy.

```bash
# 1. Clone the repository
git clone https://github.com/andesdevroot/promptc.git
cd promptc

# 2. Run the test suite
go test -v ./...

# 3. Build the optimized binary
go build -ldflags="-s -w" -o build/promptc ./cmd/promptc/main.go
```

---

## 🤝 Philosophy and License

**License**: MIT  
**Author**: Cesar Rivas - Senior Software Engineer.  
*Developed in La Serena, Chile.* 🇨🇱