package models

// PromptSource representa la entrada de nuestro compilador
// Lo mantenemos agnóstico al formato (YAML/JSON)
type PromptSource struct {
	ID          string   `yaml:"id"`
	Role        string   `yaml:"role"`
	Task        string   `yaml:"task"`
	Context     string   `yaml:"context"`
	Constraints []string `yaml:"constraints"`
}

// Result es lo que nuestro compilador entrega
type Result struct {
	RawPrompt string
	Score     int      // 0-100(Calidad del prompt)
	Warnings  []string // Sugerencias para evitar alucionaciones
}
