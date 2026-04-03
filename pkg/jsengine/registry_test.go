package jsengine

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDiscoverTagNameFast_ExplicitDefine(t *testing.T) {
	cases := []struct {
		name   string
		input  string
		want   string
		wantOK bool
	}{
		{
			name:   "single-quoted define",
			input:  `customElements.define('my-greeting', MyGreeting);`,
			want:   "my-greeting",
			wantOK: true,
		},
		{
			name:   "double-quoted define",
			input:  `customElements.define("my-card", MyCard);`,
			want:   "my-card",
			wantOK: true,
		},
		{
			name:   "whitespace around dot and paren",
			input:  `customElements . define ( 'rh-badge' , RhBadge)`,
			want:   "rh-badge",
			wantOK: true,
		},
		{
			name:   "multi-hyphen tag name",
			input:  `customElements.define('my-cool-widget', MyCoolWidget);`,
			want:   "my-cool-widget",
			wantOK: true,
		},
		{
			name:   "tag with digits",
			input:  `customElements.define('x-item2', XItem2);`,
			want:   "x-item2",
			wantOK: true,
		},
		{
			name:   "multiple defines returns last",
			input:  `customElements.define('dep-a', DepA); customElements.define('my-app', MyApp);`,
			want:   "my-app",
			wantOK: true,
		},
		{
			name:   "no define call",
			input:  `export class Foo { }`,
			want:   "",
			wantOK: false,
		},
		{
			name:   "no hyphen in name (invalid custom element)",
			input:  `customElements.define('mywidget', MyWidget);`,
			want:   "",
			wantOK: false,
		},
		{
			name:   "empty string",
			input:  "",
			want:   "",
			wantOK: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := discoverTagNameFast(tc.input)
			if ok != tc.wantOK {
				t.Fatalf("ok = %v, want %v", ok, tc.wantOK)
			}
			if got != tc.want {
				t.Errorf("tag = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestDiscoverTagNameFast_RealBundle(t *testing.T) {
	bundle := bundleMyGreeting(t)

	tag, ok := discoverTagNameFast(bundle)
	if !ok {
		t.Fatal("regex failed to find tag name in real esbuild bundle")
	}
	if tag != "my-greeting" {
		t.Errorf("tag = %q, want %q", tag, "my-greeting")
	}
}

func TestDiscoverTagNameFast_DecoratorBundle_MissesVariableDefine(t *testing.T) {
	bundle, err := BundleComponent("../../testdata/sources/my-card.ts")
	if err != nil {
		t.Fatalf("bundling my-card: %v", err)
	}

	// Decorator-compiled bundles use customElements.define(variable, ctor)
	// instead of a string literal, so the regex correctly returns false.
	_, ok := discoverTagNameFast(bundle)
	if ok {
		t.Error("regex should miss decorator bundles that use variable tag names")
	}
}

func TestDiscoverTagName_DecoratorBundle_QJSFallback(t *testing.T) {
	bundle, err := BundleComponent("../../testdata/sources/my-card.ts")
	if err != nil {
		t.Fatalf("bundling my-card: %v", err)
	}

	tag, err := DiscoverTagName(bundle)
	if err != nil {
		t.Fatalf("DiscoverTagName: %v", err)
	}
	if tag != "my-card" {
		t.Errorf("tag = %q, want %q", tag, "my-card")
	}
}

func TestDiscoverTagName_RegexFastPath(t *testing.T) {
	tag, err := discoverTagName(`customElements.define('fast-path', class extends HTMLElement{});`)
	if err != nil {
		t.Fatalf("discoverTagName: %v", err)
	}
	if tag != "fast-path" {
		t.Errorf("tag = %q, want %q", tag, "fast-path")
	}
}

func TestDiscoverTagName_FallsBackToQJS(t *testing.T) {
	bundle := bundleMyGreeting(t)

	tag, err := DiscoverTagName(bundle)
	if err != nil {
		t.Fatalf("DiscoverTagName: %v", err)
	}
	if tag != "my-greeting" {
		t.Errorf("tag = %q, want %q", tag, "my-greeting")
	}
}

func TestDiscoverTagNames_Batch(t *testing.T) {
	greeting := bundleMyGreeting(t)

	card, err := BundleComponent("../../testdata/sources/my-card.ts")
	if err != nil {
		t.Fatalf("bundling my-card: %v", err)
	}

	results, err := discoverTagNames([]string{greeting, card})
	if err != nil {
		t.Fatalf("discoverTagNames: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("len(results) = %d, want 2", len(results))
	}
	if results[0] != "my-greeting" {
		t.Errorf("results[0] = %q, want %q", results[0], "my-greeting")
	}
	if results[1] != "my-card" {
		t.Errorf("results[1] = %q, want %q", results[1], "my-card")
	}
}

func TestDiscoverTagNames_EmptySlice(t *testing.T) {
	results, err := discoverTagNames(nil)
	if err != nil {
		t.Fatalf("discoverTagNames(nil): %v", err)
	}
	if results != nil {
		t.Errorf("expected nil results, got %v", results)
	}
}

func TestDiscoverTagNames_SkipsInvalidBundles(t *testing.T) {
	greeting := bundleMyGreeting(t)

	results, err := discoverTagNames([]string{
		"var x = 42;",
		greeting,
	})
	if err != nil {
		t.Fatalf("discoverTagNames: %v", err)
	}

	if _, ok := results[0]; ok {
		t.Error("expected no result for non-component bundle")
	}
	if results[1] != "my-greeting" {
		t.Errorf("results[1] = %q, want %q", results[1], "my-greeting")
	}
}

func TestRegisterBundles_RegexAndFallback(t *testing.T) {
	greeting := bundleMyGreeting(t)

	card, err := BundleComponent("../../testdata/sources/my-card.ts")
	if err != nil {
		t.Fatalf("bundling my-card: %v", err)
	}

	reg := NewRegistry()
	if err := reg.registerBundles([]string{greeting, card}); err != nil {
		t.Fatalf("registerBundles: %v", err)
	}

	if !reg.Has("my-greeting") {
		t.Error("registry missing my-greeting")
	}
	if !reg.Has("my-card") {
		t.Error("registry missing my-card")
	}
}

func TestLoadDir_BatchDiscovery(t *testing.T) {
	tmp := t.TempDir()

	bundle := bundleMyGreeting(t)
	if err := os.WriteFile(
		filepath.Join(tmp, "my-greeting.golit.bundle.js"),
		[]byte(bundle),
		0644,
	); err != nil {
		t.Fatal(err)
	}

	reg := NewRegistry()
	if err := reg.LoadDir(tmp); err != nil {
		t.Fatalf("LoadDir: %v", err)
	}

	if !reg.Has("my-greeting") {
		t.Error("registry missing my-greeting after LoadDir")
	}
}

func TestLoadDir_SkipsNonBundleFiles(t *testing.T) {
	tmp := t.TempDir()

	os.WriteFile(filepath.Join(tmp, "readme.txt"), []byte("hello"), 0644)
	os.WriteFile(filepath.Join(tmp, "util.js"), []byte("export default 1"), 0644)

	reg := NewRegistry()
	if err := reg.LoadDir(tmp); err != nil {
		t.Fatalf("LoadDir: %v", err)
	}

	if len(reg.TagNames()) != 0 {
		t.Errorf("expected no tags, got %v", reg.TagNames())
	}
}
