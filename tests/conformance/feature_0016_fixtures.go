// Package conformance — FEATURE-0016 fixtures (TASK-F16-12).
package conformance

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
	"unsafe"

	"github.com/sanjeevksaini/sovrunn/internal/api"
	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/model"
	"github.com/sanjeevksaini/sovrunn/internal/executiontarget"
	etmodel "github.com/sanjeevksaini/sovrunn/internal/executiontarget/model"
	"github.com/sanjeevksaini/sovrunn/internal/server"
)

const (
	f16Principal         = "bootstrap-principal"
	f16PlatformUID       = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	f16ProviderUID       = "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	f16OtherProviderUID  = "cccccccccccccccccccccccccccccccc"
	f16PartUID           = "11111111111111111111111111111111"
	f16StackUID          = "22222222222222222222222222222222"
	f16FaultDomainUID    = "33333333333333333333333333333333"
	f16Collection        = "/apis/execution.sovrunn.io/v1alpha1/execution-targets"
	f16ItemBase          = f16Collection + "/"
	f16MediaJSON         = "application/json"
	f16ViolationSafeDeny = apiproblem.ViolationCode("VS0_AUTHORIZATION_SAFE_DENIAL")
)

// ---------------------------------------------------------------------------
// Grants
// ---------------------------------------------------------------------------

type f16Grants struct {
	principal string
	grants    []cloudmodel.Grant
}

func (g *f16Grants) PrincipalID() string { return g.principal }

func (g *f16Grants) CoarseLookup(action string) []cloudmodel.Grant {
	var out []cloudmodel.Grant
	for _, gr := range g.grants {
		if gr.Action == action {
			out = append(out, gr)
		}
	}
	return out
}

func (g *f16Grants) AuthorizeExact(action string, scope apimeta.ScopeIdentity, resourceUID, relatedPlatformUID string) bool {
	_ = relatedPlatformUID
	for _, gr := range g.grants {
		if gr.Action != action || gr.Scope != scope {
			continue
		}
		if gr.ResourceUID != "" && gr.ResourceUID != resourceUID {
			continue
		}
		return true
	}
	return false
}

func f16ProviderScope(uid string) apimeta.ScopeIdentity {
	return apimeta.ScopeIdentity{Kind: apimeta.ScopeCloudProvider, UID: uid}
}

func f16ExecutionTargetGrants(providerUID string, write, read, qualify bool) *f16Grants {
	g := &f16Grants{principal: f16Principal}
	scope := f16ProviderScope(providerUID)
	if write {
		g.grants = append(g.grants, cloudmodel.Grant{Action: api.ActionExecutionTargetWrite, Scope: scope})
	}
	if read {
		g.grants = append(g.grants, cloudmodel.Grant{Action: api.ActionExecutionTargetRead, Scope: scope})
	}
	if qualify {
		g.grants = append(g.grants, cloudmodel.Grant{Action: api.ActionExecutionTargetQualify, Scope: scope})
	}
	return g
}

func f16FullGrants() *f16Grants {
	return f16ExecutionTargetGrants(f16ProviderUID, true, true, true)
}

// ---------------------------------------------------------------------------
// Blocking fixture (mirrors blockingActionFixture)
// ---------------------------------------------------------------------------

type f16BlockingFixture struct {
	inner   executiontarget.FixtureSource
	started chan struct{}
	release chan struct{}
}

func (b *f16BlockingFixture) Lookup(targetUID string) (executiontarget.TargetBoundFixture, bool) {
	var (
		fx executiontarget.TargetBoundFixture
		ok bool
	)
	if b.inner != nil {
		fx, ok = b.inner.Lookup(targetUID)
	}
	select {
	case <-b.started:
	default:
		close(b.started)
	}
	<-b.release
	return fx, ok
}

// f16ArmedBlockingFixture blocks Lookup only after Arm is called, so a prior
// non-blocking qualify can establish Active/Qualified before mid-flight proofs.
type f16ArmedBlockingFixture struct {
	inner   executiontarget.FixtureSource
	armed   atomic.Bool
	started chan struct{}
	release chan struct{}
}

func (b *f16ArmedBlockingFixture) Arm() { b.armed.Store(true) }

func (b *f16ArmedBlockingFixture) Lookup(targetUID string) (executiontarget.TargetBoundFixture, bool) {
	var (
		fx executiontarget.TargetBoundFixture
		ok bool
	)
	if b.inner != nil {
		fx, ok = b.inner.Lookup(targetUID)
	}
	if !b.armed.Load() {
		return fx, ok
	}
	select {
	case <-b.started:
	default:
		close(b.started)
	}
	<-b.release
	return fx, ok
}

// ---------------------------------------------------------------------------
// Harness
// ---------------------------------------------------------------------------

type f16Harness struct {
	cloud      *cloudmodel.Store
	audit      *cloudmodel.MemoryAuditAppender
	clock      *executiontarget.FakeClock
	fixtures   *executiontarget.MapFixtureSource
	life       *executiontarget.ExecutionTargetLifecycleService
	sched      *executiontarget.Scheduler
	grants     *f16Grants
	handler    http.Handler
	qualify    *api.ExecutionTargetActionHandler
	retire     *api.ExecutionTargetActionHandler
	item       *api.ExecutionTargetItemHandler
	collection *api.ExecutionTargetCollectionHandler
}

type f16HarnessOpts struct {
	grants      *f16Grants
	observer    *executiontarget.SyntheticObserver
	fixtures    *executiontarget.MapFixtureSource
	audit       executiontarget.AuditAppender
	memoryAudit *cloudmodel.MemoryAuditAppender
	startClock  time.Time
}

func newF16Harness(t *testing.T, opts f16HarnessOpts) *f16Harness {
	t.Helper()
	if opts.grants == nil {
		opts.grants = f16FullGrants()
	}
	clockStart := opts.startClock
	if clockStart.IsZero() {
		clockStart = time.Date(2026, 8, 21, 10, 0, 0, 0, time.UTC)
	}
	clock := executiontarget.NewFakeClock(clockStart)

	var memAudit *cloudmodel.MemoryAuditAppender
	var audit executiontarget.AuditAppender
	if opts.audit != nil {
		audit = opts.audit
		memAudit = opts.memoryAudit
	} else {
		memAudit = &cloudmodel.MemoryAuditAppender{}
		audit = memAudit
	}

	fixtures := opts.fixtures
	if fixtures == nil {
		fixtures = executiontarget.NewMapFixtureSource()
	}
	observer := opts.observer
	if observer == nil {
		observer = executiontarget.NewSyntheticObserver(fixtures, func() time.Time { return clock.Now().UTC() })
	}

	cloud := cloudmodel.NewStore()
	life := executiontarget.NewExecutionTargetLifecycleService(executiontarget.LifecycleConfig{
		Audit:    audit,
		Observer: observer,
		Clock:    clock,
	})
	sched := executiontarget.NewScheduler(executiontarget.SchedulerConfig{
		Lifecycle: life,
		Clock:     clock,
	})

	collection := api.NewExecutionTargetCollectionHandler(life, cloud, opts.grants, audit)
	item := api.NewExecutionTargetItemHandler(life, cloud, opts.grants, audit)
	qualify := api.NewExecutionTargetQualifyHandler(life, cloud, opts.grants, audit)
	qualify.Fixtures = fixtures
	retire := api.NewExecutionTargetRetireHandler(life, cloud, opts.grants, audit)

	mux := server.NewExecutionTargetMux(&server.ExecutionTargetHandlers{
		Collection: collection,
		Item:       item,
		Qualify:    qualify,
		Retire:     retire,
	})
	guarded := server.ExecutionTargetTransportGuard(mux)

	return &f16Harness{
		cloud:      cloud,
		audit:      memAudit,
		clock:      clock,
		fixtures:   fixtures,
		life:       life,
		sched:      sched,
		grants:     opts.grants,
		handler:    guarded,
		qualify:    qualify,
		retire:     retire,
		item:       item,
		collection: collection,
	}
}

func (h *f16Harness) rebindGrants(grants *f16Grants) {
	h.grants = grants
	audit := executiontarget.AuditAppender(h.audit)
	h.collection = api.NewExecutionTargetCollectionHandler(h.life, h.cloud, grants, audit)
	h.item = api.NewExecutionTargetItemHandler(h.life, h.cloud, grants, audit)
	h.qualify = api.NewExecutionTargetQualifyHandler(h.life, h.cloud, grants, audit)
	h.qualify.Fixtures = h.fixtures
	h.retire = api.NewExecutionTargetRetireHandler(h.life, h.cloud, grants, audit)
	mux := server.NewExecutionTargetMux(&server.ExecutionTargetHandlers{
		Collection: h.collection,
		Item:       h.item,
		Qualify:    h.qualify,
		Retire:     h.retire,
	})
	h.handler = server.ExecutionTargetTransportGuard(mux)
}

func (h *f16Harness) auditLen() int {
	if h.audit == nil {
		return 0
	}
	return h.audit.Len()
}

func (h *f16Harness) setAuditFail(err error) {
	if h.audit != nil {
		h.audit.SetFail(err)
	}
}

// ---------------------------------------------------------------------------
// Seed helpers
// ---------------------------------------------------------------------------

func f16ProviderScopeRef(uid, name string) *apimeta.ScopeRef {
	return &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
		APIVersion: model.APIVersionCloudProvider,
		Kind:       string(apimeta.ScopeCloudProvider),
		Name:       name,
		UID:        uid,
	}}
}

func (h *f16Harness) seedPlatformAndProvider(t *testing.T, providerUID string) {
	t.Helper()
	if !h.cloud.HasCloudPlatform() {
		if _, prob := h.cloud.CreateCloudPlatform(model.CloudPlatform{
			TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionCloudPlatform, Kind: model.KindCloudPlatform},
			Metadata: apimeta.ObjectMeta{Name: "root-plat", UID: f16PlatformUID},
			Spec: model.CloudPlatformSpec{OwnerRegistration: model.OwnerRegistration{
				LegalName: "Acme", RegistrationIdentifier: "R1", JurisdictionCode: "US",
			}},
			Status: model.CloudPlatformStatus{Phase: model.CloudPlatformPhaseActive},
		}); prob != nil {
			t.Fatalf("seed platform: %#v", prob)
		}
	}
	if _, ok := h.cloud.GetCloudProvider(providerUID); ok {
		return
	}
	name := "prov-a"
	if providerUID == f16OtherProviderUID {
		name = "prov-b"
	}
	if _, prob := h.cloud.CreateCloudProvider(model.CloudProvider{
		TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionCloudProvider, Kind: model.KindCloudProvider},
		Metadata: apimeta.ObjectMeta{Name: name, UID: providerUID, ResourceVersion: "1"},
		Spec:     model.CloudProviderSpec{OperatingMarkets: []string{"US"}},
		Status:   model.CloudProviderStatus{Phase: model.CloudProviderPhaseActive},
	}); prob != nil {
		t.Fatalf("seed provider: %#v", prob)
	}
}

func (h *f16Harness) seedActiveParticipationAndStack(t *testing.T) (model.CloudProviderParticipation, model.InfrastructureStack) {
	t.Helper()
	h.seedPlatformAndProvider(t, f16ProviderUID)

	part, prob := h.cloud.CreateParticipation(model.CloudProviderParticipation{
		TypeMeta: apimeta.TypeMeta{
			APIVersion: model.APIVersionCloudProviderParticipation,
			Kind:       model.KindCloudProviderParticipation,
		},
		Metadata: apimeta.ObjectMeta{
			Name: "part-active", UID: f16PartUID, ResourceVersion: "1", Generation: 1,
			ScopeRef: &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
				APIVersion: model.APIVersionCloudPlatform, Kind: string(apimeta.ScopeCloudPlatform),
				Name: "root-plat", UID: f16PlatformUID,
			}},
		},
		Spec: model.CloudProviderParticipationSpec{
			CloudPlatformRef: apimeta.TypedRef{
				APIVersion: model.APIVersionCloudPlatform, Kind: model.KindCloudPlatform,
				Name: "root-plat", UID: f16PlatformUID,
			},
			CloudProviderRef: apimeta.TypedRef{
				APIVersion: model.APIVersionCloudProvider, Kind: model.KindCloudProvider,
				Name: "prov-a", UID: f16ProviderUID,
			},
			Environment: model.ParticipationEnvironmentDevelopment,
		},
		Status: model.CloudProviderParticipationStatus{Phase: model.ParticipationPhaseActive},
	})
	if prob != nil {
		t.Fatalf("seed participation: %#v", prob)
	}

	stack, prob := h.cloud.CreateInfrastructureStack(model.InfrastructureStack{
		TypeMeta: apimeta.TypeMeta{
			APIVersion: model.APIVersionInfrastructureStack,
			Kind:       model.KindInfrastructureStack,
		},
		Metadata: apimeta.ObjectMeta{
			Name: "stack-active", UID: f16StackUID, ResourceVersion: "1", Generation: 3,
			ScopeRef: f16ProviderScopeRef(f16ProviderUID, "prov-a"),
		},
		Spec: model.InfrastructureStackSpec{
			FaultDomainRef: apimeta.TypedRef{
				APIVersion: model.APIVersionFaultDomain, Kind: model.KindFaultDomain,
				Name: "fd", UID: f16FaultDomainUID,
			},
		},
		Status: model.InfrastructureStackStatus{Phase: model.InfrastructureStackPhaseActive},
	})
	if prob != nil {
		t.Fatalf("seed stack: %#v", prob)
	}
	return part, stack
}

func (h *f16Harness) seedOtherProviderAndStack(t *testing.T, stackUID string) model.InfrastructureStack {
	t.Helper()
	h.seedPlatformAndProvider(t, f16OtherProviderUID)
	stack, prob := h.cloud.CreateInfrastructureStack(model.InfrastructureStack{
		TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionInfrastructureStack, Kind: model.KindInfrastructureStack},
		Metadata: apimeta.ObjectMeta{
			Name: "other-stack", UID: stackUID, ResourceVersion: "1", Generation: 1,
			ScopeRef: f16ProviderScopeRef(f16OtherProviderUID, "prov-b"),
		},
		Status: model.InfrastructureStackStatus{Phase: model.InfrastructureStackPhaseActive},
	})
	if prob != nil {
		t.Fatalf("seed other stack: %#v", prob)
	}
	return stack
}

func (h *f16Harness) seedExtraStack(t *testing.T, name, uid string, generation int64) model.InfrastructureStack {
	t.Helper()
	stack, prob := h.cloud.CreateInfrastructureStack(model.InfrastructureStack{
		TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionInfrastructureStack, Kind: model.KindInfrastructureStack},
		Metadata: apimeta.ObjectMeta{
			Name: name, UID: uid, ResourceVersion: "1", Generation: generation,
			ScopeRef: f16ProviderScopeRef(f16ProviderUID, "prov-a"),
		},
		Status: model.InfrastructureStackStatus{Phase: model.InfrastructureStackPhaseActive},
	})
	if prob != nil {
		t.Fatalf("seed stack %s: %#v", name, prob)
	}
	return stack
}

func f16ActorRef(principalID string) apimeta.TypedRef {
	return apimeta.TypedRef{
		APIVersion: "identity.sovrunn.io/v1alpha1",
		Kind:       "Principal",
		Name:       principalID,
		UID:        principalID,
	}
}

// ---------------------------------------------------------------------------
// HTTP helpers
// ---------------------------------------------------------------------------

func f16CreateBody(name, partUID, stackUID string) []byte {
	return []byte(`{
		"metadata":{"name":"` + name + `"},
		"spec":{
			"cloudProviderParticipationRef":{"uid":"` + partUID + `"},
			"infrastructureStackRef":{"uid":"` + stackUID + `"},
			"targetClass":"synthetic-iaas"
		}
	}`)
}

func f16CreateHeaders(key string) map[string]string {
	return map[string]string{
		"Content-Type":    f16MediaJSON,
		"Idempotency-Key": key,
		"Accept":          "application/xml",
	}
}

func f16ActionHeaders(key, ifMatch string) map[string]string {
	return map[string]string{
		"Content-Type":    f16MediaJSON,
		"Idempotency-Key": key,
		"If-Match":        ifMatch,
		"Accept":          "application/xml",
	}
}

func f16ItemPath(uid string) string {
	return f16ItemBase + uid
}

func f16ActionPath(uid, action string) string {
	return f16ItemBase + uid + "/actions/" + action
}

func doF16(t *testing.T, h *f16Harness, method, path string, body []byte, headers map[string]string) *httptest.ResponseRecorder {
	t.Helper()
	var r *http.Request
	if body == nil {
		r = httptest.NewRequest(method, path, nil)
	} else {
		r = httptest.NewRequest(method, path, bytes.NewReader(body))
	}
	r.Header.Set("Authorization", "Bearer "+f16Principal)
	r.Header.Set("X-Sovrunn-Request-ID", "req-f16-conformance")
	for k, v := range headers {
		if k == "If-Match" && v == "" {
			continue
		}
		r.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()
	h.handler.ServeHTTP(rec, r)
	return rec
}

func doF16Unauth(t *testing.T, h *f16Harness, method, path string, body []byte, headers map[string]string) *httptest.ResponseRecorder {
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
	h.handler.ServeHTTP(rec, r)
	return rec
}

func doF16InvalidAuth(t *testing.T, h *f16Harness, method, path string, body []byte, headers map[string]string) *httptest.ResponseRecorder {
	t.Helper()
	var r *http.Request
	if body == nil {
		r = httptest.NewRequest(method, path, nil)
	} else {
		r = httptest.NewRequest(method, path, bytes.NewReader(body))
	}
	r.Header.Set("Authorization", "Bearer not-the-principal")
	r.Header.Set("X-Sovrunn-Request-ID", "req-f16-invalid-auth")
	for k, v := range headers {
		r.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()
	h.handler.ServeHTTP(rec, r)
	return rec
}

func doF16MultiIfMatch(t *testing.T, h *f16Harness, path, key string, values ...string) *httptest.ResponseRecorder {
	t.Helper()
	r := httptest.NewRequest(http.MethodPost, path, nil)
	r.Header.Set("Authorization", "Bearer "+f16Principal)
	r.Header.Set("Content-Type", f16MediaJSON)
	r.Header.Set("Idempotency-Key", key)
	for _, v := range values {
		r.Header.Add("If-Match", v)
	}
	rec := httptest.NewRecorder()
	h.handler.ServeHTTP(rec, r)
	return rec
}

func decodeF16Problem(t *testing.T, rec *httptest.ResponseRecorder) *apiproblem.Problem {
	t.Helper()
	var p apiproblem.Problem
	if err := json.Unmarshal(rec.Body.Bytes(), &p); err != nil {
		t.Fatalf("problem decode: %v body=%s", err, rec.Body.String())
	}
	return &p
}

func hasF16Violation(prob *apiproblem.Problem, code apiproblem.ViolationCode) bool {
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

func assertF16JSONNosniff(t *testing.T, rec *httptest.ResponseRecorder) {
	t.Helper()
	if got := rec.Header().Get("Content-Type"); got != f16MediaJSON {
		t.Fatalf("Content-Type=%q, want %q", got, f16MediaJSON)
	}
	if got := rec.Header().Get("X-Content-Type-Options"); got != "nosniff" {
		t.Fatalf("X-Content-Type-Options=%q, want nosniff", got)
	}
}

func assertF16ProblemNosniff(t *testing.T, rec *httptest.ResponseRecorder) {
	t.Helper()
	if got := rec.Header().Get("Content-Type"); got != apiproblem.MediaTypeProblemJSON {
		t.Fatalf("Content-Type=%q, want problem+json", got)
	}
	if got := rec.Header().Get("X-Content-Type-Options"); got != "nosniff" {
		t.Fatalf("X-Content-Type-Options=%q, want nosniff", got)
	}
}

func assertF16TransportOnly(t *testing.T, rec *httptest.ResponseRecorder, wantStatus int, wantAllow string) {
	t.Helper()
	if rec.Code != wantStatus {
		t.Fatalf("status=%d, want %d body=%s", rec.Code, wantStatus, rec.Body.String())
	}
	if rec.Body.Len() != 0 {
		t.Fatalf("transport response must have empty body, got %q", rec.Body.String())
	}
	if got := rec.Header().Get("X-Content-Type-Options"); got != "nosniff" {
		t.Fatalf("X-Content-Type-Options=%q, want nosniff", got)
	}
	if wantAllow != "" {
		if got := rec.Header().Get("Allow"); got != wantAllow {
			t.Fatalf("Allow=%q, want %q", got, wantAllow)
		}
	}
	if ct := rec.Header().Get("Content-Type"); ct != "" {
		t.Fatalf("transport response must not set Content-Type, got %q", ct)
	}
}

func assertF16Problem(t *testing.T, rec *httptest.ResponseRecorder, wantCode apiproblem.ErrorCode, wantStatus int, wantViolation apiproblem.ViolationCode) *apiproblem.Problem {
	t.Helper()
	if rec.Code != wantStatus {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	assertF16ProblemNosniff(t, rec)
	p := decodeF16Problem(t, rec)
	if p.Code != wantCode {
		t.Fatalf("code=%s, want %s", p.Code, wantCode)
	}
	if wantViolation != "" && !hasF16Violation(p, wantViolation) {
		t.Fatalf("missing violation %s in %#v", wantViolation, p.Violations)
	}
	return p
}

func assertF16NoAuditDelta(t *testing.T, h *f16Harness, before int) {
	t.Helper()
	if h.auditLen() != before {
		t.Fatalf("audit len=%d, want unchanged %d", h.auditLen(), before)
	}
}

func assertF16AuditDelta(t *testing.T, h *f16Harness, before, delta int) {
	t.Helper()
	want := before + delta
	if h.auditLen() != want {
		t.Fatalf("audit len=%d, want %d", h.auditLen(), want)
	}
}

func assertF16SafeProjectionOmitsInternals(t *testing.T, body string) {
	t.Helper()
	for _, forbidden := range []string{
		"factSetRef", "qualificationResultRef", "observerID", "NormalizedTargetFactSet",
		"TargetQualificationResult", "handle", `"facts"`, "conditions",
	} {
		if bytes.Contains([]byte(body), []byte(forbidden)) {
			t.Fatalf("safe projection must omit %q; body=%s", forbidden, body)
		}
	}
}

// ---------------------------------------------------------------------------
// Create / qualify helpers
// ---------------------------------------------------------------------------

func (h *f16Harness) createTarget(t *testing.T, name, key string) executiontarget.SafeExecutionTarget {
	t.Helper()
	h.seedActiveParticipationAndStack(t)
	rec := doF16(t, h, http.MethodPost, f16Collection, f16CreateBody(name, f16PartUID, f16StackUID), f16CreateHeaders(key))
	if rec.Code != http.StatusCreated {
		t.Fatalf("seed create status=%d body=%s", rec.Code, rec.Body.String())
	}
	var created executiontarget.SafeExecutionTarget
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatalf("decode create: %v", err)
	}
	return created
}

func (h *f16Harness) createTargetOnStack(t *testing.T, name, key, stackUID string) executiontarget.SafeExecutionTarget {
	t.Helper()
	rec := doF16(t, h, http.MethodPost, f16Collection, f16CreateBody(name, f16PartUID, stackUID), f16CreateHeaders(key))
	if rec.Code != http.StatusCreated {
		t.Fatalf("create status=%d body=%s", rec.Code, rec.Body.String())
	}
	var created executiontarget.SafeExecutionTarget
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatalf("decode create: %v", err)
	}
	return created
}

func (h *f16Harness) setPresentFixture(uid string, truths etmodel.FactSetTruths) {
	h.fixtures.Set(uid, executiontarget.TargetBoundFixture{
		Revision: "fx-1",
		Mode:     executiontarget.FixturePresent,
		Truths:   truths,
	})
}

func f16AllSupportedTruths() etmodel.FactSetTruths {
	return etmodel.FactSetTruths{
		Compute: etmodel.FactComputeTruths{VM: etmodel.FactSupported},
		Storage: etmodel.FactStorageTruths{Block: etmodel.FactSupported, Object: etmodel.FactSupported},
		Network: etmodel.FactNetworkTruths{Private: etmodel.FactSupported},
	}
}

func f16UnsupportedVMTruths() etmodel.FactSetTruths {
	t := f16AllSupportedTruths()
	t.Compute.VM = etmodel.FactUnsupported
	return t
}

func (h *f16Harness) maintenanceTrigger(et etmodel.ExecutionTarget, op executiontarget.MaintenanceTriggerOperation, auditUID string) executiontarget.MaintenanceTrigger {
	return executiontarget.MaintenanceTrigger{
		TargetUID:                      et.Metadata.UID,
		ExpectedInfrastructureStackGen: et.Status.ObservedGeneration,
		ExpectedMaintenanceEpoch:       et.Status.MaintenanceEpoch,
		Operation:                      op,
		Actor:                          executiontarget.MaintenanceSystemActor(),
		RequestID:                      "corr-" + auditUID,
		AuditUID:                       auditUID,
	}
}

func (h *f16Harness) deliverMaintenance(t *testing.T, et etmodel.ExecutionTarget, op executiontarget.MaintenanceTriggerOperation, auditUID string) etmodel.ExecutionTarget {
	t.Helper()
	got, prob := h.sched.DeliverMaintenanceTrigger(context.Background(), h.maintenanceTrigger(et, op, auditUID))
	if prob != nil {
		t.Fatalf("maintenance %s: %#v", op, prob)
	}
	return got
}

// f16StoreShadow mirrors executiontarget.Store layout for F16-75 mid-flight
// MaintenanceEpoch bumps without an active current-Maintenance marker.
type f16StoreShadow struct {
	mu                 sync.Mutex
	targets            map[string]etmodel.ExecutionTarget
	factSets           map[string]etmodel.NormalizedTargetFactSet
	results            map[string]etmodel.TargetQualificationResult
	maintenanceMarkers map[string]etmodel.CurrentMaintenanceMarker
	nameIndex          map[string]string
	tupleIndex         map[string]string
}

func (h *f16Harness) bumpMaintenanceEpochNoMarker(uid string, epoch int64) {
	// #nosec G103 -- test-only white-box fixture required to model the registered stale-epoch/no-marker case.
	shadow := (*f16StoreShadow)(unsafe.Pointer(h.life.Store()))
	shadow.mu.Lock()
	defer shadow.mu.Unlock()
	cur, ok := shadow.targets[uid]
	if !ok {
		return
	}
	cur.Status.MaintenanceEpoch = epoch
	if cur.Metadata.ResourceVersion == "" {
		cur.Metadata.ResourceVersion = "1"
	} else if n, err := strconv.ParseUint(cur.Metadata.ResourceVersion, 10, 64); err == nil {
		cur.Metadata.ResourceVersion = strconv.FormatUint(n+1, 10)
	} else {
		cur.Metadata.ResourceVersion = "1"
	}
	shadow.targets[uid] = cur
}

func (h *f16Harness) updateStackGeneration(t *testing.T, stackUID string, generation int64) {
	t.Helper()
	stack, ok := h.cloud.GetInfrastructureStack(stackUID)
	if !ok {
		t.Fatalf("stack %s missing", stackUID)
	}
	stack.Metadata.Generation = generation
	h.cloud.BeginPublication()
	staged, _ := h.cloud.StageUpdateInfrastructureStack(stack)
	h.cloud.Publish(staged)
	h.cloud.EndPublication()
}

func (h *f16Harness) setParticipationEffective(t *testing.T, partUID string, active bool) {
	t.Helper()
	part, ok := h.cloud.GetParticipation(partUID)
	if !ok {
		t.Fatalf("participation %s missing", partUID)
	}
	if active {
		part.Status.PlatformSuspended = false
		part.Status.Phase = model.ParticipationPhaseActive
	} else {
		part.Status.PlatformSuspended = true
		part.Status.Phase = model.ParticipationPhaseSuspended
	}
	h.cloud.BeginPublication()
	staged, _ := h.cloud.StageUpdateParticipation(part)
	h.cloud.Publish(staged)
	h.cloud.EndPublication()
}

func (h *f16Harness) setStackPhase(t *testing.T, stackUID string, phase model.InfrastructureStackPhase) {
	t.Helper()
	stack, ok := h.cloud.GetInfrastructureStack(stackUID)
	if !ok {
		t.Fatalf("stack %s missing", stackUID)
	}
	stack.Status.Phase = phase
	h.cloud.BeginPublication()
	staged, _ := h.cloud.StageUpdateInfrastructureStack(stack)
	h.cloud.Publish(staged)
	h.cloud.EndPublication()
}
