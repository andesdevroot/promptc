package core

import (
	"errors"

	"github.com/andesdevroot/promptc/internal/analyzer"
	"github.com/andesdevroot/promptc/internal/models"
)

var (
	ErrMissingRole = errors.New("el rol del sistema es obligatorio ")
	ErrTooVague    = errors.New("el prompt es demasiado corto y causará alucionaciones")
)

type Compiler struct{}

func NewCompiler() *Compiler {
	return &Compiler{}
}

// Compile toma el origen y devuelve el prompt optimizado o un error
func (c *Compiler) Compile(ps models.PromptSource) (*models.Result, error) {
	// 1. Validaciones Críticas
	if ps.Role == "" {
		return nil, ErrMissingRole
	}
	if len(ps.Task) < 10 {
		return nil, ErrTooVague
	}

	// 2. Analisis y Scoring
	warnings, score := analyzer.Analyze(ps)

	// 3. Generacion (BINARIO)
	output := c.render(ps)

	return &models.Result{
		RawPrompt: output,
		Score:     score,
		Warnings:  warnings,
	}, nil

}

func (c *Compiler) render(ps models.PromptSource) string {
	// Implementación básica de renderizado
	constraintsStr := ""
	for _, constraint := range ps.Constraints {
		constraintsStr += "- " + constraint + "\n"
	}
	return ps.Role + "\n\n" + ps.Task + "\n\nContexto: " + ps.Context + "\n\nRestricciones:\n" + constraintsStr
}
