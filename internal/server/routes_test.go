package server

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel"
)

func TestCloudModelRoutePatterns_ExactArithmetic(t *testing.T) {
	patterns := CloudModelRoutePatterns()
	if len(patterns) != CloudModelRouteCount {
		t.Fatalf("pattern count = %d, want %d", len(patterns), CloudModelRouteCount)
	}

	var collection, patchableItem, participationItem, actions int
	seen := make(map[string]bool, len(patterns))
	for _, p := range patterns {
		if seen[p] {
			t.Fatalf("duplicate pattern: %q", p)
		}
		seen[p] = true
		if !hasExplicitMethodPrefix(p) {
			t.Fatalf("path-only or malformed pattern: %q", p)
		}
		if strings.Contains(p, "...") || strings.Contains(p, "{path") {
			t.Fatalf("catch-all pattern forbidden: %q", p)
		}
		switch {
		case strings.Contains(p, "/actions/"):
			if !strings.HasPrefix(p, "POST ") {
				t.Fatalf("action pattern must be POST: %q", p)
			}
			if strings.Contains(p, ":") && strings.Contains(p, "/{uid}:") {
				t.Fatalf("retired /{uid}:<action> form forbidden: %q", p)
			}
			actions++
		case strings.HasPrefix(p, "GET ") && strings.HasSuffix(p, "/cloud-provider-participations/{uid}"):
			participationItem++
		case strings.Contains(p, "/{uid}") && (strings.HasPrefix(p, "GET ") || strings.HasPrefix(p, "PATCH ")):
			patchableItem++
		case strings.HasPrefix(p, "GET ") || strings.HasPrefix(p, "POST "):
			collection++
		default:
			t.Fatalf("unclassified pattern: %q", p)
		}
	}

	if collection != 14 {
		t.Fatalf("collection registrations = %d, want 14", collection)
	}
	if patchableItem != 12 {
		t.Fatalf("PATCHable-item registrations = %d, want 12", patchableItem)
	}
	if participationItem != 1 {
		t.Fatalf("participation-item registrations = %d, want 1", participationItem)
	}
	if actions != 8 {
		t.Fatalf("participation action registrations = %d, want 8", actions)
	}
	if collection+patchableItem+participationItem+actions != CloudModelRouteCount {
		t.Fatalf("arithmetic %d+%d+%d+%d != %d", collection, patchableItem, participationItem, actions, CloudModelRouteCount)
	}

	logicalPaths := map[string]struct{}{}
	for _, p := range patterns {
		logicalPaths[logicalPathOf(p)] = struct{}{}
	}
	if len(logicalPaths) != CloudModelLogicalPathCount {
		t.Fatalf("logical paths = %d, want %d", len(logicalPaths), CloudModelLogicalPathCount)
	}

	// No PATCH for participation item.
	for _, p := range patterns {
		if strings.HasPrefix(p, "PATCH ") && strings.Contains(p, "cloud-provider-participations") {
			t.Fatalf("participation PATCH must not be registered: %q", p)
		}
	}

	// No ExecutionTarget / ServiceRegion / migration routes.
	forbidden := []string{"execution-target", "service-region", "migration"}
	for _, p := range patterns {
		low := strings.ToLower(p)
		for _, f := range forbidden {
			if strings.Contains(low, f) {
				t.Fatalf("excluded route present: %q", p)
			}
		}
	}
}

func TestRegisterCloudModelRoutes_AllPatternsReachable(t *testing.T) {
	mux := http.NewServeMux()
	hits := make(map[string]int)
	h := stubCloudModelHandlers(func(label string) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			hits[label]++
			if r.Method == http.MethodHead {
				w.Header().Set("Allow", "GET")
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			w.WriteHeader(http.StatusNoContent)
		})
	})
	var logBuf bytes.Buffer
	logger := log.New(&logBuf, "", 0)
	RegisterCloudModelRoutes(mux, h, logger)

	if !strings.Contains(logBuf.String(), "count=35") || !strings.Contains(logBuf.String(), "logical_paths=22") {
		t.Fatalf("startup registration log missing: %q", logBuf.String())
	}

	for _, pattern := range CloudModelRoutePatterns() {
		method, path := splitPattern(pattern)
		path = strings.ReplaceAll(path, "{uid}", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
		req := httptest.NewRequest(method, path, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code == http.StatusNotFound {
			t.Fatalf("pattern %q not registered (got 404)", pattern)
		}
		if rec.Code != http.StatusNoContent && rec.Code != http.StatusUnauthorized && rec.Code != http.StatusMethodNotAllowed {
			// Stub returns 204; auth middleware is not applied. Allow 204 only.
			if rec.Code != http.StatusNoContent {
				t.Fatalf("pattern %q status = %d, want 204", pattern, rec.Code)
			}
		}
	}
	if len(hits) == 0 {
		t.Fatal("expected stub handlers to be invoked")
	}
}

func TestRegisterCloudModelRoutes_HEADMatchedByGETReturns405(t *testing.T) {
	mux := http.NewServeMux()
	rt := NewCloudModelRuntime(&cloudmodel.MemoryAuditAppender{})
	grants := cloudmodel.NewBootstrapGrantResolver(cloudmodel.BootstrapGrantConfig{
		PrincipalID: "bootstrap-principal",
	})
	RegisterCloudModelRoutes(mux, NewCloudModelHandlers(rt, grants), log.New(io.Discard, "", 0))

	getPaths := []string{
		"/apis/core.sovrunn.io/v1alpha1/cloud-platforms",
		"/apis/core.sovrunn.io/v1alpha1/cloud-platforms/aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		"/apis/core.sovrunn.io/v1alpha1/cloud-providers",
		"/apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations",
		"/apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations/aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		"/apis/infrastructure.sovrunn.io/v1alpha1/hosting-locations",
		"/apis/infrastructure.sovrunn.io/v1alpha1/datacenters",
		"/apis/infrastructure.sovrunn.io/v1alpha1/fault-domains",
		"/apis/infrastructure.sovrunn.io/v1alpha1/infrastructure-stacks",
	}
	for _, path := range getPaths {
		req := httptest.NewRequest(http.MethodHead, path, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusMethodNotAllowed {
			t.Fatalf("HEAD %s status = %d, want 405 (no HEAD registration)", path, rec.Code)
		}
		if rec.Body.Len() != 0 {
			t.Fatalf("HEAD %s must not write a body, got %q", path, rec.Body.String())
		}
	}

	// HEAD must not appear among the 35 registered patterns.
	for _, p := range CloudModelRoutePatterns() {
		if strings.HasPrefix(p, "HEAD ") {
			t.Fatalf("HEAD must not be registered: %q", p)
		}
	}
}

func TestRegisterCloudModelRoutes_UnsupportedMethods405(t *testing.T) {
	mux := http.NewServeMux()
	rt := NewCloudModelRuntime(&cloudmodel.MemoryAuditAppender{})
	grants := cloudmodel.NewBootstrapGrantResolver(cloudmodel.BootstrapGrantConfig{
		PrincipalID: "bootstrap-principal",
	})
	RegisterCloudModelRoutes(mux, NewCloudModelHandlers(rt, grants), nil)

	cases := []struct {
		method string
		path   string
	}{
		{http.MethodPut, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms/aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
		{http.MethodDelete, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms/aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
		{http.MethodPut, "/apis/core.sovrunn.io/v1alpha1/cloud-providers/aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
		{http.MethodDelete, "/apis/core.sovrunn.io/v1alpha1/cloud-providers/aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
		{http.MethodPut, "/apis/infrastructure.sovrunn.io/v1alpha1/hosting-locations/aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
		{http.MethodDelete, "/apis/infrastructure.sovrunn.io/v1alpha1/hosting-locations/aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
		{http.MethodPut, "/apis/infrastructure.sovrunn.io/v1alpha1/datacenters/aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
		{http.MethodDelete, "/apis/infrastructure.sovrunn.io/v1alpha1/datacenters/aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
		{http.MethodPut, "/apis/infrastructure.sovrunn.io/v1alpha1/fault-domains/aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
		{http.MethodDelete, "/apis/infrastructure.sovrunn.io/v1alpha1/fault-domains/aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
		{http.MethodPut, "/apis/infrastructure.sovrunn.io/v1alpha1/infrastructure-stacks/aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
		{http.MethodDelete, "/apis/infrastructure.sovrunn.io/v1alpha1/infrastructure-stacks/aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
		{http.MethodPatch, "/apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations/aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
		{http.MethodPut, "/apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations/aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
		{http.MethodDelete, "/apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations/aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
	}
	for _, tc := range cases {
		req := httptest.NewRequest(tc.method, tc.path, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusMethodNotAllowed {
			t.Fatalf("%s %s status = %d, want 405", tc.method, tc.path, rec.Code)
		}
	}
}

func TestRegisterCloudModelRoutes_PathValueUIDWildcard(t *testing.T) {
	mux := http.NewServeMux()
	const wantUID = "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	var gotUID string
	h := &CloudModelHandlers{
		CloudPlatformItem: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotUID = r.PathValue("uid")
			w.WriteHeader(http.StatusNoContent)
		}),
	}
	// Minimal registration probe: only the item GET pattern needed for PathValue.
	mux.Handle(routeCloudPlatformGet, h.CloudPlatformItem)

	req := httptest.NewRequest(http.MethodGet, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms/"+wantUID, nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204", rec.Code)
	}
	if gotUID != wantUID {
		t.Fatalf("PathValue(uid) = %q, want %q", gotUID, wantUID)
	}

	// Catch-all form must not be used in FEATURE-0015 patterns.
	for _, p := range CloudModelRoutePatterns() {
		if strings.Contains(p, "{uid...}") || strings.Contains(p, "{$") {
			t.Fatalf("catch-all wildcard forbidden: %q", p)
		}
		if strings.Contains(p, "{uid}") {
			// Complete single-segment wildcard only.
			if !strings.Contains(p, "/{uid}") {
				t.Fatalf("{uid} must be a complete path segment: %q", p)
			}
		}
	}
}

func TestRegisterCloudModelRoutes_ActionsFormNotColon(t *testing.T) {
	mux := http.NewServeMux()
	actionHits := map[string]int{}
	h := stubCloudModelHandlers(func(label string) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actionHits[label]++
			w.WriteHeader(http.StatusNoContent)
		})
	})
	RegisterCloudModelRoutes(mux, h, nil)

	uid := "cccccccccccccccccccccccccccccccc"
	for _, action := range []string{
		"accept", "reject", "withdraw", "suspend", "resume",
		"request-release", "accept-release", "decline-release",
	} {
		path := "/apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations/" + uid + "/actions/" + action
		req := httptest.NewRequest(http.MethodPost, path, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusNoContent {
			t.Fatalf("POST %s status = %d, want 204", path, rec.Code)
		}

		// Retired colon form must not reach an action handler.
		colon := "/apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations/" + uid + ":" + action
		before := actionHits["ParticipationAccept"] + actionHits["ParticipationReject"] +
			actionHits["ParticipationWithdraw"] + actionHits["ParticipationSuspend"] +
			actionHits["ParticipationResume"] + actionHits["ParticipationRequestRelease"] +
			actionHits["ParticipationAcceptRelease"] + actionHits["ParticipationDeclineRelease"]
		req = httptest.NewRequest(http.MethodPost, colon, nil)
		rec = httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		after := actionHits["ParticipationAccept"] + actionHits["ParticipationReject"] +
			actionHits["ParticipationWithdraw"] + actionHits["ParticipationSuspend"] +
			actionHits["ParticipationResume"] + actionHits["ParticipationRequestRelease"] +
			actionHits["ParticipationAcceptRelease"] + actionHits["ParticipationDeclineRelease"]
		if after != before {
			t.Fatalf("colon form %s must not invoke participation action handlers", colon)
		}
		if rec.Code == http.StatusNoContent {
			t.Fatalf("colon form %s must not succeed as an action route", colon)
		}
	}
}

func TestRegisterCloudModelRoutes_NoExecutionTargetRoute(t *testing.T) {
	mux := http.NewServeMux()
	RegisterCloudModelRoutes(mux, stubCloudModelHandlers(func(string) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		})
	}), nil)

	paths := []string{
		"/apis/core.sovrunn.io/v1alpha1/execution-targets",
		"/apis/infrastructure.sovrunn.io/v1alpha1/execution-targets",
		"/apis/core.sovrunn.io/v1alpha1/service-regions",
	}
	for _, path := range paths {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("%s status = %d, want 404 (excluded feature)", path, rec.Code)
		}
	}
}

func TestRoutesGo_NoPathOnlyCatchAllReflectionOrDispatch(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	srcPath := filepath.Join(filepath.Dir(thisFile), "routes.go")
	raw, err := os.ReadFile(srcPath)
	if err != nil {
		t.Fatalf("read routes.go: %v", err)
	}
	src := string(raw)

	forbidden := []string{
		"reflect.",
		"{uid...}",
		"switch r.Method",
	}
	for _, f := range forbidden {
		if strings.Contains(src, f) {
			t.Fatalf("routes.go must not contain %q", f)
		}
	}

	// Every registration must be an explicit method/path Handle call.
	handleCount := strings.Count(src, "mux.Handle(")
	if handleCount != CloudModelRouteCount {
		t.Fatalf("mux.Handle call count = %d, want %d", handleCount, CloudModelRouteCount)
	}

	// Path-only Handle registrations are forbidden (must include METHOD prefix in const).
	for _, line := range strings.Split(src, "\n") {
		trim := strings.TrimSpace(line)
		if !strings.HasPrefix(trim, "mux.Handle(") {
			continue
		}
		if strings.Contains(trim, "Handle(\"/") || strings.Contains(trim, "Handle(`/") {
			t.Fatalf("path-only registration forbidden: %s", trim)
		}
	}

	for _, p := range CloudModelRoutePatterns() {
		method, path := splitPattern(p)
		fragment := method + " " + path
		if !strings.Contains(src, fragment) {
			t.Fatalf("routes.go missing pattern literal %q", fragment)
		}
		if strings.Contains(path, "/actions/") {
			continue
		}
		// Retired colon-action form must not appear as a registration literal.
		if strings.Contains(fragment, "/{uid}:") {
			t.Fatalf("retired colon-action registration forbidden: %q", fragment)
		}
	}
}

func stubCloudModelHandlers(makeHandler func(label string) http.Handler) *CloudModelHandlers {
	return &CloudModelHandlers{
		CloudPlatformCollection:       makeHandler("CloudPlatformCollection"),
		CloudPlatformItem:             makeHandler("CloudPlatformItem"),
		CloudProviderCollection:       makeHandler("CloudProviderCollection"),
		CloudProviderItem:             makeHandler("CloudProviderItem"),
		ParticipationCollection:       makeHandler("ParticipationCollection"),
		ParticipationItem:             makeHandler("ParticipationItem"),
		HostingLocationCollection:     makeHandler("HostingLocationCollection"),
		HostingLocationItem:           makeHandler("HostingLocationItem"),
		DatacenterCollection:          makeHandler("DatacenterCollection"),
		DatacenterItem:                makeHandler("DatacenterItem"),
		FaultDomainCollection:         makeHandler("FaultDomainCollection"),
		FaultDomainItem:               makeHandler("FaultDomainItem"),
		InfrastructureStackCollection: makeHandler("InfrastructureStackCollection"),
		InfrastructureStackItem:       makeHandler("InfrastructureStackItem"),
		ParticipationAccept:           makeHandler("ParticipationAccept"),
		ParticipationReject:           makeHandler("ParticipationReject"),
		ParticipationWithdraw:         makeHandler("ParticipationWithdraw"),
		ParticipationSuspend:          makeHandler("ParticipationSuspend"),
		ParticipationResume:           makeHandler("ParticipationResume"),
		ParticipationRequestRelease:   makeHandler("ParticipationRequestRelease"),
		ParticipationAcceptRelease:    makeHandler("ParticipationAcceptRelease"),
		ParticipationDeclineRelease:   makeHandler("ParticipationDeclineRelease"),
	}
}

func hasExplicitMethodPrefix(pattern string) bool {
	for _, m := range []string{"GET ", "POST ", "PATCH ", "PUT ", "DELETE ", "HEAD "} {
		if strings.HasPrefix(pattern, m) {
			return true
		}
	}
	return false
}

func splitPattern(pattern string) (method, path string) {
	parts := strings.SplitN(pattern, " ", 2)
	if len(parts) != 2 {
		return "", pattern
	}
	return parts[0], parts[1]
}

func logicalPathOf(pattern string) string {
	_, path := splitPattern(pattern)
	return path
}
