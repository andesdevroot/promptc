package security

import (
	"testing"
)

func TestMasker_Mask(t *testing.T) {
	masker := NewMasker()

	tests := []struct {
		name     string
		input    string
		expected string
		foundPII bool
	}{
		{
			name:     "Detección de RUT",
			input:    "El cliente con RUT 12.345.678-9 solicitó un crédito.",
			expected: "El cliente con RUT [RUT_HIDDEN] solicitó un crédito.",
			foundPII: true,
		},
		{
			name:     "Detección de Email",
			input:    "Enviar reporte a cesar.rivas@empresa.cl mañana.",
			expected: "Enviar reporte a [EMAIL_HIDDEN] mañana.",
			foundPII: true,
		},
		{
			name:     "Texto Limpio",
			input:    "Optimiza este código Go.",
			expected: "Optimiza este código Go.",
			foundPII: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, detected := masker.Mask(tt.input)
			if result != tt.expected {
				t.Errorf("Mask() falló. Esperado: %s, Obtenido: %s", tt.expected, result)
			}
			if (len(detected) > 0) != tt.foundPII {
				t.Errorf("Detección de PII incorrecta. Esperado: %v, Obtenido: %v", tt.foundPII, len(detected) > 0)
			}
		})
	}
}
