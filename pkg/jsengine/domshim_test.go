package jsengine

import (
	_ "embed"
	"testing"

	"github.com/fastschema/qjs"
)

//go:embed domshim.js
var domShimJS string

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
