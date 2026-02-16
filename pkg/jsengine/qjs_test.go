package jsengine

import (
	"testing"

	"github.com/fastschema/qjs"
)

func TestQJS_BasicEval(t *testing.T) {
	rt, err := qjs.New()
	if err != nil {
		t.Fatalf("creating QJS runtime: %v", err)
	}
	defer rt.Close()

	ctx := rt.Context()
	result, err := ctx.Eval("test.js", qjs.Code(`1 + 2`))
	if err != nil {
		t.Fatalf("eval: %v", err)
	}
	if result.Int32() != 3 {
		t.Errorf("got %d, want 3", result.Int32())
	}
}

func TestQJS_TaggedTemplateLiterals(t *testing.T) {
	rt, err := qjs.New()
	if err != nil {
		t.Fatal(err)
	}
	defer rt.Close()

	ctx := rt.Context()
	result, err := ctx.Eval("test.js", qjs.Code(`
		function html(strings, ...values) {
			return { strings, values, html: strings.reduce((a,s,i) => a + s + (values[i]||''), '') };
		}
		const name = 'World';
		html`+"`<p>${name}</p>`"+`.html;
	`))
	if err != nil {
		t.Fatal(err)
	}
	if result.String() != "<p>World</p>" {
		t.Errorf("got %q, want %q", result.String(), "<p>World</p>")
	}
}

func TestQJS_WeakMapPrivateFields(t *testing.T) {
	rt, err := qjs.New()
	if err != nil {
		t.Fatal(err)
	}
	defer rt.Close()

	ctx := rt.Context()
	result, err := ctx.Eval("test.js", qjs.Code(`
		var _name = new WeakMap();
		class Person {
			constructor(name) { _name.set(this, name); }
			getName() { return _name.get(this); }
		}
		new Person('Alice').getName();
	`))
	if err != nil {
		t.Fatal(err)
	}
	if result.String() != "Alice" {
		t.Errorf("got %q, want %q", result.String(), "Alice")
	}
}
