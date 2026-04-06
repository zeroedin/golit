#!/usr/bin/env bash
set -euo pipefail

# ---------------------------------------------------------------------------
# PHP + golit performance benchmark
#
# Measures server-side (curl), startup time (run → static asset 200), cold
# first HTML request per endpoint, container memory (podman/docker stats), and
# optionally client-side (Chrome headless) metrics, then prints a comparison.
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
# Ready probe must not hit HTML routes (keeps / and /about cold for first-hit).
HEALTH_PATH="/components/my-counter.js"
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

now_ms() {
  if command -v python3 &>/dev/null; then
    python3 -c 'import time; print(int(time.time() * 1000))'
  else
    echo $(( $(date +%s) * 1000 ))
  fi
}

wait_for_healthy() {
  local url="http://localhost:$PORT${HEALTH_PATH}"
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

container_engine() {
  if command -v podman &>/dev/null; then
    printf '%s\n' podman
  elif command -v docker &>/dev/null; then
    printf '%s\n' docker
  else
    printf '%s\n' ""
  fi
}

# Maximum Mem%% numeric value over several quick samples (host sees cgroup limit).
mem_peak_percent() {
  local name="$1" eng="$2"
  local max=0 k p
  for k in 1 2 3 4 5; do
    p=$("$eng" stats "$name" --no-stream --format '{{.MemPerc}}' 2>/dev/null || true)
    p=$(echo "$p" | tr -dc '0-9.')
    [[ -z "$p" ]] && p=0
    max=$(awk -v a="$max" -v b="$p" 'BEGIN { if ((b + 0) > (a + 0)) print b + 0; else print a + 0 }')
    sleep 0.12
  done
  awk -v m="$max" 'BEGIN { printf "%.2f", m }'
}

write_mem_snapshot() {
  local name="$1" outfile="$2"
  local eng usage pct peak
  eng=$(container_engine)
  {
    echo "engine=${eng:-none}"
    if [[ -z "$eng" ]]; then
      echo "usage=n/a"
      echo "pct=n/a"
      echo "peak_pct=n/a"
      return
    fi
    usage=$("$eng" stats "$name" --no-stream --format '{{.MemUsage}}' 2>/dev/null || echo "n/a")
    pct=$("$eng" stats "$name" --no-stream --format '{{.MemPerc}}' 2>/dev/null || echo "n/a")
    peak=$(mem_peak_percent "$name" "$eng")
    echo "usage=$usage"
    echo "pct=$pct"
    echo "peak_pct=$peak"
  } > "$outfile"
}

print_container_mem_section() {
  local wf="$RESULTS_DIR/mem_with.snapshot"
  local uf="$RESULTS_DIR/mem_without.snapshot"
  [[ -f "$wf" && -f "$uf" ]] || return
  local eng u_w p_w pk_w u_u p_u pk_u
  eng=$(grep '^engine=' "$wf" | cut -d= -f2-)
  u_w=$(grep '^usage=' "$wf" | cut -d= -f2-)
  p_w=$(grep '^pct=' "$wf" | cut -d= -f2-)
  pk_w=$(grep '^peak_pct=' "$wf" | cut -d= -f2-)
  u_u=$(grep '^usage=' "$uf" | cut -d= -f2-)
  p_u=$(grep '^pct=' "$uf" | cut -d= -f2-)
  pk_u=$(grep '^peak_pct=' "$uf" | cut -d= -f2-)
  echo ""
  echo "  Container memory (host: ${eng:-n/a}; after HTTP benchmark load)"
  echo "  -------------------------------------------"
  echo "  With golit:"
  echo "    usage / limit:  $u_w"
  echo "    Mem %:          $p_w    peak ~${pk_w}% (max of 5 quick samples)"
  echo "  Without golit:"
  echo "    usage / limit:  $u_u"
  echo "    Mem %:          $p_u    peak ~${pk_u}% (max of 5 quick samples)"
  if [[ "$pk_w" != "n/a" && "$pk_u" != "n/a" ]]; then
    local d
    d=$(awk -v a="$pk_w" -v b="$pk_u" 'BEGIN { printf "%+.2f", a - b }')
    echo ""
    echo "  Peak Mem % delta (with − without): ${d}%"
  fi
}

# Single request; same columns as run_curl_bench (cold HTML after health probe).
record_first_hit() {
  local url="$1"
  local outfile="$2"
  curl -sf -o /dev/null \
    -w '%{time_starttransfer} %{time_total} %{size_download} %{http_code}\n' \
    "$url" > "$outfile"
}

first_hit_ttfb_total_ms() {
  local file="$1"
  [[ -s "$file" ]] || { echo "- -"; return; }
  awk '{ printf "%.2f %.2f", $1 * 1000, $2 * 1000 }' "$file"
}

print_startup_section() {
  local wf="$RESULTS_DIR/startup_with_ms.txt"
  local uf="$RESULTS_DIR/startup_without_ms.txt"
  [[ -f "$wf" && -f "$uf" ]] || return
  local sw so
  sw=$(cut -d= -f2 "$wf")
  so=$(cut -d= -f2 "$uf")
  echo ""
  echo "  Startup (podman run → first 200 on ${HEALTH_PATH})"
  echo "  -------------------------------------------"
  printf "  %-18s %s ms\n" "With golit:" "$sw"
  printf "  %-18s %s ms\n" "Without golit:" "$so"
  if [[ "$sw" =~ ^[0-9]+$ && "$so" =~ ^[0-9]+$ ]]; then
    local d
    d=$((sw - so))
    printf "  %-18s %+d ms (with − without)\n" "Delta:" "$d"
  fi
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
echo "  Container memory: podman or docker stats on host"
echo "  Startup + first HTML: static probe then cold page timings"
echo "============================================="

# Build the container image
echo ""
echo "Building container image..."
make -C "$SCRIPT_DIR" container 2>&1 | tail -1

# --- Run WITH golit --------------------------------------------------------
echo ""
echo ">>> Starting container WITH golit SSR..."
cleanup
_start=$(now_ms)
podman run -d --name "$CONTAINER_NAME" -p "$PORT:8080" "$IMAGE" >/dev/null
wait_for_healthy
_end=$(now_ms)
echo "ms=$((_end - _start))" > "$RESULTS_DIR/startup_with_ms.txt"

for endpoint in "${ENDPOINTS[@]}"; do
  tag=$(echo "$endpoint" | tr '/' '_')
  [[ "$tag" == "_" ]] && tag="_root"
  record_first_hit "http://localhost:$PORT$endpoint" "$RESULTS_DIR/with_first${tag}.csv"
  run_curl_bench "http://localhost:$PORT$endpoint" "$RESULTS_DIR/with${tag}.csv"
done

sleep 0.25
write_mem_snapshot "$CONTAINER_NAME" "$RESULTS_DIR/mem_with.snapshot"

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
_start=$(now_ms)
podman run -d --name "$CONTAINER_NAME" -p "$PORT:8080" \
  -e GOLIT_DISABLED=1 "$IMAGE" >/dev/null
wait_for_healthy
_end=$(now_ms)
echo "ms=$((_end - _start))" > "$RESULTS_DIR/startup_without_ms.txt"

for endpoint in "${ENDPOINTS[@]}"; do
  tag=$(echo "$endpoint" | tr '/' '_')
  [[ "$tag" == "_" ]] && tag="_root"
  record_first_hit "http://localhost:$PORT$endpoint" "$RESULTS_DIR/without_first${tag}.csv"
  run_curl_bench "http://localhost:$PORT$endpoint" "$RESULTS_DIR/without${tag}.csv"
done

sleep 0.25
write_mem_snapshot "$CONTAINER_NAME" "$RESULTS_DIR/mem_without.snapshot"

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

print_startup_section

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

  first_with="$RESULTS_DIR/with_first${tag}.csv"
  first_without="$RESULTS_DIR/without_first${tag}.csv"
  read -r cold_w_ttfb cold_w_tot <<< "$(first_hit_ttfb_total_ms "$first_with")"
  read -r cold_wo_ttfb cold_wo_tot <<< "$(first_hit_ttfb_total_ms "$first_without")"
  echo ""
  echo "  First HTML request (cold, after ${HEALTH_PATH} probe)"
  printf "  %-18s %10s %10s\n" "" "TTFB" "Total"
  printf "  %-18s %9sms %9sms\n" "With golit:" "$cold_w_ttfb" "$cold_w_tot"
  printf "  %-18s %9sms %9sms\n" "Without golit:" "$cold_wo_ttfb" "$cold_wo_tot"
  if [[ "$cold_w_ttfb" != "-" && "$cold_wo_ttfb" != "-" ]]; then
    print_delta "Cold TTFB" "$cold_w_ttfb" "$cold_wo_ttfb"
  fi
  if [[ "$cold_w_tot" != "-" && "$cold_wo_tot" != "-" ]]; then
    print_delta "Cold total" "$cold_w_tot" "$cold_wo_tot"
  fi

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

print_container_mem_section

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
