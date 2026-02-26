package security

import (
	"errors"
	"fmt"
)

// MaxPromptSize define el límite físico (100KB es un estándar seguro para prompts)
const MaxPromptSize = 100 * 1024

var ErrPromptTooLarge = errors.New("security: prompt exceeds maximum allowed size")

type ResourceValidator struct {
	MaxSize int
}

func NewResourceValidator() *ResourceValidator {
	return &ResourceValidator{MaxSize: MaxPromptSize}
}

func (v *ResourceValidator) Validate(content string) error {
	// Validamos por bytes, que es lo que realmente afecta la RAM del Mac Mini
	if len(content) > v.MaxSize {
		return fmt.Errorf("%w: size %d bytes, limit %d bytes",
			ErrPromptTooLarge, len(content), v.MaxSize)
	}
	return nil
}
