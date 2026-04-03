// Package jsengine provides a QJS-based JavaScript execution engine
// for server-side rendering Lit components with full expression fidelity.
package jsengine

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/sspriggs/golit/pkg/fileutil"
)

//go:embed domshim.js
var domShimJS string

//go:embed templatecollector.js
var templateCollectorJS string

var (
	shimOnce      sync.Once
	shimDir       string
	shimPath      string
	collectorPath string
	shimErr       error
)

// ensureShimDir writes the embedded DOM shim and template collector to a
// temp directory once per process. All esbuild invocations share these
// files via esbuild's Inject option.
func ensureShimDir() (string, string, string, error) {
	shimOnce.Do(func() {
		shimDir, shimErr = os.MkdirTemp("", "golit-shim-*")
		if shimErr != nil {
			return
		}
		shimPath = filepath.Join(shimDir, "domshim.js")
		if shimErr = os.WriteFile(shimPath, []byte(domShimJS), 0644); shimErr != nil {
			return
		}
		collectorPath = filepath.Join(shimDir, "templatecollector.js")
		shimErr = os.WriteFile(collectorPath, []byte(templateCollectorJS), 0644)
	})
	return shimPath, collectorPath, shimDir, shimErr
}

// BundleOptions configures the component bundling.
type BundleOptions struct {
	// Minify minifies the output bundle.
	Minify bool
}

// BundleComponent uses esbuild to bundle a Lit component source file
// along with its dependencies (lit, DOM shim) into a single JS string
// that can be executed in QJS.
func BundleComponent(componentPath string, opts ...BundleOptions) (string, error) {
	opt := BundleOptions{}
	if len(opts) > 0 {
		opt = opts[0]
	}
	return bundleComponent(componentPath, opt)
}

func bundleComponent(componentPath string, opt BundleOptions) (string, error) {
	esm, err := bundleComponentRaw(componentPath, opt)
	if err != nil {
		return "", err
	}

	code, hasTopLevelAwait := stripESMExports(esm)
	if hasTopLevelAwait {
		code = "(async () => {\n" + code + "\n})();\n"
	}

	return code, nil
}

// bundleComponentRaw runs esbuild and returns the raw ESM output before
// any post-processing (export stripping, async wrapping).
func bundleComponentRaw(componentPath string, opt BundleOptions) (string, error) {
	absPath, err := filepath.Abs(componentPath)
	if err != nil {
		return "", fmt.Errorf("resolving path: %w", err)
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return "", fmt.Errorf("component file not found: %s", absPath)
	}

	sp, cp, _, err := ensureShimDir()
	if err != nil {
		return "", fmt.Errorf("preparing shim files: %w", err)
	}

	nodeModulesDir := findNodeModules(absPath)
	sourceDir := filepath.Dir(absPath)

	result := api.Build(api.BuildOptions{
		EntryPoints:      []string{absPath},
		Bundle:           true,
		Format:           api.FormatESModule,
		Target:           api.ES2022,
		Platform:         api.PlatformNeutral,
		Inject:           []string{sp, cp},
		Write:            false,
		MinifyWhitespace: opt.Minify,
		MinifySyntax:     opt.Minify,
		NodePaths:        []string{nodeModulesDir},
		TsconfigRaw:      `{"compilerOptions":{"experimentalDecorators":true,"useDefineForClassFields":false}}`,
		Plugins:          buildPlugins(opt),
		Conditions:       []string{"node"},
		LogLevel:         api.LogLevelSilent,
		AbsWorkingDir:    sourceDir,
	})

	if len(result.Errors) > 0 {
		var msgs []string
		for _, e := range result.Errors {
			msgs = append(msgs, e.Text)
		}
		return "", fmt.Errorf("esbuild bundle errors: %s", strings.Join(msgs, "; "))
	}

	if len(result.OutputFiles) == 0 {
		return "", fmt.Errorf("esbuild produced no output")
	}

	return string(result.OutputFiles[0].Contents), nil
}

// stripESMExports removes export keywords from bundled ESM output so it
// can be evaluated in QJS script mode. Since esbuild has already bundled
// all dependencies, there are no import statements to handle.
// Also detects top-level await (await at column 0, not inside functions).
func stripESMExports(code string) (string, bool) {
	hasTopLevelAwait := false
	inExportBlock := false
	lines := strings.Split(code, "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Track multi-line "export { ... };" blocks
		if inExportBlock {
			lines[i] = ""
			if strings.Contains(trimmed, "};") || trimmed == "}" {
				inExportBlock = false
			}
			continue
		}

		// Single-line "export { Foo, Bar };"
		if (strings.HasPrefix(trimmed, "export {") || strings.HasPrefix(trimmed, "export{")) &&
			strings.Contains(trimmed, "}") {
			lines[i] = ""
			continue
		}
		// Multi-line "export {\n  Foo,\n  Bar\n};"
		if strings.HasPrefix(trimmed, "export {") || strings.HasPrefix(trimmed, "export{") {
			lines[i] = ""
			inExportBlock = true
			continue
		}
		// "export default" -> just the value
		if strings.HasPrefix(trimmed, "export default ") {
			lines[i] = strings.Replace(line, "export default ", "", 1)
			continue
		}
		// "export var/let/const/function/class" -> strip "export "
		if strings.HasPrefix(trimmed, "export var ") ||
			strings.HasPrefix(trimmed, "export let ") ||
			strings.HasPrefix(trimmed, "export const ") ||
			strings.HasPrefix(trimmed, "export function ") ||
			strings.HasPrefix(trimmed, "export class ") ||
			strings.HasPrefix(trimmed, "export async ") {
			lines[i] = strings.Replace(line, "export ", "", 1)
			continue
		}
		// Detect top-level await: line starts with "await " at shallow indentation
		// (esbuild emits top-level code at indent 0; method bodies are deeper)
		if !hasTopLevelAwait && strings.HasPrefix(trimmed, "await ") {
			indent := len(line) - len(strings.TrimLeft(line, " \t"))
			if indent <= 2 {
				hasTopLevelAwait = true
			}
		}
	}
	return strings.Join(lines, "\n"), hasTopLevelAwait
}

// BundleComponents bundles multiple component source files in a single
// esbuild invocation, sharing module resolution and parsed ASTs across
// all entry points. Returns a map from input path to bundle string.
// Files that fail to bundle are silently skipped.
func BundleComponents(componentPaths []string, opts ...BundleOptions) (map[string]string, error) {
	opt := BundleOptions{}
	if len(opts) > 0 {
		opt = opts[0]
	}

	if len(componentPaths) == 0 {
		return nil, nil
	}

	// Resolve all paths and filter out missing files
	type entry struct {
		absPath string
		key     string // indexed output path for mapping
	}
	entries := make([]entry, 0, len(componentPaths))
	for _, p := range componentPaths {
		absPath, err := filepath.Abs(p)
		if err != nil {
			continue
		}
		if _, err := os.Stat(absPath); os.IsNotExist(err) {
			continue
		}
		key := fmt.Sprintf("entry_%d", len(entries))
		entries = append(entries, entry{absPath: absPath, key: key})
	}

	if len(entries) == 0 {
		return nil, nil
	}

	sp, cp, sd, err := ensureShimDir()
	if err != nil {
		return nil, fmt.Errorf("preparing shim files: %w", err)
	}

	// Find node_modules from the first entry point
	nodeModulesDir := findNodeModules(entries[0].absPath)

	// Build entry points list and output-key-to-input-path map
	esbuildEntries := make([]api.EntryPoint, len(entries))
	keyToPath := make(map[string]string, len(entries))
	for i, e := range entries {
		esbuildEntries[i] = api.EntryPoint{
			InputPath:  e.absPath,
			OutputPath: e.key,
		}
		keyToPath[e.key+".js"] = e.absPath
	}

	// Single esbuild call for all entry points
	result := api.Build(api.BuildOptions{
		EntryPointsAdvanced: esbuildEntries,
		Bundle:              true,
		Format:              api.FormatESModule,
		Target:              api.ES2022,
		Platform:            api.PlatformNeutral,
		Inject:              []string{sp, cp},
		Write:               false,
		Outdir:              sd, // required for multiple entry points; not written (Write:false)
		MinifyWhitespace:    opt.Minify,
		MinifySyntax:        opt.Minify,
		NodePaths:           []string{nodeModulesDir},
		TsconfigRaw:         `{"compilerOptions":{"experimentalDecorators":true,"useDefineForClassFields":false}}`,
		Plugins:             buildPlugins(opt),
		Conditions:          []string{"node"},
		LogLevel:            api.LogLevelSilent,
	})

	if len(result.Errors) > 0 {
		var msgs []string
		for _, e := range result.Errors {
			msgs = append(msgs, e.Text)
		}
		return nil, fmt.Errorf("esbuild batch bundle errors: %s", strings.Join(msgs, "; "))
	}

	// Map output files back to input paths
	bundles := make(map[string]string, len(result.OutputFiles))
	for _, of := range result.OutputFiles {
		base := filepath.Base(of.Path)
		inputPath, ok := keyToPath[base]
		if !ok {
			continue
		}
		code := string(of.Contents)
		code, hasTopLevelAwait := stripESMExports(code)
		if hasTopLevelAwait {
			code = "(async () => {\n" + code + "\n})();\n"
		}
		bundles[inputPath] = code
	}

	return bundles, nil
}

// buildPlugins returns the shared esbuild plugins for component bundling.
func buildPlugins(opt BundleOptions) []api.Plugin {
	return []api.Plugin{
		// Plugin 1: Handle CSS import assertions
		{
			Name: "css-import",
			Setup: func(build api.PluginBuild) {
				build.OnResolve(api.OnResolveOptions{Filter: `\.css$`},
					func(args api.OnResolveArgs) (api.OnResolveResult, error) {
						resolved := args.Path
						if !filepath.IsAbs(resolved) {
							resolved = filepath.Join(args.ResolveDir, args.Path)
						}
						if _, err := os.Stat(resolved); err == nil {
							return api.OnResolveResult{
								Path:      resolved,
								Namespace: "css-module",
							}, nil
						}
						return api.OnResolveResult{}, nil
					})
				build.OnLoad(api.OnLoadOptions{Filter: ".*", Namespace: "css-module"},
					func(args api.OnLoadArgs) (api.OnLoadResult, error) {
						data, err := os.ReadFile(args.Path)
						if err != nil {
							return api.OnLoadResult{}, err
						}
						css := string(data)
						if opt.Minify {
							minResult := api.Transform(css, api.TransformOptions{
								Loader:           api.LoaderCSS,
								MinifyWhitespace: true,
								MinifySyntax:     true,
							})
							if len(minResult.Errors) == 0 {
								css = strings.TrimSpace(string(minResult.Code))
							}
						}
						escaped := strings.ReplaceAll(css, "`", "\\`")
						escaped = strings.ReplaceAll(escaped, "$", "\\$")
						contents := fmt.Sprintf(`
							const cssText = %s;
							const sheet = { cssText, _$cssResult$: true, toString() { return cssText; } };
							export default sheet;
						`, "`"+escaped+"`")
						return api.OnLoadResult{
							Contents: &contents,
							Loader:   api.LoaderJS,
						}, nil
					})
			},
		},
		// Plugin 2: Stub Lit SSR packages
		{
			Name: "ssr-stubs",
			Setup: func(build api.PluginBuild) {
				build.OnResolve(api.OnResolveOptions{Filter: `^@lit-labs/ssr(-dom-shim)?(/.*)?$`},
					func(args api.OnResolveArgs) (api.OnResolveResult, error) {
						return api.OnResolveResult{
							Path:      args.Path,
							Namespace: "ssr-stub",
						}, nil
					})
				build.OnLoad(api.OnLoadOptions{Filter: ".*", Namespace: "ssr-stub"},
					func(args api.OnLoadArgs) (api.OnLoadResult, error) {
						contents := `
							export const HTMLElement = globalThis.HTMLElement;
							export const customElements = globalThis.customElements;
							export const Element = globalThis.Element;
							export const Event = globalThis.Event;
							export const CustomEvent = globalThis.CustomEvent;
							export const EventTarget = globalThis.EventTarget;
							export const CSSStyleSheet = globalThis.CSSStyleSheet || class CSSStyleSheet {};
							export function installWindowOnGlobal() {}
							export function getWindow() { return globalThis; }
							export function ariaMixinAttributes() { return []; }
						`
						return api.OnLoadResult{
							Contents: &contents,
							Loader:   api.LoaderJS,
						}, nil
					})
			},
		},
	}
}

// BundleSource bundles a component from inline JS/TS source code (no file on disk).
// Uses esbuild's Stdin option instead of EntryPoints.
func BundleSource(source string, opts ...BundleOptions) (string, error) {
	opt := BundleOptions{}
	if len(opts) > 0 {
		opt = opts[0]
	}

	sp, cp, _, err := ensureShimDir()
	if err != nil {
		return "", fmt.Errorf("preparing shim files: %w", err)
	}

	cwd, _ := os.Getwd()
	nodeModulesDir := findNodeModules(filepath.Join(cwd, "dummy"))

	result := api.Build(api.BuildOptions{
		Stdin: &api.StdinOptions{
			Contents:   source,
			ResolveDir: cwd,
			Loader:     api.LoaderTS,
		},
		Bundle:           true,
		Format:           api.FormatESModule,
		Target:           api.ES2022,
		Platform:         api.PlatformNeutral,
		Inject:           []string{sp, cp},
		Write:            false,
		MinifyWhitespace: opt.Minify,
		MinifySyntax:     opt.Minify,
		NodePaths:        []string{nodeModulesDir},
		TsconfigRaw:      `{"compilerOptions":{"experimentalDecorators":true,"useDefineForClassFields":false}}`,
		Plugins:          buildPlugins(opt),
		Conditions:       []string{"node"},
		LogLevel:         api.LogLevelSilent,
	})

	if len(result.Errors) > 0 {
		var msgs []string
		for _, e := range result.Errors {
			msgs = append(msgs, e.Text)
		}
		return "", fmt.Errorf("esbuild bundle errors: %s", strings.Join(msgs, "; "))
	}

	if len(result.OutputFiles) == 0 {
		return "", fmt.Errorf("esbuild produced no output")
	}

	code := string(result.OutputFiles[0].Contents)
	code, hasTopLevelAwait := stripESMExports(code)
	if hasTopLevelAwait {
		code = "(async () => {\n" + code + "\n})();\n"
	}

	return code, nil
}

// BundlePreload bundles a module and wraps it so its exports are registered
// in the __preloadedModules registry under the given name.
// Unlike BundleComponent, this captures ESM exports into the registry
// instead of stripping them.
func BundlePreload(modulePath string, name string) (string, error) {
	esm, err := bundleComponentRaw(modulePath, BundleOptions{})
	if err != nil {
		return "", err
	}

	exportNames := extractESMExportNames(esm)

	code, hasTopLevelAwait := stripESMExports(esm)
	if hasTopLevelAwait {
		code = "(async () => {\n" + code + "\n})();\n"
	}

	var capture strings.Builder
	capture.WriteString(fmt.Sprintf("\nglobalThis.__preloadedModules[%q] = {", name))
	for i, n := range exportNames {
		if i > 0 {
			capture.WriteString(", ")
		}
		capture.WriteString(fmt.Sprintf("%s: (typeof %s !== 'undefined' ? %s : undefined)", n, n, n))
	}
	capture.WriteString("};\n")

	return code + capture.String(), nil
}

// extractESMExportNames parses export declarations from raw ESM output and
// returns the local names. Handles single-line and multi-line export blocks
// as well as "Foo as Bar" aliases (returning the local name).
func extractESMExportNames(esm string) []string {
	var names []string
	inBlock := false
	for _, line := range strings.Split(esm, "\n") {
		trimmed := strings.TrimSpace(line)

		if inBlock {
			if strings.Contains(trimmed, "}") {
				// Final line of multi-line export block
				body := trimmed[:strings.Index(trimmed, "}")]
				names = append(names, parseExportList(body)...)
				inBlock = false
			} else {
				names = append(names, parseExportList(trimmed)...)
			}
			continue
		}

		if !strings.HasPrefix(trimmed, "export {") && !strings.HasPrefix(trimmed, "export{") {
			continue
		}

		if strings.Contains(trimmed, "}") {
			start := strings.Index(trimmed, "{")
			end := strings.LastIndex(trimmed, "}")
			if start >= 0 && end > start {
				names = append(names, parseExportList(trimmed[start+1:end])...)
			}
		} else {
			start := strings.Index(trimmed, "{")
			if start >= 0 {
				names = append(names, parseExportList(trimmed[start+1:])...)
			}
			inBlock = true
		}
	}
	return names
}

func parseExportList(body string) []string {
	var names []string
	for _, part := range strings.Split(body, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if idx := strings.Index(part, " as "); idx >= 0 {
			names = append(names, strings.TrimSpace(part[:idx]))
		} else {
			names = append(names, part)
		}
	}
	return names
}

// ResolveModulePath resolves a bare module specifier to a file path by
// looking up its entry point in node_modules/package.json.
func ResolveModulePath(specifier string, fromDir string) (string, error) {
	// If it's already a path, return as-is
	if strings.HasPrefix(specifier, ".") || strings.HasPrefix(specifier, "/") {
		return filepath.Abs(specifier)
	}

	// Find node_modules
	nmDir := findNodeModules(filepath.Join(fromDir, "dummy"))
	if nmDir == "" {
		return "", fmt.Errorf("node_modules not found from %s", fromDir)
	}

	pkgDir := filepath.Join(nmDir, specifier)
	pkgJSON := filepath.Join(pkgDir, "package.json")

	data, err := os.ReadFile(pkgJSON)
	if err != nil {
		// No package.json — try index.js or specifier.js directly
		for _, candidate := range []string{
			filepath.Join(pkgDir, "index.js"),
			filepath.Join(pkgDir, specifier+".js"),
		} {
			if _, err := os.Stat(candidate); err == nil {
				return candidate, nil
			}
		}
		return "", fmt.Errorf("cannot resolve module %q: %w", specifier, err)
	}

	// Parse package.json for "module", "main", or "exports" entry
	type PkgJSON struct {
		Module  string `json:"module"`
		Main    string `json:"main"`
	}
	var pkg PkgJSON
	if err := json.Unmarshal(data, &pkg); err != nil {
		return "", fmt.Errorf("parsing %s: %w", pkgJSON, err)
	}

	entry := pkg.Module
	if entry == "" {
		entry = pkg.Main
	}
	if entry == "" {
		entry = "index.js"
	}

	resolved := filepath.Join(pkgDir, entry)
	if _, err := os.Stat(resolved); err != nil {
		return "", fmt.Errorf("module entry %s not found", resolved)
	}
	return resolved, nil
}

// SaveBundle writes a bundle string to a file atomically.
func SaveBundle(bundle string, path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating directory: %w", err)
	}
	return fileutil.WriteFileAtomic(path, []byte(bundle), 0644)
}


// findNodeModules walks up from the given path to find node_modules.
func findNodeModules(fromPath string) string {
	dir := filepath.Dir(fromPath)
	for {
		candidate := filepath.Join(dir, "node_modules")
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			return candidate
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ""
}
