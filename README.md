# golit

**Lit SSR in pure Go.** Server-side render Lit web components into [Declarative Shadow DOM](https://developer.chrome.com/docs/css-ui/declarative-shadow-dom) HTML -- no Node.js at build time.

golit executes your actual Lit component code using an embedded JavaScript engine (QuickJS via WebAssembly) to produce Declarative Shadow DOM HTML. classMap, directives, private fields, reactive state -- everything works because the real component `render()` method runs. Built entirely in Go with zero CGo dependencies.

## Why

Lit web components render empty shells until JavaScript loads. The official `@lit-labs/ssr` requires Node.js, which adds complexity and slows down Go-based SSG build pipelines. golit brings full-fidelity Lit SSR into a single Go binary.

**Before golit:** Users see a blank component until JS downloads, parses, and executes.

**After golit:** Users see fully rendered content immediately. Lit hydrates in the background for interactivity.

## Quick Start

### Prerequisites

- Go 1.22+
- Component source files available **locally on disk** (installed via npm)

golit bundles and executes component source code at build time using esbuild. This requires the component library and all of its dependencies to be installed locally -- typically via `npm install` in your project.

```bash
npm install @rhds/elements   # or whatever component library you use
```

**CDN import maps are not supported for SSR.** If your HTML uses a CDN-based import map (e.g. `https://ga.jspm.io/...`), golit cannot fetch and bundle from remote URLs. You must either:
- Install the packages locally and use an import map with local paths for golit
- Use `--sources` to point at the local `node_modules` directory
- Use `--importmap` with a build-time import map that resolves to local files

The browser still uses your CDN import map at runtime for hydration and interactivity -- golit only needs local files during the build.

### Install

```bash
go install github.com/zeroedin/golit/cmd/golit@latest
```

### The simplest workflow (zero config)

Write your HTML with custom elements and an import map that points to local files:

```html
<!DOCTYPE html>
<html>
<head>
  <script type="importmap">{
    "imports": {
      "@rhds/elements/": "./node_modules/@rhds/elements/elements/"
    }
  }</script>
  <script type="module">
    import '@rhds/elements/rh-badge/rh-badge.js';
  </script>
</head>
<body>
  <rh-badge state="success" number="7">7</rh-badge>
</body>
</html>
```

Then run golit:

```bash
golit transform public/
```

That's it. golit reads the import map and module imports from your HTML, resolves them to local files, bundles and executes the components, and injects Declarative Shadow DOM. No configuration files, no separate bundle step.

> **Note:** The import map paths must resolve to files on disk. If your production HTML uses CDN URLs in the import map, provide a separate build-time import map via `--importmap` that points to local `node_modules`.

## Component Discovery

golit supports four ways to discover which components to SSR. All modes are combinable.

### Mode 1: HTML Auto-Discovery (default, zero config)

golit reads `<script type="importmap">` and `<script type="module">` directly from the HTML being transformed. Components are bundled on-demand.

```bash
golit transform public/
```

### Mode 2: CLI Import Map

Pass an import map file explicitly. golit reads module imports from the HTML and resolves them through your import map.

```bash
golit transform public/ --importmap importmap.json
```

### Mode 3: Source Directory

Point at a directory of component source files. golit bundles each one and makes them available for rendering.

```bash
golit transform public/ --sources node_modules/@rhds/elements/elements/
```

### Mode 4: Pre-Bundled

For CI/CD or when you want maximum transform speed, pre-bundle components and point at the output directory. Bundling a directory automatically discovers shared dependencies, produces a shared runtime module, and thin per-component ES modules.

```bash
golit bundle node_modules/@rhds/elements/elements/ --out bundles/
golit transform public/ --defs bundles/
```

### Combining Modes

All flags are optional and combinable:

```bash
golit transform public/ --defs bundles/ --importmap importmap.json --sources extra/ --out dist/
```

## Using with Hugo

No special Hugo module or shortcode is needed. Write your HTML templates with custom elements as normal, build with Hugo, then post-process:

```bash
hugo build && golit transform public/
```

If your import map uses paths relative to the site root, golit resolves them automatically. For more control, pass an import map via CLI:

```bash
hugo build && golit transform public/ --importmap importmap.json
```

### Try the Hugo example

A complete working example using [Red Hat Design System](https://ux.redhat.com) components is included in `examples/hugo-rhds/`. To try it:

```bash
cd examples/hugo-rhds
npm install
make serve
```

This builds golit from source, runs Hugo, transforms all custom elements into Declarative Shadow DOM, and serves the result at `http://localhost:8080`. Open it in a browser to see fully server-side rendered web components -- no client JS needed for first paint.

For authoring content without SSR (faster iteration with Hugo's live-reload dev server):

```bash
make serve-dev
```

See `examples/hugo-rhds/Makefile` for the full build pipeline.

## Middleware examples (PHP and Ruby)

These examples show how to SSR Lit in a **dynamic app** (front controller or Rack) instead of batch-transforming static files.

| Example | Directory | Default URL |
| -------- | ---------- | ------------ |
| PHP (built-in server) | [`examples/php-middleware/`](examples/php-middleware/) | `http://localhost:8080` |
| Ruby (Rack) | [`examples/ruby-middleware/`](examples/ruby-middleware/) | `http://localhost:9292` |

Each demo includes a small Lit component (`<my-counter>`), `golit bundle` output under `bundles/`, and middleware that sends HTML through golit before the response is returned.

### Warm path: `golit serve`

By default, **containers** and **`make serve`** start a long-lived **`golit serve`** process that keeps one [`Renderer`](https://pkg.go.dev/github.com/zeroedin/golit#Renderer) warm. The app posts each full HTML document to `POST /render` (see environment variable **`GOLIT_SERVE_URL`**). That avoids spawning **`golit transform`** on every request.

- **`GET /health`** -- readiness probe
- **`POST /render`** -- body is full HTML; response is transformed HTML

CLI usage (same binary as `transform` / `bundle`):

```bash
golit serve --defs bundles/ --listen 127.0.0.1:9777
```

Flags: **`--defs`** (or **`GOLIT_DEFS`**), **`--listen`** (or **`GOLIT_SERVE_LISTEN`**; default `127.0.0.1:9777`), **`--sources`**, repeatable **`--ignore`**.

### Cold path: `golit transform` per request

If **`GOLIT_SERVE_URL`** is **unset**, the examples fall back to running **`golit transform`** on a temporary directory for each HTML response (simpler deployment, higher latency).

### Environment variables

| Variable | Purpose |
| -------- | -------- |
| **`GOLIT_SERVE_URL`** | Base URL of `golit serve` (e.g. `http://127.0.0.1:9777`). When set, middleware uses HTTP instead of exec. |
| **`GOLIT_DEFS`** | Directory of `.golit.module.js` files and `_runtime.golit.module.js` (used by `golit serve` and cold-path transform). |
| **`GOLIT_BIN`** | Path to `golit` binary for the cold path only. |
| **`GOLIT_DISABLED`** | If set (e.g. `1`), skip SSR and serve untransformed HTML (used by benchmarks for A/B comparison). |

QuickJS SSR (used by **`golit transform`** and **`golit serve`** rendering) exposes **`globalThis.fetch`** backed by Go **`net/http`** (similar in spirit to **`@lit/ssr`** + **`node-fetch`**). **`matchMedia`** is not defined server-side (viewport-only).

| Variable | Purpose |
| -------- | -------- |
| **`GOLIT_SSR_LOCATION`** | Base URL string for **`globalThis.location`** (default **`http://localhost/`**), aligned with Lit’s **`getWindow()`** default. |
| **`GOLIT_FETCH_ALLOWLIST`** | Optional comma-separated **hostnames** (no scheme/port). If set, **`fetch`** may only request those hosts (mitigates SSRF). If unset, only **`http:`** / **`https:`** are allowed. |
| **`GOLIT_FETCH_TIMEOUT_SEC`** | Per-request timeout in seconds (default **10**, clamped). |
| **`GOLIT_FETCH_MAX_BODY_BYTES`** | Max response body bytes read (default **16 MiB**, capped). |

After upgrading golit, regenerate pre-bundled modules (e.g. **`make bundle`** in **`examples/hugo-rhds`**) so the shared runtime and domshim match the new binary.

### PHP example

**Prerequisites:** Go (to build golit from this repo), PHP 8+, Node/npm for `npm install` during `make build`.

```bash
cd examples/php-middleware
npm install
make serve          # http://localhost:8080 — starts golit serve + PHP
```

**Container** (build from **repository root**; uses Podman in the Makefile — use `docker build` / `docker run` if you prefer):

```bash
cd examples/php-middleware
make container      # image: golit-php
make container-run  # publishes port 8080
```

**Without SSR** (static `public/` only, for comparison): `make serve-raw`.

### Ruby example

**Prerequisites:** Go, Ruby 3+, Bundler, Node/npm.

```bash
cd examples/ruby-middleware
npm install
bundle install
make serve            # http://localhost:9292 — golit serve + rackup
```

**Container:**

```bash
cd examples/ruby-middleware
make container        # image: golit-ruby
make container-run    # publishes port 9292
```

### Benchmarks (PHP and Ruby)

From each example directory, scripts compare **SSR on vs off** using **curl** timings (TTFB, total time, response size), **startup** time from **`container run`** until a static asset returns **200** (so HTML routes stay cold), **cold first HTML** request per endpoint right after that probe, plus **container memory** from **`stats`** on the host after the HTTP load. Optional tiers use **Chrome** for client metrics and traces (no extra load-test binaries).

**Requirements:** **`curl`**, **`make`**, and **`podman`** or **`docker`** on `PATH`. **`bench.sh`** picks **one** OCI binary (**`podman` if present, otherwise `docker`**) and uses it for **`run`**, **`rm`**, **`stats`**, and **`make container CONTAINER_RUNTIME=…`** (so the image build matches the runtime that runs the bench). **`python3`** is used for millisecond startup timing when available (otherwise startup is second-rounded). For **`make bench-full`** / **`bench-trace`**: Google Chrome at the default macOS path (headless). For **`make container`** / **`container-run`** when **not** using **`bench.sh`**, the default is **`podman`**; pass **`CONTAINER_RUNTIME=docker`** to use **Docker**.

```bash
cd examples/php-middleware   # or ruby-middleware
make bench          # 100 requests per endpoint (/ and /about), tier 1 only
make bench-quick    # 20 requests
make bench-full     # adds --browser and --trace (Chrome)
make bench-trace    # Chrome trace files for chrome://tracing
```

Direct script flags:

```bash
./bench.sh -n 50              # custom request count
./bench.sh --browser          # Performance API metrics via headless Chrome
./bench.sh --trace            # CPU trace JSON for flame charts
```

Results, raw CSVs, memory snapshots (`mem_with.snapshot`, `mem_without.snapshot`), startup files (`startup_with_ms.txt`, `startup_without_ms.txt`), and per-endpoint cold-request lines (`with_first_*.csv`, `without_first_*.csv`) are written to **`bench-results/`** (gitignored). The script builds the image with **`make container`** using the **same** detected **`podman`**/**`docker`** as **`run`**, then runs **with** golit (warm `golit serve` in the entrypoint), then **without** (`GOLIT_DISABLED=1`), and prints a side-by-side summary.

**Reading the numbers:** **Startup** is **instance readiness** (from **`container run`** until **`GET /components/my-counter.js`** returns 200). That window includes container boot and the entrypoint; the **with − without** delta is mostly **starting `golit serve`**, paid **once per new container instance**, not on every HTTP request. **Cold first HTML** is the first **`GET`** to each benchmarked path **after** that static probe, so it measures one cold trip through the SSR path without warming `/` or `/about` during the health wait. **Steady-state** timings are the bulk **`curl`** runs after that. In production you can **prewarm** (e.g. readiness checks or startup requests that hit real SSR URLs) to move cold cost into deploy/scale-up instead of the first user.

## CLI Reference

### `golit transform`

Post-process HTML files, expanding custom elements into Declarative Shadow DOM.

```bash
golit transform <html-dir> [options]
```

Options:
- `--defs <dir>` -- Directory of pre-bundled `.golit.module.js` files (and `_runtime.golit.module.js`)
- `--sources <dir>` -- Directory of component `.js`/`.ts` source files (auto-bundles)
- `--importmap <file>` -- Import map JSON file for resolving bare-module specifiers
- `--out <dir>` -- Output to a separate directory (default: in-place)
- `--verbose` -- Print progress to stderr
- `--dry-run` -- Process without writing
- `-j [N]` / `--concurrency [N]` -- Process files in parallel. `-j` alone uses all available CPUs; `-j 4` uses 4 workers. Default is sequential.

When no discovery flags are provided, auto-discovery from HTML is used.

### `golit bundle`

Pre-bundle Lit components for SSR. When given a directory, automatically discovers shared dependencies via esbuild Metafile analysis, produces a shared runtime module (`_runtime.golit.module.js`) containing all shared dependencies, plus thin per-component `.golit.module.js` files that import from the shared runtime.

```bash
golit bundle <src-dir/> [--out <bundles-dir/>] [--minify]
golit bundle <source.ts|js> [--out <file>] [--minify]
```

The three-pass build discovers dependencies from the actual import graph (no hardcoded package lists). The shared runtime is loaded once per QJS engine instance. Each component module contains only the component's own code and imports, avoiding duplicate classes and decorator state across components.

### `golit render`

Render a single HTML fragment to stdout. Useful for testing and scripting.
Requires pre-built modules (see `golit bundle` above).

```bash
# First, bundle the component(s) you want to render
golit bundle node_modules/@rhds/elements/elements/ --out bundles/

# Then render a fragment using the pre-built modules
golit render --defs bundles/ '<rh-badge state="success" number="7">7</rh-badge>'
```

### `golit serve`

Run an HTTP server that holds a warm **`Renderer`** and transforms full HTML documents on each request. Intended for middleware integration (PHP, Ruby, etc.); avoids per-request **`golit transform`** process startup.

```bash
golit serve --defs bundles/ [--listen host:port]
```

- **`GET /health`** returns `200` and plain text `ok`.
- **`POST /render`** accepts a full HTML document as the body; response is transformed HTML (`text/html`).

See [Middleware examples](#middleware-examples-php-and-ruby) for how the PHP/Ruby demos wire this up.

### `golit version`

Print the version.

## Library Usage

golit can be used as a Go library. Import `github.com/zeroedin/golit` and use the `Renderer` type:

```go
package main

import (
	"fmt"
	"log"

	"github.com/zeroedin/golit"
)

func main() {
	renderer, err := golit.NewRenderer(golit.RendererOptions{
		DefsDir: "bundles/",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer renderer.Close()

	output, err := renderer.RenderFragment(`<my-el name="World"></my-el>`)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(output)
}
```

The `Renderer` exposes three rendering methods:

- `RenderFragment(html)` -- Render an HTML fragment
- `RenderHTML(html)` -- Render a full HTML document
- `TransformDir(dir)` -- Process all HTML files in a directory

For lower-level control, use the `pkg/jsengine` and `pkg/transformer` packages directly.

## How It Works

```
HTML with import map + module scripts
    |
    | golit transform
    | 1. Parse HTML, find <script type="importmap"> and <script type="module">
    | 2. Resolve bare-module specifiers to file paths via import map
    | 3. Discover shared dependencies, build thin ES modules + shared runtime
    | 4. Load shared runtime (@golit/runtime) into QJS once
    | 5. Load thin component modules (import from shared runtime)
    | 6. For each custom element in HTML:
    |    a. Instantiate component, set attributes, call render()
    |    b. Collect rendered HTML + CSS
    |    c. Wrap in <template shadowroot="open" shadowrootmode="open">
    v
HTML with Declarative Shadow DOM
    |
    | Browser
    v
Instant paint -> Lit hydrates -> Interactive
```

## Architecture

golit uses three key technologies:

- **esbuild** (Go-native) -- Three-pass build: (1) discovers shared dependencies via Metafile analysis, (2) produces thin per-component ES modules with shared deps as external imports, (3) bundles the shared runtime from the discovered dependency graph. Handles imports, decorators, private fields, and module resolution. Uses Node.js conditional exports (`"node"` condition) so Lit's `isServer` is `true`.
- **QJS** (QuickJS via WebAssembly/Wazero) -- Loads the shared runtime module once via `JS_SetModuleLoaderFunc`, then evaluates thin component modules that import from it. Pure Go, no CGo, cross-compiles everywhere. ~2MB WASM module, ~400ms cold start, <1ms per render.
- **golang.org/x/net/html** -- Parses and transforms HTML documents, inserting Declarative Shadow DOM templates.

### Output Format

golit produces Lit-compatible Declarative Shadow DOM with hydration markers:

```html
<rh-badge state="success" number="7">
  <template shadowroot="open" shadowrootmode="open">
    <style>/* component CSS */</style>
    <!--lit-part hqKOHqwbWgk=-->
      <!--lit-node 0--><span class="success">
        <!--lit-part-->7<!--/lit-part-->
      </span>
      <!--lit-node 2--><slot class="success"></slot>
    <!--/lit-part-->
  </template>
  7
</rh-badge>
```

- `<!--lit-part DIGEST-->` -- Template boundary with DJB2 digest for hydration verification
- `<!--lit-node N-->` -- Marks elements with attribute bindings for hydration
- `<!--lit-part-->value<!--/lit-part-->` -- Child expression value boundaries
- `defer-hydration` -- Added to nested custom elements inside shadow roots

### Package Structure

```
cmd/golit/              CLI binary (bundle, compile, transform, render, serve, version)
pkg/jsengine/           QJS engine, esbuild bundler, DOM shim, template collector,
                        import map parser, module registry, shared runtime loader
pkg/transformer/        HTML file walker, component discovery, DSD expansion
```

## Dependencies

- `github.com/evanw/esbuild` -- TypeScript/JavaScript bundler (Go-native)
- `github.com/fastschema/qjs` -- QuickJS via WebAssembly (pure Go, no CGo)
- `golang.org/x/net/html` -- HTML5 parser

## Contributing

### Changesets

This project uses [Changesets](https://github.com/changesets/changesets) for version management and changelog generation. When making a notable change (new feature, bug fix, breaking change), add a changeset before merging:

```bash
npx changeset
```

This prompts you to select the semver bump type (major/minor/patch) and write a short description. A markdown file is created in `.changeset/` and committed with your PR.

Changes that don't warrant a release (docs, CI tweaks, refactoring with no public API impact) can skip this step.

### Releasing

Releases are fully automated via GitHub Actions:

1. When PRs with changeset files are merged to `main`, the [changesets action](https://github.com/changesets/action) opens (or updates) a "chore: prepare release" PR that bumps the version in `package.json` and updates `CHANGELOG.md`.
2. When you merge the release PR, the workflow automatically:
   - Creates a git tag (`vX.Y.Z`)
   - Cross-compiles binaries for Linux, macOS, and Windows (amd64 + arm64)
   - Publishes a GitHub Release with all binaries and SHA-256 checksums

### Building locally

```bash
make build          # Build for current platform
make test           # Run tests
make cross-compile  # Build for all platforms (output in dist/)
make help           # Show all targets
```

## License

MIT
