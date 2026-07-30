package apiconform

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apiref"
	"github.com/sanjeevksaini/sovrunn/internal/apivalid"
	"github.com/sanjeevksaini/sovrunn/internal/resources"
	"github.com/sanjeevksaini/sovrunn/internal/validation"
)

// FEATURE-0014 Task 20 positive conformance fixtures (design §9.2; I-6/I-7).
// Scenario evaluation uses Task 17/18 helpers over explicitly supplied
// in-memory state only — no production store, registry, or live handlers.

const (
	fixtureProvider                = "provider.json"
	fixtureProviderLocation        = "provider-location.json"
	fixtureProviderDatacenter      = "provider-datacenter.json"
	fixtureDatacenterFailureDomain = "datacenter-failure-domain.json"
	fixtureInfrastructureStack     = "infrastructure-stack.json"

	fixtureProviderUID                = "a2000000000000000000000000000001"
	fixtureProviderLocationUID        = "a3000000000000000000000000000001"
	fixtureProviderDatacenterUID      = "a4000000000000000000000000000001"
	fixtureDatacenterFailureDomainUID = "a5000000000000000000000000000001"
	fixtureInfrastructureStackUID     = "a6000000000000000000000000000001"
	fixtureOwnerOrganizationUID       = "a1000000000000000000000000000001"
)

type feature0014FixtureCase struct {
	file     string
	schemaID string
	newDst   func() any
	validate func(ctx context.Context, createJSON []byte, structural apivalid.StructuralValidator) *apiproblem.Problem
}

func TestFeature0014PositiveFixturesValidateAgainstSchemaAndValidator(t *testing.T) {
	t.Parallel()

	root := moduleRoot(t)
	structural := mustFeature0014Structural(t, root)
	lim := apivalid.DefaultLimits()
	pol := apivalid.PolicyFor(apivalid.ModeReadRepresentation)
	ctx := context.Background()

	cases := []feature0014FixtureCase{
		{
			file:     fixtureProvider,
			schemaID: CanonicalSchemasDir + "/provider.json",
			newDst:   func() any { return &resources.Provider{} },
			validate: validation.ValidateProvider,
		},
		{
			file:     fixtureProviderLocation,
			schemaID: CanonicalSchemasDir + "/provider-location.json",
			newDst:   func() any { return &resources.ProviderLocation{} },
			validate: validation.ValidateProviderLocation,
		},
		{
			file:     fixtureProviderDatacenter,
			schemaID: CanonicalSchemasDir + "/provider-datacenter.json",
			newDst:   func() any { return &resources.ProviderDatacenter{} },
			validate: validation.ValidateProviderDatacenter,
		},
		{
			file:     fixtureDatacenterFailureDomain,
			schemaID: CanonicalSchemasDir + "/datacenter-failure-domain.json",
			newDst:   func() any { return &resources.DatacenterFailureDomain{} },
			validate: validation.ValidateDatacenterFailureDomain,
		},
		{
			file:     fixtureInfrastructureStack,
			schemaID: CanonicalSchemasDir + "/infrastructure-stack.json",
			newDst:   func() any { return &resources.InfrastructureStack{} },
			validate: validation.ValidateInfrastructureStack,
		},
	}

	if len(cases) != 5 {
		t.Fatalf("expected exactly five FEATURE-0014 fixtures, got %d", len(cases))
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.file, func(t *testing.T) {
			t.Parallel()

			raw := mustReadFeature0014Fixture(t, root, tc.file)
			dst := tc.newDst()
			if prob := apivalid.DecodeJSON(raw, lim, pol, dst); prob != nil {
				t.Fatalf("DecodeJSON ModeReadRepresentation: code=%s detail=%s violations=%v",
					prob.Code, prob.Detail, prob.Violations)
			}
			violations, err := structural.Validate(dst, tc.schemaID)
			if err != nil {
				t.Fatalf("Validate(%s): %v", tc.schemaID, err)
			}
			if len(violations) != 0 {
				t.Fatalf("unexpected structural violations for %s: %#v", tc.file, violations)
			}

			createJSON := mustStripSystemOwnedForCreate(t, raw)
			if prob := tc.validate(ctx, createJSON, structural); prob != nil {
				t.Fatalf("offline validator rejected create-shaped %s: code=%s detail=%s violations=%v",
					tc.file, prob.Code, prob.Detail, prob.Violations)
			}
		})
	}
}

func TestFeature0014OwnerOperatorRoleSpecificIdentities(t *testing.T) {
	t.Parallel()

	// Distinct-owner/operator: Organization OwnerOrganization-A → two Providers
	// ProviderOperator-A and ProviderOperator-B with distinct UIDs (F14-REQ-11).
	ownerAUID := fixtureOwnerOrganizationUID
	providerA := resources.Provider{
		TypeMeta: apimeta.TypeMeta{APIVersion: resources.FabricAPIVersion, Kind: resources.KindProvider},
		Metadata: apimeta.ObjectMeta{
			Name: "provider-operator-a",
			UID:  fixtureProviderUID,
			ScopeRef: &apimeta.ScopeRef{TypedRef: apiref.TypedRef{
				APIVersion: "core.sovrunn.io/v1alpha1",
				Kind:       string(apimeta.ScopeOrganization),
				Name:       "owner-organization-a",
				UID:        ownerAUID,
			}},
		},
	}
	providerB := resources.Provider{
		TypeMeta: apimeta.TypeMeta{APIVersion: resources.FabricAPIVersion, Kind: resources.KindProvider},
		Metadata: apimeta.ObjectMeta{
			Name: "provider-operator-b",
			UID:  "b2000000000000000000000000000002",
			ScopeRef: &apimeta.ScopeRef{TypedRef: apiref.TypedRef{
				APIVersion: "core.sovrunn.io/v1alpha1",
				Kind:       string(apimeta.ScopeOrganization),
				Name:       "owner-organization-a",
				UID:        ownerAUID,
			}},
		},
	}

	tvA := validation.TopologyValueFromProvider(providerA)
	tvB := validation.TopologyValueFromProvider(providerB)
	if tvA.UID == "" || tvB.UID == "" || tvA.UID == tvB.UID {
		t.Fatalf("distinct operators under one owner must have distinct Provider UIDs: %q vs %q", tvA.UID, tvB.UID)
	}
	if tvA.ProviderScopeUID != fixtureProviderUID || tvB.ProviderScopeUID != "b2000000000000000000000000000002" {
		t.Fatalf("Provider scope UIDs must be the Provider resource UIDs: a=%q b=%q", tvA.ProviderScopeUID, tvB.ProviderScopeUID)
	}
	if providerA.Metadata.ScopeRef.UID != ownerAUID || providerB.Metadata.ScopeRef.UID != ownerAUID {
		t.Fatal("both Providers must remain scoped to OwnerOrganization-A")
	}

	// Same-party: Organization ProviderOperator-A → Provider ProviderOperator-A
	// are two role-specific resources with distinct UIDs (F14-REQ-11).
	samePartyOrgUID := "c1000000000000000000000000000001"
	samePartyProviderUID := "c2000000000000000000000000000001"
	samePartyProvider := resources.Provider{
		TypeMeta: apimeta.TypeMeta{APIVersion: resources.FabricAPIVersion, Kind: resources.KindProvider},
		Metadata: apimeta.ObjectMeta{
			Name: "provider-operator-a",
			UID:  samePartyProviderUID,
			ScopeRef: &apimeta.ScopeRef{TypedRef: apiref.TypedRef{
				APIVersion: "core.sovrunn.io/v1alpha1",
				Kind:       string(apimeta.ScopeOrganization),
				Name:       "provider-operator-a",
				UID:        samePartyOrgUID,
			}},
		},
	}
	sameTV := validation.TopologyValueFromProvider(samePartyProvider)
	if sameTV.UID != samePartyProviderUID {
		t.Fatalf("same-party Provider uid=%q want %q", sameTV.UID, samePartyProviderUID)
	}
	if samePartyProvider.Metadata.ScopeRef.UID != samePartyOrgUID {
		t.Fatalf("same-party Organization uid=%q want %q", samePartyProvider.Metadata.ScopeRef.UID, samePartyOrgUID)
	}
	if samePartyProvider.Metadata.ScopeRef.UID == samePartyProvider.Metadata.UID {
		t.Fatal("Organization and Provider must remain distinct role-specific resources with distinct UIDs")
	}
	if samePartyProvider.Metadata.Name != samePartyProvider.Metadata.ScopeRef.Name {
		t.Fatal("same-party case uses matching names across roles without collapsing identity")
	}
}

func TestFeature0014FullPathCompleteAndEmptyParentsIncomplete(t *testing.T) {
	t.Parallel()

	root := moduleRoot(t)
	set := mustLoadFeature0014CompletenessSet(t, root)
	provider := findFeature0014ByKind(t, set, resources.KindProvider)

	got := validation.EvaluateTopologyCompleteness(provider, set)
	if !got.Complete || got.Reason != validation.ReasonPathComplete {
		t.Fatalf("full five-level path: complete=%v reason=%q", got.Complete, got.Reason)
	}

	for _, kind := range []string{
		resources.KindProvider,
		resources.KindProviderLocation,
		resources.KindProviderDatacenter,
		resources.KindDatacenterFailureDomain,
	} {
		v := findFeature0014ByKind(t, set, kind)
		ev := validation.EvaluateTopologyCompleteness(v, set)
		if !ev.Complete || ev.Reason != validation.ReasonPathComplete {
			t.Fatalf("%s: complete=%v reason=%q want PathComplete", kind, ev.Complete, ev.Reason)
		}
	}

	stack := findFeature0014ByKind(t, set, resources.KindInfrastructureStack)
	leaf := validation.EvaluateTopologyCompleteness(stack, set)
	if !leaf.Complete || leaf.Reason != validation.ReasonLeafValid {
		t.Fatalf("leaf stack: complete=%v reason=%q want LeafValid", leaf.Complete, leaf.Reason)
	}

	// Boundary: empty registered parent at each hierarchy level is
	// valid-but-incomplete (F14-REQ-12, F14-REQ-13).
	emptyParents := []validation.CompletenessValue{
		{
			Kind:             resources.KindProvider,
			UID:              "empty-provider",
			ProviderScopeUID: "empty-provider",
			Valid:            true,
		},
		{
			Kind:             resources.KindProviderLocation,
			UID:              "empty-location",
			ProviderScopeUID: "empty-provider",
			ParentUID:        "empty-provider",
			Valid:            true,
		},
		{
			Kind:             resources.KindProviderDatacenter,
			UID:              "empty-datacenter",
			ProviderScopeUID: "empty-provider",
			ParentUID:        "empty-location",
			Valid:            true,
		},
		{
			Kind:             resources.KindDatacenterFailureDomain,
			UID:              "empty-failure-domain",
			ProviderScopeUID: "empty-provider",
			ParentUID:        "empty-datacenter",
			Valid:            true,
		},
	}
	for _, parent := range emptyParents {
		ev := validation.EvaluateTopologyCompleteness(parent, []validation.CompletenessValue{parent})
		if ev.Complete || ev.Reason != validation.ReasonPathIncomplete {
			t.Fatalf("empty %s: complete=%v reason=%q want incomplete", parent.Kind, ev.Complete, ev.Reason)
		}
		if !parent.Valid {
			t.Fatalf("empty %s must remain Valid (registered ≠ complete)", parent.Kind)
		}
		if ev.HasChildren {
			t.Fatalf("empty %s must not report children", parent.Kind)
		}
	}
}

func TestFeature0014DistinctSameTechnologyStacks(t *testing.T) {
	t.Parallel()

	root := moduleRoot(t)
	var stackA resources.InfrastructureStack
	mustDecodeFeature0014Fixture(t, root, fixtureInfrastructureStack, &stackA)

	// Boundary: two identical-technology stacks retain distinct UIDs and parents
	// (F14-REQ-15, F14-REQ-16).
	stackB := stackA
	stackB.Metadata.Name = "infrastructure-stack-cloudstack-b"
	stackB.Metadata.UID = "a7000000000000000000000000000001"
	stackB.Spec.DatacenterFailureDomainRef = resources.DatacenterFailureDomainRef{
		TypedRef: apiref.TypedRef{
			APIVersion: resources.FabricAPIVersion,
			Kind:       resources.KindDatacenterFailureDomain,
			Name:       "datacenter-failure-domain-2",
			UID:        "a8000000000000000000000000000001",
		},
	}
	// Shared descriptive technology is not identity.
	if stackA.Spec.Technology != stackB.Spec.Technology {
		t.Fatalf("technology mismatch: %q vs %q", stackA.Spec.Technology, stackB.Spec.Technology)
	}
	if stackA.Metadata.UID == stackB.Metadata.UID {
		t.Fatal("same-technology stacks must retain distinct UIDs")
	}
	if stackA.Spec.DatacenterFailureDomainRef.UID == stackB.Spec.DatacenterFailureDomainRef.UID {
		t.Fatal("same-technology stacks must retain distinct failure-domain parents")
	}

	va := validation.CompletenessValueFromInfrastructureStack(stackA)
	vb := validation.CompletenessValueFromInfrastructureStack(stackB)
	if va.UID == vb.UID || va.ParentUID == vb.ParentUID {
		t.Fatalf("projected stacks must remain distinct: a=%+v b=%+v", va, vb)
	}
}

func TestFeature0014LeafFirstDeletionSequenceAccepted(t *testing.T) {
	t.Parallel()

	root := moduleRoot(t)
	set := mustLoadFeature0014CompletenessSet(t, root)

	order := []string{
		resources.KindInfrastructureStack,
		resources.KindDatacenterFailureDomain,
		resources.KindProviderDatacenter,
		resources.KindProviderLocation,
		resources.KindProvider,
	}

	remaining := append([]validation.CompletenessValue(nil), set...)
	for _, kind := range order {
		subject := findFeature0014ByKind(t, remaining, kind)
		if validation.HasImmediateTopologyChildren(subject, remaining) {
			t.Fatalf("leaf-first delete blocked for %s while children remain", kind)
		}
		remaining = removeFeature0014ByUID(remaining, subject.UID)
	}
	if len(remaining) != 0 {
		t.Fatalf("leaf-first deletion must empty the set; remaining=%d", len(remaining))
	}
}

func TestFeature0014GeoDescriptiveOnlyNoResidencyInference(t *testing.T) {
	t.Parallel()

	root := moduleRoot(t)
	raw := mustReadFeature0014Fixture(t, root, fixtureProviderLocation)

	lower := strings.ToLower(string(raw))
	for _, banned := range []string{
		`"residency"`, `"compliance"`, `"jurisdiction"`, `"policyoutcome"`,
		`"placement"`, `"eligible"`, `"sovereigntyproof"`,
	} {
		if strings.Contains(lower, banned) {
			t.Fatalf("provider-location fixture must not carry residency/compliance inference field %s", banned)
		}
	}

	var loc resources.ProviderLocation
	mustDecodeFeature0014Fixture(t, root, fixtureProviderLocation, &loc)
	if loc.Spec.Geo == nil {
		t.Fatal("fixture must carry descriptive geo")
	}
	if loc.Spec.Geo.CountryCode != "IN" || loc.Spec.Geo.SubdivisionCode != "IN-MH" {
		t.Fatalf("geo = %+v, want IN / IN-MH", loc.Spec.Geo)
	}
	// Descriptive topology only: presence of geo does not imply residency,
	// compliance, or placement eligibility (F14-REQ-27).
}

func TestFeature0014MultiOwnerIsolationIndependentScopeUIDs(t *testing.T) {
	t.Parallel()

	// Isolation: same external operator name under two owner Organizations
	// yields independent Provider scope UIDs (F14-REQ-30). Authorization never
	// keys on operator name.
	const operatorName = "provider-operator-a"
	ownerA := fixtureOwnerOrganizationUID
	ownerB := "d1000000000000000000000000000001"
	providerUnderA := resources.Provider{
		TypeMeta: apimeta.TypeMeta{APIVersion: resources.FabricAPIVersion, Kind: resources.KindProvider},
		Metadata: apimeta.ObjectMeta{
			Name: operatorName,
			UID:  fixtureProviderUID,
			ScopeRef: &apimeta.ScopeRef{TypedRef: apiref.TypedRef{
				APIVersion: "core.sovrunn.io/v1alpha1",
				Kind:       string(apimeta.ScopeOrganization),
				Name:       "owner-organization-a",
				UID:        ownerA,
			}},
		},
	}
	providerUnderB := resources.Provider{
		TypeMeta: apimeta.TypeMeta{APIVersion: resources.FabricAPIVersion, Kind: resources.KindProvider},
		Metadata: apimeta.ObjectMeta{
			Name: operatorName,
			UID:  "d2000000000000000000000000000001",
			ScopeRef: &apimeta.ScopeRef{TypedRef: apiref.TypedRef{
				APIVersion: "core.sovrunn.io/v1alpha1",
				Kind:       string(apimeta.ScopeOrganization),
				Name:       "owner-organization-b",
				UID:        ownerB,
			}},
		},
	}

	tvA := validation.TopologyValueFromProvider(providerUnderA)
	tvB := validation.TopologyValueFromProvider(providerUnderB)
	if tvA.ProviderScopeUID == tvB.ProviderScopeUID {
		t.Fatalf("same-operator/different-owner must yield independent scope UIDs: both %q", tvA.ProviderScopeUID)
	}
	if providerUnderA.Metadata.Name != providerUnderB.Metadata.Name {
		t.Fatal("isolation fixture uses the same operator name across owners")
	}
	if providerUnderA.Metadata.ScopeRef.UID == providerUnderB.Metadata.ScopeRef.UID {
		t.Fatal("owner Organization UIDs must differ")
	}

	// Hierarchy agreement never crosses Provider scope UIDs.
	locUnderA := validation.TopologyValue{
		Kind:             resources.KindProviderLocation,
		UID:              fixtureProviderLocationUID,
		ProviderScopeUID: tvA.ProviderScopeUID,
	}
	if !validation.EvaluateScopeAndHierarchy(tvA, locUnderA) {
		t.Fatal("same-owner path must agree")
	}
	if validation.EvaluateScopeAndHierarchy(tvB, locUnderA) {
		t.Fatal("cross-owner Provider scope must not agree")
	}
}

func mustFeature0014Structural(t *testing.T, root string) *StructuralValidator {
	t.Helper()
	reg, err := NewRepositorySchemaRegistry(filepath.Join(root, CanonicalSchemasDir))
	if err != nil {
		t.Fatalf("NewRepositorySchemaRegistry: %v", err)
	}
	resolver, err := NewLocalRefResolver(reg, DefaultMaxRefDepth)
	if err != nil {
		t.Fatalf("NewLocalRefResolver: %v", err)
	}
	cfg, err := NewStructuralValidatorConfig(reg, resolver)
	if err != nil {
		t.Fatalf("NewStructuralValidatorConfig: %v", err)
	}
	v, err := NewStructuralValidator(cfg)
	if err != nil {
		t.Fatalf("NewStructuralValidator: %v", err)
	}
	return v
}

func mustReadFeature0014Fixture(t *testing.T, root, name string) []byte {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join(root, ConformanceFixturesDir, name)) // #nosec G304 -- trusted repo-local fixture
	if err != nil {
		t.Fatalf("read fixture %s: %v", name, err)
	}
	if len(raw) == 0 {
		t.Fatalf("fixture %s is empty", name)
	}
	return raw
}

func mustDecodeFeature0014Fixture(t *testing.T, root, name string, dst any) {
	t.Helper()
	raw := mustReadFeature0014Fixture(t, root, name)
	lim := apivalid.DefaultLimits()
	pol := apivalid.PolicyFor(apivalid.ModeReadRepresentation)
	if prob := apivalid.DecodeJSON(raw, lim, pol, dst); prob != nil {
		t.Fatalf("DecodeJSON %s: code=%s detail=%s", name, prob.Code, prob.Detail)
	}
}

func mustStripSystemOwnedForCreate(t *testing.T, raw []byte) []byte {
	t.Helper()
	var top map[string]json.RawMessage
	if err := json.Unmarshal(raw, &top); err != nil {
		t.Fatalf("unmarshal fixture: %v", err)
	}
	delete(top, "status")

	metaRaw, ok := top["metadata"]
	if !ok {
		t.Fatal("fixture missing metadata")
	}
	var meta map[string]json.RawMessage
	if err := json.Unmarshal(metaRaw, &meta); err != nil {
		t.Fatalf("unmarshal metadata: %v", err)
	}
	for _, key := range []string{"uid", "generation", "resourceVersion", "createdAt", "updatedAt"} {
		delete(meta, key)
	}
	metaOut, err := json.Marshal(meta)
	if err != nil {
		t.Fatalf("marshal metadata: %v", err)
	}
	top["metadata"] = metaOut

	out, err := json.Marshal(top)
	if err != nil {
		t.Fatalf("marshal create-shaped fixture: %v", err)
	}
	return out
}

func mustLoadFeature0014CompletenessSet(t *testing.T, root string) []validation.CompletenessValue {
	t.Helper()

	var provider resources.Provider
	mustDecodeFeature0014Fixture(t, root, fixtureProvider, &provider)
	var loc resources.ProviderLocation
	mustDecodeFeature0014Fixture(t, root, fixtureProviderLocation, &loc)
	var dc resources.ProviderDatacenter
	mustDecodeFeature0014Fixture(t, root, fixtureProviderDatacenter, &dc)
	var fd resources.DatacenterFailureDomain
	mustDecodeFeature0014Fixture(t, root, fixtureDatacenterFailureDomain, &fd)
	var stack resources.InfrastructureStack
	mustDecodeFeature0014Fixture(t, root, fixtureInfrastructureStack, &stack)

	set := []validation.CompletenessValue{
		validation.CompletenessValueFromProvider(provider),
		validation.CompletenessValueFromProviderLocation(loc),
		validation.CompletenessValueFromProviderDatacenter(dc),
		validation.CompletenessValueFromDatacenterFailureDomain(fd),
		validation.CompletenessValueFromInfrastructureStack(stack),
	}

	if findFeature0014ByKind(t, set, resources.KindProvider).UID != fixtureProviderUID {
		t.Fatalf("provider uid mismatch")
	}
	if findFeature0014ByKind(t, set, resources.KindProviderLocation).UID != fixtureProviderLocationUID {
		t.Fatalf("location uid mismatch")
	}
	if findFeature0014ByKind(t, set, resources.KindProviderDatacenter).UID != fixtureProviderDatacenterUID {
		t.Fatalf("datacenter uid mismatch")
	}
	if findFeature0014ByKind(t, set, resources.KindDatacenterFailureDomain).UID != fixtureDatacenterFailureDomainUID {
		t.Fatalf("failure-domain uid mismatch")
	}
	if findFeature0014ByKind(t, set, resources.KindInfrastructureStack).UID != fixtureInfrastructureStackUID {
		t.Fatalf("stack uid mismatch")
	}

	// Sanity: Valid conditions projected from fixtures must be True.
	for _, v := range set {
		if !v.Valid {
			t.Fatalf("%s fixture Valid projection is false (check generation/conditions)", v.Kind)
		}
	}
	return set
}

func findFeature0014ByKind(t *testing.T, set []validation.CompletenessValue, kind string) validation.CompletenessValue {
	t.Helper()
	for _, v := range set {
		if v.Kind == kind {
			return v
		}
	}
	t.Fatalf("kind %s not found in completeness set", kind)
	return validation.CompletenessValue{}
}

func removeFeature0014ByUID(set []validation.CompletenessValue, uid string) []validation.CompletenessValue {
	out := make([]validation.CompletenessValue, 0, len(set))
	for _, v := range set {
		if v.UID == uid {
			continue
		}
		out = append(out, v)
	}
	return out
}
