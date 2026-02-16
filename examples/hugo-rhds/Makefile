# Hugo + golit build pipeline
MAKEFLAGS += -j2
GOLIT := ../../dist/golit
BUNDLES_DIR := bundles
SOURCES_DIR := node_modules/@rhds/elements/elements

.PHONY: build bundle hugo serve serve-dev clean golit-build

## bundle: Pre-bundle RHDS components for faster transforms
bundle: golit-build
	@mkdir -p $(BUNDLES_DIR)
	$(GOLIT) bundle $(SOURCES_DIR)/ --out $(BUNDLES_DIR)/

## hugo: Build the static site with Hugo
hugo:
	hugo --minify

## build: Full pipeline — Hugo + bundle in parallel, then golit SSR transform
build: bundle hugo
	$(GOLIT) transform public/

## serve: Build then serve the SSR'd output on port 8080
serve: build
	@echo "Serving public/ at http://localhost:8080"
	python3 -m http.server 8080 --directory public

## serve-dev: Hugo dev server (no SSR, for authoring)
serve-dev:
	hugo server --buildDrafts

## golit-build: Build golit binary from source
golit-build:
	$(MAKE) -C ../.. build

## clean: Remove Hugo and bundle output
clean:
	rm -rf public/ $(BUNDLES_DIR)/
