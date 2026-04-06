#!/bin/sh
set -e

# Start warm `golit serve` unless SSR is disabled (e.g. benchmarks with GOLIT_DISABLED=1).
if [ -z "${GOLIT_DISABLED:-}" ]; then
  /usr/local/bin/golit serve --defs "${GOLIT_DEFS}" --listen 127.0.0.1:9777 &
  _golit_pid=$!
  _i=0
  _ready=0
  while [ "$_i" -lt 60 ]; do
    if wget -q -O- http://127.0.0.1:9777/health >/dev/null 2>&1; then
      _ready=1
      break
    fi
    _i=$((_i + 1))
    sleep 0.1
  done
  if [ "$_ready" != 1 ]; then
    echo "docker-entrypoint: golit serve did not become ready on /health (timeout)" >&2
    kill "$_golit_pid" 2>/dev/null || true
    exit 1
  fi
  export GOLIT_SERVE_URL=http://127.0.0.1:9777
fi

exec "$@"
