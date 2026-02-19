#!/bin/bash
set -e

# Colores para la terminal
GREEN='\033[0;32m'
CYAN='\033[0;36m'
RED='\033[0;31m'
NC='\033[0m' # Sin color

echo -e "${CYAN}🚀 Iniciando la instalación de PromptC...${NC}"

# 1. Detectar Sistema Operativo y Arquitectura
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

# Normalizar la arquitectura a la nomenclatura de Go
if [ "$ARCH" = "x86_64" ]; then
    ARCH="amd64"
elif [ "$ARCH" = "arm64" ] || [ "$ARCH" = "aarch64" ]; then
    ARCH="arm64"
else
    echo -e "${RED}❌ Arquitectura no soportada: $ARCH${NC}"
    exit 1
fi

# Validar SO soportado
if [ "$OS" != "darwin" ] && [ "$OS" != "linux" ]; then
    echo -e "${RED}❌ Sistema operativo no soportado por este instalador: $OS${NC}"
    exit 1
fi

# 2. Configurar variables de descarga
REPO="andesdevroot/promptc"
VERSION="v0.1.0-alpha" # En el futuro, este script puede consultar la API de GitHub para obtener el "latest"
BINARY_NAME="promptc-${OS}-${ARCH}"
DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${BINARY_NAME}"

# 3. Descargar el binario
echo -e "⬇️  Descargando binario para ${OS}-${ARCH}..."
if ! curl -fsSL -o promptc "$DOWNLOAD_URL"; then
    echo -e "${RED}❌ Error al descargar el binario. Verifica tu conexión o la URL: $DOWNLOAD_URL${NC}"
    exit 1
fi

# 4. Dar permisos y mover al PATH
echo -e "📦 Instalando en /usr/local/bin (puede pedir tu contraseña de administrador)..."
chmod +x promptc
sudo mv promptc /usr/local/bin/promptc

echo -e "${GREEN}✅ ¡PromptC instalado con éxito!${NC}"
echo -e "Ejecuta ${CYAN}promptc${NC} en tu terminal para comenzar."
