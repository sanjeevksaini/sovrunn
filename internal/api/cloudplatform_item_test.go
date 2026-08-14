package api

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/model"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/validate"
)

func seedCloudPlatform(t *testing.T, rt *f15Runtime, name string) model.CloudPlatform {
	t.Helper()
	mux := muxPlatform(rt)
	headers := map[string]string{"Content-Type": "application/json", "Idempotency-Key": "seed-" + name}
	rec := doMux(t, mux, http.MethodPost, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms", platformCreateBody(name), headers)
	if rec.Code != http.StatusCreated {
		t.Fatalf("seed create status=%d body=%s", rec.Code, rec.Body.String())
	}
	var cp model.CloudPlatform
	if err := json.Unmarshal(rec.Body.Bytes(), &cp); err != nil {
		t.Fatalf("seed decode: %v", err)
	}
	return cp
}

func TestCloudPlatformItem_GetPatchWriterMethods(t *testing.T) {
	rt := newF15Runtime(fullPlatformGrants(testPrincipal))
	mux := muxPlatform(rt)
	cp := seedCloudPlatform(t, rt, "plat-item")
	auditAfterCreate := rt.audit.Len()

	// Authorized GET.
	rec := doMux(t, mux, http.MethodGet, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms/"+cp.Metadata.UID, nil, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("get status=%d", rec.Code)
	}

	// Authorized PATCH description.
	patchHeaders := map[string]string{
		"Content-Type": mediaMergePatchJSON,
		"If-Match":     cp.Metadata.ResourceVersion,
	}
	rec = doMux(t, mux, http.MethodPatch, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms/"+cp.Metadata.UID,
		[]byte(`{"spec":{"description":"updated"}}`), patchHeaders)
	if rec.Code != http.StatusOK {
		t.Fatalf("patch status=%d body=%s", rec.Code, rec.Body.String())
	}
	var updated model.CloudPlatform
	_ = json.Unmarshal(rec.Body.Bytes(), &updated)
	if updated.Spec.Description != "updated" || updated.Metadata.ResourceVersion != "2" {
		t.Fatalf("patched resource=%#v", updated)
	}
	if rt.audit.Len() != auditAfterCreate+1 {
		t.Fatalf("patch audit count=%d want %d", rt.audit.Len(), auditAfterCreate+1)
	}

	// Stale If-Match → 412, no audit.
	before := rt.audit.Len()
	patchHeaders["If-Match"] = "1"
	rec = doMux(t, mux, http.MethodPatch, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms/"+cp.Metadata.UID,
		[]byte(`{"spec":{"description":"stale"}}`), patchHeaders)
	if rec.Code != http.StatusPreconditionFailed {
		t.Fatalf("stale status=%d", rec.Code)
	}
	if decodeProblem(t, rec).Code != apiproblem.CodeStaleResourceVersion {
		t.Fatalf("stale code")
	}
	if rt.audit.Len() != before {
		t.Fatal("stale must not audit")
	}

	// Missing If-Match → 412.
	delete(patchHeaders, "If-Match")
	rec = doMux(t, mux, http.MethodPatch, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms/"+cp.Metadata.UID,
		[]byte(`{"spec":{"description":"x"}}`), patchHeaders)
	if rec.Code != http.StatusPreconditionFailed {
		t.Fatalf("missing If-Match status=%d", rec.Code)
	}

	// Malformed If-Match.
	patchHeaders["If-Match"] = "bad\x00"
	rec = doMux(t, mux, http.MethodPatch, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms/"+cp.Metadata.UID,
		[]byte(`{"spec":{"description":"x"}}`), patchHeaders)
	if rec.Code != http.StatusPreconditionFailed {
		t.Fatalf("malformed If-Match status=%d", rec.Code)
	}

	// Immutable ownerRegistration → 422 VS0_PATCH_IMMUTABLE_FIELD, no audit.
	before = rt.audit.Len()
	patchHeaders["If-Match"] = updated.Metadata.ResourceVersion
	rec = doMux(t, mux, http.MethodPatch, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms/"+cp.Metadata.UID,
		[]byte(`{"spec":{"ownerRegistration":{"legalName":"Other","registrationIdentifier":"R1","jurisdictionCode":"US"}}}`), patchHeaders)
	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("immutable status=%d body=%s", rec.Code, rec.Body.String())
	}
	if !hasViolation(decodeProblem(t, rec), validate.ViolationPatchImmutableField) {
		t.Fatalf("immutable violation missing: %#v", decodeProblem(t, rec))
	}
	if rt.audit.Len() != before {
		t.Fatal("immutable must not audit")
	}

	// Immutable name.
	rec = doMux(t, mux, http.MethodPatch, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms/"+cp.Metadata.UID,
		[]byte(`{"metadata":{"name":"renamed"}}`), patchHeaders)
	if rec.Code != http.StatusUnprocessableEntity || !hasViolation(decodeProblem(t, rec), validate.ViolationPatchImmutableField) {
		t.Fatalf("name immutable: %d %#v", rec.Code, decodeProblem(t, rec))
	}

	// Client status write → audited 403.
	before = rt.audit.Len()
	rec = doMux(t, mux, http.MethodPatch, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms/"+cp.Metadata.UID,
		[]byte(`{"status":{"phase":"Active"}}`), patchHeaders)
	if rec.Code != http.StatusForbidden || !hasViolation(decodeProblem(t, rec), validate.ViolationStatusFieldWrite) {
		t.Fatalf("status patch: %d %#v", rec.Code, decodeProblem(t, rec))
	}
	if rt.audit.Len() != before+1 {
		t.Fatal("status patch must audit")
	}

	// Wrong-administrator → 403 before semantic, no audit.
	rtWrong := newF15Runtime(readOnlyPlatformGrants(testPrincipal))
	// Seed via store directly so wrong-admin runtime has a target.
	seeded, prob := rtWrong.store.CreateCloudPlatform(model.CloudPlatform{
		TypeMeta: apimetaTypeCloudPlatform(),
		Metadata: apimeta.ObjectMeta{Name: "plat-wrong", UID: "dddddddddddddddddddddddddddddddd", ResourceVersion: "1"},
		Spec: model.CloudPlatformSpec{OwnerRegistration: model.OwnerRegistration{
			LegalName: "Acme", RegistrationIdentifier: "R1", JurisdictionCode: "US",
		}},
		Status: model.CloudPlatformStatus{Phase: model.CloudPlatformPhaseActive},
	})
	if prob != nil {
		t.Fatalf("seed wrong-admin: %#v", prob)
	}
	muxWrong := muxPlatform(rtWrong)
	before = rtWrong.audit.Len()
	rec = doMux(t, muxWrong, http.MethodPatch, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms/"+seeded.Metadata.UID,
		[]byte(`{"spec":{"description":"nope"}}`), map[string]string{
			"Content-Type": mediaMergePatchJSON, "If-Match": "1",
		})
	if rec.Code != http.StatusForbidden {
		t.Fatalf("wrong-admin status=%d body=%s", rec.Code, rec.Body.String())
	}
	if rtWrong.audit.Len() != before {
		t.Fatal("wrong-admin must not audit")
	}

	// Forged bootstrap header before unsupported media (would be 415).
	rtForge := newF15Runtime(fullPlatformGrants(testPrincipal))
	cp2 := seedCloudPlatform(t, rtForge, "plat-forge")
	muxForge := muxPlatform(rtForge)
	before = rtForge.audit.Len()
	rec = doMux(t, muxForge, http.MethodPatch, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms/"+cp2.Metadata.UID,
		[]byte(`{"spec":{"description":"x"}}`), map[string]string{
			"Content-Type":                  "application/json",
			"If-Match":                      "1",
			cloudmodel.HeaderBootstrapGrant: "forged",
		})
	if rec.Code != http.StatusForbidden {
		t.Fatalf("forged header status=%d (must beat 415)", rec.Code)
	}
	if rtForge.audit.Len() != before+1 {
		t.Fatal("forged header must audit")
	}

	// PUT/DELETE → 405, no audit.
	before = rt.audit.Len()
	rec = doMux(t, mux, http.MethodPut, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms/"+cp.Metadata.UID, nil, nil)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("PUT status=%d", rec.Code)
	}
	rec = doMux(t, mux, http.MethodDelete, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms/"+cp.Metadata.UID, nil, nil)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("DELETE status=%d", rec.Code)
	}
	if rt.audit.Len() != before {
		t.Fatal("405 must not audit")
	}

	// PATCH never requires idempotency state.
	if rt.idemp.InFlightCountForTest() != 0 || rt.idemp.CompletedCountForTest() == 0 {
		// create completed records may exist; PATCH must not add InFlight.
	}
	if rt.idemp.InFlightCountForTest() != 0 {
		t.Fatal("PATCH must not leave InFlight idempotency state")
	}

	// Inaccessible GET → safe 404 + audit.
	rtNoRead := newF15Runtime(noReadPlatformGrants(testPrincipal))
	_, _ = rtNoRead.store.CreateCloudPlatform(model.CloudPlatform{
		TypeMeta: apimetaTypeCloudPlatform(),
		Metadata: apimeta.ObjectMeta{Name: "hidden", UID: "eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee", ResourceVersion: "1"},
		Spec: model.CloudPlatformSpec{OwnerRegistration: model.OwnerRegistration{
			LegalName: "Acme", RegistrationIdentifier: "R1", JurisdictionCode: "US",
		}},
		Status: model.CloudPlatformStatus{Phase: model.CloudPlatformPhaseActive},
	})
	muxNoRead := muxPlatform(rtNoRead)
	rec = doMux(t, muxNoRead, http.MethodGet, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms/eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee", nil, nil)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("inaccessible GET status=%d", rec.Code)
	}
	if !hasViolation(decodeProblem(t, rec), violationAuthorizationSafeDenial) {
		t.Fatalf("safe denial violation missing")
	}
	if rtNoRead.audit.Len() != 1 {
		t.Fatalf("inaccessible GET audit=%d", rtNoRead.audit.Len())
	}

	// StageSet plugged for PATCH.
	stages, ok := validate.StageSetFor(validate.KindCloudPlatform, validate.RequestClassPatch)
	if !ok || stages.Semantic == nil {
		t.Fatal("PATCH StageSet required")
	}
}

func apimetaTypeCloudPlatform() apimeta.TypeMeta {
	return apimeta.TypeMeta{APIVersion: model.APIVersionCloudPlatform, Kind: model.KindCloudPlatform}
}
