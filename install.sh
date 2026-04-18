#!/bin/bash
set -e

CYAN='\033[1;36m'
GREEN='\033[1;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${CYAN}"
echo "    ____  ____  ____  __  ______  __________ "
echo "   / __ \/ __ \/ __ \/  |/  / _ \/_  __/ __/ "
echo "  / /_/ / /_/ / / / / /|_/ / /_/ // / / /    "
echo " / ____/ _, _/ /_/ / /  / / ____// / / /___  "
echo "/_/   /_/ |_|\____/_/  /_/_/    /_/  \____/  "
echo -e "${NC}"
echo "=> Iniciando instalación de PROMPTC v0.3.1 (Codex Edition)..."
echo ""

if ! command -v go >/dev/null 2>&1; then
    echo -e "${RED}[FATAL] Go no está instalado.${NC} Descárgalo desde https://go.dev/dl/ e inténtalo de nuevo."
    exit 1
fi

if ! command -v git >/dev/null 2>&1; then
    echo -e "${RED}[FATAL] Git no está instalado.${NC} Es necesario para clonar PROMPTC."
    exit 1
fi

PROMPTC_DIR="$HOME/.promptc"
echo "=> Creando directorio base en $PROMPTC_DIR..."
mkdir -p "$PROMPTC_DIR"

echo -n "=> Pega tu OPENAI_API_KEY (Enter para omitir fallback cloud): "
read -r USER_OPENAI_KEY < /dev/tty || true

echo -n "=> Modelo OpenAI [gpt-5.4-mini]: "
read -r USER_OPENAI_MODEL < /dev/tty || true
USER_OPENAI_MODEL="${USER_OPENAI_MODEL:-gpt-5.4-mini}"

echo -n "=> Pega tu GEMINI_API_KEY (Enter para fallback secundario opcional): "
read -r USER_GEMINI_KEY < /dev/tty || true

DEFAULT_REMOTE_IP="100.90.6.101"
echo -n "=> PROMPTC_MACMINI_IP [$DEFAULT_REMOTE_IP]: "
read -r USER_REMOTE_IP < /dev/tty || true
USER_REMOTE_IP="${USER_REMOTE_IP:-$DEFAULT_REMOTE_IP}"

echo "=> Descargando código fuente desde GitHub (rama master)..."
TEMP_DIR=$(mktemp -d)
git clone -q -b master https://github.com/andesdevroot/promptc.git "$TEMP_DIR"

echo "=> Compilando binario estático optimizado..."
cd "$TEMP_DIR"
go build -ldflags="-s -w" -o "$PROMPTC_DIR/promptc" ./cmd/promptc/main.go

ENV_ARGS=(--env "PROMPTC_MCP_CLIENT=codex-desktop" --env "PROMPTC_MACMINI_IP=$USER_REMOTE_IP")
if [ -n "$USER_OPENAI_KEY" ]; then
    ENV_ARGS+=(--env "OPENAI_API_KEY=$USER_OPENAI_KEY" --env "OPENAI_MODEL=$USER_OPENAI_MODEL")
fi
if [ -n "$USER_GEMINI_KEY" ]; then
    ENV_ARGS+=(--env "GEMINI_API_KEY=$USER_GEMINI_KEY")
fi

MANUAL_CMD="codex mcp add PROMPTC ${ENV_ARGS[*]} -- $PROMPTC_DIR/promptc"

if command -v codex >/dev/null 2>&1; then
    echo "=> Registrando PROMPTC como servidor MCP en Codex..."
    codex mcp remove PROMPTC >/dev/null 2>&1 || true
    codex mcp add PROMPTC "${ENV_ARGS[@]}" -- "$PROMPTC_DIR/promptc"
    CONFIG_STATUS="Codex configurado automáticamente"
else
    CONFIG_STATUS="Codex no detectado; configuración manual requerida"
fi

cat > "$PROMPTC_DIR/codex-mcp-setup.sh" <<EOF
#!/bin/bash
$MANUAL_CMD
EOF
chmod +x "$PROMPTC_DIR/codex-mcp-setup.sh"

rm -rf "$TEMP_DIR"

echo ""
echo -e "${GREEN}[SUCCESS] ¡PROMPTC Codex Edition instalado exitosamente!${NC}"
echo "--------------------------------------------------------"
echo " • Binario instalado en: $PROMPTC_DIR/promptc"
echo " • Estado MCP: $CONFIG_STATUS"
if [ -n "$USER_OPENAI_KEY" ]; then
    echo " • Cloud primario: OpenAI ($USER_OPENAI_MODEL)"
fi
if [ -n "$USER_GEMINI_KEY" ]; then
    echo " • Cloud secundario: Gemini fallback"
fi
echo " • Script de apoyo: $PROMPTC_DIR/codex-mcp-setup.sh"
echo "--------------------------------------------------------"

if ! command -v codex >/dev/null 2>&1; then
    echo -e "${YELLOW}>> PASO FINAL:${NC} Instala Codex y luego ejecuta:"
    echo "   $MANUAL_CMD"
else
    echo -e "${YELLOW}>> PASO FINAL:${NC} Reinicia Codex para que tome el nuevo servidor MCP."
fi
