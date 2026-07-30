package resources

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apicond"
	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
)

// designMaxStatusConditions is the §6.4 reviewed initial bound for
// status.conditions entries. The Go type uses a slice and does not imply
// a tighter hard limit than this contract bound.
const designMaxStatusConditions = 32

func TestProviderIdentityConstants(t *testing.T) {
	t.Parallel()

	if APIVersion != "fabric.sovrunn.io/v1alpha1" {
		t.Fatalf("APIVersion = %q, want %q", APIVersion, "fabric.sovrunn.io/v1alpha1")
	}
	cases := []struct {
		name string
		got  string
		want string
	}{
		{"KindProvider", KindProvider, "Provider"},
		{"KindProviderLocation", KindProviderLocation, "ProviderLocation"},
		{"KindProviderDatacenter", KindProviderDatacenter, "ProviderDatacenter"},
		{"KindDatacenterFailureDomain", KindDatacenterFailureDomain, "DatacenterFailureDomain"},
		{"KindInfrastructureStack", KindInfrastructureStack, "InfrastructureStack"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if tc.got != tc.want {
				t.Fatalf("%s = %q, want %q", tc.name, tc.got, tc.want)
			}
		})
	}
}

func TestProviderJSONRoundTrip(t *testing.T) {
	t.Parallel()

	in := Provider{
		TypeMeta: apimeta.TypeMeta{
			APIVersion: APIVersion,
			Kind:       KindProvider,
		},
		Metadata: apimeta.ObjectMeta{
			Name:        "sovereign-provider-a",
			UID:         "uid-provider-a",
			DisplayName: "Sovereign Provider A",
			ScopeRef: &apimeta.ScopeRef{
				TypedRef: apimeta.TypedRef{
					APIVersion: "platform.sovrunn.io/v1alpha1",
					Kind:       string(apimeta.ScopeOrganization),
					Name:       "owner-org-a",
					UID:        "uid-org-a",
				},
			},
			Labels:      map[string]string{"tier": "regulated"},
			Annotations: map[string]string{"note": "fixture"},
			Generation:  3,
		},
		Spec: ProviderSpec{},
		Status: ProviderStatus{
			ObservedGeneration: 3,
			Conditions: []apicond.Condition{
				{
					Type:               "Valid",
					Status:             apicond.ConditionTrue,
					Reason:             "ValidationSucceeded",
					Message:            "Resource is valid.",
					ObservedGeneration: 3,
					LastTransitionTime: "2026-07-30T00:00:00Z",
				},
				{
					Type:               "TopologyComplete",
					Status:             apicond.ConditionFalse,
					Reason:             "PathIncomplete",
					Message:            "Required child path is incomplete.",
					ObservedGeneration: 3,
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
	if len(specObj) != 0 {
		t.Fatalf("spec must be empty object, got %v", specObj)
	}

	var out Provider
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatalf("unmarshal Provider: %v", err)
	}

	if out.APIVersion != APIVersion {
		t.Fatalf("apiVersion = %q, want %q", out.APIVersion, APIVersion)
	}
	if out.Kind != KindProvider {
		t.Fatalf("kind = %q, want %q", out.Kind, KindProvider)
	}
	if out.Metadata.Name != in.Metadata.Name {
		t.Fatalf("metadata.name = %q, want %q", out.Metadata.Name, in.Metadata.Name)
	}
	if out.Metadata.UID != in.Metadata.UID {
		t.Fatalf("metadata.uid = %q, want %q", out.Metadata.UID, in.Metadata.UID)
	}
	if out.Metadata.DisplayName != in.Metadata.DisplayName {
		t.Fatalf("metadata.displayName = %q, want %q", out.Metadata.DisplayName, in.Metadata.DisplayName)
	}
	if out.Metadata.ScopeRef == nil {
		t.Fatal("metadata.scopeRef is nil")
	}
	if out.Metadata.ScopeRef.Kind != string(apimeta.ScopeOrganization) {
		t.Fatalf("metadata.scopeRef.kind = %q, want %q", out.Metadata.ScopeRef.Kind, apimeta.ScopeOrganization)
	}
	if out.Metadata.ScopeRef.Name != "owner-org-a" || out.Metadata.ScopeRef.UID != "uid-org-a" {
		t.Fatalf("metadata.scopeRef identity mismatch: %+v", out.Metadata.ScopeRef)
	}
	if out.Metadata.Labels["tier"] != "regulated" {
		t.Fatalf("metadata.labels = %v", out.Metadata.Labels)
	}
	if out.Metadata.Annotations["note"] != "fixture" {
		t.Fatalf("metadata.annotations = %v", out.Metadata.Annotations)
	}
	if out.Metadata.Generation != 3 {
		t.Fatalf("metadata.generation = %d, want 3", out.Metadata.Generation)
	}
	if out.Spec != (ProviderSpec{}) {
		t.Fatalf("spec = %#v, want empty ProviderSpec{}", out.Spec)
	}
	if out.Status.ObservedGeneration != 3 {
		t.Fatalf("status.observedGeneration = %d, want 3", out.Status.ObservedGeneration)
	}
	if len(out.Status.Conditions) != 2 {
		t.Fatalf("status.conditions len = %d, want 2", len(out.Status.Conditions))
	}
	if out.Status.Conditions[0].Type != "Valid" || out.Status.Conditions[0].Status != apicond.ConditionTrue {
		t.Fatalf("conditions[0] = %+v", out.Status.Conditions[0])
	}
	if out.Status.Conditions[1].Type != "TopologyComplete" || out.Status.Conditions[1].Status != apicond.ConditionFalse {
		t.Fatalf("conditions[1] = %+v", out.Status.Conditions[1])
	}
}

func TestProviderRejectsForbiddenFieldsByTypeContract(t *testing.T) {
	t.Parallel()

	forbidden := []string{
		// domain / owner fields
		"owner", "ownerRef", "Owner", "OwnerRef",
		"ProviderOwner", "SupplyOwner", "providerOwner", "supplyOwner",
		"providerID", "providerId", "nativeID", "nativeId",
		// connectivity / capacity / capability
		"connectivity", "connected", "reachable", "adjacency",
		"capacity", "capability", "capabilities",
		"endpoint", "endpoints", "credential", "credentials",
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

	checkTypeAbsent(t, reflect.TypeOf(Provider{}), forbidden)
	checkTypeAbsent(t, reflect.TypeOf(ProviderSpec{}), forbidden)
	checkTypeAbsent(t, reflect.TypeOf(ProviderStatus{}), forbidden)

	specType := reflect.TypeOf(ProviderSpec{})
	if specType.NumField() != 0 {
		t.Fatalf("ProviderSpec must have zero domain fields, got %d", specType.NumField())
	}

	// Unknown domain/owner/connectivity payload keys are absent from the
	// type contract: decoding into Provider leaves Spec empty and does not
	// materialize forbidden fields on the struct.
	raw := []byte(`{
		"apiVersion":"fabric.sovrunn.io/v1alpha1",
		"kind":"Provider",
		"metadata":{"name":"p1"},
		"spec":{
			"description":"not-allowed",
			"owner":"x",
			"ownerRef":{"kind":"Organization","name":"o1"},
			"ProviderOwner":"x",
			"connectivity":true,
			"capacity":1,
			"capability":"compute"
		},
		"owner":"x",
		"connectivity":"up",
		"capacity":99
	}`)
	var got Provider
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Spec != (ProviderSpec{}) {
		t.Fatalf("spec retained domain data: %#v", got.Spec)
	}
	checkTypeAbsent(t, reflect.TypeOf(got), forbidden)
	checkTypeAbsent(t, reflect.TypeOf(got.Spec), forbidden)
}

func TestProviderStatusConditionsTypeContract(t *testing.T) {
	t.Parallel()

	status := ProviderStatus{
		ObservedGeneration: 1,
		Conditions: []apicond.Condition{
			{Type: "Valid", Status: apicond.ConditionTrue, Reason: "ValidationSucceeded", ObservedGeneration: 1, LastTransitionTime: "2026-07-30T00:00:00Z"},
			{Type: "TopologyComplete", Status: apicond.ConditionFalse, Reason: "PathIncomplete", ObservedGeneration: 1, LastTransitionTime: "2026-07-30T00:00:00Z"},
		},
	}
	if len(status.Conditions) != 2 {
		t.Fatalf("conditions must accommodate Valid and TopologyComplete; len=%d", len(status.Conditions))
	}

	// The type contract is a slice: it accommodates the §6.4 bound of 32
	// and does not encode a tighter compile-time maximum.
	bound := make([]apicond.Condition, designMaxStatusConditions)
	for i := range bound {
		bound[i] = apicond.Condition{
			Type:               "Valid",
			Status:             apicond.ConditionUnknown,
			Reason:             "BoundCheck",
			ObservedGeneration: 1,
			LastTransitionTime: "2026-07-30T00:00:00Z",
		}
	}
	status.Conditions = bound
	if len(status.Conditions) != designMaxStatusConditions {
		t.Fatalf("conditions len = %d, want design bound %d", len(status.Conditions), designMaxStatusConditions)
	}

	field, ok := reflect.TypeOf(ProviderStatus{}).FieldByName("Conditions")
	if !ok {
		t.Fatal("ProviderStatus.Conditions missing")
	}
	if field.Type.Kind() != reflect.Slice {
		t.Fatalf("Conditions kind = %v, want slice (no tighter type-level max than §6.4)", field.Type.Kind())
	}
	if field.Type.Elem() != reflect.TypeOf(apicond.Condition{}) {
		t.Fatalf("Conditions element = %v, want apicond.Condition", field.Type.Elem())
	}
}
