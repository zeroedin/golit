#!/usr/bin/env bash
set -euo pipefail

# ---------------------------------------------------------------------------
# PHP + golit performance benchmark
#
# Measures server-side (curl) and optionally client-side (Chrome headless)
# metrics with and without golit SSR, then prints a comparison.
#
# Usage:
#   ./bench.sh              # 100 requests, tier 1 only
#   ./bench.sh -n 50        # 50 requests
#   ./bench.sh --browser    # include tier 2 client metrics
#   ./bench.sh --trace      # produce Chrome trace files for flame graphs
# ---------------------------------------------------------------------------

IMAGE="golit-php"
CONTAINER_NAME="golit-php-bench"
PORT=8080
REQUESTS=100
BROWSER=false
TRACE=false
ENDPOINTS=("/" "/about")
CHROME="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"

while [[ $# -gt 0 ]]; do
  case $1 in
    -n)       REQUESTS="$2"; shift 2 ;;
    --browser) BROWSER=true; shift ;;
    --trace)   TRACE=true; shift ;;
    *)         echo "Unknown option: $1"; exit 1 ;;
  esac
done

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
RESULTS_DIR="$SCRIPT_DIR/bench-results"
mkdir -p "$RESULTS_DIR"

# --- helpers ---------------------------------------------------------------

cleanup() {
  podman rm -f "$CONTAINER_NAME" &>/dev/null || true
}
trap cleanup EXIT

wait_for_healthy() {
  local url="http://localhost:$PORT/"
  local retries=30
  while ! curl -sf -o /dev/null "$url" 2>/dev/null; do
    retries=$((retries - 1))
    if [[ $retries -le 0 ]]; then
      echo "ERROR: container failed to become healthy" >&2
      exit 1
    fi
    sleep 0.3
  done
}

# Collect curl timings into a file: ttfb total_time size_download http_code
run_curl_bench() {
  local url="$1"
  local outfile="$2"
  local i=0
  > "$outfile"
  while [[ $i -lt $REQUESTS ]]; do
    curl -sf -o /dev/null \
      -w '%{time_starttransfer} %{time_total} %{size_download} %{http_code}\n' \
      "$url" >> "$outfile"
    i=$((i + 1))
  done
}

# Compute stats from a column of numbers (seconds -> ms).
# Uses sort(1) instead of gawk's asort — portable on macOS BSD awk.
compute_stats() {
  local file="$1"
  local col="$2"
  awk -v c="$col" '{ print $c * 1000 }' "$file" | sort -n | awk '
    { v[NR] = $1; s += $1 }
    END {
      n = NR
      if (n == 0) { print "- - - - - -"; exit }
      avg = s / n
      p50 = v[int(n * 0.50) + 1]
      p95 = v[int(n * 0.95) + 1]
      p99 = v[int(n * 0.99) + 1]
      printf "%.2f %.2f %.2f %.2f %.2f %.2f", v[1], avg, p50, p95, p99, v[n]
    }
  '
}

avg_bytes() {
  local file="$1"
  awk '{ s += $3 } END { printf "%.0f", (NR > 0) ? s/NR : 0 }' "$file"
}

print_header() {
  printf "\n%-24s %8s %8s %8s %8s %8s %8s\n" "$1" "min" "avg" "p50" "p95" "p99" "max"
  printf "%-24s %8s %8s %8s %8s %8s %8s\n" "$(printf '%.0s-' {1..24})" "------" "------" "------" "------" "------" "------"
}

print_row() {
  local label="$1"; shift
  printf "%-24s" "$label"
  for v in $@; do
    printf " %7s" "${v}ms"
  done
  printf "\n"
}

print_delta() {
  local label="$1"
  local with_avg="$2"
  local without_avg="$3"
  local delta pct
  delta=$(awk "BEGIN { printf \"%.2f\", $with_avg - $without_avg }")
  if awk "BEGIN { exit ($without_avg == 0) ? 0 : 1 }" 2>/dev/null; then
    pct="N/A"
  else
    pct=$(awk "BEGIN { printf \"%.1f\", (($with_avg - $without_avg) / $without_avg) * 100 }")
  fi
  printf "  %-22s %s ms (%s%%)\n" "$label" "$delta" "$pct"
}

# Collect browser metrics via Chrome headless --dump-dom
run_browser_bench() {
  local url="$1"
  local outfile="$2"
  local runs=10
  local i=0
  > "$outfile"
  while [[ $i -lt $runs ]]; do
    local dom
    dom=$("$CHROME" --headless=new --disable-gpu --dump-dom "${url}?bench" 2>/dev/null || true)
    local json
    json=$(echo "$dom" | grep -o '<pre id="golit-perf"[^>]*>[^<]*</pre>' | sed 's/<[^>]*>//g' || true)
    if [[ -n "$json" ]]; then
      echo "$json" >> "$outfile"
    fi
    i=$((i + 1))
    sleep 0.5
  done
}

# Parse browser metrics JSON lines and compute avg for a field
browser_avg() {
  local file="$1"
  local field="$2"
  awk -F'[,:}]' -v f="\"$field\"" '
    { for(i=1;i<=NF;i++) if($i ~ f) { gsub(/[^0-9.]/, "", $(i+1)); s+=$(i+1); n++ } }
    END { if(n>0) printf "%.2f", s/n; else printf "-" }
  ' "$file"
}

# Capture Chrome trace file
run_trace() {
  local url="$1"
  local outfile="$2"
  "$CHROME" --headless=new --disable-gpu \
    --trace-startup --trace-startup-file="$outfile" \
    --trace-startup-duration=5 \
    "$url" 2>/dev/null || true
  sleep 6
}

# --- main ------------------------------------------------------------------

echo "============================================="
echo "  PHP + golit Performance Benchmark"
echo "============================================="
echo "  Requests per endpoint: $REQUESTS"
echo "  Endpoints: ${ENDPOINTS[*]}"
echo "  Browser metrics: $BROWSER"
echo "  Trace capture: $TRACE"
echo "============================================="

# Build the container image
echo ""
echo "Building container image..."
make -C "$SCRIPT_DIR" container 2>&1 | tail -1

# --- Run WITH golit --------------------------------------------------------
echo ""
echo ">>> Starting container WITH golit SSR..."
cleanup
podman run -d --name "$CONTAINER_NAME" -p "$PORT:8080" "$IMAGE" >/dev/null
wait_for_healthy

for endpoint in "${ENDPOINTS[@]}"; do
  tag=$(echo "$endpoint" | tr '/' '_')
  [[ "$tag" == "_" ]] && tag="_root"
  run_curl_bench "http://localhost:$PORT$endpoint" "$RESULTS_DIR/with${tag}.csv"
done

if $BROWSER && [[ -x "$CHROME" ]]; then
  for endpoint in "${ENDPOINTS[@]}"; do
    tag=$(echo "$endpoint" | tr '/' '_')
    [[ "$tag" == "_" ]] && tag="_root"
    run_browser_bench "http://localhost:$PORT$endpoint" "$RESULTS_DIR/with${tag}_browser.json"
  done
fi

if $TRACE && [[ -x "$CHROME" ]]; then
  run_trace "http://localhost:$PORT/" "$RESULTS_DIR/trace-with-golit.json"
fi

cleanup

# --- Run WITHOUT golit -----------------------------------------------------
echo ">>> Starting container WITHOUT golit SSR..."
podman run -d --name "$CONTAINER_NAME" -p "$PORT:8080" \
  -e GOLIT_DISABLED=1 "$IMAGE" >/dev/null
wait_for_healthy

for endpoint in "${ENDPOINTS[@]}"; do
  tag=$(echo "$endpoint" | tr '/' '_')
  [[ "$tag" == "_" ]] && tag="_root"
  run_curl_bench "http://localhost:$PORT$endpoint" "$RESULTS_DIR/without${tag}.csv"
done

if $BROWSER && [[ -x "$CHROME" ]]; then
  for endpoint in "${ENDPOINTS[@]}"; do
    tag=$(echo "$endpoint" | tr '/' '_')
    [[ "$tag" == "_" ]] && tag="_root"
    run_browser_bench "http://localhost:$PORT$endpoint" "$RESULTS_DIR/without${tag}_browser.json"
  done
fi

if $TRACE && [[ -x "$CHROME" ]]; then
  run_trace "http://localhost:$PORT/" "$RESULTS_DIR/trace-without-golit.json"
fi

cleanup

# --- Print results ---------------------------------------------------------
echo ""
echo "============================================="
echo "  RESULTS"
echo "============================================="

for endpoint in "${ENDPOINTS[@]}"; do
  tag=$(echo "$endpoint" | tr '/' '_')
  [[ "$tag" == "_" ]] && tag="_root"

  with_file="$RESULTS_DIR/with${tag}.csv"
  without_file="$RESULTS_DIR/without${tag}.csv"

  echo ""
  echo "  Endpoint: $endpoint"
  echo "  -------------------------------------------"

  with_size=$(avg_bytes "$with_file")
  without_size=$(avg_bytes "$without_file")
  size_delta=$((with_size - without_size))
  echo "  Response size: ${with_size}B (with golit) vs ${without_size}B (without) [+${size_delta}B]"

  # TTFB (column 1)
  print_header "TTFB (ms)"
  with_ttfb_stats=$(compute_stats "$with_file" 1)
  without_ttfb_stats=$(compute_stats "$without_file" 1)
  print_row "  With golit" $with_ttfb_stats
  print_row "  Without golit" $without_ttfb_stats

  # Total response time (column 2)
  print_header "Total Response (ms)"
  with_total_stats=$(compute_stats "$with_file" 2)
  without_total_stats=$(compute_stats "$without_file" 2)
  print_row "  With golit" $with_total_stats
  print_row "  Without golit" $without_total_stats

  # Overhead
  with_ttfb_avg=$(echo "$with_ttfb_stats" | awk '{print $2}')
  without_ttfb_avg=$(echo "$without_ttfb_stats" | awk '{print $2}')
  with_total_avg=$(echo "$with_total_stats" | awk '{print $2}')
  without_total_avg=$(echo "$without_total_stats" | awk '{print $2}')

  echo ""
  echo "  Overhead (golit SSR cost):"
  print_delta "TTFB" "$with_ttfb_avg" "$without_ttfb_avg"
  print_delta "Total" "$with_total_avg" "$without_total_avg"

  # Browser metrics
  if $BROWSER; then
    with_browser="$RESULTS_DIR/with${tag}_browser.json"
    without_browser="$RESULTS_DIR/without${tag}_browser.json"
    if [[ -s "$with_browser" && -s "$without_browser" ]]; then
      echo ""
      echo "  Client-side metrics (avg over 10 runs):"
      printf "  %-24s %10s %10s\n" "" "with" "without"
      printf "  %-24s %10s %10s\n" "$(printf '%.0s-' {1..24})" "--------" "--------"
      for field in fcp lcp domContentLoaded loadEvent domInteractive; do
        w=$(browser_avg "$with_browser" "$field")
        wo=$(browser_avg "$without_browser" "$field")
        printf "  %-24s %9sms %9sms\n" "$field" "$w" "$wo"
      done
    fi
  fi
done

if $TRACE; then
  echo ""
  echo "  -------------------------------------------"
  echo "  Trace files saved to:"
  echo "    $RESULTS_DIR/trace-with-golit.json"
  echo "    $RESULTS_DIR/trace-without-golit.json"
  echo ""
  echo "  Open chrome://tracing and load each file"
  echo "  to compare flame graphs side by side."
fi

echo ""
echo "  Raw data saved to: $RESULTS_DIR/"
echo "============================================="
