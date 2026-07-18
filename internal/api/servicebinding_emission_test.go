package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/registry"
	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// sbEmissionEnv wires a ServiceBindingHandler to a real OperationRegistry-backed
// emitter for emission assertions.
type sbEmissionEnv struct {
	bindingReg *registry.ServiceBindingRegistry
	siReg      *registry.ServiceInstanceRegistry
	opReg      *registry.OperationRegistry
	handler    *ServiceBindingHandler
}

func newServiceBindingEmissionEnv(t *testing.T) *sbEmissionEnv {
	t.Helper()
	bindingReg := registry.NewServiceBindingRegistry()
	siReg := registry.NewServiceInstanceRegistry()
	opReg := registry.NewOperationRegistry()
	emitter := NewRegistryEmitter(opReg, nil)

	h := NewServiceBindingHandler(bindingReg, siReg, emitter)
	return &sbEmissionEnv{
		bindingReg: bindingReg,
		siReg:      siReg,
		opReg:      opReg,
		handler:    h,
	}
}

func sbRecordedOperations(t *testing.T, env *sbEmissionEnv) []resources.Operation {
	t.Helper()
	items, err := env.opReg.ListOperations(context.Background())
	if err != nil {
		t.Fatalf("ListOperations() error = %v", err)
	}
	return items
}

func sbRequireSingleOperation(t *testing.T, env *sbEmissionEnv) resources.Operation {
	t.Helper()
	items := sbRecordedOperations(t, env)
	if len(items) != 1 {
		t.Fatalf("recorded operations = %d, want exactly 1", len(items))
	}
	return items[0]
}

func assertServiceBindingOperation(t *testing.T, op resources.Operation, want resources.OperationSpec) {
	t.Helper()
	assertOperation(t, op, want)
	if op.Spec.ServiceInstanceName != want.ServiceInstanceName {
		t.Errorf("Spec.ServiceInstanceName = %q, want %q", op.Spec.ServiceInstanceName, want.ServiceInstanceName)
	}
	if op.Spec.ServiceBindingName != want.ServiceBindingName {
		t.Errorf("Spec.ServiceBindingName = %q, want %q", op.Spec.ServiceBindingName, want.ServiceBindingName)
	}
}

func TestServiceBindingEmission_MutatingPaths(t *testing.T) {
	wantFields := resources.OperationSpec{
		ResourceKind:        resources.ServiceBindingKind,
		ResourceName:        "pg-binding",
		ServiceInstanceName: "pg-prod",
		ServiceBindingName:  "pg-binding",
	}

	cases := []struct {
		name       string
		setup      func(env *sbEmissionEnv)
		action     func(env *sbEmissionEnv) *httptest.ResponseRecorder
		wantStatus int
		wantType   string
	}{
		{
			name: "Create",
			setup: func(env *sbEmissionEnv) {
				seedServiceInstanceDirect(t, env.siReg, "pg-prod")
			},
			action: func(env *sbEmissionEnv) *httptest.ResponseRecorder {
				return createServiceBinding(env.handler, "pg-binding", "pg-prod")
			},
			wantStatus: http.StatusCreated,
			wantType:   resources.OpCreateServiceBinding,
		},
		{
			name: "Delete",
			setup: func(env *sbEmissionEnv) {
				seedServiceInstanceDirect(t, env.siReg, "pg-prod")
				seedServiceBindingDirect(t, env.bindingReg, "pg-binding", "pg-prod")
			},
			action: func(env *sbEmissionEnv) *httptest.ResponseRecorder {
				req := jsonRequest(http.MethodDelete, "/v1/service-bindings/pg-binding", nil, "")
				rec := httptest.NewRecorder()
				env.handler.HandleItem(rec, req)
				return rec
			},
			wantStatus: http.StatusNoContent,
			wantType:   resources.OpDeleteServiceBinding,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			env := newServiceBindingEmissionEnv(t)
			tc.setup(env)
			rec := tc.action(env)
			if rec.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d; body=%s", rec.Code, tc.wantStatus, rec.Body.String())
			}
			op := sbRequireSingleOperation(t, env)
			want := wantFields
			want.Type = tc.wantType
			assertServiceBindingOperation(t, op, want)
		})
	}
}

func TestServiceBindingEmission_NoEmissionOnFailedAction(t *testing.T) {
	t.Run("validation failure", func(t *testing.T) {
		env := newServiceBindingEmissionEnv(t)
		seedServiceInstanceDirect(t, env.siReg, "pg-prod")
		body := validServiceBindingBody("INVALID", "pg-prod")
		req := jsonRequest(http.MethodPost, "/v1/service-bindings", body, "application/json")
		rec := httptest.NewRecorder()
		env.handler.HandleCollection(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400", rec.Code)
		}
		if ops := sbRecordedOperations(t, env); len(ops) != 0 {
			t.Errorf("operations = %d, want 0", len(ops))
		}
	})

	t.Run("duplicate create", func(t *testing.T) {
		env := newServiceBindingEmissionEnv(t)
		seedServiceInstanceDirect(t, env.siReg, "pg-prod")
		seedServiceBindingDirect(t, env.bindingReg, "pg-binding", "pg-prod")
		rec := createServiceBinding(env.handler, "pg-binding", "pg-prod")
		if rec.Code != http.StatusConflict {
			t.Fatalf("status = %d, want 409", rec.Code)
		}
		if ops := sbRecordedOperations(t, env); len(ops) != 0 {
			t.Errorf("operations = %d, want 0", len(ops))
		}
	})

	t.Run("missing ServiceInstance ref", func(t *testing.T) {
		env := newServiceBindingEmissionEnv(t)
		// No ServiceInstance seeded.
		rec := createServiceBinding(env.handler, "pg-binding", "pg-prod")
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400", rec.Code)
		}
		if ops := sbRecordedOperations(t, env); len(ops) != 0 {
			t.Errorf("operations = %d, want 0", len(ops))
		}
	})

	t.Run("delete missing target", func(t *testing.T) {
		env := newServiceBindingEmissionEnv(t)
		req := jsonRequest(http.MethodDelete, "/v1/service-bindings/missing", nil, "")
		rec := httptest.NewRecorder()
		env.handler.HandleItem(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want 404", rec.Code)
		}
		if ops := sbRecordedOperations(t, env); len(ops) != 0 {
			t.Errorf("operations = %d, want 0", len(ops))
		}
	})
}

func TestServiceBindingEmission_NilEmitterNoPanic(t *testing.T) {
	env := newTestServiceBindingHandler() // emitter is nil
	env.seedServiceInstance(t, "pg-prod")

	if rec := createServiceBinding(env.handler, "pg-binding", "pg-prod"); rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want 201; body=%s", rec.Code, rec.Body.String())
	}

	del := jsonRequest(http.MethodDelete, "/v1/service-bindings/pg-binding", nil, "")
	delRec := httptest.NewRecorder()
	env.handler.HandleItem(delRec, del)
	if delRec.Code != http.StatusNoContent {
		t.Fatalf("delete status = %d, want 204; body=%s", delRec.Code, delRec.Body.String())
	}
}
