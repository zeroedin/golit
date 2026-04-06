# Ruby/Rack + golit SSR middleware
GOLIT      := ../../dist/golit
BUNDLES    := bundles
COMPONENTS := components
# Build/run images with podman or docker: `make container CONTAINER_RUNTIME=docker`
CONTAINER_RUNTIME ?= podman

.PHONY: build bundle serve container container-run bench bench-quick bench-full bench-trace clean

## bundle: Bundle Lit components for golit SSR
bundle: golit-build
	@mkdir -p $(BUNDLES)
	$(GOLIT) bundle $(COMPONENTS)/ --out $(BUNDLES)/

## build: Install deps and bundle components
build: bundle

## serve: Build then serve with golit SSR (warm `golit serve` + Rack) on port 9292
serve: build
	@echo "Serving at http://localhost:9292 (golit serve + Rack; warm Renderer)"
	GOLIT_BIN=$(realpath $(GOLIT)) GOLIT_DEFS=$(realpath $(BUNDLES)) \
		./dev-serve.sh bundle exec rackup --host 0.0.0.0 --port 9292

## container: Build the container image (override CONTAINER_RUNTIME=docker if needed)
container: golit-build
	$(CONTAINER_RUNTIME) build -f Containerfile -t golit-ruby ../..

## container-run: Run the container on port 9292
container-run:
	$(CONTAINER_RUNTIME) run --rm -p 9292:9292 golit-ruby

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
	rm -rf $(BUNDLES)/ bench-results/
