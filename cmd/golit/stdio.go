package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/zeroedin/golit/pkg/jsengine"
	"github.com/zeroedin/golit/pkg/transformer"
)

func runStdio(stdin io.Reader, stdout io.Writer, pool *jsengine.EnginePool, registry *jsengine.Registry, ignored map[string]bool) error {
	fmt.Fprintf(os.Stderr, "golit serve: stdio mode, reading NUL-delimited requests from stdin\n")

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	type readResult struct {
		data string
		err  error
	}

	reader := bufio.NewReader(stdin)
	writer := bufio.NewWriter(stdout)

	for {
		ch := make(chan readResult, 1)
		go func() {
			data, err := reader.ReadString('\x00')
			ch <- readResult{data, err}
		}()

		select {
		case <-sigCh:
			fmt.Fprintf(os.Stderr, "golit serve: received signal, shutting down\n")
			return nil
		case result := <-ch:
			if result.err != nil {
				if result.err == io.EOF {
					if result.data != "" {
						input := strings.TrimSpace(result.data)
						if input != "" {
							engine := pool.Get()
							out, _ := transformer.RenderHTMLWithEngine(input, engine, registry, ignored)
							pool.Put(engine)
							if _, err := writer.WriteString(out); err != nil {
								return fmt.Errorf("writing stdout: %w", err)
							}
							if err := writer.WriteByte('\x00'); err != nil {
								return fmt.Errorf("writing stdout: %w", err)
							}
							_ = writer.Flush()
						}
					}
					fmt.Fprintf(os.Stderr, "golit serve: stdin closed, shutting down\n")
					return nil
				}
				return fmt.Errorf("reading stdin: %w", result.err)
			}

			input := strings.TrimSuffix(result.data, "\x00")
			if input == "" {
				if err := writer.WriteByte('\x00'); err != nil {
					return fmt.Errorf("writing stdout: %w", err)
				}
				if err := writer.Flush(); err != nil {
					return fmt.Errorf("flushing stdout: %w", err)
				}
				continue
			}

			start := time.Now()
			engine := pool.Get()
			out, renderErr := transformer.RenderHTMLWithEngine(input, engine, registry, ignored)
			pool.Put(engine)
			dur := time.Since(start)

			fmt.Fprintf(os.Stderr, "golit serve: render %.1fms\n",
				float64(dur.Microseconds())/1000.0)

			if renderErr != nil {
				fmt.Fprintf(os.Stderr, "golit serve: render error: %v\n", renderErr)
				if err := writer.WriteByte('\x00'); err != nil {
					return fmt.Errorf("writing stdout: %w", err)
				}
				if err := writer.Flush(); err != nil {
					return fmt.Errorf("flushing stdout: %w", err)
				}
				continue
			}

			if _, err := writer.WriteString(out); err != nil {
				return fmt.Errorf("writing stdout: %w", err)
			}
			if err := writer.WriteByte('\x00'); err != nil {
				return fmt.Errorf("writing stdout: %w", err)
			}
			if err := writer.Flush(); err != nil {
				return fmt.Errorf("flushing stdout: %w", err)
			}
		}
	}
}
