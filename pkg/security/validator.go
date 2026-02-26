package security

import (
	"errors"
	"fmt"
)

const (
	// MaxPromptLength define el límite de seguridad (aprox 100KB)
	MaxPromptLength = 100000
)

var (
	ErrPromptTooLarge = errors.New("security: prompt content exceeds safety limits")
)

// ResourceValidator gestiona los límites físicos del compilador.
type ResourceValidator struct {
	MaxLength int
}

func NewResourceValidator() *ResourceValidator {
	return &ResourceValidator{
		MaxLength: MaxPromptLength,
	}
}

// Validate comprueba si el prompt es seguro para procesar.
func (v *ResourceValidator) Validate(content string) error {
	if len(content) > v.MaxLength {
		return fmt.Errorf("%w: current size %d, max allowed %d",
			ErrPromptTooLarge, len(content), v.MaxLength)
	}
	return nil
}
