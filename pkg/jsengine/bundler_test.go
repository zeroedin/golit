package jsengine

import (
	"strings"
	"testing"
)

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
