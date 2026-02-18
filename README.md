# 🚀 PromptC: The Prompt Compiler (v0.1.0-alpha)

**Ingeniería de Prompts Determinista para la Era de la Confiabilidad.**

`promptc` es un compilador de línea de comandos (CLI) de código abierto diseñado para transformar la intención del usuario en instrucciones blindadas para Modelos de Lenguaje Extensos (LLMs). A diferencia de los motores de plantillas simples, `promptc` aplica **análisis estático** para reducir alucinaciones, asegurar la estructura y estandarizar la calidad del output.

Desarrollado en **Go** para garantizar portabilidad, velocidad y un binario estático sin dependencias.

---

## 🧠 ¿Por qué PromptC?

En el desarrollo de software tradicional, usamos compiladores y linters para asegurar que el código sea válido antes de ejecutarse. Sin embargo, en la IA generativa, la mayoría de los desarrolladores envían "prosa" sin validar y esperan que el modelo no alucine.

**PromptC cambia el paradigma:**
* **Análisis Semántico:** Detecta instrucciones vagas (ej: "hazlo lo mejor posible") que aumentan la entropía del modelo.
* **Inyección de Seguridad:** Añade automáticamente capas de protección como *Chain of Thought* (CoT) y delimitadores estructurales.
* **Foco en Español/LATAM:** Diseñado para manejar las sutilezas y ambigüedades del idioma español en contextos profesionales.
* **Determinismo:** Mismo archivo fuente, mismo prompt optimizado.

---

## 🏗️ Arquitectura del Compilador



El flujo de compilación sigue los principios clásicos de la ingeniería de sistemas:
1.  **Frontend (Source):** Archivos `.yaml` que definen Rol, Contexto, Tarea y Restricciones.
2.  **Middle-end (Analyzer):** Linter de calidad que calcula un "Health Score" y detecta riesgos de alucinación.
3.  **Backend (Generator):** Transpila la intención a un prompt optimizado para proveedores específicos (OpenAI, Anthropic, Ollama).

---

## 🛠️ Roadmap del Proyecto

- [x] Estructura base y modelos de dominio en Go.
- [x] Linter básico de salud del prompt.
- [ ] **Fase 1:** Parser de archivos YAML (En progreso).
- [ ] **Fase 2:** CLI con `Cobra` para uso profesional.
- [ ] **Fase 3:** Sistema de plugins para reglas de anti-alucinación personalizadas.
- [ ] **Fase 4:** Exportación a múltiples formatos (JSON, XML, Markdown).

---

## 🚀 Instalación y Desarrollo

### Requisitos
* Go 1.21 o superior.

### Clonar y Probar
```bash
git clone [https://github.com/TU_USUARIO/promptc.git](https://github.com/TU_USUARIO/promptc.git)
cd promptc
go mod tidy
go test ./... -v