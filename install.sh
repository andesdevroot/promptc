#!/bin/bash
set -e

# --- UI & Branding ---
CYAN='\033[1;36m'
GREEN='\033[1;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${CYAN}"
echo "    ____  ____  ____  __  ______  __________ "
echo "   / __ \/ __ \/ __ \/  |/  / _ \/_  __/ __/ "
echo "  / /_/ / /_/ / / / / /|_/ / /_/ // / / /    "
echo " / ____/ _, _/ /_/ / /  / / ____// / / /___  "
echo "/_/   /_/ |_|\____/_/  /_/_/    /_/  \____/  "
echo -e "${NC}"
echo "=> Iniciando instalación de PROMPTC v0.3.1 (Community Edition)..."
echo ""

# --- 1. Verificación de Dependencias ---
if ! command -v go &> /dev/null; then
    echo -e "${RED}[FATAL] Go no está instalado.${NC} Descárgalo desde https://go.dev/dl/ e inténtalo de nuevo."
    exit 1
fi

if ! command -v python3 &> /dev/null; then
    echo -e "${RED}[FATAL] Python3 no está instalado.${NC} Es necesario para configurar Claude."
    exit 1
fi

# --- 2. Preparación de Entorno ---
PROMPTC_DIR="$HOME/.promptc"
echo "=> Creando directorio base en $PROMPTC_DIR..."
mkdir -p "$PROMPTC_DIR"

# --- 3. Obtención de la API Key (Lectura directa del TTY) ---
echo -n "=> Pega tu GEMINI_API_KEY (Presiona Enter si prefieres configurarla luego): "
read -r USER_GEMINI_KEY < /dev/tty || true
export USER_GEMINI_KEY

# --- 4. Descarga y Compilación ---
echo "=> Descargando código fuente desde GitHub (rama master)..."
TEMP_DIR=$(mktemp -d)
git clone -q -b master https://github.com/andesdevroot/promptc.git "$TEMP_DIR"

echo "=> Compilando binario estático optimizado..."
cd "$TEMP_DIR"
go build -ldflags="-s -w" -o "$PROMPTC_DIR/promptc" ./cmd/promptc/main.go

# --- 5. Inyección de Configuración en Claude Desktop ---
echo "=> Inyectando servidor MCP en Claude Desktop..."

# Usamos comillas simples para proteger el script de Python de la expansión de Bash
python3 -c '
import json, os

path = os.path.expanduser("~/Library/Application Support/Claude/claude_desktop_config.json")
data = {"mcpServers": {}}

if os.path.exists(path):
    try:
        with open(path, "r") as f:
            data = json.load(f)
    except Exception as e:
        print(f"   [WARN] No se pudo leer config previa: {e}")

if "mcpServers" not in data:
    data["mcpServers"] = {}

api_key = os.environ.get("USER_GEMINI_KEY", "")
env_vars = {}
if api_key:
    env_vars["GEMINI_API_KEY"] = api_key

data["mcpServers"]["PROMPTC"] = {
    "command": os.path.expanduser("~/.promptc/promptc"),
    "args": ["-mode=community"],
    "env": env_vars
}

os.makedirs(os.path.dirname(path), exist_ok=True)

with open(path, "w") as f:
    json.dump(data, f, indent=2)
'

# --- 6. Limpieza y Cierre ---
rm -rf "$TEMP_DIR"

echo ""
echo -e "${GREEN}[SUCCESS] ¡PROMPTC Community Edition instalado exitosamente!${NC}"
echo "--------------------------------------------------------"
echo " • Binario instalado en: $PROMPTC_DIR/promptc"
echo " • Claude configurado (Modo: Community)"
echo "--------------------------------------------------------"
echo -e "${YELLOW}>> PASO FINAL:${NC} Reinicia Claude Desktop (Cmd + Q) para aplicar los cambios."