# bench/

Shared assets for golit example benchmarks.

## perf-snippet.js

A browser-side performance measurement snippet that collects Web Vitals
via `PerformanceObserver` and the Navigation Timing API. It only activates
when the page URL contains `?bench` — zero impact on normal usage.

After `window.onload` plus a 500ms settle, it writes a JSON object into a
hidden `<pre id="golit-perf">` element with these metrics:

| Field | Source |
|-------|--------|
| `fcp` | First Contentful Paint |
| `lcp` | Largest Contentful Paint |
| `ttfb` | Time to First Byte (`responseStart`) |
| `domContentLoaded` | `domContentLoadedEventEnd` |
| `domInteractive` | `domInteractive` |
| `loadEvent` | `loadEventEnd` |
| `responseEnd` | `responseEnd` |

### How it's used

The PHP and Ruby middleware examples include this snippet in their page
layouts:

- `examples/php-middleware/public/pages/_layout.php`
- `examples/ruby-middleware/views/layout.erb`

Each middleware example has a `bench.sh` script that:

1. Starts the app container with and without golit SSR
2. Runs server-side benchmarks via `curl` (TTFB, total response, cold start)
3. Optionally runs client-side benchmarks via Chrome headless (`--browser` flag)

For step 3, Chrome loads each page with `?bench` appended to the URL and
uses `--dump-dom` to extract the JSON from `<pre id="golit-perf">`.

### Running benchmarks

From a middleware example directory:

```sh
# Server-side only (100 requests per endpoint)
./bench.sh

# Fewer requests
./bench.sh -n 50

# Include client-side Web Vitals (requires Chrome)
./bench.sh --browser

# Capture Chrome trace files for flame graphs
./bench.sh --trace
```

See `examples/php-middleware/bench.sh` or `examples/ruby-middleware/bench.sh`
for full usage details.
