package golit_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/sspriggs/golit"
	"github.com/sspriggs/golit/pkg/jsengine"
)

func setupRendererBundles(t *testing.T) string {
	t.Helper()
	bundleDir := t.TempDir()

	for _, src := range []string{
		"testdata/sources/my-greeting.js",
		"testdata/sources/my-card.ts",
	} {
		bundle, err := jsengine.BundleComponent(src)
		if err != nil {
			t.Fatalf("bundling %s: %v", src, err)
		}
		base := filepath.Base(src)
		ext := filepath.Ext(base)
		outName := strings.TrimSuffix(base, ext) + ".golit.bundle.js"
		if err := jsengine.SaveBundle(bundle, filepath.Join(bundleDir, outName)); err != nil {
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
