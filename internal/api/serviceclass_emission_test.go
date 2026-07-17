package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/registry"
	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// catalogEmissionEnv wires ServiceClass/ServicePlan handlers to a shared
// OperationRegistry-backed emitter for catalog emission assertions.
type catalogEmissionEnv struct {
	scRegistry *registry.ServiceClassRegistry
	spRegistry *registry.ServicePlanRegistry
	opRegistry *registry.OperationRegistry
	scHandler  *ServiceClassHandler
	spHandler  *ServicePlanHandler
}

func newCatalogEmissionEnv(t *testing.T) *catalogEmissionEnv {
	t.Helper()
	scReg := registry.NewServiceClassRegistry()
	spReg := registry.NewServicePlanRegistry()
	opReg := registry.NewOperationRegistry()
	emitter := NewRegistryEmitter(opReg, nil)
	blocker := registry.NewServicePlanChildBlockerChecker(spReg)
	return &catalogEmissionEnv{
		scRegistry: scReg,
		spRegistry: spReg,
		opRegistry: opReg,
		scHandler:  NewServiceClassHandler(scReg, blocker, emitter),
		spHandler:  NewServicePlanHandler(spReg, scReg, emitter),
	}
}

func catalogRecordedOperations(t *testing.T, env *catalogEmissionEnv) []resources.Operation {
	t.Helper()
	items, err := env.opRegistry.ListOperations(context.Background())
	if err != nil {
		t.Fatalf("ListOperations() error = %v", err)
	}
	return items
}

func catalogRequireSingleOperation(t *testing.T, env *catalogEmissionEnv) resources.Operation {
	t.Helper()
	items := catalogRecordedOperations(t, env)
	if len(items) != 1 {
		t.Fatalf("recorded operations = %d, want exactly 1", len(items))
	}
	return items[0]
}

func assertCatalogOperation(t *testing.T, op resources.Operation, want resources.OperationSpec) {
	t.Helper()
	assertOperation(t, op, want)
	if op.Spec.ServiceClassName != want.ServiceClassName {
		t.Errorf("Spec.ServiceClassName = %q, want %q", op.Spec.ServiceClassName, want.ServiceClassName)
	}
	if op.Spec.ServicePlanName != want.ServicePlanName {
		t.Errorf("Spec.ServicePlanName = %q, want %q", op.Spec.ServicePlanName, want.ServicePlanName)
	}
}

func seedServiceClassDirect(t *testing.T, reg *registry.ServiceClassRegistry, name string) {
	t.Helper()
	sc := resources.ServiceClass{
		APIVersion: resources.ServiceClassAPIVersion,
		Kind:       resources.ServiceClassKind,
		Metadata:   resources.Metadata{Name: name},
		Spec: resources.ServiceClassSpec{
			Category:  resources.CategoryDatabase,
			Lifecycle: resources.LifecycleActive,
		},
		Status: resources.ServiceClassStatus{Phase: resources.PhaseActive},
	}
	if _, err := reg.CreateServiceClass(context.Background(), sc); err != nil {
		t.Fatalf("seedServiceClassDirect(%s): %v", name, err)
	}
}

func seedServicePlanDirect(t *testing.T, reg *registry.ServicePlanRegistry, serviceClassName, name string) {
	t.Helper()
	sp := resources.ServicePlan{
		APIVersion: resources.ServicePlanAPIVersion,
		Kind:       resources.ServicePlanKind,
		Metadata:   resources.Metadata{Name: name},
		Spec: resources.ServicePlanSpec{
			ServiceClassName: serviceClassName,
			Tier:             resources.TierSmall,
			Lifecycle:        resources.LifecycleActive,
		},
		Status: resources.ServicePlanStatus{Phase: resources.PhaseActive},
	}
	if _, err := reg.CreateServicePlan(context.Background(), sp); err != nil {
		t.Fatalf("seedServicePlanDirect(%s/%s): %v", serviceClassName, name, err)
	}
}

func TestServiceClassEmission_MutatingPaths(t *testing.T) {
	updateBody := map[string]any{
		"metadata": map[string]any{"name": "postgres"},
		"spec": map[string]any{
			"category":    resources.CategoryDatabase,
			"lifecycle":   resources.LifecycleDeprecated,
			"description": "updated",
		},
	}

	cases := []struct {
		name       string
		setup      func(env *catalogEmissionEnv)
		action     func(env *catalogEmissionEnv) *httptest.ResponseRecorder
		wantStatus int
		wantSpec   resources.OperationSpec
	}{
		{
			name:  "Create",
			setup: func(env *catalogEmissionEnv) {},
			action: func(env *catalogEmissionEnv) *httptest.ResponseRecorder {
				return createServiceClass(env.scHandler, "postgres")
			},
			wantStatus: http.StatusCreated,
			wantSpec: resources.OperationSpec{
				Type:             resources.OpCreateServiceClass,
				ResourceKind:     resources.ServiceClassKind,
				ResourceName:     "postgres",
				ServiceClassName: "postgres",
			},
		},
		{
			name: "Update",
			setup: func(env *catalogEmissionEnv) {
				seedServiceClassDirect(t, env.scRegistry, "postgres")
			},
			action: func(env *catalogEmissionEnv) *httptest.ResponseRecorder {
				req := jsonRequest(http.MethodPut, "/v1/service-classes/postgres", updateBody, "application/json")
				rec := httptest.NewRecorder()
				env.scHandler.HandleItem(rec, req)
				return rec
			},
			wantStatus: http.StatusOK,
			wantSpec: resources.OperationSpec{
				Type:             resources.OpUpdateServiceClass,
				ResourceKind:     resources.ServiceClassKind,
				ResourceName:     "postgres",
				ServiceClassName: "postgres",
			},
		},
		{
			name: "Delete",
			setup: func(env *catalogEmissionEnv) {
				seedServiceClassDirect(t, env.scRegistry, "postgres")
			},
			action: func(env *catalogEmissionEnv) *httptest.ResponseRecorder {
				req := jsonRequest(http.MethodDelete, "/v1/service-classes/postgres", nil, "")
				rec := httptest.NewRecorder()
				env.scHandler.HandleItem(rec, req)
				return rec
			},
			wantStatus: http.StatusNoContent,
			wantSpec: resources.OperationSpec{
				Type:             resources.OpDeleteServiceClass,
				ResourceKind:     resources.ServiceClassKind,
				ResourceName:     "postgres",
				ServiceClassName: "postgres",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			env := newCatalogEmissionEnv(t)
			tc.setup(env)
			rec := tc.action(env)
			if rec.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d; body=%s", rec.Code, tc.wantStatus, rec.Body.String())
			}
			op := catalogRequireSingleOperation(t, env)
			assertCatalogOperation(t, op, tc.wantSpec)
		})
	}
}

func TestServiceClassEmission_NoEmissionOnFailedAction(t *testing.T) {
	t.Run("validation failure", func(t *testing.T) {
		env := newCatalogEmissionEnv(t)
		body := map[string]any{
			"metadata": map[string]any{"name": "INVALID"},
			"spec": map[string]any{
				"category":  resources.CategoryDatabase,
				"lifecycle": resources.LifecycleActive,
			},
		}
		req := jsonRequest(http.MethodPost, "/v1/service-classes", body, "application/json")
		rec := httptest.NewRecorder()
		env.scHandler.HandleCollection(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400", rec.Code)
		}
		if ops := catalogRecordedOperations(t, env); len(ops) != 0 {
			t.Errorf("operations = %d, want 0", len(ops))
		}
	})

	t.Run("duplicate create", func(t *testing.T) {
		env := newCatalogEmissionEnv(t)
		seedServiceClassDirect(t, env.scRegistry, "postgres")
		rec := createServiceClass(env.scHandler, "postgres")
		if rec.Code != http.StatusConflict {
			t.Fatalf("status = %d, want 409", rec.Code)
		}
		if ops := catalogRecordedOperations(t, env); len(ops) != 0 {
			t.Errorf("operations = %d, want 0", len(ops))
		}
	})

	t.Run("update missing target", func(t *testing.T) {
		env := newCatalogEmissionEnv(t)
		req := jsonRequest(http.MethodPut, "/v1/service-classes/missing", validServiceClassBody("missing"), "application/json")
		rec := httptest.NewRecorder()
		env.scHandler.HandleItem(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want 404", rec.Code)
		}
		if ops := catalogRecordedOperations(t, env); len(ops) != 0 {
			t.Errorf("operations = %d, want 0", len(ops))
		}
	})

	t.Run("delete missing target", func(t *testing.T) {
		env := newCatalogEmissionEnv(t)
		req := jsonRequest(http.MethodDelete, "/v1/service-classes/missing", nil, "")
		rec := httptest.NewRecorder()
		env.scHandler.HandleItem(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want 404", rec.Code)
		}
		if ops := catalogRecordedOperations(t, env); len(ops) != 0 {
			t.Errorf("operations = %d, want 0", len(ops))
		}
	})

	t.Run("delete blocked", func(t *testing.T) {
		env := newCatalogEmissionEnv(t)
		seedServiceClassDirect(t, env.scRegistry, "postgres")
		seedServicePlanDirect(t, env.spRegistry, "postgres", "small")
		req := jsonRequest(http.MethodDelete, "/v1/service-classes/postgres", nil, "")
		rec := httptest.NewRecorder()
		env.scHandler.HandleItem(rec, req)
		if rec.Code != http.StatusConflict {
			t.Fatalf("status = %d, want 409; body=%s", rec.Code, rec.Body.String())
		}
		if ops := catalogRecordedOperations(t, env); len(ops) != 0 {
			t.Errorf("operations = %d, want 0", len(ops))
		}
	})
}

func TestServiceClassEmission_FailureIsolation(t *testing.T) {
	newHandler := func() (*ServiceClassHandler, *registry.ServiceClassRegistry) {
		scReg := registry.NewServiceClassRegistry()
		spReg := registry.NewServicePlanRegistry()
		blocker := registry.NewServicePlanChildBlockerChecker(spReg)
		return NewServiceClassHandler(scReg, blocker, &failingEmitter{}), scReg
	}

	t.Run("create still returns 201", func(t *testing.T) {
		h, _ := newHandler()
		rec := createServiceClass(h, "postgres")
		if rec.Code != http.StatusCreated {
			t.Fatalf("status = %d, want 201; body=%s", rec.Code, rec.Body.String())
		}
	})

	t.Run("update still returns 200", func(t *testing.T) {
		h, scReg := newHandler()
		seedServiceClassDirect(t, scReg, "postgres")
		req := jsonRequest(http.MethodPut, "/v1/service-classes/postgres", validServiceClassBody("postgres"), "application/json")
		rec := httptest.NewRecorder()
		h.HandleItem(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
		}
	})

	t.Run("delete still returns 204", func(t *testing.T) {
		h, scReg := newHandler()
		seedServiceClassDirect(t, scReg, "postgres")
		req := jsonRequest(http.MethodDelete, "/v1/service-classes/postgres", nil, "")
		rec := httptest.NewRecorder()
		h.HandleItem(rec, req)
		if rec.Code != http.StatusNoContent {
			t.Fatalf("status = %d, want 204; body=%s", rec.Code, rec.Body.String())
		}
	})
}
