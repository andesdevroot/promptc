#!/bin/bash
# ==============================================================================
# PROMPTC COMMUNITY EDITION - UNIVERSAL AUTO INSTALLER (v0.3.0)
# ==============================================================================

set -e # Detener ante cualquier error

# Configuración de Origen
REPO_URL="https://github.com/andesdevroot/promptc/releases/download/v0.3.0"
PROMPTC_DIR="$HOME/.promptc"

echo "========================================================"
echo " 🚀 INICIANDO INSTALACIÓN DE PROMPTC COMMUNITY EDITION "
echo "========================================================"
echo ""

# 1. Validación de Sistema y Arquitectura
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

if [ "$ARCH" == "x86_64" ]; then
    BINARY_ARCH="amd64"
elif [ "$ARCH" == "arm64" ]; then
    BINARY_ARCH="arm64"
else
    echo "❌ Arquitectura $ARCH no soportada actualmente."
    exit 1
fi

BINARY_NAME="promptc-$OS-$BINARY_ARCH"
DOWNLOAD_URL="$REPO_URL/$BINARY_NAME"

# 2. Solicitar GEMINI_API_KEY al usuario
echo "Para usar PROMPTC necesitas tu llave de Google AI Studio."
echo "Obtenla gratis en: https://aistudio.google.com/app/apikey"
read -p "🔑 Ingresa tu GEMINI_API_KEY: " GEMINI_KEY

if [ -z "$GEMINI_KEY" ]; then
    echo "❌ Error: La API Key es obligatoria."
    exit 1
fi

# 3. Preparación del Entorno Local
echo "📂 Creando entorno en $PROMPTC_DIR..."
mkdir -p "$PROMPTC_DIR"

# 4. Descarga del Binario Real desde GitHub
echo "⚙️  Descargando motor PROMPTC ($BINARY_NAME)..."
curl -L "$DOWNLOAD_URL" -o "$PROMPTC_DIR/promptc"
chmod +x "$PROMPTC_DIR/promptc"

# 5. Inicialización de Plantillas Base
if [ ! -f "$PROMPTC_DIR/templates.json" ]; then
    echo "📄 Inicializando plantillas industriales base..."
    cat <<EOF > "$PROMPTC_DIR/templates.json"
{
  "PROMPTC_BANCA_RIESGO": {
    "description": "Protocolo de mitigación de fraudes Swift.",
    "content": "ROL: {{role}}\nCONTEXTO: {{context}}\nTAREA: {{task}}\nRESTRICCIONES: {{constraints}}\nPROTOCOLO: Analiza vectores de riesgo, evalúa controles, identifica brechas normativas CMF Chile."
  },
  "PROMPTC_MINERIA_BASE": {
    "description": "Protocolo Sernageomin para drones.",
    "content": "ROL: {{role}}\nCONTEXTO: {{context}}\nTAREA: {{task}}\nRESTRICCIONES: {{constraints}}\nPROTOCOLO: Aplica normativa Sernageomin DS132, evalúa riesgos operacionales en faena."
  }
}
EOF
fi

# 6. Inyección en Claude Desktop vía Python (Seguridad Atómica)
echo "🔗 Conectando PROMPTC con Claude Desktop..."

python3 -c '
import json, os, sys

config_path = os.path.expanduser("~/Library/Application Support/Claude/claude_desktop_config.json")
promptc_dir = os.path.expanduser("~/.promptc")
gemini_key = sys.argv[1]

try:
    with open(config_path, "r") as f:
        data = json.load(f)
except Exception:
    data = {"mcpServers": {}}

if "mcpServers" not in data:
    data["mcpServers"] = {}

data["mcpServers"]["promptc"] = {
    "command": f"{promptc_dir}/promptc",
    "args": ["-mode=community"],
    "env": {
        "GEMINI_API_KEY": gemini_key
    }
}

os.makedirs(os.path.dirname(config_path), exist_ok=True)
with open(config_path, "w") as f:
    json.dump(data, f, indent=2)
' "$GEMINI_KEY"

echo ""
echo "========================================================"
echo " 🎉 ¡INSTALACIÓN COMPLETADA EXITOSAMENTE! "
echo "========================================================"
echo "👉 1. REINICIA CLAUDE DESKTOP (Cmd + Q)."
echo "👉 2. Dile a Claude: 'Usa optimize_prompt de PROMPTC...'"
echo "👉 3. Dashboard local: http://localhost:8080"
echo "========================================================"