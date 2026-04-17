package main

import (
	"bytes"
	"io"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/zeroedin/golit/pkg/jsengine"
	"github.com/zeroedin/golit/pkg/transformer"
)

func projectRoot(t *testing.T) string {
	t.Helper()
	dir, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	return dir
}

func buildGolit(t *testing.T) string {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "golit")
	cmd := exec.Command("go", "build", "-o", bin, ".")
	cmd.Dir = filepath.Join(projectRoot(t), "cmd", "golit")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("building golit: %v\n%s", err, out)
	}
	return bin
}

func buildTestBundles(t *testing.T) string {
	t.Helper()
	root := projectRoot(t)
	bundleDir := t.TempDir()

	sources := []string{
		filepath.Join(root, "testdata", "sources", "my-greeting.js"),
	}

	nodeModulesDir := jsengine.FindNodeModules(sources[0])
	externals, err := jsengine.DiscoverExternalPackages(sources, nodeModulesDir)
	if err != nil {
		t.Fatalf("discovering externals: %v", err)
	}

	modules, err := jsengine.BundleComponentModules(sources, jsengine.BundleOptions{
		ExternalPackages: externals,
	})
	if err != nil {
		t.Fatalf("bundling modules: %v", err)
	}

	if nodeModulesDir != "" {
		rt, err := jsengine.BundleSharedRuntime(nodeModulesDir, modules)
		if err != nil {
			t.Fatalf("building shared runtime: %v", err)
		}
		if err := jsengine.SaveBundle(rt, filepath.Join(bundleDir, "_runtime.golit.module.js")); err != nil {
			t.Fatal(err)
		}
	}

	modules = jsengine.RewriteModuleImports(modules, externals)

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

func buildTestPool(t *testing.T) (*jsengine.EnginePool, *jsengine.Registry) {
	t.Helper()
	bundleDir := buildTestBundles(t)

	registry := jsengine.NewRegistry()
	if err := registry.LoadDir(bundleDir); err != nil {
		t.Fatalf("loading bundles: %v", err)
	}

	pool, err := jsengine.NewEnginePool(1)
	if err != nil {
		t.Fatalf("creating pool: %v", err)
	}
	if err := pool.PreloadAll(registry, nil); err != nil {
		t.Fatalf("preloading pool: %v", err)
	}

	return pool, registry
}

func TestRunStdio_SingleRequest(t *testing.T) {
	pool, registry := buildTestPool(t)
	defer pool.Close()

	input := "<my-greeting name=\"Stdio\"></my-greeting>\x00"
	stdin := strings.NewReader(input)
	var stdout bytes.Buffer

	err := runStdio(stdin, &stdout, pool, registry, nil)
	if err != nil {
		t.Fatalf("runStdio: %v", err)
	}

	out := stdout.String()
	if !strings.HasSuffix(out, "\x00") {
		t.Fatal("output should end with NUL")
	}
	result := strings.TrimSuffix(out, "\x00")
	if !strings.Contains(result, `shadowrootmode="open"`) {
		t.Errorf("missing DSD in output:\n%s", result)
	}
	if !strings.Contains(result, "Stdio") {
		t.Errorf("missing name in output:\n%s", result)
	}
}

func TestRunStdio_MultipleRequests(t *testing.T) {
	pool, registry := buildTestPool(t)
	defer pool.Close()

	input := "<my-greeting name=\"First\"></my-greeting>\x00" +
		"<my-greeting name=\"Second\"></my-greeting>\x00"
	stdin := strings.NewReader(input)
	var stdout bytes.Buffer

	err := runStdio(stdin, &stdout, pool, registry, nil)
	if err != nil {
		t.Fatalf("runStdio: %v", err)
	}

	parts := strings.Split(stdout.String(), "\x00")
	if len(parts) < 3 {
		t.Fatalf("expected 2 responses, got output: %q", stdout.String())
	}
	if !strings.Contains(parts[0], "First") {
		t.Errorf("first response missing 'First':\n%s", parts[0])
	}
	if !strings.Contains(parts[1], "Second") {
		t.Errorf("second response missing 'Second':\n%s", parts[1])
	}
}

func TestRunStdio_EmptyInput(t *testing.T) {
	pool, registry := buildTestPool(t)
	defer pool.Close()

	stdin := strings.NewReader("\x00")
	var stdout bytes.Buffer

	err := runStdio(stdin, &stdout, pool, registry, nil)
	if err != nil {
		t.Fatalf("runStdio: %v", err)
	}

	if stdout.String() != "\x00" {
		t.Errorf("expected bare NUL response, got: %q", stdout.String())
	}
}

func TestRunStdio_EOF(t *testing.T) {
	pool, registry := buildTestPool(t)
	defer pool.Close()

	stdin := strings.NewReader("")
	var stdout bytes.Buffer

	err := runStdio(stdin, &stdout, pool, registry, nil)
	if err != nil {
		t.Fatalf("runStdio should return nil on EOF: %v", err)
	}
}

func TestRunStdio_PassthroughHTML(t *testing.T) {
	pool, registry := buildTestPool(t)
	defer pool.Close()

	input := "<div>no custom elements</div>\x00"
	stdin := strings.NewReader(input)
	var stdout bytes.Buffer

	err := runStdio(stdin, &stdout, pool, registry, nil)
	if err != nil {
		t.Fatalf("runStdio: %v", err)
	}

	result := strings.TrimSuffix(stdout.String(), "\x00")
	if !strings.Contains(result, "no custom elements") {
		t.Errorf("passthrough HTML should be preserved:\n%s", result)
	}
}

func TestRunStdio_IgnoredTags(t *testing.T) {
	pool, registry := buildTestPool(t)
	defer pool.Close()

	ignored := map[string]bool{"my-greeting": true}
	input := "<my-greeting name=\"Skip\"></my-greeting>\x00"
	stdin := strings.NewReader(input)
	var stdout bytes.Buffer

	err := runStdio(stdin, &stdout, pool, registry, ignored)
	if err != nil {
		t.Fatalf("runStdio: %v", err)
	}

	result := strings.TrimSuffix(stdout.String(), "\x00")
	if strings.Contains(result, `shadowrootmode`) {
		t.Error("ignored element should not get DSD")
	}
}

func TestRunStdio_LargePayload(t *testing.T) {
	pool, registry := buildTestPool(t)
	defer pool.Close()

	var sb strings.Builder
	for i := 0; i < 100; i++ {
		sb.WriteString(`<my-greeting name="Item"></my-greeting>`)
	}
	sb.WriteByte('\x00')

	stdin := strings.NewReader(sb.String())
	var stdout bytes.Buffer

	err := runStdio(stdin, &stdout, pool, registry, nil)
	if err != nil {
		t.Fatalf("runStdio: %v", err)
	}

	result := strings.TrimSuffix(stdout.String(), "\x00")
	count := strings.Count(result, `shadowrootmode="open"`)
	if count < 100 {
		t.Errorf("expected >= 100 DSD templates, got %d", count)
	}
}

func TestRunStdio_FullDocument(t *testing.T) {
	pool, registry := buildTestPool(t)
	defer pool.Close()

	input := "<!DOCTYPE html><html><head><title>Test</title></head><body><my-greeting name=\"Doc\"></my-greeting></body></html>\x00"
	stdin := strings.NewReader(input)
	var stdout bytes.Buffer

	err := runStdio(stdin, &stdout, pool, registry, nil)
	if err != nil {
		t.Fatalf("runStdio: %v", err)
	}

	result := strings.TrimSuffix(stdout.String(), "\x00")
	if !strings.Contains(result, "<title>Test</title>") {
		t.Errorf("title should be preserved:\n%s", result)
	}
	if !strings.Contains(result, `shadowrootmode="open"`) {
		t.Errorf("missing DSD in output:\n%s", result)
	}
}

func TestRunStdio_BrokenStdout(t *testing.T) {
	pool, registry := buildTestPool(t)
	defer pool.Close()

	input := "<my-greeting name=\"Test\"></my-greeting>\x00"
	stdin := strings.NewReader(input)
	w := &brokenWriter{}

	err := runStdio(stdin, w, pool, registry, nil)
	if err == nil {
		t.Fatal("expected error on broken stdout")
	}
	if !strings.Contains(err.Error(), "stdout") {
		t.Errorf("expected stdout-related error, got: %v", err)
	}
}

type brokenWriter struct{}

func (w *brokenWriter) Write(p []byte) (int, error) {
	return 0, io.ErrClosedPipe
}

func TestRunStdio_PoolReuse(t *testing.T) {
	pool, registry := buildTestPool(t)
	defer pool.Close()

	engine := pool.Get()
	expected, err := transformer.RenderHTMLWithEngine(
		`<my-greeting name="Warm"></my-greeting>`, engine, registry, nil)
	pool.Put(engine)
	if err != nil {
		t.Fatal(err)
	}

	var sb strings.Builder
	for i := 0; i < 3; i++ {
		sb.WriteString("<my-greeting name=\"Warm\"></my-greeting>\x00")
	}
	stdin := strings.NewReader(sb.String())
	var stdout bytes.Buffer

	if err := runStdio(stdin, &stdout, pool, registry, nil); err != nil {
		t.Fatal(err)
	}

	parts := strings.Split(stdout.String(), "\x00")
	for i := 0; i < 3; i++ {
		if parts[i] != expected {
			t.Errorf("request %d output differs from direct render", i)
		}
	}
}

func TestE2E_ServeStdio(t *testing.T) {
	bin := buildGolit(t)
	bundleDir := buildTestBundles(t)

	cmd := exec.Command(bin, "serve", "--defs", bundleDir, "--stdio")
	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		t.Fatal(err)
	}
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	cmd.Stderr = nil

	if err := cmd.Start(); err != nil {
		t.Fatalf("starting golit serve --stdio: %v", err)
	}
	defer func() {
		stdinPipe.Close()
		cmd.Wait()
	}()

	_, err = stdinPipe.Write([]byte("<my-greeting name=\"E2E\"></my-greeting>\x00"))
	if err != nil {
		t.Fatalf("writing to stdin: %v", err)
	}

	result, err := readUntilNUL(t, stdoutPipe)
	if err != nil {
		t.Fatalf("reading response: %v", err)
	}

	if !strings.Contains(result, `shadowrootmode="open"`) {
		t.Errorf("missing DSD in e2e output:\n%s", result)
	}
	if !strings.Contains(result, "E2E") {
		t.Errorf("missing name in e2e output:\n%s", result)
	}

	_, err = stdinPipe.Write([]byte("<my-greeting name=\"Again\"></my-greeting>\x00"))
	if err != nil {
		t.Fatalf("writing second request: %v", err)
	}

	result2, err := readUntilNUL(t, stdoutPipe)
	if err != nil {
		t.Fatalf("reading second response: %v", err)
	}
	if !strings.Contains(result2, "Again") {
		t.Errorf("missing name in second response:\n%s", result2)
	}

	stdinPipe.Close()

	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()

	select {
	case err := <-done:
		if err != nil {
			t.Errorf("process exited with error: %v", err)
		}
	case <-time.After(5 * time.Second):
		cmd.Process.Kill()
		t.Fatal("process did not exit after stdin closed")
	}
}

func TestE2E_ServeStdioMutualExclusion(t *testing.T) {
	bin := buildGolit(t)
	bundleDir := buildTestBundles(t)

	cmd := exec.Command(bin, "serve", "--defs", bundleDir, "--stdio", "--listen", "127.0.0.1:9999")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected error when --stdio and --listen both set")
	}
	if !strings.Contains(string(out), "mutually exclusive") {
		t.Errorf("expected 'mutually exclusive' error, got: %s", out)
	}
}

func readUntilNUL(t *testing.T, r io.Reader) (string, error) {
	t.Helper()
	var buf bytes.Buffer
	b := make([]byte, 1)
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		n, err := r.Read(b)
		if err != nil {
			return buf.String(), err
		}
		if n > 0 {
			if b[0] == '\x00' {
				return buf.String(), nil
			}
			buf.WriteByte(b[0])
		}
	}
	return buf.String(), io.ErrNoProgress
}
