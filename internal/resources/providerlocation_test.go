package resources

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apicond"
	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
)

func TestProviderLocationKindConstant(t *testing.T) {
	t.Parallel()

	if KindProviderLocation != "ProviderLocation" {
		t.Fatalf("KindProviderLocation = %q, want %q", KindProviderLocation, "ProviderLocation")
	}
}

func TestProviderLocationJSONRoundTripWithGeo(t *testing.T) {
	t.Parallel()

	in := ProviderLocation{
		TypeMeta: apimeta.TypeMeta{
			APIVersion: FabricAPIVersion,
			Kind:       KindProviderLocation,
		},
		Metadata: apimeta.ObjectMeta{
			Name:        "loc-in",
			UID:         "uid-loc-in",
			DisplayName: "India Location",
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
		Spec: ProviderLocationSpec{
			Geo: &ProviderLocationGeo{
				CountryCode:     "IN",
				SubdivisionCode: "IN-KA",
			},
		},
		Status: ProviderLocationStatus{
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
	geoObj, ok := specObj["geo"].(map[string]any)
	if !ok {
		t.Fatalf("spec.geo missing or wrong type: %v", specObj["geo"])
	}
	if geoObj["countryCode"] != "IN" {
		t.Fatalf("spec.geo.countryCode = %v, want IN", geoObj["countryCode"])
	}
	if geoObj["subdivisionCode"] != "IN-KA" {
		t.Fatalf("spec.geo.subdivisionCode = %v, want IN-KA", geoObj["subdivisionCode"])
	}

	var out ProviderLocation
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatalf("unmarshal ProviderLocation: %v", err)
	}

	if out.APIVersion != FabricAPIVersion {
		t.Fatalf("apiVersion = %q, want %q", out.APIVersion, FabricAPIVersion)
	}
	if out.Kind != KindProviderLocation {
		t.Fatalf("kind = %q, want %q", out.Kind, KindProviderLocation)
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
	if out.Spec.Geo == nil {
		t.Fatal("spec.geo is nil after round-trip with geo")
	}
	if out.Spec.Geo.CountryCode != "IN" || out.Spec.Geo.SubdivisionCode != "IN-KA" {
		t.Fatalf("spec.geo = %+v", out.Spec.Geo)
	}
	if out.Status.ObservedGeneration != 2 {
		t.Fatalf("status.observedGeneration = %d, want 2", out.Status.ObservedGeneration)
	}
	if len(out.Status.Conditions) != 2 {
		t.Fatalf("status.conditions len = %d, want 2", len(out.Status.Conditions))
	}
}

func TestProviderLocationJSONRoundTripWithoutGeo(t *testing.T) {
	t.Parallel()

	in := ProviderLocation{
		TypeMeta: apimeta.TypeMeta{
			APIVersion: FabricAPIVersion,
			Kind:       KindProviderLocation,
		},
		Metadata: apimeta.ObjectMeta{
			Name: "loc-no-geo",
			UID:  "uid-loc-no-geo",
			ScopeRef: &apimeta.ScopeRef{
				TypedRef: apimeta.TypedRef{
					APIVersion: FabricAPIVersion,
					Kind:       string(apimeta.ScopeProvider),
					Name:       "sovereign-provider-a",
					UID:        "uid-provider-a",
				},
			},
		},
		Spec: ProviderLocationSpec{},
		Status: ProviderLocationStatus{
			ObservedGeneration: 1,
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
	if _, ok := specObj["geo"]; ok {
		t.Fatalf("spec must omit geo when absent, got %v", specObj)
	}

	var out ProviderLocation
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatalf("unmarshal ProviderLocation: %v", err)
	}
	if out.Kind != KindProviderLocation {
		t.Fatalf("kind = %q, want %q", out.Kind, KindProviderLocation)
	}
	if out.Spec.Geo != nil {
		t.Fatalf("spec.geo = %+v, want nil", out.Spec.Geo)
	}
	if out.Spec != (ProviderLocationSpec{}) {
		t.Fatalf("spec = %#v, want empty ProviderLocationSpec{}", out.Spec)
	}
}

func TestProviderLocationGeoAbsenceHasNoDefault(t *testing.T) {
	t.Parallel()

	// Zero-value and explicit empty spec both leave Geo unset (nil);
	// the type contract carries no default country or subdivision.
	var zero ProviderLocation
	if zero.Spec.Geo != nil {
		t.Fatalf("zero-value Spec.Geo = %+v, want nil (no default)", zero.Spec.Geo)
	}

	empty := ProviderLocationSpec{}
	if empty.Geo != nil {
		t.Fatalf("empty Spec.Geo = %+v, want nil (no default)", empty.Geo)
	}

	raw := []byte(`{
		"apiVersion":"fabric.sovrunn.io/v1alpha1",
		"kind":"ProviderLocation",
		"metadata":{"name":"loc-absent-geo"},
		"spec":{}
	}`)
	var got ProviderLocation
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Spec.Geo != nil {
		t.Fatalf("decoded Spec.Geo = %+v, want nil (absence carries no default)", got.Spec.Geo)
	}
}

func TestProviderLocationOmitsForbiddenFieldsFromTypeContract(t *testing.T) {
	t.Parallel()

	forbidden := []string{
		// alternate location vocabulary (F14-REQ-02, F14-AD-006, F14-R02)
		"Region", "region",
		"Zone", "zone",
		"Site", "site",
		"Area", "area",
		"AvailabilityZone", "availabilityZone",
		"Location", "location",
		"ProviderRegion", "providerRegion",
		"GeoRegion", "geoRegion",
		"Locale", "locale",
		// parent / containment aliases (DD-08)
		"Parent", "parent",
		"ParentRef", "parentRef",
		"ProviderRef", "providerRef",
		"ProviderParent", "providerParent",
		"Owner", "owner",
		"OwnerRef", "ownerRef",
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

	checkTypeAbsent(t, reflect.TypeOf(ProviderLocation{}), forbidden)
	checkTypeAbsent(t, reflect.TypeOf(ProviderLocationSpec{}), forbidden)
	checkTypeAbsent(t, reflect.TypeOf(ProviderLocationGeo{}), forbidden)
	checkTypeAbsent(t, reflect.TypeOf(ProviderLocationStatus{}), forbidden)

	specType := reflect.TypeOf(ProviderLocationSpec{})
	if specType.NumField() != 1 {
		t.Fatalf("ProviderLocationSpec must expose only optional geo, got %d fields", specType.NumField())
	}
	geoField, ok := specType.FieldByName("Geo")
	if !ok {
		t.Fatal("ProviderLocationSpec.Geo missing")
	}
	if geoField.Type != reflect.TypeOf((*ProviderLocationGeo)(nil)) {
		t.Fatalf("Geo type = %v, want *ProviderLocationGeo", geoField.Type)
	}
	if geoField.Tag.Get("json") != "geo,omitempty" {
		t.Fatalf("Geo json tag = %q, want %q", geoField.Tag.Get("json"), "geo,omitempty")
	}

	// Unknown alternate-location / parent / connectivity payload keys are
	// absent from the type contract: decoding leaves only optional geo.
	raw := []byte(`{
		"apiVersion":"fabric.sovrunn.io/v1alpha1",
		"kind":"ProviderLocation",
		"metadata":{"name":"loc1"},
		"spec":{
			"region":"apac",
			"zone":"az1",
			"site":"dc-campus",
			"parent":"provider-a",
			"parentRef":{"kind":"Provider","name":"p1"},
			"providerRef":{"kind":"Provider","name":"p1"},
			"connectivity":true,
			"capacity":1,
			"geo":{"countryCode":"IN"}
		},
		"parent":"x",
		"connectivity":"up",
		"region":"apac"
	}`)
	var got ProviderLocation
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Spec.Geo == nil || got.Spec.Geo.CountryCode != "IN" {
		t.Fatalf("spec.geo not retained: %#v", got.Spec.Geo)
	}
	if got.Spec.Geo.SubdivisionCode != "" {
		t.Fatalf("subdivisionCode = %q, want empty", got.Spec.Geo.SubdivisionCode)
	}
	checkTypeAbsent(t, reflect.TypeOf(got), forbidden)
	checkTypeAbsent(t, reflect.TypeOf(got.Spec), forbidden)
}
