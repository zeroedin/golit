package transformer

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"

	"github.com/zeroedin/golit/pkg/jsengine"
)

// RenderHTML takes an HTML string, finds custom elements, and returns
// the transformed HTML with Declarative Shadow DOM.
func RenderHTML(input string, registry *jsengine.Registry, ignored ...map[string]bool) (string, error) {
	var ign map[string]bool
	if len(ignored) > 0 {
		ign = ignored[0]
	}
	return renderHTMLWithIgnored(input, registry, ign)
}

func renderHTMLWithIgnored(input string, registry *jsengine.Registry, ignored map[string]bool) (string, error) {
	engine, err := jsengine.NewEngine()
	if err != nil {
		return "", fmt.Errorf("creating JS engine: %w", err)
	}
	defer engine.Close()
	return RenderHTMLWithEngine(input, engine, registry, ignored)
}

// transformContext carries shared state through the recursive transform walk.
type transformContext struct {
	engine       ElementRenderer
	registry     *jsengine.Registry
	ignored      map[string]bool
	file         string // current HTML file path (for error reporting)
	renderErrors []RenderError
}

// RenderHTMLWithEngine transforms HTML using a caller-provided engine,
// avoiding engine creation overhead. Use this when you have a long-lived
// engine (e.g. in a Renderer).
func RenderHTMLWithEngine(input string, engine ElementRenderer, registry *jsengine.Registry, ignored map[string]bool) (string, error) {
	ctx := &transformContext{engine: engine, registry: registry, ignored: ignored}
	output, err := renderHTMLWithContext(input, ctx)
	return output, err
}

func renderHTMLWithContext(input string, ctx *transformContext) (string, error) {
	doc, err := html.Parse(strings.NewReader(input))
	if err != nil {
		return "", fmt.Errorf("parsing HTML: %w", err)
	}

	if err := renderHTMLBatched(doc, ctx, 10); err != nil {
		return "", err
	}

	var buf bytes.Buffer
	buf.Grow(len(input) + len(input)/2)
	if err := html.Render(&buf, doc); err != nil {
		return "", fmt.Errorf("rendering HTML: %w", err)
	}

	result := buf.String()
	if !isFullDocument(input) {
		result = extractBodyContent(result)
	}
	return result, nil
}

// RenderFragment renders an HTML fragment.
func RenderFragment(input string, registry *jsengine.Registry, ignored ...map[string]bool) (string, error) {
	var ign map[string]bool
	if len(ignored) > 0 {
		ign = ignored[0]
	}
	return renderFragmentWithIgnored(input, registry, ign)
}

func renderFragmentWithIgnored(input string, registry *jsengine.Registry, ignored map[string]bool) (string, error) {
	engine, err := jsengine.NewEngine()
	if err != nil {
		return "", fmt.Errorf("creating JS engine: %w", err)
	}
	defer engine.Close()
	return RenderFragmentWithEngine(input, engine, registry, ignored)
}

// RenderFragmentWithEngine renders an HTML fragment using a caller-provided
// engine, avoiding engine creation overhead.
func RenderFragmentWithEngine(input string, engine ElementRenderer, registry *jsengine.Registry, ignored map[string]bool) (string, error) {
	nodes, err := html.ParseFragment(strings.NewReader(input), &html.Node{
		Type: html.ElementNode, Data: "body", DataAtom: atom.Body,
	})
	if err != nil {
		return "", fmt.Errorf("parsing fragment: %w", err)
	}

	ctx := &transformContext{engine: engine, registry: registry, ignored: ignored}

	wrapper := &html.Node{Type: html.ElementNode, Data: "body", DataAtom: atom.Body}
	for _, node := range nodes {
		wrapper.AppendChild(node)
	}

	if err := renderHTMLBatched(wrapper, ctx, 10); err != nil {
		return "", err
	}

	var buf bytes.Buffer
	buf.Grow(len(input) + len(input)/2)
	for child := wrapper.FirstChild; child != nil; child = child.NextSibling {
		if err := html.Render(&buf, child); err != nil {
			return "", fmt.Errorf("rendering: %w", err)
		}
	}
	return buf.String(), nil
}

// pendingElement tracks a custom element waiting to be expanded in batch mode.
type pendingElement struct {
	node  *html.Node
	depth int
}

// collectUnexpanded walks the HTML tree and returns all custom elements
// that haven't been expanded yet (no <template shadowrootmode> child).
func collectUnexpanded(node *html.Node, ctx *transformContext) []pendingElement {
	var pending []pendingElement
	var walk func(*html.Node, int)
	walk = func(n *html.Node, depth int) {
		if depth > 10 {
			return
		}
		if n.Type == html.ElementNode && strings.Contains(n.Data, "-") {
			if !ctx.ignored[n.Data] && !hasDeclarativeShadowRoot(n) {
				loaded, loadErr := ctx.engine.LoadBundleForTag(n.Data, ctx.registry)
				if loadErr != nil {
					fmt.Fprintf(os.Stderr, "golit: warning: %v\n", loadErr)
				}
				if loaded || ctx.engine.IsRegistered(n.Data) {
					pending = append(pending, pendingElement{node: n, depth: depth})
					return
				}
				ctx.registry.MarkUnregistered(n.Data)
			}
		}
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			walk(child, depth+1)
		}
	}
	walk(node, 0)
	return pending
}

// renderHTMLBatched uses BFS-by-depth to render all custom elements,
// batching all elements at each depth level into a single QJS Eval call.
func renderHTMLBatched(doc *html.Node, ctx *transformContext, maxDepth int) error {
	for depth := 0; depth < maxDepth; depth++ {
		pending := collectUnexpanded(doc, ctx)
		if len(pending) == 0 {
			break
		}

		requests := make([]jsengine.BatchRequest, len(pending))
		for i, p := range pending {
			attrs := make(map[string]string)
			for _, attr := range p.node.Attr {
				attrs[attr.Key] = attr.Val
			}
			requests[i] = jsengine.BatchRequest{
				ID:      i,
				TagName: p.node.Data,
				Attrs:   attrs,
			}
		}

		results, err := ctx.engine.RenderBatch(requests)
		if err != nil {
			return fmt.Errorf("batch render at depth %d: %w", depth, err)
		}

		resultSlice := make([]jsengine.BatchResult, len(pending))
		populated := make([]bool, len(pending))
		for _, r := range results {
			if r.ID >= 0 && r.ID < len(resultSlice) {
				resultSlice[r.ID] = r
				populated[r.ID] = true
			}
		}

		for i, p := range pending {
			if !populated[i] {
				continue
			}
			r := resultSlice[i]
			if r.Error != "" {
				ctx.renderErrors = append(ctx.renderErrors, RenderError{
					TagName: p.node.Data,
					File:    ctx.file,
					Err:     errors.New(r.Error),
				})
				continue
			}

			var shadowContent strings.Builder
			if r.CSS != "" {
				shadowContent.WriteString("<style>")
				shadowContent.WriteString(strings.TrimSpace(r.CSS))
				shadowContent.WriteString("</style>")
			}
			shadowContent.WriteString(r.HTML)

			shadowNodes, err := html.ParseFragment(
				strings.NewReader(shadowContent.String()),
				&html.Node{Type: html.ElementNode, Data: "body", DataAtom: atom.Body},
			)
			if err != nil {
				continue
			}

			templateNode := &html.Node{
				Type: html.ElementNode, Data: "template", DataAtom: atom.Template,
				Attr: []html.Attribute{
					{Key: "shadowroot", Val: "open"},
					{Key: "shadowrootmode", Val: "open"},
				},
			}

			for _, sn := range shadowNodes {
				templateNode.AppendChild(sn)
			}

			if p.node.FirstChild != nil {
				p.node.InsertBefore(templateNode, p.node.FirstChild)
			} else {
				p.node.AppendChild(templateNode)
			}
		}
	}
	return nil
}

func hasDeclarativeShadowRoot(node *html.Node) bool {
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode && child.Data == "template" {
			for _, attr := range child.Attr {
				if attr.Key == "shadowrootmode" || attr.Key == "shadowroot" {
					return true
				}
			}
		}
	}
	return false
}

func isFullDocument(input string) bool {
	s := strings.TrimLeft(input, " \t\n\r\f")
	if len(s) < 5 {
		return false
	}
	n := min(15, len(s))
	prefix := strings.ToLower(s[:n])
	return strings.HasPrefix(prefix, "<!doctype") || strings.HasPrefix(prefix, "<html")
}

// indexFold returns the index of the first case-insensitive match of
// substr in s, or -1 if not found.
func indexFold(s, substr string) int {
	n := len(substr)
	if n == 0 {
		return 0
	}
	if n > len(s) {
		return -1
	}
	for i := 0; i <= len(s)-n; i++ {
		if strings.EqualFold(s[i:i+n], substr) {
			return i
		}
	}
	return -1
}

// lastIndexFold returns the index of the last case-insensitive match
// of substr in s, or -1 if not found.
func lastIndexFold(s, substr string) int {
	n := len(substr)
	if n == 0 {
		return len(s)
	}
	if n > len(s) {
		return -1
	}
	for i := len(s) - n; i >= 0; i-- {
		if strings.EqualFold(s[i:i+n], substr) {
			return i
		}
	}
	return -1
}

func extractBodyContent(rendered string) string {
	bodyStart := indexFold(rendered, "<body>")
	if bodyStart == -1 {
		return rendered
	}
	bodyStart += len("<body>")
	bodyEnd := lastIndexFold(rendered, "</body>")
	if bodyEnd == -1 || bodyEnd < bodyStart {
		return rendered
	}
	return rendered[bodyStart:bodyEnd]
}
