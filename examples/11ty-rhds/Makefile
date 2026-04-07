# 11ty + golit serve transform pipeline
GOLIT := ../../dist/golit
BUNDLES_DIR := bundles
SOURCES_DIR := node_modules/@rhds/elements/elements

.PHONY: build bundle serve serve-dev clean golit-build

## bundle: Pre-bundle RHDS components for SSR
bundle: golit-build
	@mkdir -p $(BUNDLES_DIR)
	$(GOLIT) bundle $(SOURCES_DIR)/ --out $(BUNDLES_DIR)/

## build: Full pipeline — bundle first, then 11ty with golit transform
build: bundle
	time npx @11ty/eleventy

## serve: Build then serve the SSR'd output on port 8080
serve: build
	@echo "Serving _site/ at http://localhost:8080"
	python3 -m http.server 8080 --directory _site

## serve-dev: 11ty dev server (no SSR, for authoring)
serve-dev:
	npx @11ty/eleventy --serve

## golit-build: Build golit binary from source
golit-build:
	$(MAKE) -C ../.. build

## clean: Remove 11ty and bundle output
clean:
	rm -rf _site/ $(BUNDLES_DIR)/
