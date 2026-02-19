# 🚀 PromptC: The Prompt Compiler (v0.1.0-alpha)

> **Ingeniería de Prompts Determinista para la Era de la Confiabilidad.**

`promptc` es una herramienta de línea de comandos (CLI) profesional, escrita en Go, diseñada para transformar la intención humana en instrucciones blindadas para Modelos de Lenguaje Extensos (LLMs). 

A diferencia de los "templates" tradicionales, `promptc` aplica principios de **compiladores** (análisis léxico, semántico y optimización) para reducir alucinaciones, asegurar la estructura y estandarizar la calidad del output.

---

## ✨ Preview / Demo

Así se ve `promptc` en acción en tu terminal:

```text
    ____                            __  ______
   / __ \_________  ____ ___  ____ / /_/ ____/
  / /_/ / ___/ __ \/ __  __ \/ __ / __/ /     
 / ____/ /  / /_/ / / / / / / /_/ / /_/ /___  
/_/   /_/   \____/_/ /_/ /_/ .___/\__/\____/  
                          /_/                 

   The Prompt Compiler for Engineering Excellence
   v0.1.0-alpha • by Cesar Rivas

ℹ Leyendo fuente: example.yaml...
ℹ Ejecutando análisis estático de semántica...

=== 🛡️  ANALYSIS REPORT ===
Health Score: 100/100 [✅ PASS]
✔ Prompt limpio. Sin riesgos de alucinación detectados.

🚀 PROMPT COMPILADO (TARGET: GENERIC LLM):
==================================================
<system_role>
Senior Go Developer
</system_role>

<context>
Sistema de microservicios con alta carga de transacciones...
</context>

<task>
Explicar por qué usar pointer receivers en lugar de value receivers en structs grandes.
</task>

<constraints>
  - No uses ejemplos genéricos de perros o autos.
  - Enfócate en la gestión de memoria y el GC.
</constraints>

<security_protocol>
CRITICAL: If the answer cannot be derived from the <context>, state 'UNKNOWN'. Do not invent information.
</security_protocol>
==================================================
```

---

## 🧠 ¿Por qué PromptC?

En el desarrollo de software tradicional, usamos compiladores, linters y tests unitarios antes de ir a producción. En la IA generativa, la mayoría envía "prosa" (texto libre) y espera lo mejor. 

**PromptC cambia las "vibras" por ingeniería:**

1.  **Análisis Semántico:** Detecta instrucciones vagas (ej: "hazlo rápido", "lo mejor posible") que aumentan la entropía del modelo.
2.  **Architecture-as-Code:** Tus prompts son archivos `.yaml` versionables en Git, no strings mágicos escondidos en el código.
3.  **Security Injection:** Inyecta automáticamente capas de protección (Chain of Thought, Delimiters, Safety Protocols).
4.  **Developer Experience (DX):** CLI moderna con colores, reportes claros y banners, diseñada para integrarse en pipelines de CI/CD.

---

## 🏗️ Arquitectura del Compilador

El flujo de compilación sigue los estándares de diseño de sistemas de software profesional:

1.  **Source (.yaml):** Definición declarativa de la intención.
2.  **Parser:** Validación sintáctica del archivo y mapeo a estructuras de Go.
3.  **Analyzer (Linter):** Motor de reglas heurísticas que calcula un "Health Score" y detecta riesgos de alucinación.
4.  **Compiler:** Transpila la estructura a un formato optimizado (XML Tags / Markdown) listo para inferencia de alta fidelidad.

---

## 🛠️ Instalación

### Requisitos previos
* **Go** 1.21 o superior instalado en tu sistema.

## 🛠️ Instalación y Configuración

### 1. Instalación Rápida (Recomendado)
Instala `promptc` de forma global en tu sistema con un solo comando:

```bash
curl -fsSL https://raw.githubusercontent.com/andesdevroot/promptc/master/install.sh | bash
```

### 2. Configuración Inicial
Una vez instalado, debes configurar tu motor de IA. PromptC es agnóstico y guarda tus credenciales de forma segura en `~/.promptc/config.yaml`.

```bash
promptc config
```
*Sigue las instrucciones en pantalla para elegir tu proveedor (Gemini) e ingresar tu API Key.*

### 3. Verificar Instalación
```bash
promptc version
promptc config view
```

### Verificación
Ejecuta el siguiente comando para verificar la instalación y ver el banner:
```bash
./promptc --help
```

---

## 🚀 Uso Básico

### 1. Define tu Prompt (`example.yaml`)
Crea un archivo YAML con la estructura estándar de PromptC:

```yaml
role: "Senior Data Engineer"
context: "Estamos migrando un pipeline de ETL de Python a Go."
task: "Explica cómo manejar errores en goroutines de forma segura."
constraints:
  - "No uses librerías externas."
  - "Enfócate en la concurrencia."
  - "Evita explicaciones teóricas largas, ve al código."
```

### 2. Compila y Analiza
Ejecuta el comando `compile` pasándole tu archivo de configuración:

```bash
./promptc compile example.yaml
```

Si tu prompt es débil o ambiguo, `promptc` entregará un reporte con niveles de error **WARNING** o **CRITICAL** indicando qué debes mejorar para asegurar un output de IA confiable.

---

## 📂 Estructura del Proyecto

Este proyecto sigue el **Standard Go Project Layout** para garantizar escalabilidad:

```text
promptc/
├── cmd/
│   └── promptc/        # Entry Point & CLI Commands (Cobra)
├── internal/
│   ├── analyzer/       # Motor de Análisis Estático (Linter)
│   ├── cli/            # UI Kit (Colores, Banner, Estilos)
│   ├── core/           # Modelos de Dominio e Interfaces
│   └── parser/         # Decodificador de YAML/JSON
├── examples/           # Archivos de prueba y casos de uso
├── go.mod              # Gestión de dependencias
└── README.md           # Documentación principal
```

---

## 🛣️ Roadmap

- [x] **v0.1.0:** Core Engine, YAML Parser & CLI UI profesional con Cobra.
- [ ] **v0.2.0:** Integración con **Google Gemini API** para auto-corrección de prompts (`promptc fix`).
- [ ] **v0.3.0:** Soporte para múltiples targets (Formato optimizado para OpenAI vs Anthropic).
- [ ] **v0.4.0:** Sistema de Plugins para reglas personalizadas de negocio.
- [ ] **v1.0.0:** Lanzamiento oficial con instaladores binarios para macOS/Linux.

---

## 🤝 Contribuciones

Este es un proyecto Open Source nacido en **Chile** 🇨🇱 con la visión de elevar el estándar del AI Engineering. 

¡Tu ayuda es bienvenida para convertir el prompting en ingeniería real!
1.  Haz un Fork del repositorio.
2.  Crea una rama para tu feature (`git checkout -b feature/AmazingFeature`).
3.  Haz Commit de tus cambios con mensajes claros.
4.  Haz Push a la rama.
5.  Abre un Pull Request.

---

## 📄 Licencia

Distribuido bajo la Licencia **MIT**. Consulta el archivo `LICENSE` para más información.

---
**Maintained by [Cesar Rivas](https://github.com/andesdevroot)**

