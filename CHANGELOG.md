## v0.1.0 (2026-07-20)

- Initial minor release consolidating all features since first commit.

  **Features**
  - Library API with batch rendering, error reporting, and CLI enhancements (#1)
  - ES module loader architecture with shared runtime (#9)
  - Engine pool for concurrent request handling in `golit serve` (#12)
  - SSR fetch bridge and Lit-aligned domshim globals (#8)
  - Middleware examples, warm path serving, and SSR benchmarks (#7)
  - Server-Timing header with render duration and pool utilization (#21)
  - Reading HTML fragments from stdin (#23)
  - Stdio mode (`--stdio`) for NUL-delimited stdin/stdout protocol (#24)
  - Shared render cache across pool engines with 4096-entry cap (#28, #29, #31)
  - Configurable esbuild options: `--target`, `--format`, `--platform`, `--conditions`, `--main-fields` (#34)
  - Warn on unregistered elements and render failures with `--quiet` suppression (#35)
  - Filename-based tag discovery for bundles with inlined dependencies (#34)

  **Fixes**
  - Guard domshim globals to prevent CustomElementRegistry wipe (#6)
  - Stop injecting defer-hydration attribute during SSR (#3)
  - Non-fatal discovery errors and dynamic import resolution (#13)
  - Emit root lit-part markers for components with empty render (#17)
  - Register dynamic import targets as standalone QJS modules (#20)
  - Preserve DOCTYPE/html/head/body in full document input (#27)
  - Walk all node_modules directories for nested dependency resolution (#34)
  - Add MainFields and browser conditions for PlatformNeutral esbuild resolution (#34)

  **Performance**
  - Transform pipeline optimizations (#2)
  - Render result caching and JS execution optimizations (#28)

# golit

## Changelog
