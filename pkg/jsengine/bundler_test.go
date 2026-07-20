package jsengine

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

func TestExtractDynamicImportTargets(t *testing.T) {
	cases := []struct {
		name    string
		modules map[string]string
		want    []string
	}{
		{
			name: "finds double-quoted dynamic import",
			modules: map[string]string{
				"comp.js": `const { default: { cssText } } = await import("@rhds/tokens/css/default-theme.css.js");`,
			},
			want: []string{"@rhds/tokens/css/default-theme.css.js"},
		},
		{
			name: "finds single-quoted dynamic import",
			modules: map[string]string{
				"comp.js": `await import('@some/pkg/styles.css.js');`,
			},
			want: []string{"@some/pkg/styles.css.js"},
		},
		{
			name: "ignores local imports",
			modules: map[string]string{
				"comp.js": `import("./local.js")`,
			},
			want: nil,
		},
		{
			name: "ignores @golit/runtime",
			modules: map[string]string{
				"comp.js": `import("@golit/runtime")`,
			},
			want: nil,
		},
		{
			name: "deduplicates across modules",
			modules: map[string]string{
				"a.js": `import("@rhds/tokens/css/default-theme.css.js");`,
				"b.js": `import("@rhds/tokens/css/default-theme.css.js");`,
			},
			want: []string{"@rhds/tokens/css/default-theme.css.js"},
		},
		{
			name: "finds multiple targets across lines",
			modules: map[string]string{
				"comp.js": "import(\"@rhds/tokens/css/default-theme.css.js\");\nimport(\"@rhds/icons/ui/check.js\");",
			},
			want: []string{"@rhds/icons/ui/check.js", "@rhds/tokens/css/default-theme.css.js"},
		},
		{
			name: "finds multiple targets on same line",
			modules: map[string]string{
				"comp.js": `import("@rhds/tokens/css/default-theme.css.js"); import("@rhds/icons/ui/check.js");`,
			},
			want: []string{"@rhds/icons/ui/check.js", "@rhds/tokens/css/default-theme.css.js"},
		},
		{
			name: "handles mixed quotes on same line",
			modules: map[string]string{
				"comp.js": `import("@rhds/tokens/css/default-theme.css.js"); import('@rhds/icons/ui/check.js');`,
			},
			want: []string{"@rhds/icons/ui/check.js", "@rhds/tokens/css/default-theme.css.js"},
		},
		{
			name:    "empty modules",
			modules: map[string]string{},
			want:    nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := extractDynamicImportTargets(tc.modules)
			if len(got) != len(tc.want) {
				t.Fatalf("got %d targets %v, want %d targets %v", len(got), got, len(tc.want), tc.want)
			}
			sort.Strings(got)
			sort.Strings(tc.want)
			for i := range got {
				if got[i] != tc.want[i] {
					t.Errorf("target[%d] = %q, want %q", i, got[i], tc.want[i])
				}
			}
		})
	}
}

func TestExtractPackageName(t *testing.T) {
	cases := []struct {
		specifier string
		want      string
	}{
		{"lit", "lit"},
		{"lit/decorators.js", "lit"},
		{"@rhds/tokens", "@rhds/tokens"},
		{"@rhds/tokens/css/default-theme.css.js", "@rhds/tokens"},
		{"prism-esm/components/prism-css.js", "prism-esm"},
		{"@rhds/elements/rh-icon/rh-icon.js", "@rhds/elements"},
		{"@scope", "@scope"},
	}
	for _, tc := range cases {
		t.Run(tc.specifier, func(t *testing.T) {
			got := extractPackageName(tc.specifier)
			if got != tc.want {
				t.Errorf("extractPackageName(%q) = %q, want %q", tc.specifier, got, tc.want)
			}
		})
	}
}

func TestResolveModulePath_SubpathSpecifier(t *testing.T) {
	// This test requires node_modules to exist. Use the hugo-rhds example
	// if available, otherwise skip.
	nmDirs := FindAllNodeModules("../../examples/hugo-rhds/dummy")
	var nmDir string
	if len(nmDirs) > 0 {
		nmDir = nmDirs[0]
	}
	if nmDir == "" {
		t.Skip("node_modules not found in hugo-rhds example")
	}

	resolved, err := ResolveModulePath("@rhds/tokens/css/default-theme.css.js", "../../examples/hugo-rhds")
	if err != nil {
		t.Fatalf("expected subpath to resolve, got: %v", err)
	}
	if !strings.Contains(resolved, "default-theme.css.js") {
		t.Errorf("resolved path %q should contain default-theme.css.js", resolved)
	}

	resolved, err = ResolveModulePath("prism-esm/prism.js", "../../examples/hugo-rhds")
	if err != nil {
		t.Fatalf("expected prism-esm subpath to resolve, got: %v", err)
	}
	if !strings.Contains(resolved, "prism.js") {
		t.Errorf("resolved path %q should contain prism.js", resolved)
	}
}

func TestFindAllNodeModules_NestedDirs(t *testing.T) {
	root := t.TempDir()

	// Simulate:
	//   root/node_modules/@shoelace-style/localize/
	//   root/node_modules/@awesome.me/webawesome/node_modules/@ctrl/tinycolor/
	//   root/node_modules/@awesome.me/webawesome/dist/components/button/button.js
	nested := filepath.Join(root, "node_modules", "@awesome.me", "webawesome", "node_modules", "@ctrl", "tinycolor")
	if err := os.MkdirAll(nested, 0755); err != nil {
		t.Fatal(err)
	}
	localize := filepath.Join(root, "node_modules", "@shoelace-style", "localize")
	if err := os.MkdirAll(localize, 0755); err != nil {
		t.Fatal(err)
	}
	buttonDir := filepath.Join(root, "node_modules", "@awesome.me", "webawesome", "dist", "components", "button")
	if err := os.MkdirAll(buttonDir, 0755); err != nil {
		t.Fatal(err)
	}
	buttonFile := filepath.Join(buttonDir, "button.js")
	if err := os.WriteFile(buttonFile, []byte("// button"), 0644); err != nil {
		t.Fatal(err)
	}

	dirs := findAllNodeModules(buttonFile)

	if len(dirs) != 2 {
		t.Fatalf("expected 2 node_modules dirs, got %d: %v", len(dirs), dirs)
	}

	// First should be the nested one (closer to the source file)
	expectedNested := filepath.Join(root, "node_modules", "@awesome.me", "webawesome", "node_modules")
	expectedRoot := filepath.Join(root, "node_modules")

	if dirs[0] != expectedNested {
		t.Errorf("first dir = %q, want %q", dirs[0], expectedNested)
	}
	if dirs[1] != expectedRoot {
		t.Errorf("second dir = %q, want %q", dirs[1], expectedRoot)
	}
}

func TestFindAllNodeModules_SingleDir(t *testing.T) {
	root := t.TempDir()
	nm := filepath.Join(root, "node_modules")
	if err := os.MkdirAll(nm, 0755); err != nil {
		t.Fatal(err)
	}
	srcDir := filepath.Join(root, "src")
	if err := os.MkdirAll(srcDir, 0755); err != nil {
		t.Fatal(err)
	}

	dirs := findAllNodeModules(filepath.Join(srcDir, "app.js"))
	if len(dirs) != 1 {
		t.Fatalf("expected 1 node_modules dir, got %d: %v", len(dirs), dirs)
	}
	if dirs[0] != nm {
		t.Errorf("dir = %q, want %q", dirs[0], nm)
	}
}

func TestFindAllNodeModules_NoDirs(t *testing.T) {
	root := t.TempDir()
	dirs := findAllNodeModules(filepath.Join(root, "dummy"))
	if len(dirs) != 0 {
		t.Fatalf("expected 0 node_modules dirs, got %d: %v", len(dirs), dirs)
	}
}

func TestDiscoverExternalPackages_NonFatalErrors(t *testing.T) {
	// DiscoverExternalPackages with no valid entry points should return nil, nil
	// (not an error) since there's nothing to discover.
	result, err := DiscoverExternalPackages(nil, nil)
	if err != nil {
		t.Fatalf("expected nil error for empty input, got: %v", err)
	}
	if result != nil {
		t.Fatalf("expected nil result for empty input, got: %v", result)
	}

	// With non-existent paths, should also return nil, nil (paths are filtered out)
	result, err = DiscoverExternalPackages([]string{"/nonexistent/path.js"}, []string{"/nonexistent"})
	if err != nil {
		t.Fatalf("expected nil error for non-existent paths, got: %v", err)
	}
	if result != nil {
		t.Fatalf("expected nil result for non-existent paths, got: %v", result)
	}
}

func TestBundleComponent_MyGreeting(t *testing.T) {
	bundle, err := BundleComponent("../../testdata/sources/my-greeting.js")
	if err != nil {
		t.Fatalf("bundling: %v", err)
	}

	if len(bundle) == 0 {
		t.Fatal("empty bundle")
	}

	t.Logf("Bundle size: %d bytes", len(bundle))

	if !strings.Contains(bundle, "MyGreeting") {
		t.Error("bundle should contain MyGreeting class")
	}
	if !strings.Contains(bundle, "CustomElementRegistry") {
		t.Error("bundle should contain DOM shim")
	}
	if !strings.Contains(bundle, "__collectTemplateResult") {
		t.Error("bundle should contain template collector")
	}
}

func TestBundleComponent_MyCard(t *testing.T) {
	bundle, err := BundleComponent("../../testdata/sources/my-card.ts")
	if err != nil {
		t.Fatalf("bundling: %v", err)
	}

	t.Logf("Bundle size: %d bytes", len(bundle))

	if !strings.Contains(bundle, "MyCard") && !strings.Contains(bundle, "my-card") {
		t.Error("bundle should contain MyCard component")
	}
}
