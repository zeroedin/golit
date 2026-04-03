package transformer

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"

	"github.com/zeroedin/golit/pkg/jsengine"
)

// mockEngine implements ElementRenderer for testing without QJS.
type mockEngine struct {
	bundles    map[string]bool
	results    map[string]jsengine.BatchResult
	registered map[string]bool
}

func (m *mockEngine) LoadBundleForTag(tagName string, _ *jsengine.Registry) (bool, error) {
	if m.bundles[tagName] {
		return true, nil
	}
	return false, nil
}

func (m *mockEngine) RenderBatch(requests []jsengine.BatchRequest) ([]jsengine.BatchResult, error) {
	var out []jsengine.BatchResult
	for _, req := range requests {
		if r, ok := m.results[req.TagName]; ok {
			r.ID = req.ID
			r.TagName = req.TagName
			out = append(out, r)
		}
	}
	return out, nil
}

func (m *mockEngine) IsRegistered(tagName string) bool {
	return m.registered[tagName]
}

func newMockEngine() *mockEngine {
	return &mockEngine{
		bundles:    make(map[string]bool),
		results:    make(map[string]jsengine.BatchResult),
		registered: make(map[string]bool),
	}
}

func (m *mockEngine) addComponent(tag, renderedHTML, css string) {
	m.bundles[tag] = true
	m.results[tag] = jsengine.BatchResult{
		HTML: renderedHTML,
		CSS:  css,
	}
}

func (m *mockEngine) addComponentError(tag, errMsg string) {
	m.bundles[tag] = true
	m.results[tag] = jsengine.BatchResult{
		Error: errMsg,
	}
}

func TestRenderFragment_SingleElement(t *testing.T) {
	engine := newMockEngine()
	engine.addComponent("my-el", "<p>hello</p>", ":host{display:block}")

	registry := jsengine.NewRegistry()
	registry.Register("my-el", "fake-bundle")

	output, err := RenderFragmentWithEngine(`<my-el></my-el>`, engine, registry, nil)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(output, `shadowrootmode="open"`) {
		t.Error("missing shadowrootmode attribute")
	}
	if !strings.Contains(output, "<p>hello</p>") {
		t.Error("missing rendered HTML in shadow root")
	}
	if !strings.Contains(output, "<style>:host{display:block}</style>") {
		t.Error("missing CSS style in shadow root")
	}
}

func TestRenderFragment_NoCSSProducesNoStyleTag(t *testing.T) {
	engine := newMockEngine()
	engine.addComponent("no-css", "<span>hi</span>", "")

	registry := jsengine.NewRegistry()
	registry.Register("no-css", "fake-bundle")

	output, err := RenderFragmentWithEngine(`<no-css></no-css>`, engine, registry, nil)
	if err != nil {
		t.Fatal(err)
	}

	if strings.Contains(output, "<style>") {
		t.Error("should not have a <style> tag when CSS is empty")
	}
	if !strings.Contains(output, "<span>hi</span>") {
		t.Error("missing rendered HTML")
	}
}

func TestRenderFragment_IgnoredElement(t *testing.T) {
	engine := newMockEngine()
	engine.addComponent("my-el", "<p>rendered</p>", "")

	registry := jsengine.NewRegistry()
	registry.Register("my-el", "fake-bundle")

	ignored := map[string]bool{"my-el": true}
	output, err := RenderFragmentWithEngine(`<my-el></my-el>`, engine, registry, ignored)
	if err != nil {
		t.Fatal(err)
	}

	if strings.Contains(output, "shadowrootmode") {
		t.Error("ignored element should not be expanded")
	}
}

func TestRenderFragment_ExistingDSDSkipped(t *testing.T) {
	engine := newMockEngine()
	engine.addComponent("my-el", "<p>new</p>", "")

	registry := jsengine.NewRegistry()
	registry.Register("my-el", "fake-bundle")

	input := `<my-el><template shadowrootmode="open"><p>existing</p></template></my-el>`
	output, err := RenderFragmentWithEngine(input, engine, registry, nil)
	if err != nil {
		t.Fatal(err)
	}

	if strings.Contains(output, "new") {
		t.Error("element with existing DSD should not be re-rendered")
	}
	if !strings.Contains(output, "existing") {
		t.Error("existing DSD content should be preserved")
	}
}

func TestRenderFragment_RenderError(t *testing.T) {
	engine := newMockEngine()
	engine.addComponentError("bad-el", "render exploded")

	registry := jsengine.NewRegistry()
	registry.Register("bad-el", "fake-bundle")

	ctx := &transformContext{engine: engine, registry: registry}
	wrapper := parseFragment(t, `<bad-el></bad-el>`)

	err := renderHTMLBatched(wrapper, ctx, 1)
	if err != nil {
		t.Fatal(err)
	}

	if len(ctx.renderErrors) != 1 {
		t.Fatalf("expected 1 render error, got %d", len(ctx.renderErrors))
	}
	if ctx.renderErrors[0].TagName != "bad-el" {
		t.Errorf("error tag = %q, want %q", ctx.renderErrors[0].TagName, "bad-el")
	}
	if !strings.Contains(ctx.renderErrors[0].Err.Error(), "render exploded") {
		t.Errorf("error message = %q, want to contain 'render exploded'", ctx.renderErrors[0].Err.Error())
	}
}

func TestRenderFragment_UnregisteredElement(t *testing.T) {
	engine := newMockEngine()

	registry := jsengine.NewRegistry()

	_, err := RenderFragmentWithEngine(`<unknown-el></unknown-el>`, engine, registry, nil)
	if err != nil {
		t.Fatal(err)
	}

	unregistered := registry.Unregistered()
	found := false
	for _, tag := range unregistered {
		if tag == "unknown-el" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected 'unknown-el' in unregistered list, got %v", unregistered)
	}
}

func TestRenderFragment_NestedElements_DeferHydration(t *testing.T) {
	engine := newMockEngine()
	engine.addComponent("outer-el", `<inner-el></inner-el>`, "")
	engine.addComponent("inner-el", "<p>nested</p>", "")

	registry := jsengine.NewRegistry()
	registry.Register("outer-el", "fake-bundle")
	registry.Register("inner-el", "fake-bundle")

	output, err := RenderFragmentWithEngine(
		`<outer-el></outer-el>`, engine, registry, nil)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(output, "defer-hydration") {
		t.Error("nested element should have defer-hydration attribute")
	}
	if !strings.Contains(output, "<p>nested</p>") {
		t.Error("nested element should be rendered")
	}
}

func TestRenderFragment_PreservesSlotContent(t *testing.T) {
	engine := newMockEngine()
	engine.addComponent("my-el", "<slot></slot>", "")

	registry := jsengine.NewRegistry()
	registry.Register("my-el", "fake-bundle")

	output, err := RenderFragmentWithEngine(
		`<my-el>light content</my-el>`, engine, registry, nil)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(output, "light content") {
		t.Error("light DOM content should be preserved")
	}
	if !strings.Contains(output, "shadowrootmode") {
		t.Error("shadow root should be inserted")
	}
}

func TestRenderHTML_FullDocument(t *testing.T) {
	engine := newMockEngine()
	engine.addComponent("my-el", "<p>doc</p>", "")

	registry := jsengine.NewRegistry()
	registry.Register("my-el", "fake-bundle")

	input := `<!DOCTYPE html><html><head></head><body><my-el></my-el></body></html>`
	output, err := RenderHTMLWithEngine(input, engine, registry, nil)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.HasPrefix(output, "<html>") && !strings.HasPrefix(output, "<!DOCTYPE") {
		t.Error("full document output should start with <html> or <!DOCTYPE")
	}
	if !strings.Contains(output, "shadowrootmode") {
		t.Error("element should be expanded in full document mode")
	}
}

func TestRenderHTML_FragmentMode(t *testing.T) {
	engine := newMockEngine()
	engine.addComponent("my-el", "<p>frag</p>", "")

	registry := jsengine.NewRegistry()
	registry.Register("my-el", "fake-bundle")

	input := `<div><my-el></my-el></div>`
	output, err := RenderHTMLWithEngine(input, engine, registry, nil)
	if err != nil {
		t.Fatal(err)
	}

	if strings.Contains(output, "<html>") {
		t.Error("fragment input should not produce full document wrapper")
	}
	if !strings.Contains(output, "shadowrootmode") {
		t.Error("element should be expanded")
	}
}

func TestIsFullDocument(t *testing.T) {
	cases := []struct {
		input string
		want  bool
	}{
		{`<!DOCTYPE html><html>`, true},
		{`<!doctype HTML>`, true},
		{`<html><head>`, true},
		{`  <html>`, true},
		{`<div>hello</div>`, false},
		{`<my-element></my-element>`, false},
		{``, false},
	}
	for _, tc := range cases {
		if got := isFullDocument(tc.input); got != tc.want {
			t.Errorf("isFullDocument(%q) = %v, want %v", tc.input, got, tc.want)
		}
	}
}

func TestExtractBodyContent(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{`<html><body><p>hi</p></body></html>`, `<p>hi</p>`},
		{`<html><body></body></html>`, ``},
		{`no body tags here`, `no body tags here`},
		{`<html><BODY><p>upper</p></BODY></html>`, `<p>upper</p>`},
		{`<html><Body><p>mixed</p></Body></html>`, `<p>mixed</p>`},
	}
	for _, tc := range cases {
		if got := extractBodyContent(tc.input); got != tc.want {
			t.Errorf("extractBodyContent(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestHasDeclarativeShadowRoot(t *testing.T) {
	cases := []struct {
		name  string
		html  string
		want  bool
	}{
		{
			name: "has shadowrootmode",
			html: `<my-el><template shadowrootmode="open"></template></my-el>`,
			want: true,
		},
		{
			name: "has shadowroot (legacy)",
			html: `<my-el><template shadowroot="open"></template></my-el>`,
			want: true,
		},
		{
			name: "no template",
			html: `<my-el><p>content</p></my-el>`,
			want: false,
		},
		{
			name: "template without shadowrootmode",
			html: `<my-el><template id="foo"></template></my-el>`,
			want: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			node := parseFirstElement(t, tc.html)
			if got := hasDeclarativeShadowRoot(node); got != tc.want {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestRenderFragment_MultipleElements(t *testing.T) {
	engine := newMockEngine()
	engine.addComponent("el-a", "<p>A</p>", "")
	engine.addComponent("el-b", "<p>B</p>", "")

	registry := jsengine.NewRegistry()
	registry.Register("el-a", "fake")
	registry.Register("el-b", "fake")

	output, err := RenderFragmentWithEngine(
		`<el-a></el-a><el-b></el-b>`, engine, registry, nil)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(output, "<p>A</p>") {
		t.Error("missing rendered content for el-a")
	}
	if !strings.Contains(output, "<p>B</p>") {
		t.Error("missing rendered content for el-b")
	}
	if strings.Count(output, `shadowrootmode="open"`) != 2 {
		t.Errorf("expected 2 shadow roots, got %d", strings.Count(output, `shadowrootmode="open"`))
	}
}

func TestRenderFragment_AttributesPassedThrough(t *testing.T) {
	engine := &attrCapturingEngine{}
	registry := jsengine.NewRegistry()
	registry.Register("my-el", "fake")

	_, _ = RenderFragmentWithEngine(
		`<my-el name="World" count="5"></my-el>`, engine, registry, nil)

	if len(engine.captured) != 1 {
		t.Fatalf("expected 1 batch request, got %d", len(engine.captured))
	}
	attrs := engine.captured[0].Attrs
	if attrs["name"] != "World" {
		t.Errorf("name attr = %q, want %q", attrs["name"], "World")
	}
	if attrs["count"] != "5" {
		t.Errorf("count attr = %q, want %q", attrs["count"], "5")
	}
}

// attrCapturingEngine records the batch requests it receives.
type attrCapturingEngine struct {
	captured []jsengine.BatchRequest
}

func (e *attrCapturingEngine) LoadBundleForTag(_ string, _ *jsengine.Registry) (bool, error) {
	return true, nil
}

func (e *attrCapturingEngine) RenderBatch(requests []jsengine.BatchRequest) ([]jsengine.BatchResult, error) {
	e.captured = append(e.captured, requests...)
	var results []jsengine.BatchResult
	for _, r := range requests {
		results = append(results, jsengine.BatchResult{
			ID: r.ID, TagName: r.TagName, HTML: "<p>ok</p>",
		})
	}
	return results, nil
}

func (e *attrCapturingEngine) IsRegistered(_ string) bool { return false }

// --- test helpers ---

func parseFragment(t *testing.T, input string) *html.Node {
	t.Helper()
	nodes, err := html.ParseFragment(strings.NewReader(input), &html.Node{
		Type: html.ElementNode, Data: "body", DataAtom: atom.Body,
	})
	if err != nil {
		t.Fatal(err)
	}
	wrapper := &html.Node{Type: html.ElementNode, Data: "body", DataAtom: atom.Body}
	for _, n := range nodes {
		wrapper.AppendChild(n)
	}
	return wrapper
}

func parseFirstElement(t *testing.T, input string) *html.Node {
	t.Helper()
	nodes, err := html.ParseFragment(strings.NewReader(input), &html.Node{
		Type: html.ElementNode, Data: "body", DataAtom: atom.Body,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(nodes) == 0 {
		t.Fatal("no nodes parsed")
	}
	return nodes[0]
}
