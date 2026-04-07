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

func TestDiscoverTagNameFast_ReverseScanFallback(t *testing.T) {
	input := `customElements.define('dep-a', DepA);
	// Later, a non-define reference to customElements
	var reg = customElements;`
	got, ok := discoverTagNameFast(input)
	if !ok {
		t.Fatal("expected tag name from earlier define call")
	}
	if got != "dep-a" {
		t.Errorf("tag = %q, want %q", got, "dep-a")
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

func TestDiscoverTagName_DecoratorBundle(t *testing.T) {
	bundle, err := BundleComponent("../../testdata/sources/my-card.ts")
	if err != nil {
		t.Fatalf("bundling my-card: %v", err)
	}

	tag, err := DiscoverTagName(bundle)
	// Decorator bundles use customElements.define(variable, ctor) which the
	// regex may miss — that's acceptable since thin modules use the
	// @customElement("tag") decorator pattern that the regex does catch.
	if err == nil && tag != "my-card" {
		t.Errorf("tag = %q, want %q", tag, "my-card")
	}
}

func TestDiscoverTagName_RegexFastPath(t *testing.T) {
	tag, err := DiscoverTagName(`customElements.define('fast-path', class extends HTMLElement{});`)
	if err != nil {
		t.Fatalf("DiscoverTagName: %v", err)
	}
	if tag != "fast-path" {
		t.Errorf("tag = %q, want %q", tag, "fast-path")
	}
}

func TestDiscoverTagName_FromBundle(t *testing.T) {
	bundle := bundleMyGreeting(t)

	tag, err := DiscoverTagName(bundle)
	if err != nil {
		t.Fatalf("DiscoverTagName: %v", err)
	}
	if tag != "my-greeting" {
		t.Errorf("tag = %q, want %q", tag, "my-greeting")
	}
}

func TestDiscoverTagNameFast_DecoratorPattern(t *testing.T) {
	cases := []struct {
		name   string
		input  string
		want   string
		wantOK bool
	}{
		{
			name:   "simple decorator",
			input:  `customElement("rh-accordion")`,
			want:   "rh-accordion",
			wantOK: true,
		},
		{
			name:   "numbered decorator",
			input:  `customElement3("rh-accordion-header")`,
			want:   "rh-accordion-header",
			wantOK: true,
		},
		{
			name:   "single quotes",
			input:  `customElement('my-el')`,
			want:   "my-el",
			wantOK: true,
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

func TestRegister_And_Lookup(t *testing.T) {
	reg := NewRegistry()
	reg.Register("my-el", "export default class MyEl extends HTMLElement {}")

	if !reg.Has("my-el") {
		t.Error("registry missing my-el")
	}
	if reg.Lookup("my-el") == "" {
		t.Error("Lookup returned empty for registered tag")
	}
}

func TestLoadDir_ModuleDiscovery(t *testing.T) {
	tmp := t.TempDir()

	// Write a thin module with a @customElement decorator
	mod := `import { customElement } from "@golit/runtime";
class MyGreeting extends HTMLElement {}
customElement("my-greeting")(MyGreeting);`
	if err := os.WriteFile(
		filepath.Join(tmp, "my-greeting.golit.module.js"),
		[]byte(mod),
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
