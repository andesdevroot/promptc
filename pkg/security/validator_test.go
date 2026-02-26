package security

import (
	"strings"
	"testing"
)

func TestResourceValidator_Validate(t *testing.T) {
	v := NewResourceValidator()

	t.Run("Prompt Seguro", func(t *testing.T) {
		err := v.Validate("Este es un prompt normal.")
		if err != nil {
			t.Errorf("No debería fallar con un prompt pequeño: %v", err)
		}
	})

	t.Run("Prompt Bomba (Excede límites)", func(t *testing.T) {
		// Creamos un string que excede los 100,000 caracteres
		hugeContent := strings.Repeat("A", MaxPromptLength+1)
		err := v.Validate(hugeContent)
		if err == nil {
			t.Error("Debería haber fallado por exceso de tamaño")
		}
	})
}
