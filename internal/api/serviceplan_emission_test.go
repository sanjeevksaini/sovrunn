package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/registry"
	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

func TestServicePlanEmission_MutatingPaths(t *testing.T) {
	updateBody := map[string]any{
		"metadata": map[string]any{"name": "small"},
		"spec": map[string]any{
			"serviceClassName": "postgres",
			"tier":             resources.TierLarge,
			"lifecycle":        resources.LifecycleDeprecated,
			"description":      "updated",
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
			name: "Create",
			setup: func(env *catalogEmissionEnv) {
				seedServiceClassDirect(t, env.scRegistry, "postgres")
			},
			action: func(env *catalogEmissionEnv) *httptest.ResponseRecorder {
				return createServicePlan(env.spHandler, "postgres", "small")
			},
			wantStatus: http.StatusCreated,
			wantSpec: resources.OperationSpec{
				Type:             resources.OpCreateServicePlan,
				ResourceKind:     resources.ServicePlanKind,
				ResourceName:     "small",
				ServiceClassName: "postgres",
				ServicePlanName:  "small",
			},
		},
		{
			name: "Update",
			setup: func(env *catalogEmissionEnv) {
				seedServiceClassDirect(t, env.scRegistry, "postgres")
				seedServicePlanDirect(t, env.spRegistry, "postgres", "small")
			},
			action: func(env *catalogEmissionEnv) *httptest.ResponseRecorder {
				req := jsonRequest(http.MethodPut, "/v1/service-plans/postgres/small", updateBody, "application/json")
				rec := httptest.NewRecorder()
				env.spHandler.HandleItem(rec, req)
				return rec
			},
			wantStatus: http.StatusOK,
			wantSpec: resources.OperationSpec{
				Type:             resources.OpUpdateServicePlan,
				ResourceKind:     resources.ServicePlanKind,
				ResourceName:     "small",
				ServiceClassName: "postgres",
				ServicePlanName:  "small",
			},
		},
		{
			name: "Delete",
			setup: func(env *catalogEmissionEnv) {
				seedServiceClassDirect(t, env.scRegistry, "postgres")
				seedServicePlanDirect(t, env.spRegistry, "postgres", "small")
			},
			action: func(env *catalogEmissionEnv) *httptest.ResponseRecorder {
				req := jsonRequest(http.MethodDelete, "/v1/service-plans/postgres/small", nil, "")
				rec := httptest.NewRecorder()
				env.spHandler.HandleItem(rec, req)
				return rec
			},
			wantStatus: http.StatusNoContent,
			wantSpec: resources.OperationSpec{
				Type:             resources.OpDeleteServicePlan,
				ResourceKind:     resources.ServicePlanKind,
				ResourceName:     "small",
				ServiceClassName: "postgres",
				ServicePlanName:  "small",
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

func TestServicePlanEmission_NoEmissionOnFailedAction(t *testing.T) {
	t.Run("validation failure", func(t *testing.T) {
		env := newCatalogEmissionEnv(t)
		seedServiceClassDirect(t, env.scRegistry, "postgres")
		body := map[string]any{
			"metadata": map[string]any{"name": "INVALID"},
			"spec": map[string]any{
				"serviceClassName": "postgres",
				"tier":             resources.TierSmall,
				"lifecycle":        resources.LifecycleActive,
			},
		}
		req := jsonRequest(http.MethodPost, "/v1/service-plans", body, "application/json")
		rec := httptest.NewRecorder()
		env.spHandler.HandleCollection(rec, req)
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
		seedServicePlanDirect(t, env.spRegistry, "postgres", "small")
		rec := createServicePlan(env.spHandler, "postgres", "small")
		if rec.Code != http.StatusConflict {
			t.Fatalf("status = %d, want 409", rec.Code)
		}
		if ops := catalogRecordedOperations(t, env); len(ops) != 0 {
			t.Errorf("operations = %d, want 0", len(ops))
		}
	})

	t.Run("missing parent on create", func(t *testing.T) {
		env := newCatalogEmissionEnv(t)
		rec := createServicePlan(env.spHandler, "postgres", "small")
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400", rec.Code)
		}
		if ops := catalogRecordedOperations(t, env); len(ops) != 0 {
			t.Errorf("operations = %d, want 0", len(ops))
		}
	})

	t.Run("update missing target", func(t *testing.T) {
		env := newCatalogEmissionEnv(t)
		seedServiceClassDirect(t, env.scRegistry, "postgres")
		req := jsonRequest(http.MethodPut, "/v1/service-plans/postgres/missing", validServicePlanBody("postgres", "missing"), "application/json")
		rec := httptest.NewRecorder()
		env.spHandler.HandleItem(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want 404", rec.Code)
		}
		if ops := catalogRecordedOperations(t, env); len(ops) != 0 {
			t.Errorf("operations = %d, want 0", len(ops))
		}
	})

	t.Run("delete missing target", func(t *testing.T) {
		env := newCatalogEmissionEnv(t)
		req := jsonRequest(http.MethodDelete, "/v1/service-plans/postgres/missing", nil, "")
		rec := httptest.NewRecorder()
		env.spHandler.HandleItem(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want 404", rec.Code)
		}
		if ops := catalogRecordedOperations(t, env); len(ops) != 0 {
			t.Errorf("operations = %d, want 0", len(ops))
		}
	})
}

func TestServicePlanEmission_FailureIsolation(t *testing.T) {
	newHandler := func() (*ServicePlanHandler, *registry.ServiceClassRegistry, *registry.ServicePlanRegistry) {
		scReg := registry.NewServiceClassRegistry()
		spReg := registry.NewServicePlanRegistry()
		return NewServicePlanHandler(spReg, scReg, &failingEmitter{}), scReg, spReg
	}

	t.Run("create still returns 201", func(t *testing.T) {
		h, scReg, _ := newHandler()
		seedServiceClassDirect(t, scReg, "postgres")
		rec := createServicePlan(h, "postgres", "small")
		if rec.Code != http.StatusCreated {
			t.Fatalf("status = %d, want 201; body=%s", rec.Code, rec.Body.String())
		}
	})

	t.Run("update still returns 200", func(t *testing.T) {
		h, scReg, spReg := newHandler()
		seedServiceClassDirect(t, scReg, "postgres")
		seedServicePlanDirect(t, spReg, "postgres", "small")
		req := jsonRequest(http.MethodPut, "/v1/service-plans/postgres/small", validServicePlanBody("postgres", "small"), "application/json")
		rec := httptest.NewRecorder()
		h.HandleItem(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
		}
	})

	t.Run("delete still returns 204", func(t *testing.T) {
		h, scReg, spReg := newHandler()
		seedServiceClassDirect(t, scReg, "postgres")
		seedServicePlanDirect(t, spReg, "postgres", "small")
		req := jsonRequest(http.MethodDelete, "/v1/service-plans/postgres/small", nil, "")
		rec := httptest.NewRecorder()
		h.HandleItem(rec, req)
		if rec.Code != http.StatusNoContent {
			t.Fatalf("status = %d, want 204; body=%s", rec.Code, rec.Body.String())
		}
	})
}
