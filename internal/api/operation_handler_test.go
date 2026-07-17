package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/registry"
	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

func testOperation(id string, createdAt string) resources.Operation {
	return resources.Operation{
		APIVersion: resources.OperationAPIVersion,
		Kind:       resources.OperationKind,
		Metadata:   resources.Metadata{Name: id},
		Spec: resources.OperationSpec{
			Type:                 resources.OpCreateProject,
			ResourceKind:         resources.ProjectKind,
			ResourceName:         "project-a",
			OrganizationName:     "org-a",
			OrganizationUnitName: "ou-a",
			TenantName:           "tenant-a",
			ProjectName:          "project-a",
			Actor:                "system",
		},
		Status: resources.OperationStatus{
			Phase:       resources.OperationPhaseSucceeded,
			CreatedAt:   createdAt,
			UpdatedAt:   createdAt,
			CompletedAt: createdAt,
		},
	}
}

func newTestOperationHandler() (*OperationHandler, *registry.OperationRegistry) {
	reg := registry.NewOperationRegistry()
	return NewOperationHandler(reg), reg
}

func decodeOperationItems(t *testing.T, rec *httptest.ResponseRecorder) []resources.Operation {
	t.Helper()
	var raw map[string]json.RawMessage
	if err := json.NewDecoder(rec.Body).Decode(&raw); err != nil {
		t.Fatalf("decode list response: %v", err)
	}
	if len(raw) != 1 {
		t.Fatalf("top-level keys = %d, want only items; body=%v", len(raw), raw)
	}
	itemsRaw, ok := raw["items"]
	if !ok {
		t.Fatalf("response missing items field: %v", raw)
	}
	if string(itemsRaw) == "null" {
		t.Fatal("items is null, want JSON array")
	}
	var items []resources.Operation
	if err := json.Unmarshal(itemsRaw, &items); err != nil {
		t.Fatalf("decode items: %v", err)
	}
	return items
}

func TestOperationHandler_List_Empty(t *testing.T) {
	h, _ := newTestOperationHandler()

	req := jsonRequest(http.MethodGet, "/v1/operations", nil, "")
	rec := httptest.NewRecorder()
	h.HandleCollection(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}

	items := decodeOperationItems(t, rec)
	if items == nil {
		t.Fatal("decoded items is nil, want non-nil empty slice")
	}
	if len(items) != 0 {
		t.Fatalf("items len = %d, want 0", len(items))
	}
}

func TestOperationHandler_List_Sorted(t *testing.T) {
	h, reg := newTestOperationHandler()

	inputs := []resources.Operation{
		testOperation("c", "2026-01-02T00:00:00Z"),
		testOperation("a", "2026-01-03T00:00:00Z"),
		testOperation("b", "2026-01-02T00:00:00Z"),
		testOperation("d", "2026-01-01T00:00:00Z"),
	}
	for _, op := range inputs {
		if _, err := reg.CreateOperation(context.Background(), op); err != nil {
			t.Fatalf("CreateOperation(%s) error = %v", op.Metadata.Name, err)
		}
	}

	req := jsonRequest(http.MethodGet, "/v1/operations", nil, "")
	rec := httptest.NewRecorder()
	h.HandleCollection(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}

	items := decodeOperationItems(t, rec)
	wantOrder := []string{"d", "b", "c", "a"}
	if len(items) != len(wantOrder) {
		t.Fatalf("items len = %d, want %d", len(items), len(wantOrder))
	}
	for i, want := range wantOrder {
		if items[i].Metadata.Name != want {
			t.Errorf("items[%d].metadata.name = %q, want %q", i, items[i].Metadata.Name, want)
		}
	}
}

func TestOperationHandler_Get_Existing(t *testing.T) {
	h, reg := newTestOperationHandler()
	op := testOperation("op-1", "2026-01-01T00:00:00Z")
	if _, err := reg.CreateOperation(context.Background(), op); err != nil {
		t.Fatalf("CreateOperation() error = %v", err)
	}

	req := jsonRequest(http.MethodGet, "/v1/operations/op-1", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	var got resources.Operation
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode operation: %v", err)
	}
	if got.Metadata.Name != "op-1" {
		t.Errorf("metadata.name = %q, want op-1", got.Metadata.Name)
	}
	if got.APIVersion != resources.OperationAPIVersion || got.Kind != resources.OperationKind ||
		got.Spec != op.Spec || got.Status != op.Status {
		t.Errorf("operation = %+v, want full stored resource %+v", got, op)
	}
}

func TestOperationHandler_Get_Missing(t *testing.T) {
	h, _ := newTestOperationHandler()

	req := jsonRequest(http.MethodGet, "/v1/operations/missing", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404; body=%s", rec.Code, rec.Body.String())
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeResourceNotFound {
		t.Errorf("error.code = %q, want RESOURCE_NOT_FOUND", errBody.Code)
	}
}

func TestOperationHandler_Get_BareItemPath(t *testing.T) {
	h, _ := newTestOperationHandler()

	req := jsonRequest(http.MethodGet, "/v1/operations/", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404; body=%s", rec.Code, rec.Body.String())
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeResourceNotFound {
		t.Errorf("error.code = %q, want RESOURCE_NOT_FOUND", errBody.Code)
	}
}

func TestOperationHandler_Get_ExtraSegmentPath(t *testing.T) {
	h, _ := newTestOperationHandler()

	req := jsonRequest(http.MethodGet, "/v1/operations/op-1/extra", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404; body=%s", rec.Code, rec.Body.String())
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeResourceNotFound {
		t.Errorf("error.code = %q, want RESOURCE_NOT_FOUND", errBody.Code)
	}
}

func TestOperationHandler_CollectionPost_MethodNotAllowed(t *testing.T) {
	h, _ := newTestOperationHandler()

	req := jsonRequest(http.MethodPost, "/v1/operations", nil, "application/json")
	rec := httptest.NewRecorder()
	h.HandleCollection(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want 405; body=%s", rec.Code, rec.Body.String())
	}
}

func TestOperationHandler_ItemUnsupportedMethod(t *testing.T) {
	h, _ := newTestOperationHandler()

	req := jsonRequest(http.MethodDelete, "/v1/operations/op-1", nil, "")
	rec := httptest.NewRecorder()
	h.HandleItem(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want 405; body=%s", rec.Code, rec.Body.String())
	}
}
