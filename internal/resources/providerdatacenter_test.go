package resources

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apicond"
	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiref"
)

func TestProviderDatacenterKindConstant(t *testing.T) {
	t.Parallel()

	if KindProviderDatacenter != "ProviderDatacenter" {
		t.Fatalf("KindProviderDatacenter = %q, want %q", KindProviderDatacenter, "ProviderDatacenter")
	}
}

func TestProviderDatacenterJSONRoundTripWithProviderLocationRef(t *testing.T) {
	t.Parallel()

	in := ProviderDatacenter{
		TypeMeta: apimeta.TypeMeta{
			APIVersion: FabricAPIVersion,
			Kind:       KindProviderDatacenter,
		},
		Metadata: apimeta.ObjectMeta{
			Name:        "dc-bangalore-1",
			UID:         "uid-dc-bangalore-1",
			DisplayName: "Bangalore Datacenter 1",
			ScopeRef: &apimeta.ScopeRef{
				TypedRef: apimeta.TypedRef{
					APIVersion: FabricAPIVersion,
					Kind:       string(apimeta.ScopeCloudProvider),
					Name:       "sovereign-provider-a",
					UID:        "uid-provider-a",
				},
			},
			Generation: 2,
		},
		Spec: ProviderDatacenterSpec{
			ProviderLocationRef: ProviderLocationRef{
				TypedRef: apiref.TypedRef{
					APIVersion: FabricAPIVersion,
					Kind:       KindProviderLocation,
					Name:       "loc-in",
					UID:        "uid-loc-in",
				},
			},
		},
		Status: ProviderDatacenterStatus{
			ObservedGeneration: 2,
			Conditions: []apicond.Condition{
				{
					Type:               "Valid",
					Status:             apicond.ConditionTrue,
					Reason:             "ValidationSucceeded",
					Message:            "Resource is valid.",
					ObservedGeneration: 2,
					LastTransitionTime: "2026-07-30T00:00:00Z",
				},
				{
					Type:               "TopologyComplete",
					Status:             apicond.ConditionFalse,
					Reason:             "PathIncomplete",
					Message:            "Required child path is incomplete.",
					ObservedGeneration: 2,
					LastTransitionTime: "2026-07-30T00:00:00Z",
				},
			},
		},
	}

	raw, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var payload map[string]json.RawMessage
	if err := json.Unmarshal(raw, &payload); err != nil {
		t.Fatalf("unmarshal envelope: %v", err)
	}
	for _, key := range []string{"apiVersion", "kind", "metadata", "spec", "status"} {
		if _, ok := payload[key]; !ok {
			t.Fatalf("marshaled JSON missing %q", key)
		}
	}

	var specObj map[string]any
	if err := json.Unmarshal(payload["spec"], &specObj); err != nil {
		t.Fatalf("unmarshal spec: %v", err)
	}
	if len(specObj) != 1 {
		t.Fatalf("spec must expose exactly one domain field, got %v", specObj)
	}
	refObj, ok := specObj["providerLocationRef"].(map[string]any)
	if !ok {
		t.Fatalf("spec.providerLocationRef missing or wrong type: %v", specObj["providerLocationRef"])
	}
	if refObj["apiVersion"] != FabricAPIVersion {
		t.Fatalf("providerLocationRef.apiVersion = %v, want %s", refObj["apiVersion"], FabricAPIVersion)
	}
	if refObj["kind"] != KindProviderLocation {
		t.Fatalf("providerLocationRef.kind = %v, want %s", refObj["kind"], KindProviderLocation)
	}
	if refObj["name"] != "loc-in" {
		t.Fatalf("providerLocationRef.name = %v, want loc-in", refObj["name"])
	}
	if refObj["uid"] != "uid-loc-in" {
		t.Fatalf("providerLocationRef.uid = %v, want uid-loc-in", refObj["uid"])
	}

	var out ProviderDatacenter
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatalf("unmarshal ProviderDatacenter: %v", err)
	}

	if out.APIVersion != FabricAPIVersion {
		t.Fatalf("apiVersion = %q, want %q", out.APIVersion, FabricAPIVersion)
	}
	if out.Kind != KindProviderDatacenter {
		t.Fatalf("kind = %q, want %q", out.Kind, KindProviderDatacenter)
	}
	if out.Metadata.Name != in.Metadata.Name {
		t.Fatalf("metadata.name = %q, want %q", out.Metadata.Name, in.Metadata.Name)
	}
	if out.Metadata.ScopeRef == nil {
		t.Fatal("metadata.scopeRef is nil")
	}
	if out.Metadata.ScopeRef.Kind != string(apimeta.ScopeCloudProvider) {
		t.Fatalf("metadata.scopeRef.kind = %q, want %q", out.Metadata.ScopeRef.Kind, apimeta.ScopeCloudProvider)
	}
	if out.Spec.ProviderLocationRef.Kind != KindProviderLocation {
		t.Fatalf("spec.providerLocationRef.kind = %q, want %q", out.Spec.ProviderLocationRef.Kind, KindProviderLocation)
	}
	if out.Spec.ProviderLocationRef.Name != "loc-in" || out.Spec.ProviderLocationRef.UID != "uid-loc-in" {
		t.Fatalf("spec.providerLocationRef = %+v", out.Spec.ProviderLocationRef)
	}
	if out.Status.ObservedGeneration != 2 {
		t.Fatalf("status.observedGeneration = %d, want 2", out.Status.ObservedGeneration)
	}
	if len(out.Status.Conditions) != 2 {
		t.Fatalf("status.conditions len = %d, want 2", len(out.Status.Conditions))
	}
}

func TestProviderDatacenterProviderLocationRefOptionalUID(t *testing.T) {
	t.Parallel()

	// Human-authored input MAY omit uid; the typed-ref fragment keeps it optional.
	in := ProviderDatacenter{
		TypeMeta: apimeta.TypeMeta{
			APIVersion: FabricAPIVersion,
			Kind:       KindProviderDatacenter,
		},
		Metadata: apimeta.ObjectMeta{
			Name: "dc-no-uid",
			ScopeRef: &apimeta.ScopeRef{
				TypedRef: apimeta.TypedRef{
					APIVersion: FabricAPIVersion,
					Kind:       string(apimeta.ScopeCloudProvider),
					Name:       "sovereign-provider-a",
					UID:        "uid-provider-a",
				},
			},
		},
		Spec: ProviderDatacenterSpec{
			ProviderLocationRef: ProviderLocationRef{
				TypedRef: apiref.TypedRef{
					APIVersion: FabricAPIVersion,
					Kind:       KindProviderLocation,
					Name:       "loc-in",
				},
			},
		},
	}

	raw, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var payload map[string]json.RawMessage
	if err := json.Unmarshal(raw, &payload); err != nil {
		t.Fatalf("unmarshal envelope: %v", err)
	}
	var specObj map[string]any
	if err := json.Unmarshal(payload["spec"], &specObj); err != nil {
		t.Fatalf("unmarshal spec: %v", err)
	}
	refObj, ok := specObj["providerLocationRef"].(map[string]any)
	if !ok {
		t.Fatalf("spec.providerLocationRef missing: %v", specObj)
	}
	if _, present := refObj["uid"]; present {
		t.Fatalf("providerLocationRef must omit uid when empty, got %v", refObj["uid"])
	}

	var out ProviderDatacenter
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if out.Spec.ProviderLocationRef.UID != "" {
		t.Fatalf("uid = %q, want empty", out.Spec.ProviderLocationRef.UID)
	}
	if out.Spec.ProviderLocationRef.Kind != KindProviderLocation || out.Spec.ProviderLocationRef.Name != "loc-in" {
		t.Fatalf("providerLocationRef = %+v", out.Spec.ProviderLocationRef)
	}
}

func TestProviderDatacenterOmitsForbiddenFieldsFromTypeContract(t *testing.T) {
	t.Parallel()

	forbidden := []string{
		// second / alternate parent fields (DD-08, F14-REQ-09)
		"Parent", "parent",
		"ParentRef", "parentRef",
		"ProviderRef", "providerRef",
		"ProviderParent", "providerParent",
		"ProviderDatacenterRef", "providerDatacenterRef",
		"DatacenterFailureDomainRef", "datacenterFailureDomainRef",
		"Owner", "owner",
		"OwnerRef", "ownerRef",
		// cross-kind topology parents must not appear as sibling fields
		"ProviderLocation", "providerLocation",
		"LocationRef", "locationRef",
		"InfrastructureStackRef", "infrastructureStackRef",
		// connectivity / capacity / capability / native
		"connectivity", "connected", "reachable", "adjacency",
		"capacity", "capability", "capabilities",
		"endpoint", "endpoints", "credential", "credentials",
		"providerID", "providerId", "nativeID", "nativeId",
	}

	checkTypeAbsent := func(t *testing.T, typ reflect.Type, names []string) {
		t.Helper()
		for typ.Kind() == reflect.Pointer {
			typ = typ.Elem()
		}
		for _, name := range names {
			if _, ok := typ.FieldByName(name); ok {
				t.Fatalf("%s unexpectedly declares field %q", typ.Name(), name)
			}
			for i := 0; i < typ.NumField(); i++ {
				f := typ.Field(i)
				tag := f.Tag.Get("json")
				if tag == "" || tag == "-" {
					continue
				}
				jsonName, _, _ := strings.Cut(tag, ",")
				if jsonName == "" || jsonName == "-" {
					continue
				}
				if strings.EqualFold(jsonName, name) || jsonName == name {
					t.Fatalf("%s unexpectedly declares json tag %q on field %s", typ.Name(), tag, f.Name)
				}
			}
		}
	}

	checkTypeAbsent(t, reflect.TypeOf(ProviderDatacenter{}), forbidden)
	checkTypeAbsent(t, reflect.TypeOf(ProviderDatacenterSpec{}), forbidden)
	checkTypeAbsent(t, reflect.TypeOf(ProviderDatacenterStatus{}), forbidden)
	checkTypeAbsent(t, reflect.TypeOf(ProviderLocationRef{}), forbidden)

	specType := reflect.TypeOf(ProviderDatacenterSpec{})
	if specType.NumField() != 1 {
		t.Fatalf("ProviderDatacenterSpec must expose only providerLocationRef, got %d fields", specType.NumField())
	}
	refField, ok := specType.FieldByName("ProviderLocationRef")
	if !ok {
		t.Fatal("ProviderDatacenterSpec.ProviderLocationRef missing")
	}
	if refField.Type != reflect.TypeOf(ProviderLocationRef{}) {
		t.Fatalf("ProviderLocationRef type = %v, want ProviderLocationRef", refField.Type)
	}
	if refField.Tag.Get("json") != "providerLocationRef" {
		t.Fatalf("ProviderLocationRef json tag = %q, want %q", refField.Tag.Get("json"), "providerLocationRef")
	}

	// Unknown second-parent / cross-kind / connectivity keys are absent from
	// the type contract: decoding retains only providerLocationRef.
	raw := []byte(`{
		"apiVersion":"fabric.sovrunn.io/v1alpha1",
		"kind":"ProviderDatacenter",
		"metadata":{"name":"dc1"},
		"spec":{
			"providerLocationRef":{
				"apiVersion":"fabric.sovrunn.io/v1alpha1",
				"kind":"ProviderLocation",
				"name":"loc-in"
			},
			"parentRef":{"kind":"Provider","name":"p1"},
			"providerDatacenterRef":{"kind":"ProviderDatacenter","name":"dc-other"},
			"datacenterFailureDomainRef":{"kind":"DatacenterFailureDomain","name":"fd1"},
			"locationRef":{"kind":"ProviderLocation","name":"loc-other"},
			"connectivity":true,
			"capacity":1
		},
		"parent":"x",
		"connectivity":"up",
		"capacity":99
	}`)
	var got ProviderDatacenter
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Spec.ProviderLocationRef.Kind != KindProviderLocation || got.Spec.ProviderLocationRef.Name != "loc-in" {
		t.Fatalf("providerLocationRef not retained: %+v", got.Spec.ProviderLocationRef)
	}
	checkTypeAbsent(t, reflect.TypeOf(got), forbidden)
	checkTypeAbsent(t, reflect.TypeOf(got.Spec), forbidden)
	if reflect.TypeOf(got.Spec).NumField() != 1 {
		t.Fatalf("decoded Spec must keep exactly one parent field, got %d", reflect.TypeOf(got.Spec).NumField())
	}
}

func TestProviderDatacenterProviderLocationRefMatchesTypedRefFragment(t *testing.T) {
	t.Parallel()

	// Boundary: ProviderLocationRef shape matches api/schemas/_common/typed-ref.json
	// fields — apiVersion, kind, name, and optional uid — via apiref.TypedRef.
	refType := reflect.TypeOf(ProviderLocationRef{})
	if refType.NumField() != 1 {
		t.Fatalf("ProviderLocationRef must embed exactly one TypedRef base, got %d fields", refType.NumField())
	}
	embedded := refType.Field(0)
	if !embedded.Anonymous {
		t.Fatal("ProviderLocationRef must anonymously embed the TypedRef base")
	}
	if embedded.Type != reflect.TypeOf(apiref.TypedRef{}) {
		t.Fatalf("embedded type = %v, want apiref.TypedRef", embedded.Type)
	}

	base := reflect.TypeOf(apiref.TypedRef{})
	wantFields := map[string]string{
		"APIVersion": "apiVersion",
		"Kind":       "kind",
		"Name":       "name",
		"UID":        "uid,omitempty",
	}
	if base.NumField() != len(wantFields) {
		t.Fatalf("TypedRef field count = %d, want %d (_common/typed-ref)", base.NumField(), len(wantFields))
	}
	for name, tag := range wantFields {
		f, ok := base.FieldByName(name)
		if !ok {
			t.Fatalf("TypedRef missing field %q", name)
		}
		if f.Tag.Get("json") != tag {
			t.Fatalf("TypedRef.%s json tag = %q, want %q", name, f.Tag.Get("json"), tag)
		}
	}

	// Promoted accessors match the fragment field names after round-trip.
	ref := ProviderLocationRef{
		TypedRef: apiref.TypedRef{
			APIVersion: FabricAPIVersion,
			Kind:       KindProviderLocation,
			Name:       "loc-boundary",
			UID:        "uid-loc-boundary",
		},
	}
	raw, err := json.Marshal(ref)
	if err != nil {
		t.Fatalf("marshal ref: %v", err)
	}
	var obj map[string]any
	if err := json.Unmarshal(raw, &obj); err != nil {
		t.Fatalf("unmarshal ref map: %v", err)
	}
	for _, key := range []string{"apiVersion", "kind", "name", "uid"} {
		if _, ok := obj[key]; !ok {
			t.Fatalf("typed-ref fragment field %q missing from marshaled ProviderLocationRef: %v", key, obj)
		}
	}
	if len(obj) != 4 {
		t.Fatalf("ProviderLocationRef must serialize only typed-ref fields, got %v", obj)
	}
}
