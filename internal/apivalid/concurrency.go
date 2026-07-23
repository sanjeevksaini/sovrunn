package apivalid

import "github.com/sanjeevksaini/sovrunn/internal/apiproblem"

// CheckIfMatch compares an If-Match value with the current opaque
// resourceVersion (D-10, F12-UPDATE-002).
//
// Behavior:
//   - absent If-Match (empty string): nil — unprotected update allowed
//   - exact match with currentResourceVersion: nil
//   - mismatch: 412 STALE_RESOURCE_VERSION Problem so the write does not
//     overwrite current state
//
// Comparison is opaque string equality. Callers that accept quoted HTTP
// ETags MUST normalize the header value to the opaque resourceVersion
// before invoking this helper. RequestID/Instance are left empty so
// adopters may attach correlation identifiers without changing the stable
// problem shape. Secrets, credentials, tokens, and raw storage details
// MUST NOT be added to the returned Problem.
func CheckIfMatch(ifMatch, currentResourceVersion string) *apiproblem.Problem {
	if ifMatch == "" {
		return nil
	}
	if ifMatch == currentResourceVersion {
		return nil
	}
	return apiproblem.New(apiproblem.CodeStaleResourceVersion)
}
