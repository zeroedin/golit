package golit_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/zeroedin/golit/pkg/jsengine"
	"github.com/zeroedin/golit/pkg/transformer"
)

func setupBundles(t *testing.T) string {
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

func TestE2E_RenderFragment_MyGreeting(t *testing.T) {
	bundleDir := setupBundles(t)
	registry := jsengine.NewRegistry()
	if err := registry.LoadDir(bundleDir); err != nil {
		t.Fatal(err)
	}

	output, err := transformer.RenderFragment(
		`<my-greeting name="Hugo"></my-greeting>`,
		registry,
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Output:\n%s", output)
	assertContains(t, output, `shadowrootmode="open"`)
	assertContains(t, output, `<style>`)
	assertContains(t, output, `Hugo`)
	assertDSDCount(t, output, 1)
}

func TestE2E_RenderHTML_FullDocument(t *testing.T) {
	bundleDir := setupBundles(t)
	registry := jsengine.NewRegistry()
	registry.LoadDir(bundleDir)

	input := `<!DOCTYPE html><html><head><title>Test</title></head><body><my-greeting name="Go"></my-greeting></body></html>`
	output, err := transformer.RenderHTML(input, registry)
	if err != nil {
		t.Fatal(err)
	}

	assertContains(t, output, `<title>Test</title>`)
	assertContains(t, output, `Go`)
	assertContains(t, output, `shadowrootmode`)
	assertDSDCount(t, output, 1)
}

func TestE2E_TransformDir(t *testing.T) {
	bundleDir := setupBundles(t)

	htmlDir := t.TempDir()
	os.WriteFile(filepath.Join(htmlDir, "index.html"),
		[]byte(`<!DOCTYPE html><html><head></head><body><my-greeting name="Test"></my-greeting></body></html>`), 0644)

	result, err := transformer.TransformDir(htmlDir, transformer.Options{DefsDir: bundleDir})
	if err != nil {
		t.Fatal(err)
	}

	if result.FilesProcessed != 1 || result.FilesModified != 1 {
		t.Errorf("expected 1 processed/modified, got %d/%d", result.FilesProcessed, result.FilesModified)
	}

	data, _ := os.ReadFile(filepath.Join(htmlDir, "index.html"))
	assertContains(t, string(data), "shadowrootmode")
}

func TestE2E_TransformDir_OutDir(t *testing.T) {
	bundleDir := setupBundles(t)
	htmlDir := t.TempDir()
	outDir := t.TempDir()

	os.WriteFile(filepath.Join(htmlDir, "test.html"),
		[]byte(`<!DOCTYPE html><html><head></head><body><my-greeting name="X"></my-greeting></body></html>`), 0644)

	_, err := transformer.TransformDir(htmlDir, transformer.Options{DefsDir: bundleDir, OutDir: outDir})
	if err != nil {
		t.Fatal(err)
	}

	orig, _ := os.ReadFile(filepath.Join(htmlDir, "test.html"))
	if strings.Contains(string(orig), "shadowrootmode") {
		t.Error("original should not be modified")
	}

	out, _ := os.ReadFile(filepath.Join(outDir, "test.html"))
	assertContains(t, string(out), "shadowrootmode")
}

func TestE2E_UnregisteredElements(t *testing.T) {
	bundleDir := setupBundles(t)
	registry := jsengine.NewRegistry()
	registry.LoadDir(bundleDir)

	output, err := transformer.RenderFragment(
		`<my-greeting name="A"></my-greeting><unknown-el>hi</unknown-el>`,
		registry,
	)
	if err != nil {
		t.Fatal(err)
	}

	assertContains(t, output, `shadowrootmode`)
	assertContains(t, output, `<unknown-el>hi</unknown-el>`)

	found := false
	for _, tag := range registry.Unregistered() {
		if tag == "unknown-el" {
			found = true
		}
	}
	if !found {
		t.Error("unknown-el should be in unregistered list")
	}
}

func TestE2E_NestedComponents(t *testing.T) {
	bundleDir := setupBundles(t)
	registry := jsengine.NewRegistry()
	if err := registry.LoadDir(bundleDir); err != nil {
		t.Fatal(err)
	}

	output, err := transformer.RenderFragment(
		`<my-card subtitle="Nested"><my-greeting name="Inner"></my-greeting></my-card>`,
		registry,
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Output:\n%s", output)
	assertDSDCount(t, output, 2)
	assertContains(t, output, `Inner`)
	assertContains(t, output, `Nested`)
}

func TestE2E_Idempotent(t *testing.T) {
	bundleDir := setupBundles(t)
	registry := jsengine.NewRegistry()
	registry.LoadDir(bundleDir)

	input := `<my-greeting name="World"></my-greeting>`
	output1, err := transformer.RenderFragment(input, registry)
	if err != nil {
		t.Fatal(err)
	}

	registry2 := jsengine.NewRegistry()
	registry2.LoadDir(bundleDir)
	output2, err := transformer.RenderFragment(output1, registry2)
	if err != nil {
		t.Fatal(err)
	}

	if output1 != output2 {
		t.Errorf("not idempotent:\nFirst: %s\nSecond: %s", output1[:200], output2[:200])
	}
}

func assertContains(t *testing.T, output, substr string) {
	t.Helper()
	if !strings.Contains(output, substr) {
		t.Errorf("output does not contain %q\n\nFull output:\n%s", substr, output)
	}
}

func countDSD(output string) int {
	return strings.Count(output, `shadowrootmode="open"`)
}

func assertDSDCount(t *testing.T, output string, atLeast int) {
	t.Helper()
	n := countDSD(output)
	if n < atLeast {
		t.Errorf("expected >= %d DSD templates, got %d\n\nFull output:\n%s", atLeast, n, output)
	}
}
