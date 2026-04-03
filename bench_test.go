package golit_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sspriggs/golit/pkg/jsengine"
	"github.com/sspriggs/golit/pkg/transformer"
)

// setupBenchBundles creates pre-bundled component files and returns the bundle dir.
func setupBenchBundles(b *testing.B) string {
	b.Helper()
	bundleDir := b.TempDir()

	for _, src := range []string{
		"testdata/sources/my-greeting.js",
		"testdata/sources/my-card.ts",
	} {
		bundle, err := jsengine.BundleComponent(src)
		if err != nil {
			b.Fatalf("bundling %s: %v", src, err)
		}
		base := filepath.Base(src)
		ext := filepath.Ext(base)
		outName := strings.TrimSuffix(base, ext) + ".golit.bundle.js"
		if err := jsengine.SaveBundle(bundle, filepath.Join(bundleDir, outName)); err != nil {
			b.Fatal(err)
		}
	}

	return bundleDir
}

// setupBenchHTML creates a directory with N HTML files, each containing
// several custom elements to render.
func setupBenchHTML(b *testing.B, n int) string {
	b.Helper()
	dir := b.TempDir()

	for i := 0; i < n; i++ {
		content := fmt.Sprintf(`<!DOCTYPE html>
<html><head><title>Page %d</title></head>
<body>
  <my-greeting name="User%d"></my-greeting>
  <my-card subtitle="Card %d">
    <p>Content for page %d</p>
  </my-card>
  <my-greeting name="Footer%d"></my-greeting>
</body></html>`, i, i, i, i, i)
		path := filepath.Join(dir, fmt.Sprintf("page-%03d.html", i))
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			b.Fatal(err)
		}
	}

	return dir
}

func benchmarkTransformDir(b *testing.B, fileCount, concurrency int) {
	bundleDir := setupBenchBundles(b)
	htmlDir := setupBenchHTML(b, fileCount)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Copy HTML files fresh each iteration (transform modifies in-place)
		workDir := b.TempDir()
		entries, _ := os.ReadDir(htmlDir)
		for _, e := range entries {
			data, _ := os.ReadFile(filepath.Join(htmlDir, e.Name()))
			os.WriteFile(filepath.Join(workDir, e.Name()), data, 0644)
		}

		_, err := transformer.TransformDir(workDir, transformer.Options{
			DefsDir:     bundleDir,
			Concurrency: concurrency,
		})
		if err != nil {
			b.Fatal(err)
		}
	}
}

// 10 files
func BenchmarkTransformDir_10files_seq(b *testing.B)  { benchmarkTransformDir(b, 10, 1) }
func BenchmarkTransformDir_10files_par(b *testing.B)  { benchmarkTransformDir(b, 10, 0) }

// 50 files
func BenchmarkTransformDir_50files_seq(b *testing.B)  { benchmarkTransformDir(b, 50, 1) }
func BenchmarkTransformDir_50files_par(b *testing.B)  { benchmarkTransformDir(b, 50, 0) }

// 100 files
func BenchmarkTransformDir_100files_seq(b *testing.B) { benchmarkTransformDir(b, 100, 1) }
func BenchmarkTransformDir_100files_par(b *testing.B) { benchmarkTransformDir(b, 100, 0) }

// 200 files
func BenchmarkTransformDir_200files_seq(b *testing.B) { benchmarkTransformDir(b, 200, 1) }
func BenchmarkTransformDir_200files_par(b *testing.B) { benchmarkTransformDir(b, 200, 0) }

// 4 workers (lower overhead than 12)
func BenchmarkTransformDir_100files_4w(b *testing.B)  { benchmarkTransformDir(b, 100, 4) }
func BenchmarkTransformDir_50files_4w(b *testing.B)   { benchmarkTransformDir(b, 50, 4) }

// Engine pool creation cost
func BenchmarkEnginePoolCreate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		pool, err := jsengine.NewEnginePool(4)
		if err != nil {
			b.Fatal(err)
		}
		pool.Close()
	}
}

func BenchmarkSingleEngineCreate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		e, err := jsengine.NewEngine()
		if err != nil {
			b.Fatal(err)
		}
		e.Close()
	}
}
