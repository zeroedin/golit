package golit_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/zeroedin/golit"
	"github.com/zeroedin/golit/pkg/jsengine"
)

func setupRendererBundles(t *testing.T) string {
	t.Helper()
	bundleDir := t.TempDir()

	sources := []string{
		"testdata/sources/my-greeting.js",
		"testdata/sources/my-card.ts",
	}

	modules, err := jsengine.BundleComponentModules(sources)
	if err != nil {
		t.Fatalf("bundling modules: %v", err)
	}

	nodeModulesDir := jsengine.FindNodeModules(sources[0])
	if nodeModulesDir != "" {
		rt, err := jsengine.BundleSharedRuntime(nodeModulesDir, modules)
		if err != nil {
			t.Fatalf("building shared runtime: %v", err)
		}
		if err := jsengine.SaveBundle(rt, filepath.Join(bundleDir, "_runtime.golit.module.js")); err != nil {
			t.Fatal(err)
		}
	}

	modules = jsengine.RewriteModuleImports(modules)

	for srcPath, mod := range modules {
		base := filepath.Base(srcPath)
		ext := filepath.Ext(base)
		outName := strings.TrimSuffix(base, ext) + ".golit.module.js"
		if err := jsengine.SaveBundle(mod, filepath.Join(bundleDir, outName)); err != nil {
			t.Fatal(err)
		}
	}

	return bundleDir
}

func TestRenderer_RenderFragment(t *testing.T) {
	bundleDir := setupRendererBundles(t)

	renderer, err := golit.NewRenderer(golit.RendererOptions{
		DefsDir: bundleDir,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer renderer.Close()

	output, err := renderer.RenderFragment(`<my-greeting name="Library"></my-greeting>`)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(output, `shadowrootmode="open"`) {
		t.Error("missing DSD in output")
	}
	if !strings.Contains(output, "Library") {
		t.Error("missing name in output")
	}
}

func TestRenderer_RenderHTML(t *testing.T) {
	bundleDir := setupRendererBundles(t)

	renderer, err := golit.NewRenderer(golit.RendererOptions{
		DefsDir: bundleDir,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer renderer.Close()

	input := `<!DOCTYPE html><html><head><title>API</title></head><body><my-greeting name="Test"></my-greeting></body></html>`
	output, err := renderer.RenderHTML(input)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(output, `<title>API</title>`) {
		t.Error("title should be preserved")
	}
	if !strings.Contains(output, `shadowrootmode`) {
		t.Error("missing DSD in output")
	}
}

func TestRenderer_Ignored(t *testing.T) {
	bundleDir := setupRendererBundles(t)

	renderer, err := golit.NewRenderer(golit.RendererOptions{
		DefsDir: bundleDir,
		Ignored: []string{"my-greeting"},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer renderer.Close()

	output, err := renderer.RenderFragment(`<my-greeting name="Skip"></my-greeting>`)
	if err != nil {
		t.Fatal(err)
	}

	if strings.Contains(output, `shadowrootmode`) {
		t.Error("ignored element should not get DSD")
	}
}
