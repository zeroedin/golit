package jsengine

import (
	"strings"
	"testing"
)

func bundleMyGreeting(t *testing.T) string {
	t.Helper()
	bundle, err := BundleComponent("../../testdata/sources/my-greeting.js")
	if err != nil {
		t.Fatalf("bundling my-greeting: %v", err)
	}
	return bundle
}

func TestEngine_RenderMyGreeting(t *testing.T) {
	bundle := bundleMyGreeting(t)

	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}
	defer engine.Close()

	if err := engine.LoadBundle(bundle); err != nil {
		t.Fatal(err)
	}

	result, err := engine.RenderElement("my-greeting", map[string]string{
		"name": "World",
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("HTML: %s", result.HTML)
	t.Logf("CSS: %s", result.CSS[:min(80, len(result.CSS))])

	if !strings.Contains(result.HTML, "World") {
		t.Error("missing 'World' in output")
	}
	if !strings.Contains(result.HTML, "Hello") {
		t.Error("missing 'Hello' in output")
	}
}

func TestEngine_RenderMyGreeting_DifferentNames(t *testing.T) {
	bundle := bundleMyGreeting(t)

	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}
	defer engine.Close()
	engine.LoadBundle(bundle)

	for _, name := range []string{"Alice", "Go", "Hugo"} {
		t.Run(name, func(t *testing.T) {
			result, err := engine.RenderElement("my-greeting", map[string]string{
				"name": name,
			})
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(result.HTML, name) {
				t.Errorf("missing %q in: %s", name, result.HTML)
			}
		})
	}
}

func TestEngine_StyleExtraction(t *testing.T) {
	bundle := bundleMyGreeting(t)

	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}
	defer engine.Close()
	engine.LoadBundle(bundle)

	result, err := engine.RenderElement("my-greeting", map[string]string{
		"name": "Test",
	})
	if err != nil {
		t.Fatal(err)
	}

	if result.CSS == "" {
		t.Error("expected non-empty CSS from style extraction")
	}
	t.Logf("CSS: %s", result.CSS[:min(100, len(result.CSS))])
}

func TestEngine_UnregisteredElement(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}
	defer engine.Close()

	_, err = engine.RenderElement("unknown-element", map[string]string{})
	if err == nil {
		t.Error("expected error for unregistered element")
	}
}
