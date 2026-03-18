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
go install github.com/sspriggs/golit/cmd/golit@latest
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

For CI/CD or when you want maximum transform speed, pre-bundle components and point at the bundles directory.

```bash
golit bundle node_modules/@rhds/elements/elements/rh-badge/rh-badge.js --out bundles/
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

## CLI Reference

### `golit transform`

Post-process HTML files, expanding custom elements into Declarative Shadow DOM.

```bash
golit transform <html-dir> [options]
```

Options:
- `--defs <dir>` -- Directory of pre-bundled `.golit.bundle.js` files
- `--sources <dir>` -- Directory of component `.js`/`.ts` source files (auto-bundles)
- `--importmap <file>` -- Import map JSON file for resolving bare-module specifiers
- `--out <dir>` -- Output to a separate directory (default: in-place)
- `--verbose` -- Print progress to stderr
- `--dry-run` -- Process without writing

When no discovery flags are provided, auto-discovery from HTML is used.

### `golit bundle`

Pre-bundle a Lit component for SSR. Produces a `.golit.bundle.js` file.

```bash
golit bundle <source.ts|js> [--out <file>] [--minify]
golit bundle <src-dir/> [--out <bundles-dir/>] [--minify]
```

The bundle includes the component, Lit runtime, DOM shim, and template collector -- everything needed for the QJS engine to render the component.

### `golit render`

Render a single HTML fragment to stdout. Useful for testing and scripting.
Requires pre-built bundles (see `golit bundle` above).

```bash
# First, bundle the component(s) you want to render
golit bundle node_modules/@rhds/elements/elements/rh-badge/rh-badge.js --out bundles/

# Then render a fragment using the pre-built bundles
golit render --defs bundles/ '<rh-badge state="success" number="7">7</rh-badge>'
```

### `golit version`

Print the version.

## Library Usage

golit can be used as a Go library. Import `github.com/sspriggs/golit` and use the `Renderer` type:

```go
package main

import (
	"fmt"
	"log"

	"github.com/sspriggs/golit"
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
    | 3. Bundle each component with esbuild (component + Lit + DOM shim)
    | 4. For each custom element in HTML:
    |    a. Load bundle into QJS (QuickJS via WASM)
    |    b. Instantiate component, set attributes, call render()
    |    c. Collect rendered HTML + CSS
    |    d. Wrap in <template shadowroot="open" shadowrootmode="open">
    v
HTML with Declarative Shadow DOM
    |
    | Browser
    v
Instant paint -> Lit hydrates -> Interactive
```

## Architecture

golit uses three key technologies:

- **esbuild** (Go-native) -- Bundles TypeScript/JavaScript components with all dependencies into a single file. Handles imports, decorators, private fields, and module resolution. Uses Node.js conditional exports (`"node"` condition) so Lit's `isServer` is `true`.
- **QJS** (QuickJS via WebAssembly/Wazero) -- Executes the bundled component code in a sandboxed JavaScript environment. Pure Go, no CGo, cross-compiles everywhere. ~2MB WASM module, ~400ms cold start, <1ms per render.
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
cmd/golit/              CLI binary
pkg/jsengine/           QJS engine, esbuild bundler, DOM shim, template collector,
                        import map parser, bundle registry
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
