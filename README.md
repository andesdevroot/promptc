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




