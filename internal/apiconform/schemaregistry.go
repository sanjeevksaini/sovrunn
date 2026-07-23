package apiconform

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// SchemaRegistry loads immutable schema documents by stable schema ID
// (D-01a, D-02; F12-VALIDATION-001(4), F12-IMPL-001).
//
// Implementations MUST:
//   - be immutable after construction (no process-global mutable registry);
//   - never perform network access;
//   - reject any schema ID that would trigger a network fetch;
//   - return a defensive copy of schema bytes from Load so callers cannot
//     mutate the registry's stored documents.
type SchemaRegistry interface {
	Load(schemaID string) (schema []byte, err error)
}

// CanonicalSchemasDir is the repository-relative directory that holds
// FEATURE-0012 canonical contract schemas and the _common/ shared
// sub-schemas. Baseline snapshots live under CanonicalSchemasDir/baseline
// and are intentionally excluded from structural-validation loading.
const CanonicalSchemasDir = "api/schemas"

// ErrSchemaNotFound is returned when Load cannot find a schema for the
// requested ID in the registry snapshot.
var ErrSchemaNotFound = errors.New("apiconform: schema not found")

// ErrSchemaIDRejected is returned when a schema ID is empty, unsafe, or
// would require a network fetch.
var ErrSchemaIDRejected = errors.New("apiconform: schema ID rejected")

// MemorySchemaRegistry is an immutable in-memory SchemaRegistry for tests
// and local composition. It accepts pre-loaded schemas by stable schema ID
// at construction and never mutates afterward.
type MemorySchemaRegistry struct {
	schemas map[string][]byte
}

// NewMemorySchemaRegistry builds an immutable in-memory registry from the
// provided schema ID → document map. The input map and its byte slices are
// defensively copied; later mutation of the caller's map or slices does not
// affect the registry. Schema IDs that would trigger a network fetch are
// rejected at construction. A nil or empty map yields an empty registry.
func NewMemorySchemaRegistry(schemas map[string][]byte) (*MemorySchemaRegistry, error) {
	out := make(map[string][]byte, len(schemas))
	for id, body := range schemas {
		if err := rejectUnsafeSchemaID(id); err != nil {
			return nil, err
		}
		out[id] = copyBytes(body)
	}
	return &MemorySchemaRegistry{schemas: out}, nil
}

// Load returns a defensive copy of the schema document for schemaID.
func (r *MemorySchemaRegistry) Load(schemaID string) ([]byte, error) {
	if r == nil {
		return nil, fmt.Errorf("%w: nil memory registry", ErrSchemaNotFound)
	}
	if err := rejectUnsafeSchemaID(schemaID); err != nil {
		return nil, err
	}
	body, ok := r.schemas[schemaID]
	if !ok {
		return nil, fmt.Errorf("%w: %q", ErrSchemaNotFound, schemaID)
	}
	return copyBytes(body), nil
}

// RepositorySchemaRegistry is the repository canonical SchemaRegistry.
// At construction it loads *.json documents from schemasDir (the on-disk
// api/schemas layout) and schemasDir/_common into an immutable snapshot.
// Baseline files under schemasDir/baseline are not loaded.
//
// Stable schema IDs use the repository-relative form:
//
//	api/schemas/<name>.json
//	api/schemas/_common/<name>.json
//
// schemasDir may be a temporary test directory with the same layout; IDs
// still use the CanonicalSchemasDir prefix so callers bind to stable IDs.
type RepositorySchemaRegistry struct {
	schemas map[string][]byte
}

// NewRepositorySchemaRegistry loads canonical schema files from schemasDir
// (expected layout of api/schemas) into an immutable registry snapshot.
// It performs no network access. Missing schemasDir or a non-directory path
// returns an error. An empty but valid directory yields an empty registry.
func NewRepositorySchemaRegistry(schemasDir string) (*RepositorySchemaRegistry, error) {
	if strings.TrimSpace(schemasDir) == "" {
		return nil, fmt.Errorf("%w: empty schemas directory", ErrSchemaIDRejected)
	}
	info, err := os.Stat(schemasDir)
	if err != nil {
		return nil, fmt.Errorf("apiconform: schemas directory: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("apiconform: schemas path is not a directory: %s", schemasDir)
	}

	out := make(map[string][]byte)
	if err := loadJSONFiles(schemasDir, CanonicalSchemasDir, out); err != nil {
		return nil, err
	}
	commonDir := filepath.Join(schemasDir, "_common")
	if st, err := os.Stat(commonDir); err == nil && st.IsDir() {
		if err := loadJSONFiles(commonDir, CanonicalSchemasDir+"/_common", out); err != nil {
			return nil, err
		}
	} else if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("apiconform: _common schemas directory: %w", err)
	}

	return &RepositorySchemaRegistry{schemas: out}, nil
}

// Load returns a defensive copy of the schema document for schemaID.
func (r *RepositorySchemaRegistry) Load(schemaID string) ([]byte, error) {
	if r == nil {
		return nil, fmt.Errorf("%w: nil repository registry", ErrSchemaNotFound)
	}
	if err := rejectUnsafeSchemaID(schemaID); err != nil {
		return nil, err
	}
	body, ok := r.schemas[schemaID]
	if !ok {
		return nil, fmt.Errorf("%w: %q", ErrSchemaNotFound, schemaID)
	}
	return copyBytes(body), nil
}

func loadJSONFiles(dir, idPrefix string, out map[string][]byte) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("apiconform: read schemas directory %s: %w", dir, err)
	}
	for _, ent := range entries {
		if ent.IsDir() {
			continue
		}
		name := ent.Name()
		if !strings.HasSuffix(name, ".json") {
			continue
		}
		id := idPrefix + "/" + name
		if err := rejectUnsafeSchemaID(id); err != nil {
			return err
		}
		path := filepath.Join(dir, name)
		// Resolve symlinks and ensure the file remains under dir (no escape).
		resolved, err := filepath.EvalSymlinks(path)
		if err != nil {
			return fmt.Errorf("apiconform: resolve schema %s: %w", path, err)
		}
		absDir, err := filepath.Abs(dir)
		if err != nil {
			return fmt.Errorf("apiconform: abs schemas directory: %w", err)
		}
		absResolved, err := filepath.Abs(resolved)
		if err != nil {
			return fmt.Errorf("apiconform: abs schema path: %w", err)
		}
		rel, err := filepath.Rel(absDir, absResolved)
		if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
			return fmt.Errorf("%w: schema path escapes directory: %s", ErrSchemaIDRejected, path)
		}
		body, err := os.ReadFile(resolved)
		if err != nil {
			return fmt.Errorf("apiconform: read schema %s: %w", path, err)
		}
		if _, exists := out[id]; exists {
			return fmt.Errorf("apiconform: duplicate schema ID %q", id)
		}
		out[id] = copyBytes(body)
	}
	return nil
}

// rejectUnsafeSchemaID rejects empty IDs and any ID that indicates a
// network or remote fetch (URI schemes such as http/https/ftp/file).
// Absolute filesystem paths are also rejected so Load never treats a
// schema ID as a remote or absolute fetch target.
func rejectUnsafeSchemaID(schemaID string) error {
	if schemaID == "" {
		return fmt.Errorf("%w: empty schema ID", ErrSchemaIDRejected)
	}
	if strings.Contains(schemaID, "://") {
		return fmt.Errorf("%w: network schema ID %q", ErrSchemaIDRejected, schemaID)
	}
	// Catch scheme-like prefixes without authority (e.g. "http:project").
	if i := strings.IndexByte(schemaID, ':'); i > 0 {
		scheme := strings.ToLower(schemaID[:i])
		if isNetworkScheme(scheme) {
			return fmt.Errorf("%w: network schema ID %q", ErrSchemaIDRejected, schemaID)
		}
	}
	if parsed, err := url.Parse(schemaID); err == nil && parsed.Scheme != "" {
		if isNetworkScheme(strings.ToLower(parsed.Scheme)) {
			return fmt.Errorf("%w: network schema ID %q", ErrSchemaIDRejected, schemaID)
		}
	}
	if filepath.IsAbs(schemaID) || strings.HasPrefix(schemaID, "/") {
		return fmt.Errorf("%w: absolute schema ID %q", ErrSchemaIDRejected, schemaID)
	}
	return nil
}

func isNetworkScheme(scheme string) bool {
	switch scheme {
	case "http", "https", "ftp", "ftps", "file", "ws", "wss", "data", "blob":
		return true
	default:
		return false
	}
}

func copyBytes(in []byte) []byte {
	if in == nil {
		return nil
	}
	out := make([]byte, len(in))
	copy(out, in)
	return out
}

var (
	_ SchemaRegistry = (*MemorySchemaRegistry)(nil)
	_ SchemaRegistry = (*RepositorySchemaRegistry)(nil)
)
