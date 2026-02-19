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

📂 Estructura del Proyecto
Este proyecto sigue el Standard Go Project Layout para garantizar escalabilidad:

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


🛣️ Roadmap
[x] v0.1.0: Core Engine, YAML Parser & CLI UI profesional con Cobra.
[ ] v0.2.0: Integración con Google Gemini API para auto-corrección de prompts (promptc fix).
[ ] v0.3.0: Soporte para múltiples targets (Formato optimizado para OpenAI vs Anthropic).
[ ] v0.4.0: Sistema de Plugins para reglas personalizadas de negocio.
[ ] v1.0.0: Lanzamiento oficial con instaladores binarios para macOS/Linux.


🤝 Contribuciones
Este es un proyecto Open Source nacido en Chile 🇨🇱 con la visión de elevar el estándar del AI Engineering.

¡Tu ayuda es bienvenida para convertir el prompting en ingeniería real!

1. Haz un Fork del repositorio.

2. Crea una rama para tu feature (git checkout -b feature/AmazingFeature).

3. Haz Commit de tus cambios con mensajes claros.

4. Haz Push a la rama.

5. Abre un Pull Request.

📄 Licencia
Distribuido bajo la Licencia MIT. Consulta el archivo LICENSE para más información.

Maintained by Cesar Rivas