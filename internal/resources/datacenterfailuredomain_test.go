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

func TestDatacenterFailureDomainKindConstant(t *testing.T) {
	t.Parallel()

	if KindDatacenterFailureDomain != "DatacenterFailureDomain" {
		t.Fatalf("KindDatacenterFailureDomain = %q, want %q", KindDatacenterFailureDomain, "DatacenterFailureDomain")
	}
}

func TestDatacenterFailureDomainJSONRoundTripWithProviderDatacenterRef(t *testing.T) {
	t.Parallel()

	in := DatacenterFailureDomain{
		TypeMeta: apimeta.TypeMeta{
			APIVersion: FabricAPIVersion,
			Kind:       KindDatacenterFailureDomain,
		},
		Metadata: apimeta.ObjectMeta{
			Name:        "fd-bangalore-az1",
			UID:         "uid-fd-bangalore-az1",
			DisplayName: "Bangalore AZ1 Failure Domain",
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
		Spec: DatacenterFailureDomainSpec{
			ProviderDatacenterRef: ProviderDatacenterRef{
				TypedRef: apiref.TypedRef{
					APIVersion: FabricAPIVersion,
					Kind:       KindProviderDatacenter,
					Name:       "dc-bangalore-1",
					UID:        "uid-dc-bangalore-1",
				},
			},
		},
		Status: DatacenterFailureDomainStatus{
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
	refObj, ok := specObj["providerDatacenterRef"].(map[string]any)
	if !ok {
		t.Fatalf("spec.providerDatacenterRef missing or wrong type: %v", specObj["providerDatacenterRef"])
	}
	if refObj["apiVersion"] != FabricAPIVersion {
		t.Fatalf("providerDatacenterRef.apiVersion = %v, want %s", refObj["apiVersion"], FabricAPIVersion)
	}
	if refObj["kind"] != KindProviderDatacenter {
		t.Fatalf("providerDatacenterRef.kind = %v, want %s", refObj["kind"], KindProviderDatacenter)
	}
	if refObj["name"] != "dc-bangalore-1" {
		t.Fatalf("providerDatacenterRef.name = %v, want dc-bangalore-1", refObj["name"])
	}
	if refObj["uid"] != "uid-dc-bangalore-1" {
		t.Fatalf("providerDatacenterRef.uid = %v, want uid-dc-bangalore-1", refObj["uid"])
	}

	var out DatacenterFailureDomain
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatalf("unmarshal DatacenterFailureDomain: %v", err)
	}

	if out.APIVersion != FabricAPIVersion {
		t.Fatalf("apiVersion = %q, want %q", out.APIVersion, FabricAPIVersion)
	}
	if out.Kind != KindDatacenterFailureDomain {
		t.Fatalf("kind = %q, want %q", out.Kind, KindDatacenterFailureDomain)
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
	if out.Spec.ProviderDatacenterRef.Kind != KindProviderDatacenter {
		t.Fatalf("spec.providerDatacenterRef.kind = %q, want %q", out.Spec.ProviderDatacenterRef.Kind, KindProviderDatacenter)
	}
	if out.Spec.ProviderDatacenterRef.Name != "dc-bangalore-1" || out.Spec.ProviderDatacenterRef.UID != "uid-dc-bangalore-1" {
		t.Fatalf("spec.providerDatacenterRef = %+v", out.Spec.ProviderDatacenterRef)
	}
	if out.Status.ObservedGeneration != 2 {
		t.Fatalf("status.observedGeneration = %d, want 2", out.Status.ObservedGeneration)
	}
	if len(out.Status.Conditions) != 2 {
		t.Fatalf("status.conditions len = %d, want 2", len(out.Status.Conditions))
	}
}

func TestDatacenterFailureDomainProviderDatacenterRefOptionalUID(t *testing.T) {
	t.Parallel()

	// Human-authored input MAY omit uid; the typed-ref fragment keeps it optional.
	in := DatacenterFailureDomain{
		TypeMeta: apimeta.TypeMeta{
			APIVersion: FabricAPIVersion,
			Kind:       KindDatacenterFailureDomain,
		},
		Metadata: apimeta.ObjectMeta{
			Name: "fd-no-uid",
			ScopeRef: &apimeta.ScopeRef{
				TypedRef: apimeta.TypedRef{
					APIVersion: FabricAPIVersion,
					Kind:       string(apimeta.ScopeCloudProvider),
					Name:       "sovereign-provider-a",
					UID:        "uid-provider-a",
				},
			},
		},
		Spec: DatacenterFailureDomainSpec{
			ProviderDatacenterRef: ProviderDatacenterRef{
				TypedRef: apiref.TypedRef{
					APIVersion: FabricAPIVersion,
					Kind:       KindProviderDatacenter,
					Name:       "dc-bangalore-1",
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
	refObj, ok := specObj["providerDatacenterRef"].(map[string]any)
	if !ok {
		t.Fatalf("spec.providerDatacenterRef missing: %v", specObj)
	}
	if _, present := refObj["uid"]; present {
		t.Fatalf("providerDatacenterRef must omit uid when empty, got %v", refObj["uid"])
	}

	var out DatacenterFailureDomain
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if out.Spec.ProviderDatacenterRef.UID != "" {
		t.Fatalf("uid = %q, want empty", out.Spec.ProviderDatacenterRef.UID)
	}
	if out.Spec.ProviderDatacenterRef.Kind != KindProviderDatacenter || out.Spec.ProviderDatacenterRef.Name != "dc-bangalore-1" {
		t.Fatalf("providerDatacenterRef = %+v", out.Spec.ProviderDatacenterRef)
	}
}

func TestDatacenterFailureDomainOmitsForbiddenFieldsFromTypeContract(t *testing.T) {
	t.Parallel()

	forbidden := []string{
		// second / alternate / wrong-kind parent fields (DD-08, F14-REQ-09)
		"Parent", "parent",
		"ParentRef", "parentRef",
		"ProviderRef", "providerRef",
		"ProviderParent", "providerParent",
		"ProviderLocationRef", "providerLocationRef",
		"DatacenterFailureDomainRef", "datacenterFailureDomainRef",
		"Owner", "owner",
		"OwnerRef", "ownerRef",
		// cross-kind topology parents must not appear as sibling fields
		"ProviderLocation", "providerLocation",
		"LocationRef", "locationRef",
		"InfrastructureStackRef", "infrastructureStackRef",
		// connectivity / capacity / capability / resilience / native
		"connectivity", "connected", "reachable", "adjacency",
		"capacity", "capability", "capabilities",
		"resilience", "availability", "quorum", "correlatedRisk", "correlated-risk",
		"placement", "eligibility",
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

	checkTypeAbsent(t, reflect.TypeOf(DatacenterFailureDomain{}), forbidden)
	checkTypeAbsent(t, reflect.TypeOf(DatacenterFailureDomainSpec{}), forbidden)
	checkTypeAbsent(t, reflect.TypeOf(DatacenterFailureDomainStatus{}), forbidden)
	checkTypeAbsent(t, reflect.TypeOf(ProviderDatacenterRef{}), forbidden)

	specType := reflect.TypeOf(DatacenterFailureDomainSpec{})
	if specType.NumField() != 1 {
		t.Fatalf("DatacenterFailureDomainSpec must expose only providerDatacenterRef, got %d fields", specType.NumField())
	}
	refField, ok := specType.FieldByName("ProviderDatacenterRef")
	if !ok {
		t.Fatal("DatacenterFailureDomainSpec.ProviderDatacenterRef missing")
	}
	if refField.Type != reflect.TypeOf(ProviderDatacenterRef{}) {
		t.Fatalf("ProviderDatacenterRef type = %v, want ProviderDatacenterRef", refField.Type)
	}
	if refField.Tag.Get("json") != "providerDatacenterRef" {
		t.Fatalf("ProviderDatacenterRef json tag = %q, want %q", refField.Tag.Get("json"), "providerDatacenterRef")
	}

	// Unknown second-parent / cross-kind / connectivity / resilience keys are
	// absent from the type contract: decoding retains only providerDatacenterRef.
	raw := []byte(`{
		"apiVersion":"fabric.sovrunn.io/v1alpha1",
		"kind":"DatacenterFailureDomain",
		"metadata":{"name":"fd1"},
		"spec":{
			"providerDatacenterRef":{
				"apiVersion":"fabric.sovrunn.io/v1alpha1",
				"kind":"ProviderDatacenter",
				"name":"dc-bangalore-1"
			},
			"parentRef":{"kind":"Provider","name":"p1"},
			"providerLocationRef":{"kind":"ProviderLocation","name":"loc-in"},
			"datacenterFailureDomainRef":{"kind":"DatacenterFailureDomain","name":"fd-other"},
			"locationRef":{"kind":"ProviderLocation","name":"loc-other"},
			"connectivity":true,
			"capacity":1,
			"resilience":"high",
			"availability":"zone",
			"quorum":3,
			"placement":{}
		},
		"parent":"x",
		"connectivity":"up",
		"capacity":99,
		"resilience":"high"
	}`)
	var got DatacenterFailureDomain
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Spec.ProviderDatacenterRef.Kind != KindProviderDatacenter || got.Spec.ProviderDatacenterRef.Name != "dc-bangalore-1" {
		t.Fatalf("providerDatacenterRef not retained: %+v", got.Spec.ProviderDatacenterRef)
	}
	checkTypeAbsent(t, reflect.TypeOf(got), forbidden)
	checkTypeAbsent(t, reflect.TypeOf(got.Spec), forbidden)
	if reflect.TypeOf(got.Spec).NumField() != 1 {
		t.Fatalf("decoded Spec must keep exactly one parent field, got %d", reflect.TypeOf(got.Spec).NumField())
	}
}

func TestDatacenterFailureDomainProviderDatacenterRefMatchesTypedRefFragment(t *testing.T) {
	t.Parallel()

	// Boundary: ProviderDatacenterRef shape matches api/schemas/_common/typed-ref.json
	// fields — apiVersion, kind, name, and optional uid — via apiref.TypedRef.
	refType := reflect.TypeOf(ProviderDatacenterRef{})
	if refType.NumField() != 1 {
		t.Fatalf("ProviderDatacenterRef must embed exactly one TypedRef base, got %d fields", refType.NumField())
	}
	embedded := refType.Field(0)
	if !embedded.Anonymous {
		t.Fatal("ProviderDatacenterRef must anonymously embed the TypedRef base")
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
	ref := ProviderDatacenterRef{
		TypedRef: apiref.TypedRef{
			APIVersion: FabricAPIVersion,
			Kind:       KindProviderDatacenter,
			Name:       "dc-boundary",
			UID:        "uid-dc-boundary",
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
			t.Fatalf("typed-ref fragment field %q missing from marshaled ProviderDatacenterRef: %v", key, obj)
		}
	}
	if len(obj) != 4 {
		t.Fatalf("ProviderDatacenterRef must serialize only typed-ref fields, got %v", obj)
	}
}
