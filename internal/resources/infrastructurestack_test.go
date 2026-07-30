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

func TestInfrastructureStackKindConstant(t *testing.T) {
	t.Parallel()

	if KindInfrastructureStack != "InfrastructureStack" {
		t.Fatalf("KindInfrastructureStack = %q, want %q", KindInfrastructureStack, "InfrastructureStack")
	}
}

func TestInfrastructureStackJSONRoundTripWithTechnology(t *testing.T) {
	t.Parallel()

	in := InfrastructureStack{
		TypeMeta: apimeta.TypeMeta{
			APIVersion: FabricAPIVersion,
			Kind:       KindInfrastructureStack,
		},
		Metadata: apimeta.ObjectMeta{
			Name:        "stack-openshift-a",
			UID:         "uid-stack-openshift-a",
			DisplayName: "OpenShift Stack A",
			ScopeRef: &apimeta.ScopeRef{
				TypedRef: apimeta.TypedRef{
					APIVersion: FabricAPIVersion,
					Kind:       string(apimeta.ScopeProvider),
					Name:       "sovereign-provider-a",
					UID:        "uid-provider-a",
				},
			},
			Generation: 2,
		},
		Spec: InfrastructureStackSpec{
			DatacenterFailureDomainRef: DatacenterFailureDomainRef{
				TypedRef: apiref.TypedRef{
					APIVersion: FabricAPIVersion,
					Kind:       KindDatacenterFailureDomain,
					Name:       "fd-bangalore-az1",
					UID:        "uid-fd-bangalore-az1",
				},
			},
			Technology: "Red Hat OpenShift",
		},
		Status: InfrastructureStackStatus{
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
					Status:             apicond.ConditionTrue,
					Reason:             "LeafValid",
					Message:            "Leaf stack is valid.",
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

	var kind string
	if err := json.Unmarshal(payload["kind"], &kind); err != nil {
		t.Fatalf("unmarshal kind: %v", err)
	}
	if kind != KindInfrastructureStack {
		t.Fatalf("kind = %q, want %q", kind, KindInfrastructureStack)
	}

	var specObj map[string]any
	if err := json.Unmarshal(payload["spec"], &specObj); err != nil {
		t.Fatalf("unmarshal spec: %v", err)
	}
	if len(specObj) != 2 {
		t.Fatalf("spec must expose parent ref and technology, got %v", specObj)
	}
	refObj, ok := specObj["datacenterFailureDomainRef"].(map[string]any)
	if !ok {
		t.Fatalf("spec.datacenterFailureDomainRef missing or wrong type: %v", specObj["datacenterFailureDomainRef"])
	}
	if refObj["apiVersion"] != FabricAPIVersion {
		t.Fatalf("datacenterFailureDomainRef.apiVersion = %v, want %s", refObj["apiVersion"], FabricAPIVersion)
	}
	if refObj["kind"] != KindDatacenterFailureDomain {
		t.Fatalf("datacenterFailureDomainRef.kind = %v, want %s", refObj["kind"], KindDatacenterFailureDomain)
	}
	if refObj["name"] != "fd-bangalore-az1" {
		t.Fatalf("datacenterFailureDomainRef.name = %v, want fd-bangalore-az1", refObj["name"])
	}
	if refObj["uid"] != "uid-fd-bangalore-az1" {
		t.Fatalf("datacenterFailureDomainRef.uid = %v, want uid-fd-bangalore-az1", refObj["uid"])
	}
	if specObj["technology"] != "Red Hat OpenShift" {
		t.Fatalf("spec.technology = %v, want Red Hat OpenShift", specObj["technology"])
	}

	var out InfrastructureStack
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatalf("unmarshal InfrastructureStack: %v", err)
	}

	if out.APIVersion != FabricAPIVersion {
		t.Fatalf("apiVersion = %q, want %q", out.APIVersion, FabricAPIVersion)
	}
	if out.Kind != KindInfrastructureStack {
		t.Fatalf("kind = %q, want %q", out.Kind, KindInfrastructureStack)
	}
	if out.Metadata.Name != in.Metadata.Name {
		t.Fatalf("metadata.name = %q, want %q", out.Metadata.Name, in.Metadata.Name)
	}
	if out.Metadata.ScopeRef == nil {
		t.Fatal("metadata.scopeRef is nil")
	}
	if out.Metadata.ScopeRef.Kind != string(apimeta.ScopeProvider) {
		t.Fatalf("metadata.scopeRef.kind = %q, want %q", out.Metadata.ScopeRef.Kind, apimeta.ScopeProvider)
	}
	if out.Spec.DatacenterFailureDomainRef.Kind != KindDatacenterFailureDomain {
		t.Fatalf("spec.datacenterFailureDomainRef.kind = %q, want %q", out.Spec.DatacenterFailureDomainRef.Kind, KindDatacenterFailureDomain)
	}
	if out.Spec.DatacenterFailureDomainRef.Name != "fd-bangalore-az1" || out.Spec.DatacenterFailureDomainRef.UID != "uid-fd-bangalore-az1" {
		t.Fatalf("spec.datacenterFailureDomainRef = %+v", out.Spec.DatacenterFailureDomainRef)
	}
	if out.Spec.Technology != "Red Hat OpenShift" {
		t.Fatalf("spec.technology = %q, want %q", out.Spec.Technology, "Red Hat OpenShift")
	}
	if out.Status.ObservedGeneration != 2 {
		t.Fatalf("status.observedGeneration = %d, want 2", out.Status.ObservedGeneration)
	}
	if len(out.Status.Conditions) != 2 {
		t.Fatalf("status.conditions len = %d, want 2", len(out.Status.Conditions))
	}
}

func TestInfrastructureStackJSONRoundTripWithoutTechnology(t *testing.T) {
	t.Parallel()

	in := InfrastructureStack{
		TypeMeta: apimeta.TypeMeta{
			APIVersion: FabricAPIVersion,
			Kind:       KindInfrastructureStack,
		},
		Metadata: apimeta.ObjectMeta{
			Name: "stack-no-tech",
			UID:  "uid-stack-no-tech",
			ScopeRef: &apimeta.ScopeRef{
				TypedRef: apimeta.TypedRef{
					APIVersion: FabricAPIVersion,
					Kind:       string(apimeta.ScopeProvider),
					Name:       "sovereign-provider-a",
					UID:        "uid-provider-a",
				},
			},
		},
		Spec: InfrastructureStackSpec{
			DatacenterFailureDomainRef: DatacenterFailureDomainRef{
				TypedRef: apiref.TypedRef{
					APIVersion: FabricAPIVersion,
					Kind:       KindDatacenterFailureDomain,
					Name:       "fd-bangalore-az1",
					UID:        "uid-fd-bangalore-az1",
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
	if len(specObj) != 1 {
		t.Fatalf("spec without technology must expose only datacenterFailureDomainRef, got %v", specObj)
	}
	if _, present := specObj["technology"]; present {
		t.Fatalf("technology must be omitted when empty, got %v", specObj["technology"])
	}
	refObj, ok := specObj["datacenterFailureDomainRef"].(map[string]any)
	if !ok {
		t.Fatalf("spec.datacenterFailureDomainRef missing: %v", specObj)
	}
	if refObj["kind"] != KindDatacenterFailureDomain || refObj["name"] != "fd-bangalore-az1" {
		t.Fatalf("datacenterFailureDomainRef = %v", refObj)
	}

	var out InfrastructureStack
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if out.Kind != KindInfrastructureStack {
		t.Fatalf("kind = %q, want %q", out.Kind, KindInfrastructureStack)
	}
	if out.Spec.Technology != "" {
		t.Fatalf("technology = %q, want empty", out.Spec.Technology)
	}
	if out.Spec.DatacenterFailureDomainRef.Kind != KindDatacenterFailureDomain {
		t.Fatalf("datacenterFailureDomainRef.kind = %q", out.Spec.DatacenterFailureDomainRef.Kind)
	}
}

func TestInfrastructureStackDatacenterFailureDomainRefOptionalUID(t *testing.T) {
	t.Parallel()

	// Human-authored input MAY omit uid; the typed-ref fragment keeps it optional.
	in := InfrastructureStack{
		TypeMeta: apimeta.TypeMeta{
			APIVersion: FabricAPIVersion,
			Kind:       KindInfrastructureStack,
		},
		Metadata: apimeta.ObjectMeta{
			Name: "stack-no-uid",
			ScopeRef: &apimeta.ScopeRef{
				TypedRef: apimeta.TypedRef{
					APIVersion: FabricAPIVersion,
					Kind:       string(apimeta.ScopeProvider),
					Name:       "sovereign-provider-a",
					UID:        "uid-provider-a",
				},
			},
		},
		Spec: InfrastructureStackSpec{
			DatacenterFailureDomainRef: DatacenterFailureDomainRef{
				TypedRef: apiref.TypedRef{
					APIVersion: FabricAPIVersion,
					Kind:       KindDatacenterFailureDomain,
					Name:       "fd-bangalore-az1",
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
	refObj, ok := specObj["datacenterFailureDomainRef"].(map[string]any)
	if !ok {
		t.Fatalf("spec.datacenterFailureDomainRef missing: %v", specObj)
	}
	if _, present := refObj["uid"]; present {
		t.Fatalf("datacenterFailureDomainRef must omit uid when empty, got %v", refObj["uid"])
	}

	var out InfrastructureStack
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if out.Spec.DatacenterFailureDomainRef.UID != "" {
		t.Fatalf("uid = %q, want empty", out.Spec.DatacenterFailureDomainRef.UID)
	}
	if out.Spec.DatacenterFailureDomainRef.Kind != KindDatacenterFailureDomain || out.Spec.DatacenterFailureDomainRef.Name != "fd-bangalore-az1" {
		t.Fatalf("datacenterFailureDomainRef = %+v", out.Spec.DatacenterFailureDomainRef)
	}
}

func TestInfrastructureStackOmitsForbiddenFieldsAndSupersededKind(t *testing.T) {
	t.Parallel()

	// Superseded stack-kind name must not appear as an active type, constant,
	// tag, or alias (F14-REQ-03, F14-AD-021).
	superseded := "IaaSStack"
	if KindInfrastructureStack == superseded {
		t.Fatalf("KindInfrastructureStack must not equal superseded name %q", superseded)
	}

	forbidden := []string{
		// superseded stack-kind name (F14-REQ-03, F14-AD-021)
		superseded, "iaasStack", "IaasStack",
		// second / alternate / wrong-kind parent fields (DD-08, F14-REQ-09, F14-REQ-16)
		"Parent", "parent",
		"ParentRef", "parentRef",
		"ProviderRef", "providerRef",
		"ProviderParent", "providerParent",
		"ProviderLocationRef", "providerLocationRef",
		"ProviderDatacenterRef", "providerDatacenterRef",
		"DatacenterFailureDomainRefs", "datacenterFailureDomainRefs",
		"FailureDomainRef", "failureDomainRef",
		"FailureDomainRefs", "failureDomainRefs",
		"Owner", "owner",
		"OwnerRef", "ownerRef",
		// cross-kind topology parents must not appear as sibling fields
		"ProviderLocation", "providerLocation",
		"LocationRef", "locationRef",
		"InfrastructureStackRef", "infrastructureStackRef",
		// connectivity / capacity / capability / native
		"connectivity", "connected", "reachable", "adjacency",
		"capacity", "capability", "capabilities",
		"placement", "eligibility",
		"endpoint", "endpoints", "credential", "credentials",
		"providerID", "providerId", "nativeID", "nativeId",
		"availabilityZone", "availabilityZoneID", "nativeAvailabilityZone",
	}

	checkTypeAbsent := func(t *testing.T, typ reflect.Type, names []string) {
		t.Helper()
		for typ.Kind() == reflect.Pointer {
			typ = typ.Elem()
		}
		if strings.Contains(typ.Name(), superseded) {
			t.Fatalf("type name %q unexpectedly contains superseded stack-kind name", typ.Name())
		}
		for _, name := range names {
			if _, ok := typ.FieldByName(name); ok {
				t.Fatalf("%s unexpectedly declares field %q", typ.Name(), name)
			}
			for i := 0; i < typ.NumField(); i++ {
				f := typ.Field(i)
				if strings.Contains(f.Name, superseded) {
					t.Fatalf("%s unexpectedly declares field %q containing superseded name", typ.Name(), f.Name)
				}
				tag := f.Tag.Get("json")
				if tag == "" || tag == "-" {
					continue
				}
				jsonName, _, _ := strings.Cut(tag, ",")
				if jsonName == "" || jsonName == "-" {
					continue
				}
				if strings.Contains(jsonName, superseded) || strings.Contains(tag, superseded) {
					t.Fatalf("%s unexpectedly declares json tag %q containing superseded name", typ.Name(), tag)
				}
				if strings.EqualFold(jsonName, name) || jsonName == name {
					t.Fatalf("%s unexpectedly declares json tag %q on field %s", typ.Name(), tag, f.Name)
				}
			}
		}
	}

	checkTypeAbsent(t, reflect.TypeOf(InfrastructureStack{}), forbidden)
	checkTypeAbsent(t, reflect.TypeOf(InfrastructureStackSpec{}), forbidden)
	checkTypeAbsent(t, reflect.TypeOf(InfrastructureStackStatus{}), forbidden)
	checkTypeAbsent(t, reflect.TypeOf(DatacenterFailureDomainRef{}), forbidden)

	specType := reflect.TypeOf(InfrastructureStackSpec{})
	if specType.NumField() != 2 {
		t.Fatalf("InfrastructureStackSpec must expose parent ref and optional technology, got %d fields", specType.NumField())
	}
	refField, ok := specType.FieldByName("DatacenterFailureDomainRef")
	if !ok {
		t.Fatal("InfrastructureStackSpec.DatacenterFailureDomainRef missing")
	}
	if refField.Type != reflect.TypeOf(DatacenterFailureDomainRef{}) {
		t.Fatalf("DatacenterFailureDomainRef type = %v, want DatacenterFailureDomainRef", refField.Type)
	}
	if refField.Tag.Get("json") != "datacenterFailureDomainRef" {
		t.Fatalf("DatacenterFailureDomainRef json tag = %q, want %q", refField.Tag.Get("json"), "datacenterFailureDomainRef")
	}
	techField, ok := specType.FieldByName("Technology")
	if !ok {
		t.Fatal("InfrastructureStackSpec.Technology missing")
	}
	if techField.Type != reflect.TypeOf("") {
		t.Fatalf("Technology type = %v, want string", techField.Type)
	}
	if techField.Tag.Get("json") != "technology,omitempty" {
		t.Fatalf("Technology json tag = %q, want %q", techField.Tag.Get("json"), "technology,omitempty")
	}

	// Unknown second-parent / superseded-kind / connectivity / capacity keys
	// are absent from the type contract: decoding retains only the single
	// failure-domain parent and optional technology.
	raw := []byte(`{
		"apiVersion":"fabric.sovrunn.io/v1alpha1",
		"kind":"InfrastructureStack",
		"metadata":{"name":"stack1"},
		"spec":{
			"datacenterFailureDomainRef":{
				"apiVersion":"fabric.sovrunn.io/v1alpha1",
				"kind":"DatacenterFailureDomain",
				"name":"fd-bangalore-az1"
			},
			"technology":"Apache CloudStack",
			"datacenterFailureDomainRefs":[{"kind":"DatacenterFailureDomain","name":"fd-other"}],
			"failureDomainRef":{"kind":"DatacenterFailureDomain","name":"fd-alt"},
			"parentRef":{"kind":"Provider","name":"p1"},
			"providerLocationRef":{"kind":"ProviderLocation","name":"loc-in"},
			"providerDatacenterRef":{"kind":"ProviderDatacenter","name":"dc1"},
			"connectivity":true,
			"capacity":1,
			"placement":{},
			"nativeId":"az-1"
		},
		"kindAlias":"IaaSStack",
		"parent":"x",
		"connectivity":"up",
		"capacity":99
	}`)
	var got InfrastructureStack
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Kind != KindInfrastructureStack {
		t.Fatalf("kind = %q, want %q", got.Kind, KindInfrastructureStack)
	}
	if got.Spec.DatacenterFailureDomainRef.Kind != KindDatacenterFailureDomain || got.Spec.DatacenterFailureDomainRef.Name != "fd-bangalore-az1" {
		t.Fatalf("datacenterFailureDomainRef not retained: %+v", got.Spec.DatacenterFailureDomainRef)
	}
	if got.Spec.Technology != "Apache CloudStack" {
		t.Fatalf("technology not retained: %q", got.Spec.Technology)
	}
	checkTypeAbsent(t, reflect.TypeOf(got), forbidden)
	checkTypeAbsent(t, reflect.TypeOf(got.Spec), forbidden)
	if reflect.TypeOf(got.Spec).NumField() != 2 {
		t.Fatalf("decoded Spec must keep exactly one parent field plus technology, got %d", reflect.TypeOf(got.Spec).NumField())
	}
}

func TestInfrastructureStackIdenticalTechnologyRetainsDistinctIdentity(t *testing.T) {
	t.Parallel()

	// Boundary: shared technology never implies shared identity (F14-REQ-15,
	// F14-REQ-17, DD-05). Two stacks with identical technology retain distinct
	// uid and distinct parent references.
	tech := "AWS IaaS"
	a := InfrastructureStack{
		TypeMeta: apimeta.TypeMeta{
			APIVersion: FabricAPIVersion,
			Kind:       KindInfrastructureStack,
		},
		Metadata: apimeta.ObjectMeta{
			Name: "stack-a",
			UID:  "uid-stack-a",
			ScopeRef: &apimeta.ScopeRef{
				TypedRef: apimeta.TypedRef{
					APIVersion: FabricAPIVersion,
					Kind:       string(apimeta.ScopeProvider),
					Name:       "sovereign-provider-a",
					UID:        "uid-provider-a",
				},
			},
		},
		Spec: InfrastructureStackSpec{
			DatacenterFailureDomainRef: DatacenterFailureDomainRef{
				TypedRef: apiref.TypedRef{
					APIVersion: FabricAPIVersion,
					Kind:       KindDatacenterFailureDomain,
					Name:       "fd-a",
					UID:        "uid-fd-a",
				},
			},
			Technology: tech,
		},
	}
	b := InfrastructureStack{
		TypeMeta: apimeta.TypeMeta{
			APIVersion: FabricAPIVersion,
			Kind:       KindInfrastructureStack,
		},
		Metadata: apimeta.ObjectMeta{
			Name: "stack-b",
			UID:  "uid-stack-b",
			ScopeRef: &apimeta.ScopeRef{
				TypedRef: apimeta.TypedRef{
					APIVersion: FabricAPIVersion,
					Kind:       string(apimeta.ScopeProvider),
					Name:       "sovereign-provider-a",
					UID:        "uid-provider-a",
				},
			},
		},
		Spec: InfrastructureStackSpec{
			DatacenterFailureDomainRef: DatacenterFailureDomainRef{
				TypedRef: apiref.TypedRef{
					APIVersion: FabricAPIVersion,
					Kind:       KindDatacenterFailureDomain,
					Name:       "fd-b",
					UID:        "uid-fd-b",
				},
			},
			Technology: tech,
		},
	}

	if a.Spec.Technology != b.Spec.Technology {
		t.Fatalf("technology mismatch: %q vs %q", a.Spec.Technology, b.Spec.Technology)
	}
	if a.Metadata.UID == b.Metadata.UID {
		t.Fatalf("identical technology must not imply shared uid: both %q", a.Metadata.UID)
	}
	if a.Spec.DatacenterFailureDomainRef.Name == b.Spec.DatacenterFailureDomainRef.Name ||
		a.Spec.DatacenterFailureDomainRef.UID == b.Spec.DatacenterFailureDomainRef.UID {
		t.Fatalf("identical technology must not imply shared parent ref: a=%+v b=%+v",
			a.Spec.DatacenterFailureDomainRef, b.Spec.DatacenterFailureDomainRef)
	}

	rawA, err := json.Marshal(a)
	if err != nil {
		t.Fatalf("marshal a: %v", err)
	}
	rawB, err := json.Marshal(b)
	if err != nil {
		t.Fatalf("marshal b: %v", err)
	}
	var outA, outB InfrastructureStack
	if err := json.Unmarshal(rawA, &outA); err != nil {
		t.Fatalf("unmarshal a: %v", err)
	}
	if err := json.Unmarshal(rawB, &outB); err != nil {
		t.Fatalf("unmarshal b: %v", err)
	}
	if outA.Spec.Technology != tech || outB.Spec.Technology != tech {
		t.Fatalf("technology not preserved: a=%q b=%q", outA.Spec.Technology, outB.Spec.Technology)
	}
	if outA.Metadata.UID != "uid-stack-a" || outB.Metadata.UID != "uid-stack-b" {
		t.Fatalf("uids not preserved: a=%q b=%q", outA.Metadata.UID, outB.Metadata.UID)
	}
	if outA.Spec.DatacenterFailureDomainRef.UID != "uid-fd-a" || outB.Spec.DatacenterFailureDomainRef.UID != "uid-fd-b" {
		t.Fatalf("parent refs not preserved: a=%+v b=%+v",
			outA.Spec.DatacenterFailureDomainRef, outB.Spec.DatacenterFailureDomainRef)
	}
}

func TestInfrastructureStackDatacenterFailureDomainRefMatchesTypedRefFragment(t *testing.T) {
	t.Parallel()

	// Boundary: DatacenterFailureDomainRef shape matches
	// api/schemas/_common/typed-ref.json fields — apiVersion, kind, name,
	// and optional uid — via apiref.TypedRef.
	refType := reflect.TypeOf(DatacenterFailureDomainRef{})
	if refType.NumField() != 1 {
		t.Fatalf("DatacenterFailureDomainRef must embed exactly one TypedRef base, got %d fields", refType.NumField())
	}
	embedded := refType.Field(0)
	if !embedded.Anonymous {
		t.Fatal("DatacenterFailureDomainRef must anonymously embed the TypedRef base")
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

	ref := DatacenterFailureDomainRef{
		TypedRef: apiref.TypedRef{
			APIVersion: FabricAPIVersion,
			Kind:       KindDatacenterFailureDomain,
			Name:       "fd-boundary",
			UID:        "uid-fd-boundary",
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
			t.Fatalf("typed-ref fragment field %q missing from marshaled DatacenterFailureDomainRef: %v", key, obj)
		}
	}
	if len(obj) != 4 {
		t.Fatalf("DatacenterFailureDomainRef must serialize only typed-ref fields, got %v", obj)
	}
}
