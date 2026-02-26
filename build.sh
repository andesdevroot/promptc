#!/bin/bash
# ==============================================================================
# PROMPTC COMPILER - BUILD SCRIPT
# Genera binarios estáticos para distribución masiva
# ==============================================================================

set -e

echo "🚀 Iniciando compilación de PROMPTC v0.3.0..."

# Crear carpeta de salida limpia
rm -rf build/
mkdir -p build

# 1. Mac OS (Apple Silicon - M1/M2/M3) -> El tuyo y la mayoría moderna
echo "⚙️  Compilando para macOS (Apple Silicon)..."
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o build/promptc-darwin-arm64 ./cmd/promptc/main.go

# 2. Mac OS (Intel) -> Macs antiguas
echo "⚙️  Compilando para macOS (Intel)..."
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o build/promptc-darwin-amd64 ./cmd/promptc/main.go

# 3. Linux (AMD64) -> Servidores, Ubuntu, etc.
echo "⚙️  Compilando para Linux (x86_64)..."
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o build/promptc-linux-amd64 ./cmd/promptc/main.go

echo "✅ Compilación exitosa. Los binarios están en la carpeta 'build/':"
ls -lh build/
echo "====================================================================="
echo "👉 Siguiente paso: Sube estos 3 archivos a la sección 'Releases' de GitHub."