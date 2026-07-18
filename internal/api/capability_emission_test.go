package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/registry"
	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

func TestCapabilityEmission_MutatingPaths(t *testing.T) {
	cases := []struct {
		name       string
		setup      func(env *pluginCapabilityEmissionEnv)
		action     func(env *pluginCapabilityEmissionEnv) *httptest.ResponseRecorder
		wantStatus int
		wantSpec   resources.OperationSpec
	}{
		{
			name: "Create",
			setup: func(env *pluginCapabilityEmissionEnv) {
				seedPluginServiceClass(t, env.scRegistry, "postgres")
				seedPluginDirect(t, env.pluginRegistry, "pg-plugin")
			},
			action: func(env *pluginCapabilityEmissionEnv) *httptest.ResponseRecorder {
				return createCapability(env.capabilityHandler, "pg-provision", "pg-plugin", "postgres")
			},
			wantStatus: http.StatusCreated,
			wantSpec: resources.OperationSpec{
				Type:           resources.OpCreateCapability,
				ResourceKind:   resources.CapabilityKind,
				ResourceName:   "pg-provision",
				PluginName:     "pg-plugin",
				CapabilityName: "pg-provision",
			},
		},
		{
			name: "Delete",
			setup: func(env *pluginCapabilityEmissionEnv) {
				seedPluginServiceClass(t, env.scRegistry, "postgres")
				seedPluginDirect(t, env.pluginRegistry, "pg-plugin")
				seedCapabilityDirect(t, env.capabilityRegistry, "pg-provision", "pg-plugin", "postgres")
			},
			action: func(env *pluginCapabilityEmissionEnv) *httptest.ResponseRecorder {
				req := jsonRequest(http.MethodDelete, "/v1/capabilities/pg-provision", nil, "")
				rec := httptest.NewRecorder()
				env.capabilityHandler.HandleItem(rec, req)
				return rec
			},
			wantStatus: http.StatusNoContent,
			wantSpec: resources.OperationSpec{
				Type:           resources.OpDeleteCapability,
				ResourceKind:   resources.CapabilityKind,
				ResourceName:   "pg-provision",
				PluginName:     "pg-plugin",
				CapabilityName: "pg-provision",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			env := newPluginCapabilityEmissionEnv(t)
			tc.setup(env)
			rec := tc.action(env)
			if rec.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d; body=%s", rec.Code, tc.wantStatus, rec.Body.String())
			}
			op := pluginCapabilityRequireSingleOperation(t, env)
			assertPluginCapabilityOperation(t, op, tc.wantSpec)
		})
	}
}

func TestCapabilityEmission_NoEmissionOnFailedAction(t *testing.T) {
	t.Run("validation failure", func(t *testing.T) {
		env := newPluginCapabilityEmissionEnv(t)
		seedPluginServiceClass(t, env.scRegistry, "postgres")
		seedPluginDirect(t, env.pluginRegistry, "pg-plugin")
		body := map[string]any{
			"metadata": map[string]any{"name": "INVALID"},
			"spec": map[string]any{
				"pluginRef":       "pg-plugin",
				"serviceClassRef": "postgres",
				"operation":       resources.CapOpProvision,
				"supported":       true,
			},
		}
		req := jsonRequest(http.MethodPost, "/v1/capabilities", body, "application/json")
		rec := httptest.NewRecorder()
		env.capabilityHandler.HandleCollection(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400", rec.Code)
		}
		if ops := pluginCapabilityRecordedOperations(t, env); len(ops) != 0 {
			t.Errorf("operations = %d, want 0", len(ops))
		}
	})

	t.Run("duplicate create", func(t *testing.T) {
		env := newPluginCapabilityEmissionEnv(t)
		seedPluginServiceClass(t, env.scRegistry, "postgres")
		seedPluginDirect(t, env.pluginRegistry, "pg-plugin")
		seedCapabilityDirect(t, env.capabilityRegistry, "pg-provision", "pg-plugin", "postgres")
		rec := createCapability(env.capabilityHandler, "pg-provision", "pg-plugin", "postgres")
		if rec.Code != http.StatusConflict {
			t.Fatalf("status = %d, want 409", rec.Code)
		}
		if ops := pluginCapabilityRecordedOperations(t, env); len(ops) != 0 {
			t.Errorf("operations = %d, want 0", len(ops))
		}
	})

	t.Run("missing plugin ref", func(t *testing.T) {
		env := newPluginCapabilityEmissionEnv(t)
		seedPluginServiceClass(t, env.scRegistry, "postgres")
		rec := createCapability(env.capabilityHandler, "pg-provision", "pg-plugin", "postgres")
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400", rec.Code)
		}
		if ops := pluginCapabilityRecordedOperations(t, env); len(ops) != 0 {
			t.Errorf("operations = %d, want 0", len(ops))
		}
	})

	t.Run("missing service class ref", func(t *testing.T) {
		env := newPluginCapabilityEmissionEnv(t)
		seedPluginDirect(t, env.pluginRegistry, "pg-plugin")
		rec := createCapability(env.capabilityHandler, "pg-provision", "pg-plugin", "postgres")
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400", rec.Code)
		}
		if ops := pluginCapabilityRecordedOperations(t, env); len(ops) != 0 {
			t.Errorf("operations = %d, want 0", len(ops))
		}
	})

	t.Run("service class not declared by plugin", func(t *testing.T) {
		env := newPluginCapabilityEmissionEnv(t)
		seedPluginServiceClass(t, env.scRegistry, "postgres")
		seedPluginServiceClass(t, env.scRegistry, "redis")
		seedPluginDirect(t, env.pluginRegistry, "pg-plugin", "postgres")
		rec := createCapability(env.capabilityHandler, "pg-provision", "pg-plugin", "redis")
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400", rec.Code)
		}
		if ops := pluginCapabilityRecordedOperations(t, env); len(ops) != 0 {
			t.Errorf("operations = %d, want 0", len(ops))
		}
	})

	t.Run("delete missing target", func(t *testing.T) {
		env := newPluginCapabilityEmissionEnv(t)
		req := jsonRequest(http.MethodDelete, "/v1/capabilities/missing", nil, "")
		rec := httptest.NewRecorder()
		env.capabilityHandler.HandleItem(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want 404", rec.Code)
		}
		if ops := pluginCapabilityRecordedOperations(t, env); len(ops) != 0 {
			t.Errorf("operations = %d, want 0", len(ops))
		}
	})

	t.Run("method not allowed on put", func(t *testing.T) {
		env := newPluginCapabilityEmissionEnv(t)
		seedPluginServiceClass(t, env.scRegistry, "postgres")
		seedPluginDirect(t, env.pluginRegistry, "pg-plugin")
		seedCapabilityDirect(t, env.capabilityRegistry, "pg-provision", "pg-plugin", "postgres")
		req := jsonRequest(
			http.MethodPut,
			"/v1/capabilities/pg-provision",
			validCapabilityBody("pg-provision", "pg-plugin", "postgres"),
			"application/json",
		)
		rec := httptest.NewRecorder()
		env.capabilityHandler.HandleItem(rec, req)
		if rec.Code != http.StatusMethodNotAllowed {
			t.Fatalf("status = %d, want 405", rec.Code)
		}
		if ops := pluginCapabilityRecordedOperations(t, env); len(ops) != 0 {
			t.Errorf("operations = %d, want 0", len(ops))
		}
	})
}

func TestCapabilityEmission_FailureIsolation(t *testing.T) {
	newHandler := func() (
		*CapabilityHandler,
		*registry.ServiceClassRegistry,
		*registry.PluginRegistry,
		*registry.CapabilityRegistry,
	) {
		scReg := registry.NewServiceClassRegistry()
		pluginReg := registry.NewPluginRegistry()
		capReg := registry.NewCapabilityRegistry()
		return NewCapabilityHandler(capReg, pluginReg, scReg, &failingEmitter{}), scReg, pluginReg, capReg
	}

	t.Run("create still returns 201", func(t *testing.T) {
		h, scReg, pluginReg, _ := newHandler()
		seedPluginServiceClass(t, scReg, "postgres")
		seedPluginDirect(t, pluginReg, "pg-plugin")
		rec := createCapability(h, "pg-provision", "pg-plugin", "postgres")
		if rec.Code != http.StatusCreated {
			t.Fatalf("status = %d, want 201; body=%s", rec.Code, rec.Body.String())
		}
	})

	t.Run("delete still returns 204", func(t *testing.T) {
		h, scReg, pluginReg, capReg := newHandler()
		seedPluginServiceClass(t, scReg, "postgres")
		seedPluginDirect(t, pluginReg, "pg-plugin")
		seedCapabilityDirect(t, capReg, "pg-provision", "pg-plugin", "postgres")
		req := jsonRequest(http.MethodDelete, "/v1/capabilities/pg-provision", nil, "")
		rec := httptest.NewRecorder()
		h.HandleItem(rec, req)
		if rec.Code != http.StatusNoContent {
			t.Fatalf("status = %d, want 204; body=%s", rec.Code, rec.Body.String())
		}
	})
}
