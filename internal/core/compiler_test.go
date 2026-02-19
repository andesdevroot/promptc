package core

import (
	"testing"

	"github.com/andesdevroot/promptc/internal/models"
)

func TestCompiler_Success(t *testing.T) {
	c := NewCompiler()
	src := models.PromptSource{
		Role: "Senior Developer",
		Task: "Explicar punteros en Go de forma clara, con ejemplos técnicos",
	}

	res, err := c.Compile(src)
	if err != nil {
		t.Errorf("Se esperaba éxito, se obtuvo error: %v", err)
	}

	if res.Score < 50 {
		t.Errorf("El score es demasiado bajo para un prompt válido: %d", res.Score)
	}
}

func TestCompile_MissingRole(t *testing.T) {
	c := NewCompiler()
	src := models.PromptSource{
		Role: "Senior Developer",
		Task: "Explicar punteros en Go de forma clara, con ejemplos técnicos",
	}

	res, err := c.Compile(src)
	if err != nil {
		t.Errorf("Se esperaba éxito, se obtuvo error: %v", err)
	}

	if res.Score < 50 {
		t.Errorf("El score es demasiado bajo para un prompt válido: %d", res.Score)
	}

}
