// Package conformance -- FEATURE-0016 local conformance (TASK-F16-12).
//
// One runtime scenario per FEATURE-0016-owned registry ID:
// VS0-CF-F16-01 through VS0-CF-F16-122.
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
	"sync/atomic"
	"testing"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/api"
	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apivalid"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/model"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/validate"
	"github.com/sanjeevksaini/sovrunn/internal/executiontarget"
	etmodel "github.com/sanjeevksaini/sovrunn/internal/executiontarget/model"
)

// f16CatalogIDs lists every approved VS0-CF-F16 case ID (01..122).
var f16CatalogIDs = []string{
	"VS0-CF-F16-01",
	"VS0-CF-F16-02",
	"VS0-CF-F16-03",
	"VS0-CF-F16-04",
	"VS0-CF-F16-05",
	"VS0-CF-F16-06",
	"VS0-CF-F16-07",
	"VS0-CF-F16-08",
	"VS0-CF-F16-09",
	"VS0-CF-F16-10",
	"VS0-CF-F16-11",
	"VS0-CF-F16-12",
	"VS0-CF-F16-13",
	"VS0-CF-F16-14",
	"VS0-CF-F16-15",
	"VS0-CF-F16-16",
	"VS0-CF-F16-17",
	"VS0-CF-F16-18",
	"VS0-CF-F16-19",
	"VS0-CF-F16-20",
	"VS0-CF-F16-21",
	"VS0-CF-F16-22",
	"VS0-CF-F16-23",
	"VS0-CF-F16-24",
	"VS0-CF-F16-25",
	"VS0-CF-F16-26",
	"VS0-CF-F16-27",
	"VS0-CF-F16-28",
	"VS0-CF-F16-29",
	"VS0-CF-F16-30",
	"VS0-CF-F16-31",
	"VS0-CF-F16-32",
	"VS0-CF-F16-33",
	"VS0-CF-F16-34",
	"VS0-CF-F16-35",
	"VS0-CF-F16-36",
	"VS0-CF-F16-37",
	"VS0-CF-F16-38",
	"VS0-CF-F16-39",
	"VS0-CF-F16-40",
	"VS0-CF-F16-41",
	"VS0-CF-F16-42",
	"VS0-CF-F16-43",
	"VS0-CF-F16-44",
	"VS0-CF-F16-45",
	"VS0-CF-F16-46",
	"VS0-CF-F16-47",
	"VS0-CF-F16-48",
	"VS0-CF-F16-49",
	"VS0-CF-F16-50",
	"VS0-CF-F16-51",
	"VS0-CF-F16-52",
	"VS0-CF-F16-53",
	"VS0-CF-F16-54",
	"VS0-CF-F16-55",
	"VS0-CF-F16-56",
	"VS0-CF-F16-57",
	"VS0-CF-F16-58",
	"VS0-CF-F16-59",
	"VS0-CF-F16-60",
	"VS0-CF-F16-61",
	"VS0-CF-F16-62",
	"VS0-CF-F16-63",
	"VS0-CF-F16-64",
	"VS0-CF-F16-65",
	"VS0-CF-F16-66",
	"VS0-CF-F16-67",
	"VS0-CF-F16-68",
	"VS0-CF-F16-69",
	"VS0-CF-F16-70",
	"VS0-CF-F16-71",
	"VS0-CF-F16-72",
	"VS0-CF-F16-73",
	"VS0-CF-F16-74",
	"VS0-CF-F16-75",
	"VS0-CF-F16-76",
	"VS0-CF-F16-77",
	"VS0-CF-F16-78",
	"VS0-CF-F16-79",
	"VS0-CF-F16-80",
	"VS0-CF-F16-81",
	"VS0-CF-F16-82",
	"VS0-CF-F16-83",
	"VS0-CF-F16-84",
	"VS0-CF-F16-85",
	"VS0-CF-F16-86",
	"VS0-CF-F16-87",
	"VS0-CF-F16-88",
	"VS0-CF-F16-89",
	"VS0-CF-F16-90",
	"VS0-CF-F16-91",
	"VS0-CF-F16-92",
	"VS0-CF-F16-93",
	"VS0-CF-F16-94",
	"VS0-CF-F16-95",
	"VS0-CF-F16-96",
	"VS0-CF-F16-97",
	"VS0-CF-F16-98",
	"VS0-CF-F16-99",
	"VS0-CF-F16-100",
	"VS0-CF-F16-101",
	"VS0-CF-F16-102",
	"VS0-CF-F16-103",
	"VS0-CF-F16-104",
	"VS0-CF-F16-105",
	"VS0-CF-F16-106",
	"VS0-CF-F16-107",
	"VS0-CF-F16-108",
	"VS0-CF-F16-109",
	"VS0-CF-F16-110",
	"VS0-CF-F16-111",
	"VS0-CF-F16-112",
	"VS0-CF-F16-113",
	"VS0-CF-F16-114",
	"VS0-CF-F16-115",
	"VS0-CF-F16-116",
	"VS0-CF-F16-117",
	"VS0-CF-F16-118",
	"VS0-CF-F16-119",
	"VS0-CF-F16-120",
	"VS0-CF-F16-121",
	"VS0-CF-F16-122",
}

func TestVS0_CF_F16_CatalogCoverage(t *testing.T) {
	seen := map[string]bool{}
	for _, id := range f16CatalogIDs {
		if seen[id] {
			t.Fatalf("duplicate catalog id %s", id)
		}
		seen[id] = true
	}
	for i := 1; i <= 122; i++ {
		id := fmt.Sprintf("VS0-CF-F16-%02d", i)
		if !seen[id] {
			t.Fatalf("missing catalog id %s", id)
		}
	}
	if len(f16CatalogIDs) != 122 {
		t.Fatalf("catalog size=%d, want 122", len(f16CatalogIDs))
	}
}

// f16CategoryGroups is the corrected, exhaustive TASK-F16-12 grouping (design §7.3
// extended so every ID 01..122 appears; several IDs are intentionally cross-listed).
// Explicitly includes the prior-omission set 03,15,18,37,51..53,61..64,93..97,111.
var f16CategoryGroups = map[string][]int{
	"create-field-boundary-and-scope-derivation": {
		1, 2, 4, 54, 55, 56, 57, 84, 91, 92,
	},
	"authorization-denial-without-required-grant": {
		3, 15, 18, 37,
	},
	"safe-access-before-graph-validation": {
		9, 10, 11, 12, 13, 17, 19, 65, 66, 98,
	},
	"route-guard-transport-outcomes": {
		44, 45, 48, 105, 106, 107, 108, 109, 110, 113, 114, 115, 116, 117, 118,
	},
	"route-guard-precedence": {
		5, 6, 7, 8, 24, 25, 58, 59, 60, 70, 71, 72, 73, 74, 85, 86,
	},
	"authentication-failures-four-path-five-registration": {
		61, 62, 63, 64, 93, 94, 95, 96, 97,
	},
	"idempotency-replay-reuse-isolation-eviction-panic-waiter": {
		14, 16, 31, 32, 33, 34, 35, 46, 47, 49, 50, 76, 77, 78, 81, 82, 87, 88, 89, 99,
	},
	"malformed-duplicate-oversized-body": {
		51, 52, 53,
	},
	"observer-facts-conclusions": {
		20, 21, 22, 23, 67, 68, 69, 90,
	},
	"lifecycle-availability": {
		26, 27, 28, 29, 30, 36, 38, 39, 41, 42, 43, 75, 79, 80, 81, 82, 83, 88, 89, 100, 101, 119, 120, 121, 122,
	},
	"safe-accessible-without-grant-item-get": {
		111,
	},
	"audit-matrix-and-append-failure": {
		39, 40, 82, 99, 100, 101, 102, 103, 104, 112,
	},
	"future-boundary": {
		46,
	},
}

func TestVS0_CF_F16_CategoryCoverage(t *testing.T) {
	covered := map[int]bool{}
	for name, ids := range f16CategoryGroups {
		if len(ids) == 0 {
			t.Fatalf("category %q has no members", name)
		}
		for _, n := range ids {
			if n < 1 || n > 122 {
				t.Fatalf("category %q has out-of-range id %d", name, n)
			}
			covered[n] = true
		}
	}
	var missing []int
	for i := 1; i <= 122; i++ {
		if !covered[i] {
			missing = append(missing, i)
		}
	}
	if len(missing) != 0 {
		t.Fatalf("category groups omit case(s) %v", missing)
	}
	// Prior-omission correction set must remain explicitly categorized.
	for _, n := range []int{3, 15, 18, 37, 51, 52, 53, 61, 62, 63, 64, 93, 94, 95, 96, 97, 111} {
		if !covered[n] {
			t.Fatalf("prior-omission case %d must appear in a category group", n)
		}
	}
}

func f16WaitStarted(t *testing.T, started <-chan struct{}) {
	t.Helper()
	select {
	case <-started:
	case <-time.After(3 * time.Second):
		t.Fatal("did not reach observation")
	}
}

func f16WaitRec(t *testing.T, ch <-chan *httptest.ResponseRecorder) *httptest.ResponseRecorder {
	t.Helper()
	select {
	case rec := <-ch:
		return rec
	case <-time.After(3 * time.Second):
		t.Fatal("request did not finish")
		return nil
	}
}

func f16BlockingHarness(t *testing.T) (h *f16Harness, started, release chan struct{}) {
	t.Helper()
	fixtures := executiontarget.NewMapFixtureSource()
	started = make(chan struct{})
	release = make(chan struct{})
	blocking := &f16BlockingFixture{inner: fixtures, started: started, release: release}
	h = newF16Harness(t, f16HarnessOpts{
		fixtures: fixtures,
		observer: executiontarget.NewSyntheticObserver(blocking, nil),
	})
	return h, started, release
}

func f16QualifyNS(uid, key string) executiontarget.IdempotencyNamespace {
	return executiontarget.IdempotencyNamespace{
		PrincipalUID:     f16Principal,
		RoutePattern:     api.PatternExecutionTargetQualify,
		CloudProviderUID: f16ProviderUID,
		ActionTargetUID:  uid,
		Key:              key,
	}
}

func f16CreateNS(key string) executiontarget.IdempotencyNamespace {
	return executiontarget.IdempotencyNamespace{
		PrincipalUID:     f16Principal,
		RoutePattern:     api.PatternExecutionTargetCreate,
		CloudProviderUID: f16ProviderUID,
		Key:              key,
	}
}

func f16RetireNS(uid, key string) executiontarget.IdempotencyNamespace {
	return executiontarget.IdempotencyNamespace{
		PrincipalUID:     f16Principal,
		RoutePattern:     api.PatternExecutionTargetRetire,
		CloudProviderUID: f16ProviderUID,
		ActionTargetUID:  uid,
		Key:              key,
	}
}

func f16ResolveBacking(h *f16Harness, et etmodel.ExecutionTarget) executiontarget.BackingViability {
	out := executiontarget.BackingViability{
		ParticipationUID: et.Spec.CloudProviderParticipationRef.UID,
		StackUID:         et.Spec.InfrastructureStackRef.UID,
	}
	if part, ok := h.cloud.GetParticipation(et.Spec.CloudProviderParticipationRef.UID); ok {
		out.ParticipationEffectiveActive = part.Status.Phase == model.ParticipationPhaseActive && !part.Status.PlatformSuspended
		out.ParticipationScopeUID = part.Spec.CloudProviderRef.UID
	}
	if stack, ok := h.cloud.GetInfrastructureStack(et.Spec.InfrastructureStackRef.UID); ok {
		out.StackPhase = string(stack.Status.Phase)
		out.StackGeneration = stack.Metadata.Generation
		if stack.Metadata.ScopeRef != nil {
			out.StackScopeUID = stack.Metadata.ScopeRef.UID
		}
	}
	return out
}

func f16QualifyViaLifecycle(t *testing.T, h *f16Harness, uid, key string) (etmodel.ExecutionTarget, *apiproblem.Problem) {
	t.Helper()
	et, ok := h.life.Store().GetExecutionTarget(uid)
	if !ok {
		t.Fatalf("target %s missing", uid)
	}
	ns := f16QualifyNS(uid, key)
	if res := h.life.ReserveOrInspectReplay(ns, executiontarget.DigestAction()); res.Outcome != executiontarget.ReservationOutcomeReserved {
		t.Fatalf("reserve qualify %s: %#v", key, res)
	}
	backing := f16ResolveBacking(h, et)
	spec := et.Spec
	req := executiontarget.QualifyCommitRequest{
		TargetUID: uid,
		Backing:   backing,
		RefreshBacking: func() executiontarget.BackingViability {
			// Must not call F0016 Store getters (would re-enter store mutex under DD-02).
			return f16ResolveBacking(h, etmodel.ExecutionTarget{Spec: spec})
		},
		FactSetUID:  "fs-" + key,
		ResultUID:   "qr-" + key,
		Actor:       f16ActorRef(f16Principal),
		RequestID:   "req-" + key,
		AuditUID:    "audit-" + key,
		Idempotency: ns,
		Digest:      executiontarget.DigestAction(),
		Completion:  executiontarget.CompletedResult{StatusCode: http.StatusOK, Body: []byte(`{}`)},
	}
	return h.life.Qualify(context.Background(), req)
}

func f16DecodeSafe(t *testing.T, rec *httptest.ResponseRecorder) executiontarget.SafeExecutionTarget {
	t.Helper()
	var got executiontarget.SafeExecutionTarget
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode safe: %v body=%s", err, rec.Body.String())
	}
	return got
}

func f16PadIndex(i int) string {
	return fmt.Sprintf("%05d", i)
}

func f16AssertInternalNoDisclosure(t *testing.T, rec *httptest.ResponseRecorder, forbiddenStatus int) {
	t.Helper()
	p := decodeF16Problem(t, rec)
	if rec.Code != http.StatusInternalServerError || p.Code != apiproblem.CodeInternalError {
		t.Fatalf("got %d %#v", rec.Code, p)
	}
	if p.Status == forbiddenStatus {
		t.Fatalf("must not disclose suppressed status %d", forbiddenStatus)
	}
}

func TestVS0_CF_F16_01(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	h.seedActiveParticipationAndStack(t)
	before := h.auditLen()
	rec := doF16(t, h, http.MethodPost, f16Collection, f16CreateBody("t-01", f16PartUID, f16StackUID), f16CreateHeaders("k-01"))
	if rec.Code != http.StatusCreated {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	assertF16JSONNosniff(t, rec)
	if rec.Header().Get("ETag") == "" {
		t.Fatal("create must return ETag")
	}
	if rec.Header().Get("Location") != "" || rec.Header().Get("Idempotency-Key") != "" {
		t.Fatal("must not echo Location or Idempotency-Key")
	}
	created := f16DecodeSafe(t, rec)
	if created.Status.Lifecycle != etmodel.LifecycleActive ||
		created.Status.Qualification != etmodel.QualificationUnqualified ||
		created.EffectiveAvailability != etmodel.EffectiveUnavailable {
		t.Fatalf("projection=%#v", created)
	}
	stored, ok := h.life.Store().GetExecutionTarget(created.Metadata.UID)
	if !ok || stored.Status.MaintenanceEpoch != 0 || stored.Status.ObservedGeneration != 3 {
		t.Fatalf("stored=%#v ok=%v", stored.Status, ok)
	}
	assertF16AuditDelta(t, h, before, 1)
	assertF16SafeProjectionOmitsInternals(t, rec.Body.String())

	rec2 := doF16(t, h, http.MethodPost, f16Collection, f16CreateBody("t-01", f16PartUID, f16StackUID), f16CreateHeaders("k-01"))
	if rec2.Code != http.StatusCreated {
		t.Fatalf("replay status=%d", rec2.Code)
	}
	if !bytes.Equal(bytes.TrimSpace(rec2.Body.Bytes()), bytes.TrimSpace(rec.Body.Bytes())) {
		t.Fatal("replay body mismatch")
	}
	assertF16AuditDelta(t, h, before, 1)
}

func TestVS0_CF_F16_02(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	h.seedActiveParticipationAndStack(t)
	before := h.auditLen()
	body := []byte(`{"metadata":{"name":"t-02"},"spec":{"cloudProviderParticipationRef":{"uid":"` + f16PartUID + `"},"infrastructureStackRef":{"uid":"` + f16StackUID + `"},"targetClass":"synthetic-iaas","unknownField":1}}`)
	rec := doF16(t, h, http.MethodPost, f16Collection, body, f16CreateHeaders("k-02"))
	assertF16Problem(t, rec, apiproblem.CodeUnknownField, http.StatusBadRequest, "")
	assertF16NoAuditDelta(t, h, before)
	if len(h.life.Store().ListExecutionTargets()) != 0 {
		t.Fatal("must not publish")
	}
}

func TestVS0_CF_F16_03(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	h.seedActiveParticipationAndStack(t)
	before := h.auditLen()
	body := []byte(`{"metadata":{"name":"t-03"},"spec":{"cloudProviderParticipationRef":{"uid":"` + f16PartUID + `"},"infrastructureStackRef":{"uid":"` + f16StackUID + `"},"targetClass":"synthetic-iaas"},"status":{"lifecycle":"Retired"}}`)
	rec := doF16(t, h, http.MethodPost, f16Collection, body, f16CreateHeaders("k-03"))
	assertF16Problem(t, rec, apiproblem.CodeAuthorizationDenied, http.StatusForbidden, validate.ViolationStatusFieldWrite)
	assertF16AuditDelta(t, h, before, 1)
}

func TestVS0_CF_F16_04(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	h.seedActiveParticipationAndStack(t)
	before := h.auditLen()
	body := []byte(`{"metadata":{"name":"t-04","uid":"cccccccccccccccccccccccccccccccc"},"spec":{"cloudProviderParticipationRef":{"uid":"` + f16PartUID + `"},"infrastructureStackRef":{"uid":"` + f16StackUID + `"},"targetClass":"synthetic-iaas"}}`)
	rec := doF16(t, h, http.MethodPost, f16Collection, body, f16CreateHeaders("k-04"))
	assertF16Problem(t, rec, apiproblem.CodeAuthorizationDenied, http.StatusForbidden, validate.ViolationSystemOwnedFieldWrite)
	assertF16AuditDelta(t, h, before, 1)
}

func TestVS0_CF_F16_05(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	h.seedActiveParticipationAndStack(t)
	before := h.auditLen()
	rec := doF16(t, h, http.MethodPost, f16Collection, f16CreateBody("t-05", f16PartUID, f16StackUID), map[string]string{"Content-Type": f16MediaJSON})
	assertF16Problem(t, rec, apiproblem.CodeMalformedRequest, http.StatusBadRequest, "")
	assertF16NoAuditDelta(t, h, before)
	if len(h.life.Store().ListExecutionTargets()) != 0 {
		t.Fatal("must not publish")
	}
}

func TestVS0_CF_F16_06(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	h.seedActiveParticipationAndStack(t)
	before := h.auditLen()
	rec := doF16(t, h, http.MethodPost, f16Collection, f16CreateBody("t-06", f16PartUID, f16StackUID), map[string]string{
		"Content-Type": "text/plain", "Idempotency-Key": "k-06",
	})
	assertF16Problem(t, rec, apiproblem.CodeUnsupportedMediaType, http.StatusUnsupportedMediaType, "")
	assertF16NoAuditDelta(t, h, before)
}

func TestVS0_CF_F16_07(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	h.seedActiveParticipationAndStack(t)
	before := h.auditLen()
	rec := doF16Unauth(t, h, http.MethodPost, f16Collection, f16CreateBody("t-07", f16PartUID, f16StackUID), f16CreateHeaders("k-07"))
	assertF16Problem(t, rec, apiproblem.CodeAuthRequired, http.StatusUnauthorized, "")
	assertF16NoAuditDelta(t, h, before)
}

func TestVS0_CF_F16_08(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{grants: f16ExecutionTargetGrants(f16ProviderUID, false, true, true)})
	h.seedActiveParticipationAndStack(t)
	before := h.auditLen()
	rec := doF16(t, h, http.MethodPost, f16Collection, f16CreateBody("t-08", f16PartUID, f16StackUID), f16CreateHeaders("k-08"))
	assertF16Problem(t, rec, apiproblem.CodeAuthorizationDenied, http.StatusForbidden, "")
	assertF16AuditDelta(t, h, before, 1)
	if len(h.life.Store().ListExecutionTargets()) != 0 {
		t.Fatal("must not publish")
	}
}

func TestVS0_CF_F16_09(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	h.seedPlatformAndProvider(t, f16ProviderUID)
	before := h.auditLen()
	rec := doF16(t, h, http.MethodPost, f16Collection, f16CreateBody("t-09", f16PartUID, f16StackUID), f16CreateHeaders("k-09"))
	assertF16Problem(t, rec, apiproblem.CodeResourceNotFound, http.StatusNotFound, f16ViolationSafeDeny)
	assertF16AuditDelta(t, h, before, 1)
}

func TestVS0_CF_F16_10(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	h.seedActiveParticipationAndStack(t)
	h.setParticipationEffective(t, f16PartUID, false)
	before := h.auditLen()
	rec := doF16(t, h, http.MethodPost, f16Collection, f16CreateBody("t-10", f16PartUID, f16StackUID), f16CreateHeaders("k-10"))
	assertF16Problem(t, rec, apiproblem.CodeConflict, http.StatusConflict, executiontarget.ViolationExecutionTargetParticipationUnavailable)
	assertF16NoAuditDelta(t, h, before)
}

func TestVS0_CF_F16_11(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	h.seedActiveParticipationAndStack(t)
	inactive := h.seedExtraStack(t, "stack-inactive", "acacacacacacacacacacacacacacacac", 1)
	h.setStackPhase(t, inactive.Metadata.UID, "")
	before := h.auditLen()
	rec := doF16(t, h, http.MethodPost, f16Collection, f16CreateBody("t-11", f16PartUID, inactive.Metadata.UID), f16CreateHeaders("k-11"))
	assertF16Problem(t, rec, apiproblem.CodeConflict, http.StatusConflict, executiontarget.ViolationExecutionTargetStackUnavailable)
	assertF16NoAuditDelta(t, h, before)
}

func TestVS0_CF_F16_12(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	h.seedActiveParticipationAndStack(t)
	other := h.seedOtherProviderAndStack(t, "66666666666666666666666666666666")
	before := h.auditLen()
	rec := doF16(t, h, http.MethodPost, f16Collection, f16CreateBody("t-12", f16PartUID, other.Metadata.UID), f16CreateHeaders("k-12"))
	assertF16Problem(t, rec, apiproblem.CodeValidationFailed, http.StatusUnprocessableEntity, executiontarget.ViolationExecutionTargetScopeMismatch)
	assertF16NoAuditDelta(t, h, before)
}

func TestVS0_CF_F16_13(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	h.seedActiveParticipationAndStack(t)
	rec := doF16(t, h, http.MethodPost, f16Collection, f16CreateBody("dup-name", f16PartUID, f16StackUID), f16CreateHeaders("dup-1"))
	if rec.Code != http.StatusCreated {
		t.Fatalf("first=%d %s", rec.Code, rec.Body.String())
	}
	stackB := h.seedExtraStack(t, "stack-b", "77777777777777777777777777777777", 1)
	before := h.auditLen()
	rec = doF16(t, h, http.MethodPost, f16Collection, f16CreateBody("dup-name", f16PartUID, stackB.Metadata.UID), f16CreateHeaders("dup-2"))
	assertF16Problem(t, rec, apiproblem.CodeAlreadyExists, http.StatusConflict, "")
	assertF16NoAuditDelta(t, h, before)
}

func TestVS0_CF_F16_14(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	h.seedActiveParticipationAndStack(t)
	recB := doF16(t, h, http.MethodPost, f16Collection, f16CreateBody("target-b", f16PartUID, f16StackUID), f16CreateHeaders("et-list-b"))
	if recB.Code != http.StatusCreated {
		t.Fatalf("create b=%d", recB.Code)
	}
	createdB := f16DecodeSafe(t, recB)
	stack2 := h.seedExtraStack(t, "stack-2", "44444444444444444444444444444444", 1)
	_ = doF16(t, h, http.MethodPost, f16Collection, f16CreateBody("target-a", f16PartUID, stack2.Metadata.UID), f16CreateHeaders("et-list-a"))

	retireNS := f16RetireNS(createdB.Metadata.UID, "retire-list-b")
	if res := h.life.ReserveOrInspectReplay(retireNS, executiontarget.DigestAction()); res.Outcome != executiontarget.ReservationOutcomeReserved {
		t.Fatalf("reserve retire: %#v", res)
	}
	if _, retireProb := h.life.CommitRetire(context.Background(), executiontarget.RetireCommitRequest{
		TargetUID:               createdB.Metadata.UID,
		ExpectedResourceVersion: createdB.Metadata.ResourceVersion,
		Actor:                   f16ActorRef(f16Principal),
		RequestID:               "req-retire-list",
		AuditUID:                "audit-retire-list",
		Idempotency:             retireNS,
		Digest:                  executiontarget.DigestAction(),
		Completion:              executiontarget.CompletedResult{StatusCode: http.StatusOK, Body: []byte(`{}`)},
	}); retireProb != nil {
		t.Fatalf("retire: %#v", retireProb)
	}

	before := h.auditLen()
	rec := doF16(t, h, http.MethodGet, f16Collection, nil, map[string]string{"Idempotency-Key": "ignored-on-list", "Accept": "text/plain"})
	if rec.Code != http.StatusOK {
		t.Fatalf("list status=%d body=%s", rec.Code, rec.Body.String())
	}
	assertF16JSONNosniff(t, rec)
	if rec.Header().Get("ETag") != "" {
		t.Fatal("LIST must not return ETag")
	}
	assertF16NoAuditDelta(t, h, before)
	var items []executiontarget.SafeExecutionTarget
	if err := json.Unmarshal(rec.Body.Bytes(), &items); err != nil {
		t.Fatalf("list decode: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("list len=%d", len(items))
	}
	if items[0].Metadata.UID > items[1].Metadata.UID {
		t.Fatal("LIST must be ascending UID")
	}
	var sawActive, sawRetired bool
	for _, item := range items {
		switch item.Status.Lifecycle {
		case etmodel.LifecycleActive:
			sawActive = true
		case etmodel.LifecycleRetired:
			sawRetired = true
		}
	}
	if !sawActive || !sawRetired {
		t.Fatalf("LIST must include Active and Retired; items=%#v", items)
	}
}

func TestVS0_CF_F16_15(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{grants: f16ExecutionTargetGrants(f16ProviderUID, true, false, true)})
	before := h.auditLen()
	rec := doF16(t, h, http.MethodGet, f16Collection, nil, nil)
	assertF16Problem(t, rec, apiproblem.CodeAuthorizationDenied, http.StatusForbidden, "")
	assertF16AuditDelta(t, h, before, 1)
}

func TestVS0_CF_F16_16(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	created := h.createTarget(t, "item-get-ok", "item-create-ok")
	before := h.auditLen()
	rec := doF16(t, h, http.MethodGet, f16ItemPath(created.Metadata.UID), nil, map[string]string{
		"Idempotency-Key": "ignored-on-item-get", "Accept": "application/xml",
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	assertF16JSONNosniff(t, rec)
	if got := rec.Header().Get("ETag"); got == "" || got != created.Metadata.ResourceVersion {
		t.Fatalf("ETag=%q want %q", got, created.Metadata.ResourceVersion)
	}
	assertF16NoAuditDelta(t, h, before)
	assertF16SafeProjectionOmitsInternals(t, rec.Body.String())
	ns := f16CreateNS("ignored-on-item-get")
	inspect := h.life.ReserveOrInspectReplay(ns, executiontarget.DigestAction())
	if inspect.Outcome != executiontarget.ReservationOutcomeReserved {
		t.Fatalf("GET must not consume key; outcome=%v", inspect.Outcome)
	}
	h.life.AbortReservation(ns)
}

func TestVS0_CF_F16_17(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	before := h.auditLen()
	rec := doF16(t, h, http.MethodGet, f16ItemPath("bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"), nil, nil)
	assertF16Problem(t, rec, apiproblem.CodeResourceNotFound, http.StatusNotFound, f16ViolationSafeDeny)
	assertF16AuditDelta(t, h, before, 1)
}

func TestVS0_CF_F16_18(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{grants: f16ExecutionTargetGrants(f16ProviderUID, true, true, false)})
	created := h.createTarget(t, "no-qual", "create-no-qual")
	before := h.auditLen()
	rec := doF16(t, h, http.MethodPost, f16ActionPath(created.Metadata.UID, "qualify"), nil,
		f16ActionHeaders("q-deny", created.Metadata.ResourceVersion))
	assertF16Problem(t, rec, apiproblem.CodeAuthorizationDenied, http.StatusForbidden, "")
	assertF16AuditDelta(t, h, before, 1)
}

func TestVS0_CF_F16_19(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	h.seedActiveParticipationAndStack(t)
	before := h.auditLen()
	missing := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	rec := doF16(t, h, http.MethodPost, f16ActionPath(missing, "qualify"), nil, f16ActionHeaders("missing", "1"))
	assertF16Problem(t, rec, apiproblem.CodeResourceNotFound, http.StatusNotFound, f16ViolationSafeDeny)
	assertF16AuditDelta(t, h, before, 1)
}

func TestVS0_CF_F16_20(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	created := h.createTarget(t, "q-ok", "create-q-ok")
	h.setPresentFixture(created.Metadata.UID, f16AllSupportedTruths())
	before := h.auditLen()
	etag := created.Metadata.ResourceVersion
	rec := doF16(t, h, http.MethodPost, f16ActionPath(created.Metadata.UID, "qualify"), nil,
		f16ActionHeaders("qualify-1", etag))
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	assertF16JSONNosniff(t, rec)
	if rec.Header().Get("ETag") == "" || rec.Header().Get("ETag") == etag {
		t.Fatalf("ETag must advance; got %q", rec.Header().Get("ETag"))
	}
	proj := f16DecodeSafe(t, rec)
	if proj.Status.Qualification != etmodel.QualificationQualified ||
		proj.EffectiveAvailability != etmodel.EffectiveAvailable {
		t.Fatalf("projection=%#v", proj)
	}
	assertF16AuditDelta(t, h, before, 1)
	assertF16SafeProjectionOmitsInternals(t, rec.Body.String())
	rec2 := doF16(t, h, http.MethodPost, f16ActionPath(created.Metadata.UID, "qualify"), nil,
		f16ActionHeaders("qualify-1", etag))
	if rec2.Code != http.StatusOK || rec2.Body.String() != rec.Body.String() {
		t.Fatal("replay must match")
	}
	assertF16AuditDelta(t, h, before, 1)
}

func TestVS0_CF_F16_21(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	created := h.createTarget(t, "q-rej", "create-q-rej")
	h.setPresentFixture(created.Metadata.UID, f16UnsupportedVMTruths())
	rec := doF16(t, h, http.MethodPost, f16ActionPath(created.Metadata.UID, "qualify"), nil,
		f16ActionHeaders("qualify-rej", created.Metadata.ResourceVersion))
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d %s", rec.Code, rec.Body.String())
	}
	proj := f16DecodeSafe(t, rec)
	if proj.Status.Qualification != etmodel.QualificationRejected ||
		proj.EffectiveAvailability != etmodel.EffectiveUnavailable {
		t.Fatalf("projection=%#v", proj)
	}
}

func TestVS0_CF_F16_22(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	created := h.createTarget(t, "q-miss", "create-q-miss")
	// default: missing fixture
	rec := doF16(t, h, http.MethodPost, f16ActionPath(created.Metadata.UID, "qualify"), nil,
		f16ActionHeaders("qualify-miss", created.Metadata.ResourceVersion))
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d %s", rec.Code, rec.Body.String())
	}
	proj := f16DecodeSafe(t, rec)
	if proj.Status.Qualification != etmodel.QualificationIndeterminate ||
		proj.EffectiveAvailability != etmodel.EffectiveUnavailable {
		t.Fatalf("projection=%#v", proj)
	}
	stored, _ := h.life.Store().GetExecutionTarget(created.Metadata.UID)
	if stored.Status.FactSetRef == nil {
		t.Fatal("fact set ref required")
	}
	fs, ok := h.life.Store().GetFactSet(stored.Status.FactSetRef.UID)
	if !ok || fs.Facts.Compute.VM != etmodel.FactUnknown {
		t.Fatalf("want Unknown facts: %#v ok=%v", fs, ok)
	}
}

func TestVS0_CF_F16_23(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	created := h.createTarget(t, "q-mal", "create-q-mal")
	h.setPresentFixture(created.Metadata.UID, f16AllSupportedTruths())
	rec := doF16(t, h, http.MethodPost, f16ActionPath(created.Metadata.UID, "qualify"), nil,
		f16ActionHeaders("q-first", created.Metadata.ResourceVersion))
	if rec.Code != http.StatusOK {
		t.Fatalf("first qualify: %d", rec.Code)
	}
	qualified := f16DecodeSafe(t, rec)
	beforeRV := qualified.Metadata.ResourceVersion
	beforeAudit := h.auditLen()
	h.fixtures.Set(created.Metadata.UID, executiontarget.TargetBoundFixture{
		Revision: "bad", Mode: executiontarget.FixtureFaultMalformed, Truths: f16AllSupportedTruths(),
	})
	rec2 := doF16(t, h, http.MethodPost, f16ActionPath(created.Metadata.UID, "qualify"), nil,
		f16ActionHeaders("q-mal", beforeRV))
	assertF16Problem(t, rec2, apiproblem.CodeInternalError, http.StatusInternalServerError, "")
	assertF16NoAuditDelta(t, h, beforeAudit)
	after, _ := h.life.Store().GetExecutionTarget(created.Metadata.UID)
	if after.Metadata.ResourceVersion != beforeRV {
		t.Fatal("must leave ETag unchanged")
	}
	if after.Status.FactSetRef == nil {
		t.Fatal("must retain FactSet links")
	}
}

func TestVS0_CF_F16_24(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	created := h.createTarget(t, "grammar", "create-grammar")
	before := h.auditLen()
	rec := doF16(t, h, http.MethodPost, f16ActionPath(created.Metadata.UID, "qualify"), nil, map[string]string{
		"Content-Type": f16MediaJSON, "Idempotency-Key": "no-if",
	})
	assertF16Problem(t, rec, apiproblem.CodeStaleResourceVersion, http.StatusPreconditionFailed, "")
	assertF16NoAuditDelta(t, h, before)
}

func TestVS0_CF_F16_25(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	created := h.createTarget(t, "body", "create-body")
	before := h.auditLen()
	rec := doF16(t, h, http.MethodPost, f16ActionPath(created.Metadata.UID, "qualify"), []byte(`{}`),
		f16ActionHeaders("body", created.Metadata.ResourceVersion))
	assertF16Problem(t, rec, apiproblem.CodeMalformedRequest, http.StatusBadRequest, "")
	assertF16NoAuditDelta(t, h, before)
}

func TestVS0_CF_F16_26(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	created := h.createTarget(t, "retired", "create-retired")
	rec := doF16(t, h, http.MethodPost, f16ActionPath(created.Metadata.UID, "retire"), nil,
		f16ActionHeaders("retire-first", created.Metadata.ResourceVersion))
	if rec.Code != http.StatusOK {
		t.Fatalf("retire: %d", rec.Code)
	}
	retired := f16DecodeSafe(t, rec)
	q := doF16(t, h, http.MethodPost, f16ActionPath(created.Metadata.UID, "qualify"), nil,
		f16ActionHeaders("q-retired", retired.Metadata.ResourceVersion))
	assertF16Problem(t, q, apiproblem.CodeConflict, http.StatusConflict, executiontarget.ViolationTargetRetired)
}

func TestVS0_CF_F16_27(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	created := h.createTarget(t, "maint", "create-maint")
	et, _ := h.life.Store().GetExecutionTarget(created.Metadata.UID)
	_ = h.deliverMaintenance(t, et, executiontarget.MaintenanceTriggerEnter, "audit-maint")
	et, _ = h.life.Store().GetExecutionTarget(created.Metadata.UID)
	q := doF16(t, h, http.MethodPost, f16ActionPath(created.Metadata.UID, "qualify"), nil,
		f16ActionHeaders("q-maint", et.Metadata.ResourceVersion))
	assertF16Problem(t, q, apiproblem.CodeConflict, http.StatusConflict, executiontarget.ViolationTargetMaintenance)
}

func TestVS0_CF_F16_28(t *testing.T) {
	h, started, release := f16BlockingHarness(t)
	created := h.createTarget(t, "inprog", "create-inprog")
	go func() {
		_ = doF16(t, h, http.MethodPost, f16ActionPath(created.Metadata.UID, "qualify"), nil,
			f16ActionHeaders("owner-q", created.Metadata.ResourceVersion))
	}()
	f16WaitStarted(t, started)
	rec := doF16(t, h, http.MethodPost, f16ActionPath(created.Metadata.UID, "qualify"), nil,
		f16ActionHeaders("other-key", created.Metadata.ResourceVersion))
	assertF16Problem(t, rec, apiproblem.CodeConflict, http.StatusConflict, executiontarget.ViolationTargetQualificationInProgress)
	close(release)
}

func TestVS0_CF_F16_29(t *testing.T) {
	h, started, release := f16BlockingHarness(t)
	created := h.createTarget(t, "gen-stale", "create-gen")
	h.setPresentFixture(created.Metadata.UID, f16AllSupportedTruths())
	beforeAudit := h.auditLen()
	beforeRV := created.Metadata.ResourceVersion
	qualDone := make(chan *httptest.ResponseRecorder, 1)
	go func() {
		qualDone <- doF16(t, h, http.MethodPost, f16ActionPath(created.Metadata.UID, "qualify"), nil,
			f16ActionHeaders("q-gen", created.Metadata.ResourceVersion))
	}()
	f16WaitStarted(t, started)
	h.updateStackGeneration(t, f16StackUID, 99)
	close(release)
	qRec := f16WaitRec(t, qualDone)
	assertF16Problem(t, qRec, apiproblem.CodeStaleResourceVersion, http.StatusPreconditionFailed, executiontarget.ViolationTargetEpochStale)
	assertF16NoAuditDelta(t, h, beforeAudit)
	after, _ := h.life.Store().GetExecutionTarget(created.Metadata.UID)
	if after.Metadata.ResourceVersion != beforeRV || after.Status.FactSetRef != nil {
		t.Fatal("must leave committed target/links unchanged")
	}
}

func TestVS0_CF_F16_30(t *testing.T) {
	h, started, release := f16BlockingHarness(t)
	created := h.createTarget(t, "via-stale", "create-via")
	h.setPresentFixture(created.Metadata.UID, f16AllSupportedTruths())
	beforeAudit := h.auditLen()
	beforeRV := created.Metadata.ResourceVersion
	qualDone := make(chan *httptest.ResponseRecorder, 1)
	go func() {
		qualDone <- doF16(t, h, http.MethodPost, f16ActionPath(created.Metadata.UID, "qualify"), nil,
			f16ActionHeaders("q-via", created.Metadata.ResourceVersion))
	}()
	f16WaitStarted(t, started)
	h.setParticipationEffective(t, f16PartUID, false)
	close(release)
	qRec := f16WaitRec(t, qualDone)
	assertF16Problem(t, qRec, apiproblem.CodeStaleResourceVersion, http.StatusPreconditionFailed, executiontarget.ViolationExecutionTargetViabilityStale)
	assertF16NoAuditDelta(t, h, beforeAudit)
	after, _ := h.life.Store().GetExecutionTarget(created.Metadata.UID)
	if after.Metadata.ResourceVersion != beforeRV || after.Status.FactSetRef != nil {
		t.Fatal("must leave committed target/links unchanged")
	}
}

func TestVS0_CF_F16_31(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	h.seedActiveParticipationAndStack(t)
	headers := f16CreateHeaders("k-31")
	rec := doF16(t, h, http.MethodPost, f16Collection, f16CreateBody("t-31", f16PartUID, f16StackUID), headers)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create=%d", rec.Code)
	}
	before := h.auditLen()
	firstBody := append([]byte(nil), rec.Body.Bytes()...)
	rec2 := doF16(t, h, http.MethodPost, f16Collection, f16CreateBody("t-31", f16PartUID, f16StackUID), headers)
	if rec2.Code != http.StatusCreated {
		t.Fatalf("replay=%d", rec2.Code)
	}
	assertF16JSONNosniff(t, rec2)
	if !bytes.Equal(bytes.TrimSpace(rec2.Body.Bytes()), bytes.TrimSpace(firstBody)) {
		t.Fatal("replay body mismatch")
	}
	assertF16NoAuditDelta(t, h, before)
}

func TestVS0_CF_F16_32(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	created := h.createTarget(t, "f32", "create-f32")
	headers := f16ActionHeaders("replay-authz", created.Metadata.ResourceVersion)
	rec := doF16(t, h, http.MethodPost, f16ActionPath(created.Metadata.UID, "qualify"), nil, headers)
	if rec.Code != http.StatusOK {
		t.Fatalf("qualify: %d", rec.Code)
	}
	h.rebindGrants(f16ExecutionTargetGrants(f16ProviderUID, true, true, false))
	before := h.auditLen()
	rec2 := doF16(t, h, http.MethodPost, f16ActionPath(created.Metadata.UID, "qualify"), nil, headers)
	assertF16Problem(t, rec2, apiproblem.CodeAuthorizationDenied, http.StatusForbidden, "")
	assertF16AuditDelta(t, h, before, 1)
	if strings.Contains(rec2.Body.String(), created.Metadata.UID) && strings.Contains(rec2.Body.String(), "Qualified") {
		t.Fatal("must not disclose stored-result")
	}
}

func TestVS0_CF_F16_33(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	h.seedActiveParticipationAndStack(t)
	stackC := h.seedExtraStack(t, "stack-c", "88888888888888888888888888888888", 1)
	headers := f16CreateHeaders("reuse-key")
	rec := doF16(t, h, http.MethodPost, f16Collection, f16CreateBody("reuse-a", f16PartUID, stackC.Metadata.UID), headers)
	if rec.Code != http.StatusCreated {
		t.Fatalf("first=%d %s", rec.Code, rec.Body.String())
	}
	rec = doF16(t, h, http.MethodPost, f16Collection, f16CreateBody("reuse-b", f16PartUID, stackC.Metadata.UID), headers)
	assertF16Problem(t, rec, apiproblem.CodeConflict, http.StatusConflict, cloudmodel.ViolationIdempotencyKeyReuseMismatch)
}

func TestVS0_CF_F16_34(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	h.seedActiveParticipationAndStack(t)
	headers := f16CreateHeaders("k-34")
	rec := doF16(t, h, http.MethodPost, f16Collection, f16CreateBody("t-34", f16PartUID, f16StackUID), headers)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create=%d", rec.Code)
	}
	uid1 := f16DecodeSafe(t, rec).Metadata.UID
	h.clock.Advance(executiontarget.CompletedRetention + time.Second)
	// Same key processed as new after expiry.
	stack2 := h.seedExtraStack(t, "stack-34b", "a4a4a4a4a4a4a4a4a4a4a4a4a4a4a4a4", 1)
	rec2 := doF16(t, h, http.MethodPost, f16Collection, f16CreateBody("t-34b", f16PartUID, stack2.Metadata.UID), headers)
	if rec2.Code != http.StatusCreated {
		t.Fatalf("new after expiry status=%d %s", rec2.Code, rec2.Body.String())
	}
	uid2 := f16DecodeSafe(t, rec2).Metadata.UID
	if uid1 == uid2 {
		t.Fatal("expired completion must not replay; must create new")
	}
}

func TestVS0_CF_F16_35(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	h.seedActiveParticipationAndStack(t)
	ns := f16CreateNS("panic-key")
	stackP := h.seedExtraStack(t, "stack-panic", "abababababababababababababababab", 1)
	body := f16CreateBody("panic-t", f16PartUID, stackP.Metadata.UID)
	digest, err := executiontarget.DigestCreate(json.RawMessage(body))
	if err != nil {
		t.Fatal(err)
	}
	if res := h.life.ReserveOrInspectReplay(ns, digest); res.Outcome != executiontarget.ReservationOutcomeReserved {
		t.Fatalf("reserve: %#v", res)
	}
	done, ok := h.life.AttachWaiter(ns)
	if !ok {
		t.Fatal("attach waiter")
	}
	var woke atomic.Bool
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-done
		woke.Store(true)
	}()
	h.life.AbortReservation(ns)
	wg.Wait()
	if !woke.Load() {
		t.Fatal("waiter must wake after abort/panic cleanup")
	}
	if len(h.life.Store().ListExecutionTargets()) != 0 {
		t.Fatal("panic/abort path must not publish")
	}
}

func TestVS0_CF_F16_36(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	created := h.createTarget(t, "retire-ok", "create-r-ok")
	before := h.auditLen()
	rec := doF16(t, h, http.MethodPost, f16ActionPath(created.Metadata.UID, "retire"), nil,
		f16ActionHeaders("retire-1", created.Metadata.ResourceVersion))
	if rec.Code != http.StatusOK {
		t.Fatalf("retire status=%d %s", rec.Code, rec.Body.String())
	}
	assertF16JSONNosniff(t, rec)
	proj := f16DecodeSafe(t, rec)
	if proj.Status.Lifecycle != etmodel.LifecycleRetired ||
		proj.Status.Qualification != etmodel.QualificationUnqualified ||
		proj.EffectiveAvailability != etmodel.EffectiveUnavailable {
		t.Fatalf("projection=%#v", proj)
	}
	assertF16AuditDelta(t, h, before, 1)
	stored, ok := h.life.Store().GetExecutionTarget(created.Metadata.UID)
	if !ok || stored.Status.FactSetRef != nil || stored.Status.QualificationResultRef != nil {
		t.Fatal("retire must clear record refs")
	}
}

func TestVS0_CF_F16_37(t *testing.T) {
	hWrite := newF16Harness(t, f16HarnessOpts{})
	created := hWrite.createTarget(t, "no-write", "create-nw")
	h := newF16Harness(t, f16HarnessOpts{grants: f16ExecutionTargetGrants(f16ProviderUID, false, true, true)})
	// Share lifecycle/cloud/audit so target exists under denied grants.
	h.life, h.cloud, h.audit, h.fixtures = hWrite.life, hWrite.cloud, hWrite.audit, hWrite.fixtures
	h.rebindGrants(f16ExecutionTargetGrants(f16ProviderUID, false, true, true))
	before := hWrite.auditLen()
	rec := doF16(t, h, http.MethodPost, f16ActionPath(created.Metadata.UID, "retire"), nil,
		f16ActionHeaders("retire-deny", created.Metadata.ResourceVersion))
	assertF16Problem(t, rec, apiproblem.CodeAuthorizationDenied, http.StatusForbidden, "")
	if hWrite.auditLen() != before+1 {
		t.Fatal("retire authz denial must append exactly one AuditEvent")
	}
}

func TestVS0_CF_F16_38(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	created := h.createTarget(t, "retired2", "create-retired2")
	rec := doF16(t, h, http.MethodPost, f16ActionPath(created.Metadata.UID, "retire"), nil,
		f16ActionHeaders("retire-first", created.Metadata.ResourceVersion))
	retired := f16DecodeSafe(t, rec)
	r2 := doF16(t, h, http.MethodPost, f16ActionPath(created.Metadata.UID, "retire"), nil,
		f16ActionHeaders("retire-again", retired.Metadata.ResourceVersion))
	assertF16Problem(t, r2, apiproblem.CodeConflict, http.StatusConflict, executiontarget.ViolationTargetRetired)
}

func TestVS0_CF_F16_39(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	h.seedActiveParticipationAndStack(t)
	h.setAuditFail(errors.New("append boom"))
	rec := doF16(t, h, http.MethodPost, f16Collection, f16CreateBody("t-39", f16PartUID, f16StackUID), f16CreateHeaders("k-39"))
	assertF16Problem(t, rec, apiproblem.CodeInternalError, http.StatusInternalServerError, "")
	if h.auditLen() != 0 {
		t.Fatal("failed append must leave no durable AuditEvent")
	}
	if len(h.life.Store().ListExecutionTargets()) != 0 {
		t.Fatal("append failure must not publish")
	}
}

func TestVS0_CF_F16_40(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{grants: f16ExecutionTargetGrants(f16ProviderUID, false, true, true)})
	h.seedActiveParticipationAndStack(t)
	h.setAuditFail(errors.New("append boom"))
	rec := doF16(t, h, http.MethodPost, f16Collection, f16CreateBody("t-40", f16PartUID, f16StackUID), f16CreateHeaders("no-write-fail"))
	assertF16Problem(t, rec, apiproblem.CodeInternalError, http.StatusInternalServerError, "")
	p := decodeF16Problem(t, rec)
	if p.Status == http.StatusForbidden {
		t.Fatal("must not disclose suppressed 403")
	}
	if h.auditLen() != 0 {
		t.Fatal("failed append must leave no durable AuditEvent")
	}
}

func TestVS0_CF_F16_41(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	created := h.createTarget(t, "exp", "create-exp")
	h.setPresentFixture(created.Metadata.UID, f16AllSupportedTruths())
	rec := doF16(t, h, http.MethodPost, f16ActionPath(created.Metadata.UID, "qualify"), nil,
		f16ActionHeaders("q-exp", created.Metadata.ResourceVersion))
	if rec.Code != http.StatusOK {
		t.Fatalf("qualify=%d", rec.Code)
	}
	qualified := f16DecodeSafe(t, rec)
	etag := qualified.Metadata.ResourceVersion
	stored, _ := h.life.Store().GetExecutionTarget(created.Metadata.UID)
	fsUID := stored.Status.FactSetRef.UID
	beforeAudit := h.auditLen()
	h.clock.Advance(executiontarget.FactFreshnessWindow + time.Second)
	h.sched.ProcessExpiryTick(context.Background())
	if !h.life.FactSetExpired(fsUID) || !h.life.FactSetExpiryAudited(fsUID) {
		t.Fatal("expiry must mark freshness and audit once")
	}
	assertF16AuditDelta(t, h, beforeAudit, 1)
	get := doF16(t, h, http.MethodGet, f16ItemPath(created.Metadata.UID), nil, nil)
	proj := f16DecodeSafe(t, get)
	if proj.EffectiveAvailability != etmodel.EffectiveUnavailable {
		t.Fatalf("availability=%s", proj.EffectiveAvailability)
	}
	if get.Header().Get("ETag") != etag {
		t.Fatalf("ETag changed: %q vs %q", get.Header().Get("ETag"), etag)
	}
}

func TestVS0_CF_F16_42(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	created := h.createTarget(t, "mc-enter", "create-mc")
	h.setPresentFixture(created.Metadata.UID, f16AllSupportedTruths())
	_ = doF16(t, h, http.MethodPost, f16ActionPath(created.Metadata.UID, "qualify"), nil,
		f16ActionHeaders("q-mc", created.Metadata.ResourceVersion))
	et, _ := h.life.Store().GetExecutionTarget(created.Metadata.UID)
	beforeRV := et.Metadata.ResourceVersion
	beforeEpoch := et.Status.MaintenanceEpoch
	beforeAudit := h.auditLen()
	got := h.deliverMaintenance(t, et, executiontarget.MaintenanceTriggerEnter, "audit-enter-ok")
	if got.Status.MaintenanceEpoch != beforeEpoch+1 {
		t.Fatalf("epoch=%d", got.Status.MaintenanceEpoch)
	}
	if got.Metadata.ResourceVersion == beforeRV {
		t.Fatal("enter must advance resourceVersion")
	}
	if got.Status.FactSetRef != nil || got.Status.QualificationResultRef != nil {
		t.Fatal("enter must clear both links")
	}
	marker, ok := h.life.Store().GetMaintenanceMarker(et.Metadata.UID)
	if !ok || !marker.Active || marker.MaintenanceEpoch != got.Status.MaintenanceEpoch {
		t.Fatalf("marker=%#v ok=%v", marker, ok)
	}
	assertF16AuditDelta(t, h, beforeAudit, 1)
	get := doF16(t, h, http.MethodGet, f16ItemPath(created.Metadata.UID), nil, nil)
	proj := f16DecodeSafe(t, get)
	if proj.EffectiveAvailability != etmodel.EffectiveMaintenance {
		t.Fatalf("availability=%s", proj.EffectiveAvailability)
	}
}

func TestVS0_CF_F16_43(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	created := h.createTarget(t, "mc-clear", "create-mcc")
	et, _ := h.life.Store().GetExecutionTarget(created.Metadata.UID)
	entered := h.deliverMaintenance(t, et, executiontarget.MaintenanceTriggerEnter, "audit-enter-43")
	beforeAudit := h.auditLen()
	beforeRV := entered.Metadata.ResourceVersion
	got := h.deliverMaintenance(t, entered, executiontarget.MaintenanceTriggerClear, "audit-clear-ok")
	if got.Status.MaintenanceEpoch != entered.Status.MaintenanceEpoch+1 {
		t.Fatalf("clear epoch=%d", got.Status.MaintenanceEpoch)
	}
	if got.Metadata.ResourceVersion == beforeRV {
		t.Fatal("clear must advance RV")
	}
	if got.Status.Qualification != etmodel.QualificationUnqualified {
		t.Fatalf("qualification=%s", got.Status.Qualification)
	}
	if _, ok := h.life.Store().GetMaintenanceMarker(et.Metadata.UID); ok {
		t.Fatal("clear must remove marker")
	}
	assertF16AuditDelta(t, h, beforeAudit, 1)
	get := doF16(t, h, http.MethodGet, f16ItemPath(created.Metadata.UID), nil, nil)
	proj := f16DecodeSafe(t, get)
	if proj.Status.Lifecycle != etmodel.LifecycleActive ||
		proj.EffectiveAvailability != etmodel.EffectiveUnavailable {
		t.Fatalf("projection=%#v", proj)
	}
}

func TestVS0_CF_F16_44(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	created := h.createTarget(t, "head-item", "create-head")
	before := h.auditLen()
	rec := doF16(t, h, http.MethodHead, f16ItemPath(created.Metadata.UID), nil, nil)
	assertF16TransportOnly(t, rec, http.StatusMethodNotAllowed, "GET")
	assertF16NoAuditDelta(t, h, before)
}

func TestVS0_CF_F16_45(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	created := h.createTarget(t, "accept-ok", "create-accept")
	rec := doF16(t, h, http.MethodGet, f16ItemPath(created.Metadata.UID), nil, map[string]string{"Accept": "application/xml"})
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d", rec.Code)
	}
	assertF16JSONNosniff(t, rec)
}

func TestVS0_CF_F16_46(t *testing.T) {
	sharedAudit := &cloudmodel.MemoryAuditAppender{}
	h1 := newF16Harness(t, f16HarnessOpts{memoryAudit: sharedAudit, audit: sharedAudit})
	created := h1.createTarget(t, "restart", "create-restart")
	_ = doF16(t, h1, http.MethodPost, f16ActionPath(created.Metadata.UID, "qualify"), nil,
		f16ActionHeaders("q-restart", created.Metadata.ResourceVersion))
	if sharedAudit.Len() < 2 {
		t.Fatalf("want create+qualify audits, got %d", sharedAudit.Len())
	}
	beforeAudit := sharedAudit.Len()
	h2 := newF16Harness(t, f16HarnessOpts{memoryAudit: sharedAudit, audit: sharedAudit})
	if len(h2.life.Store().ListExecutionTargets()) != 0 {
		t.Fatal("restart must clear F0016-owned targets")
	}
	if sharedAudit.Len() != beforeAudit {
		t.Fatal("inherited AuditEvents must not be deleted")
	}
	if h2.life.Observer().ExternalCallCount() != 0 || h2.sched.ExternalCallCount() != 0 {
		t.Fatal("externalCallCount must be 0")
	}
}

func TestVS0_CF_F16_47(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	created := h.createTarget(t, "f47", "create-f47")
	headers := f16ActionHeaders("replay-safe", created.Metadata.ResourceVersion)
	rec := doF16(t, h, http.MethodPost, f16ActionPath(created.Metadata.UID, "qualify"), nil, headers)
	if rec.Code != http.StatusOK {
		t.Fatalf("qualify: %d", rec.Code)
	}
	h.qualify.Cloud = cloudmodel.NewStore()
	before := h.auditLen()
	rec2 := doF16(t, h, http.MethodPost, f16ActionPath(created.Metadata.UID, "qualify"), nil, headers)
	assertF16Problem(t, rec2, apiproblem.CodeResourceNotFound, http.StatusNotFound, f16ViolationSafeDeny)
	assertF16AuditDelta(t, h, before, 1)
}

func TestVS0_CF_F16_48(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	rec := doF16(t, h, http.MethodGet, f16ItemPath("dddddddddddddddddddddddddddddddd"), nil, map[string]string{"Accept": "text/html"})
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status=%d", rec.Code)
	}
	assertF16ProblemNosniff(t, rec)
}

func TestVS0_CF_F16_49(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	a := h.createTarget(t, "iso-a", "create-iso-a")
	stackB := h.seedExtraStack(t, "stack-iso-b", "b1b1b1b1b1b1b1b1b1b1b1b1b1b1b1b1", 1)
	b := h.createTargetOnStack(t, "iso-b", "create-iso-b", stackB.Metadata.UID)
	key := "shared-action-key"
	recA := doF16(t, h, http.MethodPost, f16ActionPath(a.Metadata.UID, "qualify"), nil,
		f16ActionHeaders(key, a.Metadata.ResourceVersion))
	if recA.Code != http.StatusOK {
		t.Fatalf("qualify a=%d", recA.Code)
	}
	recB := doF16(t, h, http.MethodPost, f16ActionPath(b.Metadata.UID, "qualify"), nil,
		f16ActionHeaders(key, b.Metadata.ResourceVersion))
	if recB.Code != http.StatusOK {
		t.Fatalf("qualify b=%d %s", recB.Code, recB.Body.String())
	}
	if recA.Body.String() == recB.Body.String() && f16DecodeSafe(t, recA).Metadata.UID == f16DecodeSafe(t, recB).Metadata.UID {
		t.Fatal("must not replay first target onto second")
	}
}

func TestVS0_CF_F16_50(t *testing.T) {
	hA := newF16Harness(t, f16HarnessOpts{})
	hA.seedActiveParticipationAndStack(t)
	key := "shared-create-key"
	recA := doF16(t, hA, http.MethodPost, f16Collection, f16CreateBody("scope-a", f16PartUID, f16StackUID), f16CreateHeaders(key))
	if recA.Code != http.StatusCreated {
		t.Fatalf("create a=%d", recA.Code)
	}
	hB := newF16Harness(t, f16HarnessOpts{grants: f16ExecutionTargetGrants(f16OtherProviderUID, true, true, true)})
	hB.seedPlatformAndProvider(t, f16OtherProviderUID)
	// Seed participation+stack for other provider.
	part, prob := hB.cloud.CreateParticipation(model.CloudProviderParticipation{
		TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionCloudProviderParticipation, Kind: model.KindCloudProviderParticipation},
		Metadata: apimeta.ObjectMeta{Name: "part-b", UID: "d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1d1", ResourceVersion: "1", Generation: 1,
			ScopeRef: &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{APIVersion: model.APIVersionCloudPlatform, Kind: string(apimeta.ScopeCloudPlatform), Name: "root-plat", UID: f16PlatformUID}}},
		Spec: model.CloudProviderParticipationSpec{
			CloudPlatformRef: apimeta.TypedRef{APIVersion: model.APIVersionCloudPlatform, Kind: model.KindCloudPlatform, Name: "root-plat", UID: f16PlatformUID},
			CloudProviderRef: apimeta.TypedRef{APIVersion: model.APIVersionCloudProvider, Kind: model.KindCloudProvider, Name: "prov-b", UID: f16OtherProviderUID},
			Environment:      model.ParticipationEnvironmentDevelopment,
		},
		Status: model.CloudProviderParticipationStatus{Phase: model.ParticipationPhaseActive},
	})
	if prob != nil {
		t.Fatalf("part: %#v", prob)
	}
	stack := hB.seedOtherProviderAndStack(t, "e1e1e1e1e1e1e1e1e1e1e1e1e1e1e1e1")
	recB := doF16(t, hB, http.MethodPost, f16Collection, f16CreateBody("scope-b", part.Metadata.UID, stack.Metadata.UID), f16CreateHeaders(key))
	if recB.Code != http.StatusCreated {
		t.Fatalf("create b=%d %s", recB.Code, recB.Body.String())
	}
	if f16DecodeSafe(t, recA).Metadata.UID == f16DecodeSafe(t, recB).Metadata.UID {
		t.Fatal("must not replay across CloudProvider scopes")
	}
}

func TestVS0_CF_F16_51(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	h.seedActiveParticipationAndStack(t)
	before := h.auditLen()
	rec := doF16(t, h, http.MethodPost, f16Collection, []byte(`{not-json`), f16CreateHeaders("k-51"))
	assertF16Problem(t, rec, apiproblem.CodeMalformedRequest, http.StatusBadRequest, "")
	assertF16NoAuditDelta(t, h, before)
}

func TestVS0_CF_F16_52(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	h.seedActiveParticipationAndStack(t)
	before := h.auditLen()
	body := []byte(`{"metadata":{"name":"t-52"},"metadata":{"name":"t-52b"},"spec":{"cloudProviderParticipationRef":{"uid":"` + f16PartUID + `"},"infrastructureStackRef":{"uid":"` + f16StackUID + `"},"targetClass":"synthetic-iaas"}}`)
	rec := doF16(t, h, http.MethodPost, f16Collection, body, f16CreateHeaders("k-52"))
	assertF16Problem(t, rec, apiproblem.CodeDuplicateField, http.StatusBadRequest, "")
	assertF16NoAuditDelta(t, h, before)
}

func TestVS0_CF_F16_53(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	h.seedActiveParticipationAndStack(t)
	before := h.auditLen()
	lim := apivalid.DefaultLimits()
	oversized := bytes.Repeat([]byte("a"), lim.MaxObjectBytes+1)
	rec := doF16(t, h, http.MethodPost, f16Collection, oversized, f16CreateHeaders("k-53"))
	assertF16Problem(t, rec, apiproblem.CodeRequestTooLarge, http.StatusBadRequest, "")
	assertF16NoAuditDelta(t, h, before)
}

func TestVS0_CF_F16_54(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	h.seedActiveParticipationAndStack(t)
	before := h.auditLen()
	body := []byte(`{"metadata":{"name":"t-54"},"spec":{"infrastructureStackRef":{"uid":"` + f16StackUID + `"},"targetClass":"synthetic-iaas"}}`)
	rec := doF16(t, h, http.MethodPost, f16Collection, body, f16CreateHeaders("k-54"))
	assertF16Problem(t, rec, apiproblem.CodeValidationFailed, http.StatusUnprocessableEntity, "")
	assertF16NoAuditDelta(t, h, before)
}

func TestVS0_CF_F16_55(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	h.seedActiveParticipationAndStack(t)
	before := h.auditLen()
	body := []byte(`{"metadata":{"name":"t-55"},"spec":{"cloudProviderParticipationRef":"not-an-object","infrastructureStackRef":{"uid":"` + f16StackUID + `"},"targetClass":"synthetic-iaas"}}`)
	rec := doF16(t, h, http.MethodPost, f16Collection, body, f16CreateHeaders("k-55"))
	assertF16Problem(t, rec, apiproblem.CodeValidationFailed, http.StatusUnprocessableEntity, "")
	assertF16NoAuditDelta(t, h, before)
}

func TestVS0_CF_F16_56(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	h.seedActiveParticipationAndStack(t)
	before := h.auditLen()
	body := []byte(`{"metadata":{"name":"t-56"},"spec":{"cloudProviderParticipationRef":{"uid":"` + f16PartUID + `"},"infrastructureStackRef":{"uid":"` + f16StackUID + `"},"targetClass":"other"}}`)
	rec := doF16(t, h, http.MethodPost, f16Collection, body, f16CreateHeaders("k-56"))
	assertF16Problem(t, rec, apiproblem.CodeValidationFailed, http.StatusUnprocessableEntity, "")
	assertF16NoAuditDelta(t, h, before)
}

func TestVS0_CF_F16_57(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	h.seedActiveParticipationAndStack(t)
	before := h.auditLen()
	body := []byte(`{"metadata":{"name":"t-57"},"spec":{"cloudProviderParticipationRef":{"uid":"` + f16PartUID + `"},"infrastructureStackRef":{"uid":"` + f16StackUID + `"},"targetClass":"synthetic-iaas","adapterAuthorityRef":{"uid":"x"}}}`)
	rec := doF16(t, h, http.MethodPost, f16Collection, body, f16CreateHeaders("k-57"))
	assertF16Problem(t, rec, apiproblem.CodeUnknownField, http.StatusBadRequest, "")
	assertF16NoAuditDelta(t, h, before)
}

func TestVS0_CF_F16_58(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	h.seedActiveParticipationAndStack(t)
	before := h.auditLen()
	rec := doF16(t, h, http.MethodPost, f16Collection, f16CreateBody("t-58", f16PartUID, f16StackUID), map[string]string{
		"Content-Type": f16MediaJSON, "Idempotency-Key": "   ",
	})
	assertF16Problem(t, rec, apiproblem.CodeMalformedRequest, http.StatusBadRequest, "")
	assertF16NoAuditDelta(t, h, before)
}

func TestVS0_CF_F16_59(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	h.seedActiveParticipationAndStack(t)
	before := h.auditLen()
	rec := doF16(t, h, http.MethodPost, f16Collection, f16CreateBody("t-59", f16PartUID, f16StackUID), map[string]string{
		"Content-Type": f16MediaJSON, "Idempotency-Key": "bad\x01key",
	})
	assertF16Problem(t, rec, apiproblem.CodeMalformedRequest, http.StatusBadRequest, "")
	assertF16NoAuditDelta(t, h, before)
}

func TestVS0_CF_F16_60(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	h.seedActiveParticipationAndStack(t)
	before := h.auditLen()
	rec := doF16(t, h, http.MethodPost, f16Collection, f16CreateBody("t-60", f16PartUID, f16StackUID), map[string]string{
		"Content-Type": f16MediaJSON, "Idempotency-Key": strings.Repeat("a", validate.MaxIdempotencyKeyLen+1),
	})
	assertF16Problem(t, rec, apiproblem.CodeMalformedRequest, http.StatusBadRequest, "")
	assertF16NoAuditDelta(t, h, before)
}

func TestVS0_CF_F16_61(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	before := h.auditLen()
	rec := doF16Unauth(t, h, http.MethodGet, f16Collection, nil, nil)
	assertF16Problem(t, rec, apiproblem.CodeAuthRequired, http.StatusUnauthorized, "")
	assertF16NoAuditDelta(t, h, before)
}

func TestVS0_CF_F16_62(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	before := h.auditLen()
	rec := doF16Unauth(t, h, http.MethodGet, f16ItemPath("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"), nil, nil)
	assertF16Problem(t, rec, apiproblem.CodeAuthRequired, http.StatusUnauthorized, "")
	assertF16NoAuditDelta(t, h, before)
}

func TestVS0_CF_F16_63(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	before := h.auditLen()
	rec := doF16Unauth(t, h, http.MethodPost, f16ActionPath("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "qualify"), nil, f16ActionHeaders("k", "1"))
	assertF16Problem(t, rec, apiproblem.CodeAuthRequired, http.StatusUnauthorized, "")
	assertF16NoAuditDelta(t, h, before)
}

func TestVS0_CF_F16_64(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	before := h.auditLen()
	rec := doF16Unauth(t, h, http.MethodPost, f16ActionPath("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "retire"), nil, f16ActionHeaders("k", "1"))
	assertF16Problem(t, rec, apiproblem.CodeAuthRequired, http.StatusUnauthorized, "")
	assertF16NoAuditDelta(t, h, before)
}

func TestVS0_CF_F16_65(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	h.seedActiveParticipationAndStack(t)
	stackT := h.seedExtraStack(t, "stack-tuple", "b1b1b1b1b1b1b1b1b1b1b1b1b1b1b1b2", 1)
	rec := doF16(t, h, http.MethodPost, f16Collection, f16CreateBody("tuple-a", f16PartUID, stackT.Metadata.UID), f16CreateHeaders("tuple-a"))
	if rec.Code != http.StatusCreated {
		t.Fatalf("first=%d", rec.Code)
	}
	before := h.auditLen()
	rec = doF16(t, h, http.MethodPost, f16Collection, f16CreateBody("tuple-b", f16PartUID, stackT.Metadata.UID), f16CreateHeaders("tuple-b"))
	assertF16Problem(t, rec, apiproblem.CodeAlreadyExists, http.StatusConflict, "")
	assertF16NoAuditDelta(t, h, before)
}

func TestVS0_CF_F16_66(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	h.seedActiveParticipationAndStack(t)
	before := h.auditLen()
	missing := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	rec := doF16(t, h, http.MethodPost, f16ActionPath(missing, "retire"), nil, f16ActionHeaders("missing-r", "1"))
	assertF16Problem(t, rec, apiproblem.CodeResourceNotFound, http.StatusNotFound, f16ViolationSafeDeny)
	assertF16AuditDelta(t, h, before, 1)
}

func TestVS0_CF_F16_67(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	created := h.createTarget(t, "dup-fact", "create-dup")
	h.setPresentFixture(created.Metadata.UID, f16AllSupportedTruths())
	rec := doF16(t, h, http.MethodPost, f16ActionPath(created.Metadata.UID, "qualify"), nil,
		f16ActionHeaders("q-ok", created.Metadata.ResourceVersion))
	qualified := f16DecodeSafe(t, rec)
	beforeRV := qualified.Metadata.ResourceVersion
	beforeAudit := h.auditLen()
	h.fixtures.Set(created.Metadata.UID, executiontarget.TargetBoundFixture{
		Revision: "dup", Mode: executiontarget.FixtureFaultDuplicate, Truths: f16AllSupportedTruths(),
	})
	rec2 := doF16(t, h, http.MethodPost, f16ActionPath(created.Metadata.UID, "qualify"), nil,
		f16ActionHeaders("q-dup", beforeRV))
	assertF16Problem(t, rec2, apiproblem.CodeInternalError, http.StatusInternalServerError, "")
	assertF16NoAuditDelta(t, h, beforeAudit)
	after, _ := h.life.Store().GetExecutionTarget(created.Metadata.UID)
	if after.Metadata.ResourceVersion != beforeRV {
		t.Fatal("must leave ETag unchanged")
	}
}

func TestVS0_CF_F16_68(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	created := h.createTarget(t, "miss-fact", "create-missf")
	h.setPresentFixture(created.Metadata.UID, f16AllSupportedTruths())
	rec := doF16(t, h, http.MethodPost, f16ActionPath(created.Metadata.UID, "qualify"), nil,
		f16ActionHeaders("q-ok", created.Metadata.ResourceVersion))
	qualified := f16DecodeSafe(t, rec)
	beforeRV := qualified.Metadata.ResourceVersion
	beforeAudit := h.auditLen()
	h.fixtures.Set(created.Metadata.UID, executiontarget.TargetBoundFixture{
		Revision: "miss", Mode: executiontarget.FixtureFaultMissing, Truths: f16AllSupportedTruths(),
		MissingName: executiontarget.FactNameNetworkPrivate,
	})
	rec2 := doF16(t, h, http.MethodPost, f16ActionPath(created.Metadata.UID, "qualify"), nil,
		f16ActionHeaders("q-miss", beforeRV))
	assertF16Problem(t, rec2, apiproblem.CodeInternalError, http.StatusInternalServerError, "")
	assertF16NoAuditDelta(t, h, beforeAudit)
	after, _ := h.life.Store().GetExecutionTarget(created.Metadata.UID)
	if after.Metadata.ResourceVersion != beforeRV {
		t.Fatal("must leave ETag unchanged")
	}
}

func TestVS0_CF_F16_69(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	created := h.createTarget(t, "mal-prov", "create-malp")
	beforeAudit := h.auditLen()
	beforeRV := created.Metadata.ResourceVersion
	h.fixtures.Set(created.Metadata.UID, executiontarget.TargetBoundFixture{
		Revision: "mal", Mode: executiontarget.FixtureFaultMalformed, Truths: f16AllSupportedTruths(),
	})
	rec := doF16(t, h, http.MethodPost, f16ActionPath(created.Metadata.UID, "qualify"), nil,
		f16ActionHeaders("q-malp", beforeRV))
	assertF16Problem(t, rec, apiproblem.CodeInternalError, http.StatusInternalServerError, "")
	assertF16NoAuditDelta(t, h, beforeAudit)
	after, _ := h.life.Store().GetExecutionTarget(created.Metadata.UID)
	if after.Metadata.ResourceVersion != beforeRV || after.Status.FactSetRef != nil {
		t.Fatal("must leave target unchanged")
	}
}

func TestVS0_CF_F16_70(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	created := h.createTarget(t, "if-mal", "create-ifmal")
	before := h.auditLen()
	rec := doF16(t, h, http.MethodPost, f16ActionPath(created.Metadata.UID, "qualify"), nil,
		f16ActionHeaders("malformed", "bad\x01"))
	assertF16Problem(t, rec, apiproblem.CodeStaleResourceVersion, http.StatusPreconditionFailed, "")
	assertF16NoAuditDelta(t, h, before)
}

func TestVS0_CF_F16_71(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	created := h.createTarget(t, "if-weak", "create-weak")
	before := h.auditLen()
	rec := doF16(t, h, http.MethodPost, f16ActionPath(created.Metadata.UID, "qualify"), nil,
		f16ActionHeaders("weak", `W/"1"`))
	assertF16Problem(t, rec, apiproblem.CodeStaleResourceVersion, http.StatusPreconditionFailed, "")
	assertF16NoAuditDelta(t, h, before)
}

func TestVS0_CF_F16_72(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	created := h.createTarget(t, "if-wild", "create-wild")
	before := h.auditLen()
	rec := doF16(t, h, http.MethodPost, f16ActionPath(created.Metadata.UID, "qualify"), nil,
		f16ActionHeaders("wild", "*"))
	assertF16Problem(t, rec, apiproblem.CodeStaleResourceVersion, http.StatusPreconditionFailed, "")
	assertF16NoAuditDelta(t, h, before)
}

func TestVS0_CF_F16_73(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	created := h.createTarget(t, "if-multi", "create-multi")
	before := h.auditLen()
	rec := doF16MultiIfMatch(t, h, f16ActionPath(created.Metadata.UID, "qualify"), "multi", "1", "2")
	assertF16Problem(t, rec, apiproblem.CodeStaleResourceVersion, http.StatusPreconditionFailed, "")
	assertF16NoAuditDelta(t, h, before)
}

func TestVS0_CF_F16_74(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	created := h.createTarget(t, "if-stale", "create-stale")
	before := h.auditLen()
	rec := doF16(t, h, http.MethodPost, f16ActionPath(created.Metadata.UID, "qualify"), nil,
		f16ActionHeaders("stale", "999"))
	assertF16Problem(t, rec, apiproblem.CodeStaleResourceVersion, http.StatusPreconditionFailed, "")
	assertF16NoAuditDelta(t, h, before)
}

func TestVS0_CF_F16_75(t *testing.T) {
	h, started, release := f16BlockingHarness(t)
	created := h.createTarget(t, "epoch-stale", "create-epoch")
	h.setPresentFixture(created.Metadata.UID, f16AllSupportedTruths())
	beforeAudit := h.auditLen()
	beforeRV := created.Metadata.ResourceVersion
	qualDone := make(chan *httptest.ResponseRecorder, 1)
	go func() {
		qualDone <- doF16(t, h, http.MethodPost, f16ActionPath(created.Metadata.UID, "qualify"), nil,
			f16ActionHeaders("q-epoch", created.Metadata.ResourceVersion))
	}()
	f16WaitStarted(t, started)
	h.bumpMaintenanceEpochNoMarker(created.Metadata.UID, 9)
	close(release)
	qRec := f16WaitRec(t, qualDone)
	assertF16Problem(t, qRec, apiproblem.CodeStaleResourceVersion, http.StatusPreconditionFailed, executiontarget.ViolationTargetEpochStale)
	assertF16NoAuditDelta(t, h, beforeAudit)
	after, _ := h.life.Store().GetExecutionTarget(created.Metadata.UID)
	if after.Status.FactSetRef != nil {
		t.Fatal("must not publish qualification refs")
	}
	if after.Metadata.ResourceVersion == beforeRV {
		t.Fatal("epoch bump should have advanced RV independently")
	}
}

func TestVS0_CF_F16_76(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	created := h.createTarget(t, "evict", "create-evict")
	const extra = 2
	total := executiontarget.MaxCompletedRecords + extra
	firstKey := "ret-" + f16PadIndex(0)
	for i := 0; i < total; i++ {
		key := "ret-" + f16PadIndex(i)
		if _, prob := f16QualifyViaLifecycle(t, h, created.Metadata.UID, key); prob != nil {
			t.Fatalf("qualify %d: %#v", i, prob)
		}
		h.clock.Advance(time.Millisecond)
	}
	// Earliest completed record unavailable → processed as new.
	nsFirst := f16QualifyNS(created.Metadata.UID, firstKey)
	res := h.life.ReserveOrInspectReplay(nsFirst, executiontarget.DigestAction())
	if res.Outcome != executiontarget.ReservationOutcomeReserved {
		t.Fatalf("earliest must be unavailable for replay; got %#v", res)
	}
	h.life.AbortReservation(nsFirst)

	inFlightKey := "inflight-keep"
	nsIF := f16QualifyNS(created.Metadata.UID, inFlightKey)
	if res := h.life.ReserveOrInspectReplay(nsIF, executiontarget.DigestAction()); res.Outcome != executiontarget.ReservationOutcomeReserved {
		t.Fatalf("inflight reserve: %#v", res)
	}
	for i := 0; i < 3; i++ {
		key := "more-" + f16PadIndex(i)
		if _, prob := f16QualifyViaLifecycle(t, h, created.Metadata.UID, key); prob != nil {
			t.Fatalf("more %d: %#v", i, prob)
		}
		h.clock.Advance(time.Millisecond)
	}
	// InFlight must still be attachable.
	if _, ok := h.life.AttachWaiter(nsIF); !ok {
		t.Fatal("InFlight must never be evicted")
	}
	h.life.AbortReservation(nsIF)
}

func TestVS0_CF_F16_77(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	created := h.createTarget(t, "shutdown", "create-sd")
	nsCreate := f16CreateNS("k-sd-inflight")
	nsQualify := f16QualifyNS(created.Metadata.UID, "k-sd-q")
	nsRetire := f16RetireNS(created.Metadata.UID, "k-sd-r")
	for _, ns := range []executiontarget.IdempotencyNamespace{nsCreate, nsQualify, nsRetire} {
		if res := h.life.ReserveOrInspectReplay(ns, executiontarget.DigestAction()); res.Outcome != executiontarget.ReservationOutcomeReserved {
			t.Fatalf("reserve %#v: %#v", ns, res)
		}
	}
	done, ok := h.life.AttachWaiter(nsCreate)
	if !ok {
		t.Fatal("attach")
	}
	woke := make(chan struct{})
	go func() {
		<-done
		close(woke)
	}()
	before, _ := h.life.Store().GetExecutionTarget(created.Metadata.UID)
	beforeRV := before.Metadata.ResourceVersion
	h.sched.StopAcceptance()
	et, _ := h.life.Store().GetExecutionTarget(created.Metadata.UID)
	_, prob := h.sched.DeliverMaintenanceTrigger(context.Background(), h.maintenanceTrigger(et, executiontarget.MaintenanceTriggerEnter, "audit-post-stop"))
	if prob == nil {
		t.Fatal("post-stop trigger must not succeed")
	}
	h.life.Shutdown()
	if !h.life.AdmissionClosed() {
		t.Fatal("admission must be closed")
	}
	select {
	case <-woke:
	case <-time.After(2 * time.Second):
		t.Fatal("shutdown must wake waiters")
	}
	after, _ := h.life.Store().GetExecutionTarget(created.Metadata.UID)
	if after.Metadata.ResourceVersion != beforeRV {
		t.Fatal("last committed projection must remain unchanged")
	}
	res := h.life.ReserveOrInspectReplay(f16CreateNS("k-sd-after"), executiontarget.DigestAction())
	if res.Outcome == executiontarget.ReservationOutcomeReserved {
		t.Fatal("no new InFlight after shutdown")
	}
}

func TestVS0_CF_F16_78(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	created := h.createTarget(t, "cancel", "create-cancel")
	ns := f16QualifyNS(created.Metadata.UID, "shared-key")
	if res := h.life.ReserveOrInspectReplay(ns, executiontarget.DigestAction()); res.Outcome != executiontarget.ReservationOutcomeReserved {
		t.Fatalf("owner reserve: %#v", res)
	}
	var wg sync.WaitGroup
	wg.Add(1)
	var waiterCode int
	go func() {
		defer wg.Done()
		rec := doF16(t, h, http.MethodPost, f16ActionPath(created.Metadata.UID, "qualify"), nil,
			f16ActionHeaders("shared-key", created.Metadata.ResourceVersion))
		waiterCode = rec.Code
	}()
	time.Sleep(30 * time.Millisecond)
	if !h.life.AbortReservation(ns) {
		t.Fatal("abort owner reservation")
	}
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	select {
	case <-done:
		if waiterCode == 0 {
			t.Fatal("empty status")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("waiter did not wake after abort")
	}
}

func TestVS0_CF_F16_79(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	created := h.createTarget(t, "retry", "create-retry")
	h.setPresentFixture(created.Metadata.UID, f16AllSupportedTruths())
	rec := doF16(t, h, http.MethodPost, f16ActionPath(created.Metadata.UID, "qualify"), nil,
		f16ActionHeaders("q-retry", created.Metadata.ResourceVersion))
	if rec.Code != http.StatusOK {
		t.Fatalf("qualify=%d", rec.Code)
	}
	qualified := f16DecodeSafe(t, rec)
	etag := qualified.Metadata.ResourceVersion
	stored, _ := h.life.Store().GetExecutionTarget(created.Metadata.UID)
	fsUID := stored.Status.FactSetRef.UID
	h.clock.Advance(executiontarget.FactFreshnessWindow + time.Second)
	h.setAuditFail(errors.New("append boom"))
	h.sched.ProcessExpiryTick(context.Background())
	if h.life.FactSetExpiryAudited(fsUID) {
		t.Fatal("failed append must not mark audited")
	}
	if !h.life.FactSetExpired(fsUID) {
		t.Fatal("freshness must still flip")
	}
	after, _ := h.life.Store().GetExecutionTarget(created.Metadata.UID)
	if after.Metadata.ResourceVersion != etag {
		t.Fatal("ETag must stay unchanged")
	}
	h.clock.Advance(time.Second)
	h.sched.ProcessExpiryTick(context.Background())
	h.setAuditFail(nil)
	h.clock.Advance(time.Second)
	beforeAudit := h.auditLen()
	h.sched.ProcessExpiryTick(context.Background())
	if !h.life.FactSetExpiryAudited(fsUID) {
		t.Fatal("successful retry must mark audited")
	}
	assertF16AuditDelta(t, h, beforeAudit, 1)
}

func TestVS0_CF_F16_80(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	created := h.createTarget(t, "stale-enter", "create-se")
	et, _ := h.life.Store().GetExecutionTarget(created.Metadata.UID)
	beforeRV := et.Metadata.ResourceVersion
	beforeAudit := h.auditLen()
	trig := h.maintenanceTrigger(et, executiontarget.MaintenanceTriggerEnter, "audit-stale-enter")
	trig.ExpectedMaintenanceEpoch = et.Status.MaintenanceEpoch + 1
	got, prob := h.sched.DeliverMaintenanceTrigger(context.Background(), trig)
	if prob != nil {
		t.Fatalf("stale enter must be silent no-op, got %#v", prob)
	}
	if got.Metadata.ResourceVersion != beforeRV {
		t.Fatal("stale enter must not advance RV")
	}
	assertF16NoAuditDelta(t, h, beforeAudit)
	if _, ok := h.life.Store().GetMaintenanceMarker(et.Metadata.UID); ok {
		t.Fatal("stale enter must not set marker")
	}
}

func TestVS0_CF_F16_81(t *testing.T) {
	h, started, release := f16BlockingHarness(t)
	created := h.createTarget(t, "waiter81", "create-81")
	h.setPresentFixture(created.Metadata.UID, f16AllSupportedTruths())
	path := f16ActionPath(created.Metadata.UID, "qualify")
	ownerDone := make(chan *httptest.ResponseRecorder, 1)
	go func() {
		ownerDone <- doF16(t, h, http.MethodPost, path, nil, f16ActionHeaders("shared-81", created.Metadata.ResourceVersion))
	}()
	f16WaitStarted(t, started)
	waiterDone := make(chan *httptest.ResponseRecorder, 1)
	go func() {
		time.Sleep(20 * time.Millisecond)
		waiterDone <- doF16(t, h, http.MethodPost, path, nil, f16ActionHeaders("shared-81", created.Metadata.ResourceVersion))
	}()
	time.Sleep(30 * time.Millisecond)
	h.updateStackGeneration(t, f16StackUID, 99)
	close(release)
	ownerRec := f16WaitRec(t, ownerDone)
	assertF16Problem(t, ownerRec, apiproblem.CodeStaleResourceVersion, http.StatusPreconditionFailed, executiontarget.ViolationTargetEpochStale)
	waiterRec := f16WaitRec(t, waiterDone)
	// No cached completion from the aborted owner: waiter re-enters and may
	// establish a new reservation (often 200 with the bumped generation).
	if waiterRec.Code == 0 {
		t.Fatal("waiter must wake and produce a response")
	}
	if waiterRec.Code == http.StatusOK {
		return
	}
	// Alternatively the waiter may also observe STALE if fences remain mismatched.
	if waiterRec.Code != http.StatusPreconditionFailed {
		t.Fatalf("waiter unexpected status=%d body=%s", waiterRec.Code, waiterRec.Body.String())
	}
}

func TestVS0_CF_F16_82(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	created := h.createTarget(t, "q-audit-fail", "create-qaf")
	h.setPresentFixture(created.Metadata.UID, f16AllSupportedTruths())
	beforeRV := created.Metadata.ResourceVersion
	beforeAudit := h.auditLen()
	h.setAuditFail(errors.New("append failed"))
	rec := doF16(t, h, http.MethodPost, f16ActionPath(created.Metadata.UID, "qualify"), nil,
		f16ActionHeaders("q-af", created.Metadata.ResourceVersion))
	assertF16Problem(t, rec, apiproblem.CodeInternalError, http.StatusInternalServerError, "")
	assertF16NoAuditDelta(t, h, beforeAudit)
	after, _ := h.life.Store().GetExecutionTarget(created.Metadata.UID)
	if after.Metadata.ResourceVersion != beforeRV || after.Status.FactSetRef != nil {
		t.Fatal("failed append must leave target/ETag unmutated")
	}
}

func TestVS0_CF_F16_83(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	created := h.createTarget(t, "stale-clear", "create-sc")
	et, _ := h.life.Store().GetExecutionTarget(created.Metadata.UID)
	entered := h.deliverMaintenance(t, et, executiontarget.MaintenanceTriggerEnter, "audit-enter-83")
	beforeAudit := h.auditLen()
	beforeRV := entered.Metadata.ResourceVersion
	trig := h.maintenanceTrigger(entered, executiontarget.MaintenanceTriggerClear, "audit-stale-clear")
	trig.ExpectedInfrastructureStackGen = entered.Status.ObservedGeneration + 9
	got, prob := h.sched.DeliverMaintenanceTrigger(context.Background(), trig)
	if prob != nil {
		t.Fatalf("stale clear: %#v", prob)
	}
	if got.Metadata.ResourceVersion != beforeRV {
		t.Fatal("stale clear must not mutate")
	}
	assertF16NoAuditDelta(t, h, beforeAudit)
}

func TestVS0_CF_F16_84(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	h.seedActiveParticipationAndStack(t)
	before := h.auditLen()
	body := []byte(`{"metadata":{"name":"secret-reject"},"spec":{"cloudProviderParticipationRef":{"uid":"` + f16PartUID + `"},"infrastructureStackRef":{"uid":"` + f16StackUID + `"},"targetClass":"synthetic-iaas","secretRef":{"name":"x"}}}`)
	rec := doF16(t, h, http.MethodPost, f16Collection, body, f16CreateHeaders("secret-reject"))
	assertF16Problem(t, rec, apiproblem.CodeUnknownField, http.StatusBadRequest, "")
	assertF16NoAuditDelta(t, h, before)
}

func TestVS0_CF_F16_85(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	created := h.createTarget(t, "stale-maint", "create-sm")
	et, _ := h.life.Store().GetExecutionTarget(created.Metadata.UID)
	_ = h.deliverMaintenance(t, et, executiontarget.MaintenanceTriggerEnter, "audit-sm")
	rec := doF16(t, h, http.MethodPost, f16ActionPath(created.Metadata.UID, "qualify"), nil,
		f16ActionHeaders("stale-before-maint", created.Metadata.ResourceVersion))
	assertF16Problem(t, rec, apiproblem.CodeStaleResourceVersion, http.StatusPreconditionFailed, "")
}

func TestVS0_CF_F16_86(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	created := h.createTarget(t, "hdr-body", "create-hb")
	before := h.auditLen()
	rec := doF16(t, h, http.MethodPost, f16ActionPath(created.Metadata.UID, "qualify"), []byte(`{"x":1}`),
		f16ActionHeaders("malformed", "bad\x01"))
	assertF16Problem(t, rec, apiproblem.CodeStaleResourceVersion, http.StatusPreconditionFailed, "")
	assertF16NoAuditDelta(t, h, before)
}

func TestVS0_CF_F16_87(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	created := h.createTarget(t, "replay-etag", "create-re")
	h.setPresentFixture(created.Metadata.UID, f16AllSupportedTruths())
	headers := f16ActionHeaders("qualify-etag", created.Metadata.ResourceVersion)
	rec := doF16(t, h, http.MethodPost, f16ActionPath(created.Metadata.UID, "qualify"), nil, headers)
	if rec.Code != http.StatusOK {
		t.Fatalf("qualify=%d", rec.Code)
	}
	// Change ETag via maintenance enter; replay consults completed record before
	// If-Match / Maintenance evaluation (If-Match exclusion).
	et, _ := h.life.Store().GetExecutionTarget(created.Metadata.UID)
	_ = h.deliverMaintenance(t, et, executiontarget.MaintenanceTriggerEnter, "audit-etag-bump")
	before := h.auditLen()
	rec2 := doF16(t, h, http.MethodPost, f16ActionPath(created.Metadata.UID, "qualify"), nil, headers)
	if rec2.Code != http.StatusOK {
		t.Fatalf("replay after ETag change=%d %s", rec2.Code, rec2.Body.String())
	}
	if rec2.Body.String() != rec.Body.String() {
		t.Fatal("must return original successful body")
	}
	assertF16NoAuditDelta(t, h, before)
}

func TestVS0_CF_F16_88(t *testing.T) {
	h, started, release := f16BlockingHarness(t)
	created := h.createTarget(t, "retire-wins", "create-rw")
	beforeAudit := h.auditLen()
	qualDone := make(chan *httptest.ResponseRecorder, 1)
	go func() {
		qualDone <- doF16(t, h, http.MethodPost, f16ActionPath(created.Metadata.UID, "qualify"), nil,
			f16ActionHeaders("q-inflight", created.Metadata.ResourceVersion))
	}()
	f16WaitStarted(t, started)
	retireRec := doF16(t, h, http.MethodPost, f16ActionPath(created.Metadata.UID, "retire"), nil,
		f16ActionHeaders("retire-win", created.Metadata.ResourceVersion))
	if retireRec.Code != http.StatusOK {
		t.Fatalf("retire: %d %s", retireRec.Code, retireRec.Body.String())
	}
	close(release)
	qRec := f16WaitRec(t, qualDone)
	assertF16Problem(t, qRec, apiproblem.CodeConflict, http.StatusConflict, executiontarget.ViolationTargetRetired)
	if h.auditLen() != beforeAudit+1 {
		t.Fatalf("audit len=%d want create+retire only (%d)", h.auditLen(), beforeAudit+1)
	}
}

func TestVS0_CF_F16_89(t *testing.T) {
	h, started, release := f16BlockingHarness(t)
	created := h.createTarget(t, "maint-wins", "create-mw")
	beforeAudit := h.auditLen()
	qualDone := make(chan *httptest.ResponseRecorder, 1)
	go func() {
		qualDone <- doF16(t, h, http.MethodPost, f16ActionPath(created.Metadata.UID, "qualify"), nil,
			f16ActionHeaders("q-maint-inflight", created.Metadata.ResourceVersion))
	}()
	f16WaitStarted(t, started)
	et, _ := h.life.Store().GetExecutionTarget(created.Metadata.UID)
	_ = h.deliverMaintenance(t, et, executiontarget.MaintenanceTriggerEnter, "audit-maint-wins")
	close(release)
	qRec := f16WaitRec(t, qualDone)
	assertF16Problem(t, qRec, apiproblem.CodeConflict, http.StatusConflict, executiontarget.ViolationTargetMaintenance)
	if h.auditLen() != beforeAudit+1 {
		t.Fatalf("audit len=%d want create+maintenance only", h.auditLen())
	}
}

func TestVS0_CF_F16_90(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	created := h.createTarget(t, "timeout", "create-to")
	h.fixtures.Set(created.Metadata.UID, executiontarget.TargetBoundFixture{
		Revision: "to-1", Mode: executiontarget.FixtureLogicalTimeout,
	})
	rec := doF16(t, h, http.MethodPost, f16ActionPath(created.Metadata.UID, "qualify"), nil,
		f16ActionHeaders("q-to", created.Metadata.ResourceVersion))
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d %s", rec.Code, rec.Body.String())
	}
	proj := f16DecodeSafe(t, rec)
	if proj.Status.Qualification != etmodel.QualificationIndeterminate ||
		proj.EffectiveAvailability != etmodel.EffectiveUnavailable {
		t.Fatalf("projection=%#v", proj)
	}
	stored, _ := h.life.Store().GetExecutionTarget(created.Metadata.UID)
	fs, _ := h.life.Store().GetFactSet(stored.Status.FactSetRef.UID)
	if fs.ObserverRevision != "to-1" {
		t.Fatalf("observerRevision=%q", fs.ObserverRevision)
	}
}

func TestVS0_CF_F16_91(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	h.seedActiveParticipationAndStack(t)
	before := h.auditLen()
	body := []byte(`{"metadata":{"name":"t-91","resourceVersion":"9"},"spec":{"cloudProviderParticipationRef":{"uid":"` + f16PartUID + `"},"infrastructureStackRef":{"uid":"` + f16StackUID + `"},"targetClass":"synthetic-iaas"}}`)
	rec := doF16(t, h, http.MethodPost, f16Collection, body, f16CreateHeaders("k-91"))
	assertF16Problem(t, rec, apiproblem.CodeAuthorizationDenied, http.StatusForbidden, validate.ViolationSystemOwnedFieldWrite)
	assertF16AuditDelta(t, h, before, 1)
}

func TestVS0_CF_F16_92(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	h.seedActiveParticipationAndStack(t)
	before := h.auditLen()
	body := []byte(`{"metadata":{"name":"t-92","generation":9},"spec":{"cloudProviderParticipationRef":{"uid":"` + f16PartUID + `"},"infrastructureStackRef":{"uid":"` + f16StackUID + `"},"targetClass":"synthetic-iaas"}}`)
	rec := doF16(t, h, http.MethodPost, f16Collection, body, f16CreateHeaders("k-92"))
	assertF16Problem(t, rec, apiproblem.CodeAuthorizationDenied, http.StatusForbidden, validate.ViolationSystemOwnedFieldWrite)
	assertF16AuditDelta(t, h, before, 1)
}

func TestVS0_CF_F16_93(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	h.seedActiveParticipationAndStack(t)
	before := h.auditLen()
	rec := doF16InvalidAuth(t, h, http.MethodPost, f16Collection, f16CreateBody("t-93", f16PartUID, f16StackUID), f16CreateHeaders("k-93"))
	assertF16Problem(t, rec, apiproblem.CodeAuthRequired, http.StatusUnauthorized, "")
	assertF16NoAuditDelta(t, h, before)
}

func TestVS0_CF_F16_94(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	before := h.auditLen()
	rec := doF16InvalidAuth(t, h, http.MethodGet, f16Collection, nil, nil)
	assertF16Problem(t, rec, apiproblem.CodeAuthRequired, http.StatusUnauthorized, "")
	assertF16NoAuditDelta(t, h, before)
}

func TestVS0_CF_F16_95(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	before := h.auditLen()
	rec := doF16InvalidAuth(t, h, http.MethodGet, f16ItemPath("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"), nil, nil)
	assertF16Problem(t, rec, apiproblem.CodeAuthRequired, http.StatusUnauthorized, "")
	assertF16NoAuditDelta(t, h, before)
}

func TestVS0_CF_F16_96(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	before := h.auditLen()
	rec := doF16InvalidAuth(t, h, http.MethodPost, f16ActionPath("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "qualify"), nil, f16ActionHeaders("k", "1"))
	assertF16Problem(t, rec, apiproblem.CodeAuthRequired, http.StatusUnauthorized, "")
	assertF16NoAuditDelta(t, h, before)
}

func TestVS0_CF_F16_97(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	before := h.auditLen()
	rec := doF16InvalidAuth(t, h, http.MethodPost, f16ActionPath("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "retire"), nil, f16ActionHeaders("k", "1"))
	assertF16Problem(t, rec, apiproblem.CodeAuthRequired, http.StatusUnauthorized, "")
	assertF16NoAuditDelta(t, h, before)
}

func TestVS0_CF_F16_98(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	created := h.createTarget(t, "name-keep", "create-nk")
	rec := doF16(t, h, http.MethodPost, f16ActionPath(created.Metadata.UID, "retire"), nil,
		f16ActionHeaders("retire-nk", created.Metadata.ResourceVersion))
	if rec.Code != http.StatusOK {
		t.Fatalf("retire=%d", rec.Code)
	}
	stack2 := h.seedExtraStack(t, "stack-nk2", "c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2", 1)
	before := h.auditLen()
	rec2 := doF16(t, h, http.MethodPost, f16Collection, f16CreateBody("name-keep", f16PartUID, stack2.Metadata.UID), f16CreateHeaders("reuse-name"))
	assertF16Problem(t, rec2, apiproblem.CodeAlreadyExists, http.StatusConflict, "")
	assertF16NoAuditDelta(t, h, before)
}

func TestVS0_CF_F16_99(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	created := h.createTarget(t, "retire-fail", "create-rf")
	beforeRV := created.Metadata.ResourceVersion
	h.setAuditFail(errors.New("retire append boom"))
	rec := doF16(t, h, http.MethodPost, f16ActionPath(created.Metadata.UID, "retire"), nil,
		f16ActionHeaders("retire-fail", created.Metadata.ResourceVersion))
	assertF16Problem(t, rec, apiproblem.CodeInternalError, http.StatusInternalServerError, "")
	after, _ := h.life.Store().GetExecutionTarget(created.Metadata.UID)
	if after.Status.Lifecycle != etmodel.LifecycleActive || after.Metadata.ResourceVersion != beforeRV {
		t.Fatal("failed retire must not publish")
	}
}

func TestVS0_CF_F16_100(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	created := h.createTarget(t, "maint-fail", "create-mf")
	et, _ := h.life.Store().GetExecutionTarget(created.Metadata.UID)
	beforeRV := et.Metadata.ResourceVersion
	beforeEpoch := et.Status.MaintenanceEpoch
	h.setAuditFail(errors.New("maint append boom"))
	_, prob := h.sched.DeliverMaintenanceTrigger(context.Background(), h.maintenanceTrigger(et, executiontarget.MaintenanceTriggerEnter, "audit-100"))
	if prob == nil || prob.Code != apiproblem.CodeInternalError {
		t.Fatalf("want INTERNAL_ERROR, got %#v", prob)
	}
	after, _ := h.life.Store().GetExecutionTarget(created.Metadata.UID)
	if after.Metadata.ResourceVersion != beforeRV || after.Status.MaintenanceEpoch != beforeEpoch {
		t.Fatal("failed enter must not advance epoch/RV")
	}
	if _, ok := h.life.Store().GetMaintenanceMarker(et.Metadata.UID); ok {
		t.Fatal("must not set marker")
	}
}

func TestVS0_CF_F16_101(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	created := h.createTarget(t, "clear-fail", "create-cf")
	et, _ := h.life.Store().GetExecutionTarget(created.Metadata.UID)
	entered := h.deliverMaintenance(t, et, executiontarget.MaintenanceTriggerEnter, "audit-enter-101")
	beforeRV := entered.Metadata.ResourceVersion
	beforeEpoch := entered.Status.MaintenanceEpoch
	h.setAuditFail(errors.New("clear append boom"))
	_, prob := h.sched.DeliverMaintenanceTrigger(context.Background(), h.maintenanceTrigger(entered, executiontarget.MaintenanceTriggerClear, "audit-101"))
	if prob == nil || prob.Code != apiproblem.CodeInternalError {
		t.Fatalf("want INTERNAL_ERROR, got %#v", prob)
	}
	after, _ := h.life.Store().GetExecutionTarget(created.Metadata.UID)
	if after.Metadata.ResourceVersion != beforeRV || after.Status.MaintenanceEpoch != beforeEpoch {
		t.Fatal("failed clear must not advance")
	}
	if _, ok := h.life.Store().GetMaintenanceMarker(et.Metadata.UID); !ok {
		t.Fatal("marker must remain")
	}
}

func TestVS0_CF_F16_102(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	h.setAuditFail(errors.New("append boom"))
	rec := doF16(t, h, http.MethodGet, f16ItemPath("cccccccccccccccccccccccccccccccc"), nil, map[string]string{"Accept": "text/plain"})
	f16AssertInternalNoDisclosure(t, rec, http.StatusNotFound)
	if h.auditLen() != 0 {
		t.Fatal("failed append must leave no durable AuditEvent")
	}
}

func TestVS0_CF_F16_103(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	created := h.createTarget(t, "f103", "create-f103")
	headers := f16ActionHeaders("replay-authz", created.Metadata.ResourceVersion)
	rec := doF16(t, h, http.MethodPost, f16ActionPath(created.Metadata.UID, "qualify"), nil, headers)
	if rec.Code != http.StatusOK {
		t.Fatalf("qualify: %d", rec.Code)
	}
	h.rebindGrants(f16ExecutionTargetGrants(f16ProviderUID, true, true, false))
	h.setAuditFail(errors.New("append boom"))
	rec2 := doF16(t, h, http.MethodPost, f16ActionPath(created.Metadata.UID, "qualify"), nil, headers)
	f16AssertInternalNoDisclosure(t, rec2, http.StatusForbidden)
}

func TestVS0_CF_F16_104(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	created := h.createTarget(t, "f104", "create-f104")
	headers := f16ActionHeaders("replay-safe", created.Metadata.ResourceVersion)
	rec := doF16(t, h, http.MethodPost, f16ActionPath(created.Metadata.UID, "qualify"), nil, headers)
	if rec.Code != http.StatusOK {
		t.Fatalf("qualify: %d", rec.Code)
	}
	h.qualify.Cloud = cloudmodel.NewStore()
	h.setAuditFail(errors.New("append boom"))
	rec2 := doF16(t, h, http.MethodPost, f16ActionPath(created.Metadata.UID, "qualify"), nil, headers)
	f16AssertInternalNoDisclosure(t, rec2, http.StatusNotFound)
}

func TestVS0_CF_F16_105(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	h.seedActiveParticipationAndStack(t)
	before := h.auditLen()
	rec := doF16(t, h, http.MethodPost, f16Collection, f16CreateBody("t-105", f16PartUID, f16StackUID), map[string]string{
		"Content-Type": f16MediaJSON, "Idempotency-Key": "ifmatch", "If-Match": "1",
	})
	assertF16Problem(t, rec, apiproblem.CodeMalformedRequest, http.StatusBadRequest, "")
	assertF16NoAuditDelta(t, h, before)
}

func TestVS0_CF_F16_106(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	created := h.createTarget(t, "query-fam", "create-qf")
	before := h.auditLen()
	members := []struct {
		method string
		path   string
		body   []byte
		hdrs   map[string]string
	}{
		{http.MethodPost, f16Collection + "?x=1", f16CreateBody("t", f16PartUID, f16StackUID), f16CreateHeaders("q-post")},
		{http.MethodGet, f16Collection + "?x=1", nil, nil},
		{http.MethodGet, f16ItemPath(created.Metadata.UID) + "?x=1", nil, nil},
		{http.MethodPost, f16ActionPath(created.Metadata.UID, "qualify") + "?x=1", nil, f16ActionHeaders("q", created.Metadata.ResourceVersion)},
		{http.MethodPost, f16ActionPath(created.Metadata.UID, "retire") + "?x=1", nil, f16ActionHeaders("r", created.Metadata.ResourceVersion)},
	}
	for _, m := range members {
		rec := doF16(t, h, m.method, m.path, m.body, m.hdrs)
		assertF16Problem(t, rec, apiproblem.CodeMalformedRequest, http.StatusBadRequest, "")
	}
	assertF16NoAuditDelta(t, h, before)
}

func TestVS0_CF_F16_107(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	before := h.auditLen()
	for _, bad := range []string{"not-a-uid", "ABCDEFABCDEFABCDEFABCDEFABCDEFAB", "short", "zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz"} {
		rec := doF16(t, h, http.MethodGet, f16ItemPath(bad), nil, nil)
		assertF16Problem(t, rec, apiproblem.CodeMalformedRequest, http.StatusBadRequest, "")
	}
	assertF16NoAuditDelta(t, h, before)
}

func TestVS0_CF_F16_108(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	created := h.createTarget(t, "slash-item", "create-si")
	before := h.auditLen()
	rec := doF16(t, h, http.MethodGet, f16ItemPath(created.Metadata.UID)+"/", nil, nil)
	assertF16TransportOnly(t, rec, http.StatusNotFound, "")
	assertF16NoAuditDelta(t, h, before)
}

func TestVS0_CF_F16_109(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	before := h.auditLen()
	rec := doF16(t, h, http.MethodPost, f16ActionPath("not-a-uid", "qualify"), nil, f16ActionHeaders("baduid", "1"))
	assertF16Problem(t, rec, apiproblem.CodeMalformedRequest, http.StatusBadRequest, "")
	assertF16NoAuditDelta(t, h, before)
}

func TestVS0_CF_F16_110(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	before := h.auditLen()
	rec := doF16(t, h, http.MethodPost, f16ActionPath("not-a-uid", "retire"), nil, f16ActionHeaders("baduid", "1"))
	assertF16Problem(t, rec, apiproblem.CodeMalformedRequest, http.StatusBadRequest, "")
	assertF16NoAuditDelta(t, h, before)
}

func TestVS0_CF_F16_111(t *testing.T) {
	hWrite := newF16Harness(t, f16HarnessOpts{grants: f16ExecutionTargetGrants(f16ProviderUID, true, false, false)})
	created := hWrite.createTarget(t, "item-no-read", "item-no-read-create")
	failAudit := &cloudmodel.MemoryAuditAppender{}
	h := newF16Harness(t, f16HarnessOpts{
		grants:      f16ExecutionTargetGrants(f16ProviderUID, true, false, false),
		memoryAudit: failAudit,
		audit:       failAudit,
	})
	h.life, h.cloud, h.fixtures = hWrite.life, hWrite.cloud, hWrite.fixtures
	h.rebindGrants(f16ExecutionTargetGrants(f16ProviderUID, true, false, false))
	rec := doF16(t, h, http.MethodGet, f16ItemPath(created.Metadata.UID), nil, nil)
	assertF16Problem(t, rec, apiproblem.CodeAuthorizationDenied, http.StatusForbidden, "")
	if failAudit.Len() != 1 {
		t.Fatalf("authz denial audit=%d", failAudit.Len())
	}
}

func TestVS0_CF_F16_112(t *testing.T) {
	// Equivalence family: every named member returns 500 INTERNAL_ERROR without disclosing 403/404.
	t.Run("LIST without read", func(t *testing.T) {
		h := newF16Harness(t, f16HarnessOpts{grants: f16ExecutionTargetGrants(f16ProviderUID, true, false, true)})
		h.setAuditFail(errors.New("append boom"))
		rec := doF16(t, h, http.MethodGet, f16Collection, nil, nil)
		f16AssertInternalNoDisclosure(t, rec, http.StatusForbidden)
		if h.auditLen() != 0 {
			t.Fatal("no durable AuditEvent")
		}
	})
	t.Run("direct GET without read", func(t *testing.T) {
		hWrite := newF16Harness(t, f16HarnessOpts{})
		created := hWrite.createTarget(t, "item-112", "item-112-create")
		failAudit := &cloudmodel.MemoryAuditAppender{}
		failAudit.SetFail(errors.New("append boom"))
		h := newF16Harness(t, f16HarnessOpts{
			grants:      f16ExecutionTargetGrants(f16ProviderUID, true, false, false),
			memoryAudit: failAudit,
			audit:       failAudit,
		})
		h.life, h.cloud, h.fixtures = hWrite.life, hWrite.cloud, hWrite.fixtures
		h.rebindGrants(f16ExecutionTargetGrants(f16ProviderUID, true, false, false))
		rec := doF16(t, h, http.MethodGet, f16ItemPath(created.Metadata.UID), nil, nil)
		f16AssertInternalNoDisclosure(t, rec, http.StatusForbidden)
		if failAudit.Len() != 0 {
			t.Fatal("no durable AuditEvent")
		}
	})
	t.Run("qualify without qualify grant", func(t *testing.T) {
		h := newF16Harness(t, f16HarnessOpts{grants: f16ExecutionTargetGrants(f16ProviderUID, true, true, false)})
		created := h.createTarget(t, "q112", "create-q112")
		h.setAuditFail(errors.New("append boom"))
		rec := doF16(t, h, http.MethodPost, f16ActionPath(created.Metadata.UID, "qualify"), nil,
			f16ActionHeaders("q112", created.Metadata.ResourceVersion))
		f16AssertInternalNoDisclosure(t, rec, http.StatusForbidden)
	})
	t.Run("retire without write grant", func(t *testing.T) {
		hWrite := newF16Harness(t, f16HarnessOpts{})
		created := hWrite.createTarget(t, "r112", "create-r112")
		h := newF16Harness(t, f16HarnessOpts{grants: f16ExecutionTargetGrants(f16ProviderUID, false, true, true)})
		h.life, h.cloud, h.audit, h.fixtures = hWrite.life, hWrite.cloud, hWrite.audit, hWrite.fixtures
		h.rebindGrants(f16ExecutionTargetGrants(f16ProviderUID, false, true, true))
		h.setAuditFail(errors.New("append boom"))
		rec := doF16(t, h, http.MethodPost, f16ActionPath(created.Metadata.UID, "retire"), nil,
			f16ActionHeaders("r112", created.Metadata.ResourceVersion))
		f16AssertInternalNoDisclosure(t, rec, http.StatusForbidden)
	})
	t.Run("inaccessible direct qualify", func(t *testing.T) {
		h := newF16Harness(t, f16HarnessOpts{})
		h.seedActiveParticipationAndStack(t)
		h.setAuditFail(errors.New("append boom"))
		rec := doF16(t, h, http.MethodPost, f16ActionPath("bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", "qualify"), nil,
			f16ActionHeaders("fail-q", "1"))
		f16AssertInternalNoDisclosure(t, rec, http.StatusNotFound)
	})
	t.Run("inaccessible direct retire", func(t *testing.T) {
		h := newF16Harness(t, f16HarnessOpts{})
		h.seedActiveParticipationAndStack(t)
		h.setAuditFail(errors.New("append boom"))
		rec := doF16(t, h, http.MethodPost, f16ActionPath("bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", "retire"), nil,
			f16ActionHeaders("fail-r", "1"))
		f16AssertInternalNoDisclosure(t, rec, http.StatusNotFound)
	})
}

func TestVS0_CF_F16_113(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	before := h.auditLen()
	rec := doF16(t, h, http.MethodHead, f16Collection, nil, nil)
	assertF16TransportOnly(t, rec, http.StatusMethodNotAllowed, "GET, POST")
	assertF16NoAuditDelta(t, h, before)
}

func TestVS0_CF_F16_114(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	created := h.createTarget(t, "head-act", "create-ha")
	before := h.auditLen()
	for _, action := range []string{"qualify", "retire"} {
		rec := doF16(t, h, http.MethodHead, f16ActionPath(created.Metadata.UID, action), nil, nil)
		assertF16TransportOnly(t, rec, http.StatusMethodNotAllowed, "POST")
	}
	assertF16NoAuditDelta(t, h, before)
}

func TestVS0_CF_F16_115(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	created := h.createTarget(t, "slash-fam", "create-sf")
	before := h.auditLen()
	members := []struct {
		method, path string
	}{
		{http.MethodPost, f16Collection + "/"},
		{http.MethodGet, f16Collection + "/"},
		{http.MethodPost, f16ActionPath(created.Metadata.UID, "qualify") + "/"},
		{http.MethodPost, f16ActionPath(created.Metadata.UID, "retire") + "/"},
	}
	for _, m := range members {
		rec := doF16(t, h, m.method, m.path, nil, nil)
		assertF16TransportOnly(t, rec, http.StatusNotFound, "")
	}
	assertF16NoAuditDelta(t, h, before)
}

func TestVS0_CF_F16_116(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	before := h.auditLen()
	rec := doF16(t, h, http.MethodPut, f16Collection, []byte(`{}`), map[string]string{"Content-Type": f16MediaJSON})
	assertF16TransportOnly(t, rec, http.StatusMethodNotAllowed, "GET, POST")
	assertF16NoAuditDelta(t, h, before)
}

func TestVS0_CF_F16_117(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	created := h.createTarget(t, "put-item", "create-pi")
	before := h.auditLen()
	rec := doF16(t, h, http.MethodPut, f16ItemPath(created.Metadata.UID), []byte(`{}`), map[string]string{"Content-Type": f16MediaJSON})
	assertF16TransportOnly(t, rec, http.StatusMethodNotAllowed, "GET")
	assertF16NoAuditDelta(t, h, before)
}

func TestVS0_CF_F16_118(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	created := h.createTarget(t, "put-act", "create-pa")
	before := h.auditLen()
	for _, action := range []string{"qualify", "retire"} {
		rec := doF16(t, h, http.MethodPut, f16ActionPath(created.Metadata.UID, action), []byte(`{}`), map[string]string{"Content-Type": f16MediaJSON})
		assertF16TransportOnly(t, rec, http.StatusMethodNotAllowed, "POST")
	}
	assertF16NoAuditDelta(t, h, before)
}

func TestVS0_CF_F16_119(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	created := h.createTarget(t, "via-down", "create-vd")
	h.setPresentFixture(created.Metadata.UID, f16AllSupportedTruths())
	rec := doF16(t, h, http.MethodPost, f16ActionPath(created.Metadata.UID, "qualify"), nil,
		f16ActionHeaders("q-vd", created.Metadata.ResourceVersion))
	qualified := f16DecodeSafe(t, rec)
	etag := qualified.Metadata.ResourceVersion
	before := h.auditLen()
	h.setParticipationEffective(t, f16PartUID, false)
	get := doF16(t, h, http.MethodGet, f16ItemPath(created.Metadata.UID), nil, nil)
	proj := f16DecodeSafe(t, get)
	if proj.EffectiveAvailability != etmodel.EffectiveUnavailable {
		t.Fatalf("availability=%s", proj.EffectiveAvailability)
	}
	if get.Header().Get("ETag") != etag {
		t.Fatal("ETag must be unchanged")
	}
	assertF16NoAuditDelta(t, h, before)
}

func TestVS0_CF_F16_120(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	created := h.createTarget(t, "via-up", "create-vu")
	h.setPresentFixture(created.Metadata.UID, f16AllSupportedTruths())
	rec := doF16(t, h, http.MethodPost, f16ActionPath(created.Metadata.UID, "qualify"), nil,
		f16ActionHeaders("q-vu", created.Metadata.ResourceVersion))
	qualified := f16DecodeSafe(t, rec)
	etag := qualified.Metadata.ResourceVersion
	before := h.auditLen()
	get := doF16(t, h, http.MethodGet, f16ItemPath(created.Metadata.UID), nil, nil)
	proj := f16DecodeSafe(t, get)
	if proj.EffectiveAvailability != etmodel.EffectiveAvailable {
		t.Fatalf("availability=%s", proj.EffectiveAvailability)
	}
	if get.Header().Get("ETag") != etag {
		t.Fatal("ETag must be unchanged")
	}
	assertF16NoAuditDelta(t, h, before)
}

func TestVS0_CF_F16_121(t *testing.T) {
	fixtures := executiontarget.NewMapFixtureSource()
	started := make(chan struct{})
	release := make(chan struct{})
	armed := &f16ArmedBlockingFixture{inner: fixtures, started: started, release: release}
	h := newF16Harness(t, f16HarnessOpts{
		fixtures: fixtures,
		observer: executiontarget.NewSyntheticObserver(armed, nil),
	})
	created := h.createTarget(t, "qualifying", "create-qualifying")
	h.setPresentFixture(created.Metadata.UID, f16AllSupportedTruths())
	// First qualify (unarmed) reaches Active/Qualified.
	rec := doF16(t, h, http.MethodPost, f16ActionPath(created.Metadata.UID, "qualify"), nil,
		f16ActionHeaders("q-first", created.Metadata.ResourceVersion))
	if rec.Code != http.StatusOK {
		t.Fatalf("first qualify=%d", rec.Code)
	}
	qualified := f16DecodeSafe(t, rec)
	beforeAudit := h.auditLen()
	beforeRV := qualified.Metadata.ResourceVersion
	storedBefore, _ := h.life.Store().GetExecutionTarget(created.Metadata.UID)
	fsBefore := ""
	if storedBefore.Status.FactSetRef != nil {
		fsBefore = storedBefore.Status.FactSetRef.UID
	}
	armed.Arm()
	qualDone := make(chan *httptest.ResponseRecorder, 1)
	go func() {
		qualDone <- doF16(t, h, http.MethodPost, f16ActionPath(created.Metadata.UID, "qualify"), nil,
			f16ActionHeaders("q-second", beforeRV))
	}()
	f16WaitStarted(t, started)
	// While Qualifying: public projection unchanged, no audit/completion yet.
	mid, _ := h.life.Store().GetExecutionTarget(created.Metadata.UID)
	if mid.Metadata.ResourceVersion != beforeRV {
		t.Fatal("Qualifying must retain public ETag")
	}
	if mid.Status.FactSetRef == nil || mid.Status.FactSetRef.UID != fsBefore {
		t.Fatal("must retain FactSet/Result links")
	}
	assertF16NoAuditDelta(t, h, beforeAudit)
	close(release)
	_ = f16WaitRec(t, qualDone)
}

func TestVS0_CF_F16_122(t *testing.T) {
	h := newF16Harness(t, f16HarnessOpts{})
	created := h.createTarget(t, "fault-from-q", "create-ffq")
	h.setPresentFixture(created.Metadata.UID, f16AllSupportedTruths())
	rec := doF16(t, h, http.MethodPost, f16ActionPath(created.Metadata.UID, "qualify"), nil,
		f16ActionHeaders("q-first", created.Metadata.ResourceVersion))
	qualified := f16DecodeSafe(t, rec)
	beforeRV := qualified.Metadata.ResourceVersion
	beforeAudit := h.auditLen()
	storedBefore, _ := h.life.Store().GetExecutionTarget(created.Metadata.UID)
	fsBefore := storedBefore.Status.FactSetRef.UID
	h.fixtures.Set(created.Metadata.UID, executiontarget.TargetBoundFixture{
		Revision: "fault", Mode: executiontarget.FixtureFaultDuplicate, Truths: f16AllSupportedTruths(),
	})
	rec2 := doF16(t, h, http.MethodPost, f16ActionPath(created.Metadata.UID, "qualify"), nil,
		f16ActionHeaders("q-fault", beforeRV))
	assertF16Problem(t, rec2, apiproblem.CodeInternalError, http.StatusInternalServerError, "")
	assertF16NoAuditDelta(t, h, beforeAudit)
	after, _ := h.life.Store().GetExecutionTarget(created.Metadata.UID)
	if after.Metadata.ResourceVersion != beforeRV {
		t.Fatal("must leave ETag unchanged")
	}
	if after.Status.FactSetRef == nil || after.Status.FactSetRef.UID != fsBefore {
		t.Fatal("must retain prior FactSet/Result links")
	}
}
