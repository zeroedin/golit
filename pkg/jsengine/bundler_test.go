package jsengine

import (
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

func TestDiscoverExternalPackages_NonFatalErrors(t *testing.T) {
	// DiscoverExternalPackages with no valid entry points should return nil, nil
	// (not an error) since there's nothing to discover.
	result, err := DiscoverExternalPackages(nil, "")
	if err != nil {
		t.Fatalf("expected nil error for empty input, got: %v", err)
	}
	if result != nil {
		t.Fatalf("expected nil result for empty input, got: %v", result)
	}

	// With non-existent paths, should also return nil, nil (paths are filtered out)
	result, err = DiscoverExternalPackages([]string{"/nonexistent/path.js"}, "/nonexistent")
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
