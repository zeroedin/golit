package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/zeroedin/golit"
)

// runServe starts an HTTP server with a long-lived Renderer for POST /render.
// Reduces per-request cold start vs shelling out to `golit transform` each time.
func runServe(args []string) error {
	var (
		defsDir    string
		sourcesDir string
		listen     = "127.0.0.1:9777"
		listenFlag bool
		ignored    []string
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
		default:
			if strings.HasPrefix(args[i], "-") {
				return fmt.Errorf("unknown flag: %s", args[i])
			}
			return fmt.Errorf("unexpected argument: %s", args[i])
		}
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

	renderer, err := golit.NewRenderer(golit.RendererOptions{
		DefsDir:    defsDir,
		SourcesDir: sourcesDir,
		Ignored:    ignored,
	})
	if err != nil {
		return fmt.Errorf("golit serve: init renderer: %w", err)
	}
	defer renderer.Close()

	var mu sync.Mutex
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

		mu.Lock()
		out, err := renderer.RenderHTML(string(body))
		mu.Unlock()
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
