package security

import (
	"regexp"
)

// SensitivePattern define una regla de detección y su reemplazo.
type SensitivePattern struct {
	Name    string
	Regex   *regexp.Regexp
	Replace string
}

// Masker es el guardián de la privacidad de PROMPTC.
type Masker struct {
	Patterns []SensitivePattern
}

// NewMasker inicializa el buscador de PII con reglas industriales.
func NewMasker() *Masker {
	return &Masker{
		Patterns: []SensitivePattern{
			{
				Name:    "RUT_CHILE",
				Regex:   regexp.MustCompile(`(?i)\b(\d{1,2}(?:\.?\d{3}){2}-?[\dkK])\b`),
				Replace: "[RUT_HIDDEN]",
			},
			{
				Name:    "EMAIL",
				Regex:   regexp.MustCompile(`(?i)\b[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}\b`),
				Replace: "[EMAIL_HIDDEN]",
			},
			{
				Name:    "CREDIT_CARD",
				Regex:   regexp.MustCompile(`\b(?:\d[ -]*?){13,16}\b`),
				Replace: "[CARD_HIDDEN]",
			},
			{
				Name:    "IP_ADDRESS",
				Regex:   regexp.MustCompile(`\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b`),
				Replace: "[IP_HIDDEN]",
			},
		},
	}
}

// Mask procesa un string y reemplaza toda la data sensible detectada.
func (m *Masker) Mask(text string) (string, []string) {
	detected := []string{}
	result := text

	for _, p := range m.Patterns {
		if p.Regex.MatchString(result) {
			result = p.Regex.ReplaceAllString(result, p.Replace)
			detected = append(detected, p.Name)
		}
	}

	return result, detected
}
