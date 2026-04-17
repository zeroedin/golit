package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/zeroedin/golit/pkg/jsengine"
	"github.com/zeroedin/golit/pkg/transformer"
)

// runServe starts an HTTP server with a pool of QJS engines for POST /render.
// Each concurrent request gets its own engine from the pool, enabling true
// parallel rendering. Defaults to runtime.NumCPU() workers.
func runServe(args []string) error {
	var (
		defsDir     string
		sourcesDir  string
		listen      = "127.0.0.1:9777"
		listenFlag  bool
		stdioMode   bool
		ignored     []string
		preload     []string
		concurrency int
	)

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--defs":
			if i+1 >= len(args) {
				return fmt.Errorf("--defs requires a directory argument")
			}
			defsDir = args[i+1]
			i++
		case "--sources":
			if i+1 >= len(args) {
				return fmt.Errorf("--sources requires a directory argument")
			}
			sourcesDir = args[i+1]
			i++
		case "--listen":
			if i+1 >= len(args) {
				return fmt.Errorf("--listen requires a host:port argument")
			}
			listen = args[i+1]
			listenFlag = true
			i++
		case "--ignore":
			if i+1 >= len(args) {
				return fmt.Errorf("--ignore requires a tag name argument")
			}
			ignored = append(ignored, args[i+1])
			i++
		case "--preload":
			if i+1 >= len(args) {
				return fmt.Errorf("--preload requires a module name argument")
			}
			preload = append(preload, args[i+1])
			i++
		case "--stdio":
			stdioMode = true
		case "--concurrency", "-j":
			if i+1 < len(args) {
				if n, err := strconv.Atoi(args[i+1]); err == nil {
					if n < 1 {
						return fmt.Errorf("--concurrency value must be a positive integer")
					}
					concurrency = n
					i++
				}
			}
			if concurrency == 0 {
				concurrency = runtime.NumCPU()
			}
			i++
		default:
			if strings.HasPrefix(args[i], "-") {
				return fmt.Errorf("unknown flag: %s", args[i])
			}
			return fmt.Errorf("unexpected argument: %s", args[i])
		}
	}

	if stdioMode && listenFlag {
		return fmt.Errorf("golit serve: --stdio and --listen are mutually exclusive")
	}

	if defsDir == "" {
		defsDir = os.Getenv("GOLIT_DEFS")
	}
	if defsDir == "" {
		return fmt.Errorf("golit serve: requires --defs <dir> or GOLIT_DEFS")
	}
	if !listenFlag {
		if v := os.Getenv("GOLIT_SERVE_LISTEN"); v != "" {
			listen = v
		}
	}
	if concurrency == 0 {
		if stdioMode {
			concurrency = 1
		} else {
			concurrency = runtime.NumCPU()
		}
	}

	registry := jsengine.NewRegistry()

	if err := registry.LoadDir(defsDir); err != nil {
		return fmt.Errorf("golit serve: loading bundles: %w", err)
	}

	if sourcesDir != "" {
		if err := registry.LoadSourceDir(sourcesDir); err != nil {
			return fmt.Errorf("golit serve: loading sources: %w", err)
		}
	}

	ignoredMap := make(map[string]bool, len(ignored))
	for _, tag := range ignored {
		ignoredMap[tag] = true
	}

	pool, err := jsengine.NewEnginePool(concurrency)
	if err != nil {
		return fmt.Errorf("golit serve: creating engine pool: %w", err)
	}
	defer pool.Close()

	if err := pool.PreloadAll(registry, preload); err != nil {
		return fmt.Errorf("golit serve: preloading pool: %w", err)
	}

	fmt.Fprintf(os.Stderr, "golit serve: initialized %d engine workers\n", concurrency)

	if stdioMode {
		return runStdio(os.Stdin, os.Stdout, pool, registry, ignoredMap)
	}

	const maxBody = 32 << 20 // 32 MiB

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("/render", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		body, err := io.ReadAll(io.LimitReader(r.Body, maxBody+1))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if len(body) > maxBody {
			http.Error(w, "request body too large", http.StatusRequestEntityTooLarge)
			return
		}

		start := time.Now()
		engine := pool.Get()
		active := pool.Size() - pool.Available()
		out, err := transformer.RenderHTMLWithEngine(string(body), engine, registry, ignoredMap)
		pool.Put(engine)
		dur := time.Since(start)

		w.Header().Set("Server-Timing", fmt.Sprintf(
			`render;dur=%.1f, pool;desc="%d/%d busy"`,
			float64(dur.Microseconds())/1000.0, active, pool.Size()))

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(out))
	})

	srv := &http.Server{
		Addr:              listen,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		fmt.Fprintf(os.Stderr, "golit serve: listening on http://%s (POST /render, GET /health)\n", listen)
		errCh <- srv.ListenAndServe()
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("golit serve: %w", err)
		}
		return nil
	case <-sigCh:
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if shutErr := srv.Shutdown(ctx); shutErr != nil {
			return fmt.Errorf("golit serve: shutdown: %w", shutErr)
		}
		serveErr := <-errCh
		if serveErr != nil && serveErr != http.ErrServerClosed {
			return serveErr
		}
		return nil
	}
}
