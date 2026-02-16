package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/sspriggs/golit/pkg/transformer"
)

// Config represents a golit.yaml configuration file.
type Config struct {
	// Transform options
	Transform TransformConfig `yaml:"transform"`

	// Bundle options
	Bundle BundleConfig `yaml:"bundle"`
}

// TransformConfig holds transform-specific settings.
type TransformConfig struct {
	// Input is the HTML directory to process.
	Input string `yaml:"input"`

	// Out is the output directory (default: in-place).
	Out string `yaml:"out"`

	// Defs is the directory of pre-bundled .golit.bundle.js files.
	Defs string `yaml:"defs"`

	// Sources is a directory of component source files to auto-bundle.
	Sources string `yaml:"sources"`

	// ImportMap is a path to an import map JSON file.
	ImportMap string `yaml:"importmap"`

	// Ignore is a list of custom element tag names to skip.
	Ignore []string `yaml:"ignore"`

	// Preload is a list of extra JS modules to bundle and load into the
	// QJS engine before component rendering. Useful for dynamic imports
	// that components need at render time (e.g. prism-esm for rh-code-block).
	Preload []string `yaml:"preload"`

	// Verbose enables progress output.
	Verbose bool `yaml:"verbose"`

	// DryRun processes without writing.
	DryRun bool `yaml:"dryRun"`
}

// BundleConfig holds bundle-specific settings.
type BundleConfig struct {
	// Input is the source file or directory to bundle.
	Input string `yaml:"input"`

	// Out is the output file or directory.
	Out string `yaml:"out"`

	// Minify enables minification.
	Minify bool `yaml:"minify"`
}

// LoadConfig reads a golit.yaml config file.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config %s: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config %s: %w", path, err)
	}

	return &cfg, nil
}

// FindConfig looks for golit.yaml in the project root or config/ directory.
// Returns the path if found, or "" if not.
func FindConfig() string {
	for _, name := range []string{"golit.yaml", "golit.yml", "config/golit.yaml", "config/golit.yml"} {
		if _, err := os.Stat(name); err == nil {
			return name
		}
	}
	return ""
}

// ToTransformOptions converts the config's transform section to Options,
// with CLI flags overriding config values.
func (c *Config) ToTransformOptions(cliOpts transformer.Options) transformer.Options {
	opts := transformer.Options{}

	// Start with config values
	if c.Transform.Defs != "" {
		opts.DefsDir = c.Transform.Defs
	}
	if c.Transform.Sources != "" {
		opts.SourcesDir = c.Transform.Sources
	}
	if c.Transform.ImportMap != "" {
		opts.ImportMapFile = c.Transform.ImportMap
	}
	if c.Transform.Out != "" {
		opts.OutDir = c.Transform.Out
	}
	if c.Transform.Verbose {
		opts.Verbose = true
	}
	if c.Transform.DryRun {
		opts.DryRun = true
	}
	if len(c.Transform.Ignore) > 0 {
		opts.Ignored = make(map[string]bool)
		for _, tag := range c.Transform.Ignore {
			opts.Ignored[tag] = true
		}
	}
	if len(c.Transform.Preload) > 0 {
		opts.Preload = append(opts.Preload, c.Transform.Preload...)
	}

	// CLI flags override config
	if cliOpts.DefsDir != "" {
		opts.DefsDir = cliOpts.DefsDir
	}
	if cliOpts.SourcesDir != "" {
		opts.SourcesDir = cliOpts.SourcesDir
	}
	if cliOpts.ImportMapFile != "" {
		opts.ImportMapFile = cliOpts.ImportMapFile
	}
	if cliOpts.OutDir != "" {
		opts.OutDir = cliOpts.OutDir
	}
	if cliOpts.Verbose {
		opts.Verbose = true
	}
	if cliOpts.DryRun {
		opts.DryRun = true
	}
	if cliOpts.Ignored != nil {
		if opts.Ignored == nil {
			opts.Ignored = make(map[string]bool)
		}
		for tag := range cliOpts.Ignored {
			opts.Ignored[tag] = true
		}
	}
	if len(cliOpts.Preload) > 0 {
		opts.Preload = append(opts.Preload, cliOpts.Preload...)
	}

	return opts
}
