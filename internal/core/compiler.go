package core

import "errors"

var (
	ErrMissingRole = errors.New("el rol del sistema es obligatorio ")
	ErrTooVague    = errors.New("el prompt es demasiado corto y causará alucionaciones")
)

type Compiler struct{}

func NewCompiler() *Compiler {
	return &Compiler{}
}

// Compile toma el origen y devuelve el prompt optimizado o un error
func (c *Compiler) Compile(ps PromptSource) (*Result, error) {
	// 1. Validaciones Críticas
	if ps.Role == "" {
		return nil, ErrMissingRole
	}
	if len(ps.Task) < 10 {
		return nil, ErrTooVague
	}

	// 2. Analisis y Scoring
	warnings := c.analyze(ps)
	score := 100 - (len(warnings) * 20)

	// 3. Generacion (BINARIO)
	output := c.render(ps)

	return &Result{
		RawPrompt: output,
		Score:     score,
		Warnings:  warnings,
	}, nil

}

func (c *Compiler) analyze(ps PromptSource) []string {
	var warnings []string

	if ps.Context == "" {
		warnings = append(warnings, "Falta contexto adicional para mejores resultados")
	}
	if len(ps.Constraints) < 2 {
		warnings = append(warnings, "Las restricciones son muy limitadas")
	}

	return warnings
}

func (c *Compiler) render(ps PromptSource) string {
	// Implementación básica de renderizado
	constraintsStr := ""
	for _, constraint := range ps.Constraints {
		constraintsStr += "- " + constraint + "\n"
	}
	return ps.Role + "\n\n" + ps.Task + "\n\nContexto: " + ps.Context + "\n\nRestricciones:\n" + constraintsStr
}
