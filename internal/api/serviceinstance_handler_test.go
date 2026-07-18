package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/registry"
	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// siTestEnv bundles a ServiceInstanceHandler with the registries used to seed
// governance and catalog parents for handler tests.
type siTestEnv struct {
	handler    *ServiceInstanceHandler
	siReg      *registry.ServiceInstanceRegistry
	orgReg     *registry.OrganizationRegistry
	ouReg      *registry.OrganizationUnitRegistry
	tenantReg  *registry.TenantRegistry
	projectReg *registry.ProjectRegistry
	scReg      *registry.ServiceClassRegistry
	spReg      *registry.ServicePlanRegistry
	bindingReg *registry.ServiceBindingRegistry
}

func newTestServiceInstanceHandler() *siTestEnv {
	siReg := registry.NewServiceInstanceRegistry()
	orgReg := registry.NewOrganizationRegistry()
	ouReg := registry.NewOrganizationUnitRegistry()
	tenantReg := registry.NewTenantRegistry()
	projectReg := registry.NewProjectRegistry()
	scReg := registry.NewServiceClassRegistry()
	spReg := registry.NewServicePlanRegistry()
	capLookup := registry.NewCapabilityLookup(registry.NewCapabilityRegistry())
	bindingReg := registry.NewServiceBindingRegistry()

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
		nil,
		nil,
	)
	return &siTestEnv{
		handler:    h,
		siReg:      siReg,
		orgReg:     orgReg,
		ouReg:      ouReg,
		tenantReg:  tenantReg,
		projectReg: projectReg,
		scReg:      scReg,
		spReg:      spReg,
		bindingReg: bindingReg,
	}
}

func (e *siTestEnv) seedOrg(t *testing.T, name string) {
	t.Helper()
	org := resources.Organization{
		APIVersion: resources.OrgAPIVersion,
		Kind:       resources.OrgKind,
		Metadata:   resources.Metadata{Name: name},
		Status:     resources.OrganizationStatus{Phase: resources.PhaseActive},
	}
	if err := e.orgReg.CreateOrganization(context.Background(), org); err != nil {
		t.Fatalf("seedOrg(%s): %v", name, err)
	}
}

func (e *siTestEnv) seedOU(t *testing.T, orgName, name string) {
	t.Helper()
	ou := resources.OrganizationUnit{
		APIVersion: resources.OUAPIVersion,
		Kind:       resources.OUKind,
		Metadata:   resources.Metadata{Name: name},
		Spec:       resources.OrganizationUnitSpec{OrganizationName: orgName},
		Status:     resources.OrganizationUnitStatus{Phase: resources.PhaseActive},
	}
	if _, err := e.ouReg.CreateOrganizationUnit(context.Background(), ou); err != nil {
		t.Fatalf("seedOU(%s/%s): %v", orgName, name, err)
	}
}

func (e *siTestEnv) seedTenant(t *testing.T, orgName, ouName, name string) {
	t.Helper()
	tenant := resources.Tenant{
		APIVersion: resources.TenantAPIVersion,
		Kind:       resources.TenantKind,
		Metadata:   resources.Metadata{Name: name},
		Spec: resources.TenantSpec{
			OrganizationName:     orgName,
			OrganizationUnitName: ouName,
		},
		Status: resources.TenantStatus{Phase: resources.PhaseActive},
	}
	if _, err := e.tenantReg.CreateTenant(context.Background(), tenant); err != nil {
		t.Fatalf("seedTenant(%s/%s/%s): %v", orgName, ouName, name, err)
	}
}

func (e *siTestEnv) seedProject(t *testing.T, orgName, ouName, tenantName, name string) {
	t.Helper()
	project := resources.Project{
		APIVersion: resources.ProjectAPIVersion,
		Kind:       resources.ProjectKind,
		Metadata:   resources.Metadata{Name: name},
		Spec: resources.ProjectSpec{
			OrganizationName:     orgName,
			OrganizationUnitName: ouName,
			TenantName:           tenantName,
		},
		Status: resources.ProjectStatus{Phase: resources.PhaseActive},
	}
	if _, err := e.projectReg.CreateProject(context.Background(), project); err != nil {
		t.Fatalf("seedProject(%s/%s/%s/%s): %v", orgName, ouName, tenantName, name, err)
	}
}

func (e *siTestEnv) seedServiceClass(t *testing.T, name string) {
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
	if _, err := e.scReg.CreateServiceClass(context.Background(), sc); err != nil {
		t.Fatalf("seedServiceClass(%s): %v", name, err)
	}
}

func (e *siTestEnv) seedServicePlan(t *testing.T, serviceClassName, name string) {
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
	if _, err := e.spReg.CreateServicePlan(context.Background(), sp); err != nil {
		t.Fatalf("seedServicePlan(%s/%s): %v", serviceClassName, name, err)
	}
}

// seedDefaults creates the common governance + catalog graph used by most tests:
// nic / ministry-health / payments / prod + postgres / small.
func (e *siTestEnv) seedDefaults(t *testing.T) {
	t.Helper()
	e.seedOrg(t, "nic")
	e.seedOU(t, "nic", "ministry-health")
	e.seedTenant(t, "nic", "ministry-health", "payments")
	e.seedProject(t, "nic", "ministry-health", "payments", "prod")
	e.seedServiceClass(t, "postgres")
	e.seedServicePlan(t, "postgres", "small")
}

func validServiceInstanceBody(name string) map[string]any {
	return map[string]any{
		"metadata": map[string]any{"name": name},
		"spec": map[string]any{
			"organizationRef":     "nic",
			"organizationUnitRef": "ministry-health",
			"tenantRef":           "payments",
			"projectRef":          "prod",
			"serviceClassRef":     "postgres",
			"servicePlanRef":      "small",
		},
	}
}

func createServiceInstance(h *ServiceInstanceHandler, name string) *httptest.ResponseRecorder {
	req := jsonRequest(http.MethodPost, "/v1/service-instances", validServiceInstanceBody(name), "application/json")
	rec := httptest.NewRecorder()
	h.HandleCollection(rec, req)
	return rec
}

func TestServiceInstanceHandler_Create_Valid(t *testing.T) {
	env := newTestServiceInstanceHandler()
	env.seedDefaults(t)

	rec := createServiceInstance(env.handler, "pg-prod")
	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body=%s", rec.Code, rec.Body.String())
	}
	var si resources.ServiceInstance
	if err := json.NewDecoder(rec.Body).Decode(&si); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if si.APIVersion != resources.ServiceInstanceAPIVersion || si.Kind != resources.ServiceInstanceKind {
		t.Errorf("apiVersion/kind = %q/%q, want server-owned", si.APIVersion, si.Kind)
	}
	if si.Metadata.Name != "pg-prod" {
		t.Errorf("name = %q, want pg-prod", si.Metadata.Name)
	}
	if si.Status.Phase != "Ready" {
		t.Errorf("phase = %q, want Ready", si.Status.Phase)
	}
	if si.Status.Message == "" {
		t.Errorf("status.message = empty, want non-empty")
	}
	if si.Spec.OrganizationRef != "nic" || si.Spec.TenantRef != "payments" ||
		si.Spec.ProjectRef != "prod" || si.Spec.ServiceClassRef != "postgres" ||
		si.Spec.ServicePlanRef != "small" {
		t.Errorf("unexpected spec: %+v", si.Spec)
	}
}

func TestServiceInstanceHandler_Create_Duplicate(t *testing.T) {
	env := newTestServiceInstanceHandler()
	env.seedDefaults(t)
	if rec := createServiceInstance(env.handler, "pg-prod"); rec.Code != http.StatusCreated {
		t.Fatalf("first create status = %d", rec.Code)
	}
	rec := createServiceInstance(env.handler, "pg-prod")
	if rec.Code != http.StatusConflict {
		t.Fatalf("status = %d, want 409", rec.Code)
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeResourceAlreadyExists {
		t.Errorf("code = %q, want RESOURCE_ALREADY_EXISTS", errBody.Code)
	}
}

func TestServiceInstanceHandler_Create_DuplicateAcrossGovernanceRefs(t *testing.T) {
	env := newTestServiceInstanceHandler()
	env.seedDefaults(t)
	env.seedProject(t, "nic", "ministry-health", "payments", "staging")

	if rec := createServiceInstance(env.handler, "pg-prod"); rec.Code != http.StatusCreated {
		t.Fatalf("first create status = %d", rec.Code)
	}

	body := validServiceInstanceBody("pg-prod")
	body["spec"].(map[string]any)["projectRef"] = "staging"
	req := jsonRequest(http.MethodPost, "/v1/service-instances", body, "application/json")
	rec := httptest.NewRecorder()
	env.handler.HandleCollection(rec, req)
	if rec.Code != http.StatusConflict {
		t.Fatalf("status = %d, want 409; body=%s", rec.Code, rec.Body.String())
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeResourceAlreadyExists {
		t.Errorf("code = %q, want RESOURCE_ALREADY_EXISTS", errBody.Code)
	}
}

func TestServiceInstanceHandler_Create_InvalidFields(t *testing.T) {
	env := newTestServiceInstanceHandler()
	env.seedDefaults(t)

	body := validServiceInstanceBody("INVALID")
	req := jsonRequest(http.MethodPost, "/v1/service-instances", body, "application/json")
	rec := httptest.NewRecorder()
	env.handler.HandleCollection(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeValidationFailed || errBody.Field != "metadata.name" {
		t.Errorf("error = %+v, want VALIDATION_FAILED metadata.name", errBody)
	}
}

func TestServiceInstanceHandler_Create_StatusFieldRejected(t *testing.T) {
	env := newTestServiceInstanceHandler()
	env.seedDefaults(t)

	payload := `{"metadata":{"name":"pg-prod"},"spec":{"organizationRef":"nic","organizationUnitRef":"ministry-health","tenantRef":"payments","projectRef":"prod","serviceClassRef":"postgres","servicePlanRef":"small"},"status":{}}`
	req := withRequestID(httptest.NewRequest(http.MethodPost, "/v1/service-instances", strings.NewReader(payload)), "id")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	env.handler.HandleCollection(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeValidationFailed || errBody.Field != "status" {
		t.Errorf("error = %+v, want VALIDATION_FAILED field=status", errBody)
	}
}

func TestServiceInstanceHandler_Create_BadJSON(t *testing.T) {
	env := newTestServiceInstanceHandler()
	req := withRequestID(httptest.NewRequest(http.MethodPost, "/v1/service-instances", strings.NewReader("{")), "id")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	env.handler.HandleCollection(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
}

func TestServiceInstanceHandler_Create_UnknownField(t *testing.T) {
	env := newTestServiceInstanceHandler()
	env.seedDefaults(t)

	payload := `{"metadata":{"name":"pg-prod"},"spec":{"organizationRef":"nic","organizationUnitRef":"ministry-health","tenantRef":"payments","projectRef":"prod","serviceClassRef":"postgres","servicePlanRef":"small"},"bogus":true}`
	req := withRequestID(httptest.NewRequest(http.MethodPost, "/v1/service-instances", strings.NewReader(payload)), "id")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	env.handler.HandleCollection(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeValidationFailed {
		t.Errorf("code = %q, want VALIDATION_FAILED", errBody.Code)
	}
}

func TestServiceInstanceHandler_Create_OversizedBody(t *testing.T) {
	env := newTestServiceInstanceHandler()
	env.seedDefaults(t)

	large := strings.Repeat("a", 1<<20+1)
	payload := fmt.Sprintf(
		`{"metadata":{"name":"pg-prod","displayName":"%s"},"spec":{"organizationRef":"nic","organizationUnitRef":"ministry-health","tenantRef":"payments","projectRef":"prod","serviceClassRef":"postgres","servicePlanRef":"small"}}`,
		large,
	)
	req := withRequestID(httptest.NewRequest(http.MethodPost, "/v1/service-instances", strings.NewReader(payload)), "id")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	env.handler.HandleCollection(rec, req)
	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("status = %d, want 413", rec.Code)
	}
}

func TestServiceInstanceHandler_Create_MissingOrganization(t *testing.T) {
	env := newTestServiceInstanceHandler()
	env.seedOU(t, "nic", "ministry-health")
	env.seedTenant(t, "nic", "ministry-health", "payments")
	env.seedProject(t, "nic", "ministry-health", "payments", "prod")
	env.seedServiceClass(t, "postgres")
	env.seedServicePlan(t, "postgres", "small")

	rec := createServiceInstance(env.handler, "pg-prod")
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body=%s", rec.Code, rec.Body.String())
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeValidationFailed || errBody.Field != "spec.organizationRef" {
		t.Errorf("error = %+v, want VALIDATION_FAILED spec.organizationRef", errBody)
	}
}

func TestServiceInstanceHandler_Create_MissingTenant(t *testing.T) {
	env := newTestServiceInstanceHandler()
	env.seedOrg(t, "nic")
	env.seedOU(t, "nic", "ministry-health")
	env.seedProject(t, "nic", "ministry-health", "payments", "prod")
	env.seedServiceClass(t, "postgres")
	env.seedServicePlan(t, "postgres", "small")

	rec := createServiceInstance(env.handler, "pg-prod")
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body=%s", rec.Code, rec.Body.String())
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeValidationFailed || errBody.Field != "spec.tenantRef" {
		t.Errorf("error = %+v, want VALIDATION_FAILED spec.tenantRef", errBody)
	}
}

func TestServiceInstanceHandler_Create_MissingProject(t *testing.T) {
	env := newTestServiceInstanceHandler()
	env.seedOrg(t, "nic")
	env.seedOU(t, "nic", "ministry-health")
	env.seedTenant(t, "nic", "ministry-health", "payments")
	env.seedServiceClass(t, "postgres")
	env.seedServicePlan(t, "postgres", "small")

	rec := createServiceInstance(env.handler, "pg-prod")
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body=%s", rec.Code, rec.Body.String())
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeValidationFailed || errBody.Field != "spec.projectRef" {
		t.Errorf("error = %+v, want VALIDATION_FAILED spec.projectRef", errBody)
	}
}

func TestServiceInstanceHandler_Create_MissingServiceClass(t *testing.T) {
	env := newTestServiceInstanceHandler()
	env.seedOrg(t, "nic")
	env.seedOU(t, "nic", "ministry-health")
	env.seedTenant(t, "nic", "ministry-health", "payments")
	env.seedProject(t, "nic", "ministry-health", "payments", "prod")
	env.seedServicePlan(t, "postgres", "small")

	rec := createServiceInstance(env.handler, "pg-prod")
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body=%s", rec.Code, rec.Body.String())
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeValidationFailed || errBody.Field != "spec.serviceClassRef" {
		t.Errorf("error = %+v, want VALIDATION_FAILED spec.serviceClassRef", errBody)
	}
}

func TestServiceInstanceHandler_Create_MissingServicePlan(t *testing.T) {
	env := newTestServiceInstanceHandler()
	env.seedOrg(t, "nic")
	env.seedOU(t, "nic", "ministry-health")
	env.seedTenant(t, "nic", "ministry-health", "payments")
	env.seedProject(t, "nic", "ministry-health", "payments", "prod")
	env.seedServiceClass(t, "postgres")

	rec := createServiceInstance(env.handler, "pg-prod")
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body=%s", rec.Code, rec.Body.String())
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeValidationFailed || errBody.Field != "spec.servicePlanRef" {
		t.Errorf("error = %+v, want VALIDATION_FAILED spec.servicePlanRef", errBody)
	}
}

func TestServiceInstanceHandler_Create_ServicePlanNotMatchingServiceClass(t *testing.T) {
	env := newTestServiceInstanceHandler()
	env.seedDefaults(t)
	env.seedServiceClass(t, "redis")
	env.seedServicePlan(t, "redis", "small")

	// Plan "small" exists under redis, but not under a mismatched lookup when
	// we request postgres with a plan name that only exists under redis after
	// removing the postgres/small plan — use a plan name unique to redis.
	env.seedServicePlan(t, "redis", "cache-small")

	body := validServiceInstanceBody("pg-prod")
	body["spec"].(map[string]any)["serviceClassRef"] = "postgres"
	body["spec"].(map[string]any)["servicePlanRef"] = "cache-small"
	req := jsonRequest(http.MethodPost, "/v1/service-instances", body, "application/json")
	rec := httptest.NewRecorder()
	env.handler.HandleCollection(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body=%s", rec.Code, rec.Body.String())
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeValidationFailed || errBody.Field != "spec.servicePlanRef" {
		t.Errorf("error = %+v, want VALIDATION_FAILED spec.servicePlanRef", errBody)
	}
}

func TestServiceInstanceHandler_Create_GovernanceHierarchyInconsistency(t *testing.T) {
	env := newTestServiceInstanceHandler()
	env.seedDefaults(t)
	env.seedOU(t, "nic", "ministry-finance")

	// Tenant "payments" exists under ministry-health, not ministry-finance.
	body := validServiceInstanceBody("pg-prod")
	body["spec"].(map[string]any)["organizationUnitRef"] = "ministry-finance"
	req := jsonRequest(http.MethodPost, "/v1/service-instances", body, "application/json")
	rec := httptest.NewRecorder()
	env.handler.HandleCollection(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body=%s", rec.Code, rec.Body.String())
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeValidationFailed || errBody.Field != "spec.tenantRef" {
		t.Errorf("error = %+v, want VALIDATION_FAILED spec.tenantRef", errBody)
	}
}

func TestServiceInstanceHandler_Get_Found(t *testing.T) {
	env := newTestServiceInstanceHandler()
	env.seedDefaults(t)
	if rec := createServiceInstance(env.handler, "pg-prod"); rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d", rec.Code)
	}

	req := jsonRequest(http.MethodGet, "/v1/service-instances/pg-prod", nil, "")
	rec := httptest.NewRecorder()
	env.handler.HandleItem(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	var si resources.ServiceInstance
	if err := json.NewDecoder(rec.Body).Decode(&si); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if si.Metadata.Name != "pg-prod" {
		t.Errorf("name = %q, want pg-prod", si.Metadata.Name)
	}
}

func TestServiceInstanceHandler_Get_NotFound(t *testing.T) {
	env := newTestServiceInstanceHandler()
	req := jsonRequest(http.MethodGet, "/v1/service-instances/missing", nil, "")
	rec := httptest.NewRecorder()
	env.handler.HandleItem(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeResourceNotFound {
		t.Errorf("code = %q, want RESOURCE_NOT_FOUND", errBody.Code)
	}
}

func TestServiceInstanceHandler_Get_InvalidPathSegment(t *testing.T) {
	env := newTestServiceInstanceHandler()
	req := jsonRequest(http.MethodGet, "/v1/service-instances/INVALID", nil, "")
	rec := httptest.NewRecorder()
	env.handler.HandleItem(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeValidationFailed {
		t.Errorf("code = %q, want VALIDATION_FAILED", errBody.Code)
	}
}

func TestServiceInstanceHandler_List_Empty(t *testing.T) {
	env := newTestServiceInstanceHandler()
	req := jsonRequest(http.MethodGet, "/v1/service-instances", nil, "")
	rec := httptest.NewRecorder()
	env.handler.HandleCollection(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var resp serviceInstanceListResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Items == nil {
		t.Fatal("items = nil, want non-nil empty slice")
	}
	if len(resp.Items) != 0 {
		t.Errorf("len(items) = %d, want 0", len(resp.Items))
	}
}

func TestServiceInstanceHandler_List_Sorted(t *testing.T) {
	env := newTestServiceInstanceHandler()
	env.seedDefaults(t)
	for _, name := range []string{"zeta", "alpha", "middle"} {
		if rec := createServiceInstance(env.handler, name); rec.Code != http.StatusCreated {
			t.Fatalf("create %s status = %d", name, rec.Code)
		}
	}

	req := jsonRequest(http.MethodGet, "/v1/service-instances", nil, "")
	rec := httptest.NewRecorder()
	env.handler.HandleCollection(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var resp serviceInstanceListResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(resp.Items) != 3 {
		t.Fatalf("len(items) = %d, want 3", len(resp.Items))
	}
	want := []string{"alpha", "middle", "zeta"}
	for i, name := range want {
		if resp.Items[i].Metadata.Name != name {
			t.Errorf("items[%d].name = %q, want %q", i, resp.Items[i].Metadata.Name, name)
		}
	}
}

func TestServiceInstanceHandler_List_QueryFilters(t *testing.T) {
	env := newTestServiceInstanceHandler()
	env.seedDefaults(t)
	env.seedTenant(t, "nic", "ministry-health", "billing")
	env.seedProject(t, "nic", "ministry-health", "billing", "prod")
	env.seedProject(t, "nic", "ministry-health", "payments", "staging")

	createNamed := func(name, tenant, project string) {
		t.Helper()
		body := validServiceInstanceBody(name)
		body["spec"].(map[string]any)["tenantRef"] = tenant
		body["spec"].(map[string]any)["projectRef"] = project
		req := jsonRequest(http.MethodPost, "/v1/service-instances", body, "application/json")
		rec := httptest.NewRecorder()
		env.handler.HandleCollection(rec, req)
		if rec.Code != http.StatusCreated {
			t.Fatalf("create %s status = %d; body=%s", name, rec.Code, rec.Body.String())
		}
	}
	createNamed("si-payments-prod", "payments", "prod")
	createNamed("si-payments-staging", "payments", "staging")
	createNamed("si-billing-prod", "billing", "prod")

	t.Run("tenantRef", func(t *testing.T) {
		req := jsonRequest(http.MethodGet, "/v1/service-instances?tenantRef=payments", nil, "")
		rec := httptest.NewRecorder()
		env.handler.HandleCollection(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", rec.Code)
		}
		var resp serviceInstanceListResponse
		if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if len(resp.Items) != 2 {
			t.Fatalf("len(items) = %d, want 2", len(resp.Items))
		}
		for _, item := range resp.Items {
			if item.Spec.TenantRef != "payments" {
				t.Errorf("tenantRef = %q, want payments", item.Spec.TenantRef)
			}
		}
	})

	t.Run("projectRef", func(t *testing.T) {
		req := jsonRequest(http.MethodGet, "/v1/service-instances?projectRef=prod", nil, "")
		rec := httptest.NewRecorder()
		env.handler.HandleCollection(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", rec.Code)
		}
		var resp serviceInstanceListResponse
		if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if len(resp.Items) != 2 {
			t.Fatalf("len(items) = %d, want 2", len(resp.Items))
		}
	})

	t.Run("tenantRefAndProjectRef", func(t *testing.T) {
		req := jsonRequest(http.MethodGet, "/v1/service-instances?tenantRef=payments&projectRef=prod", nil, "")
		rec := httptest.NewRecorder()
		env.handler.HandleCollection(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", rec.Code)
		}
		var resp serviceInstanceListResponse
		if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if len(resp.Items) != 1 || resp.Items[0].Metadata.Name != "si-payments-prod" {
			t.Fatalf("items = %+v, want [si-payments-prod]", resp.Items)
		}
	})
}

func TestServiceInstanceHandler_Update_Success(t *testing.T) {
	env := newTestServiceInstanceHandler()
	env.seedDefaults(t)
	if rec := createServiceInstance(env.handler, "pg-prod"); rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d", rec.Code)
	}

	body := validServiceInstanceBody("pg-prod")
	body["metadata"].(map[string]any)["displayName"] = "Postgres Prod"
	body["metadata"].(map[string]any)["labels"] = map[string]any{"env": "prod"}
	body["spec"].(map[string]any)["parameters"] = map[string]any{"storage": "100Gi"}
	req := jsonRequest(http.MethodPut, "/v1/service-instances/pg-prod", body, "application/json")
	rec := httptest.NewRecorder()
	env.handler.HandleItem(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	var si resources.ServiceInstance
	if err := json.NewDecoder(rec.Body).Decode(&si); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if si.Metadata.DisplayName != "Postgres Prod" {
		t.Errorf("displayName = %q, want Postgres Prod", si.Metadata.DisplayName)
	}
	if si.Metadata.Labels["env"] != "prod" {
		t.Errorf("labels = %+v, want env=prod", si.Metadata.Labels)
	}
	if si.Spec.Parameters["storage"] != "100Gi" {
		t.Errorf("parameters = %+v, want storage=100Gi", si.Spec.Parameters)
	}
}

func TestServiceInstanceHandler_Update_PreservesStoredStatus(t *testing.T) {
	env := newTestServiceInstanceHandler()
	env.seedDefaults(t)
	if rec := createServiceInstance(env.handler, "pg-prod"); rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d", rec.Code)
	}

	body := validServiceInstanceBody("pg-prod")
	body["metadata"].(map[string]any)["displayName"] = "updated"
	req := jsonRequest(http.MethodPut, "/v1/service-instances/pg-prod", body, "application/json")
	rec := httptest.NewRecorder()
	env.handler.HandleItem(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	var si resources.ServiceInstance
	if err := json.NewDecoder(rec.Body).Decode(&si); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if si.Status.Phase != "Ready" {
		t.Errorf("phase = %q, want Ready", si.Status.Phase)
	}
	if si.Status.Message != "Registered only; no real provisioning in Phase 1" {
		t.Errorf("message = %q, want Phase 1 message", si.Status.Message)
	}
	if si.Metadata.DisplayName != "updated" {
		t.Errorf("displayName = %q, want updated", si.Metadata.DisplayName)
	}
}

func TestServiceInstanceHandler_Update_NotFound(t *testing.T) {
	env := newTestServiceInstanceHandler()
	env.seedDefaults(t)

	req := jsonRequest(http.MethodPut, "/v1/service-instances/missing", validServiceInstanceBody("missing"), "application/json")
	rec := httptest.NewRecorder()
	env.handler.HandleItem(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeResourceNotFound {
		t.Errorf("code = %q, want RESOURCE_NOT_FOUND", errBody.Code)
	}
}

func TestServiceInstanceHandler_Update_NameMismatch(t *testing.T) {
	env := newTestServiceInstanceHandler()
	env.seedDefaults(t)
	if rec := createServiceInstance(env.handler, "pg-prod"); rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d", rec.Code)
	}

	req := jsonRequest(http.MethodPut, "/v1/service-instances/pg-prod", validServiceInstanceBody("other"), "application/json")
	rec := httptest.NewRecorder()
	env.handler.HandleItem(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Field != "metadata.name" {
		t.Errorf("field = %q, want metadata.name", errBody.Field)
	}
}

func TestServiceInstanceHandler_Update_StatusFieldRejected(t *testing.T) {
	env := newTestServiceInstanceHandler()
	env.seedDefaults(t)
	if rec := createServiceInstance(env.handler, "pg-prod"); rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d", rec.Code)
	}

	payload := `{"metadata":{"name":"pg-prod"},"spec":{"organizationRef":"nic","organizationUnitRef":"ministry-health","tenantRef":"payments","projectRef":"prod","serviceClassRef":"postgres","servicePlanRef":"small"},"status":{"phase":"Ready"}}`
	req := withRequestID(httptest.NewRequest(http.MethodPut, "/v1/service-instances/pg-prod", strings.NewReader(payload)), "id")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	env.handler.HandleItem(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeValidationFailed || errBody.Field != "status" {
		t.Errorf("error = %+v, want VALIDATION_FAILED field=status", errBody)
	}
}

func TestServiceInstanceHandler_Update_ImmutableFields(t *testing.T) {
	cases := []struct {
		field string
		apply func(spec map[string]any)
	}{
		{"spec.organizationRef", func(spec map[string]any) { spec["organizationRef"] = "other-org" }},
		{"spec.organizationUnitRef", func(spec map[string]any) { spec["organizationUnitRef"] = "other-ou" }},
		{"spec.tenantRef", func(spec map[string]any) { spec["tenantRef"] = "other-tenant" }},
		{"spec.projectRef", func(spec map[string]any) { spec["projectRef"] = "other-project" }},
		{"spec.serviceClassRef", func(spec map[string]any) { spec["serviceClassRef"] = "other-class" }},
		{"spec.servicePlanRef", func(spec map[string]any) { spec["servicePlanRef"] = "other-plan" }},
	}

	for _, tc := range cases {
		t.Run(tc.field, func(t *testing.T) {
			env := newTestServiceInstanceHandler()
			env.seedDefaults(t)
			if rec := createServiceInstance(env.handler, "pg-prod"); rec.Code != http.StatusCreated {
				t.Fatalf("create status = %d", rec.Code)
			}

			body := validServiceInstanceBody("pg-prod")
			tc.apply(body["spec"].(map[string]any))
			req := jsonRequest(http.MethodPut, "/v1/service-instances/pg-prod", body, "application/json")
			rec := httptest.NewRecorder()
			env.handler.HandleItem(rec, req)
			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want 400; body=%s", rec.Code, rec.Body.String())
			}
			errBody := decodeAPIError(t, rec)
			if errBody.Code != resources.ErrCodeValidationFailed || errBody.Field != tc.field {
				t.Errorf("error = %+v, want VALIDATION_FAILED %s", errBody, tc.field)
			}
		})
	}
}

func TestServiceInstanceHandler_Delete_Success(t *testing.T) {
	env := newTestServiceInstanceHandler()
	env.seedDefaults(t)
	if rec := createServiceInstance(env.handler, "pg-prod"); rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d", rec.Code)
	}

	req := jsonRequest(http.MethodDelete, "/v1/service-instances/pg-prod", nil, "")
	rec := httptest.NewRecorder()
	env.handler.HandleItem(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204", rec.Code)
	}
	if rec.Body.Len() != 0 {
		t.Errorf("body = %q, want empty", rec.Body.String())
	}
}

func TestServiceInstanceHandler_Delete_NotFound(t *testing.T) {
	env := newTestServiceInstanceHandler()
	req := jsonRequest(http.MethodDelete, "/v1/service-instances/missing", nil, "")
	rec := httptest.NewRecorder()
	env.handler.HandleItem(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeResourceNotFound {
		t.Errorf("code = %q, want RESOURCE_NOT_FOUND", errBody.Code)
	}
}

func TestServiceInstanceHandler_Delete_BlockedByServiceBindings(t *testing.T) {
	env := newTestServiceInstanceHandler()
	env.seedDefaults(t)
	if rec := createServiceInstance(env.handler, "pg-prod"); rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d", rec.Code)
	}

	sb := resources.ServiceBinding{
		APIVersion: resources.ServiceBindingAPIVersion,
		Kind:       resources.ServiceBindingKind,
		Metadata:   resources.Metadata{Name: "pg-bind"},
		Spec: resources.ServiceBindingSpec{
			ServiceInstanceRef: "pg-prod",
			ConsumerRef:        &resources.ConsumerRef{Kind: "Application", Name: "app-1"},
			BindingType:        resources.BindingTypeCredentials,
		},
		Status: resources.ServiceBindingStatus{Phase: "Ready"},
	}
	if _, err := env.bindingReg.CreateServiceBinding(context.Background(), sb); err != nil {
		t.Fatalf("seed ServiceBinding: %v", err)
	}

	req := jsonRequest(http.MethodDelete, "/v1/service-instances/pg-prod", nil, "")
	rec := httptest.NewRecorder()
	env.handler.HandleItem(rec, req)
	if rec.Code != http.StatusConflict {
		t.Fatalf("status = %d, want 409; body=%s", rec.Code, rec.Body.String())
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeDeleteBlocked {
		t.Errorf("code = %q, want DELETE_BLOCKED", errBody.Code)
	}
	if !strings.Contains(errBody.Message, "ServiceBinding") {
		t.Errorf("message = %q, want it to mention ServiceBinding", errBody.Message)
	}
}

func TestServiceInstanceHandler_WrongPathShape(t *testing.T) {
	env := newTestServiceInstanceHandler()
	req := jsonRequest(http.MethodGet, "/v1/service-instances/pg-prod/extra", nil, "")
	rec := httptest.NewRecorder()
	env.handler.HandleItem(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
	errBody := decodeAPIError(t, rec)
	if errBody.Code != resources.ErrCodeResourceNotFound {
		t.Errorf("code = %q, want RESOURCE_NOT_FOUND", errBody.Code)
	}
}
