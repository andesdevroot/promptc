package core

// PromptSource representa la estructura base del archivo YAML de entrada.
// Este es el modelo de dominio central que usan el Parser, el Analyzer y el LLM.
type PromptSource struct {
	Role        string   `yaml:"role" json:"role"`
	Context     string   `yaml:"context" json:"context"`
	Task        string   `yaml:"task" json:"task"`
	Constraints []string `yaml:"constraints" json:"constraints"`
}
