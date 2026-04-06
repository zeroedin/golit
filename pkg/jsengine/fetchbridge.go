package jsengine

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fastschema/qjs"
)

// SSR fetch defaults (override with GOLIT_FETCH_TIMEOUT_SEC, GOLIT_FETCH_MAX_BODY_BYTES).
const (
	defaultFetchTimeout  = 10 * time.Second
	defaultFetchMaxBody  = 16 << 20 // 16 MiB
	minFetchTimeout      = 1 * time.Second
	maxFetchTimeout      = 5 * time.Minute
	maxFetchMaxBody      = 64 << 20
	envFetchAllowlist    = "GOLIT_FETCH_ALLOWLIST"
	envFetchTimeoutSec   = "GOLIT_FETCH_TIMEOUT_SEC"
	envFetchMaxBodyBytes = "GOLIT_FETCH_MAX_BODY_BYTES"
	envSSRLocation       = "GOLIT_SSR_LOCATION"
)

type fetchInitJSON struct {
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

type fetchResultJSON struct {
	OK         bool   `json:"ok"`
	Status     int    `json:"status"`
	StatusText string `json:"statusText"`
	Body       string `json:"body"`
	Error      string `json:"error,omitempty"`
}

// fetchHostAllowlist returns lowercase hostnames from GOLIT_FETCH_ALLOWLIST (comma-separated).
// nil means no allowlist (only http/https are still required).
func fetchHostAllowlist() map[string]bool {
	raw := strings.TrimSpace(os.Getenv(envFetchAllowlist))
	if raw == "" {
		return nil
	}
	m := make(map[string]bool)
	for _, h := range strings.Split(raw, ",") {
		h = strings.TrimSpace(strings.ToLower(h))
		if h != "" {
			m[h] = true
		}
	}
	return m
}

func fetchTimeout() time.Duration {
	s := strings.TrimSpace(os.Getenv(envFetchTimeoutSec))
	if s == "" {
		return defaultFetchTimeout
	}
	sec, err := strconv.Atoi(s)
	if err != nil || sec <= 0 {
		return defaultFetchTimeout
	}
	d := time.Duration(sec) * time.Second
	if d < minFetchTimeout {
		return minFetchTimeout
	}
	if d > maxFetchTimeout {
		return maxFetchTimeout
	}
	return d
}

func fetchMaxBody() int64 {
	s := strings.TrimSpace(os.Getenv(envFetchMaxBodyBytes))
	if s == "" {
		return defaultFetchMaxBody
	}
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil || n <= 0 {
		return defaultFetchMaxBody
	}
	if n > maxFetchMaxBody {
		return maxFetchMaxBody
	}
	return n
}

func ssrLocationHref() string {
	h := strings.TrimSpace(os.Getenv(envSSRLocation))
	if h == "" {
		return "http://localhost/"
	}
	return h
}

func (e *Engine) injectSSRStringGlobals() error {
	href, err := json.Marshal(ssrLocationHref())
	if err != nil {
		return err
	}
	_, err = e.ctx.Eval("golit-ssr-location.js", qjs.Code("globalThis.__golitLocationHref="+string(href)))
	return err
}

// installFetchBridge registers globalThis.__golitFetch(url, initJson) for domshim fetch().
func (e *Engine) installFetchBridge() {
	allow := fetchHostAllowlist()
	timeout := fetchTimeout()
	maxBody := fetchMaxBody()

	e.ctx.SetFunc("__golitFetch", func(th *qjs.This) (*qjs.Value, error) {
		c := th.Context()
		args := th.Args()
		if len(args) < 1 {
			return fetchErrorString(c, "__golitFetch: missing url")
		}
		if args[0].IsUndefined() || args[0].IsNull() {
			return fetchErrorString(c, "__golitFetch: url is null or undefined")
		}
		urlStr := args[0].String()
		initJSON := "{}"
		if len(args) >= 2 && !args[1].IsUndefined() && !args[1].IsNull() {
			initJSON = args[1].String()
		}

		out, err := runGolitFetch(urlStr, initJSON, allow, timeout, maxBody)
		if err != nil {
			return fetchErrorString(c, err.Error())
		}
		// Return a JSON string so JS can JSON.parse (matches domshim fetch wrapper).
		return c.NewString(string(out)), nil
	})
}

func fetchErrorString(c *qjs.Context, msg string) (*qjs.Value, error) {
	b, err := json.Marshal(fetchResultJSON{Error: msg})
	if err != nil {
		return nil, err
	}
	return c.NewString(string(b)), nil
}

func runGolitFetch(urlStr, initJSON string, allow map[string]bool, timeout time.Duration, maxBody int64) ([]byte, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("invalid url: %w", err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, fmt.Errorf("only http and https URLs are allowed")
	}
	host := strings.ToLower(u.Hostname())
	if allow != nil && !allow[host] {
		return nil, fmt.Errorf("fetch host %q not in %s", host, envFetchAllowlist)
	}

	var init fetchInitJSON
	if strings.TrimSpace(initJSON) != "" && initJSON != "{}" {
		if err := json.Unmarshal([]byte(initJSON), &init); err != nil {
			return nil, fmt.Errorf("invalid fetch init JSON: %w", err)
		}
	}
	method := strings.ToUpper(strings.TrimSpace(init.Method))
	if method == "" {
		method = http.MethodGet
	}
	if method != http.MethodGet && method != http.MethodHead && method != http.MethodPost &&
		method != http.MethodPut && method != http.MethodPatch && method != http.MethodDelete {
		return nil, fmt.Errorf("unsupported fetch method %q", method)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var bodyReader io.Reader
	if init.Body != "" && (method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch) {
		bodyReader = strings.NewReader(init.Body)
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), bodyReader)
	if err != nil {
		return nil, err
	}
	for k, v := range init.Headers {
		k = strings.TrimSpace(k)
		if k == "" {
			continue
		}
		req.Header.Set(k, v)
	}
	if init.Body != "" && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "text/plain;charset=UTF-8")
	}

	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	limited := io.LimitReader(resp.Body, maxBody+1)
	data, err := io.ReadAll(limited)
	if err != nil {
		return nil, err
	}
	if int64(len(data)) > maxBody {
		return nil, fmt.Errorf("response body exceeds max size (%d bytes)", maxBody)
	}

	statusText := http.StatusText(resp.StatusCode)
	if statusText == "" {
		statusText = resp.Status
	}
	res := fetchResultJSON{
		OK:         resp.StatusCode >= 200 && resp.StatusCode < 300,
		Status:     resp.StatusCode,
		StatusText: statusText,
		Body:       string(data),
	}
	return json.Marshal(res)
}
