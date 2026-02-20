package core

import "context"

// Optimizer define el contrato para cualquier IA que quiera mejorar un prompt.
type Optimizer interface {
	Name() string
	Optimize(ctx context.Context, p Prompt, issues []string) (string, error)
}
