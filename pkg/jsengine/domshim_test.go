package jsengine

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/fastschema/qjs"
)

func loadShim(t *testing.T) (*qjs.Runtime, *qjs.Context) {
	t.Helper()
	rt, err := qjs.New()
	if err != nil {
		t.Fatal(err)
	}
	ctx := rt.Context()
	if _, err := ctx.Eval("domshim.js", qjs.Code(domShimJS)); err != nil {
		rt.Close()
		t.Fatalf("loading DOM shim: %v", err)
	}
	return rt, ctx
}

func evalJSON(t *testing.T, ctx *qjs.Context, script string) map[string]interface{} {
	t.Helper()
	result, err := ctx.Eval("test.js", qjs.Code(script))
	if err != nil {
		t.Fatalf("eval: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(result.String()), &m); err != nil {
		t.Fatalf("json parse %q: %v", result.String(), err)
	}
	return m
}

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

// --- Event system tests ---

func TestDOMShim_EventProperties(t *testing.T) {
	rt, ctx := loadShim(t)
	defer rt.Close()

	m := evalJSON(t, ctx, `
		const e = new Event('click', { bubbles: true, composed: true, cancelable: true });
		const ce = new CustomEvent('my-event', { detail: { key: 42 }, bubbles: false });
		JSON.stringify({
			type: e.type,
			bubbles: e.bubbles,
			composed: e.composed,
			cancelable: e.cancelable,
			defaultPrevented: e.defaultPrevented,
			isTrusted: e.isTrusted,
			hasTimeStamp: typeof e.timeStamp === 'number',
			hasComposedPath: typeof e.composedPath === 'function',
			hasStopProp: typeof e.stopPropagation === 'function',
			hasStopImmediate: typeof e.stopImmediatePropagation === 'function',
			hasPreventDefault: typeof e.preventDefault === 'function',
			ceType: ce.type,
			ceBubbles: ce.bubbles,
			ceDetail: ce.detail.key,
			phaseNone: Event.NONE,
			phaseCapture: Event.CAPTURING_PHASE,
			phaseTarget: Event.AT_TARGET,
			phaseBubble: Event.BUBBLING_PHASE,
		});
	`)

	checks := map[string]interface{}{
		"type": "click", "bubbles": true, "composed": true, "cancelable": true,
		"defaultPrevented": false, "isTrusted": false,
		"hasTimeStamp": true, "hasComposedPath": true, "hasStopProp": true,
		"hasStopImmediate": true, "hasPreventDefault": true,
		"ceType": "my-event", "ceBubbles": false, "ceDetail": float64(42),
		"phaseNone": float64(0), "phaseCapture": float64(1),
		"phaseTarget": float64(2), "phaseBubble": float64(3),
	}
	for k, want := range checks {
		if m[k] != want {
			t.Errorf("%s: got %v (%T), want %v (%T)", k, m[k], m[k], want, want)
		}
	}
}

func TestDOMShim_ComposedPathEmpty(t *testing.T) {
	rt, ctx := loadShim(t)
	defer rt.Close()

	m := evalJSON(t, ctx, `
		const e = new Event('test');
		JSON.stringify({ path: e.composedPath(), len: e.composedPath().length });
	`)

	if length, ok := m["len"].(float64); !ok || length != 0 {
		t.Errorf("composedPath before dispatch should be empty, got len=%v", m["len"])
	}
}

func TestDOMShim_DispatchEventBasic(t *testing.T) {
	rt, ctx := loadShim(t)
	defer rt.Close()

	m := evalJSON(t, ctx, `
		const el = new HTMLElement();
		let captured = {};
		el.addEventListener('test', (e) => {
			captured.type = e.type;
			captured.hasTarget = e.target === el;
			captured.hasCurrent = e.currentTarget === el;
			captured.pathLen = e.composedPath().length;
			captured.pathHasEl = e.composedPath()[0] === el;
		});
		const returned = el.dispatchEvent(new Event('test'));
		captured.returned = returned;
		JSON.stringify(captured);
	`)

	checks := map[string]interface{}{
		"type": "test", "hasTarget": true, "hasCurrent": true,
		"pathLen": float64(1), "pathHasEl": true, "returned": true,
	}
	for k, want := range checks {
		if m[k] != want {
			t.Errorf("%s: got %v, want %v", k, m[k], want)
		}
	}
}

func TestDOMShim_StopPropagation(t *testing.T) {
	rt, ctx := loadShim(t)
	defer rt.Close()

	m := evalJSON(t, ctx, `
		const child = new HTMLElement();
		const parent = new HTMLElement();
		child.__eventTargetParent = parent;

		let calls = [];
		child.addEventListener('test', (e) => {
			calls.push('child');
			e.stopPropagation();
		});
		parent.addEventListener('test', () => calls.push('parent'));

		child.dispatchEvent(new Event('test', { bubbles: true, composed: true }));
		JSON.stringify({ calls: calls });
	`)

	callsRaw, ok := m["calls"].([]interface{})
	if !ok || len(callsRaw) != 1 || callsRaw[0] != "child" {
		t.Errorf("stopPropagation should prevent parent: got %v", m["calls"])
	}
}

func TestDOMShim_StopImmediatePropagation(t *testing.T) {
	rt, ctx := loadShim(t)
	defer rt.Close()

	m := evalJSON(t, ctx, `
		const el = new HTMLElement();
		let calls = [];
		el.addEventListener('test', (e) => {
			calls.push('first');
			e.stopImmediatePropagation();
		});
		el.addEventListener('test', () => calls.push('second'));

		el.dispatchEvent(new Event('test'));
		JSON.stringify({ calls: calls });
	`)

	callsRaw, ok := m["calls"].([]interface{})
	if !ok || len(callsRaw) != 1 || callsRaw[0] != "first" {
		t.Errorf("stopImmediatePropagation should prevent second listener: got %v", m["calls"])
	}
}

func TestDOMShim_PreventDefault(t *testing.T) {
	rt, ctx := loadShim(t)
	defer rt.Close()

	m := evalJSON(t, ctx, `
		const el = new HTMLElement();
		el.addEventListener('test', (e) => e.preventDefault());
		const returned = el.dispatchEvent(new Event('test', { cancelable: true }));
		JSON.stringify({ returned: returned });
	`)

	if m["returned"] != false {
		t.Errorf("dispatchEvent should return false after preventDefault, got %v", m["returned"])
	}
}

func TestDOMShim_CaptureVsBubbleOrder(t *testing.T) {
	rt, ctx := loadShim(t)
	defer rt.Close()

	m := evalJSON(t, ctx, `
		const el = new HTMLElement();
		let order = [];
		el.addEventListener('test', () => order.push('bubble'));
		el.addEventListener('test', () => order.push('capture'), true);

		el.dispatchEvent(new Event('test'));
		JSON.stringify({ order: order });
	`)

	orderRaw, ok := m["order"].([]interface{})
	if !ok || len(orderRaw) != 2 {
		t.Fatalf("expected 2 calls, got %v", m["order"])
	}
	if orderRaw[0] != "capture" || orderRaw[1] != "bubble" {
		t.Errorf("capture should fire before bubble: got %v", orderRaw)
	}
}

func TestDOMShim_OnceListener(t *testing.T) {
	rt, ctx := loadShim(t)
	defer rt.Close()

	m := evalJSON(t, ctx, `
		const el = new HTMLElement();
		let count = 0;
		el.addEventListener('test', () => count++, { once: true });
		el.dispatchEvent(new Event('test'));
		el.dispatchEvent(new Event('test'));
		JSON.stringify({ count: count });
	`)

	if m["count"] != float64(1) {
		t.Errorf("once listener should fire only once, got count=%v", m["count"])
	}
}

func TestDOMShim_ComposedPathParentChain(t *testing.T) {
	rt, ctx := loadShim(t)
	defer rt.Close()

	m := evalJSON(t, ctx, `
		const child = new HTMLElement();
		const parent = new HTMLElement();
		const root = new HTMLElement();
		child.__eventTargetParent = parent;
		parent.__eventTargetParent = root;

		let pathLen = 0;
		child.addEventListener('test', (e) => { pathLen = e.composedPath().length; });

		child.dispatchEvent(new Event('test', { composed: true }));
		JSON.stringify({ pathLen: pathLen });
	`)

	if m["pathLen"] != float64(3) {
		t.Errorf("composed path should have 3 elements (child->parent->root), got %v", m["pathLen"])
	}
}

func TestDOMShim_ContextRequestPattern(t *testing.T) {
	rt, ctx := loadShim(t)
	defer rt.Close()

	m := evalJSON(t, ctx, `
		const provider = new HTMLElement();
		let result = { composedPathCalled: false, noError: true, targetMatch: false };

		provider.addEventListener('context-request', (e) => {
			try {
				const path = e.composedPath();
				result.composedPathCalled = true;
				const target = e.contextTarget ?? path[0];
				result.targetMatch = target === provider;
				e.stopPropagation();
			} catch (err) {
				result.noError = false;
				result.error = err.message;
			}
		});

		provider.dispatchEvent(new Event('context-request', { bubbles: true, composed: true }));
		JSON.stringify(result);
	`)

	if m["composedPathCalled"] != true {
		t.Errorf("composedPath should have been called")
	}
	if m["noError"] != true {
		t.Errorf("should not throw, got error: %v", m["error"])
	}
	if m["targetMatch"] != true {
		t.Errorf("target from composedPath should match provider")
	}
}

func TestDOMShim_LitServerRoot(t *testing.T) {
	rt, ctx := loadShim(t)
	defer rt.Close()

	m := evalJSON(t, ctx, `
		let called = false;
		globalThis.litServerRoot.addEventListener('context-provider', () => { called = true; });
		globalThis.litServerRoot.dispatchEvent(new Event('context-provider', { bubbles: true, composed: true }));
		JSON.stringify({
			isEventTarget: typeof globalThis.litServerRoot.addEventListener === 'function',
			localName: globalThis.litServerRoot.localName,
			called: called,
		});
	`)

	if m["isEventTarget"] != true {
		t.Errorf("litServerRoot should be a proper EventTarget")
	}
	if m["localName"] != "lit-server-root" {
		t.Errorf("litServerRoot.localName should be 'lit-server-root', got %v", m["localName"])
	}
	if m["called"] != true {
		t.Errorf("litServerRoot should dispatch events to registered listeners")
	}
}

func TestDOMShim_DuplicateDefineIgnored(t *testing.T) {
	rt, ctx := loadShim(t)
	defer rt.Close()

	m := evalJSON(t, ctx, `
		class FirstEl extends HTMLElement { static get tag() { return 'first'; } }
		class SecondEl extends HTMLElement { static get tag() { return 'second'; } }
		customElements.define('dup-el', FirstEl);
		customElements.define('dup-el', SecondEl);
		const Ctor = customElements.get('dup-el');
		const el = new Ctor();
		JSON.stringify({
			kept: Ctor.tag,
			isFirst: Ctor === FirstEl,
		});
	`)

	if m["kept"] != "first" {
		t.Errorf("duplicate define should keep first class, got tag=%v", m["kept"])
	}
	if m["isFirst"] != true {
		t.Errorf("duplicate define should not overwrite: isFirst=%v", m["isFirst"])
	}
}
