// Package conformance — FEATURE-0015 local conformance (Task 17).
//
// One runtime scenario per FEATURE-0015-owned registry ID:
// VS0-CF-X03 and VS0-CF-F15-01 through VS0-CF-F15-41.
//
// Anti-drift/non-runtime gates AC-F15-05 and AC-F15-07 remain verified by
// named repository checks (make phase2r-drift-check / make vs000-contract-check),
// not invented runtime conformance (ADH-2026-046 decision 3).
package conformance

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/apiconform"
	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/model"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/validate"
	"github.com/sanjeevksaini/sovrunn/internal/server"
)

const (
	f15Principal         = "bootstrap-principal"
	f15PlatformUID       = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	f15ProviderUID       = "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	f15OtherProviderUID  = "cccccccccccccccccccccccccccccccc"
	f15MediaMergePatch   = "application/merge-patch+json"
	f15ViolationSafeDeny = apiproblem.ViolationCode("VS0_AUTHORIZATION_SAFE_DENIAL")
	f15ViolationRootReq  = apiproblem.ViolationCode("VS0_CLOUDPLATFORM_ROOT_REQUIRED")
)

// ---------------------------------------------------------------------------
// Grant resolver and harness
// ---------------------------------------------------------------------------

type f15Grants struct {
	principal string
	grants    []cloudmodel.Grant
}

func (g *f15Grants) PrincipalID() string { return g.principal }

func (g *f15Grants) CoarseLookup(action string) []cloudmodel.Grant {
	var out []cloudmodel.Grant
	for _, gr := range g.grants {
		if gr.Action == action {
			out = append(out, gr)
		}
	}
	return out
}

func (g *f15Grants) AuthorizeExact(action string, scope apimeta.ScopeIdentity, resourceUID, relatedPlatformUID string) bool {
	for _, gr := range g.grants {
		if gr.Action != action || gr.Scope != scope {
			continue
		}
		if gr.ResourceUID != "" && gr.ResourceUID != resourceUID {
			continue
		}
		if action == cloudmodel.ActionParticipationRequestProvider {
			if relatedPlatformUID != "" && gr.RelatedCloudPlatformUID != relatedPlatformUID {
				continue
			}
		}
		return true
	}
	return false
}

func providerScope(uid string) apimeta.ScopeIdentity {
	return apimeta.ScopeIdentity{Kind: apimeta.ScopeCloudProvider, UID: uid}
}

func platformScope(uid string) apimeta.ScopeIdentity {
	return apimeta.ScopeIdentity{Kind: apimeta.ScopeCloudPlatform, UID: uid}
}

func fullAdminGrants(providerUID, platformUID string) *f15Grants {
	g := &f15Grants{principal: f15Principal}
	g.grants = []cloudmodel.Grant{
		{Action: cloudmodel.ActionCloudPlatformWrite, Scope: cloudmodel.PlatformRootScope},
		{Action: cloudmodel.ActionCloudPlatformRead, Scope: cloudmodel.PlatformRootScope},
		{Action: cloudmodel.ActionCloudProviderWrite, Scope: cloudmodel.PlatformRootScope},
		{Action: cloudmodel.ActionCloudProviderRead, Scope: cloudmodel.PlatformRootScope},
		{Action: cloudmodel.ActionTopologyWrite, Scope: providerScope(providerUID)},
		{Action: cloudmodel.ActionTopologyRead, Scope: providerScope(providerUID)},
		{Action: cloudmodel.ActionParticipationRead, Scope: platformScope(platformUID)},
		{Action: cloudmodel.ActionParticipationAcceptPlatform, Scope: platformScope(platformUID)},
		{Action: cloudmodel.ActionParticipationRejectPlatform, Scope: platformScope(platformUID)},
		{Action: cloudmodel.ActionParticipationRequestReleasePlatform, Scope: platformScope(platformUID)},
		{Action: cloudmodel.ActionParticipationWithdrawProvider, Scope: providerScope(providerUID)},
		{Action: cloudmodel.ActionParticipationAcceptReleaseProvider, Scope: providerScope(providerUID)},
		{Action: cloudmodel.ActionParticipationDeclineReleaseProvider, Scope: providerScope(providerUID)},
		{
			Action:                  cloudmodel.ActionParticipationRequestProvider,
			Scope:                   providerScope(providerUID),
			RelatedCloudPlatformUID: platformUID,
		},
		{Action: cloudmodel.ActionSuspendPlatform, Scope: platformScope(platformUID)},
		{Action: cloudmodel.ActionResumePlatform, Scope: platformScope(platformUID)},
	}
	return g
}

func platformOnlyGrants() *f15Grants {
	return &f15Grants{
		principal: f15Principal,
		grants: []cloudmodel.Grant{
			{Action: cloudmodel.ActionCloudPlatformWrite, Scope: cloudmodel.PlatformRootScope},
			{Action: cloudmodel.ActionCloudPlatformRead, Scope: cloudmodel.PlatformRootScope},
			{Action: cloudmodel.ActionCloudProviderWrite, Scope: cloudmodel.PlatformRootScope},
			{Action: cloudmodel.ActionCloudProviderRead, Scope: cloudmodel.PlatformRootScope},
		},
	}
}

type f15Harness struct {
	audit  *cloudmodel.MemoryAuditAppender
	rt     *server.CloudModelRuntime
	mux    *http.ServeMux
	grants *f15Grants
}

func newF15Harness(grants *f15Grants) *f15Harness {
	audit := &cloudmodel.MemoryAuditAppender{}
	rt := server.NewCloudModelRuntime(audit)
	mux := http.NewServeMux()
	server.RegisterCloudModelRoutes(mux, server.NewCloudModelHandlers(rt, grants), nil)
	return &f15Harness{audit: audit, rt: rt, mux: mux, grants: grants}
}

func (h *f15Harness) rebind() {
	h.mux = http.NewServeMux()
	server.RegisterCloudModelRoutes(h.mux, server.NewCloudModelHandlers(h.rt, h.grants), nil)
}

func doF15(t *testing.T, h *f15Harness, method, path string, body []byte, headers map[string]string) *httptest.ResponseRecorder {
	t.Helper()
	var r *http.Request
	if body == nil {
		r = httptest.NewRequest(method, path, nil)
	} else {
		r = httptest.NewRequest(method, path, bytes.NewReader(body))
	}
	r.Header.Set("Authorization", "Bearer "+f15Principal)
	r.Header.Set("X-Sovrunn-Request-ID", "req-f15-conformance")
	for k, v := range headers {
		r.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()
	h.mux.ServeHTTP(rec, r)
	return rec
}

func doF15Unauth(t *testing.T, h *f15Harness, method, path string, body []byte, headers map[string]string) *httptest.ResponseRecorder {
	t.Helper()
	var r *http.Request
	if body == nil {
		r = httptest.NewRequest(method, path, nil)
	} else {
		r = httptest.NewRequest(method, path, bytes.NewReader(body))
	}
	for k, v := range headers {
		r.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()
	h.mux.ServeHTTP(rec, r)
	return rec
}

func decodeF15Problem(t *testing.T, rec *httptest.ResponseRecorder) *apiproblem.Problem {
	t.Helper()
	var p apiproblem.Problem
	if err := json.Unmarshal(rec.Body.Bytes(), &p); err != nil {
		t.Fatalf("problem decode: %v body=%s", err, rec.Body.String())
	}
	return &p
}

func hasF15Violation(prob *apiproblem.Problem, code apiproblem.ViolationCode) bool {
	if prob == nil {
		return false
	}
	for _, v := range prob.Violations {
		if v.Code == code {
			return true
		}
	}
	return false
}

func assertProblem(t *testing.T, rec *httptest.ResponseRecorder, wantCode apiproblem.ErrorCode, wantStatus int, wantViolation apiproblem.ViolationCode) *apiproblem.Problem {
	t.Helper()
	if rec.Code != wantStatus {
		t.Fatalf("status=%d want %d body=%s", rec.Code, wantStatus, rec.Body.String())
	}
	p := decodeF15Problem(t, rec)
	if p.Code != wantCode {
		t.Fatalf("code=%s want %s (%#v)", p.Code, wantCode, p)
	}
	if wantViolation != "" && !hasF15Violation(p, wantViolation) {
		t.Fatalf("missing violation %s in %#v", wantViolation, p.Violations)
	}
	return p
}

func assertNoAuditDelta(t *testing.T, h *f15Harness, before int) {
	t.Helper()
	if h.audit.Len() != before {
		t.Fatalf("audit len=%d want %d", h.audit.Len(), before)
	}
}

func assertAuditDelta(t *testing.T, h *f15Harness, before, delta int) {
	t.Helper()
	if h.audit.Len() != before+delta {
		t.Fatalf("audit len=%d want %d", h.audit.Len(), before+delta)
	}
}

func lastAudit(t *testing.T, h *f15Harness) apiconform.AuditEvent {
	t.Helper()
	events := h.audit.Events()
	if len(events) == 0 {
		t.Fatal("expected AuditEvent")
	}
	return events[len(events)-1]
}

func assertRedactedAudit(t *testing.T, ev apiconform.AuditEvent, requestID string) {
	t.Helper()
	if ev.Record.RequestID != requestID && requestID != "" {
		// Request ID may be middleware-generated; require non-empty correlation.
		if ev.Record.RequestID == "" {
			t.Fatal("AuditEvent missing request correlation")
		}
	}
	raw, _ := json.Marshal(ev)
	s := string(raw)
	for _, secret := range []string{"Bearer ", "password", "raw-secret", "bootstrapGrant"} {
		if strings.Contains(s, secret) && secret != "bootstrapGrant" {
			t.Fatalf("audit leaked %q", secret)
		}
	}
	// Forged-grant claim values must never appear; member name in reason is OK.
	if strings.Contains(s, `"id":"g1"`) || strings.Contains(s, "forged-claim-value") {
		t.Fatal("audit disclosed forged grant claim")
	}
}

// ---------------------------------------------------------------------------
// Bodies and seeding
// ---------------------------------------------------------------------------

func platformBody(name string) []byte {
	return []byte(`{
		"metadata":{"name":"` + name + `"},
		"spec":{"ownerRegistration":{"legalName":"Acme","registrationIdentifier":"R1","jurisdictionCode":"US"}}
	}`)
}

func providerBody(name string) []byte {
	return []byte(`{
		"metadata":{"name":"` + name + `"},
		"spec":{"operatingMarkets":["US"]}
	}`)
}

func hostingLocationBody(name string) []byte {
	return []byte(`{
		"metadata":{"name":"` + name + `"},
		"spec":{"countryCode":"US","locality":"Austin","administrativeAreaCode":"US-TX"}
	}`)
}

func datacenterBody(name, hlName, hlUID string) []byte {
	return []byte(`{
		"metadata":{"name":"` + name + `"},
		"spec":{"hostingLocationRef":{
			"apiVersion":"infrastructure.sovrunn.io/v1alpha1",
			"kind":"HostingLocation","name":"` + hlName + `","uid":"` + hlUID + `"
		}}
	}`)
}

func faultDomainBody(name, dcName, dcUID string) []byte {
	return []byte(`{
		"metadata":{"name":"` + name + `"},
		"spec":{"datacenterRef":{
			"apiVersion":"infrastructure.sovrunn.io/v1alpha1",
			"kind":"Datacenter","name":"` + dcName + `","uid":"` + dcUID + `"
		}}
	}`)
}

func stackBody(name, fdName, fdUID string) []byte {
	return []byte(`{
		"metadata":{"name":"` + name + `"},
		"spec":{"faultDomainRef":{
			"apiVersion":"infrastructure.sovrunn.io/v1alpha1",
			"kind":"FaultDomain","name":"` + fdName + `","uid":"` + fdUID + `"
		}}
	}`)
}

func participationBody(name, platformUID, providerUID string) []byte {
	return []byte(`{
		"metadata":{"name":"` + name + `"},
		"spec":{
			"cloudPlatformRef":{"apiVersion":"core.sovrunn.io/v1alpha1","kind":"CloudPlatform","name":"root-plat","uid":"` + platformUID + `"},
			"cloudProviderRef":{"apiVersion":"core.sovrunn.io/v1alpha1","kind":"CloudProvider","name":"prov","uid":"` + providerUID + `"},
			"environment":"development"
		}
	}`)
}

func jsonHeaders(key string) map[string]string {
	return map[string]string{"Content-Type": "application/json", "Idempotency-Key": key}
}

func patchHeaders(rv string) map[string]string {
	return map[string]string{"Content-Type": f15MediaMergePatch, "If-Match": rv}
}

func actionHeaders(rv, key string) map[string]string {
	return map[string]string{"If-Match": rv, "Idempotency-Key": key}
}

func actionPath(uid string, action cloudmodel.ParticipationAction) string {
	return "/apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations/" + uid + "/actions/" + string(action)
}

func seedPlatformProvider(t *testing.T, h *f15Harness, platformUID, providerUID string) {
	t.Helper()
	if !h.rt.Store.HasCloudPlatform() {
		if _, prob := h.rt.Store.CreateCloudPlatform(model.CloudPlatform{
			TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionCloudPlatform, Kind: model.KindCloudPlatform},
			Metadata: apimeta.ObjectMeta{Name: "root-plat", UID: platformUID, ResourceVersion: "1"},
			Spec: model.CloudPlatformSpec{OwnerRegistration: model.OwnerRegistration{
				LegalName: "Acme", RegistrationIdentifier: "R1", JurisdictionCode: "US",
			}},
			Status: model.CloudPlatformStatus{Phase: model.CloudPlatformPhaseActive},
		}); prob != nil {
			t.Fatalf("seed platform: %#v", prob)
		}
	}
	if _, ok := h.rt.Store.GetCloudProvider(providerUID); ok {
		return
	}
	if _, prob := h.rt.Store.CreateCloudProvider(model.CloudProvider{
		TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionCloudProvider, Kind: model.KindCloudProvider},
		Metadata: apimeta.ObjectMeta{Name: "prov-" + providerUID[:4], UID: providerUID, ResourceVersion: "1"},
		Spec:     model.CloudProviderSpec{OperatingMarkets: []string{"US"}},
		Status:   model.CloudProviderStatus{Phase: model.CloudProviderPhaseActive},
	}); prob != nil {
		t.Fatalf("seed provider: %#v", prob)
	}
}

func seedOtherProvider(t *testing.T, h *f15Harness, providerUID string) {
	t.Helper()
	if _, ok := h.rt.Store.GetCloudProvider(providerUID); ok {
		return
	}
	if _, prob := h.rt.Store.CreateCloudProvider(model.CloudProvider{
		TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionCloudProvider, Kind: model.KindCloudProvider},
		Metadata: apimeta.ObjectMeta{Name: "prov-" + providerUID[:8], UID: providerUID, ResourceVersion: "1"},
		Spec:     model.CloudProviderSpec{OperatingMarkets: []string{"IN"}},
		Status:   model.CloudProviderStatus{Phase: model.CloudProviderPhaseActive},
	}); prob != nil {
		t.Fatalf("seed other provider: %#v", prob)
	}
}

func seedPendingParticipation(t *testing.T, h *f15Harness, name string) model.CloudProviderParticipation {
	t.Helper()
	seedPlatformProvider(t, h, f15PlatformUID, f15ProviderUID)
	rec := doF15(t, h, http.MethodPost, "/apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations",
		participationBody(name, f15PlatformUID, f15ProviderUID), jsonHeaders("seed-"+name))
	if rec.Code != http.StatusCreated {
		t.Fatalf("seed participation status=%d body=%s", rec.Code, rec.Body.String())
	}
	var p model.CloudProviderParticipation
	if err := json.Unmarshal(rec.Body.Bytes(), &p); err != nil {
		t.Fatalf("seed decode: %v", err)
	}
	return p
}

func seedHostingLocation(t *testing.T, h *f15Harness, name string) model.HostingLocation {
	t.Helper()
	seedPlatformProvider(t, h, f15PlatformUID, f15ProviderUID)
	rec := doF15(t, h, http.MethodPost, "/apis/infrastructure.sovrunn.io/v1alpha1/hosting-locations",
		hostingLocationBody(name), jsonHeaders("seed-hl-"+name))
	if rec.Code != http.StatusCreated {
		t.Fatalf("seed HL status=%d body=%s", rec.Code, rec.Body.String())
	}
	var hl model.HostingLocation
	_ = json.Unmarshal(rec.Body.Bytes(), &hl)
	return hl
}

func seedTopologyChain(t *testing.T, h *f15Harness, prefix string) (model.HostingLocation, model.Datacenter, model.FaultDomain, model.InfrastructureStack) {
	t.Helper()
	hl := seedHostingLocation(t, h, prefix+"-loc")
	rec := doF15(t, h, http.MethodPost, "/apis/infrastructure.sovrunn.io/v1alpha1/datacenters",
		datacenterBody(prefix+"-dc", hl.Metadata.Name, hl.Metadata.UID), jsonHeaders("seed-dc-"+prefix))
	if rec.Code != http.StatusCreated {
		t.Fatalf("seed DC: %d %s", rec.Code, rec.Body.String())
	}
	var dc model.Datacenter
	_ = json.Unmarshal(rec.Body.Bytes(), &dc)
	rec = doF15(t, h, http.MethodPost, "/apis/infrastructure.sovrunn.io/v1alpha1/fault-domains",
		faultDomainBody(prefix+"-fd", dc.Metadata.Name, dc.Metadata.UID), jsonHeaders("seed-fd-"+prefix))
	if rec.Code != http.StatusCreated {
		t.Fatalf("seed FD: %d %s", rec.Code, rec.Body.String())
	}
	var fd model.FaultDomain
	_ = json.Unmarshal(rec.Body.Bytes(), &fd)
	rec = doF15(t, h, http.MethodPost, "/apis/infrastructure.sovrunn.io/v1alpha1/infrastructure-stacks",
		stackBody(prefix+"-stack", fd.Metadata.Name, fd.Metadata.UID), jsonHeaders("seed-stack-"+prefix))
	if rec.Code != http.StatusCreated {
		t.Fatalf("seed stack: %d %s", rec.Code, rec.Body.String())
	}
	var st model.InfrastructureStack
	_ = json.Unmarshal(rec.Body.Bytes(), &st)
	return hl, dc, fd, st
}

func providerScopeRef(uid, name string) *apimeta.ScopeRef {
	return &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
		APIVersion: model.APIVersionCloudProvider,
		Kind:       string(apimeta.ScopeCloudProvider),
		Name:       name,
		UID:        uid,
	}}
}

// ---------------------------------------------------------------------------
// VS0-CF-X03 — cross-provider safe denial
// ---------------------------------------------------------------------------

func TestVS0_CF_X03(t *testing.T) {
	// inputs: cross-provider-target-ref
	// expectedState: unchanged
	// expectedError: RESOURCE_NOT_FOUND
	// expectedViolation: VS0_AUTHORIZATION_SAFE_DENIAL
	// expectedSideEffects: exactly one redacted AuditEvent; append failure → INTERNAL_ERROR/500
	// gate: security
	h := newF15Harness(fullAdminGrants(f15ProviderUID, f15PlatformUID))
	seedPlatformProvider(t, h, f15PlatformUID, f15ProviderUID)
	seedOtherProvider(t, h, f15OtherProviderUID)
	otherHL, prob := h.rt.Store.CreateHostingLocation(model.HostingLocation{
		TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionHostingLocation, Kind: model.KindHostingLocation},
		Metadata: apimeta.ObjectMeta{
			Name: "loc-other", UID: "eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee", ResourceVersion: "1",
			ScopeRef: providerScopeRef(f15OtherProviderUID, "prov-b"),
		},
		Spec:   model.HostingLocationSpec{CountryCode: "IN", Locality: "Mumbai"},
		Status: model.HostingLocationStatus{Phase: model.HostingLocationPhaseActive},
	})
	if prob != nil {
		t.Fatalf("seed other HL: %#v", prob)
	}

	before := h.audit.Len()
	rec := doF15(t, h, http.MethodPost, "/apis/infrastructure.sovrunn.io/v1alpha1/datacenters",
		datacenterBody("dc-x", otherHL.Metadata.Name, otherHL.Metadata.UID), jsonHeaders("x03-inacc"))
	assertProblem(t, rec, apiproblem.CodeResourceNotFound, http.StatusNotFound, f15ViolationSafeDeny)
	assertAuditDelta(t, h, before, 1)
	ev := lastAudit(t, h)
	assertRedactedAudit(t, ev, "req-f15-conformance")
	if strings.Contains(rec.Body.String(), f15OtherProviderUID) || strings.Contains(rec.Body.String(), "loc-other") {
		t.Fatal("safe denial must not disclose target existence")
	}
	for _, d := range h.rt.Store.ListDatacenters() {
		if d.Metadata.Name == "dc-x" {
			t.Fatal("expectedState unchanged: no resource persisted")
		}
	}

	h.audit.SetFail(errors.New("append boom"))
	before = h.audit.Len()
	rec = doF15(t, h, http.MethodPost, "/apis/infrastructure.sovrunn.io/v1alpha1/datacenters",
		datacenterBody("dc-x2", otherHL.Metadata.Name, otherHL.Metadata.UID), jsonHeaders("x03-fail"))
	assertProblem(t, rec, apiproblem.CodeInternalError, http.StatusInternalServerError, "")
	assertNoAuditDelta(t, h, before)
	h.audit.SetFail(nil)
}

// ---------------------------------------------------------------------------
// VS0-CF-F15-01 — authorized creation at registry-declared scopes
// ---------------------------------------------------------------------------

func TestVS0_CF_F15_01(t *testing.T) {
	h := newF15Harness(fullAdminGrants(f15ProviderUID, f15PlatformUID))
	seedPlatformProvider(t, h, f15PlatformUID, f15ProviderUID)

	// Create remaining kinds via HTTP (platform/provider already seeded with known UIDs).
	hl, dc, fd, st := seedTopologyChain(t, h, "f15-01")
	if hl.Metadata.ScopeRef == nil || hl.Metadata.ScopeRef.UID != f15ProviderUID {
		t.Fatalf("HostingLocation scope: %#v", hl.Metadata.ScopeRef)
	}
	if dc.Spec.HostingLocationRef.UID != hl.Metadata.UID || dc.Metadata.ScopeRef.UID != f15ProviderUID {
		t.Fatalf("Datacenter refs/scope: %#v %#v", dc.Spec, dc.Metadata.ScopeRef)
	}
	if fd.Spec.DatacenterRef.UID != dc.Metadata.UID || st.Spec.FaultDomainRef.UID != fd.Metadata.UID {
		t.Fatal("topology chain uid-pinned references required")
	}

	part := seedPendingParticipation(t, h, "f15-01-part")
	if part.Metadata.ScopeRef == nil || part.Metadata.ScopeRef.UID != f15PlatformUID {
		t.Fatalf("participation scope: %#v", part.Metadata.ScopeRef)
	}

	// FEATURE-0015 ends at InfrastructureStack — no ExecutionTarget route.
	for _, p := range server.CloudModelRoutePatterns() {
		if strings.Contains(strings.ToLower(p), "execution-target") {
			t.Fatalf("ExecutionTarget route introduced: %s", p)
		}
	}
}

// ---------------------------------------------------------------------------
// VS0-CF-F15-02 — scope kind outside registry-declared subset
// ---------------------------------------------------------------------------

func TestVS0_CF_F15_02(t *testing.T) {
	h := newF15Harness(platformOnlyGrants())
	body := []byte(`{
		"metadata":{"name":"plat-bad-scope","scopeRef":{"apiVersion":"v","kind":"Organization","name":"o","uid":"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"}},
		"spec":{"ownerRegistration":{"legalName":"Acme","registrationIdentifier":"R1","jurisdictionCode":"US"}}
	}`)
	before := h.audit.Len()
	rec := doF15(t, h, http.MethodPost, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms", body, jsonHeaders("f15-02"))
	assertProblem(t, rec, apiproblem.CodeValidationFailed, http.StatusUnprocessableEntity, validate.ViolationScopeKindInvalid)
	assertNoAuditDelta(t, h, before)
}

// ---------------------------------------------------------------------------
// VS0-CF-F15-03 — provider request + platform accept → Active
// ---------------------------------------------------------------------------

func TestVS0_CF_F15_03(t *testing.T) {
	h := newF15Harness(fullAdminGrants(f15ProviderUID, f15PlatformUID))
	part := seedPendingParticipation(t, h, "f15-03")
	rec := doF15(t, h, http.MethodPost, actionPath(part.Metadata.UID, cloudmodel.ParticipationActionAccept), nil,
		actionHeaders(part.Metadata.ResourceVersion, "f15-03-accept"))
	if rec.Code != http.StatusOK {
		t.Fatalf("accept status=%d body=%s", rec.Code, rec.Body.String())
	}
	_ = json.Unmarshal(rec.Body.Bytes(), &part)
	if part.Status.Phase != model.ParticipationPhaseActive {
		t.Fatalf("expectedState participation Active, got %q", part.Status.Phase)
	}
}

// ---------------------------------------------------------------------------
// VS0-CF-F15-04 — duplicate non-terminal participation pair
// ---------------------------------------------------------------------------

func TestVS0_CF_F15_04(t *testing.T) {
	h := newF15Harness(fullAdminGrants(f15ProviderUID, f15PlatformUID))
	_ = seedPendingParticipation(t, h, "f15-04-a")
	before := h.audit.Len()
	rec := doF15(t, h, http.MethodPost, "/apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations",
		participationBody("f15-04-b", f15PlatformUID, f15ProviderUID), jsonHeaders("f15-04-dup"))
	assertProblem(t, rec, apiproblem.CodeAlreadyExists, http.StatusConflict, validate.ViolationParticipationDuplicate)
	assertNoAuditDelta(t, h, before)
	count := 0
	for _, p := range h.rt.Store.ListParticipations() {
		if p.Spec.CloudPlatformRef.UID == f15PlatformUID && p.Spec.CloudProviderRef.UID == f15ProviderUID {
			count++
		}
	}
	if count != 1 {
		t.Fatalf("no second participation persisted, got %d", count)
	}
}

// ---------------------------------------------------------------------------
// VS0-CF-F15-05 — client status write denied + audited
// ---------------------------------------------------------------------------

func TestVS0_CF_F15_05(t *testing.T) {
	h := newF15Harness(fullAdminGrants(f15ProviderUID, f15PlatformUID))
	seedPlatformProvider(t, h, f15PlatformUID, f15ProviderUID)
	plat, _ := h.rt.Store.GetCloudPlatform(f15PlatformUID)
	before := h.audit.Len()
	rec := doF15(t, h, http.MethodPatch, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms/"+plat.Metadata.UID,
		[]byte(`{"status":{"phase":"Active"}}`), patchHeaders(plat.Metadata.ResourceVersion))
	assertProblem(t, rec, apiproblem.CodeAuthorizationDenied, http.StatusForbidden, validate.ViolationStatusFieldWrite)
	assertAuditDelta(t, h, before, 1)
	assertRedactedAudit(t, lastAudit(t, h), "")

	h.audit.SetFail(errors.New("down"))
	before = h.audit.Len()
	rec = doF15(t, h, http.MethodPatch, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms/"+plat.Metadata.UID,
		[]byte(`{"status":{"phase":"Active"}}`), patchHeaders(plat.Metadata.ResourceVersion))
	assertProblem(t, rec, apiproblem.CodeInternalError, http.StatusInternalServerError, "")
	assertNoAuditDelta(t, h, before)
	got, _ := h.rt.Store.GetCloudPlatform(f15PlatformUID)
	if got.Metadata.ResourceVersion != plat.Metadata.ResourceVersion {
		t.Fatal("append failure must leave resource unchanged")
	}
	h.audit.SetFail(nil)
}

// ---------------------------------------------------------------------------
// VS0-CF-F15-06 — client system-owned metadata write denied + audited
// ---------------------------------------------------------------------------

func TestVS0_CF_F15_06(t *testing.T) {
	h := newF15Harness(platformOnlyGrants())
	before := h.audit.Len()
	body := []byte(`{
		"metadata":{"name":"plat-sys","uid":"dddddddddddddddddddddddddddddddd"},
		"spec":{"ownerRegistration":{"legalName":"Acme","registrationIdentifier":"R1","jurisdictionCode":"US"}}
	}`)
	rec := doF15(t, h, http.MethodPost, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms", body, jsonHeaders("f15-06"))
	assertProblem(t, rec, apiproblem.CodeAuthorizationDenied, http.StatusForbidden, validate.ViolationSystemOwnedFieldWrite)
	assertAuditDelta(t, h, before, 1)

	h.audit.SetFail(errors.New("down"))
	before = h.audit.Len()
	rec = doF15(t, h, http.MethodPost, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms", body, jsonHeaders("f15-06b"))
	assertProblem(t, rec, apiproblem.CodeInternalError, http.StatusInternalServerError, "")
	assertNoAuditDelta(t, h, before)
	h.audit.SetFail(nil)
}

// ---------------------------------------------------------------------------
// VS0-CF-F15-07 — PUT/DELETE → 405; only merge-patch PATCH accepted
// ---------------------------------------------------------------------------

func TestVS0_CF_F15_07(t *testing.T) {
	h := newF15Harness(fullAdminGrants(f15ProviderUID, f15PlatformUID))
	seedPlatformProvider(t, h, f15PlatformUID, f15ProviderUID)
	uid := f15PlatformUID
	paths := []string{
		"/apis/core.sovrunn.io/v1alpha1/cloud-platforms/" + uid,
		"/apis/core.sovrunn.io/v1alpha1/cloud-providers/" + f15ProviderUID,
	}
	hl := seedHostingLocation(t, h, "f15-07-loc")
	paths = append(paths,
		"/apis/infrastructure.sovrunn.io/v1alpha1/hosting-locations/"+hl.Metadata.UID,
	)
	for _, path := range paths {
		for _, method := range []string{http.MethodPut, http.MethodDelete} {
			rec := doF15(t, h, method, path, nil, nil)
			if rec.Code != http.StatusMethodNotAllowed {
				t.Fatalf("%s %s status=%d want 405", method, path, rec.Code)
			}
		}
	}
	// Unsupported PATCH media type → 415 (covered also by F15-14).
	rec := doF15(t, h, http.MethodPatch, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms/"+uid,
		[]byte(`{"spec":{"description":"x"}}`), map[string]string{
			"Content-Type": "application/json", "If-Match": "1",
		})
	if rec.Code != http.StatusUnsupportedMediaType {
		t.Fatalf("unsupported PATCH media status=%d", rec.Code)
	}
}

// ---------------------------------------------------------------------------
// VS0-CF-F15-08 — resume clears only acting hold
// ---------------------------------------------------------------------------

func TestVS0_CF_F15_08(t *testing.T) {
	// Platform suspend + provider suspend, then platform resume → still Suspended.
	grants := fullAdminGrants(f15ProviderUID, f15PlatformUID)
	// Add provider suspend/resume for dual-hold setup via separate harness steps.
	grants.grants = append(grants.grants,
		cloudmodel.Grant{Action: cloudmodel.ActionSuspendProvider, Scope: providerScope(f15ProviderUID)},
		cloudmodel.Grant{Action: cloudmodel.ActionResumeProvider, Scope: providerScope(f15ProviderUID)},
	)
	// Exactly one suspend party at a time: use platform-only for first suspend,
	// then swap to provider-only for second, then platform resume.
	h := newF15Harness(fullAdminGrants(f15ProviderUID, f15PlatformUID))
	part := seedPendingParticipation(t, h, "f15-08")
	rec := doF15(t, h, http.MethodPost, actionPath(part.Metadata.UID, cloudmodel.ParticipationActionAccept), nil,
		actionHeaders(part.Metadata.ResourceVersion, "f15-08-accept"))
	_ = json.Unmarshal(rec.Body.Bytes(), &part)

	rec = doF15(t, h, http.MethodPost, actionPath(part.Metadata.UID, cloudmodel.ParticipationActionSuspend), nil,
		actionHeaders(part.Metadata.ResourceVersion, "f15-08-sus-plat"))
	if rec.Code != http.StatusOK {
		t.Fatalf("platform suspend: %d %s", rec.Code, rec.Body.String())
	}
	_ = json.Unmarshal(rec.Body.Bytes(), &part)
	if !part.Status.PlatformSuspended || part.Status.ProviderSuspended {
		t.Fatalf("after platform suspend: %#v", part.Status)
	}

	// Switch to provider-only suspend/resume grants (exact-one).
	h.grants = &f15Grants{
		principal: f15Principal,
		grants: append(baseReadWrite(f15ProviderUID, f15PlatformUID),
			cloudmodel.Grant{Action: cloudmodel.ActionSuspendProvider, Scope: providerScope(f15ProviderUID)},
			cloudmodel.Grant{Action: cloudmodel.ActionResumeProvider, Scope: providerScope(f15ProviderUID)},
			cloudmodel.Grant{Action: cloudmodel.ActionParticipationAcceptPlatform, Scope: platformScope(f15PlatformUID)},
		),
	}
	h.rebind()
	rec = doF15(t, h, http.MethodPost, actionPath(part.Metadata.UID, cloudmodel.ParticipationActionSuspend), nil,
		actionHeaders(part.Metadata.ResourceVersion, "f15-08-sus-prov"))
	if rec.Code != http.StatusOK {
		t.Fatalf("provider suspend: %d %s", rec.Code, rec.Body.String())
	}
	_ = json.Unmarshal(rec.Body.Bytes(), &part)
	if !part.Status.PlatformSuspended || !part.Status.ProviderSuspended {
		t.Fatalf("both holds: %#v", part.Status)
	}

	// Platform resume only.
	h.grants = fullAdminGrants(f15ProviderUID, f15PlatformUID)
	h.rebind()
	rec = doF15(t, h, http.MethodPost, actionPath(part.Metadata.UID, cloudmodel.ParticipationActionResume), nil,
		actionHeaders(part.Metadata.ResourceVersion, "f15-08-res-plat"))
	if rec.Code != http.StatusOK {
		t.Fatalf("platform resume: %d %s", rec.Code, rec.Body.String())
	}
	_ = json.Unmarshal(rec.Body.Bytes(), &part)
	if part.Status.PlatformSuspended || !part.Status.ProviderSuspended {
		t.Fatalf("only acting hold cleared: %#v", part.Status)
	}
	if part.Status.Phase != model.ParticipationPhaseSuspended {
		t.Fatalf("remains Suspended, got %q", part.Status.Phase)
	}
}

func baseReadWrite(providerUID, platformUID string) []cloudmodel.Grant {
	return []cloudmodel.Grant{
		{Action: cloudmodel.ActionCloudPlatformWrite, Scope: cloudmodel.PlatformRootScope},
		{Action: cloudmodel.ActionCloudPlatformRead, Scope: cloudmodel.PlatformRootScope},
		{Action: cloudmodel.ActionCloudProviderWrite, Scope: cloudmodel.PlatformRootScope},
		{Action: cloudmodel.ActionCloudProviderRead, Scope: cloudmodel.PlatformRootScope},
		{Action: cloudmodel.ActionTopologyWrite, Scope: providerScope(providerUID)},
		{Action: cloudmodel.ActionTopologyRead, Scope: providerScope(providerUID)},
		{Action: cloudmodel.ActionParticipationRead, Scope: platformScope(platformUID)},
		{
			Action: cloudmodel.ActionParticipationRequestProvider, Scope: providerScope(providerUID),
			RelatedCloudPlatformUID: platformUID,
		},
	}
}

// ---------------------------------------------------------------------------
// VS0-CF-F15-09 — accept without participation.accept.platform
// ---------------------------------------------------------------------------

func TestVS0_CF_F15_09(t *testing.T) {
	h := newF15Harness(fullAdminGrants(f15ProviderUID, f15PlatformUID))
	part := seedPendingParticipation(t, h, "f15-09")
	// Drop accept.platform
	h.grants = &f15Grants{principal: f15Principal, grants: baseReadWrite(f15ProviderUID, f15PlatformUID)}
	h.rebind()
	before := h.audit.Len()
	rec := doF15(t, h, http.MethodPost, actionPath(part.Metadata.UID, cloudmodel.ParticipationActionAccept), nil,
		actionHeaders(part.Metadata.ResourceVersion, "f15-09-deny"))
	assertProblem(t, rec, apiproblem.CodeAuthorizationDenied, http.StatusForbidden, "")
	assertNoAuditDelta(t, h, before)
	got, _ := h.rt.Store.GetParticipation(part.Metadata.UID)
	if got.Status.Phase != model.ParticipationPhasePending {
		t.Fatalf("remains Pending, got %q", got.Status.Phase)
	}
}

// ---------------------------------------------------------------------------
// VS0-CF-F15-10 — suspend/resume from invalid phases
// ---------------------------------------------------------------------------

func TestVS0_CF_F15_10(t *testing.T) {
	h := newF15Harness(fullAdminGrants(f15ProviderUID, f15PlatformUID))
	part := seedPendingParticipation(t, h, "f15-10")
	before := h.audit.Len()
	rec := doF15(t, h, http.MethodPost, actionPath(part.Metadata.UID, cloudmodel.ParticipationActionSuspend), nil,
		actionHeaders(part.Metadata.ResourceVersion, "f15-10-sus"))
	assertProblem(t, rec, apiproblem.CodeConflict, http.StatusConflict, cloudmodel.ViolationParticipationStateInvalid)
	assertNoAuditDelta(t, h, before)

	// Terminal Rejected → resume denied
	rec = doF15(t, h, http.MethodPost, actionPath(part.Metadata.UID, cloudmodel.ParticipationActionReject), nil,
		actionHeaders(part.Metadata.ResourceVersion, "f15-10-rej"))
	_ = json.Unmarshal(rec.Body.Bytes(), &part)
	before = h.audit.Len()
	rec = doF15(t, h, http.MethodPost, actionPath(part.Metadata.UID, cloudmodel.ParticipationActionResume), nil,
		actionHeaders(part.Metadata.ResourceVersion, "f15-10-res"))
	assertProblem(t, rec, apiproblem.CodeConflict, http.StatusConflict, cloudmodel.ViolationParticipationStateInvalid)
	assertNoAuditDelta(t, h, before)
}

// ---------------------------------------------------------------------------
// VS0-CF-F15-11 — cloudPlatformRef.uid != scopeRef.uid
// ---------------------------------------------------------------------------

func TestVS0_CF_F15_11(t *testing.T) {
	// Server create always derives matching scope; prove invariant rejection
	// via the same StageSet reference validator handlers invoke.
	p := model.CloudProviderParticipation{
		Metadata: apimeta.ObjectMeta{
			Name: "part",
			ScopeRef: &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
				APIVersion: model.APIVersionCloudPlatform, Kind: model.KindCloudPlatform,
				Name: "plat", UID: f15PlatformUID,
			}},
		},
		Spec: model.CloudProviderParticipationSpec{
			CloudPlatformRef: apimeta.TypedRef{
				APIVersion: model.APIVersionCloudPlatform, Kind: model.KindCloudPlatform,
				Name: "plat", UID: "dddddddddddddddddddddddddddddddd",
			},
			CloudProviderRef: apimeta.TypedRef{
				APIVersion: model.APIVersionCloudProvider, Kind: model.KindCloudProvider,
				Name: "prov", UID: f15ProviderUID,
			},
			Environment: model.ParticipationEnvironmentDevelopment,
		},
	}
	prob := validate.ValidateParticipationScopeUIDInvariant(p)
	if prob == nil || prob.Code != apiproblem.CodeValidationFailed {
		t.Fatalf("want VALIDATION_FAILED, got %#v", prob)
	}
	if !hasF15Violation(prob, validate.ViolationScopeReferenceMismatch) {
		t.Fatalf("want VS0_SCOPE_REFERENCE_MISMATCH, got %#v", prob.Violations)
	}
	if prob.Status != http.StatusUnprocessableEntity {
		t.Fatalf("status=%d want 422", prob.Status)
	}

	// Create path never persists a mismatched participation.
	h := newF15Harness(fullAdminGrants(f15ProviderUID, f15PlatformUID))
	created := seedPendingParticipation(t, h, "f15-11-ok")
	if created.Spec.CloudPlatformRef.UID != created.Metadata.ScopeRef.UID {
		t.Fatal("create must keep uid-pinned scope invariant")
	}
}

// ---------------------------------------------------------------------------
// VS0-CF-F15-12 — CloudPlatform-root required
// ---------------------------------------------------------------------------

func TestVS0_CF_F15_12(t *testing.T) {
	h := newF15Harness(fullAdminGrants(f15ProviderUID, f15PlatformUID))
	before := h.audit.Len()
	rec := doF15(t, h, http.MethodPost, "/apis/core.sovrunn.io/v1alpha1/cloud-providers",
		providerBody("prov-early"), jsonHeaders("f15-12-prov"))
	assertProblem(t, rec, apiproblem.CodeConflict, http.StatusConflict, f15ViolationRootReq)
	assertAuditDelta(t, h, before, 1)

	h.audit.SetFail(errors.New("down"))
	before = h.audit.Len()
	rec = doF15(t, h, http.MethodPost, "/apis/infrastructure.sovrunn.io/v1alpha1/hosting-locations",
		hostingLocationBody("loc-early"), jsonHeaders("f15-12-hl"))
	assertProblem(t, rec, apiproblem.CodeInternalError, http.StatusInternalServerError, "")
	assertNoAuditDelta(t, h, before)
	h.audit.SetFail(nil)
}

// ---------------------------------------------------------------------------
// VS0-CF-F15-13 — valid merge-patch of mutable field
// ---------------------------------------------------------------------------

func TestVS0_CF_F15_13(t *testing.T) {
	h := newF15Harness(fullAdminGrants(f15ProviderUID, f15PlatformUID))
	seedPlatformProvider(t, h, f15PlatformUID, f15ProviderUID)
	plat, _ := h.rt.Store.GetCloudPlatform(f15PlatformUID)
	beforeRV := plat.Metadata.ResourceVersion
	rec := doF15(t, h, http.MethodPatch, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms/"+plat.Metadata.UID,
		[]byte(`{"spec":{"description":"updated"}}`), patchHeaders(plat.Metadata.ResourceVersion))
	if rec.Code != http.StatusOK {
		t.Fatalf("patch status=%d body=%s", rec.Code, rec.Body.String())
	}
	_ = json.Unmarshal(rec.Body.Bytes(), &plat)
	if plat.Spec.Description != "updated" || plat.Metadata.ResourceVersion == beforeRV {
		t.Fatalf("mutable field/resourceVersion: %#v", plat)
	}
	if plat.Spec.OwnerRegistration.LegalName != "Acme" {
		t.Fatal("no other field changes")
	}
}

// ---------------------------------------------------------------------------
// VS0-CF-F15-14 — unsupported media / immutable / status PATCH
// ---------------------------------------------------------------------------

func TestVS0_CF_F15_14(t *testing.T) {
	h := newF15Harness(fullAdminGrants(f15ProviderUID, f15PlatformUID))
	seedPlatformProvider(t, h, f15PlatformUID, f15ProviderUID)
	plat, _ := h.rt.Store.GetCloudPlatform(f15PlatformUID)

	rec := doF15(t, h, http.MethodPatch, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms/"+plat.Metadata.UID,
		[]byte(`{"spec":{"description":"x"}}`), map[string]string{
			"Content-Type": "text/plain", "If-Match": plat.Metadata.ResourceVersion,
		})
	assertProblem(t, rec, apiproblem.CodeUnsupportedMediaType, http.StatusUnsupportedMediaType, "")

	before := h.audit.Len()
	rec = doF15(t, h, http.MethodPatch, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms/"+plat.Metadata.UID,
		[]byte(`{"spec":{"ownerRegistration":{"legalName":"X","registrationIdentifier":"R1","jurisdictionCode":"US"}}}`),
		patchHeaders(plat.Metadata.ResourceVersion))
	assertProblem(t, rec, apiproblem.CodeValidationFailed, http.StatusUnprocessableEntity, validate.ViolationPatchImmutableField)
	assertNoAuditDelta(t, h, before)

	before = h.audit.Len()
	rec = doF15(t, h, http.MethodPatch, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms/"+plat.Metadata.UID,
		[]byte(`{"status":{"phase":"Active"}}`), patchHeaders(plat.Metadata.ResourceVersion))
	assertProblem(t, rec, apiproblem.CodeAuthorizationDenied, http.StatusForbidden, validate.ViolationStatusFieldWrite)
	assertAuditDelta(t, h, before, 1)
}

// ---------------------------------------------------------------------------
// VS0-CF-F15-15 — participation create initializes Pending + 7-day expiry
// ---------------------------------------------------------------------------

func TestVS0_CF_F15_15(t *testing.T) {
	h := newF15Harness(fullAdminGrants(f15ProviderUID, f15PlatformUID))
	part := seedPendingParticipation(t, h, "f15-15")
	if part.Status.Phase != model.ParticipationPhasePending {
		t.Fatalf("phase=%q", part.Status.Phase)
	}
	if part.Status.PlatformSuspended || part.Status.ProviderSuspended {
		t.Fatalf("holds must be false: %#v", part.Status)
	}
	createdAt, _ := time.Parse(time.RFC3339, part.Metadata.CreatedAt)
	expires, _ := time.Parse(time.RFC3339, part.Status.RequestExpiresAt)
	if !expires.Equal(createdAt.Add(7 * 24 * time.Hour)) {
		t.Fatalf("requestExpiresAt=%v want createdAt+7d=%v", expires, createdAt.Add(7*24*time.Hour))
	}
}

// ---------------------------------------------------------------------------
// VS0-CF-F15-16 — scheduler expiry → Expired, idempotent, one AuditEvent
// ---------------------------------------------------------------------------

func TestVS0_CF_F15_16(t *testing.T) {
	h := newF15Harness(fullAdminGrants(f15ProviderUID, f15PlatformUID))
	part := seedPendingParticipation(t, h, "f15-16")
	before := h.audit.Len()
	if err := h.rt.Publication.ExpireParticipation(context.Background(), part.Metadata.UID, part.Metadata.ResourceVersion, "sched-f15-16"); err != nil {
		t.Fatalf("ExpireParticipation: %v", err)
	}
	got, _ := h.rt.Store.GetParticipation(part.Metadata.UID)
	if got.Status.Phase != model.ParticipationPhaseExpired {
		t.Fatalf("expected Expired, got %q", got.Status.Phase)
	}
	assertAuditDelta(t, h, before, 1)
	ev := lastAudit(t, h)
	if ev.Record.Action != cloudmodel.AuditActionParticipationExpiry {
		t.Fatalf("action=%s", ev.Record.Action)
	}
	if ev.Record.RequestID != "sched-f15-16" {
		t.Fatalf("scheduler execution ID correlation=%q", ev.Record.RequestID)
	}
	// Idempotent repeat
	before = h.audit.Len()
	_ = h.rt.Publication.ExpireParticipation(context.Background(), part.Metadata.UID, got.Metadata.ResourceVersion, "sched-f15-16-b")
	assertNoAuditDelta(t, h, before)
}

// ---------------------------------------------------------------------------
// VS0-CF-F15-17 — missing/malformed/stale If-Match on action
// ---------------------------------------------------------------------------

func TestVS0_CF_F15_17(t *testing.T) {
	h := newF15Harness(fullAdminGrants(f15ProviderUID, f15PlatformUID))
	part := seedPendingParticipation(t, h, "f15-17")
	before := h.audit.Len()

	rec := doF15(t, h, http.MethodPost, actionPath(part.Metadata.UID, cloudmodel.ParticipationActionAccept), nil,
		map[string]string{"Idempotency-Key": "f15-17-miss"})
	assertProblem(t, rec, apiproblem.CodeStaleResourceVersion, http.StatusPreconditionFailed, "")
	assertNoAuditDelta(t, h, before)

	rec = doF15(t, h, http.MethodPost, actionPath(part.Metadata.UID, cloudmodel.ParticipationActionAccept), nil,
		actionHeaders("stale", "f15-17-stale"))
	assertProblem(t, rec, apiproblem.CodeStaleResourceVersion, http.StatusPreconditionFailed, "")
	assertNoAuditDelta(t, h, before)
}

// ---------------------------------------------------------------------------
// VS0-CF-F15-18 — same-key/same-digest replay
// ---------------------------------------------------------------------------

func TestVS0_CF_F15_18(t *testing.T) {
	h := newF15Harness(platformOnlyGrants())
	headers := jsonHeaders("f15-18-plat")
	rec := doF15(t, h, http.MethodPost, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms", platformBody("f15-18"), headers)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create: %d %s", rec.Code, rec.Body.String())
	}
	first := append([]byte(nil), rec.Body.Bytes()...)
	before := h.audit.Len()
	rec = doF15(t, h, http.MethodPost, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms", platformBody("f15-18"), headers)
	if rec.Code != http.StatusCreated || !bytes.Equal(rec.Body.Bytes(), first) {
		t.Fatalf("replay status=%d body mismatch", rec.Code)
	}
	assertNoAuditDelta(t, h, before)
	if len(h.rt.Store.ListCloudPlatforms()) != 1 {
		t.Fatal("no duplicate resource")
	}

	// Participation action replay
	h2 := newF15Harness(fullAdminGrants(f15ProviderUID, f15PlatformUID))
	part := seedPendingParticipation(t, h2, "f15-18-part")
	ah := actionHeaders(part.Metadata.ResourceVersion, "f15-18-accept")
	rec = doF15(t, h2, http.MethodPost, actionPath(part.Metadata.UID, cloudmodel.ParticipationActionAccept), nil, ah)
	if rec.Code != http.StatusOK {
		t.Fatalf("accept: %d", rec.Code)
	}
	first = append([]byte(nil), rec.Body.Bytes()...)
	before = h2.audit.Len()
	rec = doF15(t, h2, http.MethodPost, actionPath(part.Metadata.UID, cloudmodel.ParticipationActionAccept), nil, ah)
	if rec.Code != http.StatusOK || !bytes.Equal(rec.Body.Bytes(), first) {
		t.Fatalf("action replay mismatch")
	}
	assertNoAuditDelta(t, h2, before)

	// Revoked grant denies without stored-result disclosure
	h2.grants = &f15Grants{principal: f15Principal, grants: baseReadWrite(f15ProviderUID, f15PlatformUID)}
	h2.rebind()
	rec = doF15(t, h2, http.MethodPost, actionPath(part.Metadata.UID, cloudmodel.ParticipationActionAccept), nil, ah)
	if rec.Code == http.StatusOK || bytes.Equal(rec.Body.Bytes(), first) {
		t.Fatal("revoked grant must not disclose stored result")
	}
}

// ---------------------------------------------------------------------------
// VS0-CF-F15-19 — same-key/changed-digest conflict
// ---------------------------------------------------------------------------

func TestVS0_CF_F15_19(t *testing.T) {
	h := newF15Harness(platformOnlyGrants())
	headers := jsonHeaders("f15-19")
	rec := doF15(t, h, http.MethodPost, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms", platformBody("f15-19-a"), headers)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create: %d", rec.Code)
	}
	rec = doF15(t, h, http.MethodPost, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms", platformBody("f15-19-b"), headers)
	assertProblem(t, rec, apiproblem.CodeConflict, http.StatusConflict, cloudmodel.ViolationIdempotencyKeyReuseMismatch)
	if len(h.rt.Store.ListCloudPlatforms()) != 1 {
		t.Fatal("original result retained")
	}
}

// ---------------------------------------------------------------------------
// VS0-CF-F15-20 — authorized lifecycle actions + invalid source state
// ---------------------------------------------------------------------------

func TestVS0_CF_F15_20(t *testing.T) {
	h := newF15Harness(fullAdminGrants(f15ProviderUID, f15PlatformUID))
	part := seedPendingParticipation(t, h, "f15-20-rej")
	rec := doF15(t, h, http.MethodPost, actionPath(part.Metadata.UID, cloudmodel.ParticipationActionReject), nil,
		actionHeaders(part.Metadata.ResourceVersion, "f15-20-rej"))
	_ = json.Unmarshal(rec.Body.Bytes(), &part)
	if part.Status.Phase != model.ParticipationPhaseRejected {
		t.Fatalf("reject → Rejected, got %q", part.Status.Phase)
	}

	part = seedPendingParticipationOnProvider(t, h, "f15-20-wd", "dddddddddddddddddddddddddddddddd")
	rec = doF15(t, h, http.MethodPost, actionPath(part.Metadata.UID, cloudmodel.ParticipationActionWithdraw), nil,
		actionHeaders(part.Metadata.ResourceVersion, "f15-20-wd"))
	_ = json.Unmarshal(rec.Body.Bytes(), &part)
	if part.Status.Phase != model.ParticipationPhaseWithdrawn {
		t.Fatalf("withdraw → Withdrawn, got %q", part.Status.Phase)
	}

	part = seedPendingParticipationOnProvider(t, h, "f15-20-rel", "eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee")
	rec = doF15(t, h, http.MethodPost, actionPath(part.Metadata.UID, cloudmodel.ParticipationActionAccept), nil,
		actionHeaders(part.Metadata.ResourceVersion, "f15-20-acc"))
	_ = json.Unmarshal(rec.Body.Bytes(), &part)
	rec = doF15(t, h, http.MethodPost, actionPath(part.Metadata.UID, cloudmodel.ParticipationActionRequestRelease), nil,
		actionHeaders(part.Metadata.ResourceVersion, "f15-20-rr"))
	_ = json.Unmarshal(rec.Body.Bytes(), &part)
	if part.Status.Phase != model.ParticipationPhaseTerminating {
		t.Fatalf("request-release → Terminating, got %q", part.Status.Phase)
	}
	rec = doF15(t, h, http.MethodPost, actionPath(part.Metadata.UID, cloudmodel.ParticipationActionDeclineRelease), nil,
		actionHeaders(part.Metadata.ResourceVersion, "f15-20-dr"))
	_ = json.Unmarshal(rec.Body.Bytes(), &part)
	if part.Status.Phase != model.ParticipationPhaseActive {
		t.Fatalf("decline-release → Active, got %q", part.Status.Phase)
	}

	// Invalid source: reject from Active
	before := h.audit.Len()
	rec = doF15(t, h, http.MethodPost, actionPath(part.Metadata.UID, cloudmodel.ParticipationActionReject), nil,
		actionHeaders(part.Metadata.ResourceVersion, "f15-20-bad"))
	assertProblem(t, rec, apiproblem.CodeConflict, http.StatusConflict, cloudmodel.ViolationParticipationStateInvalid)
	assertNoAuditDelta(t, h, before)
}

func seedPendingParticipationOnProvider(t *testing.T, h *f15Harness, name, providerUID string) model.CloudProviderParticipation {
	t.Helper()
	seedPlatformProvider(t, h, f15PlatformUID, f15ProviderUID)
	if _, ok := h.rt.Store.GetCloudProvider(providerUID); !ok {
		seedOtherProvider(t, h, providerUID)
	}
	// Extend grants for the new provider relation.
	h.grants.grants = append(h.grants.grants,
		cloudmodel.Grant{Action: cloudmodel.ActionTopologyWrite, Scope: providerScope(providerUID)},
		cloudmodel.Grant{Action: cloudmodel.ActionTopologyRead, Scope: providerScope(providerUID)},
		cloudmodel.Grant{Action: cloudmodel.ActionParticipationWithdrawProvider, Scope: providerScope(providerUID)},
		cloudmodel.Grant{Action: cloudmodel.ActionParticipationAcceptReleaseProvider, Scope: providerScope(providerUID)},
		cloudmodel.Grant{Action: cloudmodel.ActionParticipationDeclineReleaseProvider, Scope: providerScope(providerUID)},
		cloudmodel.Grant{
			Action: cloudmodel.ActionParticipationRequestProvider, Scope: providerScope(providerUID),
			RelatedCloudPlatformUID: f15PlatformUID,
		},
	)
	h.rebind()
	rec := doF15(t, h, http.MethodPost, "/apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations",
		participationBody(name, f15PlatformUID, providerUID), jsonHeaders("seed-"+name))
	if rec.Code != http.StatusCreated {
		t.Fatalf("seed part on provider: %d %s", rec.Code, rec.Body.String())
	}
	var p model.CloudProviderParticipation
	_ = json.Unmarshal(rec.Body.Bytes(), &p)
	return p
}

// ---------------------------------------------------------------------------
// VS0-CF-F15-21 — LIST/GET authorization boundary
// ---------------------------------------------------------------------------

func TestVS0_CF_F15_21(t *testing.T) {
	h := newF15Harness(platformOnlyGrants())
	rec := doF15(t, h, http.MethodPost, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms", platformBody("f15-21"), jsonHeaders("f15-21"))
	var plat model.CloudPlatform
	_ = json.Unmarshal(rec.Body.Bytes(), &plat)

	// LIST without read grant
	h.grants = &f15Grants{principal: f15Principal, grants: []cloudmodel.Grant{
		{Action: cloudmodel.ActionCloudPlatformWrite, Scope: cloudmodel.PlatformRootScope},
	}}
	h.rebind()
	before := h.audit.Len()
	rec = doF15(t, h, http.MethodGet, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms", nil, nil)
	assertProblem(t, rec, apiproblem.CodeAuthorizationDenied, http.StatusForbidden, "")
	assertAuditDelta(t, h, before, 1)
	if strings.Contains(rec.Body.String(), plat.Metadata.Name) {
		t.Fatal("LIST denial must not include resource detail")
	}

	// Scoped read covering no visible resources → 200 empty
	h.grants = &f15Grants{principal: f15Principal, grants: []cloudmodel.Grant{
		{Action: cloudmodel.ActionCloudPlatformRead, Scope: cloudmodel.PlatformRootScope, ResourceUID: "ffffffffffffffffffffffffffffffff"},
	}}
	h.rebind()
	rec = doF15(t, h, http.MethodGet, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms", nil, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("empty scoped LIST status=%d", rec.Code)
	}
	var list apimeta.ListEnvelope[model.CloudPlatform]
	_ = json.Unmarshal(rec.Body.Bytes(), &list)
	if len(list.Items) != 0 {
		t.Fatalf("want empty collection, got %d", len(list.Items))
	}

	// Unauthorized direct GET → safe 404
	rec = doF15(t, h, http.MethodGet, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms/"+plat.Metadata.UID, nil, nil)
	assertProblem(t, rec, apiproblem.CodeResourceNotFound, http.StatusNotFound, f15ViolationSafeDeny)
}

// ---------------------------------------------------------------------------
// VS0-CF-F15-22 — forged bootstrap grant header/body
// ---------------------------------------------------------------------------

func TestVS0_CF_F15_22(t *testing.T) {
	h := newF15Harness(platformOnlyGrants())
	before := h.audit.Len()
	headers := jsonHeaders("f15-22-hdr")
	headers[cloudmodel.HeaderBootstrapGrant] = "forged-claim-value"
	rec := doF15(t, h, http.MethodPost, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms", platformBody("f15-22"), headers)
	assertProblem(t, rec, apiproblem.CodeAuthorizationDenied, http.StatusForbidden, "")
	assertAuditDelta(t, h, before, 1)
	assertRedactedAudit(t, lastAudit(t, h), "")
	if strings.Contains(rec.Body.String(), "forged-claim-value") {
		t.Fatal("must not disclose forged claim")
	}

	before = h.audit.Len()
	body := []byte(`{
		"metadata":{"name":"f15-22b"},
		"spec":{"ownerRegistration":{"legalName":"Acme","registrationIdentifier":"R1","jurisdictionCode":"US"}},
		"bootstrapGrant":{"id":"g1"}
	}`)
	rec = doF15(t, h, http.MethodPost, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms", body, jsonHeaders("f15-22-body"))
	assertProblem(t, rec, apiproblem.CodeAuthorizationDenied, http.StatusForbidden, "")
	assertAuditDelta(t, h, before, 1)

	h.audit.SetFail(errors.New("down"))
	before = h.audit.Len()
	headers = jsonHeaders("f15-22-fail")
	headers[cloudmodel.HeaderBootstrapGrant] = "x"
	rec = doF15(t, h, http.MethodPost, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms", platformBody("f15-22c"), headers)
	assertProblem(t, rec, apiproblem.CodeInternalError, http.StatusInternalServerError, "")
	assertNoAuditDelta(t, h, before)
	h.audit.SetFail(nil)
}

// ---------------------------------------------------------------------------
// VS0-CF-F15-23 — inaccessible cross-provider + authorized topology mismatch
// ---------------------------------------------------------------------------

func TestVS0_CF_F15_23(t *testing.T) {
	// Inaccessible cross-provider (same as X03 path)
	h := newF15Harness(fullAdminGrants(f15ProviderUID, f15PlatformUID))
	seedPlatformProvider(t, h, f15PlatformUID, f15ProviderUID)
	seedOtherProvider(t, h, f15OtherProviderUID)
	otherHL, _ := h.rt.Store.CreateHostingLocation(model.HostingLocation{
		TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionHostingLocation, Kind: model.KindHostingLocation},
		Metadata: apimeta.ObjectMeta{
			Name: "loc-x", UID: "eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee", ResourceVersion: "1",
			ScopeRef: providerScopeRef(f15OtherProviderUID, "prov-b"),
		},
		Spec:   model.HostingLocationSpec{CountryCode: "IN", Locality: "Mumbai"},
		Status: model.HostingLocationStatus{Phase: model.HostingLocationPhaseActive},
	})
	before := h.audit.Len()
	rec := doF15(t, h, http.MethodPost, "/apis/infrastructure.sovrunn.io/v1alpha1/datacenters",
		datacenterBody("dc-inacc", otherHL.Metadata.Name, otherHL.Metadata.UID), jsonHeaders("f15-23-a"))
	assertProblem(t, rec, apiproblem.CodeResourceNotFound, http.StatusNotFound, f15ViolationSafeDeny)
	assertAuditDelta(t, h, before, 1)

	// Authorized for both: visible parent under other provider without topology.write → mismatch
	h2 := newF15Harness(&f15Grants{
		principal: f15Principal,
		grants: []cloudmodel.Grant{
			{Action: cloudmodel.ActionCloudPlatformWrite, Scope: cloudmodel.PlatformRootScope},
			{Action: cloudmodel.ActionCloudPlatformRead, Scope: cloudmodel.PlatformRootScope},
			{Action: cloudmodel.ActionCloudProviderWrite, Scope: cloudmodel.PlatformRootScope},
			{Action: cloudmodel.ActionCloudProviderRead, Scope: cloudmodel.PlatformRootScope},
			{Action: cloudmodel.ActionTopologyWrite, Scope: providerScope(f15ProviderUID)},
			{Action: cloudmodel.ActionTopologyRead, Scope: providerScope(f15ProviderUID)},
			{Action: cloudmodel.ActionTopologyRead, Scope: providerScope(f15OtherProviderUID)},
		},
	})
	seedPlatformProvider(t, h2, f15PlatformUID, f15ProviderUID)
	seedOtherProvider(t, h2, f15OtherProviderUID)
	vis, _ := h2.rt.Store.CreateHostingLocation(model.HostingLocation{
		TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionHostingLocation, Kind: model.KindHostingLocation},
		Metadata: apimeta.ObjectMeta{
			Name: "loc-vis", UID: "ffffffffffffffffffffffffffffffff", ResourceVersion: "1",
			ScopeRef: providerScopeRef(f15OtherProviderUID, "prov-b"),
		},
		Spec:   model.HostingLocationSpec{CountryCode: "IN", Locality: "Delhi"},
		Status: model.HostingLocationStatus{Phase: model.HostingLocationPhaseActive},
	})
	before = h2.audit.Len()
	rec = doF15(t, h2, http.MethodPost, "/apis/infrastructure.sovrunn.io/v1alpha1/datacenters",
		datacenterBody("dc-mm", vis.Metadata.Name, vis.Metadata.UID), jsonHeaders("f15-23-b"))
	assertProblem(t, rec, apiproblem.CodeValidationFailed, http.StatusUnprocessableEntity, validate.ViolationTopologyProviderMismatch)
	assertNoAuditDelta(t, h2, before)
}

// ---------------------------------------------------------------------------
// VS0-CF-F15-24 — audit atomicity / append-failure INTERNAL_ERROR
// ---------------------------------------------------------------------------

func TestVS0_CF_F15_24(t *testing.T) {
	h := newF15Harness(platformOnlyGrants())
	h.audit.SetFail(errors.New("append boom"))
	rec := doF15(t, h, http.MethodPost, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms",
		platformBody("f15-24"), jsonHeaders("f15-24"))
	assertProblem(t, rec, apiproblem.CodeInternalError, http.StatusInternalServerError, "")
	if h.audit.Len() != 0 {
		t.Fatal("append failure must not publish AuditEvent")
	}
	if len(h.rt.Store.ListCloudPlatforms()) != 0 {
		t.Fatal("mutation must not publish on append failure")
	}
	if p := decodeF15Problem(t, rec); p.Code == "DEPENDENCY_UNAVAILABLE" {
		t.Fatal("DEPENDENCY_UNAVAILABLE must not be used")
	}
	h.audit.SetFail(nil)

	// Successful create publishes exactly one correlated AuditEvent
	before := h.audit.Len()
	rec = doF15(t, h, http.MethodPost, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms",
		platformBody("f15-24-ok"), jsonHeaders("f15-24-ok"))
	if rec.Code != http.StatusCreated {
		t.Fatalf("create: %d", rec.Code)
	}
	assertAuditDelta(t, h, before, 1)
	assertRedactedAudit(t, lastAudit(t, h), "")
}

// ---------------------------------------------------------------------------
// VS0-CF-F15-25 — owned route surface; no ExecutionTarget/migration
// ---------------------------------------------------------------------------

func TestVS0_CF_F15_25(t *testing.T) {
	patterns := server.CloudModelRoutePatterns()
	if len(patterns) != server.CloudModelRouteCount {
		t.Fatalf("registrations=%d want %d", len(patterns), server.CloudModelRouteCount)
	}
	logical := map[string]struct{}{}
	for _, p := range patterns {
		parts := strings.SplitN(p, " ", 2)
		if len(parts) != 2 {
			t.Fatalf("malformed pattern %q", p)
		}
		logical[parts[1]] = struct{}{}
		low := strings.ToLower(p)
		for _, bad := range []string{"execution-target", "migration", "service-region"} {
			if strings.Contains(low, bad) {
				t.Fatalf("excluded route: %s", p)
			}
		}
		if strings.HasPrefix(p, "PUT ") || strings.HasPrefix(p, "DELETE ") || strings.HasPrefix(p, "HEAD ") {
			t.Fatalf("owned surface is POST/GET/LIST/PATCH only: %s", p)
		}
		if strings.Contains(p, "cloud-provider-participations/{uid}") && strings.HasPrefix(p, "PATCH ") {
			t.Fatal("participation must not expose PATCH")
		}
	}
	if len(logical) != server.CloudModelLogicalPathCount {
		t.Fatalf("logical paths=%d want %d", len(logical), server.CloudModelLogicalPathCount)
	}

	h := newF15Harness(fullAdminGrants(f15ProviderUID, f15PlatformUID))
	part := seedPendingParticipation(t, h, "f15-25")
	rec := doF15(t, h, http.MethodPost, actionPath(part.Metadata.UID, cloudmodel.ParticipationActionAccept),
		nil, actionHeaders(part.Metadata.ResourceVersion, "f15-25-empty"))
	if rec.Code != http.StatusOK {
		t.Fatalf("empty action body must succeed: %d %s", rec.Code, rec.Body.String())
	}
}

// ---------------------------------------------------------------------------
// VS0-CF-F15-26 — per-kind collection-create closed contract
// ---------------------------------------------------------------------------

func TestVS0_CF_F15_26(t *testing.T) {
	h := newF15Harness(platformOnlyGrants())
	// Unknown field rejected
	before := h.audit.Len()
	rec := doF15(t, h, http.MethodPost, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms",
		[]byte(`{"metadata":{"name":"x"},"spec":{"ownerRegistration":{"legalName":"A","registrationIdentifier":"R","jurisdictionCode":"US"},"extra":1}}`),
		jsonHeaders("f15-26-unk"))
	if rec.Code != http.StatusBadRequest && rec.Code != http.StatusUnprocessableEntity {
		// UNKNOWN_FIELD is typically 400
		p := decodeF15Problem(t, rec)
		if p.Code != apiproblem.CodeUnknownField && p.Code != apiproblem.CodeValidationFailed {
			t.Fatalf("unknown field: %d %#v", rec.Code, p)
		}
	}
	assertNoAuditDelta(t, h, before)

	rec = doF15(t, h, http.MethodPost, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms",
		platformBody("f15-26"), jsonHeaders("f15-26-ok"))
	if rec.Code != http.StatusCreated {
		t.Fatalf("create: %d %s", rec.Code, rec.Body.String())
	}
	var plat model.CloudPlatform
	_ = json.Unmarshal(rec.Body.Bytes(), &plat)
	if plat.Status.Phase != model.CloudPlatformPhaseActive {
		t.Fatalf("initial status=%q", plat.Status.Phase)
	}
	if h.audit.Len() < 1 {
		t.Fatal("required AuditEvent on success")
	}
}

// ---------------------------------------------------------------------------
// VS0-CF-F15-27 — server-side scope derivation
// ---------------------------------------------------------------------------

func TestVS0_CF_F15_27(t *testing.T) {
	h := newF15Harness(fullAdminGrants(f15ProviderUID, f15PlatformUID))
	seedPlatformProvider(t, h, f15PlatformUID, f15ProviderUID)
	hl := seedHostingLocation(t, h, "f15-27")
	if hl.Metadata.ScopeRef == nil || hl.Metadata.ScopeRef.UID != f15ProviderUID {
		t.Fatalf("HostingLocation scope from topology.write grant: %#v", hl.Metadata.ScopeRef)
	}
	part := seedPendingParticipation(t, h, "f15-27-part")
	if part.Metadata.ScopeRef.UID != f15PlatformUID {
		t.Fatalf("participation scope from CloudPlatform: %#v", part.Metadata.ScopeRef)
	}
	plat, _ := h.rt.Store.GetCloudPlatform(f15PlatformUID)
	if plat.Metadata.ScopeRef != nil {
		t.Fatalf("Platform-root scope must be nil, got %#v", plat.Metadata.ScopeRef)
	}
}

// ---------------------------------------------------------------------------
// VS0-CF-F15-28 — participation create fields vs empty action body
// ---------------------------------------------------------------------------

func TestVS0_CF_F15_28(t *testing.T) {
	h := newF15Harness(fullAdminGrants(f15ProviderUID, f15PlatformUID))
	seedPlatformProvider(t, h, f15PlatformUID, f15ProviderUID)
	// F0021 deferred field rejected
	body := []byte(`{
		"metadata":{"name":"f15-28"},
		"spec":{
			"cloudPlatformRef":{"apiVersion":"core.sovrunn.io/v1alpha1","kind":"CloudPlatform","name":"root-plat","uid":"` + f15PlatformUID + `"},
			"cloudProviderRef":{"apiVersion":"core.sovrunn.io/v1alpha1","kind":"CloudProvider","name":"prov","uid":"` + f15ProviderUID + `"},
			"environment":"development",
			"providerSelectionModes":["any"]
		}
	}`)
	rec := doF15(t, h, http.MethodPost, "/apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations", body, jsonHeaders("f15-28-def"))
	if rec.Code == http.StatusCreated {
		t.Fatal("F0021 fields must be rejected")
	}

	part := seedPendingParticipation(t, h, "f15-28-ok")
	// Non-empty ordinary body without bootstrapGrant → MALFORMED_REQUEST
	rec = doF15(t, h, http.MethodPost, actionPath(part.Metadata.UID, cloudmodel.ParticipationActionAccept),
		[]byte(`{"reason":"x"}`), actionHeaders(part.Metadata.ResourceVersion, "f15-28-body"))
	assertProblem(t, rec, apiproblem.CodeMalformedRequest, http.StatusBadRequest, "")
}

// ---------------------------------------------------------------------------
// VS0-CF-F15-29 — unassigned ISO code ZZ
// ---------------------------------------------------------------------------

func TestVS0_CF_F15_29(t *testing.T) {
	h := newF15Harness(platformOnlyGrants())
	body := []byte(`{
		"metadata":{"name":"f15-29"},
		"spec":{"ownerRegistration":{"legalName":"Acme","registrationIdentifier":"R1","jurisdictionCode":"ZZ"}}
	}`)
	before := h.audit.Len()
	rec := doF15(t, h, http.MethodPost, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms", body, jsonHeaders("f15-29"))
	assertProblem(t, rec, apiproblem.CodeValidationFailed, http.StatusUnprocessableEntity, "")
	assertNoAuditDelta(t, h, before)
}

// ---------------------------------------------------------------------------
// VS0-CF-F15-30 — malformed Idempotency-Key / If-Match on create / non-empty action body
// ---------------------------------------------------------------------------

func TestVS0_CF_F15_30(t *testing.T) {
	h := newF15Harness(platformOnlyGrants())
	before := h.audit.Len()
	rec := doF15(t, h, http.MethodPost, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms",
		platformBody("f15-30"), map[string]string{"Content-Type": "application/json"})
	assertProblem(t, rec, apiproblem.CodeMalformedRequest, http.StatusBadRequest, "")
	assertNoAuditDelta(t, h, before)

	h2 := newF15Harness(fullAdminGrants(f15ProviderUID, f15PlatformUID))
	seedPlatformProvider(t, h2, f15PlatformUID, f15ProviderUID)
	headers := jsonHeaders("f15-30-if")
	headers["If-Match"] = "1"
	before = h2.audit.Len()
	rec = doF15(t, h2, http.MethodPost, "/apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations",
		participationBody("f15-30", f15PlatformUID, f15ProviderUID), headers)
	assertProblem(t, rec, apiproblem.CodeMalformedRequest, http.StatusBadRequest, "")
	assertNoAuditDelta(t, h2, before)
}

// ---------------------------------------------------------------------------
// VS0-CF-F15-31 — AUTH_REQUIRED
// ---------------------------------------------------------------------------

func TestVS0_CF_F15_31(t *testing.T) {
	h := newF15Harness(platformOnlyGrants())
	before := h.audit.Len()
	rec := doF15Unauth(t, h, http.MethodGet, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms", nil, nil)
	assertProblem(t, rec, apiproblem.CodeAuthRequired, http.StatusUnauthorized, "")
	assertNoAuditDelta(t, h, before)

	rec = doF15Unauth(t, h, http.MethodPost, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms",
		platformBody("x"), jsonHeaders("f15-31"))
	assertProblem(t, rec, apiproblem.CodeAuthRequired, http.StatusUnauthorized, "")
	assertNoAuditDelta(t, h, before)
	if len(h.rt.Store.ListCloudPlatforms()) != 0 {
		t.Fatal("no mutation without auth")
	}
}

// ---------------------------------------------------------------------------
// VS0-CF-F15-32 — duplicate name ALREADY_EXISTS
// ---------------------------------------------------------------------------

func TestVS0_CF_F15_32(t *testing.T) {
	h := newF15Harness(platformOnlyGrants())
	rec := doF15(t, h, http.MethodPost, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms",
		platformBody("f15-32"), jsonHeaders("f15-32-a"))
	if rec.Code != http.StatusCreated {
		t.Fatalf("create: %d", rec.Code)
	}
	before := h.audit.Len()
	rec = doF15(t, h, http.MethodPost, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms",
		platformBody("f15-32"), jsonHeaders("f15-32-b"))
	assertProblem(t, rec, apiproblem.CodeAlreadyExists, http.StatusConflict, "")
	assertNoAuditDelta(t, h, before)
}

// ---------------------------------------------------------------------------
// VS0-CF-F15-33 — schema semantic validation failure
// ---------------------------------------------------------------------------

func TestVS0_CF_F15_33(t *testing.T) {
	h := newF15Harness(platformOnlyGrants())
	// Empty operatingMarkets
	seedPlatformProvider(t, h, f15PlatformUID, f15ProviderUID)
	h.grants = fullAdminGrants(f15ProviderUID, f15PlatformUID)
	h.rebind()
	before := h.audit.Len()
	rec := doF15(t, h, http.MethodPost, "/apis/core.sovrunn.io/v1alpha1/cloud-providers",
		[]byte(`{"metadata":{"name":"bad-markets"},"spec":{"operatingMarkets":[]}}`), jsonHeaders("f15-33"))
	assertProblem(t, rec, apiproblem.CodeValidationFailed, http.StatusUnprocessableEntity, "")
	assertNoAuditDelta(t, h, before)
}

// ---------------------------------------------------------------------------
// VS0-CF-F15-34 — HEAD matched to GET → 405, no Problem code
// ---------------------------------------------------------------------------

func TestVS0_CF_F15_34(t *testing.T) {
	h := newF15Harness(platformOnlyGrants())
	before := h.audit.Len()
	rec := doF15(t, h, http.MethodHead, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms", nil, nil)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("HEAD status=%d want 405", rec.Code)
	}
	if rec.Body.Len() != 0 {
		t.Fatalf("HEAD must not write body: %q", rec.Body.String())
	}
	assertNoAuditDelta(t, h, before)
	if len(server.CloudModelRoutePatterns()) != 35 {
		t.Fatal("exactly 35 registrations unchanged")
	}
}

// ---------------------------------------------------------------------------
// VS0-CF-F15-35 — forged header precedes media validation
// ---------------------------------------------------------------------------

func TestVS0_CF_F15_35(t *testing.T) {
	h := newF15Harness(fullAdminGrants(f15ProviderUID, f15PlatformUID))
	seedPlatformProvider(t, h, f15PlatformUID, f15ProviderUID)
	plat, _ := h.rt.Store.GetCloudPlatform(f15PlatformUID)
	before := h.audit.Len()
	headers := map[string]string{
		"Content-Type":                  "text/plain",
		"If-Match":                      plat.Metadata.ResourceVersion,
		cloudmodel.HeaderBootstrapGrant: "x",
	}
	rec := doF15(t, h, http.MethodPatch, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms/"+plat.Metadata.UID,
		[]byte(`not-json`), headers)
	assertProblem(t, rec, apiproblem.CodeAuthorizationDenied, http.StatusForbidden, "")
	assertAuditDelta(t, h, before, 1)

	// Without forged header, unsupported media remains 415
	rec = doF15(t, h, http.MethodPatch, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms/"+plat.Metadata.UID,
		[]byte(`{}`), map[string]string{"Content-Type": "text/plain", "If-Match": plat.Metadata.ResourceVersion})
	assertProblem(t, rec, apiproblem.CodeUnsupportedMediaType, http.StatusUnsupportedMediaType, "")
}

// ---------------------------------------------------------------------------
// VS0-CF-F15-36 — retention / eviction / cancel waiter
// ---------------------------------------------------------------------------

type f15Clock struct {
	mu  sync.Mutex
	now time.Time
}

func (c *f15Clock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.now
}

func (c *f15Clock) Set(t time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now = t.UTC()
}

func TestVS0_CF_F15_36(t *testing.T) {
	clock := &f15Clock{now: time.Date(2026, 8, 14, 0, 0, 0, 0, time.UTC)}
	c := cloudmodel.NewIdempotencyCoordinator(cloudmodel.IdempotencyCoordinatorConfig{Clock: clock})

	ns := cloudmodel.IdempotencyNamespace{PrincipalID: "p1", Pattern: "POST /apis/core.sovrunn.io/v1alpha1/cloud-platforms", Key: "k-24h"}
	d, err := cloudmodel.CanonicalDigest(cloudmodel.CanonicalDigestInput{Body: json.RawMessage(`{"metadata":{"name":"x"}}`), Pattern: ns.Pattern})
	if err != nil {
		t.Fatal(err)
	}
	_ = c.Reserve(context.Background(), ns, d, nil)
	c.Complete(ns, d, cloudmodel.CompletedResult{StatusCode: 201, Body: []byte(`{}`)})
	clock.Set(clock.Now().Add(cloudmodel.CompletedRetention + time.Second))
	if _, ok := c.StateForTest(ns); ok {
		t.Fatal("expired completed record unavailable for replay")
	}
	if res := c.Reserve(context.Background(), ns, d, nil); res.Outcome != cloudmodel.ReserveOutcomeReserved {
		t.Fatalf("next request processed as new: %#v", res)
	}

	// Capacity eviction
	clock.Set(time.Date(2026, 8, 14, 12, 0, 0, 0, time.UTC))
	c2 := cloudmodel.NewIdempotencyCoordinator(cloudmodel.IdempotencyCoordinatorConfig{Clock: clock})
	for i := 0; i < cloudmodel.MaxCompletedRecords+2; i++ {
		nsi := cloudmodel.IdempotencyNamespace{PrincipalID: "p1", Pattern: "POST /v1/cloud-platforms", Key: "ret-" + padF15(i)}
		di, _ := cloudmodel.CanonicalDigest(cloudmodel.CanonicalDigestInput{Body: json.RawMessage(`{"n":1}`), Pattern: nsi.Pattern})
		_ = c2.Reserve(context.Background(), nsi, di, nil)
		clock.Set(clock.Now().Add(time.Millisecond))
		c2.Complete(nsi, di, cloudmodel.CompletedResult{StatusCode: 201, Body: []byte(`{}`)})
	}
	if c2.CompletedCountForTest() != cloudmodel.MaxCompletedRecords {
		t.Fatalf("cap=%d", c2.CompletedCountForTest())
	}
	first := cloudmodel.IdempotencyNamespace{PrincipalID: "p1", Pattern: "POST /v1/cloud-platforms", Key: "ret-" + padF15(0)}
	if _, ok := c2.StateForTest(first); ok {
		t.Fatal("earliest completed must be evicted")
	}

	// Cancelled waiter detaches only itself; InFlight remains
	ns3 := cloudmodel.IdempotencyNamespace{PrincipalID: "p1", Pattern: "POST /v1/x", Key: "inflight"}
	d3, _ := cloudmodel.CanonicalDigest(cloudmodel.CanonicalDigestInput{Body: json.RawMessage(`{}`), Pattern: ns3.Pattern})
	_ = c2.Reserve(context.Background(), ns3, d3, nil)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan cloudmodel.ReserveResult, 1)
	go func() { done <- c2.Reserve(ctx, ns3, d3, nil) }()
	time.Sleep(20 * time.Millisecond)
	cancel()
	res := <-done
	if res.Outcome != cloudmodel.ReserveOutcomeDenied {
		t.Fatalf("cancelled waiter: %#v", res)
	}
	if st, ok := c2.StateForTest(ns3); !ok || st != cloudmodel.ReservationInFlight {
		t.Fatalf("InFlight must remain: %s ok=%v", st, ok)
	}
}

func padF15(i int) string {
	return fmt.Sprintf("%05d", i)
}

// ---------------------------------------------------------------------------
// VS0-CF-F15-37 — RFC 7396 null + LIST ascending uid
// ---------------------------------------------------------------------------

func TestVS0_CF_F15_37(t *testing.T) {
	h := newF15Harness(fullAdminGrants(f15ProviderUID, f15PlatformUID))
	seedPlatformProvider(t, h, f15PlatformUID, f15ProviderUID)
	plat, _ := h.rt.Store.GetCloudPlatform(f15PlatformUID)
	rec := doF15(t, h, http.MethodPatch, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms/"+plat.Metadata.UID,
		[]byte(`{"spec":{"description":"tmp"}}`), patchHeaders(plat.Metadata.ResourceVersion))
	_ = json.Unmarshal(rec.Body.Bytes(), &plat)
	rec = doF15(t, h, http.MethodPatch, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms/"+plat.Metadata.UID,
		[]byte(`{"spec":{"description":null}}`), patchHeaders(plat.Metadata.ResourceVersion))
	if rec.Code != http.StatusOK {
		t.Fatalf("null optional remove: %d %s", rec.Code, rec.Body.String())
	}
	var cleared model.CloudPlatform
	if err := json.Unmarshal(rec.Body.Bytes(), &cleared); err != nil {
		t.Fatalf("decode cleared: %v", err)
	}
	if cleared.Spec.Description != "" {
		t.Fatalf("description not removed: %q", cleared.Spec.Description)
	}

	// null required mutable field on provider
	prov, _ := h.rt.Store.GetCloudProvider(f15ProviderUID)
	rec = doF15(t, h, http.MethodPatch, "/apis/core.sovrunn.io/v1alpha1/cloud-providers/"+prov.Metadata.UID,
		[]byte(`{"spec":{"operatingMarkets":null}}`), patchHeaders(prov.Metadata.ResourceVersion))
	assertProblem(t, rec, apiproblem.CodeValidationFailed, http.StatusUnprocessableEntity, "")

	// Create second platform for LIST order (via store + platform grants)
	h2 := newF15Harness(platformOnlyGrants())
	_ = doF15(t, h2, http.MethodPost, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms", platformBody("z-plat"), jsonHeaders("z"))
	_ = doF15(t, h2, http.MethodPost, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms", platformBody("a-plat"), jsonHeaders("a"))
	rec = doF15(t, h2, http.MethodGet, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms", nil, nil)
	var list apimeta.ListEnvelope[model.CloudPlatform]
	_ = json.Unmarshal(rec.Body.Bytes(), &list)
	for i := 1; i < len(list.Items); i++ {
		if list.Items[i-1].Metadata.UID > list.Items[i].Metadata.UID {
			t.Fatal("LIST must be ascending metadata.uid")
		}
	}
}

// ---------------------------------------------------------------------------
// VS0-CF-F15-38 — PATCH If-Match missing/malformed/stale + lock recheck
// ---------------------------------------------------------------------------

func TestVS0_CF_F15_38(t *testing.T) {
	h := newF15Harness(fullAdminGrants(f15ProviderUID, f15PlatformUID))
	seedPlatformProvider(t, h, f15PlatformUID, f15ProviderUID)
	plat, _ := h.rt.Store.GetCloudPlatform(f15PlatformUID)
	before := h.audit.Len()
	rec := doF15(t, h, http.MethodPatch, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms/"+plat.Metadata.UID,
		[]byte(`{"spec":{"description":"x"}}`), map[string]string{"Content-Type": f15MediaMergePatch})
	assertProblem(t, rec, apiproblem.CodeStaleResourceVersion, http.StatusPreconditionFailed, "")
	assertNoAuditDelta(t, h, before)

	rec = doF15(t, h, http.MethodPatch, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms/"+plat.Metadata.UID,
		[]byte(`{"spec":{"description":"x"}}`), patchHeaders("stale"))
	assertProblem(t, rec, apiproblem.CodeStaleResourceVersion, http.StatusPreconditionFailed, "")
	assertNoAuditDelta(t, h, before)

	rec = doF15(t, h, http.MethodPatch, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms/"+plat.Metadata.UID,
		[]byte(`{"spec":{"description":"ok"}}`), patchHeaders(plat.Metadata.ResourceVersion))
	if rec.Code != http.StatusOK {
		t.Fatalf("matching If-Match publish: %d %s", rec.Code, rec.Body.String())
	}
}

// ---------------------------------------------------------------------------
// VS0-CF-F15-39 — GET/LIST exact read action/scope
// ---------------------------------------------------------------------------

func TestVS0_CF_F15_39(t *testing.T) {
	h := newF15Harness(fullAdminGrants(f15ProviderUID, f15PlatformUID))
	seedPlatformProvider(t, h, f15PlatformUID, f15ProviderUID)
	plat, _ := h.rt.Store.GetCloudPlatform(f15PlatformUID)

	rec := doF15(t, h, http.MethodGet, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms/"+plat.Metadata.UID, nil, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("authorized GET: %d", rec.Code)
	}
	rec = doF15(t, h, http.MethodGet, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms", nil, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("authorized LIST: %d", rec.Code)
	}

	h.grants = &f15Grants{principal: f15Principal, grants: []cloudmodel.Grant{
		{Action: cloudmodel.ActionCloudPlatformWrite, Scope: cloudmodel.PlatformRootScope},
	}}
	h.rebind()
	rec = doF15(t, h, http.MethodGet, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms", nil, nil)
	assertProblem(t, rec, apiproblem.CodeAuthorizationDenied, http.StatusForbidden, "")
	rec = doF15(t, h, http.MethodGet, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms/"+plat.Metadata.UID, nil, nil)
	assertProblem(t, rec, apiproblem.CodeResourceNotFound, http.StatusNotFound, f15ViolationSafeDeny)
}

// ---------------------------------------------------------------------------
// VS0-CF-F15-40 — replay headers / Aborted re-reservation
// ---------------------------------------------------------------------------

func TestVS0_CF_F15_40(t *testing.T) {
	h := newF15Harness(platformOnlyGrants())
	headers := jsonHeaders("f15-40")
	rec := doF15(t, h, http.MethodPost, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms", platformBody("f15-40"), headers)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create: %d", rec.Code)
	}
	loc := rec.Header().Get("Location")

	rec = doF15(t, h, http.MethodPost, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms", platformBody("f15-40"), headers)
	if rec.Code != http.StatusCreated {
		t.Fatalf("replay: %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
		t.Fatalf("replay Content-Type=%q", ct)
	}
	// Location must not be replayed as a stored transport header.
	if loc != "" && rec.Header().Get("Location") == loc {
		t.Fatal("Location must not be replayed from stored completion")
	}

	// Aborted namespace re-reserves as new
	ns := cloudmodel.IdempotencyNamespace{PrincipalID: f15Principal, Pattern: apiPatternCloudPlatformCreate(), Key: "abort-key"}
	d, _ := cloudmodel.CanonicalDigest(cloudmodel.CanonicalDigestInput{Body: json.RawMessage(platformBody("x")), Pattern: ns.Pattern})
	_ = h.rt.Idempotency.Reserve(context.Background(), ns, d, nil)
	h.rt.Idempotency.Abort(ns)
	if res := h.rt.Idempotency.Reserve(context.Background(), ns, d, nil); res.Outcome != cloudmodel.ReserveOutcomeReserved {
		t.Fatalf("after Aborted removal: %#v", res)
	}
}

func apiPatternCloudPlatformCreate() string {
	return "POST /apis/core.sovrunn.io/v1alpha1/cloud-platforms"
}

// ---------------------------------------------------------------------------
// VS0-CF-F15-41 — writer boundary denial without exact admin writer
// ---------------------------------------------------------------------------

func TestVS0_CF_F15_41(t *testing.T) {
	h := newF15Harness(fullAdminGrants(f15ProviderUID, f15PlatformUID))
	seedPlatformProvider(t, h, f15PlatformUID, f15ProviderUID)
	plat, _ := h.rt.Store.GetCloudPlatform(f15PlatformUID)
	hl := seedHostingLocation(t, h, "f15-41-hl")

	// Principal without cloudplatform.write
	h.grants = &f15Grants{principal: f15Principal, grants: []cloudmodel.Grant{
		{Action: cloudmodel.ActionCloudPlatformRead, Scope: cloudmodel.PlatformRootScope},
		{Action: cloudmodel.ActionCloudProviderRead, Scope: cloudmodel.PlatformRootScope},
		{Action: cloudmodel.ActionCloudProviderWrite, Scope: cloudmodel.PlatformRootScope},
		{Action: cloudmodel.ActionTopologyWrite, Scope: providerScope(f15ProviderUID)},
		{Action: cloudmodel.ActionTopologyRead, Scope: providerScope(f15ProviderUID)},
	}}
	h.rebind()
	before := h.audit.Len()
	rec := doF15(t, h, http.MethodPatch, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms/"+plat.Metadata.UID,
		[]byte(`{"spec":{"description":"nope"}}`), patchHeaders(plat.Metadata.ResourceVersion))
	assertProblem(t, rec, apiproblem.CodeAuthorizationDenied, http.StatusForbidden, "")
	assertNoAuditDelta(t, h, before)
	got, _ := h.rt.Store.GetCloudPlatform(f15PlatformUID)
	if got.Spec.Description == "nope" {
		t.Fatal("no resource mutation")
	}

	// Without topology.write for topology PATCH
	h.grants = &f15Grants{principal: f15Principal, grants: []cloudmodel.Grant{
		{Action: cloudmodel.ActionCloudPlatformRead, Scope: cloudmodel.PlatformRootScope},
		{Action: cloudmodel.ActionCloudProviderRead, Scope: cloudmodel.PlatformRootScope},
		{Action: cloudmodel.ActionTopologyRead, Scope: providerScope(f15ProviderUID)},
	}}
	h.rebind()
	before = h.audit.Len()
	rec = doF15(t, h, http.MethodPatch, "/apis/infrastructure.sovrunn.io/v1alpha1/hosting-locations/"+hl.Metadata.UID,
		[]byte(`{"spec":{"description":"nope"}}`), patchHeaders(hl.Metadata.ResourceVersion))
	assertProblem(t, rec, apiproblem.CodeAuthorizationDenied, http.StatusForbidden, "")
	assertNoAuditDelta(t, h, before)
}
