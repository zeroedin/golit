package jsengine

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ImportMap represents a parsed import map per the WICG specification.
// https://github.com/WICG/import-maps
type ImportMap struct {
	// Imports maps bare-module specifiers to resolved URLs/paths.
	Imports map[string]string `json:"imports"`

	// Scopes provides context-dependent import resolution overrides.
	Scopes map[string]map[string]string `json:"scopes,omitempty"`

	// BaseDir is the directory relative to which map values are resolved.
	BaseDir string `json:"-"`
}

// ParseImportMap parses an import map from a JSON string.
// baseDir is the directory relative to which paths in the map are resolved.
func ParseImportMap(jsonStr string, baseDir string) (*ImportMap, error) {
	var im ImportMap
	if err := json.Unmarshal([]byte(jsonStr), &im); err != nil {
		return nil, fmt.Errorf("parsing import map JSON: %w", err)
	}
	if im.Imports == nil {
		im.Imports = make(map[string]string)
	}
	im.BaseDir = baseDir
	return &im, nil
}

// LoadImportMapFile reads and parses an import map from a JSON file.
func LoadImportMapFile(path string) (*ImportMap, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading import map file %s: %w", path, err)
	}
	absDir, err := filepath.Abs(filepath.Dir(path))
	if err != nil {
		return nil, err
	}
	return ParseImportMap(string(data), absDir)
}

// Resolve resolves a bare-module specifier using the import map.
// Returns the resolved file path, or "" if the specifier doesn't match.
//
// Resolution follows the import map spec:
//   - Exact match: "lit" -> "./node_modules/lit/index.js"
//   - Prefix match: "@rhds/elements/" -> "./node_modules/@rhds/elements/elements/"
//     so "@rhds/elements/rh-badge/rh-badge.js" resolves to
//     "./node_modules/@rhds/elements/elements/rh-badge/rh-badge.js"
func (im *ImportMap) Resolve(specifier string) string {
	if im == nil || im.Imports == nil {
		return ""
	}

	// 1. Try exact match first
	if target, ok := im.Imports[specifier]; ok {
		return im.resolveToAbsPath(target)
	}

	// 2. Try prefix match (keys ending with "/")
	// Sort by length descending so longer/more-specific prefixes match first
	prefixes := make([]string, 0)
	for key := range im.Imports {
		if strings.HasSuffix(key, "/") {
			prefixes = append(prefixes, key)
		}
	}
	sort.Slice(prefixes, func(i, j int) bool {
		return len(prefixes[i]) > len(prefixes[j])
	})

	for _, prefix := range prefixes {
		if strings.HasPrefix(specifier, prefix) {
			suffix := specifier[len(prefix):]
			target := im.Imports[prefix]
			resolved := target + suffix
			return im.resolveToAbsPath(resolved)
		}
	}

	return ""
}

// resolveToAbsPath resolves a potentially relative path against BaseDir.
// URLs (http://, https://) are returned as-is.
// Paths starting with "/" are treated as site-root-relative (joined with BaseDir),
// not as filesystem-absolute paths.
func (im *ImportMap) resolveToAbsPath(target string) string {
	if target == "" {
		return ""
	}

	// URLs pass through unchanged -- golit requires local files for bundling
	// but we preserve the URL so the caller can detect and warn about it.
	if strings.HasPrefix(target, "http://") || strings.HasPrefix(target, "https://") {
		return target
	}

	// All non-URL paths are resolved relative to BaseDir.
	// In HTML import maps, "/" means "site root", not filesystem root.
	// BaseDir should be set to the site root directory so "/foo" resolves
	// to "<site-root>/foo".
	if im.BaseDir != "" {
		// Strip leading "/" so filepath.Join treats it as relative to BaseDir
		cleaned := target
		if strings.HasPrefix(cleaned, "/") {
			cleaned = cleaned[1:]
		}
		return filepath.Join(im.BaseDir, cleaned)
	}

	return target
}

// ResolveAll resolves a list of specifiers and returns the file paths
// that could be resolved. Unresolvable specifiers are silently skipped.
func (im *ImportMap) ResolveAll(specifiers []string) []string {
	var paths []string
	seen := make(map[string]bool)

	for _, spec := range specifiers {
		resolved := im.Resolve(spec)
		if resolved != "" && !seen[resolved] {
			seen[resolved] = true
			paths = append(paths, resolved)
		}
	}

	return paths
}
