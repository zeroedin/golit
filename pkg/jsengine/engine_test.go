package jsengine

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/fastschema/qjs"
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

func TestEngine_RenderElement_EmptyRender_HasRootMarkers(t *testing.T) {
	source := `
		import { LitElement } from 'lit';
		class EmptyEl extends LitElement {}
		customElements.define('empty-el', EmptyEl);
	`
	bundle, err := BundleSource(source)
	if err != nil {
		t.Fatalf("bundling: %v", err)
	}

	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}
	defer engine.Close()

	if err := engine.LoadBundle(bundle); err != nil {
		t.Fatalf("loading bundle: %v", err)
	}

	result, err := engine.RenderElement("empty-el", map[string]string{})
	if err != nil {
		t.Fatalf("render: %v", err)
	}

	if !strings.Contains(result.HTML, "<!--lit-part-->") {
		t.Errorf("expected root <!--lit-part--> marker in HTML for component with no render(), got: %q", result.HTML)
	}
	if !strings.Contains(result.HTML, "<!--/lit-part-->") {
		t.Errorf("expected root <!--/lit-part--> marker in HTML for component with no render(), got: %q", result.HTML)
	}
}

func TestEngine_RenderElement_WithRender_HasDigestMarkers(t *testing.T) {
	bundle := bundleMyGreeting(t)

	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}
	defer engine.Close()

	if err := engine.LoadBundle(bundle); err != nil {
		t.Fatal(err)
	}

	result, err := engine.RenderElement("my-greeting", map[string]string{"name": "Test"})
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(result.HTML, "<!--lit-part ") {
		t.Errorf("expected <!--lit-part DIGEST--> marker in HTML, got: %q", result.HTML[:min(200, len(result.HTML))])
	}
	if !strings.Contains(result.HTML, "<!--/lit-part-->") {
		t.Errorf("expected <!--/lit-part--> closing marker in HTML, got: %q", result.HTML[:min(200, len(result.HTML))])
	}
}

func TestBundleStandaloneModule_DefaultExport(t *testing.T) {
	esm, err := BundleStandaloneModule("../../testdata/sources/css-module.js")
	if err != nil {
		t.Fatalf("bundling standalone module: %v", err)
	}

	if !strings.Contains(esm, "export") {
		t.Error("standalone module should preserve ESM exports")
	}

	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}
	defer engine.Close()

	if err := engine.LoadModule("test-css-module", esm); err != nil {
		t.Fatalf("loading module: %v", err)
	}

	script := `
		import mod from "test-css-module";
		globalThis.__testResult = JSON.stringify({
			hasDefault: mod !== undefined,
			cssText: mod ? mod.cssText : "",
		});
	`
	if _, err := engine.ctx.Eval("test-import.js", qjs.Code(script), qjs.TypeModule()); err != nil {
		t.Fatalf("eval module: %v", err)
	}
	result, err := engine.ctx.Eval("read-result.js", qjs.Code(`globalThis.__testResult`))
	if err != nil {
		t.Fatalf("eval: %v", err)
	}

	var output struct {
		HasDefault bool   `json:"hasDefault"`
		CSSText    string `json:"cssText"`
	}
	if err := json.Unmarshal([]byte(result.String()), &output); err != nil {
		t.Fatalf("parsing result: %v (raw: %s)", err, result.String())
	}
	if !output.HasDefault {
		t.Fatal("standalone module default export was not accessible via import")
	}
	if !strings.Contains(output.CSSText, "color: red") {
		t.Errorf("expected cssText to contain 'color: red', got: %q", output.CSSText)
	}
}

func TestEngine_RenderBatch_CacheHit(t *testing.T) {
	bundle := bundleMyGreeting(t)

	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}
	defer engine.Close()
	if err := engine.LoadBundle(bundle); err != nil {
		t.Fatal(err)
	}

	attrs := map[string]string{"name": "Cache"}
	requests := []BatchRequest{
		{ID: 1, TagName: "my-greeting", Attrs: attrs},
		{ID: 2, TagName: "my-greeting", Attrs: attrs},
		{ID: 3, TagName: "my-greeting", Attrs: attrs},
	}

	results, err := engine.RenderBatch(requests)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}

	for i, r := range results {
		if r.Error != "" {
			t.Errorf("result[%d] error: %s", i, r.Error)
		}
		if !strings.Contains(r.HTML, "Cache") {
			t.Errorf("result[%d] missing 'Cache' in HTML", i)
		}
	}

	if results[0].HTML != results[1].HTML || results[1].HTML != results[2].HTML {
		t.Error("identical requests should produce identical HTML via cache")
	}
}

func TestEngine_RenderBatch_CacheMiss_DifferentAttrs(t *testing.T) {
	bundle := bundleMyGreeting(t)

	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}
	defer engine.Close()
	if err := engine.LoadBundle(bundle); err != nil {
		t.Fatal(err)
	}

	requests := []BatchRequest{
		{ID: 1, TagName: "my-greeting", Attrs: map[string]string{"name": "Alice"}},
		{ID: 2, TagName: "my-greeting", Attrs: map[string]string{"name": "Bob"}},
	}

	results, err := engine.RenderBatch(requests)
	if err != nil {
		t.Fatal(err)
	}

	if results[0].HTML == results[1].HTML {
		t.Error("different attrs should produce different HTML")
	}
	if !strings.Contains(results[0].HTML, "Alice") {
		t.Error("result[0] missing 'Alice'")
	}
	if !strings.Contains(results[1].HTML, "Bob") {
		t.Error("result[1] missing 'Bob'")
	}
}

func TestEngine_RenderBatch_ErrorPropagation(t *testing.T) {
	bundle := bundleMyGreeting(t)

	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}
	defer engine.Close()
	if err := engine.LoadBundle(bundle); err != nil {
		t.Fatal(err)
	}

	requests := []BatchRequest{
		{ID: 1, TagName: "my-greeting", Attrs: map[string]string{"name": "OK"}},
		{ID: 2, TagName: "nonexistent-element", Attrs: map[string]string{}},
	}

	results, err := engine.RenderBatch(requests)
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].Error != "" {
		t.Errorf("result[0] unexpected error: %s", results[0].Error)
	}
	if !strings.Contains(results[0].HTML, "OK") {
		t.Error("result[0] missing 'OK' in HTML")
	}
	if results[1].Error == "" {
		t.Error("expected error for unregistered element in batch")
	}
}

func TestEngine_RenderBatch_Empty(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}
	defer engine.Close()

	results, err := engine.RenderBatch(nil)
	if err != nil {
		t.Fatal(err)
	}
	if results != nil {
		t.Errorf("expected nil results for empty batch, got %v", results)
	}
}

func TestRenderCacheKey(t *testing.T) {
	k1 := renderCacheKey("my-el", map[string]string{"a": "1", "b": "2"})
	k2 := renderCacheKey("my-el", map[string]string{"b": "2", "a": "1"})
	if k1 != k2 {
		t.Error("cache key should be deterministic regardless of map iteration order")
	}

	k3 := renderCacheKey("my-el", map[string]string{})
	k4 := renderCacheKey("my-el", nil)
	if k3 != "my-el" || k4 != "my-el" {
		t.Error("empty attrs should produce tag-name-only key")
	}
}
