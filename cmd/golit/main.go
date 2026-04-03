// golit is a CLI tool for server-side rendering Lit web components
// into Declarative Shadow DOM HTML using a QJS JavaScript engine.
//
// Usage:
//
//	golit bundle <source.ts|js> [--out <file>] [--minify]
//	golit transform <html-dir> --defs <bundles-dir> [--out <dir>] [--verbose]
//	golit render --defs <bundles-dir> '<html-fragment>'
//	golit version
package main

import (
	"fmt"
	"os"
)

// version is overridden at build time via ldflags -X main.version=...
// The Makefile reads the canonical version from package.json (managed by changesets).
// "dev" is the fallback for uninjected builds, e.g. `go run ./cmd/golit`.
var version = "dev"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "bundle":
		if err := runBundle(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "compile":
		if err := runCompile(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "transform":
		if err := runTransform(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "render":
		if err := runRender(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "version":
		fmt.Printf("golit %s\n", version)
	case "help", "--help", "-h":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `golit - Lit SSR in pure Go (QJS engine)

Usage:
  golit bundle <source.ts|js> [--out <file.golit.bundle.js>] [--minify]
  golit bundle <src-dir/> [--out <bundles-dir/>] [--minify]
  golit compile --defs <bundles-dir> [--out <file.golit.compiled.js>] [--minify]
  golit transform <html-dir> [--defs <dir>] [--compiled <file>] [--sources <dir>] [--importmap <file>] [--out <dir>]
  golit render --defs <bundles-dir> '<html-fragment>'
  golit render --component-js '<source>' '<html-fragment>'
  golit version

Commands:
  bundle      Bundle a Lit component source file into a .golit.bundle.js
              file for SSR rendering. Includes Lit, DOM shim, and all deps.
  compile     Combine all bundles from a --defs directory into a single
              .golit.compiled.js artifact with a tag registry manifest.
  transform   Post-process HTML files, expanding custom elements into
              Declarative Shadow DOM using bundled components.
  render      Render a single HTML fragment to stdout.
  version     Print version information.

Options:
  --out <path>       Output path for bundle/transform
  --defs <dir>       Directory containing pre-bundled .golit.bundle.js files
  --compiled <file>  Single pre-compiled .golit.compiled.js artifact
  --sources <dir>    Directory of component .js/.ts source files (auto-bundles)
  --importmap <file> Import map JSON file for resolving bare-module specifiers
  --ignore <tag>     Skip SSR for this custom element (repeatable)
  --minify           Minify the output bundle
  --verbose          Print progress to stderr
  --dry-run          Process files without writing changes
  --strict           Exit with error if any components fail to render
  --concurrency, -j  Parallel workers for transform (default: sequential)
                     -j alone uses all CPUs, -j N uses N workers
  --component-js     Inline JS/TS component source for render command (repeatable)

Component Discovery (transform command):
  golit discovers which components to SSR using four modes (combinable):
  1. --defs      Pre-bundled .golit.bundle.js files
  2. --sources   Directory of source files (bundled on-demand)
  3. --importmap Import map file + module scripts in HTML
  4. Auto        Parse <script type="importmap"> and <script type="module">
                 from the HTML itself (default when no flags given)

Examples:
  # Auto-discover from HTML (zero config)
  hugo build && golit transform public/

  # Use pre-bundled components
  golit bundle src/components/ --out bundles/
  golit transform public/ --defs bundles/

  # Point at installed package sources
  golit transform public/ --sources node_modules/@rhds/elements/elements/

  # Use an import map
  golit transform public/ --importmap importmap.json

  # Render a single element
  golit render --defs bundles/ '<my-greeting name="World"></my-greeting>'
`)
}
