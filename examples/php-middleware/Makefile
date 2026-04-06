# PHP + golit SSR middleware
GOLIT      := ../../dist/golit
BUNDLES    := bundles
COMPONENTS := components

.PHONY: build bundle serve serve-raw container container-run bench bench-quick bench-full bench-trace clean

## bundle: Bundle Lit components for golit SSR
bundle: golit-build
	@mkdir -p $(BUNDLES)
	$(GOLIT) bundle $(COMPONENTS)/ --out $(BUNDLES)/

## build: Install deps, bundle components, link assets
build: bundle
	@ln -sfn ../../components public/components

## serve: Build then serve with golit SSR (warm `golit serve` + PHP) on port 8080
serve: build
	@echo "Serving at http://localhost:8080 (golit serve + PHP; warm Renderer)"
	GOLIT_BIN=$(realpath $(GOLIT)) GOLIT_DEFS=$(realpath $(BUNDLES)) \
		./dev-serve.sh php -S 0.0.0.0:8080 -t public public/router.php

## serve-raw: Serve without SSR (for comparison)
serve-raw:
	@echo "Serving at http://localhost:8080 (no SSR — client-side only)"
	php -S 0.0.0.0:8080 -t public

## container: Build the container image (works with podman or docker)
container: golit-build
	podman build -f Containerfile -t golit-php ../..

## container-run: Run the container on port 8080
container-run:
	podman run --rm -p 8080:8080 golit-php

## golit-build: Build golit binary from repo source
golit-build:
	$(MAKE) -C ../.. build

## bench: Run performance benchmark (100 requests, server metrics only)
bench:
	./bench.sh -n 100

## bench-quick: Quick sanity check (20 requests)
bench-quick:
	./bench.sh -n 20

## bench-full: Full benchmark with browser metrics and flame graphs
bench-full:
	./bench.sh -n 100 --browser --trace

## bench-trace: Capture Chrome trace files for flame graph comparison
bench-trace:
	./bench.sh -n 1 --trace

## clean: Remove build artifacts
clean:
	rm -rf $(BUNDLES)/ public/components bench-results/
