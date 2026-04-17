package main

import (
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/zeroedin/golit/pkg/jsengine"
)

func buildGolit(t *testing.T) string {
	t.Helper()
	name := "golit"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	bin := filepath.Join(t.TempDir(), name)
	cmd := exec.Command("go", "build", "-o", bin, ".")
	cmd.Dir = filepath.Join(projectRoot(t), "cmd", "golit")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("building golit: %v\n%s", err, out)
	}
	return bin
}

func projectRoot(t *testing.T) string {
	t.Helper()
	dir, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	return dir
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

func TestRender_StdinPipe(t *testing.T) {
	bin := buildGolit(t)
	bundleDir := buildTestBundles(t)

	cmd := exec.Command(bin, "render", "--defs", bundleDir)
	cmd.Stdin = strings.NewReader(`<my-greeting name="Stdin"></my-greeting>`)
	out, err := cmd.Output()
	if err != nil {
		stderr := ""
		if ee, ok := err.(*exec.ExitError); ok {
			stderr = string(ee.Stderr)
		}
		t.Fatalf("render via stdin failed: %v\nstderr: %s", err, stderr)
	}

	output := string(out)
	if !strings.Contains(output, `shadowrootmode="open"`) {
		t.Errorf("missing DSD in stdin output:\n%s", output)
	}
	if !strings.Contains(output, "Stdin") {
		t.Errorf("missing name in stdin output:\n%s", output)
	}
}

func TestRender_StdinLargeFragment(t *testing.T) {
	bin := buildGolit(t)
	bundleDir := buildTestBundles(t)

	var sb strings.Builder
	for i := 0; i < 100; i++ {
		sb.WriteString(`<my-greeting name="Item"></my-greeting>`)
	}
	largeInput := sb.String()

	cmd := exec.Command(bin, "render", "--defs", bundleDir)
	cmd.Stdin = strings.NewReader(largeInput)
	out, err := cmd.Output()
	if err != nil {
		stderr := ""
		if ee, ok := err.(*exec.ExitError); ok {
			stderr = string(ee.Stderr)
		}
		t.Fatalf("render large stdin failed: %v\nstderr: %s", err, stderr)
	}

	output := string(out)
	count := strings.Count(output, `shadowrootmode="open"`)
	if count < 100 {
		t.Errorf("expected >= 100 DSD templates, got %d", count)
	}
}

func TestRender_ArgTakesPrecedenceOverStdin(t *testing.T) {
	bin := buildGolit(t)
	bundleDir := buildTestBundles(t)

	cmd := exec.Command(bin, "render", "--defs", bundleDir, `<my-greeting name="Arg"></my-greeting>`)
	cmd.Stdin = strings.NewReader(`<my-greeting name="Stdin"></my-greeting>`)
	out, err := cmd.Output()
	if err != nil {
		stderr := ""
		if ee, ok := err.(*exec.ExitError); ok {
			stderr = string(ee.Stderr)
		}
		t.Fatalf("render with arg+stdin failed: %v\nstderr: %s", err, stderr)
	}

	output := string(out)
	if !strings.Contains(output, "Arg") {
		t.Errorf("arg should take precedence, but 'Arg' not in output:\n%s", output)
	}
	if strings.Contains(output, "Stdin") {
		t.Errorf("stdin should be ignored when arg is provided, but 'Stdin' found in output:\n%s", output)
	}
}

func TestRender_NoInputError(t *testing.T) {
	bin := buildGolit(t)
	bundleDir := buildTestBundles(t)

	cmd := exec.Command(bin, "render", "--defs", bundleDir)
	cmd.Stdin = strings.NewReader("")
	var stderrBuf strings.Builder
	cmd.Stderr = &stderrBuf
	err := cmd.Run()
	if err == nil {
		t.Fatal("expected error when no fragment and no stdin pipe")
	}
	if !strings.Contains(stderrBuf.String(), "missing HTML fragment") {
		t.Errorf("expected 'missing HTML fragment' error, got: %s", stderrBuf.String())
	}
}

func TestRender_StdinWithWhitespace(t *testing.T) {
	bin := buildGolit(t)
	bundleDir := buildTestBundles(t)

	cmd := exec.Command(bin, "render", "--defs", bundleDir)
	cmd.Stdin = strings.NewReader("\n  <my-greeting name=\"Trimmed\"></my-greeting>\n\n")
	out, err := cmd.Output()
	if err != nil {
		stderr := ""
		if ee, ok := err.(*exec.ExitError); ok {
			stderr = string(ee.Stderr)
		}
		t.Fatalf("render stdin with whitespace failed: %v\nstderr: %s", err, stderr)
	}

	output := string(out)
	if !strings.Contains(output, `shadowrootmode="open"`) {
		t.Errorf("missing DSD in trimmed stdin output:\n%s", output)
	}
	if !strings.Contains(output, "Trimmed") {
		t.Errorf("missing name in trimmed stdin output:\n%s", output)
	}
}
