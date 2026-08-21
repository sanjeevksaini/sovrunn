package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
)

const (
	etUID          = "01HZXEXAMPLEUID0000000001"
	etCollection   = "/apis/execution.sovrunn.io/v1alpha1/execution-targets"
	etItem         = etCollection + "/" + etUID
	etQualify      = etItem + "/actions/qualify"
	etRetire       = etItem + "/actions/retire"
	etMarkerStatus = 218
	etMarkerBody   = "next-handler-body"
	etMarkerHeader = "X-Next-Handler"
)

// representativeNonAllowedMethods prove the general not-in-allowed-set → 405
// invariant. This is not an exhaustive enumeration of every method name.
var representativeNonAllowedMethods = []string{
	http.MethodHead,
	http.MethodPut,
	http.MethodDelete,
	http.MethodPatch,
	http.MethodOptions,
	http.MethodTrace,
	http.MethodConnect,
	"FOO",
	"PROPFIND",
}

type etPathCase struct {
	name         string
	path         string
	allowLiteral string
	allowed      []string
}

func etFourPathPatterns() []etPathCase {
	return []etPathCase{
		{
			name:         "collection",
			path:         etCollection,
			allowLiteral: "GET, POST",
			allowed:      []string{http.MethodGet, http.MethodPost},
		},
		{
			name:         "item",
			path:         etItem,
			allowLiteral: "GET",
			allowed:      []string{http.MethodGet},
		},
		{
			name:         "qualify",
			path:         etQualify,
			allowLiteral: "POST",
			allowed:      []string{http.MethodPost},
		},
		{
			name:         "retire",
			path:         etRetire,
			allowLiteral: "POST",
			allowed:      []string{http.MethodPost},
		},
	}
}

func newCountingNext(t *testing.T) (http.Handler, *atomic.Int32) {
	t.Helper()
	var calls atomic.Int32
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.Header().Set(etMarkerHeader, "1")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(etMarkerStatus)
		_, _ = io.WriteString(w, etMarkerBody)
	})
	return next, &calls
}

func assertTransportOnly(t *testing.T, rec *httptest.ResponseRecorder, wantStatus int, wantAllow string) {
	t.Helper()
	if rec.Code != wantStatus {
		t.Fatalf("status=%d, want %d", rec.Code, wantStatus)
	}
	if got := rec.Header().Get("Allow"); got != wantAllow {
		t.Fatalf("Allow=%q, want exact literal %q", got, wantAllow)
	}
	if got := rec.Header().Get("X-Content-Type-Options"); got != "nosniff" {
		t.Fatalf("X-Content-Type-Options=%q, want nosniff", got)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "" {
		t.Fatalf("Content-Type=%q, want empty (no Problem media)", ct)
	}
	if loc := rec.Header().Get("Location"); loc != "" {
		t.Fatalf("Location=%q, want empty (no redirect)", loc)
	}
	if body := rec.Body.String(); body != "" {
		t.Fatalf("body=%q, want empty transport-only body", body)
	}
	// Side-effect-free transport: no auth challenge, no audit/idempotency markers.
	if got := rec.Header().Get("WWW-Authenticate"); got != "" {
		t.Fatalf("WWW-Authenticate=%q, want empty (no authentication)", got)
	}
	if got := rec.Header().Get("X-Sovrunn-Audit-Event"); got != "" {
		t.Fatalf("unexpected audit marker header %q", got)
	}
	if got := rec.Header().Get("Idempotency-Key"); got != "" {
		t.Fatalf("unexpected Idempotency-Key response header %q", got)
	}
}

func TestExecutionTargetTransportGuard_DisallowedMethod_Transport405(t *testing.T) {
	for _, pc := range etFourPathPatterns() {
		pc := pc
		t.Run(pc.name, func(t *testing.T) {
			allowed := map[string]struct{}{}
			for _, m := range pc.allowed {
				allowed[m] = struct{}{}
			}
			for _, method := range representativeNonAllowedMethods {
				method := method
				if _, ok := allowed[method]; ok {
					continue
				}
				t.Run(method, func(t *testing.T) {
					next, calls := newCountingNext(t)
					guard := ExecutionTargetTransportGuard(next)
					req := httptest.NewRequest(method, pc.path, nil)
					rec := httptest.NewRecorder()
					guard.ServeHTTP(rec, req)

					assertTransportOnly(t, rec, http.StatusMethodNotAllowed, pc.allowLiteral)
					if calls.Load() != 0 {
						t.Fatalf("next called %d times; transport 405 must not delegate", calls.Load())
					}
				})
			}
		})
	}
}

func TestExecutionTargetTransportGuard_TrailingSlash_Transport404(t *testing.T) {
	trailing := []struct {
		name string
		path string
	}{
		{name: "collection", path: etCollection + "/"},
		{name: "item", path: etItem + "/"},
		{name: "qualify", path: etQualify + "/"},
		{name: "retire", path: etRetire + "/"},
	}
	// Any method on a trailing-slash variant is transport-only 404.
	methods := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodHead,
		http.MethodPut,
		http.MethodDelete,
	}
	for _, tc := range trailing {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			for _, method := range methods {
				method := method
				t.Run(method, func(t *testing.T) {
					next, calls := newCountingNext(t)
					guard := ExecutionTargetTransportGuard(next)
					req := httptest.NewRequest(method, tc.path, nil)
					rec := httptest.NewRecorder()
					guard.ServeHTTP(rec, req)

					assertTransportOnly(t, rec, http.StatusNotFound, "")
					if calls.Load() != 0 {
						t.Fatalf("next called %d times; trailing-slash 404 must not delegate", calls.Load())
					}
				})
			}
		})
	}
}

func TestExecutionTargetTransportGuard_ValidPairs_DelegateOnceUnchanged(t *testing.T) {
	valid := []struct {
		method string
		path   string
	}{
		{http.MethodGet, etCollection},
		{http.MethodPost, etCollection},
		{http.MethodGet, etItem},
		{http.MethodPost, etQualify},
		{http.MethodPost, etRetire},
	}
	for _, tc := range valid {
		tc := tc
		t.Run(tc.method+"_"+tc.path, func(t *testing.T) {
			next, calls := newCountingNext(t)
			guard := ExecutionTargetTransportGuard(next)
			req := httptest.NewRequest(tc.method, tc.path, strings.NewReader(`{}`))
			rec := httptest.NewRecorder()
			guard.ServeHTTP(rec, req)

			if calls.Load() != 1 {
				t.Fatalf("next called %d times, want exactly 1", calls.Load())
			}
			if rec.Code != etMarkerStatus {
				t.Fatalf("status=%d, want next-handler status %d (no guard mutation)", rec.Code, etMarkerStatus)
			}
			if got := rec.Body.String(); got != etMarkerBody {
				t.Fatalf("body=%q, want next-handler body (no guard mutation)", got)
			}
			if got := rec.Header().Get(etMarkerHeader); got != "1" {
				t.Fatalf("%s=%q, want next-handler marker", etMarkerHeader, got)
			}
			// Guard must not inject Allow or rewrite nosniff on a delegated response.
			if got := rec.Header().Get("Allow"); got != "" {
				t.Fatalf("Allow=%q on delegated response; guard must not author response", got)
			}
		})
	}
}

func TestExecutionTargetTransportGuard_NonF0016_PassThroughOnce(t *testing.T) {
	paths := []string{
		"/healthz",
		"/readyz",
		"/v1/organizations",
		"/apis/core.sovrunn.io/v1alpha1/cloud-platforms",
		"/apis/execution.sovrunn.io/v1alpha1/other",
		etCollection + "/extra/segment",
		etItem + "/actions",
		etItem + "/actions/",
		etItem + "/actions/other",
		etItem + "/actions/qualify/extra",
		etCollection + "//actions/qualify",
		"/apis/execution.sovrunn.io/v1alpha1/execution-targets-extra",
	}
	methods := []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodHead}
	for _, path := range paths {
		path := path
		for _, method := range methods {
			method := method
			t.Run(method+"_"+path, func(t *testing.T) {
				next, calls := newCountingNext(t)
				guard := ExecutionTargetTransportGuard(next)
				req := httptest.NewRequest(method, path, nil)
				rec := httptest.NewRecorder()
				guard.ServeHTTP(rec, req)

				if calls.Load() != 1 {
					t.Fatalf("next called %d times, want exactly 1 for non-F0016 path", calls.Load())
				}
				if rec.Code != etMarkerStatus {
					t.Fatalf("status=%d, want next-handler status %d", rec.Code, etMarkerStatus)
				}
				if got := rec.Body.String(); got != etMarkerBody {
					t.Fatalf("body=%q, want next-handler body", got)
				}
			})
		}
	}
}

func TestExecutionTargetTransportGuard_AllowLiteralExactSerialization(t *testing.T) {
	// Bounded allowed-method-set invariant: exact literal Allow values, not an
	// unordered set serialization.
	cases := []struct {
		path   string
		method string
		allow  string
	}{
		{etCollection, http.MethodHead, "GET, POST"},
		{etItem, http.MethodPut, "GET"},
		{etQualify, http.MethodGet, "POST"},
		{etRetire, http.MethodDelete, "POST"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.path, func(t *testing.T) {
			next, calls := newCountingNext(t)
			guard := ExecutionTargetTransportGuard(next)
			req := httptest.NewRequest(tc.method, tc.path, nil)
			rec := httptest.NewRecorder()
			guard.ServeHTTP(rec, req)

			assertTransportOnly(t, rec, http.StatusMethodNotAllowed, tc.allow)
			if calls.Load() != 0 {
				t.Fatalf("next called %d times", calls.Load())
			}
			// Reject common alternate serializations of the same method set.
			got := rec.Header().Get("Allow")
			switch tc.allow {
			case "GET, POST":
				if got == "POST, GET" || got == "GET,POST" || got == "GET POST" {
					t.Fatalf("Allow=%q is not the exact literal %q", got, tc.allow)
				}
			}
		})
	}
}

func TestExecutionTargetTransportGuard_NilNext_Safe(t *testing.T) {
	guard := ExecutionTargetTransportGuard(nil)
	req := httptest.NewRequest(http.MethodHead, etCollection, nil)
	rec := httptest.NewRecorder()
	guard.ServeHTTP(rec, req)
	assertTransportOnly(t, rec, http.StatusMethodNotAllowed, "GET, POST")
}

func TestClassifyExecutionTargetPath(t *testing.T) {
	cases := []struct {
		path         string
		wantAllow    string
		wantMatched  bool
		wantTrailing bool
	}{
		{etCollection, "GET, POST", true, false},
		{etCollection + "/", "", true, true},
		{etItem, "GET", true, false},
		{etItem + "/", "", true, true},
		{etQualify, "POST", true, false},
		{etQualify + "/", "", true, true},
		{etRetire, "POST", true, false},
		{etRetire + "/", "", true, true},
		{"/healthz", "", false, false},
		{etItem + "/actions/other", "", false, false},
		{etCollection + "/", "", true, true},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.path, func(t *testing.T) {
			allow, matched, trailing := classifyExecutionTargetPath(tc.path)
			if allow != tc.wantAllow || matched != tc.wantMatched || trailing != tc.wantTrailing {
				t.Fatalf("classify(%q)=(%q,%v,%v), want (%q,%v,%v)",
					tc.path, allow, matched, trailing, tc.wantAllow, tc.wantMatched, tc.wantTrailing)
			}
		})
	}
}
