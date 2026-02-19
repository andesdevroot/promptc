package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// AppConfig representa la estructura del archivo ~/.promptc/config.yaml
type AppConfig struct {
	Provider string `yaml:"provider"`
	APIKey   string `yaml:"api_key"`
}

// getConfigPath resuelve la ruta absoluta al archivo de configuración del usuario
func getConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	dir := filepath.Join(home, ".promptc")
	// Creamos el directorio si no existe (0755 = rwxr-xr-x)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	return filepath.Join(dir, "config.yaml"), nil
}

// Load lee la configuración del disco. Si no existe, devuelve una estructura vacía.
func Load() (AppConfig, error) {
	var cfg AppConfig
	path, err := getConfigPath()
	if err != nil {
		return cfg, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil // Es válido que no exista aún
		}
		return cfg, err
	}

	err = yaml.Unmarshal(data, &cfg)
	return cfg, err
}

// Save escribe la configuración en el disco de forma segura.
func Save(cfg AppConfig) error {
	path, err := getConfigPath()
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(&cfg)
	if err != nil {
		return err
	}

	// 0600 = rw------- (Solo el dueño puede leer/escribir. Fundamental para guardar API Keys)
	return os.WriteFile(path, data, 0600)
}
