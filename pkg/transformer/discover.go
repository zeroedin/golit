package transformer

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/net/html"

	"github.com/zeroedin/golit/pkg/jsengine"
)

// importRe matches ES module import statements to extract bare-module specifiers.
// Handles: import 'x', import "x", import {...} from 'x', import x from 'x'
var importRe = regexp.MustCompile(`import\s+(?:[^'"]*\s+from\s+)?['"]([^'"]+)['"]`)

// collectSourcePaths extracts import maps and module import specifiers from
// HTML content, resolves them to local file paths, and returns the paths
// that are not yet known to the registry. Does NOT build or register anything.
// siteRoot is the top-level directory passed to TransformDir (e.g. "public/"),
// used to resolve absolute paths like "/node_modules/..." in import maps.
func collectSourcePaths(htmlContent string, htmlDir string, siteRoot string, registry *jsengine.Registry, cliImportMap *jsengine.ImportMap, verbose bool) []string {
	if !strings.Contains(htmlContent, `type="importmap"`) &&
		!strings.Contains(htmlContent, `type='importmap'`) &&
		!strings.Contains(htmlContent, `type=importmap`) &&
		!strings.Contains(htmlContent, `type="module"`) &&
		!strings.Contains(htmlContent, `type='module'`) &&
		!strings.Contains(htmlContent, `type=module`) {
		return nil
	}

	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil
	}

	var htmlImportMap *jsengine.ImportMap
	var moduleSpecifiers []string

	var walkNode func(*html.Node)
	walkNode = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "script" {
			scriptType := getAttr(n, "type")

			if scriptType == "importmap" {
				content := getTextContent(n)
				if content != "" {
					absSiteRoot, _ := filepath.Abs(siteRoot)
					im, err := jsengine.ParseImportMap(content, absSiteRoot)
					if err == nil {
						htmlImportMap = im
					}
				}
			} else if scriptType == "module" {
				content := getTextContent(n)
				if content != "" {
					matches := importRe.FindAllStringSubmatch(content, -1)
					for _, match := range matches {
						if len(match) >= 2 {
							moduleSpecifiers = append(moduleSpecifiers, match[1])
						}
					}
				}
			}
		}

		for child := n.FirstChild; child != nil; child = child.NextSibling {
			walkNode(child)
		}
	}
	walkNode(doc)

	activeMap := cliImportMap
	if activeMap == nil {
		activeMap = htmlImportMap
	}

	if activeMap == nil || len(moduleSpecifiers) == 0 {
		if verbose {
			fmt.Fprintf(os.Stderr, "  golit: discovery: activeMap=%v specifiers=%v\n", activeMap != nil, moduleSpecifiers)
		}
		return nil
	}

	resolvedPaths := activeMap.ResolveAll(moduleSpecifiers)
	if verbose {
		fmt.Fprintf(os.Stderr, "  golit: discovery: %d specifiers -> %d resolved paths\n", len(moduleSpecifiers), len(resolvedPaths))
		for _, p := range resolvedPaths {
			fmt.Fprintf(os.Stderr, "    %s\n", p)
		}
	}

	var localPaths []string
	for _, path := range resolvedPaths {
		if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
			fmt.Fprintf(os.Stderr, "  golit: skipping %s (CDN URL)\n", path)
			fmt.Fprintf(os.Stderr, "         Use --importmap with local paths or --sources for SSR\n")
			continue
		}
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}
		if !registry.HasPath(path) {
			localPaths = append(localPaths, path)
		}
	}

	return localPaths
}

// buildDiscoveredModules takes the union of all source paths collected across
// HTML files, builds a shared runtime + thin modules, and registers them.
func buildDiscoveredModules(paths []string, registry *jsengine.Registry, verbose bool) {
	if len(paths) == 0 {
		return
	}

	// Deduplicate paths.
	seen := make(map[string]bool, len(paths))
	var unique []string
	for _, p := range paths {
		if !seen[p] {
			seen[p] = true
			unique = append(unique, p)
		}
	}
	paths = unique

	nodeModulesDir := jsengine.FindNodeModules(paths[0])

	externals, err := jsengine.DiscoverExternalPackages(paths, nodeModulesDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  golit: warning: external discovery failed: %v\n", err)
		return
	}

	modules, err := jsengine.BundleComponentModules(paths, jsengine.BundleOptions{
		ExternalPackages: externals,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "  golit: warning: batch module build failed: %v\n", err)
		return
	}

	if registry.SharedRuntime() == "" && nodeModulesDir != "" {
		rt, rtErr := jsengine.BundleSharedRuntime(nodeModulesDir, modules)
		if rtErr != nil {
			fmt.Fprintf(os.Stderr, "  golit: warning: shared runtime build failed: %v\n", rtErr)
		} else {
			registry.SetSharedRuntime(rt)
		}
	}

	modules = jsengine.RewriteModuleImports(modules, externals)

	for path, mod := range modules {
		tagName, err := jsengine.DiscoverTagName(mod)
		if err != nil {
			if verbose {
				fmt.Fprintf(os.Stderr, "  golit: warning: could not discover tag in %s: %v\n", path, err)
			}
			registry.MarkPath(path)
			continue
		}

		if !registry.Has(tagName) {
			registry.RegisterModule(tagName, mod)
			fmt.Fprintf(os.Stderr, "  golit: auto-discovered <%s> from %s\n", tagName, path)
		}
		registry.MarkPath(path)
	}
}

func getAttr(n *html.Node, key string) string {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func getTextContent(n *html.Node) string {
	var buf strings.Builder
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.TextNode {
			buf.WriteString(child.Data)
		}
	}
	return buf.String()
}
