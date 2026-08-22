package server

import (
	"net/http"
	"strings"
)

// FEATURE-0016 exact path patterns (four patterns; five method/path registrations).
// The collection path registers both GET and POST; item is GET-only; actions are POST-only.
const (
	executionTargetCollectionPath = "/apis/execution.sovrunn.io/v1alpha1/execution-targets"
	executionTargetCollectionBase = executionTargetCollectionPath + "/"

	executionTargetAllowCollection = "GET, POST"
	executionTargetAllowItem       = "GET"
	executionTargetAllowAction     = "POST"

	headerAllow               = "Allow"
	headerContentTypeOptions  = "X-Content-Type-Options"
	contentTypeOptionsNosniff = "nosniff"
)

// ExecutionTargetTransportGuard is the FEATURE-0016 pre-ServeMux method/path
// guard (DD-08 / REQ-F16-03). It is a standalone http.Handler decorator that
// enforces the closed four-path surface:
//
//   - trailing-slash variants of the four F0016 path patterns → transport-only 404
//   - exact F0016 path with a method outside that path's allowed-method set →
//     transport-only 405 with the path's exact literal Allow value
//   - exact F0016 path with an allowed method → delegate to next unchanged
//   - every non-F0016 method/path → delegate to next unchanged
//
// Transport outcomes set X-Content-Type-Options: nosniff, write an empty body,
// and perform no authentication, authorization, audit, idempotency, or
// lifecycle effect. This task does not register routes or handlers; later
// tasks wrap the five real ServeMux registrations with this guard.
func ExecutionTargetTransportGuard(next http.Handler) http.Handler {
	if next == nil {
		next = http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allow, matched, trailingSlash := classifyExecutionTargetPath(r.URL.Path)
		if !matched {
			next.ServeHTTP(w, r)
			return
		}
		if trailingSlash {
			writeExecutionTargetTransport(w, http.StatusNotFound, "")
			return
		}
		if !executionTargetMethodAllowed(r.Method, allow) {
			writeExecutionTargetTransport(w, http.StatusMethodNotAllowed, allow)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// classifyExecutionTargetPath maps a request path onto the closed F0016 surface.
// matched is true only for one of the four exact path patterns or their
// trailing-slash variants. allow is the exact literal Allow header value for
// an exact match (empty when trailingSlash is true or matched is false).
func classifyExecutionTargetPath(path string) (allow string, matched bool, trailingSlash bool) {
	if path == executionTargetCollectionPath {
		return executionTargetAllowCollection, true, false
	}
	if path == executionTargetCollectionBase {
		return "", true, true
	}
	if !strings.HasPrefix(path, executionTargetCollectionBase) {
		return "", false, false
	}

	rest := path[len(executionTargetCollectionBase):]
	uid, after, cut := strings.Cut(rest, "/")
	if uid == "" {
		return "", false, false
	}
	if !cut {
		// Exact item path: .../execution-targets/{uid}
		return executionTargetAllowItem, true, false
	}
	if after == "" {
		// Trailing-slash item path: .../execution-targets/{uid}/
		return "", true, true
	}
	switch after {
	case "actions/qualify":
		return executionTargetAllowAction, true, false
	case "actions/qualify/":
		return "", true, true
	case "actions/retire":
		return executionTargetAllowAction, true, false
	case "actions/retire/":
		return "", true, true
	default:
		return "", false, false
	}
}

// executionTargetMethodAllowed reports whether method is a member of the
// bounded allowed-method set encoded by the exact literal Allow value.
func executionTargetMethodAllowed(method, allowLiteral string) bool {
	switch allowLiteral {
	case executionTargetAllowCollection:
		return method == http.MethodGet || method == http.MethodPost
	case executionTargetAllowItem:
		return method == http.MethodGet
	case executionTargetAllowAction:
		return method == http.MethodPost
	default:
		return false
	}
}

func writeExecutionTargetTransport(w http.ResponseWriter, status int, allowLiteral string) {
	if allowLiteral != "" {
		w.Header().Set(headerAllow, allowLiteral)
	}
	w.Header().Set(headerContentTypeOptions, contentTypeOptionsNosniff)
	w.WriteHeader(status)
}
