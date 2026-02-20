package sdk

import (
	"context"
	"log"

	"github.com/andesdevroot/promptc/pkg/core"
	"github.com/andesdevroot/promptc/pkg/engine"
	"github.com/andesdevroot/promptc/pkg/provider"
)

type PromptC struct {
	Engine     *engine.CompilerEngine
	Optimizers []core.Optimizer
}

func (s *PromptC) Optimize(ctx context.Context, p core.Prompt) (any, any) {
	panic("unimplemented")
}

func (s *PromptC) Analyze(p core.Prompt) any {
	panic("unimplemented")
}

// NewSDK ahora acepta 3 argumentos para incluir tu nodo de Tailscale
func NewSDK(ctx context.Context, geminiKey string, remoteIP string) (*PromptC, error) {
	eng := engine.New()
	var optimizers []core.Optimizer

	// Prioridad: Nodo local Mac mini (Soberanía de datos)
	if remoteIP != "" {
		optimizers = append(optimizers, provider.NewOllamaProvider(remoteIP))
	}

	// Respaldo: Gemini Cloud
	if geminiKey != "" {
		g, err := provider.NewGeminiProvider(ctx, geminiKey)
		if err == nil {
			optimizers = append(optimizers, g)
		}
	}

	return &PromptC{
		Engine:     eng,
		Optimizers: optimizers,
	}, nil
}

// CompileAndOptimize es el método que main.go intentaba llamar
func (s *PromptC) CompileAndOptimize(ctx context.Context, p core.Prompt) (string, error) {
	analysis := s.Engine.Analyze(p)

	// Si el prompt es perfecto, no gastamos ciclos de GPU
	if analysis.IsReliable {
		return s.Engine.Compile(p)
	}

	// Intentamos optimizar con los proveedores disponibles
	for _, opt := range s.Optimizers {
		log.Printf("[SDK] Intentando con: %s", opt.Name())
		optimized, err := opt.Optimize(ctx, p, analysis.Issues)
		if err == nil {
			return optimized, nil
		}
		log.Printf("[SDK] Error con %s: %v", opt.Name(), err)
	}

	// Fallback: Si todo falla, devolvemos la compilación base
	return s.Engine.Compile(p)
}
