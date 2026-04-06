package jsengine

import (
	"strings"
	"testing"

	"github.com/fastschema/qjs"
)

func TestDOMShim_CustomElements(t *testing.T) {
	rt, err := qjs.New()
	if err != nil {
		t.Fatal(err)
	}
	defer rt.Close()

	ctx := rt.Context()

	// Load DOM shim
	_, err = ctx.Eval("domshim.js", qjs.Code(domShimJS))
	if err != nil {
		t.Fatalf("loading DOM shim: %v", err)
	}

	// Define and use a custom element
	result, err := ctx.Eval("test.js", qjs.Code(`
		class MyElement extends HTMLElement {
			constructor() {
				super();
				this.greeting = 'Hello';
			}
		}
		customElements.define('my-element', MyElement);

		const Ctor = customElements.get('my-element');
		const el = new Ctor();
		el.setAttribute('name', 'World');

		JSON.stringify({
			tagName: el.localName,
			name: el.getAttribute('name'),
			greeting: el.greeting,
			hasName: el.hasAttribute('name'),
		});
	`))
	if err != nil {
		t.Fatalf("eval: %v", err)
	}

	t.Logf("Result: %s", result.String())
	expected := `{"tagName":"my-element","name":"World","greeting":"Hello","hasName":true}`
	if result.String() != expected {
		t.Errorf("got %s, want %s", result.String(), expected)
	}
}

func TestDOMShim_MultipleLoadsPreserveRegistry(t *testing.T) {
	rt, err := qjs.New()
	if err != nil {
		t.Fatal(err)
	}
	defer rt.Close()

	ctx := rt.Context()

	// Each bundle wraps the domshim in a function scope (via esbuild),
	// so we wrap in an IIFE to match real bundle structure.
	wrapped := "(function(){" + domShimJS + "})();"

	_, err = ctx.Eval("bundle1.js", qjs.Code(wrapped))
	if err != nil {
		t.Fatalf("first bundle load: %v", err)
	}

	_, err = ctx.Eval("define.js", qjs.Code(`
		class FirstEl extends HTMLElement {}
		customElements.define('first-el', FirstEl);
	`))
	if err != nil {
		t.Fatalf("defining first-el: %v", err)
	}

	// Load a second bundle with its own domshim copy
	_, err = ctx.Eval("bundle2.js", qjs.Code(wrapped))
	if err != nil {
		t.Fatalf("second bundle load: %v", err)
	}

	result, err := ctx.Eval("check.js", qjs.Code(
		`customElements.get('first-el') !== undefined ? 'true' : 'false'`))
	if err != nil {
		t.Fatalf("checking registration: %v", err)
	}
	if result.String() != "true" {
		t.Error("first-el should still be registered after loading domshim a second time")
	}
}

func TestDOMShim_BareLocationBinding(t *testing.T) {
	rt, err := qjs.New()
	if err != nil {
		t.Fatal(err)
	}
	defer rt.Close()

	ctx := rt.Context()
	_, err = ctx.Eval("pre.js", qjs.Code(`globalThis.__golitLocationHref='http://localhost/';`))
	if err != nil {
		t.Fatal(err)
	}
	_, err = ctx.Eval("domshim.js", qjs.Code(domShimJS))
	if err != nil {
		t.Fatalf("loading DOM shim: %v", err)
	}

	result, err := ctx.Eval("test.js", qjs.Code(
		`JSON.stringify({ gt: typeof globalThis.location, bare: typeof location })`))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result.String())
	if strings.Contains(result.String(), `"bare":"undefined"`) {
		t.Fatalf("bare location should resolve on global: %s", result.String())
	}
}

func TestDOMShim_ShadowRoot(t *testing.T) {
	rt, err := qjs.New()
	if err != nil {
		t.Fatal(err)
	}
	defer rt.Close()

	ctx := rt.Context()
	_, err = ctx.Eval("domshim.js", qjs.Code(domShimJS))
	if err != nil {
		t.Fatal(err)
	}

	result, err := ctx.Eval("test.js", qjs.Code(`
		class MyEl extends HTMLElement {
			constructor() {
				super();
				this.attachShadow({ mode: 'open' });
			}
		}
		customElements.define('my-el', MyEl);
		const el = new (customElements.get('my-el'))();
		JSON.stringify({
			hasShadowRoot: el.shadowRoot !== null,
			hostIsElement: el.shadowRoot.host === el,
		});
	`))
	if err != nil {
		t.Fatalf("eval: %v", err)
	}

	t.Logf("Result: %s", result.String())
	if result.String() != `{"hasShadowRoot":true,"hostIsElement":true}` {
		t.Errorf("unexpected result: %s", result.String())
	}
}
