package apiconform

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestMemorySchemaRegistry_LoadKnownID(t *testing.T) {
	t.Parallel()

	const id = "api/schemas/project.json"
	body := []byte(`{"$schema":"https://json-schema.org/draft/2020-12/schema","type":"object"}`)
	reg, err := NewMemorySchemaRegistry(map[string][]byte{id: body})
	if err != nil {
		t.Fatalf("NewMemorySchemaRegistry: %v", err)
	}

	got, err := reg.Load(id)
	if err != nil {
		t.Fatalf("Load(%q): %v", id, err)
	}
	if string(got) != string(body) {
		t.Fatalf("Load body = %q, want %q", got, body)
	}
}

func TestMemorySchemaRegistry_LoadUnknownID(t *testing.T) {
	t.Parallel()

	reg, err := NewMemorySchemaRegistry(map[string][]byte{
		"api/schemas/project.json": []byte(`{"type":"object"}`),
	})
	if err != nil {
		t.Fatalf("NewMemorySchemaRegistry: %v", err)
	}

	_, err = reg.Load("api/schemas/missing.json")
	if !errors.Is(err, ErrSchemaNotFound) {
		t.Fatalf("Load unknown: err = %v, want ErrSchemaNotFound", err)
	}
}

func TestMemorySchemaRegistry_ImmutableAfterConstruction(t *testing.T) {
	t.Parallel()

	const id = "api/schemas/project.json"
	original := []byte(`{"type":"object"}`)
	input := map[string][]byte{id: original}

	reg, err := NewMemorySchemaRegistry(input)
	if err != nil {
		t.Fatalf("NewMemorySchemaRegistry: %v", err)
	}

	// Mutate caller's map and byte slice after construction.
	input[id][0] = 'X'
	input["api/schemas/injected.json"] = []byte(`{"injected":true}`)
	delete(input, id)

	got, err := reg.Load(id)
	if err != nil {
		t.Fatalf("Load after caller mutation: %v", err)
	}
	if string(got) != `{"type":"object"}` {
		t.Fatalf("registry mutated via caller map; got %q", got)
	}
	if _, err := reg.Load("api/schemas/injected.json"); !errors.Is(err, ErrSchemaNotFound) {
		t.Fatalf("injected key visible after construction: err=%v", err)
	}

	// Mutating the returned slice must not affect subsequent Load.
	got[0] = 'Y'
	again, err := reg.Load(id)
	if err != nil {
		t.Fatalf("second Load: %v", err)
	}
	if again[0] == 'Y' {
		t.Fatal("returned schema bytes alias registry storage")
	}
	if string(again) != `{"type":"object"}` {
		t.Fatalf("second Load = %q, want original document", again)
	}
}

func TestMemorySchemaRegistry_RejectsNetworkSchemaID(t *testing.T) {
	t.Parallel()

	networkIDs := []string{
		"https://example.com/schemas/project.json",
		"http://example.com/schemas/project.json",
		"ftp://example.com/schemas/project.json",
		"file:///tmp/project.json",
		"https:project.json",
	}
	for _, id := range networkIDs {
		_, err := NewMemorySchemaRegistry(map[string][]byte{id: []byte(`{}`)})
		if !errors.Is(err, ErrSchemaIDRejected) {
			t.Fatalf("construct with %q: err=%v, want ErrSchemaIDRejected", id, err)
		}
	}

	reg, err := NewMemorySchemaRegistry(map[string][]byte{
		"api/schemas/project.json": []byte(`{}`),
	})
	if err != nil {
		t.Fatalf("NewMemorySchemaRegistry: %v", err)
	}
	for _, id := range networkIDs {
		_, err := reg.Load(id)
		if !errors.Is(err, ErrSchemaIDRejected) {
			t.Fatalf("Load(%q): err=%v, want ErrSchemaIDRejected", id, err)
		}
	}
}

func TestMemorySchemaRegistry_RejectsAbsoluteSchemaID(t *testing.T) {
	t.Parallel()

	_, err := NewMemorySchemaRegistry(map[string][]byte{
		"/etc/passwd": []byte(`{}`),
	})
	if !errors.Is(err, ErrSchemaIDRejected) {
		t.Fatalf("absolute ID construct: err=%v, want ErrSchemaIDRejected", err)
	}

	reg, err := NewMemorySchemaRegistry(nil)
	if err != nil {
		t.Fatalf("empty registry: %v", err)
	}
	_, err = reg.Load("/abs/path.json")
	if !errors.Is(err, ErrSchemaIDRejected) {
		t.Fatalf("Load absolute: err=%v, want ErrSchemaIDRejected", err)
	}
}

func TestRepositorySchemaRegistry_LoadKnownAndUnknown(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	projectBody := []byte(`{"title":"project"}`)
	commonBody := []byte(`{"title":"typed-ref"}`)
	if err := os.WriteFile(filepath.Join(root, "project.json"), projectBody, 0o600); err != nil {
		t.Fatalf("write project.json: %v", err)
	}
	if err := os.Mkdir(filepath.Join(root, "_common"), 0o700); err != nil {
		t.Fatalf("mkdir _common: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "_common", "typed-ref.json"), commonBody, 0o600); err != nil {
		t.Fatalf("write typed-ref.json: %v", err)
	}
	// baseline/ must not be loaded into the structural registry.
	if err := os.Mkdir(filepath.Join(root, "baseline"), 0o700); err != nil {
		t.Fatalf("mkdir baseline: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "baseline", "project.json"), []byte(`{"baseline":true}`), 0o600); err != nil {
		t.Fatalf("write baseline: %v", err)
	}

	reg, err := NewRepositorySchemaRegistry(root)
	if err != nil {
		t.Fatalf("NewRepositorySchemaRegistry: %v", err)
	}

	got, err := reg.Load("api/schemas/project.json")
	if err != nil {
		t.Fatalf("Load project: %v", err)
	}
	if string(got) != string(projectBody) {
		t.Fatalf("project body = %q, want %q", got, projectBody)
	}

	got, err = reg.Load("api/schemas/_common/typed-ref.json")
	if err != nil {
		t.Fatalf("Load common: %v", err)
	}
	if string(got) != string(commonBody) {
		t.Fatalf("common body = %q, want %q", got, commonBody)
	}

	if _, err := reg.Load("api/schemas/baseline/project.json"); !errors.Is(err, ErrSchemaNotFound) {
		t.Fatalf("baseline should not load: err=%v", err)
	}
	if _, err := reg.Load("api/schemas/missing.json"); !errors.Is(err, ErrSchemaNotFound) {
		t.Fatalf("unknown: err=%v, want ErrSchemaNotFound", err)
	}
}

func TestRepositorySchemaRegistry_ImmutableAfterConstruction(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	path := filepath.Join(root, "project.json")
	if err := os.WriteFile(path, []byte(`{"v":1}`), 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}

	reg, err := NewRepositorySchemaRegistry(root)
	if err != nil {
		t.Fatalf("NewRepositorySchemaRegistry: %v", err)
	}

	// Mutating the on-disk file after construction must not change Load.
	if err := os.WriteFile(path, []byte(`{"v":2}`), 0o600); err != nil {
		t.Fatalf("rewrite: %v", err)
	}
	got, err := reg.Load("api/schemas/project.json")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if string(got) != `{"v":1}` {
		t.Fatalf("registry not immutable snapshot; got %q", got)
	}

	got[2] = '9'
	again, err := reg.Load("api/schemas/project.json")
	if err != nil {
		t.Fatalf("second Load: %v", err)
	}
	if string(again) != `{"v":1}` {
		t.Fatalf("returned bytes aliased storage; got %q", again)
	}
}

func TestRepositorySchemaRegistry_RejectsNetworkLoad(t *testing.T) {
	t.Parallel()

	reg, err := NewRepositorySchemaRegistry(t.TempDir())
	if err != nil {
		t.Fatalf("NewRepositorySchemaRegistry: %v", err)
	}
	_, err = reg.Load("https://example.com/schema.json")
	if !errors.Is(err, ErrSchemaIDRejected) {
		t.Fatalf("network Load: err=%v, want ErrSchemaIDRejected", err)
	}
}

func TestNewRepositorySchemaRegistry_MissingDir(t *testing.T) {
	t.Parallel()

	_, err := NewRepositorySchemaRegistry(filepath.Join(t.TempDir(), "does-not-exist"))
	if err == nil {
		t.Fatal("expected error for missing schemas directory")
	}
}
