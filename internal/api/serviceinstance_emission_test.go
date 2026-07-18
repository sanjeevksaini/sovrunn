package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/registry"
	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// siEmissionEnv wires a ServiceInstanceHandler to a real OperationRegistry-backed
// emitter for emission assertions.
type siEmissionEnv struct {
	siReg      *registry.ServiceInstanceRegistry
	orgReg     *registry.OrganizationRegistry
	ouReg      *registry.OrganizationUnitRegistry
	tenantReg  *registry.TenantRegistry
	projectReg *registry.ProjectRegistry
	scReg      *registry.ServiceClassRegistry
	spReg      *registry.ServicePlanRegistry
	bindingReg *registry.ServiceBindingRegistry
	opReg      *registry.OperationRegistry
	handler    *ServiceInstanceHandler
}

func newServiceInstanceEmissionEnv(t *testing.T) *siEmissionEnv {
	t.Helper()
	siReg := registry.NewServiceInstanceRegistry()
	orgReg := registry.NewOrganizationRegistry()
	ouReg := registry.NewOrganizationUnitRegistry()
	tenantReg := registry.NewTenantRegistry()
	projectReg := registry.NewProjectRegistry()
	scReg := registry.NewServiceClassRegistry()
	spReg := registry.NewServicePlanRegistry()
	capLookup := registry.NewCapabilityLookup(registry.NewCapabilityRegistry())
	bindingReg := registry.NewServiceBindingRegistry()
	opReg := registry.NewOperationRegistry()
	emitter := NewRegistryEmitter(opReg, nil)

	h := NewServiceInstanceHandler(
		siReg,
		orgReg,
		ouReg,
		tenantReg,
		projectReg,
		scReg,
		spReg,
		capLookup,
		bindingReg,
		emitter,
		nil,
	)
	return &siEmissionEnv{
		siReg:      siReg,
		orgReg:     orgReg,
		ouReg:      ouReg,
		tenantReg:  tenantReg,
		projectReg: projectReg,
		scReg:      scReg,
		spReg:      spReg,
		bindingReg: bindingReg,
		opReg:      opReg,
		handler:    h,
	}
}

func (e *siEmissionEnv) seedDefaults(t *testing.T) {
	t.Helper()
	seedOrg(t, e.orgReg, "nic")
	seedOU(t, e.ouReg, "nic", "ministry-health")
	seedTenant(t, e.tenantReg, "nic", "ministry-health", "payments")
	seedProjectDirect(t, e.projectReg, "nic", "ministry-health", "payments", "prod")
	seedServiceClassDirect(t, e.scReg, "postgres")
	seedServicePlanDirect(t, e.spReg, "postgres", "small")
}

func seedServiceInstanceDirect(t *testing.T, reg *registry.ServiceInstanceRegistry, name string) {
	t.Helper()
	si := resources.ServiceInstance{
		APIVersion: resources.ServiceInstanceAPIVersion,
		Kind:       resources.ServiceInstanceKind,
		Metadata:   resources.Metadata{Name: name},
		Spec: resources.ServiceInstanceSpec{
			OrganizationRef:     "nic",
			OrganizationUnitRef: "ministry-health",
			TenantRef:           "payments",
			ProjectRef:          "prod",
			ServiceClassRef:     "postgres",
			ServicePlanRef:      "small",
		},
		Status: resources.ServiceInstanceStatus{
			Phase:   "Ready",
			Message: "Registered only; no real provisioning in Phase 1",
		},
	}
	if _, err := reg.CreateServiceInstance(context.Background(), si); err != nil {
		t.Fatalf("seedServiceInstanceDirect(%s): %v", name, err)
	}
}

func seedServiceBindingDirect(
	t *testing.T,
	reg *registry.ServiceBindingRegistry,
	name, serviceInstanceRef string,
) {
	t.Helper()
	sb := resources.ServiceBinding{
		APIVersion: resources.ServiceBindingAPIVersion,
		Kind:       resources.ServiceBindingKind,
		Metadata:   resources.Metadata{Name: name},
		Spec: resources.ServiceBindingSpec{
			ServiceInstanceRef: serviceInstanceRef,
			ConsumerRef:        &resources.ConsumerRef{Kind: "Application", Name: "app-1"},
			BindingType:        resources.BindingTypeCredentials,
		},
		Status: resources.ServiceBindingStatus{Phase: "Ready"},
	}
	if _, err := reg.CreateServiceBinding(context.Background(), sb); err != nil {
		t.Fatalf("seedServiceBindingDirect(%s): %v", name, err)
	}
}

func siRecordedOperations(t *testing.T, env *siEmissionEnv) []resources.Operation {
	t.Helper()
	items, err := env.opReg.ListOperations(context.Background())
	if err != nil {
		t.Fatalf("ListOperations() error = %v", err)
	}
	return items
}

func siRequireSingleOperation(t *testing.T, env *siEmissionEnv) resources.Operation {
	t.Helper()
	items := siRecordedOperations(t, env)
	if len(items) != 1 {
		t.Fatalf("recorded operations = %d, want exactly 1", len(items))
	}
	return items[0]
}

func assertServiceInstanceOperation(t *testing.T, op resources.Operation, want resources.OperationSpec) {
	t.Helper()
	assertOperation(t, op, want)
	if op.Spec.ServiceInstanceName != want.ServiceInstanceName {
		t.Errorf("Spec.ServiceInstanceName = %q, want %q", op.Spec.ServiceInstanceName, want.ServiceInstanceName)
	}
	if op.Spec.ServiceBindingName != "" {
		t.Errorf("Spec.ServiceBindingName = %q, want empty", op.Spec.ServiceBindingName)
	}
}

func TestServiceInstanceEmission_MutatingPaths(t *testing.T) {
	updateBody := validServiceInstanceBody("pg-prod")
	updateBody["metadata"].(map[string]any)["displayName"] = "Postgres Prod"
	updateBody["spec"].(map[string]any)["parameters"] = map[string]any{"storage": "100Gi"}

	wantGovernance := resources.OperationSpec{
		ResourceKind:         resources.ServiceInstanceKind,
		ResourceName:         "pg-prod",
		OrganizationName:     "nic",
		OrganizationUnitName: "ministry-health",
		TenantName:           "payments",
		ProjectName:          "prod",
		ServiceInstanceName:  "pg-prod",
	}

	cases := []struct {
		name       string
		setup      func(env *siEmissionEnv)
		action     func(env *siEmissionEnv) *httptest.ResponseRecorder
		wantStatus int
		wantType   string
	}{
		{
			name: "Create",
			setup: func(env *siEmissionEnv) {
				env.seedDefaults(t)
			},
			action: func(env *siEmissionEnv) *httptest.ResponseRecorder {
				return createServiceInstance(env.handler, "pg-prod")
			},
			wantStatus: http.StatusCreated,
			wantType:   resources.OpCreateServiceInstance,
		},
		{
			name: "Update",
			setup: func(env *siEmissionEnv) {
				env.seedDefaults(t)
				seedServiceInstanceDirect(t, env.siReg, "pg-prod")
			},
			action: func(env *siEmissionEnv) *httptest.ResponseRecorder {
				req := jsonRequest(http.MethodPut, "/v1/service-instances/pg-prod", updateBody, "application/json")
				rec := httptest.NewRecorder()
				env.handler.HandleItem(rec, req)
				return rec
			},
			wantStatus: http.StatusOK,
			wantType:   resources.OpUpdateServiceInstance,
		},
		{
			name: "Delete",
			setup: func(env *siEmissionEnv) {
				env.seedDefaults(t)
				seedServiceInstanceDirect(t, env.siReg, "pg-prod")
			},
			action: func(env *siEmissionEnv) *httptest.ResponseRecorder {
				req := jsonRequest(http.MethodDelete, "/v1/service-instances/pg-prod", nil, "")
				rec := httptest.NewRecorder()
				env.handler.HandleItem(rec, req)
				return rec
			},
			wantStatus: http.StatusNoContent,
			wantType:   resources.OpDeleteServiceInstance,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			env := newServiceInstanceEmissionEnv(t)
			tc.setup(env)
			rec := tc.action(env)
			if rec.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d; body=%s", rec.Code, tc.wantStatus, rec.Body.String())
			}
			op := siRequireSingleOperation(t, env)
			want := wantGovernance
			want.Type = tc.wantType
			assertServiceInstanceOperation(t, op, want)
		})
	}
}

func TestServiceInstanceEmission_NoEmissionOnFailedAction(t *testing.T) {
	t.Run("validation failure", func(t *testing.T) {
		env := newServiceInstanceEmissionEnv(t)
		env.seedDefaults(t)
		body := validServiceInstanceBody("INVALID")
		req := jsonRequest(http.MethodPost, "/v1/service-instances", body, "application/json")
		rec := httptest.NewRecorder()
		env.handler.HandleCollection(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400", rec.Code)
		}
		if ops := siRecordedOperations(t, env); len(ops) != 0 {
			t.Errorf("operations = %d, want 0", len(ops))
		}
	})

	t.Run("duplicate create", func(t *testing.T) {
		env := newServiceInstanceEmissionEnv(t)
		env.seedDefaults(t)
		seedServiceInstanceDirect(t, env.siReg, "pg-prod")
		rec := createServiceInstance(env.handler, "pg-prod")
		if rec.Code != http.StatusConflict {
			t.Fatalf("status = %d, want 409", rec.Code)
		}
		if ops := siRecordedOperations(t, env); len(ops) != 0 {
			t.Errorf("operations = %d, want 0", len(ops))
		}
	})

	t.Run("missing reference", func(t *testing.T) {
		env := newServiceInstanceEmissionEnv(t)
		// No governance/catalog parents seeded.
		rec := createServiceInstance(env.handler, "pg-prod")
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400", rec.Code)
		}
		if ops := siRecordedOperations(t, env); len(ops) != 0 {
			t.Errorf("operations = %d, want 0", len(ops))
		}
	})

	t.Run("update missing target", func(t *testing.T) {
		env := newServiceInstanceEmissionEnv(t)
		env.seedDefaults(t)
		req := jsonRequest(
			http.MethodPut,
			"/v1/service-instances/missing",
			validServiceInstanceBody("missing"),
			"application/json",
		)
		rec := httptest.NewRecorder()
		env.handler.HandleItem(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want 404", rec.Code)
		}
		if ops := siRecordedOperations(t, env); len(ops) != 0 {
			t.Errorf("operations = %d, want 0", len(ops))
		}
	})

	t.Run("delete missing target", func(t *testing.T) {
		env := newServiceInstanceEmissionEnv(t)
		req := jsonRequest(http.MethodDelete, "/v1/service-instances/missing", nil, "")
		rec := httptest.NewRecorder()
		env.handler.HandleItem(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want 404", rec.Code)
		}
		if ops := siRecordedOperations(t, env); len(ops) != 0 {
			t.Errorf("operations = %d, want 0", len(ops))
		}
	})

	t.Run("delete blocked", func(t *testing.T) {
		env := newServiceInstanceEmissionEnv(t)
		env.seedDefaults(t)
		seedServiceInstanceDirect(t, env.siReg, "pg-prod")
		seedServiceBindingDirect(t, env.bindingReg, "pg-bind", "pg-prod")
		req := jsonRequest(http.MethodDelete, "/v1/service-instances/pg-prod", nil, "")
		rec := httptest.NewRecorder()
		env.handler.HandleItem(rec, req)
		if rec.Code != http.StatusConflict {
			t.Fatalf("status = %d, want 409; body=%s", rec.Code, rec.Body.String())
		}
		if ops := siRecordedOperations(t, env); len(ops) != 0 {
			t.Errorf("operations = %d, want 0", len(ops))
		}
	})
}

func TestServiceInstanceEmission_NilEmitterNoPanic(t *testing.T) {
	env := newTestServiceInstanceHandler() // emitter is nil
	env.seedDefaults(t)

	if rec := createServiceInstance(env.handler, "pg-prod"); rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want 201; body=%s", rec.Code, rec.Body.String())
	}

	body := validServiceInstanceBody("pg-prod")
	body["metadata"].(map[string]any)["displayName"] = "updated"
	req := jsonRequest(http.MethodPut, "/v1/service-instances/pg-prod", body, "application/json")
	rec := httptest.NewRecorder()
	env.handler.HandleItem(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("update status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}

	del := jsonRequest(http.MethodDelete, "/v1/service-instances/pg-prod", nil, "")
	delRec := httptest.NewRecorder()
	env.handler.HandleItem(delRec, del)
	if delRec.Code != http.StatusNoContent {
		t.Fatalf("delete status = %d, want 204; body=%s", delRec.Code, delRec.Body.String())
	}
}
