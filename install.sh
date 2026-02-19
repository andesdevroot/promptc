#!/bin/bash
set -e

GREEN='\033[0;32m'
CYAN='\033[0;36m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${CYAN}🚀 Iniciando la instalación de PromptC...${NC}"

# 1. Detectar Sistema Operativo y Arquitectura
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

if [ "$ARCH" = "x86_64" ]; then
    ARCH="amd64"
elif [ "$ARCH" = "arm64" ] || [ "$ARCH" = "aarch64" ]; then
    ARCH="arm64"
else
    echo -e "${RED}❌ Arquitectura no soportada: $ARCH${NC}"
    exit 1
fi

# 2. Configurar variables
REPO="andesdevroot/promptc"
VERSION="v0.1.0-alpha"
BINARY_NAME="promptc-${OS}-${ARCH}"
DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${BINARY_NAME}"

# 3. Descargar usando un directorio temporal (Solución al error 56)
TMP_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_DIR"' EXIT # Limpiar al terminar

echo -e "⬇️  Descargando binario para ${OS}-${ARCH}..."
if ! curl -fsSL -L -o "$TMP_DIR/promptc" "$DOWNLOAD_URL"; then
    echo -e "${RED}❌ Error al descargar el binario. Verifica la URL: $DOWNLOAD_URL${NC}"
    exit 1
fi

# 4. Instalación
echo -e "📦 Instalando en /usr/local/bin (se requiere sudo)..."
chmod +x "$TMP_DIR/promptc"
sudo mv "$TMP_DIR/promptc" /usr/local/bin/promptc

echo -e "${GREEN}✅ ¡PromptC instalado con éxito!${NC}"
echo -e "Ejecuta ${CYAN}promptc version${NC} para verificar."