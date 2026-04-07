package jsengine

import (
	"strings"
	"testing"
)

func bundleMyGreeting(t *testing.T) string {
	t.Helper()
	bundle, err := BundleComponent("../../testdata/sources/my-greeting.js")
	if err != nil {
		t.Fatalf("bundling my-greeting: %v", err)
	}
	return bundle
}

func TestEngine_RenderMyGreeting(t *testing.T) {
	bundle := bundleMyGreeting(t)

	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}
	defer engine.Close()

	if err := engine.LoadBundle(bundle); err != nil {
		t.Fatal(err)
	}

	result, err := engine.RenderElement("my-greeting", map[string]string{
		"name": "World",
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("HTML: %s", result.HTML)
	t.Logf("CSS: %s", result.CSS[:min(80, len(result.CSS))])

	if !strings.Contains(result.HTML, "World") {
		t.Error("missing 'World' in output")
	}
	if !strings.Contains(result.HTML, "Hello") {
		t.Error("missing 'Hello' in output")
	}
}

func TestEngine_RenderMyGreeting_DifferentNames(t *testing.T) {
	bundle := bundleMyGreeting(t)

	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}
	defer engine.Close()
	if err := engine.LoadBundle(bundle); err != nil {
		t.Fatal(err)
	}

	for _, name := range []string{"Alice", "Go", "Hugo"} {
		t.Run(name, func(t *testing.T) {
			result, err := engine.RenderElement("my-greeting", map[string]string{
				"name": name,
			})
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(result.HTML, name) {
				t.Errorf("missing %q in: %s", name, result.HTML)
			}
		})
	}
}

func TestEngine_StyleExtraction(t *testing.T) {
	bundle := bundleMyGreeting(t)

	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}
	defer engine.Close()
	if err := engine.LoadBundle(bundle); err != nil {
		t.Fatal(err)
	}

	result, err := engine.RenderElement("my-greeting", map[string]string{
		"name": "Test",
	})
	if err != nil {
		t.Fatal(err)
	}

	if result.CSS == "" {
		t.Error("expected non-empty CSS from style extraction")
	}
	t.Logf("CSS: %s", result.CSS[:min(100, len(result.CSS))])
}

func TestShimDynamicImports(t *testing.T) {
	e := &Engine{preloadModules: []string{"prism-esm", "other-mod"}}

	cases := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "no imports",
			input: `var x = 1;`,
			want:  `var x = 1;`,
		},
		{
			name:  "double-quoted import",
			input: `import("prism-esm")`,
			want:  `Promise.resolve(globalThis.__preloadedModules["prism-esm"] || {})/*golit-shimmed:import("prism-esm")*/`,
		},
		{
			name:  "single-quoted import",
			input: `import('prism-esm')`,
			want:  `Promise.resolve(globalThis.__preloadedModules["prism-esm"] || {})/*golit-shimmed:import('prism-esm')*/`,
		},
		{
			name:  "subpath import",
			input: `import("prism-esm/components/prism-css.js")`,
			want:  `Promise.resolve(globalThis.__preloadedModules["prism-esm"] || {})/*golit-shimmed:import("prism-esm/components/prism-css.js")*/`,
		},
		{
			name:  "non-preloaded module unchanged",
			input: `import("unknown-mod")`,
			want:  `import("unknown-mod")`,
		},
		{
			name:  "second module matched",
			input: `import("other-mod")`,
			want:  `Promise.resolve(globalThis.__preloadedModules["other-mod"] || {})/*golit-shimmed:import("other-mod")*/`,
		},
		{
			name:  "multiple imports in one string",
			input: `import("prism-esm"); import('other-mod');`,
			want:  `Promise.resolve(globalThis.__preloadedModules["prism-esm"] || {})/*golit-shimmed:import("prism-esm")*/; Promise.resolve(globalThis.__preloadedModules["other-mod"] || {})/*golit-shimmed:import('other-mod')*/;`,
		},
		{
			name:  "prefix of preloaded module is not matched",
			input: `import("prism-esm-extra")`,
			want:  `import("prism-esm-extra")`,
		},
		{
			name:  "prefix with single quotes not matched",
			input: `import('other-module')`,
			want:  `import('other-module')`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := e.shimDynamicImports(tc.input)
			if got != tc.want {
				t.Errorf("got:\n  %s\nwant:\n  %s", got, tc.want)
			}
		})
	}
}

func TestShimDynamicImports_PrefixBoundary(t *testing.T) {
	e := &Engine{preloadModules: []string{"lit"}}

	cases := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "exact match is shimmed",
			input: `import("lit")`,
			want:  `Promise.resolve(globalThis.__preloadedModules["lit"] || {})/*golit-shimmed:import("lit")*/`,
		},
		{
			name:  "subpath is shimmed",
			input: `import("lit/decorators.js")`,
			want:  `Promise.resolve(globalThis.__preloadedModules["lit"] || {})/*golit-shimmed:import("lit/decorators.js")*/`,
		},
		{
			name:  "different module with same prefix is NOT shimmed",
			input: `import("lit-html")`,
			want:  `import("lit-html")`,
		},
		{
			name:  "hyphenated suffix not shimmed single quotes",
			input: `import('lit-element')`,
			want:  `import('lit-element')`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := e.shimDynamicImports(tc.input)
			if got != tc.want {
				t.Errorf("got:\n  %s\nwant:\n  %s", got, tc.want)
			}
		})
	}
}

func TestShimDynamicImports_NoModules(t *testing.T) {
	e := &Engine{}
	input := `import("something")`
	if got := e.shimDynamicImports(input); got != input {
		t.Errorf("expected no change when preloadModules is empty, got %q", got)
	}
}

func TestShimDynamicImports_RuntimeExternals(t *testing.T) {
	e := &Engine{runtimeExternals: []string{
		"@rhds/tokens", "@rhds/tokens/*",
		"@rhds/icons", "@rhds/icons/*",
	}}

	cases := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "exact match rewrites to runtime",
			input: `import("@rhds/tokens")`,
			want:  `import("@golit/runtime")/*golit-runtime:import("@rhds/tokens")*/`,
		},
		{
			name:  "subpath rewrites to runtime",
			input: `import("@rhds/tokens/css/default-theme.css.js")`,
			want:  `import("@golit/runtime")/*golit-runtime:import("@rhds/tokens/css/default-theme.css.js")*/`,
		},
		{
			name:  "single quotes work too",
			input: `import('@rhds/icons/ui/check.js')`,
			want:  `import('@golit/runtime')/*golit-runtime:import('@rhds/icons/ui/check.js')*/`,
		},
		{
			name:  "non-external package unchanged",
			input: `import("@other/pkg")`,
			want:  `import("@other/pkg")`,
		},
		{
			name:  "local import unchanged",
			input: `import("./local.js")`,
			want:  `import("./local.js")`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := e.shimDynamicImports(tc.input)
			if got != tc.want {
				t.Errorf("got:\n  %s\nwant:\n  %s", got, tc.want)
			}
		})
	}
}

func TestShimDynamicImports_PreloadTakesPrecedence(t *testing.T) {
	e := &Engine{
		preloadModules:   []string{"prism-esm"},
		runtimeExternals: []string{"prism-esm", "prism-esm/*", "@rhds/tokens", "@rhds/tokens/*"},
	}

	got := e.shimDynamicImports(`import("prism-esm/components/prism-css.js")`)
	want := `Promise.resolve(globalThis.__preloadedModules["prism-esm"] || {})/*golit-shimmed:import("prism-esm/components/prism-css.js")*/`
	if got != want {
		t.Errorf("preload should take precedence over runtime externals\ngot:\n  %s\nwant:\n  %s", got, want)
	}

	got = e.shimDynamicImports(`import("@rhds/tokens/css/default-theme.css.js")`)
	want = `import("@golit/runtime")/*golit-runtime:import("@rhds/tokens/css/default-theme.css.js")*/`
	if got != want {
		t.Errorf("non-preloaded externals should use runtime rewrite\ngot:\n  %s\nwant:\n  %s", got, want)
	}
}

func TestShimDynamicImports_BothEmpty(t *testing.T) {
	e := &Engine{}
	input := `import("@rhds/tokens/css/default-theme.css.js")`
	if got := e.shimDynamicImports(input); got != input {
		t.Errorf("expected no change when both preload and externals are empty, got %q", got)
	}
}

func TestEngine_UnregisteredElement(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}
	defer engine.Close()

	_, err = engine.RenderElement("unknown-element", map[string]string{})
	if err == nil {
		t.Error("expected error for unregistered element")
	}
}
