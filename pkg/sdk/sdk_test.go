package sdk

import (
	"context"
	"strings"
	"testing"

	"github.com/andesdevroot/promptc/pkg/core"
)

func TestSDKInitialization(t *testing.T) {
	// Evaluamos la creación segura sin dependencias externas
	app, err := NewSDK(context.Background(), "", "")

	if err != nil {
		t.Fatalf("Error inesperado al inicializar SDK vacío: %v", err)
	}

	if app.Engine == nil {
		t.Errorf("Fallo crítico: El CompilerEngine no fue instanciado en el SDK")
	}

	if len(app.Optimizers) != 0 {
		t.Errorf("Se esperaban 0 optimizers (claves vacías), se obtuvieron %d", len(app.Optimizers))
	}
}

func TestSDKCompileAndOptimizeBypass(t *testing.T) {
	// Instanciamos el SDK vacío (sin optimizadores configurados)
	app, _ := NewSDK(context.Background(), "", "")

	prompt := core.Prompt{
		Role: "Arquitecto Cloud",
		Task: "Prueba de enrutamiento y bypass.",
	}

	// Como el Engine actual retorna IsReliable=true (Mock de calidad al 100%),
	// el orquestador DEBE bypassar las llamadas de red y retornar la compilación cruda.
	result, err := app.CompileAndOptimize(context.Background(), prompt)

	if err != nil {
		t.Fatalf("CompileAndOptimize arrojó error durante el bypass: %v", err)
	}

	// Validamos que el resultado provenga efectivamente del CompilerEngine (formato Markdown)
	expectedRoleBlock := "### ROLE\nArquitecto Cloud"
	if !strings.Contains(result, expectedRoleBlock) {
		t.Errorf("El enrutador no devolvió el string compilado correctamente.\nFalta: %q\nObtuvo: %q", expectedRoleBlock, result)
	}
}
