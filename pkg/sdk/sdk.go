package sdk

import (
	"context"

	"github.com/andesdevroot/promptc/pkg/core"
	"github.com/andesdevroot/promptc/pkg/engine"
	"github.com/andesdevroot/promptc/pkg/provider"
)

type PromptC struct {
	Engine    *engine.CompilerEngine
	Optimizer core.Optimizer
}

// NewSDK inicializa el SDK con un motor y un optimizador opcional.
func NewSDK(ctx context.Context, apiKey string) (*PromptC, error) {
	gemini, err := provider.NewGeminiProvider(ctx, apiKey)
	if err != nil {
		return nil, err
	}

	return &PromptC{
		Engine:    engine.New(),
		Optimizer: gemini,
	}, nil
}

// Analyze delega al motor interno
func (s *PromptC) Analyze(p core.Prompt) core.Result {
	return s.Engine.Analyze(p)
}

// Optimize usa la IA para reparar el prompt
func (s *PromptC) Optimize(ctx context.Context, p core.Prompt) (string, error) {
	res := s.Analyze(p)
	if res.IsReliable {
		return s.Engine.Compile(p)
	}
	return s.Optimizer.Optimize(ctx, p, res.Issues)
}
