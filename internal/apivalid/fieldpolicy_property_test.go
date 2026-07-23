package apivalid

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
)

// Deterministic seed for Property 3 reproducibility
// (F12-VALIDATION-002, F12-META-002, F12-OWNER-002).
const property3Seed int64 = 20260723

const property3Iterations = 100

var (
	property3SystemKeys = []string{
		"uid",
		"generation",
		"resourceVersion",
		"createdAt",
		"updatedAt",
	}
	property3AllModes = []DecodeMode{
		ModeCreateRequest,
		ModeReplaceRequest,
		ModeStatusUpdate,
		ModeInternalObject,
		ModeReadRepresentation,
	}
	property3CustomerModes = []DecodeMode{
		ModeCreateRequest,
		ModeReplaceRequest,
	}
	property3AcceptingModes = []DecodeMode{
		ModeStatusUpdate,
		ModeInternalObject,
		ModeReadRepresentation,
	}
)

// property3Object describes a randomly generated object that always carries
// at least one status or system-owned field (the Property 3 precondition).
type property3Object struct {
	HasStatus bool
	System    []string // subset of systemOwnedMetadataKeys
	HasSpec   bool
	JSON      []byte
}

// Feature: api-resource-naming-status-and-validation-standard, Property 3: Operation-aware field ownership
//
// For any object carrying status or system-owned fields, customer mutation
// modes (ModeCreateRequest, ModeReplaceRequest) reject while ModeStatusUpdate,
// ModeInternalObject, and ModeReadRepresentation accept under Matrix C2
// ownership rules. Field acceptance is a deterministic function of DecodeMode.
//
// Validates: Requirements 4.3, 4.7, 4.9 (F12-VALIDATION-002, F12-META-002, F12-OWNER-002)
func TestProperty3_OperationAwareFieldOwnership(t *testing.T) {
	t.Parallel()

	rng := rand.New(rand.NewSource(property3Seed))
	for i := 0; i < property3Iterations; i++ {
		obj := generateProperty3Object(rng, i)
		if err := checkProperty3Object(obj, i); err != nil {
			t.Fatalf("property 3 failed at iteration %d (seed %d): %v", i, property3Seed, err)
		}
	}
}

func generateProperty3Object(rng *rand.Rand, iteration int) property3Object {
	// Always include at least one forbidden-under-customer field.
	hasStatus := rng.Intn(2) == 0
	nSys := rng.Intn(len(property3SystemKeys) + 1) // 0..len
	if !hasStatus && nSys == 0 {
		nSys = 1 + rng.Intn(len(property3SystemKeys))
	}
	sys := make([]string, 0, nSys)
	perm := rng.Perm(len(property3SystemKeys))
	for i := 0; i < nSys; i++ {
		sys = append(sys, property3SystemKeys[perm[i]])
	}

	// Status-update rejects spec; include spec only sometimes so Matrix C2
	// AllowSpecMutation is exercised for internal/read and denied for status-update.
	hasSpec := rng.Intn(2) == 0

	meta := map[string]any{
		"name": fmt.Sprintf("obj-%d-%d", iteration, rng.Intn(10000)),
	}
	for _, k := range sys {
		switch k {
		case "uid":
			meta[k] = fmt.Sprintf("%032x", rng.Uint64())
		case "generation":
			meta[k] = int64(1 + rng.Intn(100))
		case "resourceVersion":
			meta[k] = fmt.Sprintf("%d", 1+rng.Intn(1000))
		case "createdAt", "updatedAt":
			meta[k] = "2026-07-23T00:00:00Z"
		}
	}

	doc := map[string]any{
		"apiVersion": "v1",
		"kind":       "Project",
		"metadata":   meta,
	}
	if hasSpec {
		doc["spec"] = map[string]any{
			"displayName": fmt.Sprintf("Demo-%d", rng.Intn(1000)),
		}
	}
	if hasStatus {
		doc["status"] = map[string]any{
			"phase": "Ready",
		}
	}

	raw, err := json.Marshal(doc)
	if err != nil {
		// Generator must never fail; panic surfaces as a test crash with seed.
		panic(fmt.Sprintf("property3 marshal failed (seed %d iteration %d): %v", property3Seed, iteration, err))
	}

	return property3Object{
		HasStatus: hasStatus,
		System:    sys,
		HasSpec:   hasSpec,
		JSON:      raw,
	}
}

func checkProperty3Object(obj property3Object, iteration int) error {
	if !obj.HasStatus && len(obj.System) == 0 {
		return fmt.Errorf("iteration %d: generator violated precondition (no status/system fields)", iteration)
	}

	// PolicyFor must be a deterministic function of DecodeMode alone.
	for _, mode := range property3AllModes {
		a := PolicyFor(mode)
		b := PolicyFor(mode)
		if a != b {
			return fmt.Errorf("iteration %d: PolicyFor(%v) non-deterministic: %#v vs %#v", iteration, mode, a, b)
		}
		want := expectedPolicyFor(mode)
		if a != want {
			return fmt.Errorf("iteration %d: PolicyFor(%v) = %#v, want %#v", iteration, mode, a, want)
		}
	}

	// Customer mutation modes always reject objects with status/system fields.
	for _, mode := range property3CustomerModes {
		var dst decodeSample
		prob := DecodeJSON(obj.JSON, testLimits, PolicyFor(mode), &dst)
		if prob == nil {
			return fmt.Errorf("iteration %d: mode %v must reject status/system object, got nil (body=%s)",
				iteration, mode, obj.JSON)
		}
		if prob.Code != apiproblem.CodeValidationFailed {
			return fmt.Errorf("iteration %d: mode %v Code = %q, want %q",
				iteration, mode, prob.Code, apiproblem.CodeValidationFailed)
		}
		if len(prob.Violations) == 0 {
			return fmt.Errorf("iteration %d: mode %v expected violations, got none", iteration, mode)
		}
		field := prob.Violations[0].Field
		if !isProperty3OwnershipPointer(field) {
			return fmt.Errorf("iteration %d: mode %v unexpected field pointer %q (body=%s)",
				iteration, mode, field, obj.JSON)
		}
	}

	// Status-update / internal / read: accept under Matrix C2 ownership rules.
	// Status-update rejects when spec is present (AllowSpecMutation=false).
	for _, mode := range property3AcceptingModes {
		pol := PolicyFor(mode)
		var dst decodeSample
		prob := DecodeJSON(obj.JSON, testLimits, pol, &dst)

		expectReject := mode == ModeStatusUpdate && obj.HasSpec
		if expectReject {
			if prob == nil {
				return fmt.Errorf("iteration %d: ModeStatusUpdate must reject when spec present (body=%s)",
					iteration, obj.JSON)
			}
			if prob.Code != apiproblem.CodeValidationFailed {
				return fmt.Errorf("iteration %d: ModeStatusUpdate Code = %q, want %q",
					iteration, prob.Code, apiproblem.CodeValidationFailed)
			}
			if len(prob.Violations) == 0 || prob.Violations[0].Field != "/spec" {
				return fmt.Errorf("iteration %d: ModeStatusUpdate want /spec violation, got %#v",
					iteration, prob.Violations)
			}
			continue
		}

		if prob != nil {
			return fmt.Errorf("iteration %d: mode %v must accept under Matrix C2: %#v (body=%s)",
				iteration, mode, prob, obj.JSON)
		}
		if obj.HasStatus {
			if dst.Status["phase"] != "Ready" {
				return fmt.Errorf("iteration %d: mode %v status not decoded: %#v", iteration, mode, dst.Status)
			}
		}
		for _, k := range obj.System {
			if !metadataHasSystemField(dst.Metadata, k) {
				return fmt.Errorf("iteration %d: mode %v system field %q not decoded: %#v",
					iteration, mode, k, dst.Metadata)
			}
		}
	}

	return nil
}

func expectedPolicyFor(mode DecodeMode) FieldPolicy {
	switch mode {
	case ModeCreateRequest:
		return FieldPolicy{Mode: ModeCreateRequest, AllowSpecMutation: true}
	case ModeReplaceRequest:
		return FieldPolicy{Mode: ModeReplaceRequest, AllowSpecMutation: true}
	case ModeStatusUpdate:
		return FieldPolicy{Mode: ModeStatusUpdate, AllowStatus: true, AllowSystemOwned: true}
	case ModeInternalObject:
		return FieldPolicy{Mode: ModeInternalObject, AllowStatus: true, AllowSystemOwned: true, AllowSpecMutation: true}
	case ModeReadRepresentation:
		return FieldPolicy{Mode: ModeReadRepresentation, AllowStatus: true, AllowSystemOwned: true, AllowSpecMutation: true}
	default:
		return FieldPolicy{Mode: mode}
	}
}

func isProperty3OwnershipPointer(field string) bool {
	if field == "/status" {
		return true
	}
	switch field {
	case "/metadata/uid",
		"/metadata/generation",
		"/metadata/resourceVersion",
		"/metadata/createdAt",
		"/metadata/updatedAt":
		return true
	}
	return false
}

func metadataHasSystemField(m decodeMeta, key string) bool {
	switch key {
	case "uid":
		return m.UID != ""
	case "generation":
		return m.Generation != 0
	case "resourceVersion":
		return m.ResourceVersion != ""
	case "createdAt":
		return m.CreatedAt != ""
	case "updatedAt":
		return m.UpdatedAt != ""
	default:
		return false
	}
}
