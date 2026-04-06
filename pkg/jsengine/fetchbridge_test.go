package jsengine

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/fastschema/qjs"
)

func TestRunGolitFetch_GET(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte("hello"))
	}))
	t.Cleanup(srv.Close)

	raw, err := runGolitFetch(srv.URL, "{}", nil, defaultFetchTimeout, defaultFetchMaxBody)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(raw), `"ok":true`) || !strings.Contains(string(raw), "hello") {
		t.Fatalf("unexpected JSON: %s", raw)
	}
}

func TestRunGolitFetch_AllowlistDeny(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("x"))
	}))
	t.Cleanup(srv.Close)

	pu, err := url.Parse(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	host := strings.ToLower(pu.Hostname())

	allow := map[string]bool{"some-other-host.invalid": true}
	_, err = runGolitFetch(srv.URL, "{}", allow, defaultFetchTimeout, defaultFetchMaxBody)
	if err == nil || !strings.Contains(err.Error(), "not in") {
		t.Fatalf("expected allowlist error, got %v", err)
	}

	allow[host] = true
	raw, err := runGolitFetch(srv.URL, "{}", allow, defaultFetchTimeout, defaultFetchMaxBody)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(raw), `"ok":true`) {
		t.Fatalf("expected ok: %s", raw)
	}
}

func TestRunGolitFetch_SchemeReject(t *testing.T) {
	_, err := runGolitFetch("file:///etc/passwd", "{}", nil, defaultFetchTimeout, defaultFetchMaxBody)
	if err == nil || !strings.Contains(err.Error(), "only http") {
		t.Fatalf("expected scheme error, got %v", err)
	}
}

func TestDOMShim_FetchViaBridge(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"a":1}`))
	}))
	t.Cleanup(srv.Close)

	t.Setenv(envFetchAllowlist, "")
	rt, err := qjs.New()
	if err != nil {
		t.Fatal(err)
	}
	defer rt.Close()

	ctx := rt.Context()
	eng := &Engine{runtime: rt, ctx: ctx, loaded: make(map[string]bool)}
	if err := eng.injectSSRStringGlobals(); err != nil {
		t.Fatal(err)
	}
	eng.installFetchBridge()

	_, err = ctx.Eval("domshim.js", qjs.Code(domShimJS))
	if err != nil {
		t.Fatalf("domshim: %v", err)
	}

	js := `JSON.stringify({
		hasFetch: typeof fetch === 'function',
		loc: typeof location !== 'undefined' && location.href.length > 0,
		mm: typeof matchMedia
	});`
	res, err := ctx.Eval("chk.js", qjs.Code(js))
	if err != nil {
		t.Fatal(err)
	}
	s := res.String()
	if !strings.Contains(s, `"hasFetch":true`) {
		t.Fatalf("expected fetch: %s", s)
	}
	if !strings.Contains(s, `"loc":true`) {
		t.Fatalf("expected location: %s", s)
	}
	if !strings.Contains(s, `"mm":"undefined"`) {
		t.Fatalf("matchMedia should be undefined, got %s", s)
	}

	// fetch() returns a Promise; verify __golitFetch JSON and fetch is a thenable factory.
	fetchJS := `(function() {
		const u = ` + strconv.Quote(srv.URL) + `;
		const raw = globalThis.__golitFetch(u, '{}');
		const d = JSON.parse(raw);
		if (!d.ok || d.status !== 200 || d.body !== '{"a":1}') return 'bad-direct';
		if (typeof fetch !== 'function') return 'bad-fetch-fn';
		const p = fetch(u);
		if (!p || typeof p.then !== 'function') return 'bad-promise';
		return 'ok';
	})();`
	res2, err := ctx.Eval("fetch-test.js", qjs.Code(fetchJS))
	if err != nil {
		t.Fatal(err)
	}
	if res2.String() != "ok" {
		t.Fatalf("fetch wrapper: %s", res2.String())
	}
}
