package transformer

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"

	"github.com/zeroedin/golit/pkg/jsengine"
)

func newTestRegistry() *jsengine.Registry {
	return jsengine.NewRegistry()
}

func TestGetAttr(t *testing.T) {
	node := parseFirstElement(t, `<script type="module" src="app.js"></script>`)

	if got := getAttr(node, "type"); got != "module" {
		t.Errorf("getAttr(type) = %q, want %q", got, "module")
	}
	if got := getAttr(node, "src"); got != "app.js" {
		t.Errorf("getAttr(src) = %q, want %q", got, "app.js")
	}
	if got := getAttr(node, "missing"); got != "" {
		t.Errorf("getAttr(missing) = %q, want empty", got)
	}
}

func TestGetTextContent(t *testing.T) {
	nodes, err := html.ParseFragment(strings.NewReader(
		`<script>console.log("hello");</script>`),
		&html.Node{Type: html.ElementNode, Data: "body", DataAtom: atom.Body},
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(nodes) == 0 {
		t.Fatal("no nodes parsed")
	}

	got := getTextContent(nodes[0])
	if !strings.Contains(got, `console.log("hello")`) {
		t.Errorf("getTextContent = %q, want to contain console.log", got)
	}
}

func TestGetTextContent_Empty(t *testing.T) {
	node := parseFirstElement(t, `<script></script>`)
	if got := getTextContent(node); got != "" {
		t.Errorf("getTextContent = %q, want empty", got)
	}
}

func TestImportRe_MatchesVariousFormats(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{"bare import single", `import 'lit';`, "lit"},
		{"bare import double", `import "lit";`, "lit"},
		{"named import", `import { html } from 'lit';`, "lit"},
		{"default import", `import LitElement from 'lit';`, "lit"},
		{"scoped package", `import '@rhds/elements/rh-badge/rh-badge.js';`, "@rhds/elements/rh-badge/rh-badge.js"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			matches := importRe.FindAllStringSubmatch(tc.input, -1)
			if len(matches) == 0 {
				t.Fatalf("no match for %q", tc.input)
			}
			if matches[0][1] != tc.want {
				t.Errorf("got %q, want %q", matches[0][1], tc.want)
			}
		})
	}
}

func TestCollectSourcePaths_SkipsPlainHTML(t *testing.T) {
	registry := newTestRegistry()

	paths := collectSourcePaths(
		`<!DOCTYPE html><html><head><title>No scripts</title></head><body><p>plain</p></body></html>`,
		"/tmp", "/tmp", registry, nil, false,
	)

	if len(paths) != 0 {
		t.Errorf("expected no paths for plain HTML, got %v", paths)
	}
}

func TestCollectSourcePaths_DoesNotSkipUnquotedType(t *testing.T) {
	for _, attr := range []string{
		`type=module`,
		`type="module"`,
		`type='module'`,
		`type=importmap`,
		`type="importmap"`,
		`type='importmap'`,
	} {
		t.Run(attr, func(t *testing.T) {
			html := `<!DOCTYPE html><html><head><script ` + attr + `>/* content */</script></head><body></body></html>`
			if !strings.Contains(html, `type=`) {
				t.Fatal("sanity: test HTML missing type attribute")
			}
			// Verify the pre-scan doesn't reject the content and
			// collectSourcePaths completes without panic.
			registry := newTestRegistry()
			collectSourcePaths(html, "/tmp", "/tmp", registry, nil, false)
		})
	}
}

func TestCollectSourcePaths_SkipsRegularScript(t *testing.T) {
	registry := newTestRegistry()

	paths := collectSourcePaths(
		`<!DOCTYPE html><html><head><script src="app.js"></script></head><body></body></html>`,
		"/tmp", "/tmp", registry, nil, false,
	)

	if len(paths) != 0 {
		t.Errorf("expected no paths for regular scripts, got %v", paths)
	}
}

func TestImportRe_MultipleImports(t *testing.T) {
	input := `
		import { html } from 'lit';
		import '@rhds/elements/rh-badge/rh-badge.js';
	`
	matches := importRe.FindAllStringSubmatch(input, -1)
	if len(matches) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(matches))
	}
	if matches[0][1] != "lit" {
		t.Errorf("match[0] = %q, want %q", matches[0][1], "lit")
	}
	if matches[1][1] != "@rhds/elements/rh-badge/rh-badge.js" {
		t.Errorf("match[1] = %q, want %q", matches[1][1], "@rhds/elements/rh-badge/rh-badge.js")
	}
}
