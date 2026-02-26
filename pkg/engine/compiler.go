package engine

import (
	"fmt"
	"strings"

	"github.com/andesdevroot/promptc/pkg/core"
	"github.com/andesdevroot/promptc/pkg/security"
)

type CompilerEngine struct {
	masker *security.Masker
}

func New() *CompilerEngine {
	return &CompilerEngine{
		masker: security.NewMasker(),
	}
}

func (e *CompilerEngine) Compile(p core.Prompt) (string, error) {
	// 1. Sanitización de entrada (PII Masking)
	p.Role, _ = e.masker.Mask(p.Role)
	p.Context, _ = e.masker.Mask(p.Context)
	p.Task, _ = e.masker.Mask(p.Task)

	for k, v := range p.Variables {
		p.Variables[k], _ = e.masker.Mask(v)
	}

	// 2. Resolver Task
	task := p.Task
	if p.Role != "" {
		task = strings.ReplaceAll(task, "{{role}}", p.Role)
	}
	if p.Context != "" {
		task = strings.ReplaceAll(task, "{{context}}", p.Context)
	}
	if p.Variables != nil {
		for k, v := range p.Variables {
			task = strings.ReplaceAll(task, "{{"+k+"}}", v)
		}
	}

	// 3. Construir Markdown
	var blocks []string
	if p.Role != "" {
		blocks = append(blocks, "### ROLE\n"+p.Role)
	}
	if p.Context != "" {
		blocks = append(blocks, "### CONTEXT\n"+p.Context)
	}
	blocks = append(blocks, "### TASK\n"+task)

	if len(p.Constraints) > 0 {
		var cb strings.Builder
		cb.WriteString("### CONSTRAINTS\n")
		for i, c := range p.Constraints {
			maskedC, _ := e.masker.Mask(c)
			cb.WriteString("- " + maskedC)
			if i < len(p.Constraints)-1 {
				cb.WriteString("\n")
			}
		}
		blocks = append(blocks, cb.String())
	}

	return strings.Join(blocks, "\n\n"), nil
}

func (e *CompilerEngine) Analyze(p core.Prompt) core.Result {
	var issues []string
	score := 0

	// Análisis de PII (Nueva métrica de seguridad)
	_, detected := e.masker.Mask(p.Task)
	if len(detected) > 0 {
		issues = append(issues, fmt.Sprintf("SEGURIDAD: Se detectaron datos sensibles (%s). Serán ofuscados.", strings.Join(detected, ", ")))
		// No bajamos score, pero alertamos en los issues
	}

	// Heurísticas de calidad (Pesos)
	if len(strings.TrimSpace(p.Role)) >= 5 {
		score += 20
	} else {
		issues = append(issues, "Falta ROL.")
	}
	if len(strings.TrimSpace(p.Context)) >= 10 {
		score += 20
	} else {
		issues = append(issues, "Falta CONTEXTO.")
	}

	taskContent := strings.TrimSpace(p.Task)
	if len(taskContent) >= 50 {
		score += 40
	} else {
		issues = append(issues, "TAREA débil.")
	}

	if len(p.Constraints) >= 2 {
		score += 20
	} else {
		issues = append(issues, "Sin RESTRICCIONES.")
	}

	return core.Result{
		IsReliable: score >= 75,
		Score:      score,
		Issues:     issues,
	}
}
