#!/bin/sh
# Local dev: start `golit serve` then exec the app (same warm path as the container).
# Makefile sets GOLIT_BIN and GOLIT_DEFS.
set -e
if [ -z "${GOLIT_BIN:-}" ] || [ ! -f "$GOLIT_BIN" ]; then
  echo "dev-serve.sh: set GOLIT_BIN to the golit binary (run make serve from this directory)" >&2
  exit 1
fi
if [ -z "${GOLIT_DEFS:-}" ]; then
  echo "dev-serve.sh: set GOLIT_DEFS to the bundles directory" >&2
  exit 1
fi
"$GOLIT_BIN" serve --defs "$GOLIT_DEFS" --listen 127.0.0.1:9777 &
SERVE_PID=$!
cleanup() { kill "$SERVE_PID" 2>/dev/null || true; }
trap cleanup EXIT INT TERM
i=0
ready=0
while [ "$i" -lt 60 ]; do
  if curl -sf http://127.0.0.1:9777/health >/dev/null 2>&1; then
    ready=1
    break
  fi
  i=$((i + 1))
  sleep 0.1
done
if [ "$ready" != 1 ]; then
  echo "dev-serve.sh: golit serve did not become ready on /health (timeout)" >&2
  exit 1
fi
export GOLIT_SERVE_URL=http://127.0.0.1:9777
exec "$@"
