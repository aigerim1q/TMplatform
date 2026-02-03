#!/bin/bash
# Run the Go chatbot with HTTP API enabled (for frontend connection).
# From workspace root (TMplatform): ./run_chatbot_http.sh
# Or: cd silencendo/TMplatform && CHATBOT_HTTP=1 go run .

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$SCRIPT_DIR/silencendo/TMplatform" || exit 1

export CHATBOT_HTTP=1
export CHATBOT_HTTP_PORT="${CHATBOT_HTTP_PORT:-8080}"

if [ -f "$SCRIPT_DIR/.env" ]; then
  set -a
  # shellcheck source=/dev/null
  . "$SCRIPT_DIR/.env"
  set +a
fi

echo "Starting chatbot API on port ${CHATBOT_HTTP_PORT} (POST /api/chat, GET /api/health)"
exec go run .
