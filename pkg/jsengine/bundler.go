// Package jsengine provides a QJS-based JavaScript execution engine
// for server-side rendering Lit components with full expression fidelity.
package jsengine

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/zeroedin/golit/pkg/fileutil"
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

	// ExternalPackages lists package specifiers to mark as external.
	// Discovered via DiscoverExternalPackages and passed to
	// BundleComponentModule(s) so shared deps stay as import statements.
	ExternalPackages []string
}

// DiscoverExternalPackages runs esbuild with Metafile enabled on all
// component entry points to discover which node_modules packages they
// depend on. Returns esbuild-compatible external patterns (pkg + pkg/*).
func DiscoverExternalPackages(componentPaths []string, nodeModulesDir string, opts ...BundleOptions) ([]string, error) {
	opt := BundleOptions{}
	if len(opts) > 0 {
		opt = opts[0]
	}

	sp, cp, _, err := ensureShimDir()
	if err != nil {
		return nil, fmt.Errorf("preparing shim files: %w", err)
	}

	type entry struct {
		absPath string
		key     string
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
		key := fmt.Sprintf("disc_%d", len(entries))
		entries = append(entries, entry{absPath: absPath, key: key})
	}

	if len(entries) == 0 {
		return nil, nil
	}

	_, _, sd, err := ensureShimDir()
	if err != nil {
		return nil, fmt.Errorf("preparing shim files: %w", err)
	}

	esbuildEntries := make([]api.EntryPoint, len(entries))
	for i, e := range entries {
		esbuildEntries[i] = api.EntryPoint{
			InputPath:  e.absPath,
			OutputPath: e.key,
		}
	}

	result := api.Build(api.BuildOptions{
		EntryPointsAdvanced: esbuildEntries,
		Bundle:              true,
		Format:              api.FormatESModule,
		Target:              api.ES2022,
		Platform:            api.PlatformNeutral,
		Inject:              []string{sp, cp},
		Metafile:            true,
		Write:               false,
		Outdir:              sd,
		NodePaths:           []string{nodeModulesDir},
		TsconfigRaw:         `{"compilerOptions":{"experimentalDecorators":true,"useDefineForClassFields":false}}`,
		Plugins:             buildPlugins(opt),
		Conditions:          []string{"node"},
		LogLevel:            api.LogLevelSilent,
	})

	if len(result.Errors) > 0 {
		for _, e := range result.Errors {
			fmt.Fprintf(os.Stderr, "golit: warning: discovery: %s\n", e.Text)
		}
	}

	if result.Metafile == "" {
		return nil, nil
	}
	return parseMetafilePackages(result.Metafile)
}

// parseMetafilePackages extracts node_modules package names from esbuild
// metafile JSON and returns esbuild external patterns (pkg + pkg/*).
func parseMetafilePackages(metafile string) ([]string, error) {
	var meta struct {
		Inputs map[string]json.RawMessage `json:"inputs"`
	}
	if err := json.Unmarshal([]byte(metafile), &meta); err != nil {
		return nil, fmt.Errorf("parsing metafile: %w", err)
	}

	pkgs := make(map[string]bool)
	for inputPath := range meta.Inputs {
		pkg := extractNodeModulesPackage(inputPath)
		if pkg != "" {
			pkgs[pkg] = true
		}
	}

	patterns := make([]string, 0, len(pkgs)*2)
	for pkg := range pkgs {
		patterns = append(patterns, pkg, pkg+"/*")
	}
	sort.Strings(patterns)
	return patterns, nil
}

// extractNodeModulesPackage extracts the package name from a path that
// goes through node_modules. Returns "" if the path is not from node_modules.
// Handles scoped packages (@scope/pkg) and unscoped (pkg).
func extractNodeModulesPackage(inputPath string) string {
	inputPath = filepath.ToSlash(inputPath)
	const marker = "node_modules/"
	idx := strings.LastIndex(inputPath, marker)
	if idx < 0 {
		return ""
	}
	rest := inputPath[idx+len(marker):]
	if rest == "" {
		return ""
	}

	parts := strings.SplitN(rest, "/", 3)
	if strings.HasPrefix(parts[0], "@") {
		if len(parts) < 2 {
			return parts[0]
		}
		return parts[0] + "/" + parts[1]
	}
	return parts[0]
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

	code, hasTopLevelAwait, _ := stripESMExports(esm)
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
// can be evaluated in QJS script mode. Returns the stripped code, whether
// top-level await was found, and whether a default export was found.
// Default exports are assigned to __golit_default_export for later capture.
func stripESMExports(code string) (string, bool, bool) {
	hasTopLevelAwait := false
	hasDefaultExport := false
	inExportBlock := false
	var b strings.Builder
	b.Grow(len(code))

	lineNum := 0
	rest := code
	for len(rest) > 0 {
		nlIdx := strings.IndexByte(rest, '\n')
		var line string
		if nlIdx >= 0 {
			line = rest[:nlIdx]
			rest = rest[nlIdx+1:]
		} else {
			line = rest
			rest = ""
		}

		if lineNum > 0 {
			b.WriteByte('\n')
		}
		lineNum++

		trimmed := strings.TrimSpace(line)

		if inExportBlock {
			if strings.Contains(trimmed, "};") || trimmed == "}" {
				inExportBlock = false
			}
			continue
		}

		if (strings.HasPrefix(trimmed, "export {") || strings.HasPrefix(trimmed, "export{")) &&
			strings.Contains(trimmed, "}") {
			continue
		}
		if strings.HasPrefix(trimmed, "export {") || strings.HasPrefix(trimmed, "export{") {
			inExportBlock = true
			continue
		}
		if strings.HasPrefix(trimmed, "export default ") {
			b.WriteString(strings.Replace(line, "export default ", "var __golit_default_export = ", 1))
			hasDefaultExport = true
			continue
		}
		if strings.HasPrefix(trimmed, "export var ") ||
			strings.HasPrefix(trimmed, "export let ") ||
			strings.HasPrefix(trimmed, "export const ") ||
			strings.HasPrefix(trimmed, "export function ") ||
			strings.HasPrefix(trimmed, "export class ") ||
			strings.HasPrefix(trimmed, "export async ") {
			b.WriteString(strings.Replace(line, "export ", "", 1))
			continue
		}
		if !hasTopLevelAwait && strings.HasPrefix(trimmed, "await ") {
			indent := len(line) - len(strings.TrimLeft(line, " \t"))
			if indent <= 2 {
				hasTopLevelAwait = true
			}
		}

		b.WriteString(line)
	}

	return b.String(), hasTopLevelAwait, hasDefaultExport
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

// BundleStandaloneModule bundles a module file into a self-contained ES module
// with exports preserved. Used for registering dynamic import targets as named
// QJS modules so import("specifier") resolves natively.
func BundleStandaloneModule(modulePath string) (string, error) {
	return bundleComponentRaw(modulePath, BundleOptions{})
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
	code, hasTopLevelAwait, _ := stripESMExports(code)
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

	code, hasTopLevelAwait, hasDefaultExport := stripESMExports(esm)
	if hasTopLevelAwait {
		code = "(async () => {\n" + code + "\n})();\n"
	}

	var capture strings.Builder
	capture.WriteString(fmt.Sprintf("\nglobalThis.__preloadedModules[%q] = {", name))
	first := true
	if hasDefaultExport {
		capture.WriteString("default: (typeof __golit_default_export !== 'undefined' ? __golit_default_export : undefined)")
		first = false
	}
	for _, n := range exportNames {
		if !first {
			capture.WriteString(", ")
		}
		capture.WriteString(fmt.Sprintf("%s: (typeof %s !== 'undefined' ? %s : undefined)", n, n, n))
		first = false
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

// BundleSharedRuntime produces a single ES module containing all shared
// dependencies that the thin component modules import. The output keeps
// ESM export syntax intact so it can be loaded via
// Engine.LoadModule("@golit/runtime", source).
// modules is the map of thin module sources (from BundleComponentModules)
// used to discover which external specifiers to include.
func BundleSharedRuntime(nodeModulesDir string, modules map[string]string, opts ...BundleOptions) (string, error) {
	opt := BundleOptions{}
	if len(opts) > 0 {
		opt = opts[0]
	}

	sp, cp, _, err := ensureShimDir()
	if err != nil {
		return "", fmt.Errorf("preparing shim files: %w", err)
	}

	entrySource := buildRuntimeEntryFromModules(modules)

	entryDir, err := os.MkdirTemp("", "golit-runtime-entry-*")
	if err != nil {
		return "", fmt.Errorf("creating temp dir: %w", err)
	}
	defer os.RemoveAll(entryDir)

	entryPath := filepath.Join(entryDir, "_runtime_entry.js")
	if err := os.WriteFile(entryPath, []byte(entrySource), 0644); err != nil {
		return "", fmt.Errorf("writing runtime entry: %w", err)
	}

	result := api.Build(api.BuildOptions{
		EntryPoints:      []string{entryPath},
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
		return "", fmt.Errorf("esbuild shared runtime errors: %s", strings.Join(msgs, "; "))
	}

	if len(result.OutputFiles) == 0 {
		return "", fmt.Errorf("esbuild produced no output for shared runtime")
	}

	return string(result.OutputFiles[0].Contents), nil
}

// BundleComponentModule produces an ES module for a component with shared
// dependencies marked as external. The output keeps import/export statements
// so it can be loaded via Engine.EvalModule. Callers should use
// RewriteModuleImports to rewrite external specifiers to "@golit/runtime"
// after extracting import specifiers for runtime building.
// The domshim and template collector are injected so the module works
// both with and without a shared runtime.
func BundleComponentModule(componentPath string, opts ...BundleOptions) (string, error) {
	opt := BundleOptions{}
	if len(opts) > 0 {
		opt = opts[0]
	}
	return bundleComponentModule(componentPath, opt)
}

func bundleComponentModule(componentPath string, opt BundleOptions) (string, error) {
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
		External:         opt.ExternalPackages,
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
		return "", fmt.Errorf("esbuild module errors: %s", strings.Join(msgs, "; "))
	}

	if len(result.OutputFiles) == 0 {
		return "", fmt.Errorf("esbuild produced no output")
	}

	return string(result.OutputFiles[0].Contents), nil
}

// BundleComponentModules produces thin ES modules for multiple components in
// a single esbuild invocation. Returns a map from input path to module source.
func BundleComponentModules(componentPaths []string, opts ...BundleOptions) (map[string]string, error) {
	opt := BundleOptions{}
	if len(opts) > 0 {
		opt = opts[0]
	}

	if len(componentPaths) == 0 {
		return nil, nil
	}

	type entry struct {
		absPath string
		key     string
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

	nodeModulesDir := findNodeModules(entries[0].absPath)

	esbuildEntries := make([]api.EntryPoint, len(entries))
	keyToPath := make(map[string]string, len(entries))
	for i, e := range entries {
		esbuildEntries[i] = api.EntryPoint{
			InputPath:  e.absPath,
			OutputPath: e.key,
		}
		keyToPath[e.key+".js"] = e.absPath
	}

	result := api.Build(api.BuildOptions{
		EntryPointsAdvanced: esbuildEntries,
		Bundle:              true,
		Format:              api.FormatESModule,
		Target:              api.ES2022,
		Platform:            api.PlatformNeutral,
		External:            opt.ExternalPackages,
		Inject:              []string{sp, cp},
		Write:               false,
		Outdir:              sd,
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
		return nil, fmt.Errorf("esbuild batch module errors: %s", strings.Join(msgs, "; "))
	}

	modules := make(map[string]string, len(result.OutputFiles))
	for _, of := range result.OutputFiles {
		base := filepath.Base(of.Path)
		inputPath, ok := keyToPath[base]
		if !ok {
			continue
		}
		modules[inputPath] = string(of.Contents)
	}

	return modules, nil
}

// RewriteModuleImports rewrites import specifiers for shared packages
// in all module sources to "@golit/runtime". The externals list determines
// which specifiers are considered shared. Call this after
// extractExternalImports if you need the raw specifiers for runtime building.
func RewriteModuleImports(modules map[string]string, externals []string) map[string]string {
	rewritten := make(map[string]string, len(modules))
	for k, v := range modules {
		rewritten[k] = rewriteImportsToRuntime(v, externals)
	}
	return rewritten
}

func rewriteImportsToRuntime(code string, externals []string) string {
	lines := strings.Split(code, "\n")
	var b strings.Builder
	b.Grow(len(code))

	for i, line := range lines {
		if i > 0 {
			b.WriteByte('\n')
		}
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "import ") {
			if strings.Contains(trimmed, " from ") {
				if rewritten, ok := rewriteImportLine(line, externals); ok {
					b.WriteString(rewritten)
					continue
				}
			} else if rewritten, ok := rewriteSideEffectImport(line, externals); ok {
				b.WriteString(rewritten)
				continue
			}
		}
		b.WriteString(line)
	}

	return b.String()
}

func rewriteImportLine(line string, externals []string) (string, bool) {
	fromIdx := strings.LastIndex(line, " from ")
	if fromIdx < 0 {
		return "", false
	}
	specPart := strings.TrimSpace(line[fromIdx+6:])
	specPart = strings.TrimSuffix(specPart, ";")
	specPart = strings.TrimSpace(specPart)

	if len(specPart) < 2 {
		return "", false
	}
	quote := specPart[0]
	if quote != '\'' && quote != '"' {
		return "", false
	}
	endQuote := strings.IndexByte(specPart[1:], quote)
	if endQuote < 0 {
		return "", false
	}
	specifier := specPart[1 : 1+endQuote]

	if matchesExternals(specifier, externals) {
		if isDefaultImport(line[:fromIdx]) {
			return "", false
		}
		return line[:fromIdx] + " from " + string(quote) + "@golit/runtime" + string(quote) + ";", true
	}
	return "", false
}

// isDefaultImport checks if the import clause is a default import
// (e.g. "import styles" or "import foo, { bar }").
func isDefaultImport(importClause string) bool {
	trimmed := strings.TrimSpace(importClause)
	trimmed = strings.TrimPrefix(trimmed, "import ")
	trimmed = strings.TrimSpace(trimmed)
	if trimmed == "" {
		return false
	}
	if strings.HasPrefix(trimmed, "{") || strings.HasPrefix(trimmed, "*") {
		return false
	}
	return true
}

func rewriteSideEffectImport(line string, externals []string) (string, bool) {
	trimmed := strings.TrimSpace(line)
	trimmed = strings.TrimPrefix(trimmed, "import ")
	trimmed = strings.TrimSuffix(trimmed, ";")
	trimmed = strings.TrimSpace(trimmed)

	if len(trimmed) < 2 {
		return "", false
	}
	quote := trimmed[0]
	if quote != '\'' && quote != '"' {
		return "", false
	}
	endQuote := strings.IndexByte(trimmed[1:], quote)
	if endQuote < 0 {
		return "", false
	}
	specifier := trimmed[1 : 1+endQuote]

	if matchesExternals(specifier, externals) {
		return "/* side-effect import included in @golit/runtime: " + specifier + " */", true
	}
	return "", false
}

func matchesExternals(specifier string, externals []string) bool {
	for _, pkg := range externals {
		if strings.HasSuffix(pkg, "/*") {
			prefix := strings.TrimSuffix(pkg, "/*")
			if specifier == prefix || strings.HasPrefix(specifier, prefix+"/") {
				return true
			}
		} else if specifier == pkg {
			return true
		}
	}
	return false
}

// extractExternalImports scans thin module sources for import statements
// that reference external packages (not @golit/runtime or relative paths).
// Returns two sets: re-export specifiers (from named imports) and
// side-effect specifiers (from bare imports).
func extractExternalImports(modules map[string]string) (reexports []string, sideEffects []string) {
	reexportSet := make(map[string]bool)
	sideEffectSet := make(map[string]bool)

	for _, source := range modules {
		for _, line := range strings.Split(source, "\n") {
			trimmed := strings.TrimSpace(line)
			if !strings.HasPrefix(trimmed, "import ") {
				continue
			}

			if strings.Contains(trimmed, " from ") {
				spec := extractFromSpecifier(trimmed)
				if spec != "" && !isLocalOrRuntime(spec) {
					reexportSet[spec] = true
				}
			} else {
				spec := extractBareImportSpecifier(trimmed)
				if spec != "" && !isLocalOrRuntime(spec) {
					sideEffectSet[spec] = true
				}
			}
		}
	}

	reexports = make([]string, 0, len(reexportSet))
	for spec := range reexportSet {
		reexports = append(reexports, spec)
	}
	sort.Strings(reexports)

	sideEffects = make([]string, 0, len(sideEffectSet))
	for spec := range sideEffectSet {
		sideEffects = append(sideEffects, spec)
	}
	sort.Strings(sideEffects)

	return reexports, sideEffects
}

func extractFromSpecifier(line string) string {
	fromIdx := strings.LastIndex(line, " from ")
	if fromIdx < 0 {
		return ""
	}
	spec := strings.TrimSpace(line[fromIdx+6:])
	spec = strings.TrimSuffix(spec, ";")
	spec = strings.TrimSpace(spec)
	if len(spec) < 2 {
		return ""
	}
	quote := spec[0]
	if quote != '\'' && quote != '"' {
		return ""
	}
	end := strings.IndexByte(spec[1:], quote)
	if end < 0 {
		return ""
	}
	return spec[1 : 1+end]
}

func extractBareImportSpecifier(line string) string {
	trimmed := strings.TrimSpace(line)
	trimmed = strings.TrimPrefix(trimmed, "import ")
	trimmed = strings.TrimSuffix(trimmed, ";")
	trimmed = strings.TrimSpace(trimmed)
	if len(trimmed) < 2 {
		return ""
	}
	quote := trimmed[0]
	if quote != '\'' && quote != '"' {
		return ""
	}
	end := strings.IndexByte(trimmed[1:], quote)
	if end < 0 {
		return ""
	}
	return trimmed[1 : 1+end]
}

func isLocalOrRuntime(specifier string) bool {
	return specifier == "@golit/runtime" ||
		strings.HasPrefix(specifier, "./") ||
		strings.HasPrefix(specifier, "../") ||
		strings.HasPrefix(specifier, "/")
}

// buildRuntimeEntryFromModules generates the shared runtime entry source
// by scanning thin module sources for their external import specifiers.
func buildRuntimeEntryFromModules(modules map[string]string) string {
	reexports, sideEffects := extractExternalImports(modules)

	var b strings.Builder
	for _, spec := range sideEffects {
		b.WriteString(fmt.Sprintf("import '%s';\n", spec))
	}
	for _, spec := range reexports {
		b.WriteString(fmt.Sprintf("export * from '%s';\n", spec))
	}

	if b.Len() == 0 {
		return "export {};\n"
	}
	return b.String()
}

// ExtractDynamicImportTargets scans thin module sources for dynamic import()
// calls and returns non-local, non-runtime specifiers.
// Matches patterns like: import("@rhds/tokens/css/default-theme.css.js")
func ExtractDynamicImportTargets(modules map[string]string) []string {
	return extractDynamicImportTargets(modules)
}

// ExtractUnrewrittenImports scans rewritten thin module sources for static
// imports that still reference external specifiers (not @golit/runtime, not
// relative). These are default imports that were intentionally not rewritten
// and need standalone module registration.
func ExtractUnrewrittenImports(modules map[string]string) []string {
	seen := make(map[string]bool)
	var specifiers []string

	for _, source := range modules {
		for _, line := range strings.Split(source, "\n") {
			trimmed := strings.TrimSpace(line)
			if !strings.HasPrefix(trimmed, "import ") || !strings.Contains(trimmed, " from ") {
				continue
			}
			spec := extractFromSpecifier(trimmed)
			if spec != "" && !isLocalOrRuntime(spec) && !seen[spec] {
				seen[spec] = true
				specifiers = append(specifiers, spec)
			}
		}
	}
	sort.Strings(specifiers)
	return specifiers
}

func extractDynamicImportTargets(modules map[string]string) []string {
	seen := make(map[string]bool)
	var targets []string

	for _, source := range modules {
		pos := 0
		for pos < len(source) {
			idx := strings.Index(source[pos:], "import(")
			if idx < 0 {
				break
			}
			start := pos + idx + len("import(")
			if start >= len(source) {
				break
			}
			quote := source[start]
			if quote != '"' && quote != '\'' {
				pos = start
				continue
			}
			closeStr := string(quote) + ")"
			end := strings.Index(source[start+1:], closeStr)
			if end < 0 {
				pos = start + 1
				continue
			}
			spec := source[start+1 : start+1+end]
			if !isLocalOrRuntime(spec) && !seen[spec] {
				seen[spec] = true
				targets = append(targets, spec)
			}
			pos = start + 1 + end + len(closeStr)
		}
	}
	sort.Strings(targets)
	return targets
}

// ResolveModulePath resolves a module specifier to a file path.
// Handles bare package names (e.g. "lit"), scoped packages ("@rhds/tokens"),
// and subpath specifiers ("@rhds/tokens/css/default-theme.css.js").
func ResolveModulePath(specifier string, fromDir string) (string, error) {
	if strings.HasPrefix(specifier, ".") || strings.HasPrefix(specifier, "/") {
		return filepath.Abs(specifier)
	}

	nmDir := findNodeModules(filepath.Join(fromDir, "dummy"))
	if nmDir == "" {
		return "", fmt.Errorf("node_modules not found from %s", fromDir)
	}

	// Try the specifier as a direct file path under node_modules first.
	// This handles subpath specifiers like "@rhds/tokens/css/default-theme.css.js"
	// or "prism-esm/components/prism-css.js".
	directPath := filepath.Join(nmDir, specifier)
	if info, err := os.Stat(directPath); err == nil && !info.IsDir() {
		return directPath, nil
	}

	// Extract the package name for package.json lookup.
	pkgName := extractPackageName(specifier)
	pkgDir := filepath.Join(nmDir, pkgName)
	pkgJSON := filepath.Join(pkgDir, "package.json")

	data, err := os.ReadFile(pkgJSON)
	if err != nil {
		for _, candidate := range []string{
			filepath.Join(pkgDir, "index.js"),
			filepath.Join(pkgDir, pkgName+".js"),
		} {
			if _, err := os.Stat(candidate); err == nil {
				return candidate, nil
			}
		}
		return "", fmt.Errorf("cannot resolve module %q: %w", specifier, err)
	}

	type PkgJSON struct {
		Module string `json:"module"`
		Main   string `json:"main"`
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

// extractPackageName returns the npm package name from a specifier.
// Handles scoped packages: "@scope/pkg/sub/path" -> "@scope/pkg"
// and unscoped: "pkg/sub/path" -> "pkg"
func extractPackageName(specifier string) string {
	parts := strings.SplitN(specifier, "/", 3)
	if strings.HasPrefix(parts[0], "@") {
		if len(parts) >= 2 {
			return parts[0] + "/" + parts[1]
		}
		return parts[0]
	}
	return parts[0]
}

// SaveBundle writes a bundle string to a file atomically.
func SaveBundle(bundle string, path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating directory: %w", err)
	}
	return fileutil.WriteFileAtomic(path, []byte(bundle), 0644)
}

// FindNodeModules walks up from the given path to find node_modules.
func FindNodeModules(fromPath string) string {
	return findNodeModules(fromPath)
}

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
