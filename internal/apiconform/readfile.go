package apiconform

import "os"

// readRepoFile reads a trusted repository-local path used by fitness and
// conformance helpers. Callers join a resolved module root with allowlisted
// relative segments; paths are not derived from untrusted request input.
func readRepoFile(path string) ([]byte, error) {
	return os.ReadFile(path) // #nosec G304 -- trusted repo-local fitness/conformance path, not request input.
}
