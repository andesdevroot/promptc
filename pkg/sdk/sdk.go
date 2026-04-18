package sdk

import (
	"context"
	"fmt"
	"log"

	"github.com/andesdevroot/promptc/pkg/core"
	"github.com/andesdevroot/promptc/pkg/engine"
	"github.com/andesdevroot/promptc/pkg/provider"
)

type PromptC struct {
	Engine     *engine.CompilerEngine
	Optimizers []core.Optimizer
}

func (s *PromptC) Optimize(ctx context.Context, p core.Prompt) (string, error) {
	if s == nil || s.Engine == nil {
		return "", fmt.Errorf("sdk no inicializado")
	}

	analysis := s.Engine.Analyze(p)

	for _, opt := range s.Optimizers {
		log.Printf("[SDK] Intentando con: %s", opt.Name())
		optimized, err := opt.Optimize(ctx, p, analysis.Issues)
		if err == nil {
			return optimized, nil
		}
		log.Printf("[SDK] Error con %s: %v", opt.Name(), err)
	}

	return s.Engine.Compile(p)
}

func (s *PromptC) Analyze(p core.Prompt) core.Result {
	if s == nil || s.Engine == nil {
		return core.Result{
			Score:       0,
			IsReliable:  false,
			Issues:      []string{"sdk no inicializado"},
			Suggestions: []string{"Inicializa PROMPTC antes de ejecutar Analyze."},
		}
	}

	return s.Engine.Analyze(p)
}

// NewSDK inicializa la cadena de optimizadores:
// nodo local -> OpenAI -> Gemini.
func NewSDK(ctx context.Context, openAIKey string, openAIModel string, geminiKey string, remoteIP string) (*PromptC, error) {
	eng := engine.New()
	var optimizers []core.Optimizer

	// Prioridad: Nodo local Mac mini (Soberanía de datos)
	if remoteIP != "" {
		optimizers = append(optimizers, provider.NewOllamaProvider(remoteIP))
	}

	// Respaldo cloud primario: OpenAI
	if openAIKey != "" {
		optimizers = append(optimizers, provider.NewOpenAIProvider(openAIKey, openAIModel))
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
	if s == nil || s.Engine == nil {
		return "", fmt.Errorf("sdk no inicializado")
	}

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
