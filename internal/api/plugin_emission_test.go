package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/registry"
	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// pluginCapabilityEmissionEnv wires Plugin/Capability handlers to a shared
// OperationRegistry-backed emitter for emission assertions.
type pluginCapabilityEmissionEnv struct {
	scRegistry         *registry.ServiceClassRegistry
	pluginRegistry     *registry.PluginRegistry
	capabilityRegistry *registry.CapabilityRegistry
	opRegistry         *registry.OperationRegistry
	pluginHandler      *PluginHandler
	capabilityHandler  *CapabilityHandler
}

func newPluginCapabilityEmissionEnv(t *testing.T) *pluginCapabilityEmissionEnv {
	t.Helper()
	scReg := registry.NewServiceClassRegistry()
	pluginReg := registry.NewPluginRegistry()
	capReg := registry.NewCapabilityRegistry()
	opReg := registry.NewOperationRegistry()
	emitter := NewRegistryEmitter(opReg, nil)
	blocker := registry.NewCapabilityChildBlockerChecker(capReg)
	return &pluginCapabilityEmissionEnv{
		scRegistry:         scReg,
		pluginRegistry:     pluginReg,
		capabilityRegistry: capReg,
		opRegistry:         opReg,
		pluginHandler:      NewPluginHandler(pluginReg, scReg, blocker, emitter),
		capabilityHandler:  NewCapabilityHandler(capReg, pluginReg, scReg, emitter),
	}
}

func pluginCapabilityRecordedOperations(t *testing.T, env *pluginCapabilityEmissionEnv) []resources.Operation {
	t.Helper()
	items, err := env.opRegistry.ListOperations(context.Background())
	if err != nil {
		t.Fatalf("ListOperations() error = %v", err)
	}
	return items
}

func pluginCapabilityRequireSingleOperation(t *testing.T, env *pluginCapabilityEmissionEnv) resources.Operation {
	t.Helper()
	items := pluginCapabilityRecordedOperations(t, env)
	if len(items) != 1 {
		t.Fatalf("recorded operations = %d, want exactly 1", len(items))
	}
	return items[0]
}

func assertPluginCapabilityOperation(t *testing.T, op resources.Operation, want resources.OperationSpec) {
	t.Helper()
	assertOperation(t, op, want)
	if op.Spec.PluginName != want.PluginName {
		t.Errorf("Spec.PluginName = %q, want %q", op.Spec.PluginName, want.PluginName)
	}
	if op.Spec.CapabilityName != want.CapabilityName {
		t.Errorf("Spec.CapabilityName = %q, want %q", op.Spec.CapabilityName, want.CapabilityName)
	}
}

func seedPluginDirect(t *testing.T, reg *registry.PluginRegistry, name string, serviceClassRefs ...string) {
	t.Helper()
	if serviceClassRefs == nil {
		serviceClassRefs = []string{"postgres"}
	}
	p := resources.Plugin{
		APIVersion: resources.PluginAPIVersion,
		Kind:       resources.PluginKind,
		Metadata:   resources.Metadata{Name: name},
		Spec: resources.PluginSpec{
			PluginType:       resources.PluginTypeDStoreOps,
			Version:          "1.0.0",
			ServiceClassRefs: serviceClassRefs,
			DeploymentMode:   resources.DeploymentModeCompiledIn,
		},
		Status: resources.PluginStatus{Phase: resources.PhaseActive},
	}
	if _, err := reg.CreatePlugin(context.Background(), p); err != nil {
		t.Fatalf("seedPluginDirect(%s): %v", name, err)
	}
}

func seedCapabilityDirect(
	t *testing.T,
	reg *registry.CapabilityRegistry,
	name, pluginRef, serviceClassRef string,
) {
	t.Helper()
	c := resources.Capability{
		APIVersion: resources.CapabilityAPIVersion,
		Kind:       resources.CapabilityKind,
		Metadata:   resources.Metadata{Name: name},
		Spec: resources.CapabilitySpec{
			PluginRef:       pluginRef,
			ServiceClassRef: serviceClassRef,
			Operation:       resources.CapOpProvision,
			Supported:       true,
		},
		Status: resources.CapabilityStatus{Phase: resources.PhaseActive},
	}
	if _, err := reg.CreateCapability(context.Background(), c); err != nil {
		t.Fatalf("seedCapabilityDirect(%s): %v", name, err)
	}
}

func TestPluginEmission_MutatingPaths(t *testing.T) {
	updateBody := map[string]any{
		"metadata": map[string]any{"name": "pg-plugin"},
		"spec": map[string]any{
			"pluginType":       resources.PluginTypeDStoreOps,
			"version":          "2.0.0",
			"serviceClassRefs": []string{"postgres"},
			"deploymentMode":   resources.DeploymentModeCompiledIn,
			"description":      "updated",
		},
	}

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
			},
			action: func(env *pluginCapabilityEmissionEnv) *httptest.ResponseRecorder {
				return createPlugin(env.pluginHandler, "pg-plugin")
			},
			wantStatus: http.StatusCreated,
			wantSpec: resources.OperationSpec{
				Type:         resources.OpCreatePlugin,
				ResourceKind: resources.PluginKind,
				ResourceName: "pg-plugin",
				PluginName:   "pg-plugin",
			},
		},
		{
			name: "Update",
			setup: func(env *pluginCapabilityEmissionEnv) {
				seedPluginServiceClass(t, env.scRegistry, "postgres")
				seedPluginDirect(t, env.pluginRegistry, "pg-plugin")
			},
			action: func(env *pluginCapabilityEmissionEnv) *httptest.ResponseRecorder {
				req := jsonRequest(http.MethodPut, "/v1/plugins/pg-plugin", updateBody, "application/json")
				rec := httptest.NewRecorder()
				env.pluginHandler.HandleItem(rec, req)
				return rec
			},
			wantStatus: http.StatusOK,
			wantSpec: resources.OperationSpec{
				Type:         resources.OpUpdatePlugin,
				ResourceKind: resources.PluginKind,
				ResourceName: "pg-plugin",
				PluginName:   "pg-plugin",
			},
		},
		{
			name: "Delete",
			setup: func(env *pluginCapabilityEmissionEnv) {
				seedPluginServiceClass(t, env.scRegistry, "postgres")
				seedPluginDirect(t, env.pluginRegistry, "pg-plugin")
			},
			action: func(env *pluginCapabilityEmissionEnv) *httptest.ResponseRecorder {
				req := jsonRequest(http.MethodDelete, "/v1/plugins/pg-plugin", nil, "")
				rec := httptest.NewRecorder()
				env.pluginHandler.HandleItem(rec, req)
				return rec
			},
			wantStatus: http.StatusNoContent,
			wantSpec: resources.OperationSpec{
				Type:         resources.OpDeletePlugin,
				ResourceKind: resources.PluginKind,
				ResourceName: "pg-plugin",
				PluginName:   "pg-plugin",
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

func TestPluginEmission_NoEmissionOnFailedAction(t *testing.T) {
	t.Run("validation failure", func(t *testing.T) {
		env := newPluginCapabilityEmissionEnv(t)
		seedPluginServiceClass(t, env.scRegistry, "postgres")
		body := map[string]any{
			"metadata": map[string]any{"name": "INVALID"},
			"spec": map[string]any{
				"pluginType":       resources.PluginTypeDStoreOps,
				"version":          "1.0.0",
				"serviceClassRefs": []string{"postgres"},
				"deploymentMode":   resources.DeploymentModeCompiledIn,
			},
		}
		req := jsonRequest(http.MethodPost, "/v1/plugins", body, "application/json")
		rec := httptest.NewRecorder()
		env.pluginHandler.HandleCollection(rec, req)
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
		rec := createPlugin(env.pluginHandler, "pg-plugin")
		if rec.Code != http.StatusConflict {
			t.Fatalf("status = %d, want 409", rec.Code)
		}
		if ops := pluginCapabilityRecordedOperations(t, env); len(ops) != 0 {
			t.Errorf("operations = %d, want 0", len(ops))
		}
	})

	t.Run("missing service class ref", func(t *testing.T) {
		env := newPluginCapabilityEmissionEnv(t)
		rec := createPlugin(env.pluginHandler, "pg-plugin")
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400", rec.Code)
		}
		if ops := pluginCapabilityRecordedOperations(t, env); len(ops) != 0 {
			t.Errorf("operations = %d, want 0", len(ops))
		}
	})

	t.Run("update missing target", func(t *testing.T) {
		env := newPluginCapabilityEmissionEnv(t)
		seedPluginServiceClass(t, env.scRegistry, "postgres")
		req := jsonRequest(http.MethodPut, "/v1/plugins/missing", validPluginBody("missing"), "application/json")
		rec := httptest.NewRecorder()
		env.pluginHandler.HandleItem(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want 404", rec.Code)
		}
		if ops := pluginCapabilityRecordedOperations(t, env); len(ops) != 0 {
			t.Errorf("operations = %d, want 0", len(ops))
		}
	})

	t.Run("delete missing target", func(t *testing.T) {
		env := newPluginCapabilityEmissionEnv(t)
		req := jsonRequest(http.MethodDelete, "/v1/plugins/missing", nil, "")
		rec := httptest.NewRecorder()
		env.pluginHandler.HandleItem(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want 404", rec.Code)
		}
		if ops := pluginCapabilityRecordedOperations(t, env); len(ops) != 0 {
			t.Errorf("operations = %d, want 0", len(ops))
		}
	})

	t.Run("delete blocked", func(t *testing.T) {
		env := newPluginCapabilityEmissionEnv(t)
		seedPluginServiceClass(t, env.scRegistry, "postgres")
		seedPluginDirect(t, env.pluginRegistry, "pg-plugin")
		seedCapabilityDirect(t, env.capabilityRegistry, "pg-provision", "pg-plugin", "postgres")
		req := jsonRequest(http.MethodDelete, "/v1/plugins/pg-plugin", nil, "")
		rec := httptest.NewRecorder()
		env.pluginHandler.HandleItem(rec, req)
		if rec.Code != http.StatusConflict {
			t.Fatalf("status = %d, want 409; body=%s", rec.Code, rec.Body.String())
		}
		if ops := pluginCapabilityRecordedOperations(t, env); len(ops) != 0 {
			t.Errorf("operations = %d, want 0", len(ops))
		}
	})
}

func TestPluginEmission_FailureIsolation(t *testing.T) {
	newHandler := func() (*PluginHandler, *registry.ServiceClassRegistry, *registry.PluginRegistry) {
		scReg := registry.NewServiceClassRegistry()
		pluginReg := registry.NewPluginRegistry()
		capReg := registry.NewCapabilityRegistry()
		blocker := registry.NewCapabilityChildBlockerChecker(capReg)
		return NewPluginHandler(pluginReg, scReg, blocker, &failingEmitter{}), scReg, pluginReg
	}

	t.Run("create still returns 201", func(t *testing.T) {
		h, scReg, _ := newHandler()
		seedPluginServiceClass(t, scReg, "postgres")
		rec := createPlugin(h, "pg-plugin")
		if rec.Code != http.StatusCreated {
			t.Fatalf("status = %d, want 201; body=%s", rec.Code, rec.Body.String())
		}
	})

	t.Run("update still returns 200", func(t *testing.T) {
		h, scReg, pluginReg := newHandler()
		seedPluginServiceClass(t, scReg, "postgres")
		seedPluginDirect(t, pluginReg, "pg-plugin")
		req := jsonRequest(http.MethodPut, "/v1/plugins/pg-plugin", validPluginBody("pg-plugin"), "application/json")
		rec := httptest.NewRecorder()
		h.HandleItem(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
		}
	})

	t.Run("delete still returns 204", func(t *testing.T) {
		h, scReg, pluginReg := newHandler()
		seedPluginServiceClass(t, scReg, "postgres")
		seedPluginDirect(t, pluginReg, "pg-plugin")
		req := jsonRequest(http.MethodDelete, "/v1/plugins/pg-plugin", nil, "")
		rec := httptest.NewRecorder()
		h.HandleItem(rec, req)
		if rec.Code != http.StatusNoContent {
			t.Fatalf("status = %d, want 204; body=%s", rec.Code, rec.Body.String())
		}
	})
}
