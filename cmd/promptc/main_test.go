package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestMetricsInitialization(t *testing.T) {
	// 1. Setup de entorno temporal
	tempDir := t.TempDir()
	metricsPath = filepath.Join(tempDir, "metrics.json")

	// 2. Simular un estado guardado
	snap := MetricsSnapshot{
		InferenceCount:   150,
		InferenceSuccess: 145,
		TotalTokens:      50000,
		SavedAt:          time.Now(),
	}

	data, err := json.Marshal(snap)
	if err != nil {
		t.Fatalf("Error preparando test: %v", err)
	}
	os.WriteFile(metricsPath, data, 0644)

	// 3. Ejecutar función a probar
	loadMetrics()

	// 4. Validar que la memoria en ejecución cargó los datos
	if metrics.InferenceCount != 150 {
		t.Errorf("Esperaba 150 inferencias, obtuve %d", metrics.InferenceCount)
	}
	if metrics.TotalTokens != 50000 {
		t.Errorf("Esperaba 50000 tokens, obtuve %d", metrics.TotalTokens)
	}
}

func TestSnapshotAPI(t *testing.T) {
	metrics.Lock()
	metrics.CurrentMode = "community"
	metrics.NodeOnline = false
	metrics.Unlock()

	apiResult := metricsSnapshotAPI()

	if apiResult["mode"] != "community" {
		t.Errorf("El modo esperado era 'community', obtuve %v", apiResult["mode"])
	}
	if apiResult["node_online"] != false {
		t.Errorf("El nodo debía estar offline en modo community")
	}
}
