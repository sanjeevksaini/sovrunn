package apiconform

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"path"
	"strings"
)

// CommonSchemasPrefix is the stable registry ID prefix for shared sub-schemas
// under api/schemas/_common. Approved $ref targets MUST resolve under this
// prefix (D-01a; F12-VALIDATION-001(4)).
const CommonSchemasPrefix = CanonicalSchemasDir + "/_common/"

// DefaultMaxRefDepth is the default finite $ref follow depth enforced by
// LocalRefResolver when maxDepth is unset or non-positive at construction.
const DefaultMaxRefDepth = 10

// ErrRefRejected is returned when a $ref is empty, remote, absolute,
// fragment-bearing, or resolves outside api/schemas/_common.
var ErrRefRejected = errors.New("apiconform: $ref rejected")

// ErrRefCycle is returned when following $ref edges encounters a repeated
// schema ID in the active resolution stack (A → B → A).
var ErrRefCycle = errors.New("apiconform: $ref cycle")

// ErrRefDepthExceeded is returned when following $ref edges would exceed the
// configured finite depth limit.
var ErrRefDepthExceeded = errors.New("apiconform: $ref depth exceeded")

// RefResolver safely resolves approved relative $ref values under
// api/schemas/_common against an immutable SchemaRegistry (D-01a).
//
// Any rejection returns an error so structural validation fails closed at
// LayerStructural when the adapter (task 8.2) consumes this contract.
type RefResolver interface {
	// Resolve resolves ref relative to baseSchemaID and returns the
	// canonical registry schema ID plus a defensive copy of the target
	// document. It rejects remote URIs, absolute paths, traversal outside
	// api/schemas/_common, missing targets, reference cycles reachable from
	// the resolved document, and depth overflow.
	Resolve(baseSchemaID, ref string) (schemaID string, schema []byte, err error)
}

// LocalRefResolver is the repository-local RefResolver implementation.
// It performs no network access and never reads the filesystem directly;
// targets are loaded only through the configured SchemaRegistry.
type LocalRefResolver struct {
	registry SchemaRegistry
	maxDepth int
}

// NewLocalRefResolver builds a LocalRefResolver bound to registry.
// A nil registry is rejected. Non-positive maxDepth selects DefaultMaxRefDepth.
func NewLocalRefResolver(registry SchemaRegistry, maxDepth int) (*LocalRefResolver, error) {
	if registry == nil {
		return nil, fmt.Errorf("%w: nil schema registry", ErrRefRejected)
	}
	if maxDepth <= 0 {
		maxDepth = DefaultMaxRefDepth
	}
	return &LocalRefResolver{registry: registry, maxDepth: maxDepth}, nil
}

// MaxDepth returns the configured finite $ref follow depth.
func (r *LocalRefResolver) MaxDepth() int {
	if r == nil {
		return 0
	}
	return r.maxDepth
}

// Resolve implements RefResolver.
func (r *LocalRefResolver) Resolve(baseSchemaID, ref string) (string, []byte, error) {
	if r == nil {
		return "", nil, fmt.Errorf("%w: nil ref resolver", ErrRefRejected)
	}
	id, body, err := r.resolve(baseSchemaID, ref, nil)
	if err != nil {
		return "", nil, err
	}
	// Walk $ref edges reachable from the target to detect cycles and
	// depth overflow before returning success to the caller.
	if err := r.walkRefs(id, body, []string{id}); err != nil {
		return "", nil, err
	}
	return id, body, nil
}

func (r *LocalRefResolver) resolve(baseSchemaID, ref string, stack []string) (string, []byte, error) {
	id, err := canonicalizeLocalRef(baseSchemaID, ref)
	if err != nil {
		return "", nil, err
	}
	if len(stack) >= r.maxDepth {
		return "", nil, fmt.Errorf("%w: maxDepth=%d at %q", ErrRefDepthExceeded, r.maxDepth, id)
	}
	for _, seen := range stack {
		if seen == id {
			return "", nil, fmt.Errorf("%w: %q", ErrRefCycle, id)
		}
	}
	body, err := r.registry.Load(id)
	if err != nil {
		if errors.Is(err, ErrSchemaNotFound) {
			return "", nil, fmt.Errorf("%w: %q", ErrSchemaNotFound, id)
		}
		if errors.Is(err, ErrSchemaIDRejected) {
			return "", nil, fmt.Errorf("%w: %v", ErrRefRejected, err)
		}
		return "", nil, err
	}
	return id, body, nil
}

// walkRefs depth-first follows file $ref values found in schema JSON,
// detecting cycles and enforcing maxDepth. Fragment-only refs are ignored
// here because canonicalizeLocalRef already rejects non-empty fragments on
// the Resolve entry path; nested fragment-only values are skipped as they
// are not file references under _common.
func (r *LocalRefResolver) walkRefs(baseID string, body []byte, stack []string) error {
	refs, err := collectFileRefs(body)
	if err != nil {
		return fmt.Errorf("%w: schema %q: %v", ErrRefRejected, baseID, err)
	}
	for _, ref := range refs {
		if isFragmentOnlyRef(ref) {
			continue
		}
		id, nested, err := r.resolve(baseID, ref, stack)
		if err != nil {
			return err
		}
		nextStack := append(append([]string(nil), stack...), id)
		if err := r.walkRefs(id, nested, nextStack); err != nil {
			return err
		}
	}
	return nil
}

func canonicalizeLocalRef(baseSchemaID, ref string) (string, error) {
	if strings.TrimSpace(baseSchemaID) == "" {
		return "", fmt.Errorf("%w: empty base schema ID", ErrRefRejected)
	}
	if err := rejectUnsafeSchemaID(baseSchemaID); err != nil {
		return "", fmt.Errorf("%w: unsafe base schema ID %q", ErrRefRejected, baseSchemaID)
	}
	if strings.TrimSpace(ref) == "" {
		return "", fmt.Errorf("%w: empty $ref", ErrRefRejected)
	}

	refPath, frag, hasFrag := strings.Cut(ref, "#")
	if hasFrag && frag != "" {
		return "", fmt.Errorf("%w: JSON Pointer fragments are not supported in $ref %q", ErrRefRejected, ref)
	}
	if refPath == "" {
		return "", fmt.Errorf("%w: empty $ref path in %q", ErrRefRejected, ref)
	}

	if err := rejectUnsafeRefPath(refPath); err != nil {
		return "", err
	}

	baseDir := path.Dir(baseSchemaID)
	joined := path.Join(baseDir, refPath)
	cleaned := path.Clean(joined)
	// path.Clean may produce values without a stable slash form; normalize
	// to forward-slash registry IDs (schema IDs never use OS separators).
	cleaned = strings.ReplaceAll(cleaned, "\\", "/")

	if cleaned == ".." || strings.HasPrefix(cleaned, "../") {
		return "", fmt.Errorf("%w: $ref %q escapes via traversal", ErrRefRejected, ref)
	}
	if strings.Contains(cleaned, "/../") || strings.HasSuffix(cleaned, "/..") {
		return "", fmt.Errorf("%w: $ref %q retains traversal", ErrRefRejected, ref)
	}
	if !strings.HasPrefix(cleaned, CommonSchemasPrefix) {
		return "", fmt.Errorf("%w: $ref %q resolves to %q outside %s", ErrRefRejected, ref, cleaned, CommonSchemasPrefix)
	}
	rel := strings.TrimPrefix(cleaned, CommonSchemasPrefix)
	if rel == "" || rel == "." || strings.Contains(rel, "/") {
		// Only flat files under _common are approved (no nested dirs).
		if rel == "" || rel == "." {
			return "", fmt.Errorf("%w: $ref %q does not name a _common schema file", ErrRefRejected, ref)
		}
		// Nested paths under _common are still under the prefix; allow a
		// single path segment only to keep the approved surface tight.
		return "", fmt.Errorf("%w: $ref %q targets nested path %q under _common", ErrRefRejected, ref, cleaned)
	}
	if !strings.HasSuffix(rel, ".json") {
		return "", fmt.Errorf("%w: $ref %q must target a .json schema under _common", ErrRefRejected, ref)
	}
	return cleaned, nil
}

func rejectUnsafeRefPath(refPath string) error {
	if strings.Contains(refPath, "://") {
		return fmt.Errorf("%w: remote $ref %q", ErrRefRejected, refPath)
	}
	if i := strings.IndexByte(refPath, ':'); i > 0 {
		scheme := strings.ToLower(refPath[:i])
		if isNetworkScheme(scheme) {
			return fmt.Errorf("%w: remote $ref %q", ErrRefRejected, refPath)
		}
	}
	if parsed, err := url.Parse(refPath); err == nil && parsed.Scheme != "" {
		if isNetworkScheme(strings.ToLower(parsed.Scheme)) {
			return fmt.Errorf("%w: remote $ref %q", ErrRefRejected, refPath)
		}
	}
	// Absolute filesystem / URI paths (leading slash) are rejected.
	if path.IsAbs(refPath) || strings.HasPrefix(refPath, "/") {
		return fmt.Errorf("%w: absolute $ref %q", ErrRefRejected, refPath)
	}
	// Windows-style absolute paths.
	if len(refPath) >= 2 && refPath[1] == ':' {
		return fmt.Errorf("%w: absolute $ref %q", ErrRefRejected, refPath)
	}
	return nil
}

func collectFileRefs(body []byte) ([]string, error) {
	var root any
	if err := json.Unmarshal(body, &root); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}
	var out []string
	collectRefValues(root, &out)
	return out, nil
}

func collectRefValues(node any, out *[]string) {
	switch v := node.(type) {
	case map[string]any:
		if ref, ok := v["$ref"]; ok {
			if s, ok := ref.(string); ok {
				*out = append(*out, s)
			}
		}
		for _, child := range v {
			collectRefValues(child, out)
		}
	case []any:
		for _, child := range v {
			collectRefValues(child, out)
		}
	}
}

func isFragmentOnlyRef(ref string) bool {
	return strings.HasPrefix(ref, "#")
}

var _ RefResolver = (*LocalRefResolver)(nil)
