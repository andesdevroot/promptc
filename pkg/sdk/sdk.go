package sdk

import (
	"context"
	"log"

	"github.com/andesdevroot/promptc/pkg/core"
	"github.com/andesdevroot/promptc/pkg/engine"
	"github.com/andesdevroot/promptc/pkg/provider"
)

// PromptC es el corazón del sistema. Coordina el análisis local,
// la compilación y la optimización distribuida (Local o Cloud).
type PromptC struct {
	Engine     *engine.CompilerEngine
	Optimizers []core.Optimizer
}

// NewSDK inicializa el ecosistema PROMPTC detectando los nodos disponibles.
func NewSDK(ctx context.Context, geminiKey string, remoteIP string) (*PromptC, error) {
	eng := engine.New()
	var optimizers []core.Optimizer

	// 1. Capa Enterprise: Nodo Local (Ollama) via Tailscale
	if remoteIP != "" {
		optimizers = append(optimizers, provider.NewOllamaProvider(remoteIP))
	}

	// 2. Capa Community: Gemini Cloud Fallback
	if geminiKey != "" {
		g, err := provider.NewGeminiProvider(ctx, geminiKey)
		if err == nil {
			optimizers = append(optimizers, g)
		} else {
			log.Printf("[SDK_INIT_WARN] Gemini Provider no disponible: %v", err)
		}
	}

	return &PromptC{
		Engine:     eng,
		Optimizers: optimizers,
	}, nil
}

// CompileAndOptimize ejecuta el pipeline de decisión L7.
func (s *PromptC) CompileAndOptimize(ctx context.Context, p core.Prompt) (string, error) {
	// Fase 1: Análisis Heurístico Local
	analysis := s.Engine.Analyze(p)

	// Bypass: Si el prompt ya es industrial, compilamos y entregamos (Zero-Latency)
	if analysis.IsReliable {
		log.Println("[SDK_ROUTING] Calidad óptima detectada. Ejecutando BYPASS hacia Compilación Local.")
		return s.Engine.Compile(p)
	}

	// Fase 2: Optimización Inteligente
	if len(s.Optimizers) == 0 {
		log.Println("[SDK_WARN] No hay optimizadores configurados. Entregando versión cruda.")
		return s.Engine.Compile(p)
	}

	log.Printf("[SDK_ROUTING] Prompt débil detectado (Score: %d). Iniciando optimización...", analysis.Score)
	for _, issue := range analysis.Issues {
		log.Printf("  - Issue detectado: %s", issue)
	}

	for _, opt := range s.Optimizers {
		log.Printf("[SDK_INFRA] Delegando corrección a: %s", opt.Name())
		optimized, err := opt.Optimize(ctx, p, analysis.Issues)
		if err == nil {
			return optimized, nil
		}
		log.Printf("[SDK_RETRY] Fallo en %s: %v. Intentando siguiente nodo...", opt.Name(), err)
	}

	// Fallback de seguridad: Si todo falla, compilamos lo que tenemos
	log.Println("[SDK_CRITICAL] Todos los proveedores de optimización fallaron. Ejecutando compilación básica.")
	return s.Engine.Compile(p)
}
