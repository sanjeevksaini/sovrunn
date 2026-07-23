package apiconform

import (
	"errors"
	"strings"
	"testing"
)

func testCommonRegistry(t *testing.T, schemas map[string][]byte) SchemaRegistry {
	t.Helper()
	reg, err := NewMemorySchemaRegistry(schemas)
	if err != nil {
		t.Fatalf("NewMemorySchemaRegistry: %v", err)
	}
	return reg
}

func TestLocalRefResolver_ValidLocalRef(t *testing.T) {
	t.Parallel()

	const (
		base   = "api/schemas/project.json"
		ref    = "_common/typed-ref.json"
		wantID = "api/schemas/_common/typed-ref.json"
	)
	body := []byte(`{"type":"object","title":"typed-ref"}`)
	reg := testCommonRegistry(t, map[string][]byte{
		wantID: body,
	})
	r, err := NewLocalRefResolver(reg, DefaultMaxRefDepth)
	if err != nil {
		t.Fatalf("NewLocalRefResolver: %v", err)
	}

	gotID, gotBody, err := r.Resolve(base, ref)
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if gotID != wantID {
		t.Fatalf("schemaID = %q, want %q", gotID, wantID)
	}
	if string(gotBody) != string(body) {
		t.Fatalf("body = %q, want %q", gotBody, body)
	}

	// Relative ref from within _common also resolves.
	gotID, gotBody, err = r.Resolve(wantID, "typed-ref.json")
	if err != nil {
		t.Fatalf("Resolve from _common: %v", err)
	}
	if gotID != wantID {
		t.Fatalf("self-relative schemaID = %q, want %q", gotID, wantID)
	}
	if string(gotBody) != string(body) {
		t.Fatalf("self-relative body = %q, want %q", gotBody, body)
	}
}

func TestLocalRefResolver_MissingRef(t *testing.T) {
	t.Parallel()

	reg := testCommonRegistry(t, map[string][]byte{
		"api/schemas/_common/typed-ref.json": []byte(`{"type":"object"}`),
	})
	r, err := NewLocalRefResolver(reg, DefaultMaxRefDepth)
	if err != nil {
		t.Fatalf("NewLocalRefResolver: %v", err)
	}

	_, _, err = r.Resolve("api/schemas/project.json", "_common/missing.json")
	if !errors.Is(err, ErrSchemaNotFound) {
		t.Fatalf("missing ref: err=%v, want ErrSchemaNotFound", err)
	}
}

func TestLocalRefResolver_RemoteURI(t *testing.T) {
	t.Parallel()

	reg := testCommonRegistry(t, map[string][]byte{
		"api/schemas/_common/typed-ref.json": []byte(`{"type":"object"}`),
	})
	r, err := NewLocalRefResolver(reg, DefaultMaxRefDepth)
	if err != nil {
		t.Fatalf("NewLocalRefResolver: %v", err)
	}

	remote := []string{
		"https://example.com/schemas/typed-ref.json",
		"http://example.com/schemas/typed-ref.json",
		"ftp://example.com/schemas/typed-ref.json",
		"file:///tmp/typed-ref.json",
		"https:typed-ref.json",
	}
	for _, ref := range remote {
		_, _, err := r.Resolve("api/schemas/project.json", ref)
		if !errors.Is(err, ErrRefRejected) {
			t.Fatalf("remote %q: err=%v, want ErrRefRejected", ref, err)
		}
	}
}

func TestLocalRefResolver_AbsolutePath(t *testing.T) {
	t.Parallel()

	reg := testCommonRegistry(t, map[string][]byte{
		"api/schemas/_common/typed-ref.json": []byte(`{"type":"object"}`),
	})
	r, err := NewLocalRefResolver(reg, DefaultMaxRefDepth)
	if err != nil {
		t.Fatalf("NewLocalRefResolver: %v", err)
	}

	absolutes := []string{
		"/etc/passwd",
		"/api/schemas/_common/typed-ref.json",
		"C:/schemas/typed-ref.json",
	}
	for _, ref := range absolutes {
		_, _, err := r.Resolve("api/schemas/project.json", ref)
		if !errors.Is(err, ErrRefRejected) {
			t.Fatalf("absolute %q: err=%v, want ErrRefRejected", ref, err)
		}
	}
}

func TestLocalRefResolver_Traversal(t *testing.T) {
	t.Parallel()

	reg := testCommonRegistry(t, map[string][]byte{
		"api/schemas/_common/typed-ref.json": []byte(`{"type":"object"}`),
		"api/schemas/project.json":           []byte(`{"type":"object"}`),
	})
	r, err := NewLocalRefResolver(reg, DefaultMaxRefDepth)
	if err != nil {
		t.Fatalf("NewLocalRefResolver: %v", err)
	}

	cases := []string{
		"../project.json",
		"../../etc/passwd",
		"../_common/typed-ref.json", // cleans outside api/schemas/_common
		"_common/../project.json",
		"_common/../../etc/passwd",
	}
	for _, ref := range cases {
		_, _, err := r.Resolve("api/schemas/project.json", ref)
		if !errors.Is(err, ErrRefRejected) {
			t.Fatalf("traversal %q: err=%v, want ErrRefRejected", ref, err)
		}
	}

	// Escape from inside _common back to a sibling schema.
	_, _, err = r.Resolve("api/schemas/_common/typed-ref.json", "../project.json")
	if !errors.Is(err, ErrRefRejected) {
		t.Fatalf("escape from _common: err=%v, want ErrRefRejected", err)
	}
}

func TestLocalRefResolver_Cycle(t *testing.T) {
	t.Parallel()

	const (
		aID = "api/schemas/_common/a.json"
		bID = "api/schemas/_common/b.json"
	)
	reg := testCommonRegistry(t, map[string][]byte{
		aID: []byte(`{"$ref":"b.json"}`),
		bID: []byte(`{"$ref":"a.json"}`),
	})
	r, err := NewLocalRefResolver(reg, DefaultMaxRefDepth)
	if err != nil {
		t.Fatalf("NewLocalRefResolver: %v", err)
	}

	_, _, err = r.Resolve("api/schemas/project.json", "_common/a.json")
	if !errors.Is(err, ErrRefCycle) {
		t.Fatalf("cycle: err=%v, want ErrRefCycle", err)
	}
}

func TestLocalRefResolver_DepthOverflow(t *testing.T) {
	t.Parallel()

	// Chain: a → b → c. With maxDepth=2, loading a then following to b is
	// allowed, but following from b to c must fail.
	const (
		aID = "api/schemas/_common/a.json"
		bID = "api/schemas/_common/b.json"
		cID = "api/schemas/_common/c.json"
	)
	reg := testCommonRegistry(t, map[string][]byte{
		aID: []byte(`{"properties":{"x":{"$ref":"b.json"}}}`),
		bID: []byte(`{"properties":{"y":{"$ref":"c.json"}}}`),
		cID: []byte(`{"type":"object"}`),
	})
	r, err := NewLocalRefResolver(reg, 2)
	if err != nil {
		t.Fatalf("NewLocalRefResolver: %v", err)
	}
	if r.MaxDepth() != 2 {
		t.Fatalf("MaxDepth = %d, want 2", r.MaxDepth())
	}

	_, _, err = r.Resolve("api/schemas/project.json", "_common/a.json")
	if !errors.Is(err, ErrRefDepthExceeded) {
		t.Fatalf("depth overflow: err=%v, want ErrRefDepthExceeded", err)
	}

	// Same chain succeeds with a larger depth budget.
	rOK, err := NewLocalRefResolver(reg, 3)
	if err != nil {
		t.Fatalf("NewLocalRefResolver: %v", err)
	}
	gotID, _, err := rOK.Resolve("api/schemas/project.json", "_common/a.json")
	if err != nil {
		t.Fatalf("Resolve with depth=3: %v", err)
	}
	if gotID != aID {
		t.Fatalf("schemaID = %q, want %q", gotID, aID)
	}
}

func TestLocalRefResolver_DefaultMaxDepth(t *testing.T) {
	t.Parallel()

	reg := testCommonRegistry(t, nil)
	r, err := NewLocalRefResolver(reg, 0)
	if err != nil {
		t.Fatalf("NewLocalRefResolver: %v", err)
	}
	if r.MaxDepth() != DefaultMaxRefDepth {
		t.Fatalf("MaxDepth = %d, want default %d", r.MaxDepth(), DefaultMaxRefDepth)
	}
}

func TestLocalRefResolver_NilRegistryRejected(t *testing.T) {
	t.Parallel()

	_, err := NewLocalRefResolver(nil, DefaultMaxRefDepth)
	if !errors.Is(err, ErrRefRejected) {
		t.Fatalf("nil registry: err=%v, want ErrRefRejected", err)
	}
}

func TestLocalRefResolver_FragmentRejected(t *testing.T) {
	t.Parallel()

	reg := testCommonRegistry(t, map[string][]byte{
		"api/schemas/_common/typed-ref.json": []byte(`{"type":"object"}`),
	})
	r, err := NewLocalRefResolver(reg, DefaultMaxRefDepth)
	if err != nil {
		t.Fatalf("NewLocalRefResolver: %v", err)
	}

	_, _, err = r.Resolve("api/schemas/project.json", "_common/typed-ref.json#/properties/name")
	if !errors.Is(err, ErrRefRejected) {
		t.Fatalf("fragment ref: err=%v, want ErrRefRejected", err)
	}
	if err == nil || !strings.Contains(err.Error(), "fragments") {
		t.Fatalf("fragment error should mention fragments: %v", err)
	}
}

func TestLocalRefResolver_ReturnedBytesDefensiveCopy(t *testing.T) {
	t.Parallel()

	const id = "api/schemas/_common/typed-ref.json"
	reg := testCommonRegistry(t, map[string][]byte{
		id: []byte(`{"type":"object"}`),
	})
	r, err := NewLocalRefResolver(reg, DefaultMaxRefDepth)
	if err != nil {
		t.Fatalf("NewLocalRefResolver: %v", err)
	}

	_, body, err := r.Resolve("api/schemas/project.json", "_common/typed-ref.json")
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	body[0] = 'X'
	_, again, err := r.Resolve("api/schemas/project.json", "_common/typed-ref.json")
	if err != nil {
		t.Fatalf("second Resolve: %v", err)
	}
	if again[0] == 'X' {
		t.Fatal("returned schema bytes alias registry storage")
	}
}
